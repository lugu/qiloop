// file generated. DO NOT EDIT.
package stage2

import (
	bus "github.com/lugu/qiloop/bus"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
)

type NewServices struct {
	session bus.Session
}

func Services(s bus.Session) NewServices {
	return NewServices{session: s}
}

type Server interface {
	bus.Proxy
	Authenticate(P0 map[string]value.Value) (map[string]value.Value, error)
}
type Object interface {
	object.Object
	bus.Proxy
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
