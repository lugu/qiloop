// Package tester contains a generated proxy
// File generated. DO NOT EDIT.
package tester

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	object1 "github.com/lugu/qiloop/bus/client/object"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	"io"
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

// Dummy is a proxy object to the remote service
type Dummy interface {
	object.Object
	bus.Proxy
	// Hello calls the remote procedure
	Hello() error
	// SubscribePing subscribe to a remote signal
	SubscribePing() (func(), chan string, error)
	// SubscribeStatus regusters to a property
	SubscribeStatus() (func(), chan string, error)
	// SubscribeCoordinate regusters to a property
	SubscribeCoordinate() (func(), chan Coordinate, error)
}

// DummyProxy implements Dummy
type DummyProxy struct {
	object1.ObjectProxy
}

// NewDummy constructs Dummy
func NewDummy(ses bus.Session, obj uint32) (Dummy, error) {
	proxy, err := ses.Proxy("Dummy", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &DummyProxy{object1.ObjectProxy{proxy}}, nil
}

// Dummy retruns a proxy to a remote service
func (s Constructor) Dummy() (Dummy, error) {
	return NewDummy(s.session, 1)
}

// Hello calls the remote procedure
func (p *DummyProxy) Hello() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("hello", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call hello failed: %s", err)
	}
	return nil
}

// SubscribePing subscribe to a remote property
func (p *DummyProxy) SubscribePing() (func(), chan string, error) {
	propertyID, err := p.SignalID("ping")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "ping", err)
	}

	handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register event for %s: %s", "ping", err)
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

// SubscribeStatus subscribe to a remote property
func (p *DummyProxy) SubscribeStatus() (func(), chan string, error) {
	propertyID, err := p.PropertyID("status")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "status", err)
	}

	handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register event for %s: %s", "status", err)
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

// SubscribeCoordinate subscribe to a remote property
func (p *DummyProxy) SubscribeCoordinate() (func(), chan Coordinate, error) {
	propertyID, err := p.PropertyID("coordinate")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "coordinate", err)
	}

	handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register event for %s: %s", "coordinate", err)
	}
	ch := make(chan Coordinate)
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
			e, err := ReadCoordinate(buf)
			if err != nil {
				log.Printf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// Coordinate is serializable
type Coordinate struct {
	X int32
	Y int32
}

// ReadCoordinate unmarshalls Coordinate
func ReadCoordinate(r io.Reader) (s Coordinate, err error) {
	if s.X, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("failed to read X field: " + err.Error())
	}
	if s.Y, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("failed to read Y field: " + err.Error())
	}
	return s, nil
}

// WriteCoordinate marshalls Coordinate
func WriteCoordinate(s Coordinate, w io.Writer) (err error) {
	if err := basic.WriteInt32(s.X, w); err != nil {
		return fmt.Errorf("failed to write X field: " + err.Error())
	}
	if err := basic.WriteInt32(s.Y, w); err != nil {
		return fmt.Errorf("failed to write Y field: " + err.Error())
	}
	return nil
}
