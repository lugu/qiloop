package bus

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/value"
)

type signalUser struct {
	signalID  uint32
	messageID uint32
	userID    uint64
	context   Channel
	contextID int
}

// signalHandler implements the Actor interface. It implements the
// registerEvent and unregisterEvent methods and provides a means to
// update a signal.
type signalHandler struct {
	signals      []signalUser
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
		signals: make([]signalUser, 0, 10),
		tracer:  nil,
	}
}

// addSignalUser register the context as a user of event signalID.
// TODO: check if the signalID is valid
func (o *signalHandler) addSignalUser(userID uint64, signalID, messageID uint32,
	from Channel) error {

	newUser := signalUser{
		signalID:  signalID,
		messageID: messageID,
		userID:    userID,
		context:   from,
		contextID: 0,
	}

	e := from.EndPoint()
	f := func(hdr *net.Header) (bool, bool) {
		return false, true
	}
	q := make(chan<- *net.Message)
	cl := func(err error) {
		// unregister user on disconnection
		o.removeSignalUser(userID, from)
	}
	newUser.contextID = e.MakeHandler(f, q, cl)

	o.signalsMutex.Lock()

	for _, user := range o.signals {
		if user.userID == userID {
			o.signalsMutex.Unlock()
			user.context.EndPoint().RemoveHandler(user.contextID)
			return fmt.Errorf("user %d already exists", userID)
		}
	}
	o.signals = append(o.signals, newUser)
	o.signalsMutex.Unlock()
	return nil

}

// removeSignalUser unregister the given contex to events.
func (o *signalHandler) removeSignalUser(userID uint64, from Channel) error {
	o.signalsMutex.Lock()

	for i, user := range o.signals {
		if user.userID == userID {
			if from.EndPoint() == user.context.EndPoint() {
				o.signals[i] = o.signals[len(o.signals)-1]
				o.signals = o.signals[:len(o.signals)-1]
				o.signalsMutex.Unlock()
				user.context.EndPoint().RemoveHandler(user.contextID)
				return nil
			}
		}
	}
	o.signalsMutex.Unlock()
	return fmt.Errorf("unknown user id %d", userID)
}

func (o *signalHandler) RegisterEvent(msg *net.Message, from Channel) error {

	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read object uid: %s", err)
		return from.SendError(msg, err)
	}
	// remote objects don't know their real object id.
	if objectID != 0 && o.objectID < (1<<31) && objectID != o.objectID {
		err := fmt.Errorf("wrong object id, expecting %d, got %d",
			msg.Header.Object, objectID)
		return from.SendError(msg, err)
	}
	signalID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read signal uid: %s", err)
		return from.SendError(msg, err)
	}
	userID, err := basic.ReadUint64(buf)
	if err != nil {
		err = fmt.Errorf("cannot read user uid: %s", err)
		return from.SendError(msg, err)
	}
	messageID := msg.Header.ID
	err = o.addSignalUser(userID, signalID, messageID, from)
	if err != nil {
		err = fmt.Errorf("cannot register user uid %d: %s",
			userID, err)
		return from.SendError(msg, err)
	}
	var out bytes.Buffer
	err = basic.WriteUint64(userID, &out)
	if err != nil {
		err = fmt.Errorf("cannot write user uid: %s", err)
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
	if objectID != 0 && o.objectID < (1<<31) && objectID != o.objectID {
		err := fmt.Errorf("wrong object id, expecting %d, got %d",
			msg.Header.Object, objectID)
		return from.SendError(msg, err)
	}
	_, err = basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read action uid: %s", err)
		return from.SendError(msg, err)
	}
	userID, err := basic.ReadUint64(buf)
	if err != nil {
		err = fmt.Errorf("cannot read user uid: %s", err)
		return from.SendError(msg, err)
	}
	err = o.removeSignalUser(userID, from)
	if err != nil {
		return from.SendError(msg, err)
	}
	var out bytes.Buffer
	return from.SendReply(msg, out.Bytes())
}

// UpdateSignal informs the registered clients of the new state.
func (o *signalHandler) UpdateSignal(signalID uint32, data []byte) error {
	var ret error
	signals := make([]signalUser, 0)

	o.signalsMutex.RLock()
	for _, user := range o.signals {
		if user.signalID == signalID {
			signals = append(signals, user)
		}
	}
	o.signalsMutex.RUnlock()

	for _, user := range signals {
		err := o.replyEvent(&user, signalID, data)
		if err == io.EOF {
			err := o.removeSignalUser(user.userID, user.context)
			if err != nil && ret == nil {
				ret = err
			}
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

func (o *signalHandler) replyEvent(user *signalUser, signal uint32,
	value []byte) error {

	hdr := o.newHeader(net.Event, signal, user.messageID)
	msg := net.NewMessage(hdr, value)
	o.trace(&msg)
	return user.context.Send(&msg)
}

func (o *signalHandler) sendTerminate(user *signalUser, signal uint32) error {
	var buf bytes.Buffer
	val := value.String(ErrTerminate.Error())
	val.Write(&buf)
	hdr := o.newHeader(net.Error, user.signalID, user.messageID)
	msg := net.NewMessage(hdr, buf.Bytes())

	o.trace(&msg)
	return user.context.Send(&msg)
}

func (o *signalHandler) OnTerminate() {
	// close all subscribers
	o.signalsMutex.Lock()
	signals := o.signals
	o.signals = []signalUser{}
	o.signalsMutex.Unlock()
	for _, user := range signals {
		o.sendTerminate(&user, user.signalID)
		user.context.EndPoint().RemoveHandler(user.contextID)
	}
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
