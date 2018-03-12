package dummy

import (
	"fmt"
	"github.com/lugu/qiloop/object"
	"github.com/lugu/qiloop/session"
)

// Proxy is the parent strucuture for Service. It wraps Client and
// capture the service name.
type Proxy struct {
	meta    object.MetaObject
	methods map[string]uint32
	client  session.Client
	service uint32
	object  uint32
}

// CallID construct a call message and send it to the client endpoint.
func (p Proxy) CallID(actionID uint32, payload []byte) ([]byte, error) {
	return p.client.Call(p.service, p.object, actionID, payload)
}

// Call translates the name into an action id and send it to the client endpoint.
func (p Proxy) Call(action string, payload []byte) ([]byte, error) {
	id, err := p.meta.MethodUid(action)
	if err != nil {
		return nil, fmt.Errorf("failed to find call %s: %s", action, err)
	}
	return p.CallID(id, payload)
}

// ServiceID returns the service identifier.
func (p Proxy) ServiceID() uint32 {
	return p.service
}

// ObjectID returns the object identifier within the service.
func (p Proxy) ObjectID() uint32 {
	return p.object
}

// SignalStream returns a channel with the values of a signal
func (p Proxy) SignalStreamID(signal uint32, cancel chan int) (chan []byte, error) {
	return p.client.Stream(p.service, p.object, signal, cancel)
}

// SignalStream returns a channel with the values of a signal
func (p Proxy) SignalStream(signal string, cancel chan int) (chan []byte, error) {
	id, err := p.meta.SignalUid(signal)
	if err != nil {
		return nil, fmt.Errorf("failed to find signal %s: %s", signal, err)
	}
	return p.client.Stream(p.service, p.object, id, cancel)
}

// NewProxy construct a Proxy.
func NewProxy(client session.Client, meta object.MetaObject, service uint32, object uint32) Proxy {
	methods := make(map[string]uint32)
	for id, method := range meta.Methods {
		methods[method.Name] = id
	}
	return Proxy{meta, methods, client, service, object}
}
