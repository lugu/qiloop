// Package proxy contains a generated proxy
// .

package proxy

import (
	"bytes"
	"context"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	basic "github.com/lugu/qiloop/type/basic"
)

// ALMotionProxy represents a proxy object to the service
type ALMotionProxy interface {
	WakeUp() error
	Rest() error
	MoveInit() error
	MoveTo(x float32, y float32, theta float32) error
	WaitUntilMoveIsFinished() error
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
	WithContext(ctx context.Context) ALMotionProxy
}

// proxyALMotion implements ALMotionProxy
type proxyALMotion struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeALMotion returns a specialized proxy.
func MakeALMotion(sess bus.Session, proxy bus.Proxy) ALMotionProxy {
	return &proxyALMotion{bus.MakeObject(proxy), sess}
}

// ALMotion returns a proxy to a remote service
func ALMotion(session bus.Session) (ALMotionProxy, error) {
	proxy, err := session.Proxy("ALMotion", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeALMotion(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyALMotion) WithContext(ctx context.Context) ALMotionProxy {
	return MakeALMotion(p.session, p.Proxy().WithContext(ctx))
}

// WakeUp calls the remote procedure
func (p *proxyALMotion) WakeUp() error {
	var err error
	var buf bytes.Buffer
	methodID, err := p.Proxy().MetaObject().MethodID("wakeUp", "()", "v")
	if err != nil {
		return err
	}
	_, err = p.Proxy().CallID(methodID, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call wakeUp failed: %s", err)
	}
	return nil
}

// Rest calls the remote procedure
func (p *proxyALMotion) Rest() error {
	var err error
	var buf bytes.Buffer
	methodID, err := p.Proxy().MetaObject().MethodID("rest", "()", "v")
	if err != nil {
		return err
	}
	_, err = p.Proxy().CallID(methodID, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call rest failed: %s", err)
	}
	return nil
}

// MoveInit calls the remote procedure
func (p *proxyALMotion) MoveInit() error {
	var err error
	var buf bytes.Buffer
	methodID, err := p.Proxy().MetaObject().MethodID("moveInit", "()", "v")
	if err != nil {
		return err
	}
	_, err = p.Proxy().CallID(methodID, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call moveInit failed: %s", err)
	}
	return nil
}

// MoveTo calls the remote procedure
func (p *proxyALMotion) MoveTo(x float32, y float32, theta float32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteFloat32(x, &buf); err != nil {
		return fmt.Errorf("serialize x: %s", err)
	}
	if err = basic.WriteFloat32(y, &buf); err != nil {
		return fmt.Errorf("serialize y: %s", err)
	}
	if err = basic.WriteFloat32(theta, &buf); err != nil {
		return fmt.Errorf("serialize theta: %s", err)
	}
	methodID, err := p.Proxy().MetaObject().MethodID("moveTo", "(fff)", "v")
	if err != nil {
		return err
	}
	_, err = p.Proxy().CallID(methodID, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call moveTo failed: %s", err)
	}
	return nil
}

// WaitUntilMoveIsFinished calls the remote procedure
func (p *proxyALMotion) WaitUntilMoveIsFinished() error {
	var err error
	var buf bytes.Buffer
	methodID, err := p.Proxy().MetaObject().MethodID("waitUntilMoveIsFinished", "()", "v")
	if err != nil {
		return err
	}
	_, err = p.Proxy().CallID(methodID, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call waitUntilMoveIsFinished failed: %s", err)
	}
	return nil
}
