package bus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/lugu/qiloop/type/object"
)

// proxy is the parent strucuture for Service. It wraps Client and
// capture the service name.
type proxy struct {
	meta    object.MetaObject
	ctx     context.Context
	methods map[string]uint32
	client  Client
	service uint32
	object  uint32
}

// CallID construct a call message and send it to the client endpoint.
func (p proxy) CallID(actionID uint32, payload []byte) ([]byte, error) {
	return p.client.Call(p.ctx.Done(), p.service, p.object, actionID, payload)
}

// Call translates the name into an action id and send it to the client endpoint.
func (p proxy) Call(method, param, ret string, payload []byte) ([]byte, error) {
	id, err := p.meta.MethodID(method, param, ret)
	if err != nil {
		return nil, fmt.Errorf("find call %s: %s", method, err)
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

// ServiceID returns the service identifier.
func (p proxy) ServiceID() uint32 {
	return p.service
}

// ObjectID returns the object identifier within the service.
func (p proxy) ObjectID() uint32 {
	return p.object
}

// OnDisconnect registers a callback which is called when the network
// connection is unavailable.
func (p proxy) OnDisconnect(cb func(error)) error {
	return p.client.OnDisconnect(cb)
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

// MetaObject can be used to introspect the proyx.
func (p proxy) MetaObject() *object.MetaObject {
	return &p.meta
}

// WithContext returns a Proxy with a lifespan associated with the
// context. It can be used for deadline and canceallations.
func (p proxy) WithContext(ctx context.Context) Proxy {
	return proxy{p.meta, ctx, p.methods, p.client, p.service, p.object}
}

// NewProxy construct a Proxy.
// TODO: return a *proxy instead of a proxy.
func NewProxy(client Client, meta object.MetaObject, service, object uint32) Proxy {
	methods := make(map[string]uint32)
	for id, method := range meta.Methods {
		methods[method.Name] = id
	}
	return proxy{meta, context.Background(), methods, client, service, object}
}
