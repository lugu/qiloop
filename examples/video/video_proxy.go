// Package main contains a generated proxy
// .

package main

import (
	"bytes"
	"context"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	basic "github.com/lugu/qiloop/type/basic"
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
	return MakeALVideoDevice(p.session, bus.WithContext(p.Proxy(), ctx))
}

// SubscribeCamera calls the remote procedure
func (p *proxyALVideoDevice) SubscribeCamera(name string, cameraIndex int32, resolution int32, colorSpace int32, fps int32) (string, error) {
	var err error
	var ret string
	var buf bytes.Buffer
	if err = basic.WriteString(name, &buf); err != nil {
		return ret, fmt.Errorf("serialize name: %s", err)
	}
	if err = basic.WriteInt32(cameraIndex, &buf); err != nil {
		return ret, fmt.Errorf("serialize cameraIndex: %s", err)
	}
	if err = basic.WriteInt32(resolution, &buf); err != nil {
		return ret, fmt.Errorf("serialize resolution: %s", err)
	}
	if err = basic.WriteInt32(colorSpace, &buf); err != nil {
		return ret, fmt.Errorf("serialize colorSpace: %s", err)
	}
	if err = basic.WriteInt32(fps, &buf); err != nil {
		return ret, fmt.Errorf("serialize fps: %s", err)
	}
	response, err := p.Proxy().Call("subscribeCamera", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call subscribeCamera failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadString(resp)
	if err != nil {
		return ret, fmt.Errorf("parse subscribeCamera response: %s", err)
	}
	return ret, nil
}

// Unsubscribe calls the remote procedure
func (p *proxyALVideoDevice) Unsubscribe(nameId string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(nameId, &buf); err != nil {
		return ret, fmt.Errorf("serialize nameId: %s", err)
	}
	response, err := p.Proxy().Call("unsubscribe", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call unsubscribe failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse unsubscribe response: %s", err)
	}
	return ret, nil
}

// GetImageRemote calls the remote procedure
func (p *proxyALVideoDevice) GetImageRemote(name string) (value.Value, error) {
	var err error
	var ret value.Value
	var buf bytes.Buffer
	if err = basic.WriteString(name, &buf); err != nil {
		return ret, fmt.Errorf("serialize name: %s", err)
	}
	response, err := p.Proxy().Call("getImageRemote", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getImageRemote failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = value.NewValue(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getImageRemote response: %s", err)
	}
	return ret, nil
}
