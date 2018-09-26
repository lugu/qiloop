package bus

import (
	"github.com/lugu/qiloop/type/object"
)

// Client represents a client connection to a service.
type Client interface {
	Call(serviceID uint32, objectID uint32, actionID uint32, payload []byte) ([]byte, error)
	Subscribe(serviceID, objectID, signalID uint32, cancel chan int) (chan []byte, error)
}

type Proxy interface {
	Call(action string, payload []byte) ([]byte, error)
	CallID(action uint32, payload []byte) ([]byte, error)

	// SignalSubscribe returns a channel with the values of a signal
	SubscribeSignal(signal string, cancel chan int) (chan []byte, error)
	SubscribeID(signal uint32, cancel chan int) (chan []byte, error)

	MethodUid(name string) (uint32, error)
	SignalUid(name string) (uint32, error)

	// ServiceID returns the related service identifier
	ServiceID() uint32
	// ServiceID returns object identifier within the service
	// namespace.
	ObjectID() uint32

	// Disconnect stop the network connection to the remote object.
	Disconnect() error
}

type Session interface {
	Proxy(name string, objectID uint32) (Proxy, error)
	Object(ref object.ObjectReference) (object.Object, error)
	Destroy() error
}

type ActionWrapper func(payload []byte) ([]byte, error)
type Wrapper map[uint32]ActionWrapper
