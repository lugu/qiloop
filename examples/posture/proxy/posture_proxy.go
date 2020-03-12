// Package proxy contains a generated proxy
// .

package proxy

import (
	"context"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
)

// ALRobotPostureProxy represents a proxy object to the service
type ALRobotPostureProxy interface {
	GetPostureFamily() (string, error)
	GoToPosture(postureName string, maxSpeedFraction float32) (bool, error)
	ApplyPosture(postureName string, maxSpeedFraction float32) (bool, error)
	StopMove() error
	GetPostureList() ([]string, error)
	GetPostureFamilyList() ([]string, error)
	SetMaxTryNumber(pMaxTryNumber int32) error
	GetPosture() (string, error)
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
	WithContext(ctx context.Context) ALRobotPostureProxy
}

// proxyALRobotPosture implements ALRobotPostureProxy
type proxyALRobotPosture struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeALRobotPosture returns a specialized proxy.
func MakeALRobotPosture(sess bus.Session, proxy bus.Proxy) ALRobotPostureProxy {
	return &proxyALRobotPosture{bus.MakeObject(proxy), sess}
}

// ALRobotPosture returns a proxy to a remote service
func ALRobotPosture(session bus.Session) (ALRobotPostureProxy, error) {
	proxy, err := session.Proxy("ALRobotPosture", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeALRobotPosture(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyALRobotPosture) WithContext(ctx context.Context) ALRobotPostureProxy {
	return MakeALRobotPosture(p.session, p.Proxy().WithContext(ctx))
}

// GetPostureFamily calls the remote procedure
func (p *proxyALRobotPosture) GetPostureFamily() (string, error) {
	var ret string
	args := bus.NewParams("()")
	resp := bus.NewResponse("s", &ret)
	err := p.Proxy().Call2("getPostureFamily", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call getPostureFamily failed: %s", err)
	}
	return ret, nil
}

// GoToPosture calls the remote procedure
func (p *proxyALRobotPosture) GoToPosture(postureName string, maxSpeedFraction float32) (bool, error) {
	var ret bool
	args := bus.NewParams("(sf)", postureName, maxSpeedFraction)
	resp := bus.NewResponse("b", &ret)
	err := p.Proxy().Call2("goToPosture", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call goToPosture failed: %s", err)
	}
	return ret, nil
}

// ApplyPosture calls the remote procedure
func (p *proxyALRobotPosture) ApplyPosture(postureName string, maxSpeedFraction float32) (bool, error) {
	var ret bool
	args := bus.NewParams("(sf)", postureName, maxSpeedFraction)
	resp := bus.NewResponse("b", &ret)
	err := p.Proxy().Call2("applyPosture", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call applyPosture failed: %s", err)
	}
	return ret, nil
}

// StopMove calls the remote procedure
func (p *proxyALRobotPosture) StopMove() error {
	var ret struct{}
	args := bus.NewParams("()")
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("stopMove", args, resp)
	if err != nil {
		return fmt.Errorf("call stopMove failed: %s", err)
	}
	return nil
}

// GetPostureList calls the remote procedure
func (p *proxyALRobotPosture) GetPostureList() ([]string, error) {
	var ret []string
	args := bus.NewParams("()")
	resp := bus.NewResponse("[s]", &ret)
	err := p.Proxy().Call2("getPostureList", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call getPostureList failed: %s", err)
	}
	return ret, nil
}

// GetPostureFamilyList calls the remote procedure
func (p *proxyALRobotPosture) GetPostureFamilyList() ([]string, error) {
	var ret []string
	args := bus.NewParams("()")
	resp := bus.NewResponse("[s]", &ret)
	err := p.Proxy().Call2("getPostureFamilyList", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call getPostureFamilyList failed: %s", err)
	}
	return ret, nil
}

// SetMaxTryNumber calls the remote procedure
func (p *proxyALRobotPosture) SetMaxTryNumber(pMaxTryNumber int32) error {
	var ret struct{}
	args := bus.NewParams("(i)", pMaxTryNumber)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("setMaxTryNumber", args, resp)
	if err != nil {
		return fmt.Errorf("call setMaxTryNumber failed: %s", err)
	}
	return nil
}

// GetPosture calls the remote procedure
func (p *proxyALRobotPosture) GetPosture() (string, error) {
	var ret string
	args := bus.NewParams("()")
	resp := bus.NewResponse("s", &ret)
	err := p.Proxy().Call2("getPosture", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call getPosture failed: %s", err)
	}
	return ret, nil
}
