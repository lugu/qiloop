// Package pingpong contains a generated stub
// File generated. DO NOT EDIT.

package pingpong

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	object1 "github.com/lugu/qiloop/bus/client/object"
	net "github.com/lugu/qiloop/bus/net"
	server "github.com/lugu/qiloop/bus/server"
	generic "github.com/lugu/qiloop/bus/server/generic"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	"log"
	"strings"
)

// PingPongImplementor interface of the service implementation
type PingPongImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals an properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation server.Activation, helper PingPongSignalHelper) error
	OnTerminate()
	Hello(a string) (string, error)
	Ping(a string) error
}

// PingPongSignalHelper provided to PingPong a companion object
type PingPongSignalHelper interface {
	SignalPong(a string) error
}

// stubPingPong implements server.ServerObject.
type stubPingPong struct {
	obj     generic.Object
	impl    PingPongImplementor
	session bus.Session
}

// PingPongObject returns an object using PingPongImplementor
func PingPongObject(impl PingPongImplementor) server.ServerObject {
	var stb stubPingPong
	stb.impl = impl
	stb.obj = generic.NewObject(stb.metaObject(), stb.onPropertyChange)
	stb.obj.Wrap(uint32(0x64), stb.Hello)
	stb.obj.Wrap(uint32(0x65), stb.Ping)
	return &stb
}
func (p *stubPingPong) Activate(activation server.Activation) error {
	p.session = activation.Session
	p.obj.Activate(activation)
	return p.impl.Activate(activation, p)
}
func (p *stubPingPong) OnTerminate() {
	p.impl.OnTerminate()
	p.obj.OnTerminate()
}
func (p *stubPingPong) Receive(msg *net.Message, from *server.Context) error {
	return p.obj.Receive(msg, from)
}
func (p *stubPingPong) onPropertyChange(name string, data []byte) error {
	switch strings.Title(name) {
	default:
		return fmt.Errorf("unknown property %s", name)
	}
	return nil
}
func (p *stubPingPong) Hello(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	a, err := basic.ReadString(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read a: %s", err)
	}
	ret, callErr := p.impl.Hello(a)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteString(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubPingPong) Ping(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	a, err := basic.ReadString(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read a: %s", err)
	}
	callErr := p.impl.Ping(a)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubPingPong) SignalPong(a string) error {
	var buf bytes.Buffer
	if err := basic.WriteString(a, &buf); err != nil {
		return fmt.Errorf("failed to serialize a: %s", err)
	}
	err := p.obj.UpdateSignal(uint32(0x66), buf.Bytes())

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
		Signals: map[uint32]object.MetaSignal{uint32(0x66): {
			Name:      "pong",
			Signature: "(s)",
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

// PingPong represents a proxy object to the service
type PingPongProxy interface {
	object.Object
	bus.Proxy
	PingPong
}

// proxyPingPong implements PingPongProxy
type proxyPingPong struct {
	object1.ObjectProxy
	session bus.Session
}

// MakePingPong constructs PingPongProxy
func MakePingPong(sess bus.Session, proxy bus.Proxy) PingPongProxy {
	return &proxyPingPong{object1.MakeObject(proxy), sess}
}

// PingPong retruns a proxy to a remote service
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
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(a, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize a: %s", err)
	}
	response, err := p.Call("hello", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call hello failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse hello response: %s", err)
	}
	return ret, nil
}

// Ping calls the remote procedure
func (p *proxyPingPong) Ping(a string) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(a, buf); err != nil {
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
