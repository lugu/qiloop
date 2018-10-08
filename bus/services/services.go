// file generated. DO NOT EDIT.
package services

import (
	bus "github.com/lugu/qiloop/bus"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
)

type Object interface {
	object.Object
	bus.Proxy
}
type Server interface {
	bus.Proxy
	Authenticate(P0 map[string]value.Value) (map[string]value.Value, error)
}
type ServiceDirectory interface {
	object.Object
	bus.Proxy
	IsStatsEnabled() (bool, error)
	EnableStats(P0 bool) error
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
	SignalTraceObject(cancel chan int) (chan struct {
		P0 EventTrace
	}, error)
	SignalServiceAdded(cancel chan int) (chan struct {
		P0 uint32
		P1 string
	}, error)
	SignalServiceRemoved(cancel chan int) (chan struct {
		P0 uint32
		P1 string
	}, error)
}
type LogManager interface {
	object.Object
	bus.Proxy
	IsStatsEnabled() (bool, error)
	EnableStats(P0 bool) error
	Stats() (map[uint32]MethodStatistics, error)
	ClearStats() error
	IsTraceEnabled() (bool, error)
	EnableTrace(P0 bool) error
	Log(P0 []LogMessage) error
	CreateListener() (object.ObjectReference, error)
	GetListener() (object.ObjectReference, error)
	AddProvider(P0 object.ObjectReference) (int32, error)
	RemoveProvider(P0 int32) error
	SignalTraceObject(cancel chan int) (chan struct {
		P0 EventTrace
	}, error)
}
