// file generated. DO NOT EDIT.
package pingpong

import (
	"bytes"
	"fmt"
	net "github.com/lugu/qiloop/bus/net"
	server "github.com/lugu/qiloop/bus/server"
	session "github.com/lugu/qiloop/bus/session"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
)

type PingPong interface {
	Activate(sess *session.Session, serviceID, objectID uint32, signal PingPongSignalHelper)
	Ping(a string) error
}
type PingPongSignalHelper interface {
	SignalPong(a string) error
}
type stubPingPong struct {
	obj  *server.BasicObject
	impl PingPong
}

func PingPongObject(impl PingPong) server.Object {
	var stb stubPingPong
	stb.impl = impl
	stb.obj = server.NewObject(stb.metaObject())
	stb.obj.Wrapper[uint32(0x64)] = stb.Ping
	return &stb
}
func (s *stubPingPong) Activate(sess *session.Session, serviceID, objectID uint32) {
	s.obj.Activate(sess, serviceID, objectID)
	s.impl.Activate(sess, serviceID, objectID, s)
}
func (s *stubPingPong) Receive(msg *net.Message, from *server.Context) error {
	return s.obj.Receive(msg, from)
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
	err := s.obj.UpdateSignal(uint32(0x65), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalPong: %s", err)
	}
	return nil
}
func (s *stubPingPong) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "PingPong",
		Methods: map[uint32]object.MetaMethod{uint32(0x64): object.MetaMethod{
			Name:                "ping",
			ParametersSignature: "(s)",
			ReturnSignature:     "v",
			Uid:                 uint32(0x64),
		}},
		Signals: map[uint32]object.MetaSignal{uint32(0x65): object.MetaSignal{
			Name:      "pong",
			Signature: "(s)",
			Uid:       uint32(0x65),
		}},
	}
}
