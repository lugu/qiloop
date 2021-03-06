package space

import (
	"bytes"
	"context"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	net "github.com/lugu/qiloop/bus/net"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
	"log"
)

// BombImplementor interface of the service implementation
type BombImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals and properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation bus.Activation, helper BombSignalHelper) error
	OnTerminate()
	// OnDelayChange is called when the property is updated.
	// Returns an error if the property value is not allowed
	OnDelayChange(duration int32) error
}

// BombSignalHelper provided to Bomb a companion object
type BombSignalHelper interface {
	SignalBoom(energy int32) error
	UpdateDelay(duration int32) error
}

// stubBomb implements server.Actor.
type stubBomb struct {
	impl      BombImplementor
	session   bus.Session
	service   bus.Service
	serviceID uint32
	signal    bus.SignalHandler
}

// BombObject returns an object using BombImplementor
func BombObject(impl BombImplementor) bus.Actor {
	var stb stubBomb
	stb.impl = impl
	obj := bus.NewBasicObject(&stb, stb.metaObject(), stb.onPropertyChange)
	stb.signal = obj
	return obj
}

// CreateBomb registers a new object to a service
// and returns a proxy to the newly created object
func CreateBomb(session bus.Session, service bus.Service, impl BombImplementor) (BombProxy, error) {
	obj := BombObject(impl)
	objectID, err := service.Add(obj)
	if err != nil {
		return nil, err
	}
	stb := &stubBomb{}
	meta := object.FullMetaObject(stb.metaObject())
	client := bus.DirectClient(obj)
	proxy := bus.NewProxy(client, meta, service.ServiceID(), objectID)
	return MakeBomb(session, proxy), nil
}
func (p *stubBomb) Activate(activation bus.Activation) error {
	p.session = activation.Session
	p.service = activation.Service
	p.serviceID = activation.ServiceID
	return p.impl.Activate(activation, p)
}
func (p *stubBomb) OnTerminate() {
	p.impl.OnTerminate()
}
func (p *stubBomb) Receive(msg *net.Message, from bus.Channel) error {
	// action dispatch
	switch msg.Header.Action {
	default:
		return from.SendError(msg, bus.ErrActionNotFound)
	}
}
func (p *stubBomb) onPropertyChange(name string, data []byte) error {
	switch name {
	case "delay":
		buf := bytes.NewBuffer(data)
		prop, err := basic.ReadInt32(buf)
		if err != nil {
			return fmt.Errorf("cannot read Delay: %s", err)
		}
		return p.impl.OnDelayChange(prop)
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}
func (p *stubBomb) SignalBoom(energy int32) error {
	var buf bytes.Buffer
	if err := basic.WriteInt32(energy, &buf); err != nil {
		return fmt.Errorf("serialize energy: %s", err)
	}
	err := p.signal.UpdateSignal(100, buf.Bytes())

	if err != nil {
		return fmt.Errorf("update SignalBoom: %s", err)
	}
	return nil
}
func (p *stubBomb) UpdateDelay(duration int32) error {
	var buf bytes.Buffer
	if err := basic.WriteInt32(duration, &buf); err != nil {
		return fmt.Errorf("serialize duration: %s", err)
	}
	err := p.signal.UpdateProperty(101, "i", buf.Bytes())

	if err != nil {
		return fmt.Errorf("update UpdateDelay: %s", err)
	}
	return nil
}
func (p *stubBomb) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "Bomb",
		Methods:     map[uint32]object.MetaMethod{},
		Properties: map[uint32]object.MetaProperty{101: {
			Name:      "delay",
			Signature: "i",
			Uid:       101,
		}},
		Signals: map[uint32]object.MetaSignal{100: {
			Name:      "boom",
			Signature: "i",
			Uid:       100,
		}},
	}
}

// SpacecraftImplementor interface of the service implementation
type SpacecraftImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals and properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation bus.Activation, helper SpacecraftSignalHelper) error
	OnTerminate()
	Shoot() (BombProxy, error)
	Ammo(ammo BombProxy) error
}

// SpacecraftSignalHelper provided to Spacecraft a companion object
type SpacecraftSignalHelper interface{}

// stubSpacecraft implements server.Actor.
type stubSpacecraft struct {
	impl      SpacecraftImplementor
	session   bus.Session
	service   bus.Service
	serviceID uint32
	signal    bus.SignalHandler
}

