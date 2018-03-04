package net

type Proxy struct {
	client  net.Client
	service uint32
	object  uint32
}

func (p Proxy) Call(action uint32, payload []byte) ([]byte, error) {
	return p.client.Call(p.service, p.object, action, payload)
}

func NewProxy(c net.Client, service, object uint32) Proxy {
	return Proxy{c, service, object}
}
