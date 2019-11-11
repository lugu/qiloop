package bus

import (
	"bytes"
	"fmt"
	net "github.com/lugu/qiloop/bus/net"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
	"io"
	"log"
)

// ServiceZeroImplementor interface of the service implementation
type ServiceZeroImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals and properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation Activation, helper ServiceZeroSignalHelper) error
	OnTerminate()
	Authenticate(capability map[string]value.Value) (map[string]value.Value, error)
}

// ServiceZeroSignalHelper provided to ServiceZero a companion object
type ServiceZeroSignalHelper interface{}

// stubServiceZero implements server.Actor.
type stubServiceZero struct {
	impl      ServiceZeroImplementor
	session   Session
	service   Service
	serviceID uint32
	signal    SignalHandler
}

// ServiceZeroObject returns an object using ServiceZeroImplementor
func ServiceZeroObject(impl ServiceZeroImplementor) Actor {
	var stb stubServiceZero
	stb.impl = impl
	obj := NewBasicObject(&stb, stb.metaObject(), stb.onPropertyChange)
	stb.signal = obj
	return obj
}
func (p *stubServiceZero) Activate(activation Activation) error {
	p.session = activation.Session
	p.service = activation.Service
	p.serviceID = activation.ServiceID
	return p.impl.Activate(activation, p)
}
func (p *stubServiceZero) OnTerminate() {
	p.impl.OnTerminate()
}
func (p *stubServiceZero) Receive(msg *net.Message, from Channel) error {
	// action dispatch
	switch msg.Header.Action {
	case 8:
		return p.Authenticate(msg, from)
	default:
		return from.SendError(msg, ErrActionNotFound)
	}
}
func (p *stubServiceZero) onPropertyChange(name string, data []byte) error {
	switch name {
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}
func (p *stubServiceZero) Authenticate(msg *net.Message, c Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	capability, err := func() (m map[string]value.Value, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return m, fmt.Errorf("read map size: %s", err)
		}
		m = make(map[string]value.Value, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(buf)
			if err != nil {
				return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
			}
			v, err := value.NewValue(buf)
			if err != nil {
				return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read capability: %s", err))
	}
	ret, callErr := p.impl.Authenticate(capability)

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := func() error {
		err := basic.WriteUint32(uint32(len(ret)), &out)
		if err != nil {
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range ret {
			err = basic.WriteString(k, &out)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = v.Write(&out)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}()
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubServiceZero) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "ServiceZero",
		Methods: map[uint32]object.MetaMethod{8: {
			Name:                "authenticate",
			ParametersSignature: "({sm})",
			ReturnSignature:     "{sm}",
			Uid:                 8,
		}},
		Properties: map[uint32]object.MetaProperty{},
		Signals:    map[uint32]object.MetaSignal{},
	}
}

// ObjectImplementor interface of the service implementation
type ObjectImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals and properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation Activation, helper ObjectSignalHelper) error
	OnTerminate()
	RegisterEvent(msg *net.Message, c Channel) error
	UnregisterEvent(msg *net.Message, c Channel) error
	MetaObject(objectID uint32) (object.MetaObject, error)
	Terminate(objectID uint32) error
	Property(name value.Value) (value.Value, error)
	SetProperty(name value.Value, value value.Value) error
	Properties() ([]string, error)
	RegisterEventWithSignature(objectID uint32, actionID uint32, handler uint64, P3 string) (uint64, error)
	IsStatsEnabled() (bool, error)
	EnableStats(enabled bool) error
	Stats() (map[uint32]MethodStatistics, error)
	ClearStats() error
	IsTraceEnabled() (bool, error)
	EnableTrace(traced bool) error
	Tracer(msg *net.Message, from Channel) Channel
}

// ObjectSignalHelper provided to Object a companion object
type ObjectSignalHelper interface {
	SignalTraceObject(event EventTrace) error
}

// stubObject implements server.Actor.
type stubObject struct {
	impl      ObjectImplementor
	session   Session
	service   Service
	serviceID uint32
	signal    *signalHandler
	obj       Actor
}