// SpacecraftObject returns an object using SpacecraftImplementor
func SpacecraftObject(impl SpacecraftImplementor) bus.Actor {
	var stb stubSpacecraft
	stb.impl = impl
	obj := bus.NewBasicObject(&stb, stb.metaObject(), stb.onPropertyChange)
	stb.signal = obj
	return obj
}

// CreateSpacecraft registers a new object to a service
// and returns a proxy to the newly created object
func CreateSpacecraft(session bus.Session, service bus.Service, impl SpacecraftImplementor) (SpacecraftProxy, error) {
	obj := SpacecraftObject(impl)
	objectID, err := service.Add(obj)
	if err != nil {
		return nil, err
	}
	stb := &stubSpacecraft{}
	meta := object.FullMetaObject(stb.metaObject())
	client := bus.DirectClient(obj)
	proxy := bus.NewProxy(client, meta, service.ServiceID(), objectID)
	return MakeSpacecraft(session, proxy), nil
}
func (p *stubSpacecraft) Activate(activation bus.Activation) error {
	p.session = activation.Session
	p.service = activation.Service
	p.serviceID = activation.ServiceID
	return p.impl.Activate(activation, p)
}
func (p *stubSpacecraft) OnTerminate() {
	p.impl.OnTerminate()
}
func (p *stubSpacecraft) Receive(msg *net.Message, from bus.Channel) error {
	// action dispatch
	switch msg.Header.Action {
	case 100:
		return p.Shoot(msg, from)
	case 101:
		return p.Ammo(msg, from)
	default:
		return from.SendError(msg, bus.ErrActionNotFound)
	}
}
func (p *stubSpacecraft) onPropertyChange(name string, data []byte) error {
	switch name {
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}
func (p *stubSpacecraft) Shoot(msg *net.Message, c bus.Channel) error {
	ret, callErr := p.impl.Shoot()

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := func() error {
		meta, err := ret.MetaObject(ret.Proxy().ObjectID())
		if err != nil {
			return fmt.Errorf("get meta: %s", err)
		}
		ref := object.ObjectReference{
			MetaObject: meta,
			ServiceID:  ret.Proxy().ServiceID(),
			ObjectID:   ret.Proxy().ObjectID(),
		}
		return object.WriteObjectReference(ref, &out)
	}()
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubSpacecraft) Ammo(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	ammo, err := func() (BombProxy, error) {
		ref, err := object.ReadObjectReference(buf)
		if err != nil {
			return nil, fmt.Errorf("get meta: %s", err)
		}
		if ref.ServiceID == p.serviceID && ref.ObjectID >= (1<<31) {
			actor := bus.NewClientObject(ref.ObjectID, c)
			ref.ObjectID, err = p.service.Add(actor)
			if err != nil {
				return nil, fmt.Errorf("add client object: %s", err)
			}
		}
		proxy, err := p.session.Object(ref)
		if err != nil {
			return nil, fmt.Errorf("get proxy: %s", err)
		}
		return MakeBomb(p.session, proxy), nil
	}()
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read ammo: %s", err))
	}
	callErr := p.impl.Ammo(ammo)

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	return c.SendReply(msg, out.Bytes())
}
func (p *stubSpacecraft) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "Spacecraft",
		Methods: map[uint32]object.MetaMethod{
			100: {
				Name:                "shoot",
				ParametersSignature: "()",
				ReturnSignature:     "o",
				Uid:                 100,
			},
			101: {
				Name:                "ammo",
				ParametersSignature: "(o)",
				ReturnSignature:     "v",
				Uid:                 101,
			},
		},
		Properties: map[uint32]object.MetaProperty{},
		Signals:    map[uint32]object.MetaSignal{},
	}
}

// BombProxy represents a proxy object to the service
type BombProxy interface {
	SubscribeBoom() (unsubscribe func(), updates chan int32, err error)
	GetDelay() (int32, error)
	SetDelay(int32) error
	SubscribeDelay() (unsubscribe func(), updates chan int32, err error)
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
	WithContext(ctx context.Context) BombProxy
}

// proxyBomb implements BombProxy
type proxyBomb struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeBomb returns a specialized proxy.
func MakeBomb(sess bus.Session, proxy bus.Proxy) BombProxy {
	return &proxyBomb{bus.MakeObject(proxy), sess}
}

