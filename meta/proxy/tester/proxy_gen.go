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
	value "github.com/lugu/qiloop/type/value"
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
	SubscribePing() (unsubscribe func(), updates chan string, err error)
	// GetStatus returns the property value
	GetStatus() (string, error)
	// SetStatus sets the property value
	SetStatus(string) error
	// SubscribeStatus regusters to a property
	SubscribeStatus() (unsubscribe func(), updates chan string, err error)
	// GetCoordinate returns the property value
	GetCoordinate() (Coordinate, error)
	// SetCoordinate sets the property value
	SetCoordinate(Coordinate) error
	// SubscribeCoordinate regusters to a property
	SubscribeCoordinate() (unsubscribe func(), updates chan Coordinate, err error)
}

// DummyProxy implements Dummy
type DummyProxy struct {
	object1.ObjectProxy
	session bus.Session
}

// NewDummy constructs Dummy
func NewDummy(ses bus.Session, obj uint32) (Dummy, error) {
	proxy, err := ses.Proxy("Dummy", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &DummyProxy{object1.ObjectProxy{proxy}, ses}, nil
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

// GetStatus updates the property value
func (p *DummyProxy) GetStatus() (ret string, err error) {
	name := value.String("status")
	value, err := p.Property(name)
	if err != nil {
		return ret, fmt.Errorf("get property: %s", err)
	}
	var buf bytes.Buffer
	err = value.Write(&buf)
	if err != nil {
		return ret, fmt.Errorf("read response: %s", err)
	}
	s, err := basic.ReadString(&buf)
	if err != nil {
		return ret, fmt.Errorf("read signature: %s", err)
	}
	// check the signature
	sig := "s"
	if sig != s {
		return ret, fmt.Errorf("unexpected signature: %s instead of %s",
			s, sig)
	}
	ret, err = basic.ReadString(&buf)
	return ret, err
}

// SetStatus updates the property value
func (p *DummyProxy) SetStatus(update string) error {
	name := value.String("status")
	var buf bytes.Buffer
	err := basic.WriteString(update, &buf)
	if err != nil {
		return fmt.Errorf("marshall error: %s", err)
	}
	val := value.Opaque("s", buf.Bytes())
	return p.SetProperty(name, val)
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

// GetCoordinate updates the property value
func (p *DummyProxy) GetCoordinate() (ret Coordinate, err error) {
	name := value.String("coordinate")
	value, err := p.Property(name)
	if err != nil {
		return ret, fmt.Errorf("get property: %s", err)
	}
	var buf bytes.Buffer
	err = value.Write(&buf)
	if err != nil {
		return ret, fmt.Errorf("read response: %s", err)
	}
	s, err := basic.ReadString(&buf)
	if err != nil {
		return ret, fmt.Errorf("read signature: %s", err)
	}
	// check the signature
	sig := "(ii)<coordinate,x,y>"
	if sig != s {
		return ret, fmt.Errorf("unexpected signature: %s instead of %s",
			s, sig)
	}
	ret, err = ReadCoordinate(&buf)
	return ret, err
}

// SetCoordinate updates the property value
func (p *DummyProxy) SetCoordinate(update Coordinate) error {
	name := value.String("coordinate")
	var buf bytes.Buffer
	err := WriteCoordinate(update, &buf)
	if err != nil {
		return fmt.Errorf("marshall error: %s", err)
	}
	val := value.Opaque("(ii)<coordinate,x,y>", buf.Bytes())
	return p.SetProperty(name, val)
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
