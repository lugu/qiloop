package pong

import (
	"bytes"
	"context"
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

// CreatePingPong registers a new object to a service
// and returns a proxy to the newly created object
func CreatePingPong(session bus.Session, service bus.Service, impl PingPongImplementor) (PingPongProxy, error) {
	obj := PingPongObject(impl)
	objectID, err := service.Add(obj)
	if err != nil {
		return nil, err
	}
	stb := &stubPingPong{}
	meta := object.FullMetaObject(stb.metaObject())
	client := bus.DirectClient(obj)
	proxy := bus.NewProxy(client, meta, service.ServiceID(), objectID)
	return MakePingPong(session, proxy), nil
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
func (p *stubPingPong) Receive(msg *net.Message, from bus.Channel) error {
	// action dispatch
	switch msg.Header.Action {
	case 100:
		return p.Hello(msg, from)
	case 101:
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
func (p *stubPingPong) Hello(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	a, err := basic.ReadString(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read a: %s", err))
	}
	ret, callErr := p.impl.Hello(a)

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
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
func (p *stubPingPong) Ping(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	a, err := basic.ReadString(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read a: %s", err))
	}
	callErr := p.impl.Ping(a)

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	return c.SendReply(msg, out.Bytes())
}
func (p *stubPingPong) SignalPong(a string) error {
	var buf bytes.Buffer
	if err := basic.WriteString(a, &buf); err != nil {
		return fmt.Errorf("serialize a: %s", err)
	}
	err := p.signal.UpdateSignal(102, buf.Bytes())

	if err != nil {
		return fmt.Errorf("update SignalPong: %s", err)
	}
	return nil
}
func (p *stubPingPong) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "PingPong",
		Methods: map[uint32]object.MetaMethod{
			100: {
				Name:                "hello",
				ParametersSignature: "(s)",
				ReturnSignature:     "s",
				Uid:                 100,
			},
			101: {
				Name:                "ping",
				ParametersSignature: "(s)",
				ReturnSignature:     "v",
				Uid:                 101,
			},
		},
		Properties: map[uint32]object.MetaProperty{},
		Signals: map[uint32]object.MetaSignal{102: {
			Name:      "pong",
			Signature: "s",
			Uid:       102,
		}},
	}
}

// PingPongProxy represents a proxy object to the service
type PingPongProxy interface {
	Hello(a string) (string, error)
	Ping(a string) error
	SubscribePong() (unsubscribe func(), updates chan string, err error)
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
	WithContext(ctx context.Context) PingPongProxy
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
func PingPong(session bus.Session) (PingPongProxy, error) {
	proxy, err := session.Proxy("PingPong", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakePingPong(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyPingPong) WithContext(ctx context.Context) PingPongProxy {
	return MakePingPong(p.session, p.Proxy().WithContext(ctx))
}

// Hello calls the remote procedure
func (p *proxyPingPong) Hello(a string) (string, error) {
	var err error
	var ret string
	var buf bytes.Buffer
	if err = basic.WriteString(a, &buf); err != nil {
		return ret, fmt.Errorf("serialize a: %s", err)
	}
	response, err := p.Proxy().Call("hello", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call hello failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadString(resp)
	if err != nil {
		return ret, fmt.Errorf("parse hello response: %s", err)
	}
	return ret, nil
}

// Ping calls the remote procedure
func (p *proxyPingPong) Ping(a string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(a, &buf); err != nil {
		return fmt.Errorf("serialize a: %s", err)
	}
	_, err = p.Proxy().Call("ping", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call ping failed: %s", err)
	}
	return nil
}

// SubscribePong subscribe to a remote property
func (p *proxyPingPong) SubscribePong() (func(), chan string, error) {
	propertyID, err := p.Proxy().MetaObject().SignalID("pong")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "pong", err)
	}
	ch := make(chan string)
	cancel, chPay, err := p.Proxy().SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := basic.ReadString(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}
