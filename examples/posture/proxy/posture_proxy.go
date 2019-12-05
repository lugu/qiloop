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
	var err error
	var ret string
	var buf bytes.Buffer
	methodID, err := p.Proxy().MetaObject().MethodID("getPostureFamily", "()", "s")
	if err != nil {
		return ret, err
	}
	response, err := p.Proxy().CallID(methodID, buf.Bytes())
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
	methodID, err := p.Proxy().MetaObject().MethodID("goToPosture", "(sf)", "b")
	if err != nil {
		return ret, err
	}
	response, err := p.Proxy().CallID(methodID, buf.Bytes())
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
	methodID, err := p.Proxy().MetaObject().MethodID("applyPosture", "(sf)", "b")
	if err != nil {
		return ret, err
	}
	response, err := p.Proxy().CallID(methodID, buf.Bytes())
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
	methodID, err := p.Proxy().MetaObject().MethodID("stopMove", "()", "v")
	if err != nil {
		return err
	}
	_, err = p.Proxy().CallID(methodID, buf.Bytes())
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
	methodID, err := p.Proxy().MetaObject().MethodID("getPostureList", "()", "[s]")
	if err != nil {
		return ret, err
	}
	response, err := p.Proxy().CallID(methodID, buf.Bytes())
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
	methodID, err := p.Proxy().MetaObject().MethodID("getPostureFamilyList", "()", "[s]")
	if err != nil {
		return ret, err
	}
	response, err := p.Proxy().CallID(methodID, buf.Bytes())
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
	methodID, err := p.Proxy().MetaObject().MethodID("setMaxTryNumber", "(i)", "v")
	if err != nil {
		return err
	}
	_, err = p.Proxy().CallID(methodID, buf.Bytes())
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
	methodID, err := p.Proxy().MetaObject().MethodID("getPosture", "()", "s")
	if err != nil {
		return ret, err
	}
	response, err := p.Proxy().CallID(methodID, buf.Bytes())
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
