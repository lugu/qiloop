// Package main contains a generated proxy
// .

package main

import (
	"context"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	value "github.com/lugu/qiloop/type/value"
)

// ALVideoDeviceProxy represents a proxy object to the service
type ALVideoDeviceProxy interface {
	SubscribeCamera(name string, cameraIndex int32, resolution int32, colorSpace int32, fps int32) (string, error)
	Unsubscribe(nameId string) (bool, error)
	GetImageRemote(name string) (value.Value, error)
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
	WithContext(ctx context.Context) ALVideoDeviceProxy
}

// proxyALVideoDevice implements ALVideoDeviceProxy
type proxyALVideoDevice struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeALVideoDevice returns a specialized proxy.
func MakeALVideoDevice(sess bus.Session, proxy bus.Proxy) ALVideoDeviceProxy {
	return &proxyALVideoDevice{bus.MakeObject(proxy), sess}
}

// ALVideoDevice returns a proxy to a remote service
func ALVideoDevice(session bus.Session) (ALVideoDeviceProxy, error) {
	proxy, err := session.Proxy("ALVideoDevice", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeALVideoDevice(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyALVideoDevice) WithContext(ctx context.Context) ALVideoDeviceProxy {
	return MakeALVideoDevice(p.session, p.Proxy().WithContext(ctx))
}

// SubscribeCamera calls the remote procedure
func (p *proxyALVideoDevice) SubscribeCamera(name string, cameraIndex int32, resolution int32, colorSpace int32, fps int32) (string, error) {
	var ret string
	args := bus.NewParams("(siiii)", name, cameraIndex, resolution, colorSpace, fps)
	resp := bus.NewResponse("s", &ret)
	err := p.Proxy().Call2("subscribeCamera", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call subscribeCamera failed: %s", err)
	}
	return ret, nil
}

// Unsubscribe calls the remote procedure
func (p *proxyALVideoDevice) Unsubscribe(nameId string) (bool, error) {
	var ret bool
	args := bus.NewParams("(s)", nameId)
	resp := bus.NewResponse("b", &ret)
	err := p.Proxy().Call2("unsubscribe", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call unsubscribe failed: %s", err)
	}
	return ret, nil
}

// GetImageRemote calls the remote procedure
func (p *proxyALVideoDevice) GetImageRemote(name string) (value.Value, error) {
	var ret value.Value
	args := bus.NewParams("(s)", name)
	resp := bus.NewResponse("m", &ret)
	err := p.Proxy().Call2("getImageRemote", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call getImageRemote failed: %s", err)
	}
	return ret, nil
}
