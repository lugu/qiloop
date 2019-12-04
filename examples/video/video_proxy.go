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

// Constructor gives access to remote services
type Constructor struct {
	session bus.Session
}

// Services gives access to the services constructor
func Services(s bus.Session) Constructor {
	return Constructor{session: s}
}

// ALVideoDevice is the abstract interface of the service
type ALVideoDevice interface {
	// SubscribeCamera calls the remote procedure
	SubscribeCamera(name string, cameraIndex int32, resolution int32, colorSpace int32, fps int32) (string, error)
	// Unsubscribe calls the remote procedure
	Unsubscribe(nameId string) (bool, error)
	// GetImageRemote calls the remote procedure
	GetImageRemote(name string) (value.Value, error)
}

// ALVideoDeviceProxy represents a proxy object to the service
type ALVideoDeviceProxy interface {
	bus.ObjectProxy
	ALVideoDevice
	// WithContext returns a new proxy. Calls to this proxy can be
	// cancelled by the context
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
func (c Constructor) ALVideoDevice() (ALVideoDeviceProxy, error) {
	proxy, err := c.session.Proxy("ALVideoDevice", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeALVideoDevice(c.session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyALVideoDevice) WithContext(ctx context.Context) ALVideoDeviceProxy {
	return MakeALVideoDevice(p.session, bus.WithContext(p.FIXMEProxy(), ctx))
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
	response, err := p.FIXMEProxy().Call("subscribeCamera", buf.Bytes())
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
	response, err := p.FIXMEProxy().Call("unsubscribe", buf.Bytes())
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
	response, err := p.FIXMEProxy().Call("getImageRemote", buf.Bytes())
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
