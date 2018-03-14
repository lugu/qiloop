// file generated. DO NOT EDIT.
package services

import (
	bus "github.com/lugu/qiloop/bus"
	object "github.com/lugu/qiloop/object"
)

type IServer interface {
	bus.Proxy
	object.Object
}
type IObject interface {
	bus.Proxy
	object.Object
}
type IServiceDirectory interface {
	bus.Proxy
	object.Object
	Stats() (map[uint32]MethodStatistics, error)
	ClearStats() error
	IsTraceEnabled() (bool, error)
	EnableTrace(P0 bool) error
	Service(P0 string) (ServiceInfo, error)
	Services() ([]ServiceInfo, error)
	RegisterService(P0 ServiceInfo) (uint32, error)
	UnregisterService(P0 uint32) error
	ServiceReady(P0 uint32) error
	UpdateServiceInfo(P0 ServiceInfo) error
	MachineId() (string, error)
	_socketOfService(P0 uint32) (object.ObjectReference, error)
}
type ILogManager interface {
	bus.Proxy
	object.Object
	Stats() (map[uint32]MethodStatistics, error)
	ClearStats() error
	IsTraceEnabled() (bool, error)
	EnableTrace(P0 bool) error
	Log(P0 []LogMessage) error
	CreateListener() (object.ObjectReference, error)
	GetListener() (object.ObjectReference, error)
	AddProvider(P0 object.ObjectReference) (int32, error)
	RemoveProvider(P0 int32) error
}
