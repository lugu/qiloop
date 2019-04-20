package bus

import (
	"fmt"
	"github.com/lugu/qiloop/type/object"
)

// proxy is the parent strucuture for Service. It wraps Client and
// capture the service name.
type proxy struct {
	meta    object.MetaObject
	methods map[string]uint32
	client  Client
	service uint32
	object  uint32
}

// CallID construct a call message and send it to the client endpoint.
func (p proxy) CallID(actionID uint32, payload []byte) ([]byte, error) {
	return p.client.Call(p.service, p.object, actionID, payload)
}

// Call translates the name into an action id and send it to the client endpoint.
func (p proxy) Call(action string, payload []byte) ([]byte, error) {
	id, err := p.MethodID(action)
	if err != nil {
		return nil, fmt.Errorf("failed to find call %s: %s", action, err)
	}
	return p.CallID(id, payload)
}

// ServiceID returns the service identifier.
func (p proxy) ServiceID() uint32 {
	return p.service
}

// ObjectID returns the object identifier within the service.
func (p proxy) ObjectID() uint32 {
	return p.object
}

// SubscribeID returns a channel with the values of a signal or a
// property.
func (p proxy) SubscribeID(action uint32) (func(), chan []byte, error) {
	return p.client.Subscribe(p.service, p.object, action)
}

// Subscribe returns a channel with the values of a signal or a
// property.
func (p proxy) Subscribe(action string) (func(), chan []byte, error) {
	id, err := p.SignalID(action)
	if err != nil {
		id, err = p.PropertyID(action)
		if err != nil {
			return nil, nil,
				fmt.Errorf("cannot find signal or property %s",
					action)
		}
	}
	return p.client.Subscribe(p.service, p.object, id)
}

// MethodID resolve the name of the method using the meta object and
// returns the method id.
func (p proxy) MethodID(name string) (uint32, error) {
	return p.meta.MethodID(name)
}

// Disconnect closes the connection.
func (p proxy) Disconnect() error {
	// TODO
	return fmt.Errorf("proxy.Disconnect not yet implemented")
}

// SignalID resolve the name of the signal using the meta object and
// returns the signal id.
func (p proxy) SignalID(name string) (uint32, error) {
	return p.meta.SignalID(name)
}

// PropertyID resolve the name of the property using the meta object and
// returns the property id.
func (p proxy) PropertyID(name string) (uint32, error) {
	return p.meta.PropertyID(name)
}

// NewProxy construct a Proxy.
func NewProxy(client Client, meta object.MetaObject, service, object uint32) Proxy {
	methods := make(map[string]uint32)
	for id, method := range meta.Methods {
		methods[method.Name] = id
	}
	return proxy{meta, methods, client, service, object}
}
