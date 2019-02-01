package client

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/type/object"
)

// Proxy is the parent strucuture for Service. It wraps Client and
// capture the service name.
type Proxy struct {
	meta    object.MetaObject
	methods map[string]uint32
	client  bus.Client
	service uint32
	object  uint32
}

// CallID construct a call message and send it to the client endpoint.
func (p Proxy) CallID(actionID uint32, payload []byte) ([]byte, error) {
	return p.client.Call(p.service, p.object, actionID, payload)
}

// Call translates the name into an action id and send it to the client endpoint.
func (p Proxy) Call(action string, payload []byte) ([]byte, error) {
	id, err := p.MethodID(action)
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

// SubscribeID returns a channel with the values of a signal or a
// property.
func (p Proxy) SubscribeID(action uint32, cancel chan int) (chan []byte, error) {
	return p.client.Subscribe(p.service, p.object, action, cancel)
}

// Subscribe returns a channel with the values of a signal or a
// property.
func (p Proxy) Subscribe(action string, cancel chan int) (chan []byte, error) {
	id, err := p.SignalID(action)
	if err != nil {
		id, err = p.PropertyID(action)
		if err != nil {
			return nil, fmt.Errorf("cannot find signal or property %s",
				action)
		}
	}
	return p.client.Subscribe(p.service, p.object, id, cancel)
}

// MethodID resolve the name of the method using the meta object and
// returns the method id.
func (p Proxy) MethodID(name string) (uint32, error) {
	return p.meta.MethodID(name)
}

// Disconnect closes the connection.
func (p Proxy) Disconnect() error {
	return fmt.Errorf("Proxy.Disconnect not yet implemented")
}

// SignalID resolve the name of the signal using the meta object and
// returns the signal id.
func (p Proxy) SignalID(name string) (uint32, error) {
	return p.meta.SignalID(name)
}

// PropertyID resolve the name of the property using the meta object and
// returns the property id.
func (p Proxy) PropertyID(name string) (uint32, error) {
	return p.meta.PropertyID(name)
}

// NewProxy construct a Proxy.
func NewProxy(client bus.Client, meta object.MetaObject, service, object uint32) Proxy {
	methods := make(map[string]uint32)
	for id, method := range meta.Methods {
		methods[method.Name] = id
	}
	return Proxy{meta, methods, client, service, object}
}
