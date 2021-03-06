package bus

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/conversion"
	"github.com/lugu/qiloop/type/encoding"
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

func ObjectReference(p Proxy) object.ObjectReference {
	return object.ObjectReference{
		MetaObject: *p.MetaObject(),
		ServiceID:  p.ServiceID(),
		ObjectID:   p.ObjectID(),
	}
}

// CallID construct a call message and send it to the client endpoint.
func (p proxy) CallID(actionID uint32, payload []byte) ([]byte, error) {
	return p.client.Call(p.ctx.Done(), p.service, p.object, actionID, payload)
}

// Call translates the name into an action id and send it to the client endpoint.
func (p proxy) Call(method, param, ret string, payload []byte) ([]byte, error) {
	id, sig, err := p.meta.MethodID(method, param)
	if err != nil {
		return nil, fmt.Errorf("missing method %s: %s", method, err)
	} else if ret != sig {
		return nil, fmt.Errorf("unexpected return: %s, expecting %s",
			sig, ret)
	}
	return p.CallID(id, payload)
}

func (p proxy) Call2(method string, args Params, ret Response) error {
	methodID, sig, err := p.MetaObject().MethodID(method, args.Signature())
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	permission := p.client.Channel().Cap()
	var e = encoding.NewEncoder(permission, &buf)
	if err := args.Write(e); err != nil {
		return fmt.Errorf("encode param: %w", err)
	}
	res, err := p.CallID(methodID, buf.Bytes())
	if err != nil {
		return err
	}
	buf2 := bytes.NewBuffer(res)
	dec := encoding.NewDecoder(permission, buf2)
	if sig == ret.Signature() {
		err = ret.Read(dec)
		if err != nil {
			return fmt.Errorf("decode result: %w", err)
		}
	} else {
		typ, err := signature.Parse(sig)
		if err != nil {
			return fmt.Errorf("failed to parse return signature: %w", err)
		}
		err = conversion.DecodeFrom(dec, ret.resp, typ.Type())
		if err != nil {
			return fmt.Errorf("failed to convert %s into %s: %w", sig,
				ret.Signature(), err)
		}
	}
	return nil
}

// SubscribeID returns a channel with the values of a signal or a
// property.
func (p proxy) SubscribeID(action uint32) (func(), chan []byte, error) {
	_, ok := p.meta.Signals[action]
	if !ok {
		_, ok := p.meta.Properties[action]
		if !ok {
			return nil, nil, fmt.Errorf("unknown signal: %d", action)
		}
	}
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
			err := obj.UnregisterEvent(p.object, action, uint64(handler))
			if err != nil {
				log.Printf("failed to unregister event action %d: %s", action, err)
			}
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
