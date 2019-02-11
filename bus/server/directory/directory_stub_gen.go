// Package directory contains a generated stub
// File generated.
package directory

import (
	"bytes"
	"fmt"
	net "github.com/lugu/qiloop/bus/net"
	server "github.com/lugu/qiloop/bus/server"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	"io"
)

// ServiceDirectory interface of the service implementation
type ServiceDirectory interface {
	Activate(activation server.Activation, helper ServiceDirectorySignalHelper) error
	OnTerminate()
	Service(name string) (ServiceInfo, error)
	Services() ([]ServiceInfo, error)
	RegisterService(info ServiceInfo) (uint32, error)
	UnregisterService(serviceID uint32) error
	ServiceReady(serviceID uint32) error
	UpdateServiceInfo(info ServiceInfo) error
	MachineId() (string, error)
	_socketOfService(serviceID uint32) (object.ObjectReference, error)
}

// ServiceDirectorySignalHelper provided to ServiceDirectory a companion object
type ServiceDirectorySignalHelper interface {
	SignalServiceAdded(serviceID uint32, name string) error
	SignalServiceRemoved(serviceID uint32, name string) error
}

// stubServiceDirectory implements server.ServerObject.
type stubServiceDirectory struct {
	obj  *server.BasicObject
	impl ServiceDirectory
}

