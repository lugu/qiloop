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

// PingPong is a proxy object to the remote service
type PingPong interface {
	object.Object
	bus.Proxy
	// Hello calls the remote procedure
	Hello(P0 string) (string, error)
	// Ping calls the remote procedure
	Ping(P0 string) error
	// SignalPong subscribe to a remote signal
	SignalPong(cancel chan int) (chan struct {
		P0 string
	}, error)
}

// PingPongProxy implements PingPong
type PingPongProxy struct {
	object1.ObjectProxy
}

// NewPingPong constructs PingPong
func NewPingPong(ses bus.Session, obj uint32) (PingPong, error) {
	proxy, err := ses.Proxy("PingPong", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &PingPongProxy{object1.ObjectProxy{proxy}}, nil
}

// PingPong retruns a proxy to a remote service
func (s Constructor) PingPong() (PingPong, error) {
	return NewPingPong(s.session, 1)
}

// Hello calls the remote procedure
func (p *PingPongProxy) Hello(P0 string) (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
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
func (p *PingPongProxy) Ping(P0 string) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("ping", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call ping failed: %s", err)
	}
	return nil
}

// SignalPong subscribe to a remote signal
func (p *PingPongProxy) SignalPong(cancel chan int) (chan struct {
	P0 string
}, error) {
	signalID, err := p.SignalID("pong")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "pong", err)
	}

	handlerID := uint64(signalID)<<32 + 1 // FIXME: read it from proxy
	_, err = p.RegisterEvent(p.ObjectID(), signalID, handlerID)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "pong", err)
	}
	ch := make(chan struct {
		P0 string
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				p.UnregisterEvent(p.ObjectID(), signalID, handlerID)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 string
			}, err error) {
				s.P0, err = basic.ReadString(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				log.Printf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
