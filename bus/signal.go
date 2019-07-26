package bus

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"sync"

	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/value"
)

type signalUser struct {
	signalID  uint32
	messageID uint32
	clientID  uint64
	context   Channel
}

// signalHandler implements the Actor interface. It implements the
// registerEvent and unregisterEvent methods and provides a means to
// update a signal.
type signalHandler struct {
	signals      map[uint64]signalUser
	signalsMutex sync.RWMutex
	serviceID    uint32
	objectID     uint32
	tracer       func(*net.Message)
}

// SignalHandler is an helper object for service implementation: it
// is passed during the activation call to allow a service implementor
// to manipilate its signals and properties.
type SignalHandler interface {
	UpdateSignal(signal uint32, data []byte) error
	UpdateProperty(property uint32, signature string, data []byte) error
}

// BasicObject implements the common behavior to all objects (including
// the signal/properties subscriptions) in an abstract way.
type BasicObject interface {
	Actor
	SignalHandler
}

// newSignalHandler returns a helper to deal with signals.
func newSignalHandler() *signalHandler {
	return &signalHandler{
		signals: make(map[uint64]signalUser),
		tracer:  nil,
	}
}

// addSignalUser register the context as a client of event signalID.
// TODO: check if the signalID is valid
func (o *signalHandler) addSignalUser(signalID, messageID uint32,
	from Channel) uint64 {

	clientID := rand.Uint64()
	newUser := signalUser{
		signalID,
		messageID,
		clientID,
		from,
	}
	o.signalsMutex.Lock()
	_, ok := o.signals[clientID]
	if !ok {
		o.signals[clientID] = newUser
		o.signalsMutex.Unlock()
		return clientID
	}
	o.signalsMutex.Unlock()
	// pick another random number
	return o.addSignalUser(signalID, messageID, from)
}

// removeSignalUser unregister the given contex to events.
func (o *signalHandler) removeSignalUser(clientID uint64) error {
	o.signalsMutex.Lock()
	delete(o.signals, clientID)
	defer o.signalsMutex.Unlock()
	return nil
}

func (o *signalHandler) RegisterEvent(msg *net.Message, from Channel) error {

	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read object uid: %s", err)
		return from.SendError(msg, err)
	}
	// remote objects don't know their real object id.
	if o.objectID < (1<<31) && objectID != o.objectID {
		err := fmt.Errorf("wrong object id, expecting %d, got %d",
			msg.Header.Object, objectID)
		return from.SendError(msg, err)
	}
	signalID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read signal uid: %s", err)
		return from.SendError(msg, err)
	}
	_, err = basic.ReadUint64(buf)
	if err != nil {
		err = fmt.Errorf("cannot read client uid: %s", err)
		return from.SendError(msg, err)
	}
	messageID := msg.Header.ID
	clientID := o.addSignalUser(signalID, messageID, from)
	var out bytes.Buffer
	err = basic.WriteUint64(clientID, &out)
	if err != nil {
		err = fmt.Errorf("cannot write client uid: %s", err)
		return from.SendError(msg, err)
	}
	return from.SendReply(msg, out.Bytes())
}

func (o *signalHandler) UnregisterEvent(msg *net.Message, from Channel) error {

	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read object uid: %s", err)
		return from.SendError(msg, err)
	}
	// remote objects don't know their real object id.
	if o.objectID < (1<<31) && objectID != o.objectID {
		err := fmt.Errorf("wrong object id, expecting %d, got %d",
			msg.Header.Object, objectID)
		return from.SendError(msg, err)
	}
	_, err = basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read action uid: %s", err)
		return from.SendError(msg, err)
	}
	clientID, err := basic.ReadUint64(buf)
	if err != nil {
		err = fmt.Errorf("cannot read client uid: %s", err)
		return from.SendError(msg, err)
	}
	err = o.removeSignalUser(clientID)
	if err != nil {
		return from.SendError(msg, err)
	}
	var out bytes.Buffer
	return from.SendReply(msg, out.Bytes())
}

// UpdateSignal informs the registered clients of the new state.
func (o *signalHandler) UpdateSignal(id uint32, data []byte) error {
	var ret error
	signals := make([]signalUser, 0)

	o.signalsMutex.RLock()
	for _, client := range o.signals {
		if client.signalID == id {
			signals = append(signals, client)
		}
	}
	o.signalsMutex.RUnlock()

	for _, client := range signals {
		err := o.replyEvent(&client, id, data)
		if err == io.EOF {
			o.removeSignalUser(client.clientID)
		} else if err != nil {
			ret = err
		}
	}
	return ret
}

// UpdateProperty informs the registered clients of the property change
func (o *signalHandler) UpdateProperty(id uint32, sig string, data []byte) error {
	return o.UpdateSignal(id, data)
}

func (o *signalHandler) trace(msg *net.Message) {
	if o.tracer != nil {
		o.tracer(msg)
	}
}

func (o *signalHandler) replyEvent(client *signalUser, signal uint32,
	value []byte) error {

	hdr := o.newHeader(net.Event, signal, client.messageID)
	msg := net.NewMessage(hdr, value)
	o.trace(&msg)
	return client.context.Send(&msg)
}

func (o *signalHandler) sendTerminate(client *signalUser, signal uint32) error {
	var buf bytes.Buffer
	val := value.String(ErrTerminate.Error())
	val.Write(&buf)
	hdr := o.newHeader(net.Error, client.signalID, client.messageID)
	msg := net.NewMessage(hdr, buf.Bytes())

	o.trace(&msg)
	return client.context.Send(&msg)
}

func (o *signalHandler) OnTerminate() {
	// close all subscribers
	o.signalsMutex.Lock()
	for _, client := range o.signals {
		o.sendTerminate(&client, client.signalID)
		delete(o.signals, client.clientID)
	}
	o.signalsMutex.Unlock()
}

func (o *signalHandler) newHeader(typ uint8, action, id uint32) net.Header {
	return net.NewHeader(typ, o.serviceID, o.objectID, action, id)
}

// Activate informs the object when it becomes online. After this
// method returns, the object will start receiving incomming messages.
func (o *signalHandler) Activate(activation Activation) error {
	o.serviceID = activation.ServiceID
	o.objectID = activation.ObjectID
	return nil
}
