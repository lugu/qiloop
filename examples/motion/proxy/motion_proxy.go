// Package proxy contains a generated proxy
// .

package proxy

import (
	"context"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
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
	var ret struct{}
	args := bus.NewParams("()")
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("wakeUp", args, resp)
	if err != nil {
		return fmt.Errorf("call wakeUp failed: %s", err)
	}
	return nil
}

// Rest calls the remote procedure
func (p *proxyALMotion) Rest() error {
	var ret struct{}
	args := bus.NewParams("()")
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("rest", args, resp)
	if err != nil {
		return fmt.Errorf("call rest failed: %s", err)
	}
	return nil
}

// MoveInit calls the remote procedure
func (p *proxyALMotion) MoveInit() error {
	var ret struct{}
	args := bus.NewParams("()")
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("moveInit", args, resp)
	if err != nil {
		return fmt.Errorf("call moveInit failed: %s", err)
	}
	return nil
}

// MoveTo calls the remote procedure
func (p *proxyALMotion) MoveTo(x float32, y float32, theta float32) error {
	var ret struct{}
	args := bus.NewParams("(fff)", x, y, theta)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("moveTo", args, resp)
	if err != nil {
		return fmt.Errorf("call moveTo failed: %s", err)
	}
	return nil
}

// WaitUntilMoveIsFinished calls the remote procedure
func (p *proxyALMotion) WaitUntilMoveIsFinished() error {
	var ret struct{}
	args := bus.NewParams("()")
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("waitUntilMoveIsFinished", args, resp)
	if err != nil {
		return fmt.Errorf("call waitUntilMoveIsFinished failed: %s", err)
	}
	return nil
}
