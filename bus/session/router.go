package session

import (
	"github.com/lugu/qiloop/bus/net"
)

// scenario:
// - object for service zero is created
// - namespce for service zero is created
// - router is created
// - service zero namespace is added to the router
// - firewall is instanciated
// - server is instanciated and listen for incomming connection
// - object for the service directory is created
// - namespace for the service directory is created
// - service directory is added to the router
// - client A connects the server (a new endpoint is created)
// - server receive incomming message from endpoint
// - message is passed to the firewall
// - message is passed to the router
// - router pass the message to service zero
// - service zero namespace call service zero object
// - service zero object handle the authenticate calls
// - service zero object informs the firewall the client is authenticated
// - service zero namespace returns the response message
// - server receive incomming message from endpoint
// - message is passed to the firewall
// - message is passed to the router
// - router pass the message to service directory

// Server listen from incomming connections, set-up the end points and
// waits for messages.
type Server interface {
}

// Firewall ensures an endpoint talks only to the autorized service.
// Especially, it ensure authentication is passed.
type Firewall interface {
	Inspect(m net.Message, from net.EndPoint) error
}

// Router dispatch the incomming messages.
type Router interface {
	Add(n Namespace) (uint32, error)
	Remove(serviceID uint32) error
	Dispatch(m net.Message, from net.EndPoint) error
}

// Namespace represents a service
type Namespace interface {
	Add(o Object) (uint32, error)
	Remove(objectID uint32) error
	Dispatch(m net.Message, from net.EndPoint) error
	Ref(objectID uint32) error
	Unref(objectID uint32) error
}

type Object interface {
	CallID(action uint32, payload []byte) ([]byte, error)

	// SignalSubscribe returns a channel with the values of a signal
	SubscribeID(signal uint32, cancel chan int) (chan []byte, error)
}