// ObjectObject returns an object using ObjectImplementor
func ObjectObject(impl ObjectImplementor) Actor {
	var stb stubObject
	stb.impl = impl
	stb.signal = newSignalHandler()
	return &stb
}
func (p *stubObject) Activate(activation Activation) error {
	p.signal.Activate(activation)
	p.session = activation.Session
	err := p.impl.Activate(activation, p)
	if err != nil {
		p.signal.OnTerminate()
		return err
	}
	err = p.obj.Activate(activation)
	if err != nil {
		p.impl.OnTerminate()
		p.signal.OnTerminate()
		return err
	}
	return nil
}
func (p *stubObject) OnTerminate() {
	p.obj.OnTerminate()
	p.impl.OnTerminate()
	p.signal.OnTerminate()
}
func (p *stubObject) Receive(msg *net.Message, from Channel) error {
	from = p.impl.Tracer(msg, from)
	switch msg.Header.Action {
	case 0:
		return p.RegisterEvent(msg, from)
	case 1:
		return p.UnregisterEvent(msg, from)
	case 2:
		return p.MetaObject(msg, from)
	case 3:
		return p.Terminate(msg, from)
	case 5:
		return p.Property(msg, from)
	case 6:
		return p.SetProperty(msg, from)
	case 7:
		return p.Properties(msg, from)
	case 8:
		return p.RegisterEventWithSignature(msg, from)
	case 80:
		return p.IsStatsEnabled(msg, from)
	case 81:
		return p.EnableStats(msg, from)
	case 82:
		return p.Stats(msg, from)
	case 83:
		return p.ClearStats(msg, from)
	case 84:
		return p.IsTraceEnabled(msg, from)
	case 85:
		return p.EnableTrace(msg, from)
	default:
		return p.obj.Receive(msg, from)
	}
}
func (p *stubObject) onPropertyChange(name string, data []byte) error {
	switch name {
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}

func (p *stubObject) RegisterEvent(msg *net.Message, c Channel) error {
	return p.impl.RegisterEvent(msg, c)
}

func (p *stubObject) UnregisterEvent(msg *net.Message, c Channel) error {
	return p.impl.UnregisterEvent(msg, c)
}
func (p *stubObject) MetaObject(msg *net.Message, c Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read objectID: %s", err))
	}
	ret, callErr := p.impl.MetaObject(objectID)

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := object.WriteMetaObject(ret, &out)
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) Terminate(msg *net.Message, c Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read objectID: %s", err))
	}
	callErr := p.impl.Terminate(objectID)

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
func (p *stubObject) Property(msg *net.Message, c Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	name, err := value.NewValue(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read name: %s", err))
	}
	ret, callErr := p.impl.Property(name)

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := ret.Write(&out)
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) SetProperty(msg *net.Message, c Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	name, err := value.NewValue(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read name: %s", err))
	}
	value, err := value.NewValue(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read value: %s", err))
	}
	callErr := p.impl.SetProperty(name, value)

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
func (p *stubObject) Properties(msg *net.Message, c Channel) error {
	ret, callErr := p.impl.Properties()

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := func() error {
		err := basic.WriteUint32(uint32(len(ret)), &out)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range ret {
			err = basic.WriteString(v, &out)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}()
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) RegisterEventWithSignature(msg *net.Message, c Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read objectID: %s", err))
	}
	actionID, err := basic.ReadUint32(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read actionID: %s", err))
	}
	handler, err := basic.ReadUint64(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read handler: %s", err))
	}
	P3, err := basic.ReadString(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read P3: %s", err))
	}
	ret, callErr := p.impl.RegisterEventWithSignature(objectID, actionID, handler, P3)

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := basic.WriteUint64(ret, &out)
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) IsStatsEnabled(msg *net.Message, c Channel) error {
	ret, callErr := p.impl.IsStatsEnabled()

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := basic.WriteBool(ret, &out)
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) EnableStats(msg *net.Message, c Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	enabled, err := basic.ReadBool(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read enabled: %s", err))
	}
	callErr := p.impl.EnableStats(enabled)

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
func (p *stubObject) Stats(msg *net.Message, c Channel) error {
	ret, callErr := p.impl.Stats()

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := func() error {
		err := basic.WriteUint32(uint32(len(ret)), &out)
		if err != nil {
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range ret {
			err = basic.WriteUint32(k, &out)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = writeMethodStatistics(v, &out)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}()
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) ClearStats(msg *net.Message, c Channel) error {
	callErr := p.impl.ClearStats()

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
func (p *stubObject) IsTraceEnabled(msg *net.Message, c Channel) error {
	ret, callErr := p.impl.IsTraceEnabled()

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := basic.WriteBool(ret, &out)
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) EnableTrace(msg *net.Message, c Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	traced, err := basic.ReadBool(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read traced: %s", err))
	}
	callErr := p.impl.EnableTrace(traced)

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
func (p *stubObject) SignalTraceObject(event EventTrace) error {
	var buf bytes.Buffer
	if err := writeEventTrace(event, &buf); err != nil {
		return fmt.Errorf("serialize event: %s", err)
	}
	err := p.signal.UpdateSignal(86, buf.Bytes())

	if err != nil {
		return fmt.Errorf("update SignalTraceObject: %s", err)
	}
	return nil
}
func (p *stubObject) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "Object",
		Methods: map[uint32]object.MetaMethod{
			0: {
				Name:                "registerEvent",
				ParametersSignature: "(IIL)",
				ReturnSignature:     "L",
				Uid:                 0,
			},
			1: {
				Name:                "unregisterEvent",
				ParametersSignature: "(IIL)",
				ReturnSignature:     "v",
				Uid:                 1,
			},
			2: {
				Name:                "metaObject",
				ParametersSignature: "(I)",
				ReturnSignature:     "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>",
				Uid:                 2,
			},
			3: {
				Name:                "terminate",
				ParametersSignature: "(I)",
				ReturnSignature:     "v",
				Uid:                 3,
			},
			5: {
				Name:                "property",
				ParametersSignature: "(m)",
				ReturnSignature:     "m",
				Uid:                 5,
			},
			6: {
				Name:                "setProperty",
				ParametersSignature: "(mm)",
				ReturnSignature:     "v",
				Uid:                 6,
			},
			7: {
				Name:                "properties",
				ParametersSignature: "()",
				ReturnSignature:     "[s]",
				Uid:                 7,
			},
			8: {
				Name:                "registerEventWithSignature",
				ParametersSignature: "(IILs)",
				ReturnSignature:     "L",
				Uid:                 8,
			},
			80: {
				Name:                "isStatsEnabled",
				ParametersSignature: "()",
				ReturnSignature:     "b",
				Uid:                 80,
			},
			81: {
				Name:                "enableStats",
				ParametersSignature: "(b)",
				ReturnSignature:     "v",
				Uid:                 81,
			},
			82: {
				Name:                "stats",
				ParametersSignature: "()",
				ReturnSignature:     "{I(I(fff)<MinMaxSum,minValue,maxValue,cumulatedValue>(fff)<MinMaxSum,minValue,maxValue,cumulatedValue>(fff)<MinMaxSum,minValue,maxValue,cumulatedValue>)<MethodStatistics,count,wall,user,system>}",
				Uid:                 82,
			},
			83: {
				Name:                "clearStats",
				ParametersSignature: "()",
				ReturnSignature:     "v",
				Uid:                 83,
			},
			84: {
				Name:                "isTraceEnabled",
				ParametersSignature: "()",
				ReturnSignature:     "b",
				Uid:                 84,
			},
			85: {
				Name:                "enableTrace",
				ParametersSignature: "(b)",
				ReturnSignature:     "v",
				Uid:                 85,
			},
		},
		Properties: map[uint32]object.MetaProperty{},
		Signals: map[uint32]object.MetaSignal{86: {
			Name:      "traceObject",
			Signature: "(IiIm(ll)<timeval,tv_sec,tv_usec>llII)<EventTrace,id,kind,slotId,arguments,timestamp,userUsTime,systemUsTime,callerContext,calleeContext>",
			Uid:       86,
		}},
	}
}

// Constructor gives access to remote services
type Constructor struct {
	session Session
}

// Services gives access to the services constructor
func Services(s Session) Constructor {
	return Constructor{session: s}
}

// ServiceZero is the abstract interface of the service
type ServiceZero interface {
	// Authenticate calls the remote procedure
	Authenticate(capability map[string]value.Value) (map[string]value.Value, error)
}

// ServiceZeroProxy represents a proxy object to the service
type ServiceZeroProxy interface {
	Proxy
	ServiceZero
}

// proxyServiceZero implements ServiceZeroProxy
type proxyServiceZero struct {
	Proxy
}

// ServiceZero returns a proxy to a remote service
func (c Constructor) ServiceZero(closer func(error)) (ServiceZeroProxy, error) {
	proxy, err := c.session.Proxy("ServiceZero", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}

	err = proxy.OnDisconnect(closer)
	if err != nil {
		return nil, err
	}
	return &proxyServiceZero{proxy}, nil
}

// Authenticate calls the remote procedure
func (p *proxyServiceZero) Authenticate(capability map[string]value.Value) (map[string]value.Value, error) {
	var err error
	var ret map[string]value.Value
	var buf bytes.Buffer
	if err = func() error {
		err := basic.WriteUint32(uint32(len(capability)), &buf)
		if err != nil {
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range capability {
			err = basic.WriteString(k, &buf)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = v.Write(&buf)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return ret, fmt.Errorf("serialize capability: %s", err)
	}
	response, err := p.Call("authenticate", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call authenticate failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (m map[string]value.Value, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return m, fmt.Errorf("read map size: %s", err)
		}
		m = make(map[string]value.Value, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(resp)
			if err != nil {
				return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
			}
			v, err := value.NewValue(resp)
			if err != nil {
				return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse authenticate response: %s", err)
	}
	return ret, nil
}

// Object is the abstract interface of the service
type Object interface {
	// IsStatsEnabled calls the remote procedure
	IsStatsEnabled() (bool, error)
	// EnableStats calls the remote procedure
	EnableStats(enabled bool) error
	// Stats calls the remote procedure
	Stats() (map[uint32]MethodStatistics, error)
	// ClearStats calls the remote procedure
	ClearStats() error
	// IsTraceEnabled calls the remote procedure
	IsTraceEnabled() (bool, error)
	// EnableTrace calls the remote procedure
	EnableTrace(traced bool) error
	// SubscribeTraceObject subscribe to a remote signal
	SubscribeTraceObject() (unsubscribe func(), updates chan EventTrace, err error)
}

// ObjectProxy represents a proxy object to the service
type ObjectProxy interface {
	object.Object
	Proxy
	Object
}

// proxyObject implements ObjectProxy
type proxyObject struct {
	Proxy
}

// Object returns a proxy to a remote service
func (c Constructor) Object(closer func(error)) (ObjectProxy, error) {
	proxy, err := c.session.Proxy("Object", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}

	err = proxy.OnDisconnect(closer)
	if err != nil {
		return nil, err
	}
	return &proxyObject{proxy}, nil
}

// RegisterEvent calls the remote procedure
func (p *proxyObject) RegisterEvent(objectID uint32, actionID uint32, handler uint64) (uint64, error) {
	var err error
	var ret uint64
	var buf bytes.Buffer
	if err = basic.WriteUint32(objectID, &buf); err != nil {
		return ret, fmt.Errorf("serialize objectID: %s", err)
	}
	if err = basic.WriteUint32(actionID, &buf); err != nil {
		return ret, fmt.Errorf("serialize actionID: %s", err)
	}
	if err = basic.WriteUint64(handler, &buf); err != nil {
		return ret, fmt.Errorf("serialize handler: %s", err)
	}
	response, err := p.Call("registerEvent", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerEvent failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(resp)
	if err != nil {
		return ret, fmt.Errorf("parse registerEvent response: %s", err)
	}
	return ret, nil
}

// UnregisterEvent calls the remote procedure
func (p *proxyObject) UnregisterEvent(objectID uint32, actionID uint32, handler uint64) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteUint32(objectID, &buf); err != nil {
		return fmt.Errorf("serialize objectID: %s", err)
	}
	if err = basic.WriteUint32(actionID, &buf); err != nil {
		return fmt.Errorf("serialize actionID: %s", err)
	}
	if err = basic.WriteUint64(handler, &buf); err != nil {
		return fmt.Errorf("serialize handler: %s", err)
	}
	_, err = p.Call("unregisterEvent", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterEvent failed: %s", err)
	}
	return nil
}

// MetaObject calls the remote procedure
func (p *proxyObject) MetaObject(objectID uint32) (object.MetaObject, error) {
	var err error
	var ret object.MetaObject
	var buf bytes.Buffer
	if err = basic.WriteUint32(objectID, &buf); err != nil {
		return ret, fmt.Errorf("serialize objectID: %s", err)
	}
	response, err := p.Call("metaObject", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call metaObject failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = object.ReadMetaObject(resp)
	if err != nil {
		return ret, fmt.Errorf("parse metaObject response: %s", err)
	}
	return ret, nil
}

// Terminate calls the remote procedure
func (p *proxyObject) Terminate(objectID uint32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteUint32(objectID, &buf); err != nil {
		return fmt.Errorf("serialize objectID: %s", err)
	}
	_, err = p.Call("terminate", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call terminate failed: %s", err)
	}
	return nil
}

// Property calls the remote procedure
func (p *proxyObject) Property(name value.Value) (value.Value, error) {
	var err error
	var ret value.Value
	var buf bytes.Buffer
	if err = name.Write(&buf); err != nil {
		return ret, fmt.Errorf("serialize name: %s", err)
	}
	response, err := p.Call("property", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call property failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = value.NewValue(resp)
	if err != nil {
		return ret, fmt.Errorf("parse property response: %s", err)
	}
	return ret, nil
}

// SetProperty calls the remote procedure
func (p *proxyObject) SetProperty(name value.Value, value value.Value) error {
	var err error
	var buf bytes.Buffer
	if err = name.Write(&buf); err != nil {
		return fmt.Errorf("serialize name: %s", err)
	}
	if err = value.Write(&buf); err != nil {
		return fmt.Errorf("serialize value: %s", err)
	}
	_, err = p.Call("setProperty", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setProperty failed: %s", err)
	}
	return nil
}

// Properties calls the remote procedure
func (p *proxyObject) Properties() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.Call("properties", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call properties failed: %s", err)
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
		return ret, fmt.Errorf("parse properties response: %s", err)
	}
	return ret, nil
}

// RegisterEventWithSignature calls the remote procedure
func (p *proxyObject) RegisterEventWithSignature(objectID uint32, actionID uint32, handler uint64, P3 string) (uint64, error) {
	var err error
	var ret uint64
	var buf bytes.Buffer
	if err = basic.WriteUint32(objectID, &buf); err != nil {
		return ret, fmt.Errorf("serialize objectID: %s", err)
	}
	if err = basic.WriteUint32(actionID, &buf); err != nil {
		return ret, fmt.Errorf("serialize actionID: %s", err)
	}
	if err = basic.WriteUint64(handler, &buf); err != nil {
		return ret, fmt.Errorf("serialize handler: %s", err)
	}
	if err = basic.WriteString(P3, &buf); err != nil {
		return ret, fmt.Errorf("serialize P3: %s", err)
	}
	response, err := p.Call("registerEventWithSignature", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerEventWithSignature failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(resp)
	if err != nil {
		return ret, fmt.Errorf("parse registerEventWithSignature response: %s", err)
	}
	return ret, nil
}

// IsStatsEnabled calls the remote procedure
func (p *proxyObject) IsStatsEnabled() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.Call("isStatsEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isStatsEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse isStatsEnabled response: %s", err)
	}
	return ret, nil
}

// EnableStats calls the remote procedure
func (p *proxyObject) EnableStats(enabled bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(enabled, &buf); err != nil {
		return fmt.Errorf("serialize enabled: %s", err)
	}
	_, err = p.Call("enableStats", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableStats failed: %s", err)
	}
	return nil
}

// Stats calls the remote procedure
func (p *proxyObject) Stats() (map[uint32]MethodStatistics, error) {
	var err error
	var ret map[uint32]MethodStatistics
	var buf bytes.Buffer
	response, err := p.Call("stats", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call stats failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (m map[uint32]MethodStatistics, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return m, fmt.Errorf("read map size: %s", err)
		}
		m = make(map[uint32]MethodStatistics, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(resp)
			if err != nil {
				return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
			}
			v, err := readMethodStatistics(resp)
			if err != nil {
				return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse stats response: %s", err)
	}
	return ret, nil
}

// ClearStats calls the remote procedure
func (p *proxyObject) ClearStats() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("clearStats", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearStats failed: %s", err)
	}
	return nil
}

// IsTraceEnabled calls the remote procedure
func (p *proxyObject) IsTraceEnabled() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.Call("isTraceEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isTraceEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse isTraceEnabled response: %s", err)
	}
	return ret, nil
}

// EnableTrace calls the remote procedure
func (p *proxyObject) EnableTrace(traced bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(traced, &buf); err != nil {
		return fmt.Errorf("serialize traced: %s", err)
	}
	_, err = p.Call("enableTrace", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableTrace failed: %s", err)
	}
	return nil
}

// SubscribeTraceObject subscribe to a remote property
func (p *proxyObject) SubscribeTraceObject() (func(), chan EventTrace, error) {
	propertyID, err := p.SignalID("traceObject")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "traceObject", err)
	}

	handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("register event for %s: %s", "traceObject", err)
	}
	ch := make(chan EventTrace)
	cancel, chPay, err := p.SubscribeID(propertyID)
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
			e, err := readEventTrace(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()

	return func() {
		p.UnregisterEvent(p.ObjectID(), propertyID, handlerID)
		cancel()
	}, ch, nil
}

// MetaMethodParameter is serializable
type MetaMethodParameter struct {
	Name        string
	Description string
}

// readMetaMethodParameter unmarshalls MetaMethodParameter
func readMetaMethodParameter(r io.Reader) (s MetaMethodParameter, err error) {
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Description field: %s", err)
	}
	return s, nil
}

// writeMetaMethodParameter marshalls MetaMethodParameter
func writeMetaMethodParameter(s MetaMethodParameter, w io.Writer) (err error) {
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("write Description field: %s", err)
	}
	return nil
}

// MetaMethod is serializable
type MetaMethod struct {
	Uid                 uint32
	ReturnSignature     string
	Name                string
	ParametersSignature string
	Description         string
	Parameters          []MetaMethodParameter
	ReturnDescription   string
}

// readMetaMethod unmarshalls MetaMethod
func readMetaMethod(r io.Reader) (s MetaMethod, err error) {
	if s.Uid, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read Uid field: %s", err)
	}
	if s.ReturnSignature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read ReturnSignature field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	if s.ParametersSignature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read ParametersSignature field: %s", err)
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Description field: %s", err)
	}
	if s.Parameters, err = func() (b []MetaMethodParameter, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]MetaMethodParameter, size)
		for i := 0; i < int(size); i++ {
			b[i], err = readMetaMethodParameter(r)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("read Parameters field: %s", err)
	}
	if s.ReturnDescription, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read ReturnDescription field: %s", err)
	}
	return s, nil
}

