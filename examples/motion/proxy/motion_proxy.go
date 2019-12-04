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

// Constructor gives access to remote services
type Constructor struct {
	session bus.Session
}

// Services gives access to the services constructor
func Services(s bus.Session) Constructor {
	return Constructor{session: s}
}

// ALMotion is the abstract interface of the service
type ALMotion interface {
	// WakeUp calls the remote procedure
	WakeUp() error
	// Rest calls the remote procedure
	Rest() error
	// MoveInit calls the remote procedure
	MoveInit() error
	// MoveTo calls the remote procedure
	MoveTo(x float32, y float32, theta float32) error
	// WaitUntilMoveIsFinished calls the remote procedure
	WaitUntilMoveIsFinished() error
}

// ALMotionProxy represents a proxy object to the service
type ALMotionProxy interface {
	bus.ObjectProxy
	ALMotion
	// WithContext returns a new proxy. Calls to this proxy can be
	// cancelled by the context
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
func (c Constructor) ALMotion() (ALMotionProxy, error) {
	proxy, err := c.session.Proxy("ALMotion", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeALMotion(c.session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyALMotion) WithContext(ctx context.Context) ALMotionProxy {
	return MakeALMotion(p.session, bus.WithContext(p.Proxy(), ctx))
}

// WakeUp calls the remote procedure
func (p *proxyALMotion) WakeUp() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Proxy().Call("wakeUp", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call wakeUp failed: %s", err)
	}
	return nil
}

// Rest calls the remote procedure
func (p *proxyALMotion) Rest() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Proxy().Call("rest", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call rest failed: %s", err)
	}
	return nil
}

// MoveInit calls the remote procedure
func (p *proxyALMotion) MoveInit() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Proxy().Call("moveInit", buf.Bytes())
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
	_, err = p.Proxy().Call("moveTo", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call moveTo failed: %s", err)
	}
	return nil
}

// WaitUntilMoveIsFinished calls the remote procedure
func (p *proxyALMotion) WaitUntilMoveIsFinished() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Proxy().Call("waitUntilMoveIsFinished", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call waitUntilMoveIsFinished failed: %s", err)
	}
	return nil
}
