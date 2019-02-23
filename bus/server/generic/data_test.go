// Package generic contains a generated stub
// File generated. DO NOT EDIT.
package generic

import (
	"bytes"
	"fmt"
	net "github.com/lugu/qiloop/bus/net"
	server "github.com/lugu/qiloop/bus/server"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	"io"
)

// Point is serializable
type Point struct {
	X int32
	Y int32
}

// ReadPoint unmarshalls Point
func ReadPoint(r io.Reader) (s Point, err error) {
	if s.X, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("failed to read X field: " + err.Error())
	}
	if s.Y, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("failed to read Y field: " + err.Error())
	}
	return s, nil
}

// WritePoint marshalls Point
func WritePoint(s Point, w io.Writer) (err error) {
	if err := basic.WriteInt32(s.X, w); err != nil {
		return fmt.Errorf("failed to write X field: " + err.Error())
	}
	if err := basic.WriteInt32(s.Y, w); err != nil {
		return fmt.Errorf("failed to write Y field: " + err.Error())
	}
	return nil
}

// Spacecraft interface of the service implementation
type Spacecraft interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals an properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation server.Activation, helper SpacecraftSignalHelper) error
	OnTerminate()
}

// SpacecraftSignalHelper provided to Spacecraft a companion object
type SpacecraftSignalHelper interface {
	SignalBoost(energy int32) error
	UpdatePosition(p Point) error
	UpdateSpeed(value int32) error
}

// stubSpacecraft implements server.ServerObject.
type stubSpacecraft struct {
	obj  Object
	impl Spacecraft
}

// SpacecraftObject returns an object using Spacecraft
func SpacecraftObject(impl Spacecraft) server.ServerObject {
	var stb stubSpacecraft
	stb.impl = impl
	stb.obj = NewObject(stb.metaObject())
	return &stb
}
func (s *stubSpacecraft) Activate(activation server.Activation) error {
	s.obj.Activate(activation)
	return s.impl.Activate(activation, s)
}
func (s *stubSpacecraft) OnTerminate() {
	s.impl.OnTerminate()
	s.obj.OnTerminate()
}
func (s *stubSpacecraft) Receive(msg *net.Message, from *server.Context) error {
	return s.obj.Receive(msg, from)
}
func (s *stubSpacecraft) SignalBoost(energy int32) error {
	var buf bytes.Buffer
	if err := basic.WriteInt32(energy, &buf); err != nil {
		return fmt.Errorf("failed to serialize energy: %s", err)
	}
	err := s.obj.UpdateSignal(uint32(0x66), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalBoost: %s", err)
	}
	return nil
}
func (s *stubSpacecraft) UpdatePosition(p Point) error {
	var buf bytes.Buffer
	if err := WritePoint(p, &buf); err != nil {
		return fmt.Errorf("failed to serialize p: %s", err)
	}
	err := s.obj.UpdateProperty(uint32(0x64), "((ii)<Point,x,y>)", buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update UpdatePosition: %s", err)
	}
	return nil
}
func (s *stubSpacecraft) UpdateSpeed(value int32) error {
	var buf bytes.Buffer
	if err := basic.WriteInt32(value, &buf); err != nil {
		return fmt.Errorf("failed to serialize value: %s", err)
	}
	err := s.obj.UpdateProperty(uint32(0x65), "(i)", buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update UpdateSpeed: %s", err)
	}
	return nil
}
func (s *stubSpacecraft) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "Spacecraft",
		Methods:     map[uint32]object.MetaMethod{},
		Signals: map[uint32]object.MetaSignal{uint32(0x66): {
			Name:      "boost",
			Signature: "(i)",
			Uid:       uint32(0x66),
		}},
	}
}
