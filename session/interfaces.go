package session

// Client represents a client connection to a service.
type Client interface {
	Call(service uint32, object uint32, action uint32, payload []byte) ([]byte, error)
}

// Proxy is the parent strucuture for Service. It wraps Client and
// capture the service name.
type Proxy struct {
	client  Client
	service uint32
	object  uint32
}

// Call construct a call message and send it to the client endpoint.
func (p Proxy) Call(action uint32, payload []byte) ([]byte, error) {
	return p.client.Call(p.service, p.object, action, payload)
}

// NewProxy construct a Proxy.
func NewProxy(client Client, service, object uint32) Proxy {
	return Proxy{client, service, object}
}

type Session interface {
	Proxy(name string, object uint32) (Proxy, error)
}
