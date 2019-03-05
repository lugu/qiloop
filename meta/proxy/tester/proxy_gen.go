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

// Dummy is the abstract interface of the service
type Dummy interface {
	// Hello calls the remote procedure
	Hello() (BombProxy, error)
	// Attack calls the remote procedure
	Attack(b BombProxy) error
	// SubscribePing subscribe to a remote signal
	SubscribePing() (unsubscribe func(), updates chan string, err error)
	// GetStatus returns the property value
	GetStatus() (map[string]int32, error)
	// SetStatus sets the property value
	SetStatus(map[string]int32) error
	// SubscribeStatus regusters to a property
	SubscribeStatus() (unsubscribe func(), updates chan map[string]int32, err error)
	// GetCoordinate returns the property value
	GetCoordinate() (Coordinate, error)
	// SetCoordinate sets the property value
	SetCoordinate(Coordinate) error
	// SubscribeCoordinate regusters to a property
	SubscribeCoordinate() (unsubscribe func(), updates chan Coordinate, err error)
}

// Dummy represents a proxy object to the service
type DummyProxy interface {
	object.Object
	bus.Proxy
	Dummy
}

// proxyDummy implements DummyProxy
type proxyDummy struct {
	object1.ObjectProxy
	session bus.Session
}

// MakeDummy constructs DummyProxy
func MakeDummy(sess bus.Session, proxy bus.Proxy) DummyProxy {
	return &proxyDummy{object1.MakeObject(proxy), sess}
}

// NewDummy constructs DummyProxy
func NewDummy(sess bus.Session, obj uint32) (DummyProxy, error) {
	proxy, err := sess.Proxy("Dummy", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return MakeDummy(sess, proxy), nil
}

// Dummy retruns a proxy to a remote service
func (s Constructor) Dummy() (DummyProxy, error) {
	return NewDummy(s.session, 1)
}

// Hello calls the remote procedure
func (p *proxyDummy) Hello() (BombProxy, error) {
	var err error
	var ret BombProxy
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("hello", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call hello failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (BombProxy, error) {
		ref, err := object.ReadObjectReference(buf)
		if err != nil {
			return nil, fmt.Errorf("failed to get meta: %s", err)
		}
		obj, err := p.session.Object(ref)
		if err != nil {
			return nil, fmt.Errorf("failed to get proxy: %s", err)
		}
		return &proxyBomb{object1.MakeObject(obj), p.session}, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse hello response: %s", err)
	}
	return ret, nil
}

// Attack calls the remote procedure
func (p *proxyDummy) Attack(b BombProxy) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = func() error {
		meta, err := b.MetaObject(b.ObjectID())
		if err != nil {
			return fmt.Errorf("failed to get meta: %s", err)
		}
		ref := object.ObjectReference{
			true,
			meta,
			0,
			b.ServiceID(),
			b.ObjectID(),
		}
		return object.WriteObjectReference(ref, buf)
	}(); err != nil {
		return fmt.Errorf("failed to serialize b: %s", err)
	}
	_, err = p.Call("attack", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call attack failed: %s", err)
	}
	return nil
}

// SubscribePing subscribe to a remote property
func (p *proxyDummy) SubscribePing() (func(), chan string, error) {
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
func (p *proxyDummy) GetStatus() (ret map[string]int32, err error) {
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
	sig := "{si}"
	if sig != s {
		return ret, fmt.Errorf("unexpected signature: %s instead of %s",
			s, sig)
	}
	ret, err = func() (m map[string]int32, err error) {
		size, err := basic.ReadUint32(&buf)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]int32, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(&buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := basic.ReadInt32(&buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}()
	return ret, err
}

// SetStatus updates the property value
func (p *proxyDummy) SetStatus(update map[string]int32) error {
	name := value.String("status")
	var buf bytes.Buffer
	err := func() error {
		err := basic.WriteUint32(uint32(len(update)), &buf)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range update {
			err = basic.WriteString(k, &buf)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = basic.WriteInt32(v, &buf)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}()
	if err != nil {
		return fmt.Errorf("marshall error: %s", err)
	}
	val := value.Opaque("{si}", buf.Bytes())
	return p.SetProperty(name, val)
}

// SubscribeStatus subscribe to a remote property
func (p *proxyDummy) SubscribeStatus() (func(), chan map[string]int32, error) {
	propertyID, err := p.PropertyID("status")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "status", err)
	}

	handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register event for %s: %s", "status", err)
	}
	ch := make(chan map[string]int32)
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
			e, err := func() (m map[string]int32, err error) {
				size, err := basic.ReadUint32(buf)
				if err != nil {
					return m, fmt.Errorf("failed to read map size: %s", err)
				}
				m = make(map[string]int32, size)
				for i := 0; i < int(size); i++ {
					k, err := basic.ReadString(buf)
					if err != nil {
						return m, fmt.Errorf("failed to read map key: %s", err)
					}
					v, err := basic.ReadInt32(buf)
					if err != nil {
						return m, fmt.Errorf("failed to read map value: %s", err)
					}
					m[k] = v
				}
				return m, nil
			}()
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
func (p *proxyDummy) GetCoordinate() (ret Coordinate, err error) {
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
func (p *proxyDummy) SetCoordinate(update Coordinate) error {
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
func (p *proxyDummy) SubscribeCoordinate() (func(), chan Coordinate, error) {
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

// Bomb is the abstract interface of the service
type Bomb interface {
	// SubscribeBoom subscribe to a remote signal
	SubscribeBoom() (unsubscribe func(), updates chan int32, err error)
}

// Bomb represents a proxy object to the service
type BombProxy interface {
	object.Object
	bus.Proxy
	Bomb
}

// proxyBomb implements BombProxy
type proxyBomb struct {
	object1.ObjectProxy
	session bus.Session
}

// MakeBomb constructs BombProxy
func MakeBomb(sess bus.Session, proxy bus.Proxy) BombProxy {
	return &proxyBomb{object1.MakeObject(proxy), sess}
}

// NewBomb constructs BombProxy
func NewBomb(sess bus.Session, obj uint32) (BombProxy, error) {
	proxy, err := sess.Proxy("Bomb", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return MakeBomb(sess, proxy), nil
}

// Bomb retruns a proxy to a remote service
func (s Constructor) Bomb() (BombProxy, error) {
	return NewBomb(s.session, 1)
}

// SubscribeBoom subscribe to a remote property
func (p *proxyBomb) SubscribeBoom() (func(), chan int32, error) {
	propertyID, err := p.SignalID("boom")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "boom", err)
	}

	handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register event for %s: %s", "boom", err)
	}
	ch := make(chan int32)
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
			e, err := basic.ReadInt32(buf)
			if err != nil {
				log.Printf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}
