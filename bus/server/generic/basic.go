package generic

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/basic"
	"math/rand"
)

type signalUser struct {
	signalID  uint32
	messageID uint32
	context   *server.Context
	clientID  uint64
}

// BasicObject implements the ServerObject interface. It handles the
// generic methods and signals common to all objects. Services
// implementation embedded a BasicObject and fill it with the extra
// actions they wish to handle using the Wrap method. See
// type/object.Object for a list of the default methods.
type BasicObject struct {
	signals   []signalUser
	wrapper   server.Wrapper
	serviceID uint32
	objectID  uint32
}

type Object interface {
	server.ServerObject
	UpdateSignal(signal uint32, value []byte) error
	Wrap(id uint32, fn server.ActionWrapper)
}

// NewBasicObject construct a BasicObject from a MetaObject.
func NewBasicObject() *BasicObject {
	return &BasicObject{
		signals: make([]signalUser, 0),
		wrapper: make(map[uint32]server.ActionWrapper),
	}
}

// Wrap let a BasicObject owner extend it with custom actions.
func (o *BasicObject) Wrap(id uint32, fn server.ActionWrapper) {
	o.wrapper[id] = fn
}

func (o *BasicObject) addSignalUser(signalID, messageID uint32,
	from *server.Context) uint64 {

	clientID := rand.Uint64()
	newUser := signalUser{
		signalID,
		messageID,
		from,
		clientID,
	}
	o.signals = append(o.signals, newUser)
	return clientID
}

func (o *BasicObject) removeSignalUser(id uint64) error {
	return nil
}

func (o *BasicObject) handleRegisterEvent(from *server.Context,
	msg *net.Message) error {

	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read object uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	if objectID != o.objectID {
		err := fmt.Errorf("wrong object id, expecting %d, got %d",
			msg.Header.Object, objectID)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	signalID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read signal uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	_, err = basic.ReadUint64(buf)
	if err != nil {
		err = fmt.Errorf("cannot read client uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	messageID := msg.Header.ID
	clientID := o.addSignalUser(signalID, messageID, from)
	var out bytes.Buffer
	err = basic.WriteUint64(clientID, &out)
	if err != nil {
		err = fmt.Errorf("cannot write client uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	return o.reply(from, msg, out.Bytes())
}

func (o *BasicObject) handleUnregisterEvent(from *server.Context,
	msg *net.Message) error {

	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read object uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	if objectID != o.objectID {
		err := fmt.Errorf("wrong object id, expecting %d, got %d",
			msg.Header.Object, objectID)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	_, err = basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read action uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	clientID, err := basic.ReadUint64(buf)
	if err != nil {
		err = fmt.Errorf("cannot read client uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	err = o.removeSignalUser(clientID)
	if err != nil {
		return util.ReplyError(from.EndPoint, msg, err)
	}
	var out bytes.Buffer
	return o.reply(from, msg, out.Bytes())
}

// UpdateSignal informs the registered clients of the new state.
func (o *BasicObject) UpdateSignal(signal uint32, value []byte) error {
	var ret error
	for _, client := range o.signals {
		if client.signalID == signal {
			hdr := o.newHeader(net.Event, signal, client.messageID)
			msg := net.NewMessage(hdr, value)
			// FIXME: catch writing to close connection
			err := client.context.EndPoint.Send(msg)
			if err != nil {
				ret = err
			}
		}
	}
	return ret
}

func (o *BasicObject) reply(from *server.Context, m *net.Message,
	response []byte) error {

	hdr := o.newHeader(net.Reply, m.Header.Action, m.Header.ID)
	reply := net.NewMessage(hdr, response)
	return from.EndPoint.Send(reply)
}

func (o *BasicObject) handleDefault(from *server.Context,
	msg *net.Message) error {

	fn, ok := o.wrapper[msg.Header.Action]
	if !ok {
		return util.ReplyError(from.EndPoint, msg,
			server.ErrActionNotFound)
	}
	response, err := fn(msg.Payload)
	if err != nil {
		return util.ReplyError(from.EndPoint, msg, err)
	}
	return o.reply(from, msg, response)
}

// Receive processes the incoming message and responds to the client.
// The returned error is not destinated to the client which have
// already be replied.
func (o *BasicObject) Receive(m *net.Message, from *server.Context) error {
	switch m.Header.Action {
	case 0x0:
		return o.handleRegisterEvent(from, m)
	case 0x1:
		return o.handleUnregisterEvent(from, m)
	default:
		return o.handleDefault(from, m)
	}
}

func (o *BasicObject) OnTerminate() {
}

func (o *BasicObject) newHeader(typ uint8, action, id uint32) net.Header {
	return net.NewHeader(typ, o.serviceID, o.objectID, action, id)
}

// Activate informs the object when it becomes online. After this
// method returns, the object will start receiving incomming messages.
func (o *BasicObject) Activate(activation server.Activation) error {
	o.serviceID = activation.ServiceID
	o.objectID = activation.ObjectID
	return nil
}