// ServiceDirectoryObject returns an object using ServiceDirectory
func ServiceDirectoryObject(impl ServiceDirectory) server.ServerObject {
	var stb stubServiceDirectory
	stb.impl = impl
	stb.obj = server.NewObject(stb.metaObject())
	stb.obj.Wrap(uint32(0x64), stb.Service)
	stb.obj.Wrap(uint32(0x65), stb.Services)
	stb.obj.Wrap(uint32(0x66), stb.RegisterService)
	stb.obj.Wrap(uint32(0x67), stb.UnregisterService)
	stb.obj.Wrap(uint32(0x68), stb.ServiceReady)
	stb.obj.Wrap(uint32(0x69), stb.UpdateServiceInfo)
	stb.obj.Wrap(uint32(0x6c), stb.MachineId)
	stb.obj.Wrap(uint32(0x6d), stb._socketOfService)
	return &stb
}
func (s *stubServiceDirectory) Activate(activation server.Activation) error {
	s.obj.Activate(activation)
	return s.impl.Activate(activation, s)
}
func (s *stubServiceDirectory) OnTerminate() {
	s.impl.OnTerminate()
	s.obj.OnTerminate()
}
func (s *stubServiceDirectory) Receive(msg *net.Message, from *server.Context) error {
	return s.obj.Receive(msg, from)
}
func (s *stubServiceDirectory) Service(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	name, err := basic.ReadString(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read name: %s", err)
	}
	ret, callErr := s.impl.Service(name)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := WriteServiceInfo(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (s *stubServiceDirectory) Services(payload []byte) ([]byte, error) {
	ret, callErr := s.impl.Services()
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := func() error {
		err := basic.WriteUint32(uint32(len(ret)), &out)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range ret {
			err = WriteServiceInfo(v, &out)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}()
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (s *stubServiceDirectory) RegisterService(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	info, err := ReadServiceInfo(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read info: %s", err)
	}
	ret, callErr := s.impl.RegisterService(info)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteUint32(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (s *stubServiceDirectory) UnregisterService(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	serviceID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read serviceID: %s", err)
	}
	callErr := s.impl.UnregisterService(serviceID)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (s *stubServiceDirectory) ServiceReady(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	serviceID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read serviceID: %s", err)
	}
	callErr := s.impl.ServiceReady(serviceID)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (s *stubServiceDirectory) UpdateServiceInfo(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	info, err := ReadServiceInfo(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read info: %s", err)
	}
	callErr := s.impl.UpdateServiceInfo(info)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (s *stubServiceDirectory) MachineId(payload []byte) ([]byte, error) {
	ret, callErr := s.impl.MachineId()
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
func (s *stubServiceDirectory) _socketOfService(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	serviceID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read serviceID: %s", err)
	}
	ret, callErr := s.impl._socketOfService(serviceID)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := object.WriteObjectReference(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (s *stubServiceDirectory) SignalServiceAdded(serviceID uint32, name string) error {
	var buf bytes.Buffer
	if err := basic.WriteUint32(serviceID, &buf); err != nil {
		return fmt.Errorf("failed to serialize serviceID: %s", err)
	}
	if err := basic.WriteString(name, &buf); err != nil {
		return fmt.Errorf("failed to serialize name: %s", err)
	}
	err := s.obj.UpdateSignal(uint32(0x6a), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalServiceAdded: %s", err)
	}
	return nil
}
func (s *stubServiceDirectory) SignalServiceRemoved(serviceID uint32, name string) error {
	var buf bytes.Buffer
	if err := basic.WriteUint32(serviceID, &buf); err != nil {
		return fmt.Errorf("failed to serialize serviceID: %s", err)
	}
	if err := basic.WriteString(name, &buf); err != nil {
		return fmt.Errorf("failed to serialize name: %s", err)
	}
	err := s.obj.UpdateSignal(uint32(0x6b), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalServiceRemoved: %s", err)
	}
	return nil
}
func (s *stubServiceDirectory) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "ServiceDirectory",
		Methods: map[uint32]object.MetaMethod{
			uint32(0x64): {
				Name:                "service",
				ParametersSignature: "(s)",
				ReturnSignature:     "(sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>",
				Uid:                 uint32(0x64),
			},
			uint32(0x65): {
				Name:                "services",
				ParametersSignature: "()",
				ReturnSignature:     "[(sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>]",
				Uid:                 uint32(0x65),
			},
			uint32(0x66): {
				Name:                "registerService",
				ParametersSignature: "((sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>)",
				ReturnSignature:     "I",
				Uid:                 uint32(0x66),
			},
			uint32(0x67): {
				Name:                "unregisterService",
				ParametersSignature: "(I)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x67),
			},
			uint32(0x68): {
				Name:                "serviceReady",
				ParametersSignature: "(I)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x68),
			},
			uint32(0x69): {
				Name:                "updateServiceInfo",
				ParametersSignature: "((sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x69),
			},
			uint32(0x6c): {
				Name:                "machineId",
				ParametersSignature: "()",
				ReturnSignature:     "s",
				Uid:                 uint32(0x6c),
			},
			uint32(0x6d): {
				Name:                "_socketOfService",
				ParametersSignature: "(I)",
				ReturnSignature:     "o",
				Uid:                 uint32(0x6d),
			},
		},
		Signals: map[uint32]object.MetaSignal{
			uint32(0x6a): {
				Name:      "serviceAdded",
				Signature: "(Is)",
				Uid:       uint32(0x6a),
			},
			uint32(0x6b): {
				Name:      "serviceRemoved",
				Signature: "(Is)",
				Uid:       uint32(0x6b),
			},
		},
	}
}

// ServiceInfo is serializable
type ServiceInfo struct {
	Name      string
	ServiceId uint32
	MachineId string
	ProcessId uint32
	Endpoints []string
	SessionId string
}

// ReadServiceInfo unmarshalls ServiceInfo
func ReadServiceInfo(r io.Reader) (s ServiceInfo, err error) {
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	if s.ServiceId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read ServiceId field: " + err.Error())
	}
	if s.MachineId, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read MachineId field: " + err.Error())
	}
	if s.ProcessId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read ProcessId field: " + err.Error())
	}
	if s.Endpoints, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(r)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Endpoints field: " + err.Error())
	}
	if s.SessionId, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read SessionId field: " + err.Error())
	}
	return s, nil
}

// WriteServiceInfo marshalls ServiceInfo
func WriteServiceInfo(s ServiceInfo, w io.Writer) (err error) {
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
	}
	if err := basic.WriteUint32(s.ServiceId, w); err != nil {
		return fmt.Errorf("failed to write ServiceId field: " + err.Error())
	}
	if err := basic.WriteString(s.MachineId, w); err != nil {
		return fmt.Errorf("failed to write MachineId field: " + err.Error())
	}
	if err := basic.WriteUint32(s.ProcessId, w); err != nil {
		return fmt.Errorf("failed to write ProcessId field: " + err.Error())
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Endpoints)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.Endpoints {
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Endpoints field: " + err.Error())
	}
	if err := basic.WriteString(s.SessionId, w); err != nil {
		return fmt.Errorf("failed to write SessionId field: " + err.Error())
	}
	return nil
}
