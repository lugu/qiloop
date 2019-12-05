package bus

import (
	"context"
	"errors"

	"github.com/lugu/qiloop/type/object"
)

// ErrCancelled is returned when the call was cancelled.
var ErrCancelled = errors.New("Cancelled")

// Client represents a client connection to a service.
type Client interface {

	// Call initiates a remote procedure call. If cancel is
	// closed, the call is cancelled. cancel can be nil.
	// ErrCancelled is returned if the call was cancelled.
	Call(cancel <-chan struct{}, serviceID, objectID, methodID uint32, payload []byte) ([]byte, error)

	// Subscribe registers to a signal or a property. Returns a
	// cancel callback, a channel to receive the payload and an
	// error.
	Subscribe(serviceID, objectID, actionID uint32) (cancel func(), events chan []byte, err error)

	// OnDisconnect registers a callback which is called when the
	// network connection is closed. If the closure of the network
	// connection is initiated by the remote side, a non nil error
	// is passed to the call back. A nil callback returns nil.
	OnDisconnect(cb func(error)) error

	// Signal state machine: count the number of subscription to a
	// given signal in order to mutualize the calls to
	// RegisterEvent and UnregisterEvent.
	State(signal string, increment int) int
}

// Proxy represents a reference to a remote service. It allows to
// call methods and subscribe signals.
type Proxy interface {
	// CallID send a call message.
	// ErrCancelled is returned if the call was cancelled.
	CallID(action uint32, payload []byte) ([]byte, error)
	// Call calls CallID with the appropriate action ID.
	// ErrCancelled is returned if the call was cancelled.
	Call(action string, payload []byte) ([]byte, error)

	// SubscribeID returns a channel with the values of a
	// signal. Subscribe calls RegisterEvent and UnregisterEvent on
	// behalf of the user.
	Subscribe(action string) (cancel func(), events chan []byte, err error)
	// Subscribe calls Subscribe with the appropriate action ID.
	SubscribeID(action uint32) (cancel func(), events chan []byte, err error)

	// MetaObject can be used to introspect the proyx.
	MetaObject() object.MetaObject

	// MethodID returns the associate method ID.
	MethodID(name string) (uint32, error)
	// MethodID returns the associate signal ID.
	SignalID(name string) (uint32, error)
	// MethodID returns the associate property ID.
	PropertyID(name string) (uint32, error)

	// ServiceID returns the related service identifier
	ServiceID() uint32
	// ServiceID returns object identifier within the service
	// namespace.
	ObjectID() uint32

	// WithContext returns a Proxy with a lifespan associated with the
	// context. It can be used for deadline and canceallations.
	WithContext(ctx context.Context) Proxy

	// OnDisconnect registers a callback which is called when the
	// network connection is closed. If the closure of the network
	// connection is initiated by the remote side, a non nil error
	// is passed to the call back. A nil callback returns nil.
	OnDisconnect(cb func(error)) error

	// ProxyService returns a reference to a remote service which can be used
	// to create client side objects.
	ProxyService(s Session) Service
}

// Session represents a connection to a bus: it is used to instanciate
// proxies to services and create new services.
type Session interface {
	Proxy(name string, objectID uint32) (Proxy, error)
	Object(ref object.ObjectReference) (Proxy, error)
	Terminate() error
}

// Server represents a local server. A server can provide with clients
// already connected.
type Server interface {
	// NewService register a new service to the service directory.
	NewService(name string, object Actor) (Service, error)
	// Session returns a local session object which can be used to
	// access the server without authentication.
	Session() Session
	// Terminate stops the server.
	Terminate() error
	// Returns a channel to wait for the server terminaison.
	WaitTerminate() chan error

	// Client returns a direct Client.
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
