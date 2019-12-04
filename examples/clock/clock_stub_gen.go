package clock

import (
	"bytes"
	"context"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	net "github.com/lugu/qiloop/bus/net"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
)

// TimestampImplementor interface of the service implementation
type TimestampImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals and properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation bus.Activation, helper TimestampSignalHelper) error
	OnTerminate()
	Nanoseconds() (int64, error)
}

// TimestampSignalHelper provided to Timestamp a companion object
type TimestampSignalHelper interface{}

// stubTimestamp implements server.Actor.
type stubTimestamp struct {
	impl      TimestampImplementor
	session   bus.Session
	service   bus.Service
	serviceID uint32
	signal    bus.SignalHandler
}

// TimestampObject returns an object using TimestampImplementor
func TimestampObject(impl TimestampImplementor) bus.Actor {
	var stb stubTimestamp
	stb.impl = impl
	obj := bus.NewBasicObject(&stb, stb.metaObject(), stb.onPropertyChange)
	stb.signal = obj
	return obj
}

// CreateTimestamp registers a new object to a service
// and returns a proxy to the newly created object
func CreateTimestamp(session bus.Session, service bus.Service, impl TimestampImplementor) (TimestampProxy, error) {
	obj := TimestampObject(impl)
	objectID, err := service.Add(obj)
	if err != nil {
		return nil, err
	}
	stb := &stubTimestamp{}
	meta := object.FullMetaObject(stb.metaObject())
	client := bus.DirectClient(obj)
	proxy := bus.NewProxy(client, meta, service.ServiceID(), objectID)
	return MakeTimestamp(session, proxy), nil
}
func (p *stubTimestamp) Activate(activation bus.Activation) error {
	p.session = activation.Session
	p.service = activation.Service
	p.serviceID = activation.ServiceID
	return p.impl.Activate(activation, p)
}
func (p *stubTimestamp) OnTerminate() {
	p.impl.OnTerminate()
}
func (p *stubTimestamp) Receive(msg *net.Message, from bus.Channel) error {
	// action dispatch
	switch msg.Header.Action {
	case 100:
		return p.Nanoseconds(msg, from)
	default:
		return from.SendError(msg, bus.ErrActionNotFound)
	}
}
func (p *stubTimestamp) onPropertyChange(name string, data []byte) error {
	switch name {
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}
func (p *stubTimestamp) Nanoseconds(msg *net.Message, c bus.Channel) error {
	ret, callErr := p.impl.Nanoseconds()

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := basic.WriteInt64(ret, &out)
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubTimestamp) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "Timestamp",
		Methods: map[uint32]object.MetaMethod{100: {
			Name:                "nanoseconds",
			ParametersSignature: "()",
			ReturnSignature:     "l",
			Uid:                 100,
		}},
		Properties: map[uint32]object.MetaProperty{},
		Signals:    map[uint32]object.MetaSignal{},
	}
}

// TimestampProxy represents a proxy object to the service
type TimestampProxy interface {
	Nanoseconds() (int64, error)
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
	WithContext(ctx context.Context) TimestampProxy
}

// proxyTimestamp implements TimestampProxy
type proxyTimestamp struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeTimestamp returns a specialized proxy.
func MakeTimestamp(sess bus.Session, proxy bus.Proxy) TimestampProxy {
	return &proxyTimestamp{bus.MakeObject(proxy), sess}
}

// Timestamp returns a proxy to a remote service
func Timestamp(session bus.Session) (TimestampProxy, error) {
	proxy, err := session.Proxy("Timestamp", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeTimestamp(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyTimestamp) WithContext(ctx context.Context) TimestampProxy {
	return MakeTimestamp(p.session, bus.WithContext(p.Proxy(), ctx))
}

// Nanoseconds calls the remote procedure
func (p *proxyTimestamp) Nanoseconds() (int64, error) {
	var err error
	var ret int64
	var buf bytes.Buffer
	response, err := p.Proxy().Call("nanoseconds", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call nanoseconds failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadInt64(resp)
	if err != nil {
		return ret, fmt.Errorf("parse nanoseconds response: %s", err)
	}
	return ret, nil
}
