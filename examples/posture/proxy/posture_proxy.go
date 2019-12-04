// Package proxy contains a generated proxy
// .

package proxy

import (
	"bytes"
	"context"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
)

// Constructor gives access to remote services
type Constructor struct {
	session bus.Session
}

// Services gives access to the services constructor
func Services(s bus.Session) Constructor {
	return Constructor{session: s}
}

// ALRobotPosture is the abstract interface of the service
type ALRobotPosture interface {
	// GetPostureFamily calls the remote procedure
	GetPostureFamily() (string, error)
	// GoToPosture calls the remote procedure
	GoToPosture(postureName string, maxSpeedFraction float32) (bool, error)
	// ApplyPosture calls the remote procedure
	ApplyPosture(postureName string, maxSpeedFraction float32) (bool, error)
	// StopMove calls the remote procedure
	StopMove() error
	// GetPostureList calls the remote procedure
	GetPostureList() ([]string, error)
	// GetPostureFamilyList calls the remote procedure
	GetPostureFamilyList() ([]string, error)
	// SetMaxTryNumber calls the remote procedure
	SetMaxTryNumber(pMaxTryNumber int32) error
	// GetPosture calls the remote procedure
	GetPosture() (string, error)
}

// ALRobotPostureProxy represents a proxy object to the service
type ALRobotPostureProxy interface {
	object.Object
	bus.Proxy
	ALRobotPosture
	// WithContext returns a new proxy. Calls to this proxy can be
	// cancelled by the context
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
func (c Constructor) ALRobotPosture() (ALRobotPostureProxy, error) {
	proxy, err := c.session.Proxy("ALRobotPosture", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeALRobotPosture(c.session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyALRobotPosture) WithContext(ctx context.Context) ALRobotPostureProxy {
	return MakeALRobotPosture(p.session, bus.WithContext(p, ctx))
}

// GetPostureFamily calls the remote procedure
func (p *proxyALRobotPosture) GetPostureFamily() (string, error) {
	var err error
	var ret string
	var buf bytes.Buffer
	response, err := p.FIXMEProxy().Call("getPostureFamily", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getPostureFamily failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadString(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getPostureFamily response: %s", err)
	}
	return ret, nil
}

// GoToPosture calls the remote procedure
func (p *proxyALRobotPosture) GoToPosture(postureName string, maxSpeedFraction float32) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(postureName, &buf); err != nil {
		return ret, fmt.Errorf("serialize postureName: %s", err)
	}
	if err = basic.WriteFloat32(maxSpeedFraction, &buf); err != nil {
		return ret, fmt.Errorf("serialize maxSpeedFraction: %s", err)
	}
	response, err := p.FIXMEProxy().Call("goToPosture", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call goToPosture failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse goToPosture response: %s", err)
	}
	return ret, nil
}

// ApplyPosture calls the remote procedure
func (p *proxyALRobotPosture) ApplyPosture(postureName string, maxSpeedFraction float32) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(postureName, &buf); err != nil {
		return ret, fmt.Errorf("serialize postureName: %s", err)
	}
	if err = basic.WriteFloat32(maxSpeedFraction, &buf); err != nil {
		return ret, fmt.Errorf("serialize maxSpeedFraction: %s", err)
	}
	response, err := p.FIXMEProxy().Call("applyPosture", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call applyPosture failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse applyPosture response: %s", err)
	}
	return ret, nil
}

// StopMove calls the remote procedure
func (p *proxyALRobotPosture) StopMove() error {
	var err error
	var buf bytes.Buffer
	_, err = p.FIXMEProxy().Call("stopMove", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call stopMove failed: %s", err)
	}
	return nil
}

// GetPostureList calls the remote procedure
func (p *proxyALRobotPosture) GetPostureList() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.FIXMEProxy().Call("getPostureList", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getPostureList failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getPostureList response: %s", err)
	}
	return ret, nil
}

// GetPostureFamilyList calls the remote procedure
func (p *proxyALRobotPosture) GetPostureFamilyList() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.FIXMEProxy().Call("getPostureFamilyList", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getPostureFamilyList failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getPostureFamilyList response: %s", err)
	}
	return ret, nil
}

// SetMaxTryNumber calls the remote procedure
func (p *proxyALRobotPosture) SetMaxTryNumber(pMaxTryNumber int32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteInt32(pMaxTryNumber, &buf); err != nil {
		return fmt.Errorf("serialize pMaxTryNumber: %s", err)
	}
	_, err = p.FIXMEProxy().Call("setMaxTryNumber", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setMaxTryNumber failed: %s", err)
	}
	return nil
}

// GetPosture calls the remote procedure
func (p *proxyALRobotPosture) GetPosture() (string, error) {
	var err error
	var ret string
	var buf bytes.Buffer
	response, err := p.FIXMEProxy().Call("getPosture", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getPosture failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadString(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getPosture response: %s", err)
	}
	return ret, nil
}
