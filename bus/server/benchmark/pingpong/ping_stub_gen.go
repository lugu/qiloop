// Package pingpong contains a generated stub
// File generated. DO NOT EDIT.
package pingpong

import (
	"bytes"
	"fmt"
	net "github.com/lugu/qiloop/bus/net"
	server "github.com/lugu/qiloop/bus/server"
	generic "github.com/lugu/qiloop/bus/server/generic"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
)

// PingPongImplementor interface of the service implementation
type PingPongImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals an properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation server.Activation, helper PingPongSignalHelper) error
	OnTerminate()
	Hello(a string) (string, error)
	Ping(a string) error
}

// PingPongSignalHelper provided to PingPong a companion object
type PingPongSignalHelper interface {
	SignalPong(a string) error
}

// stubPingPong implements server.ServerObject.
type stubPingPong struct {
	obj  generic.Object
	impl PingPongImplementor
}

// PingPongObject returns an object using PingPongImplementor
func PingPongObject(impl PingPongImplementor) server.ServerObject {
	var stb stubPingPong
	stb.impl = impl
	stb.obj = generic.NewObject(stb.metaObject())
	stb.obj.Wrap(uint32(0x64), stb.Hello)
	stb.obj.Wrap(uint32(0x65), stb.Ping)
	return &stb
}
func (s *stubPingPong) Activate(activation server.Activation) error {
	s.obj.Activate(activation)
	return s.impl.Activate(activation, s)
}
func (s *stubPingPong) OnTerminate() {
	s.impl.OnTerminate()
	s.obj.OnTerminate()
}
func (s *stubPingPong) Receive(msg *net.Message, from *server.Context) error {
	return s.obj.Receive(msg, from)
}
func (s *stubPingPong) Hello(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	a, err := basic.ReadString(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read a: %s", err)
	}
	ret, callErr := s.impl.Hello(a)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteString(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (s *stubPingPong) Ping(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	a, err := basic.ReadString(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read a: %s", err)
	}
	callErr := s.impl.Ping(a)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (s *stubPingPong) SignalPong(a string) error {
	var buf bytes.Buffer
	if err := basic.WriteString(a, &buf); err != nil {
		return fmt.Errorf("failed to serialize a: %s", err)
	}
	err := s.obj.UpdateSignal(uint32(0x66), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalPong: %s", err)
	}
	return nil
}
func (s *stubPingPong) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "PingPong",
		Methods: map[uint32]object.MetaMethod{
			uint32(0x64): {
				Name:                "hello",
				ParametersSignature: "(s)",
				ReturnSignature:     "s",
				Uid:                 uint32(0x64),
			},
			uint32(0x65): {
				Name:                "ping",
				ParametersSignature: "(s)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x65),
			},
		},
		Signals: map[uint32]object.MetaSignal{uint32(0x66): {
			Name:      "pong",
			Signature: "(s)",
			Uid:       uint32(0x66),
		}},
	}
}