// writeMetaMethod marshalls MetaMethod
func writeMetaMethod(s MetaMethod, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("write Uid field: %s", err)
	}
	if err := basic.WriteString(s.ReturnSignature, w); err != nil {
		return fmt.Errorf("write ReturnSignature field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	if err := basic.WriteString(s.ParametersSignature, w); err != nil {
		return fmt.Errorf("write ParametersSignature field: %s", err)
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("write Description field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Parameters)), w)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range s.Parameters {
			err = writeMetaMethodParameter(v, w)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("write Parameters field: %s", err)
	}
	if err := basic.WriteString(s.ReturnDescription, w); err != nil {
		return fmt.Errorf("write ReturnDescription field: %s", err)
	}
	return nil
}

// MetaSignal is serializable
type MetaSignal struct {
	Uid       uint32
	Name      string
	Signature string
}

// readMetaSignal unmarshalls MetaSignal
func readMetaSignal(r io.Reader) (s MetaSignal, err error) {
	if s.Uid, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read Uid field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	if s.Signature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Signature field: %s", err)
	}
	return s, nil
}

// writeMetaSignal marshalls MetaSignal
func writeMetaSignal(s MetaSignal, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("write Uid field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	if err := basic.WriteString(s.Signature, w); err != nil {
		return fmt.Errorf("write Signature field: %s", err)
	}
	return nil
}

