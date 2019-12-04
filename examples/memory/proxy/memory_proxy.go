// Package proxy contains a generated proxy
// .

package proxy

import (
	"bytes"
	"context"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
	"log"
)

// Constructor gives access to remote services
type Constructor struct {
	session bus.Session
}

// Services gives access to the services constructor
func Services(s bus.Session) Constructor {
	return Constructor{session: s}
}

// ALTextToSpeech is the abstract interface of the service
type ALTextToSpeech interface {
	// Say calls the remote procedure
	Say(stringToSay string) error
}

// ALTextToSpeechProxy represents a proxy object to the service
type ALTextToSpeechProxy interface {
	bus.ObjectProxy
	ALTextToSpeech
	// WithContext returns a new proxy. Calls to this proxy can be
	// cancelled by the context
	WithContext(ctx context.Context) ALTextToSpeechProxy
}

// proxyALTextToSpeech implements ALTextToSpeechProxy
type proxyALTextToSpeech struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeALTextToSpeech returns a specialized proxy.
func MakeALTextToSpeech(sess bus.Session, proxy bus.Proxy) ALTextToSpeechProxy {
	return &proxyALTextToSpeech{bus.MakeObject(proxy), sess}
}

// ALTextToSpeech returns a proxy to a remote service
func (c Constructor) ALTextToSpeech() (ALTextToSpeechProxy, error) {
	proxy, err := c.session.Proxy("ALTextToSpeech", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeALTextToSpeech(c.session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyALTextToSpeech) WithContext(ctx context.Context) ALTextToSpeechProxy {
	return MakeALTextToSpeech(p.session, bus.WithContext(p.FIXMEProxy(), ctx))
}

// Say calls the remote procedure
func (p *proxyALTextToSpeech) Say(stringToSay string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(stringToSay, &buf); err != nil {
		return fmt.Errorf("serialize stringToSay: %s", err)
	}
	_, err = p.FIXMEProxy().Call("say", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call say failed: %s", err)
	}
	return nil
}

// ALMemory is the abstract interface of the service
type ALMemory interface {
	// GetEventList calls the remote procedure
	GetEventList() ([]string, error)
	// Subscriber calls the remote procedure
	Subscriber(eventName string) (SubscriberProxy, error)
}

// ALMemoryProxy represents a proxy object to the service
type ALMemoryProxy interface {
	bus.ObjectProxy
	ALMemory
	// WithContext returns a new proxy. Calls to this proxy can be
	// cancelled by the context
	WithContext(ctx context.Context) ALMemoryProxy
}

// proxyALMemory implements ALMemoryProxy
type proxyALMemory struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeALMemory returns a specialized proxy.
func MakeALMemory(sess bus.Session, proxy bus.Proxy) ALMemoryProxy {
	return &proxyALMemory{bus.MakeObject(proxy), sess}
}

// ALMemory returns a proxy to a remote service
func (c Constructor) ALMemory() (ALMemoryProxy, error) {
	proxy, err := c.session.Proxy("ALMemory", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeALMemory(c.session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyALMemory) WithContext(ctx context.Context) ALMemoryProxy {
	return MakeALMemory(p.session, bus.WithContext(p.FIXMEProxy(), ctx))
}

// GetEventList calls the remote procedure
func (p *proxyALMemory) GetEventList() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.FIXMEProxy().Call("getEventList", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getEventList failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getEventList response: %s", err)
	}
	return ret, nil
}

// Subscriber calls the remote procedure
func (p *proxyALMemory) Subscriber(eventName string) (SubscriberProxy, error) {
	var err error
	var ret SubscriberProxy
	var buf bytes.Buffer
	if err = basic.WriteString(eventName, &buf); err != nil {
		return ret, fmt.Errorf("serialize eventName: %s", err)
	}
	response, err := p.FIXMEProxy().Call("subscriber", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call subscriber failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (SubscriberProxy, error) {
		ref, err := object.ReadObjectReference(resp)
		if err != nil {
			return nil, fmt.Errorf("get meta: %s", err)
		}
		proxy, err := p.session.Object(ref)
		if err != nil {
			return nil, fmt.Errorf("get proxy: %s", err)
		}
		return MakeSubscriber(p.session, proxy), nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse subscriber response: %s", err)
	}
	return ret, nil
}

// Subscriber is the abstract interface of the service
type Subscriber interface {
	// SubscribeSignal subscribe to a remote signal
	SubscribeSignal() (unsubscribe func(), updates chan value.Value, err error)
}

// SubscriberProxy represents a proxy object to the service
type SubscriberProxy interface {
	bus.ObjectProxy
	Subscriber
	// WithContext returns a new proxy. Calls to this proxy can be
	// cancelled by the context
	WithContext(ctx context.Context) SubscriberProxy
}

// proxySubscriber implements SubscriberProxy
type proxySubscriber struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeSubscriber returns a specialized proxy.
func MakeSubscriber(sess bus.Session, proxy bus.Proxy) SubscriberProxy {
	return &proxySubscriber{bus.MakeObject(proxy), sess}
}

// Subscriber returns a proxy to a remote service
func (c Constructor) Subscriber() (SubscriberProxy, error) {
	proxy, err := c.session.Proxy("Subscriber", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeSubscriber(c.session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxySubscriber) WithContext(ctx context.Context) SubscriberProxy {
	return MakeSubscriber(p.session, bus.WithContext(p.FIXMEProxy(), ctx))
}

// SubscribeSignal subscribe to a remote property
func (p *proxySubscriber) SubscribeSignal() (func(), chan value.Value, error) {
	propertyID, err := p.FIXMEProxy().SignalID("signal")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "signal", err)
	}
	ch := make(chan value.Value)
	cancel, chPay, err := p.FIXMEProxy().SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := value.NewValue(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}
