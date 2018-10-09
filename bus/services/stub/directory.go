package stub

import (
	"bytes"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/services/impl"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/basic"
)

// ServiceDirectoryStub implements session.Object
type ServiceDirectoryStub struct {
	session.ObjectDispather
	impl impl.ServiceDirectoryInterface
}

func NewServiceDirectory(impl impl.ServiceDirectoryInterface) session.Object {
	var stub ServiceDirectoryStub
	stub.Wrapper = bus.Wrapper(make(map[uint32]bus.ActionWrapper))
	stub.Wrapper[100] = stub.Service
	return &stub
}

// Service unmarshall the argument, call the real implementation
// and marshall the response
func (s ServiceDirectoryStub) Service(payload []byte) ([]byte, error) {

	buf := bytes.NewBuffer(payload)
	serviceID, err := basic.ReadString(buf)
	if err != nil {
		return util.ErrorPaylad(err), nil
	}

	info, err := s.impl.Service(serviceID)
	if err != nil {
		return util.ErrorPaylad(err), nil
	}

	buf = bytes.NewBuffer(make([]byte, 0))
	err = services.WriteServiceInfo(info, buf)
	if err != nil {
		return util.ErrorPaylad(err), nil
	}

	return buf.Bytes(), nil
}

// FIXME: the rest of the methods to bind:
//
// 	Services() ([]ServiceInfo, error)
// 	RegisterService(P0 ServiceInfo) (uint32, error)
// 	UnregisterService(P0 uint32) error
// 	ServiceReady(P0 uint32) error
// 	UpdateServiceInfo(P0 ServiceInfo) error
// 	MachineId() (string, error)
// 	_socketOfService(P0 uint32) (object.ObjectReference, error)
