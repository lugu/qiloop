package generic

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
	"math/rand"
)

type signalUser struct {
	signalID  uint32
	messageID uint32
	context   *server.Context
	clientID  uint64
}

// BasicObject implements the ServerObject interface. It handles the generic
// method and signal. Services implementation embedded a BasicObject
// and fill it with the extra actions they wish to handle using the
// Wrap method. See type/object.Object for a list of the default
// methods.
//
// TODO: make BasicObject the strict minimum. Build GenericObject on
// top of it.
type BasicObject struct {
	meta      object.MetaObject
	signals   []signalUser
	Wrapper   server.Wrapper
	serviceID uint32
	objectID  uint32
	terminate server.Terminator
}

// NewObject construct a BasicObject from a MetaObject.
func NewObject(meta object.MetaObject) *BasicObject {
	var obj BasicObject
	obj.meta = object.FullMetaObject(meta)
	obj.signals = make([]signalUser, 0)
	obj.Wrapper = make(map[uint32]server.ActionWrapper)
	obj.Wrap(uint32(0x2), obj.wrapMetaObject)
	obj.Wrap(uint32(0x3), obj.wrapTerminate)
	// obj.Wrapper[uint32(0x5)] = obj.Property
	// obj.Wrapper[uint32(0x6)] = obj.SetProperty
	// obj.Wrapper[uint32(0x7)] = obj.Properties
	// obj.Wrapper[uint32(0x8)] = obj.RegisterEventWithSignature
	return &obj
}

// Wrap let a BasicObject owner extend it with custom actions.
func (o *BasicObject) Wrap(id uint32, fn server.ActionWrapper) {
	o.Wrapper[id] = fn
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

func (o *BasicObject) wrapMetaObject(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, err
	}
	if objectID != o.objectID {
		return nil, fmt.Errorf("invalid object id: %d instead of %d",
			objectID, o.objectID)
	}
	var out bytes.Buffer
	err = object.WriteMetaObject(o.meta, &out)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
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

	a, ok := o.Wrapper[msg.Header.Action]
	if !ok {
		return util.ReplyError(from.EndPoint, msg,
			server.ErrActionNotFound)
	}
	response, err := a(msg.Payload)
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

// OnTerminate is called when the object is terminated.
func (o *BasicObject) wrapTerminate(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("Terminate failed: %s", err)
	}
	if objectID != o.objectID {
		return nil, fmt.Errorf("cannot terminate %d, only %d",
			objectID, o.objectID)
	}
	o.terminate()
	return make([]byte, 0), nil
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
	o.terminate = activation.Terminate
	return nil
}
