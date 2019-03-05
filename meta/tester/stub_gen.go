// Package tester contains a generated stub
// File generated. DO NOT EDIT.
package tester

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	net "github.com/lugu/qiloop/bus/net"
	server "github.com/lugu/qiloop/bus/server"
	generic "github.com/lugu/qiloop/bus/server/generic"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
)

// DummyImplementor interface of the service implementation
type DummyImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals an properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation server.Activation, helper DummySignalHelper) error
	OnTerminate()
	Hello() (BombProxy, error)
	Attack(b BombProxy) error
}

// DummySignalHelper provided to Dummy a companion object
type DummySignalHelper interface {
	SignalPing(msg string) error
	UpdateStatus(msg map[string]int32) error
	UpdateCoordinate(x int32, y int32) error
}

// stubDummy implements server.ServerObject.
type stubDummy struct {
	obj     generic.Object
	impl    DummyImplementor
	session bus.Session
}

// DummyObject returns an object using DummyImplementor
func DummyObject(impl DummyImplementor) server.ServerObject {
	var stb stubDummy
	stb.impl = impl
	stb.obj = generic.NewObject(stb.metaObject())
	stb.obj.Wrap(uint32(0x64), stb.Hello)
	stb.obj.Wrap(uint32(0x65), stb.Attack)
	return &stb
}
func (p *stubDummy) Activate(activation server.Activation) error {
	p.session = activation.Session
	p.obj.Activate(activation)
	return p.impl.Activate(activation, p)
}
func (p *stubDummy) OnTerminate() {
	p.impl.OnTerminate()
	p.obj.OnTerminate()
}
func (p *stubDummy) Receive(msg *net.Message, from *server.Context) error {
	return p.obj.Receive(msg, from)
}
func (p *stubDummy) Hello(payload []byte) ([]byte, error) {
	ret, callErr := p.impl.Hello()
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := func() error {
		meta, err := ret.MetaObject(ret.ObjectID())
		if err != nil {
			return fmt.Errorf("failed to get meta: %s", err)
		}
		ref := object.ObjectReference{
			true,
			meta,
			0,
			ret.ServiceID(),
			ret.ObjectID(),
		}
		return object.WriteObjectReference(ref, &out)
	}()
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubDummy) Attack(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	b, err := func() (BombProxy, error) {
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
		return nil, fmt.Errorf("cannot read b: %s", err)
	}
	callErr := p.impl.Attack(b)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubDummy) SignalPing(msg string) error {
	var buf bytes.Buffer
	if err := basic.WriteString(msg, &buf); err != nil {
		return fmt.Errorf("failed to serialize msg: %s", err)
	}
	err := p.obj.UpdateSignal(uint32(0x66), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalPing: %s", err)
	}
	return nil
}
func (p *stubDummy) UpdateStatus(msg map[string]int32) error {
	var buf bytes.Buffer
	if err := func() error {
		err := basic.WriteUint32(uint32(len(msg)), &buf)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range msg {
			err = basic.WriteString(k, &buf)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = basic.WriteInt32(v, &buf)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to serialize msg: %s", err)
	}
	err := p.obj.UpdateProperty(uint32(0x67), "({si})", buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update UpdateStatus: %s", err)
	}
	return nil
}
func (p *stubDummy) UpdateCoordinate(x int32, y int32) error {
	var buf bytes.Buffer
	if err := basic.WriteInt32(x, &buf); err != nil {
		return fmt.Errorf("failed to serialize x: %s", err)
	}
	if err := basic.WriteInt32(y, &buf); err != nil {
		return fmt.Errorf("failed to serialize y: %s", err)
	}
	err := p.obj.UpdateProperty(uint32(0x68), "(ii)", buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update UpdateCoordinate: %s", err)
	}
	return nil
}
func (p *stubDummy) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "Dummy",
		Methods: map[uint32]object.MetaMethod{
			uint32(0x64): {
				Name:                "hello",
				ParametersSignature: "()",
				ReturnSignature:     "o",
				Uid:                 uint32(0x64),
			},
			uint32(0x65): {
				Name:                "attack",
				ParametersSignature: "(o)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x65),
			},
		},
		Signals: map[uint32]object.MetaSignal{uint32(0x66): {
			Name:      "ping",
			Signature: "(s)",
			Uid:       uint32(0x66),
		}},
	}
}

// BombImplementor interface of the service implementation
type BombImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals an properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation server.Activation, helper BombSignalHelper) error
	OnTerminate()
}

// BombSignalHelper provided to Bomb a companion object
type BombSignalHelper interface {
	SignalBoom(energy int32) error
}

// stubBomb implements server.ServerObject.
type stubBomb struct {
	obj     generic.Object
	impl    BombImplementor
	session bus.Session
}

// BombObject returns an object using BombImplementor
func BombObject(impl BombImplementor) server.ServerObject {
	var stb stubBomb
	stb.impl = impl
	stb.obj = generic.NewObject(stb.metaObject())
	return &stb
}
func (p *stubBomb) Activate(activation server.Activation) error {
	p.session = activation.Session
	p.obj.Activate(activation)
	return p.impl.Activate(activation, p)
}
func (p *stubBomb) OnTerminate() {
	p.impl.OnTerminate()
	p.obj.OnTerminate()
}
func (p *stubBomb) Receive(msg *net.Message, from *server.Context) error {
	return p.obj.Receive(msg, from)
}
func (p *stubBomb) SignalBoom(energy int32) error {
	var buf bytes.Buffer
	if err := basic.WriteInt32(energy, &buf); err != nil {
		return fmt.Errorf("failed to serialize energy: %s", err)
	}
	err := p.obj.UpdateSignal(uint32(0x64), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalBoom: %s", err)
	}
	return nil
}
func (p *stubBomb) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "Bomb",
		Methods:     map[uint32]object.MetaMethod{},
		Signals: map[uint32]object.MetaSignal{uint32(0x64): {
			Name:      "boom",
			Signature: "(i)",
			Uid:       uint32(0x64),
		}},
	}
}