// Bomb returns a proxy to a remote service
func Bomb(session bus.Session) (BombProxy, error) {
	proxy, err := session.Proxy("Bomb", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeBomb(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyBomb) WithContext(ctx context.Context) BombProxy {
	return MakeBomb(p.session, p.Proxy().WithContext(ctx))
}

// SubscribeBoom subscribe to a remote property
func (p *proxyBomb) SubscribeBoom() (func(), chan int32, error) {
	signalID, err := p.Proxy().MetaObject().SignalID("boom", "i")
	if err != nil {
		return nil, nil, fmt.Errorf("%s not available: %s", "boom", err)
	}
	ch := make(chan int32)
	cancel, chPay, err := p.Proxy().SubscribeID(signalID)
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
			e, err := basic.ReadInt32(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// GetDelay updates the property value
func (p *proxyBomb) GetDelay() (ret int32, err error) {
	name := value.String("delay")
	value, err := p.Property(name)
	if err != nil {
		return ret, fmt.Errorf("get property: %s", err)
	}
	var buf bytes.Buffer
	err = value.Write(&buf)
	if err != nil {
		return ret, fmt.Errorf("read response: %s", err)
	}
	s, err := basic.ReadString(&buf)
	if err != nil {
		return ret, fmt.Errorf("read signature: %s", err)
	}
	// check the signature
	sig := "i"
	if sig != s {
		return ret, fmt.Errorf("unexpected signature: %s instead of %s",
			s, sig)
	}
	ret, err = basic.ReadInt32(&buf)
	return ret, err
}

// SetDelay updates the property value
func (p *proxyBomb) SetDelay(update int32) error {
	name := value.String("delay")
	var buf bytes.Buffer
	err := basic.WriteInt32(update, &buf)
	if err != nil {
		return fmt.Errorf("marshall error: %s", err)
	}
	val := value.Opaque("i", buf.Bytes())
	return p.SetProperty(name, val)
}

// SubscribeDelay subscribe to a remote property
func (p *proxyBomb) SubscribeDelay() (func(), chan int32, error) {
	signalID, err := p.Proxy().MetaObject().PropertyID("delay", "i")
	if err != nil {
		return nil, nil, fmt.Errorf("%s not available: %s", "delay", err)
	}
	ch := make(chan int32)
	cancel, chPay, err := p.Proxy().SubscribeID(signalID)
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
			e, err := basic.ReadInt32(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// SpacecraftProxy represents a proxy object to the service
type SpacecraftProxy interface {
	Shoot() (BombProxy, error)
	Ammo(ammo BombProxy) error
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
	WithContext(ctx context.Context) SpacecraftProxy
}

// proxySpacecraft implements SpacecraftProxy
type proxySpacecraft struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeSpacecraft returns a specialized proxy.
func MakeSpacecraft(sess bus.Session, proxy bus.Proxy) SpacecraftProxy {
	return &proxySpacecraft{bus.MakeObject(proxy), sess}
}

// Spacecraft returns a proxy to a remote service
func Spacecraft(session bus.Session) (SpacecraftProxy, error) {
	proxy, err := session.Proxy("Spacecraft", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeSpacecraft(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxySpacecraft) WithContext(ctx context.Context) SpacecraftProxy {
	return MakeSpacecraft(p.session, p.Proxy().WithContext(ctx))
}

// Shoot calls the remote procedure
func (p *proxySpacecraft) Shoot() (BombProxy, error) {
	var ret BombProxy
	args := bus.NewParams("()")
	var retRef object.ObjectReference
	resp := bus.NewResponse("o", &retRef)
	err := p.Proxy().Call2("shoot", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call shoot failed: %s", err)
	}
	proxy, err := p.session.Object(retRef)
	if err != nil {
		return nil, fmt.Errorf("proxy: %s", err)
	}
	ret = MakeBomb(p.session, proxy)
	return ret, nil
}

// Ammo calls the remote procedure
func (p *proxySpacecraft) Ammo(ammo BombProxy) error {
	var ret struct{}
	args := bus.NewParams("(o)", bus.ObjectReference(ammo.Proxy()))
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("ammo", args, resp)
	if err != nil {
		return fmt.Errorf("call ammo failed: %s", err)
	}
	return nil
}
