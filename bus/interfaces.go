package bus

import (
	"github.com/lugu/qiloop/type/object"
)

// Client represents a client connection to a service.
type Client interface {
	Call(serviceID uint32, objectID uint32, methodID uint32, payload []byte) ([]byte, error)
	// Subscribe registers to a signal or a property. Returns a
	// cancel callback, a channel to receive the payload and an
	// error.
	Subscribe(serviceID, objectID, actionID uint32) (func(), chan []byte, error)
}

// Proxy represents a reference to a remote service. It allows to
// call methods and subscribe signals.
type Proxy interface {
	Call(action string, payload []byte) ([]byte, error)
	CallID(action uint32, payload []byte) ([]byte, error)

	// SignalSubscribe returns a channel with the values of a signal
	Subscribe(action string) (func(), chan []byte, error)
	SubscribeID(action uint32) (func(), chan []byte, error)

	MethodID(name string) (uint32, error)
	SignalID(name string) (uint32, error)
	PropertyID(name string) (uint32, error)

	// TODO
	// Properties() ([]string, error)
	// Property(name string) (value.Value, error)
	// SetProperty(name string, value value.Value) error

	// ServiceID returns the related service identifier
	ServiceID() uint32
	// ServiceID returns object identifier within the service
	// namespace.
	ObjectID() uint32

	// Disconnect stop the network connection to the remote object.
	Disconnect() error
}

// Session represents a connection to a bus: it is used to instanciate
// proxies to services and create new services.
type Session interface {
	Proxy(name string, objectID uint32) (Proxy, error)
	Object(ref object.ObjectReference) (object.Object, error)
	Destroy() error
}

// Server represents a local server. A server can provide with clients
// already connected.
type Server interface {
	Client() Client
}

// Namespace is used by a server to register new services. It allows
// custom name resolution. An implementation connects to the
// servide directory. Custom implementation is possible to allow
// interroperability with various naming services.
type Namespace interface {
	Reserve(name string) (uint32, error)
	Remove(serviceID uint32) error
	Enable(serviceID uint32) error
	Resolve(name string) (uint32, error)
	Session(s Server) Session
}
