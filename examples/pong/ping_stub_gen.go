package pong

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	net "github.com/lugu/qiloop/bus/net"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	"log"
)

// PingPongImplementor interface of the service implementation
type PingPongImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals and properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation bus.Activation, helper PingPongSignalHelper) error
	OnTerminate()
	Hello(a string) (string, error)
	Ping(a string) error
}

// PingPongSignalHelper provided to PingPong a companion object
type PingPongSignalHelper interface {
	SignalPong(a string) error
}

// stubPingPong implements server.Actor.
type stubPingPong struct {
	impl      PingPongImplementor
	session   bus.Session
	service   bus.Service
	serviceID uint32
	signal    bus.SignalHandler
}

// PingPongObject returns an object using PingPongImplementor
func PingPongObject(impl PingPongImplementor) bus.Actor {
	var stb stubPingPong
	stb.impl = impl
	obj := bus.NewBasicObject(&stb, stb.metaObject(), stb.onPropertyChange)
	stb.signal = obj
	return obj
}

// NewPingPong registers a new object to a service
// and returns a proxy to the newly created object
func (c Constructor) NewPingPong(service bus.Service, impl PingPongImplementor) (PingPongProxy, error) {
	obj := PingPongObject(impl)
	objectID, err := service.Add(obj)
	if err != nil {
		return nil, err
	}
	stb := &stubPingPong{}
	meta := object.FullMetaObject(stb.metaObject())
	client := bus.DirectClient(obj)
	proxy := bus.NewProxy(client, meta, service.ServiceID(), objectID)
	return MakePingPong(c.session, proxy), nil
}
func (p *stubPingPong) Activate(activation bus.Activation) error {
	p.session = activation.Session
	p.service = activation.Service
	p.serviceID = activation.ServiceID
	return p.impl.Activate(activation, p)
}
func (p *stubPingPong) OnTerminate() {
	p.impl.OnTerminate()
}
func (p *stubPingPong) Receive(msg *net.Message, from *bus.Channel) error {
	switch msg.Header.Action {
	case uint32(0x64):
		return p.Hello(msg, from)
	case uint32(0x65):
		return p.Ping(msg, from)
	default:
		return from.SendError(msg, bus.ErrActionNotFound)
	}
}
func (p *stubPingPong) onPropertyChange(name string, data []byte) error {
	switch name {
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}
func (p *stubPingPong) Hello(msg *net.Message, c *bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	a, err := basic.ReadString(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read a: %s", err))
	}
	ret, callErr := p.impl.Hello(a)
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := basic.WriteString(ret, &out)
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubPingPong) Ping(msg *net.Message, c *bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	a, err := basic.ReadString(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read a: %s", err))
	}
	callErr := p.impl.Ping(a)
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	return c.SendReply(msg, out.Bytes())
}
func (p *stubPingPong) SignalPong(a string) error {
	var buf bytes.Buffer
	if err := basic.WriteString(a, &buf); err != nil {
		return fmt.Errorf("failed to serialize a: %s", err)
	}
	err := p.signal.UpdateSignal(uint32(0x66), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalPong: %s", err)
	}
	return nil
}
func (p *stubPingPong) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "PingPong",
		Methods: map[uint32]object.MetaMethod{
			uint32(0x64): {
				Name:                "hello",
				ParametersSignature: "(s)",
				ReturnSignature:     "s",
				Uid:                 uint32(0x64),
			},
			uint32(0x65): {
				Name:                "ping",
				ParametersSignature: "(s)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x65),
			},
		},
		Properties: map[uint32]object.MetaProperty{},
		Signals: map[uint32]object.MetaSignal{uint32(0x66): {
			Name:      "pong",
			Signature: "s",
			Uid:       uint32(0x66),
		}},
	}
}

// Constructor gives access to remote services
type Constructor struct {
	session bus.Session
}

// Services gives access to the services constructor
func Services(s bus.Session) Constructor {
	return Constructor{session: s}
}

// PingPong is the abstract interface of the service
type PingPong interface {
	// Hello calls the remote procedure
	Hello(a string) (string, error)
	// Ping calls the remote procedure
	Ping(a string) error
	// SubscribePong subscribe to a remote signal
	SubscribePong() (unsubscribe func(), updates chan string, err error)
}

// PingPongProxy represents a proxy object to the service
type PingPongProxy interface {
	object.Object
	bus.Proxy
	PingPong
}

// proxyPingPong implements PingPongProxy
type proxyPingPong struct {
	bus.ObjectProxy
	session bus.Session
}

// MakePingPong returns a specialized proxy.
func MakePingPong(sess bus.Session, proxy bus.Proxy) PingPongProxy {
	return &proxyPingPong{bus.MakeObject(proxy), sess}
}

// PingPong returns a proxy to a remote service
func (s Constructor) PingPong() (PingPongProxy, error) {
	proxy, err := s.session.Proxy("PingPong", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return MakePingPong(s.session, proxy), nil
}

// Hello calls the remote procedure
func (p *proxyPingPong) Hello(a string) (string, error) {
	var err error
	var ret string
	var buf bytes.Buffer
	if err = basic.WriteString(a, &buf); err != nil {
		return ret, fmt.Errorf("failed to serialize a: %s", err)
	}
	response, err := p.Call("hello", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call hello failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadString(resp)
	if err != nil {
		return ret, fmt.Errorf("failed to parse hello response: %s", err)
	}
	return ret, nil
}

// Ping calls the remote procedure
func (p *proxyPingPong) Ping(a string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(a, &buf); err != nil {
		return fmt.Errorf("failed to serialize a: %s", err)
	}
	_, err = p.Call("ping", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call ping failed: %s", err)
	}
	return nil
}

// SubscribePong subscribe to a remote property
func (p *proxyPingPong) SubscribePong() (func(), chan string, error) {
	propertyID, err := p.SignalID("pong")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "pong", err)
	}

	handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register event for %s: %s", "pong", err)
	}
	ch := make(chan string)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				p.UnregisterEvent(p.ObjectID(), propertyID, handlerID)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := basic.ReadString(buf)
			if err != nil {
				log.Printf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}