package server

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/value"
	"io"
	"math/rand"
	"sync"
)

type signalUser struct {
	signalID  uint32
	messageID uint32
	clientID  uint64
	context   *Context
}

// BasicObject implements the ServerObject interface. It handles the
// generic methods and signals common to all objects. Services
// implementation embedded a BasicObject and fill it with the extra
// actions they wish to handle using the Wrap method. See
// type/object.Object for a list of the default methods.
type BasicObject struct {
	signals      map[uint64]signalUser
	signalsMutex sync.RWMutex
	wrapper      Wrapper
	serviceID    uint32
	objectID     uint32
	tracer       func(*net.Message)
}

type Object interface {
	ServerObject
	UpdateSignal(signal uint32, data []byte) error
	UpdateProperty(property uint32, signature string, data []byte) error
	Wrap(id uint32, fn ActionWrapper)
}

// NewBasicObject construct a BasicObject from a MetaObject.
func NewBasicObject() *BasicObject {
	return &BasicObject{
		signals: make(map[uint64]signalUser),
		wrapper: make(map[uint32]ActionWrapper),
		tracer:  nil,
	}
}

// Wrap let a BasicObject owner extend it with custom actions.
func (o *BasicObject) Wrap(id uint32, fn ActionWrapper) {
	o.wrapper[id] = fn
}

// addSignalUser register the context as a client of event signalID.
// TODO: check if the signalID is valid
func (o *BasicObject) addSignalUser(signalID, messageID uint32,
	from *Context) uint64 {

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
func (o *BasicObject) removeSignalUser(clientID uint64) error {
	o.signalsMutex.Lock()
	delete(o.signals, clientID)
	defer o.signalsMutex.Unlock()
	return nil
}

func (o *BasicObject) handleRegisterEvent(from *Context,
	msg *net.Message) error {

	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read object uid: %s", err)
		return o.replyError(from, msg, err)
	}
	if objectID != o.objectID {
		err := fmt.Errorf("wrong object id, expecting %d, got %d",
			msg.Header.Object, objectID)
		return o.replyError(from, msg, err)
	}
	signalID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read signal uid: %s", err)
		return o.replyError(from, msg, err)
	}
	_, err = basic.ReadUint64(buf)
	if err != nil {
		err = fmt.Errorf("cannot read client uid: %s", err)
		return o.replyError(from, msg, err)
	}
	messageID := msg.Header.ID
	clientID := o.addSignalUser(signalID, messageID, from)
	var out bytes.Buffer
	err = basic.WriteUint64(clientID, &out)
	if err != nil {
		err = fmt.Errorf("cannot write client uid: %s", err)
		return o.replyError(from, msg, err)
	}
	return o.reply(from, msg, out.Bytes())
}

func (o *BasicObject) handleUnregisterEvent(from *Context,
	msg *net.Message) error {

	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read object uid: %s", err)
		return o.replyError(from, msg, err)
	}
	if objectID != o.objectID {
		err := fmt.Errorf("wrong object id, expecting %d, got %d",
			msg.Header.Object, objectID)
		return o.replyError(from, msg, err)
	}
	_, err = basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read action uid: %s", err)
		return o.replyError(from, msg, err)
	}
	clientID, err := basic.ReadUint64(buf)
	if err != nil {
		err = fmt.Errorf("cannot read client uid: %s", err)
		return o.replyError(from, msg, err)
	}
	err = o.removeSignalUser(clientID)
	if err != nil {
		return o.replyError(from, msg, err)
	}
	var out bytes.Buffer
	return o.reply(from, msg, out.Bytes())
}

// UpdateSignal informs the registered clients of the new state.
func (o *BasicObject) UpdateSignal(id uint32, data []byte) error {
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
func (o *BasicObject) UpdateProperty(id uint32, sig string, data []byte) error {
	return o.UpdateSignal(id, data)
}

func (o *BasicObject) trace(msg *net.Message) {
	if o.tracer != nil {
		o.tracer(msg)
	}
}

func (o *BasicObject) replyEvent(client *signalUser, signal uint32,
	value []byte) error {

	hdr := o.newHeader(net.Event, signal, client.messageID)
	msg := net.NewMessage(hdr, value)
	o.trace(&msg)
	return client.context.EndPoint.Send(msg)
}

func (o *BasicObject) sendTerminate(client *signalUser, signal uint32) error {
	var buf bytes.Buffer
	val := value.String(ErrTerminate.Error())
	val.Write(&buf)
	hdr := o.newHeader(net.Error, client.signalID, client.messageID)
	msg := net.NewMessage(hdr, buf.Bytes())

	o.trace(&msg)
	return client.context.EndPoint.Send(msg)
}

func (o *BasicObject) replyError(from *Context, msg *net.Message,
	err error) error {

	o.trace(msg)
	return util.ReplyError(from.EndPoint, msg, err)
}

func (o *BasicObject) reply(from *Context, msg *net.Message,
	response []byte) error {

	hdr := o.newHeader(net.Reply, msg.Header.Action, msg.Header.ID)
	reply := net.NewMessage(hdr, response)
	o.trace(&reply)
	return from.EndPoint.Send(reply)
}

func (o *BasicObject) handleDefault(from *Context,
	msg *net.Message) error {

	fn, ok := o.wrapper[msg.Header.Action]
	if !ok {
		return o.replyError(from, msg, ErrActionNotFound)
	}
	response, err := fn(msg.Payload)
	if err != nil {
		return o.replyError(from, msg, err)
	}
	return o.reply(from, msg, response)
}

// Receive processes the incoming message and responds to the client.
// The returned error is not destinated to the client which have
// already be replied.
func (o *BasicObject) Receive(msg *net.Message, from *Context) error {
	o.trace(msg)

	if msg.Header.Type != net.Call {
		err := fmt.Errorf("unsupported message type: %d",
			msg.Header.Type)
		return o.replyError(from, msg, err)
	}
	switch msg.Header.Action {
	case 0x0:
		return o.handleRegisterEvent(from, msg)
	case 0x1:
		return o.handleUnregisterEvent(from, msg)
	default:
		return o.handleDefault(from, msg)
	}
}

func (o *BasicObject) OnTerminate() {
	// close all subscribers
	o.signalsMutex.Lock()
	for _, client := range o.signals {
		o.sendTerminate(&client, client.signalID)
		delete(o.signals, client.clientID)
	}
	o.signalsMutex.Unlock()
}

func (o *BasicObject) newHeader(typ uint8, action, id uint32) net.Header {
	return net.NewHeader(typ, o.serviceID, o.objectID, action, id)
}

// Activate informs the object when it becomes online. After this
// method returns, the object will start receiving incomming messages.
func (o *BasicObject) Activate(activation Activation) error {
	o.serviceID = activation.ServiceID
	o.objectID = activation.ObjectID
	return nil
}
