package bus

import (
	"fmt"
	"math/rand"

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
		return nil, fmt.Errorf("find call %s: %s", action, err)
	}
	return p.CallID(id, payload)
}

// SubscribeID returns a channel with the values of a signal or a
// property.
func (p proxy) SubscribeID(action uint32) (func(), chan []byte, error) {
	cancel, bytes, err := p.client.Subscribe(p.service, p.object, action)
	if err != nil {
		return nil, nil, err

	}
	subscriptions := p.client.State(fmt.Sprintf("%d.%d.%d", p.service, p.object, action), 1)
	if subscriptions == 1 {
		handler := rand.Int()
		p.client.State(fmt.Sprintf("%d.%d.%d.handler", p.service, p.object, action), handler)
		obj := proxyObject{p}
		_, err := obj.RegisterEvent(p.object, action, uint64(handler))
		if err != nil {
			return nil, nil, err
		}
	}
	return func() {
		subscriptions := p.client.State(fmt.Sprintf("%d.%d.%d", p.service, p.object, action), -1)
		if subscriptions == 0 {
			handler := p.client.State(fmt.Sprintf("%d.%d.%d.handler", p.service, p.object, action), 0)
			p.client.State(fmt.Sprintf("%d.%d.%d.handler", p.service, p.object, action), -handler)
			obj := proxyObject{p}
			obj.UnregisterEvent(p.object, action, uint64(handler))
		}
		cancel()
	}, bytes, nil
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
	return p.SubscribeID(id)
}

// ServiceID returns the service identifier.
func (p proxy) ServiceID() uint32 {
	return p.service
}

// ObjectID returns the object identifier within the service.
func (p proxy) ObjectID() uint32 {
	return p.object
}

// MethodID resolve the name of the method using the meta object and
// returns the method id.
func (p proxy) MethodID(name string) (uint32, error) {
	return p.meta.MethodID(name)
}

// OnDisconnect registers a callback which is called when the network
// connection is unavailable.
func (p proxy) OnDisconnect(cb func(error)) error {
	return p.client.OnDisconnect(cb)
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

// ProxyService returns a reference to a remote service which can be used
// to create client side objects.
func (p proxy) ProxyService(sess Session) Service {
	c, ok := p.client.(*client)
	if !ok {
		panic("unexpected client implementation")
	}
	return NewServiceReference(sess, c.endpoint, p.service)
}

// NewProxy construct a Proxy.
func NewProxy(client Client, meta object.MetaObject, service, object uint32) Proxy {
	methods := make(map[string]uint32)
	for id, method := range meta.Methods {
		methods[method.Name] = id
	}
	return proxy{meta, methods, client, service, object}
}
