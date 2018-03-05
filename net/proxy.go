package net

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
func NewProxy(c Client, service, object uint32) Proxy {
	return Proxy{c, service, object}
}
