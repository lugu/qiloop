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

// Spacecraft is the abstract interface of the service
type Spacecraft interface {
	// Shoot calls the remote procedure
	Shoot() (BombProxy, error)
	// Ammo calls the remote procedure
	Ammo(ammo BombProxy) error
}

// Spacecraft represents a proxy object to the service
type SpacecraftProxy interface {
	object.Object
	bus.Proxy
	Spacecraft
}

// proxySpacecraft implements SpacecraftProxy
type proxySpacecraft struct {
	object1.ObjectProxy
	session bus.Session
}

// MakeSpacecraft constructs SpacecraftProxy
func MakeSpacecraft(sess bus.Session, proxy bus.Proxy) SpacecraftProxy {
	return &proxySpacecraft{object1.MakeObject(proxy), sess}
}

// NewSpacecraft constructs SpacecraftProxy
func NewSpacecraft(sess bus.Session, obj uint32) (SpacecraftProxy, error) {
	proxy, err := sess.Proxy("Spacecraft", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return MakeSpacecraft(sess, proxy), nil
}

// Spacecraft retruns a proxy to a remote service
func (s Constructor) Spacecraft() (SpacecraftProxy, error) {
	return NewSpacecraft(s.session, 1)
}

// Shoot calls the remote procedure
func (p *proxySpacecraft) Shoot() (BombProxy, error) {
	var err error
	var ret BombProxy
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("shoot", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call shoot failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (BombProxy, error) {
		ref, err := object.ReadObjectReference(buf)
		if err != nil {
			return nil, fmt.Errorf("failed to get meta: %s", err)
		}
		proxy, err := p.session.Object(ref)
		if err != nil {
			return nil, fmt.Errorf("failed to get proxy: %s", err)
		}
		return MakeBomb(p.session, proxy), nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse shoot response: %s", err)
	}
	return ret, nil
}

// Ammo calls the remote procedure
func (p *proxySpacecraft) Ammo(ammo BombProxy) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = func() error {
		meta, err := ammo.MetaObject(ammo.ObjectID())
		if err != nil {
			return fmt.Errorf("failed to get meta: %s", err)
		}
		ref := object.ObjectReference{
			true,
			meta,
			0,
			ammo.ServiceID(),
			ammo.ObjectID(),
		}
		return object.WriteObjectReference(ref, buf)
	}(); err != nil {
		return fmt.Errorf("failed to serialize ammo: %s", err)
	}
	_, err = p.Call("ammo", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call ammo failed: %s", err)
	}
	return nil
}