// MetaProperty is serializable
type MetaProperty struct {
	Uid       uint32
	Name      string
	Signature string
}

// readMetaProperty unmarshalls MetaProperty
func readMetaProperty(r io.Reader) (s MetaProperty, err error) {
	if s.Uid, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read Uid field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	if s.Signature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Signature field: %s", err)
	}
	return s, nil
}

// writeMetaProperty marshalls MetaProperty
func writeMetaProperty(s MetaProperty, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("write Uid field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	if err := basic.WriteString(s.Signature, w); err != nil {
		return fmt.Errorf("write Signature field: %s", err)
	}
	return nil
}

// MetaObject is serializable
type MetaObject struct {
	Methods     map[uint32]MetaMethod
	Signals     map[uint32]MetaSignal
	Properties  map[uint32]MetaProperty
	Description string
}

// readMetaObject unmarshalls MetaObject
func readMetaObject(r io.Reader) (s MetaObject, err error) {
	if s.Methods, err = func() (m map[uint32]MetaMethod, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("read map size: %s", err)
		}
		m = make(map[uint32]MetaMethod, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(r)
			if err != nil {
				return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
			}
			v, err := readMetaMethod(r)
			if err != nil {
				return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("read Methods field: %s", err)
	}
	if s.Signals, err = func() (m map[uint32]MetaSignal, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("read map size: %s", err)
		}
		m = make(map[uint32]MetaSignal, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(r)
			if err != nil {
				return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
			}
			v, err := readMetaSignal(r)
			if err != nil {
				return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("read Signals field: %s", err)
	}
	if s.Properties, err = func() (m map[uint32]MetaProperty, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("read map size: %s", err)
		}
		m = make(map[uint32]MetaProperty, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(r)
			if err != nil {
				return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
			}
			v, err := readMetaProperty(r)
			if err != nil {
				return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("read Properties field: %s", err)
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Description field: %s", err)
	}
	return s, nil
}

// writeMetaObject marshalls MetaObject
func writeMetaObject(s MetaObject, w io.Writer) (err error) {
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Methods)), w)
		if err != nil {
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range s.Methods {
			err = basic.WriteUint32(k, w)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = writeMetaMethod(v, w)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("write Methods field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Signals)), w)
		if err != nil {
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range s.Signals {
			err = basic.WriteUint32(k, w)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = writeMetaSignal(v, w)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("write Signals field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Properties)), w)
		if err != nil {
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range s.Properties {
			err = basic.WriteUint32(k, w)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = writeMetaProperty(v, w)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("write Properties field: %s", err)
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("write Description field: %s", err)
	}
	return nil
}

// MinMaxSum is serializable
type MinMaxSum struct {
	MinValue       float32
	MaxValue       float32
	CumulatedValue float32
}

// readMinMaxSum unmarshalls MinMaxSum
func readMinMaxSum(r io.Reader) (s MinMaxSum, err error) {
	if s.MinValue, err = basic.ReadFloat32(r); err != nil {
		return s, fmt.Errorf("read MinValue field: %s", err)
	}
	if s.MaxValue, err = basic.ReadFloat32(r); err != nil {
		return s, fmt.Errorf("read MaxValue field: %s", err)
	}
	if s.CumulatedValue, err = basic.ReadFloat32(r); err != nil {
		return s, fmt.Errorf("read CumulatedValue field: %s", err)
	}
	return s, nil
}

// writeMinMaxSum marshalls MinMaxSum
func writeMinMaxSum(s MinMaxSum, w io.Writer) (err error) {
	if err := basic.WriteFloat32(s.MinValue, w); err != nil {
		return fmt.Errorf("write MinValue field: %s", err)
	}
	if err := basic.WriteFloat32(s.MaxValue, w); err != nil {
		return fmt.Errorf("write MaxValue field: %s", err)
	}
	if err := basic.WriteFloat32(s.CumulatedValue, w); err != nil {
		return fmt.Errorf("write CumulatedValue field: %s", err)
	}
	return nil
}

// MethodStatistics is serializable
type MethodStatistics struct {
	Count  uint32
	Wall   MinMaxSum
	User   MinMaxSum
	System MinMaxSum
}

// readMethodStatistics unmarshalls MethodStatistics
func readMethodStatistics(r io.Reader) (s MethodStatistics, err error) {
	if s.Count, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read Count field: %s", err)
	}
	if s.Wall, err = readMinMaxSum(r); err != nil {
		return s, fmt.Errorf("read Wall field: %s", err)
	}
	if s.User, err = readMinMaxSum(r); err != nil {
		return s, fmt.Errorf("read User field: %s", err)
	}
	if s.System, err = readMinMaxSum(r); err != nil {
		return s, fmt.Errorf("read System field: %s", err)
	}
	return s, nil
}

// writeMethodStatistics marshalls MethodStatistics
func writeMethodStatistics(s MethodStatistics, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Count, w); err != nil {
		return fmt.Errorf("write Count field: %s", err)
	}
	if err := writeMinMaxSum(s.Wall, w); err != nil {
		return fmt.Errorf("write Wall field: %s", err)
	}
	if err := writeMinMaxSum(s.User, w); err != nil {
		return fmt.Errorf("write User field: %s", err)
	}
	if err := writeMinMaxSum(s.System, w); err != nil {
		return fmt.Errorf("write System field: %s", err)
	}
	return nil
}

// Timeval is serializable
type Timeval struct {
	Tv_sec  int64
	Tv_usec int64
}

// readTimeval unmarshalls Timeval
func readTimeval(r io.Reader) (s Timeval, err error) {
	if s.Tv_sec, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("read Tv_sec field: %s", err)
	}
	if s.Tv_usec, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("read Tv_usec field: %s", err)
	}
	return s, nil
}

// writeTimeval marshalls Timeval
func writeTimeval(s Timeval, w io.Writer) (err error) {
	if err := basic.WriteInt64(s.Tv_sec, w); err != nil {
		return fmt.Errorf("write Tv_sec field: %s", err)
	}
	if err := basic.WriteInt64(s.Tv_usec, w); err != nil {
		return fmt.Errorf("write Tv_usec field: %s", err)
	}
	return nil
}

// EventTrace is serializable
type EventTrace struct {
	Id            uint32
	Kind          int32
	SlotId        uint32
	Arguments     value.Value
	Timestamp     Timeval
	UserUsTime    int64
	SystemUsTime  int64
	CallerContext uint32
	CalleeContext uint32
}

// readEventTrace unmarshalls EventTrace
func readEventTrace(r io.Reader) (s EventTrace, err error) {
	if s.Id, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read Id field: %s", err)
	}
	if s.Kind, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("read Kind field: %s", err)
	}
	if s.SlotId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read SlotId field: %s", err)
	}
	if s.Arguments, err = value.NewValue(r); err != nil {
		return s, fmt.Errorf("read Arguments field: %s", err)
	}
	if s.Timestamp, err = readTimeval(r); err != nil {
		return s, fmt.Errorf("read Timestamp field: %s", err)
	}
	if s.UserUsTime, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("read UserUsTime field: %s", err)
	}
	if s.SystemUsTime, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("read SystemUsTime field: %s", err)
	}
	if s.CallerContext, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read CallerContext field: %s", err)
	}
	if s.CalleeContext, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read CalleeContext field: %s", err)
	}
	return s, nil
}

