// Package proxy contains a generated proxy
// File generated. DO NOT EDIT.
package proxy

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	object1 "github.com/lugu/qiloop/bus/client/object"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	"log"
)

// Constructor gives access to remote services
type Constructor struct {
	session bus.Session
}

// Services gives access to the services constructor
func Services(s bus.Session) Constructor {
	return Constructor{session: s}
}

// PingPongObject is the abstract interface of the service
type PingPong interface {
	// Hello calls the remote procedure
	Hello(a string) (string, error)
	// Ping calls the remote procedure
	Ping(a string) error
	// SubscribePong subscribe to a remote signal
	SubscribePong() (unsubscribe func(), updates chan string, err error)
}

// PingPong represents a proxy object to the service
type PingPongObject interface {
	object.Object
	bus.Proxy
	PingPong
}

// proxyPingPong implements PingPongObject
type proxyPingPong struct {
	object1.ObjectObject
	session bus.Session
}

// NewPingPong constructs PingPongObject
func NewPingPong(ses bus.Session, obj uint32) (PingPongObject, error) {
	proxy, err := ses.Proxy("PingPong", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &proxyPingPong{object1.MakeObject(proxy), ses}, nil
}

// PingPong retruns a proxy to a remote service
func (s Constructor) PingPong() (PingPongObject, error) {
	return NewPingPong(s.session, 1)
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