// writeEventTrace marshalls EventTrace
func writeEventTrace(s EventTrace, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Id, w); err != nil {
		return fmt.Errorf("write Id field: %s", err)
	}
	if err := basic.WriteInt32(s.Kind, w); err != nil {
		return fmt.Errorf("write Kind field: %s", err)
	}
	if err := basic.WriteUint32(s.SlotId, w); err != nil {
		return fmt.Errorf("write SlotId field: %s", err)
	}
	if err := s.Arguments.Write(w); err != nil {
		return fmt.Errorf("write Arguments field: %s", err)
	}
	if err := writeTimeval(s.Timestamp, w); err != nil {
		return fmt.Errorf("write Timestamp field: %s", err)
	}
	if err := basic.WriteInt64(s.UserUsTime, w); err != nil {
		return fmt.Errorf("write UserUsTime field: %s", err)
	}
	if err := basic.WriteInt64(s.SystemUsTime, w); err != nil {
		return fmt.Errorf("write SystemUsTime field: %s", err)
	}
	if err := basic.WriteUint32(s.CallerContext, w); err != nil {
		return fmt.Errorf("write CallerContext field: %s", err)
	}
	if err := basic.WriteUint32(s.CalleeContext, w); err != nil {
		return fmt.Errorf("write CalleeContext field: %s", err)
	}
	return nil
}
