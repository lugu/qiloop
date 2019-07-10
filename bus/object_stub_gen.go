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
func (p *stubServiceZero) Receive(msg *net.Message, from *Channel) error {
	switch msg.Header.Action {
	case uint32(0x8):
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
func (p *stubServiceZero) Authenticate(msg *net.Message, c *Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	capability, err := func() (m map[string]value.Value, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]value.Value, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := value.NewValue(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read capability: %s", err))
	}
	ret, callErr := p.impl.Authenticate(capability)
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := func() error {
		err := basic.WriteUint32(uint32(len(ret)), &out)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range ret {
			err = basic.WriteString(k, &out)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = v.Write(&out)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
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
		Methods: map[uint32]object.MetaMethod{uint32(0x8): {
			Name:                "authenticate",
			ParametersSignature: "({sm})",
			ReturnSignature:     "{sm}",
			Uid:                 uint32(0x8),
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
	RegisterEvent(objectID uint32, actionID uint32, handler uint64) (uint64, error)
	UnregisterEvent(objectID uint32, actionID uint32, handler uint64) error
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
	stb.signal = NewSignalHandler()
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
func (p *stubObject) Receive(msg *net.Message, from *Channel) error {
	switch msg.Header.Action {
	case uint32(0x0):
		return p.signal.RegisterEvent(msg, from)
	case uint32(0x1):
		return p.signal.UnregisterEvent(msg, from)
	case uint32(0x2):
		return p.MetaObject(msg, from)
	case uint32(0x3):
		return p.Terminate(msg, from)
	case uint32(0x5):
		return p.Property(msg, from)
	case uint32(0x6):
		return p.SetProperty(msg, from)
	case uint32(0x7):
		return p.Properties(msg, from)
	case uint32(0x8):
		return p.RegisterEventWithSignature(msg, from)
	case uint32(0x50):
		return p.IsStatsEnabled(msg, from)
	case uint32(0x51):
		return p.EnableStats(msg, from)
	case uint32(0x52):
		return p.Stats(msg, from)
	case uint32(0x53):
		return p.ClearStats(msg, from)
	case uint32(0x54):
		return p.IsTraceEnabled(msg, from)
	case uint32(0x55):
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
func (p *stubObject) RegisterEvent(msg *net.Message, c *Channel) error {
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
	ret, callErr := p.impl.RegisterEvent(objectID, actionID, handler)
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
func (p *stubObject) UnregisterEvent(msg *net.Message, c *Channel) error {
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
	callErr := p.impl.UnregisterEvent(objectID, actionID, handler)
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) MetaObject(msg *net.Message, c *Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read objectID: %s", err))
	}
	ret, callErr := p.impl.MetaObject(objectID)
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
func (p *stubObject) Terminate(msg *net.Message, c *Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read objectID: %s", err))
	}
	callErr := p.impl.Terminate(objectID)
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) Property(msg *net.Message, c *Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	name, err := value.NewValue(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read name: %s", err))
	}
	ret, callErr := p.impl.Property(name)
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
func (p *stubObject) SetProperty(msg *net.Message, c *Channel) error {
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
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) Properties(msg *net.Message, c *Channel) error {
	ret, callErr := p.impl.Properties()
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := func() error {
		err := basic.WriteUint32(uint32(len(ret)), &out)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range ret {
			err = basic.WriteString(v, &out)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}()
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) RegisterEventWithSignature(msg *net.Message, c *Channel) error {
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
func (p *stubObject) IsStatsEnabled(msg *net.Message, c *Channel) error {
	ret, callErr := p.impl.IsStatsEnabled()
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
func (p *stubObject) EnableStats(msg *net.Message, c *Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	enabled, err := basic.ReadBool(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read enabled: %s", err))
	}
	callErr := p.impl.EnableStats(enabled)
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) Stats(msg *net.Message, c *Channel) error {
	ret, callErr := p.impl.Stats()
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := func() error {
		err := basic.WriteUint32(uint32(len(ret)), &out)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range ret {
			err = basic.WriteUint32(k, &out)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = writeMethodStatistics(v, &out)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}()
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) ClearStats(msg *net.Message, c *Channel) error {
	callErr := p.impl.ClearStats()
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) IsTraceEnabled(msg *net.Message, c *Channel) error {
	ret, callErr := p.impl.IsTraceEnabled()
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
func (p *stubObject) EnableTrace(msg *net.Message, c *Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	traced, err := basic.ReadBool(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read traced: %s", err))
	}
	callErr := p.impl.EnableTrace(traced)
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	return c.SendReply(msg, out.Bytes())
}
func (p *stubObject) SignalTraceObject(event EventTrace) error {
	var buf bytes.Buffer
	if err := writeEventTrace(event, &buf); err != nil {
		return fmt.Errorf("failed to serialize event: %s", err)
	}
	err := p.signal.UpdateSignal(uint32(0x56), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalTraceObject: %s", err)
	}
	return nil
}
func (p *stubObject) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "Object",
		Methods: map[uint32]object.MetaMethod{
			uint32(0x0): {
				Name:                "registerEvent",
				ParametersSignature: "(IIL)",
				ReturnSignature:     "L",
				Uid:                 uint32(0x0),
			},
			uint32(0x1): {
				Name:                "unregisterEvent",
				ParametersSignature: "(IIL)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x1),
			},
			uint32(0x2): {
				Name:                "metaObject",
				ParametersSignature: "(I)",
				ReturnSignature:     "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>",
				Uid:                 uint32(0x2),
			},
			uint32(0x3): {
				Name:                "terminate",
				ParametersSignature: "(I)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x3),
			},
			uint32(0x5): {
				Name:                "property",
				ParametersSignature: "(m)",
				ReturnSignature:     "m",
				Uid:                 uint32(0x5),
			},
			uint32(0x50): {
				Name:                "isStatsEnabled",
				ParametersSignature: "()",
				ReturnSignature:     "b",
				Uid:                 uint32(0x50),
			},
			uint32(0x51): {
				Name:                "enableStats",
				ParametersSignature: "(b)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x51),
			},
			uint32(0x52): {
				Name:                "stats",
				ParametersSignature: "()",
				ReturnSignature:     "{I(I(fff)<MinMaxSum,minValue,maxValue,cumulatedValue>(fff)<MinMaxSum,minValue,maxValue,cumulatedValue>(fff)<MinMaxSum,minValue,maxValue,cumulatedValue>)<MethodStatistics,count,wall,user,system>}",
				Uid:                 uint32(0x52),
			},
			uint32(0x53): {
				Name:                "clearStats",
				ParametersSignature: "()",
				ReturnSignature:     "v",
				Uid:                 uint32(0x53),
			},
			uint32(0x54): {
				Name:                "isTraceEnabled",
				ParametersSignature: "()",
				ReturnSignature:     "b",
				Uid:                 uint32(0x54),
			},
			uint32(0x55): {
				Name:                "enableTrace",
				ParametersSignature: "(b)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x55),
			},
			uint32(0x6): {
				Name:                "setProperty",
				ParametersSignature: "(mm)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x6),
			},
			uint32(0x7): {
				Name:                "properties",
				ParametersSignature: "()",
				ReturnSignature:     "[s]",
				Uid:                 uint32(0x7),
			},
			uint32(0x8): {
				Name:                "registerEventWithSignature",
				ParametersSignature: "(IILs)",
				ReturnSignature:     "L",
				Uid:                 uint32(0x8),
			},
		},
		Properties: map[uint32]object.MetaProperty{},
		Signals: map[uint32]object.MetaSignal{uint32(0x56): {
			Name:      "traceObject",
			Signature: "(IiIm(ll)<timeval,tv_sec,tv_usec>llII)<EventTrace,id,kind,slotId,arguments,timestamp,userUsTime,systemUsTime,callerContext,calleeContext>",
			Uid:       uint32(0x56),
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
func (s Constructor) ServiceZero() (ServiceZeroProxy, error) {
	proxy, err := s.session.Proxy("ServiceZero", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
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
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range capability {
			err = basic.WriteString(k, &buf)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = v.Write(&buf)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return ret, fmt.Errorf("failed to serialize capability: %s", err)
	}
	response, err := p.Call("authenticate", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call authenticate failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (m map[string]value.Value, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]value.Value, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(resp)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := value.NewValue(resp)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse authenticate response: %s", err)
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
func (s Constructor) Object() (ObjectProxy, error) {
	proxy, err := s.session.Proxy("Object", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &proxyObject{proxy}, nil
}

// RegisterEvent calls the remote procedure
func (p *proxyObject) RegisterEvent(objectID uint32, actionID uint32, handler uint64) (uint64, error) {
	var err error
	var ret uint64
	var buf bytes.Buffer
	if err = basic.WriteUint32(objectID, &buf); err != nil {
		return ret, fmt.Errorf("failed to serialize objectID: %s", err)
	}
	if err = basic.WriteUint32(actionID, &buf); err != nil {
		return ret, fmt.Errorf("failed to serialize actionID: %s", err)
	}
	if err = basic.WriteUint64(handler, &buf); err != nil {
		return ret, fmt.Errorf("failed to serialize handler: %s", err)
	}
	response, err := p.Call("registerEvent", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerEvent failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(resp)
	if err != nil {
		return ret, fmt.Errorf("failed to parse registerEvent response: %s", err)
	}
	return ret, nil
}

// UnregisterEvent calls the remote procedure
func (p *proxyObject) UnregisterEvent(objectID uint32, actionID uint32, handler uint64) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteUint32(objectID, &buf); err != nil {
		return fmt.Errorf("failed to serialize objectID: %s", err)
	}
	if err = basic.WriteUint32(actionID, &buf); err != nil {
		return fmt.Errorf("failed to serialize actionID: %s", err)
	}
	if err = basic.WriteUint64(handler, &buf); err != nil {
		return fmt.Errorf("failed to serialize handler: %s", err)
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
		return ret, fmt.Errorf("failed to serialize objectID: %s", err)
	}
	response, err := p.Call("metaObject", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call metaObject failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = object.ReadMetaObject(resp)
	if err != nil {
		return ret, fmt.Errorf("failed to parse metaObject response: %s", err)
	}
	return ret, nil
}

// Terminate calls the remote procedure
func (p *proxyObject) Terminate(objectID uint32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteUint32(objectID, &buf); err != nil {
		return fmt.Errorf("failed to serialize objectID: %s", err)
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
		return ret, fmt.Errorf("failed to serialize name: %s", err)
	}
	response, err := p.Call("property", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call property failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = value.NewValue(resp)
	if err != nil {
		return ret, fmt.Errorf("failed to parse property response: %s", err)
	}
	return ret, nil
}

// SetProperty calls the remote procedure
func (p *proxyObject) SetProperty(name value.Value, value value.Value) error {
	var err error
	var buf bytes.Buffer
	if err = name.Write(&buf); err != nil {
		return fmt.Errorf("failed to serialize name: %s", err)
	}
	if err = value.Write(&buf); err != nil {
		return fmt.Errorf("failed to serialize value: %s", err)
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
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse properties response: %s", err)
	}
	return ret, nil
}

// RegisterEventWithSignature calls the remote procedure
func (p *proxyObject) RegisterEventWithSignature(objectID uint32, actionID uint32, handler uint64, P3 string) (uint64, error) {
	var err error
	var ret uint64
	var buf bytes.Buffer
	if err = basic.WriteUint32(objectID, &buf); err != nil {
		return ret, fmt.Errorf("failed to serialize objectID: %s", err)
	}
	if err = basic.WriteUint32(actionID, &buf); err != nil {
		return ret, fmt.Errorf("failed to serialize actionID: %s", err)
	}
	if err = basic.WriteUint64(handler, &buf); err != nil {
		return ret, fmt.Errorf("failed to serialize handler: %s", err)
	}
	if err = basic.WriteString(P3, &buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P3: %s", err)
	}
	response, err := p.Call("registerEventWithSignature", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerEventWithSignature failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(resp)
	if err != nil {
		return ret, fmt.Errorf("failed to parse registerEventWithSignature response: %s", err)
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
		return ret, fmt.Errorf("failed to parse isStatsEnabled response: %s", err)
	}
	return ret, nil
}

// EnableStats calls the remote procedure
func (p *proxyObject) EnableStats(enabled bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(enabled, &buf); err != nil {
		return fmt.Errorf("failed to serialize enabled: %s", err)
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
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[uint32]MethodStatistics, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(resp)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := readMethodStatistics(resp)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse stats response: %s", err)
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
		return ret, fmt.Errorf("failed to parse isTraceEnabled response: %s", err)
	}
	return ret, nil
}

// EnableTrace calls the remote procedure
func (p *proxyObject) EnableTrace(traced bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(traced, &buf); err != nil {
		return fmt.Errorf("failed to serialize traced: %s", err)
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
		return nil, nil, fmt.Errorf("failed to register event for %s: %s", "traceObject", err)
	}
	ch := make(chan EventTrace)
	cancel, chPay, err := p.SubscribeID(propertyID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to request property: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				p.UnregisterEvent(p.ObjectID(), propertyID, handlerID)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := readEventTrace(buf)
			if err != nil {
				log.Printf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// MetaMethodParameter is serializable
type MetaMethodParameter struct {
	Name        string
	Description string
}

// readMetaMethodParameter unmarshalls MetaMethodParameter
func readMetaMethodParameter(r io.Reader) (s MetaMethodParameter, err error) {
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Description field: " + err.Error())
	}
	return s, nil
}

// writeMetaMethodParameter marshalls MetaMethodParameter
func writeMetaMethodParameter(s MetaMethodParameter, w io.Writer) (err error) {
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("failed to write Description field: " + err.Error())
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
		return s, fmt.Errorf("failed to read Uid field: " + err.Error())
	}
	if s.ReturnSignature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read ReturnSignature field: " + err.Error())
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	if s.ParametersSignature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read ParametersSignature field: " + err.Error())
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Description field: " + err.Error())
	}
	if s.Parameters, err = func() (b []MetaMethodParameter, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]MetaMethodParameter, size)
		for i := 0; i < int(size); i++ {
			b[i], err = readMetaMethodParameter(r)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Parameters field: " + err.Error())
	}
	if s.ReturnDescription, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read ReturnDescription field: " + err.Error())
	}
	return s, nil
}

// writeMetaMethod marshalls MetaMethod
func writeMetaMethod(s MetaMethod, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("failed to write Uid field: " + err.Error())
	}
	if err := basic.WriteString(s.ReturnSignature, w); err != nil {
		return fmt.Errorf("failed to write ReturnSignature field: " + err.Error())
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
	}
	if err := basic.WriteString(s.ParametersSignature, w); err != nil {
		return fmt.Errorf("failed to write ParametersSignature field: " + err.Error())
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("failed to write Description field: " + err.Error())
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Parameters)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.Parameters {
			err = writeMetaMethodParameter(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Parameters field: " + err.Error())
	}
	if err := basic.WriteString(s.ReturnDescription, w); err != nil {
		return fmt.Errorf("failed to write ReturnDescription field: " + err.Error())
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
		return s, fmt.Errorf("failed to read Uid field: " + err.Error())
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	if s.Signature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Signature field: " + err.Error())
	}
	return s, nil
}

// writeMetaSignal marshalls MetaSignal
func writeMetaSignal(s MetaSignal, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("failed to write Uid field: " + err.Error())
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
	}
	if err := basic.WriteString(s.Signature, w); err != nil {
		return fmt.Errorf("failed to write Signature field: " + err.Error())
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
		return s, fmt.Errorf("failed to read Uid field: " + err.Error())
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	if s.Signature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Signature field: " + err.Error())
	}
	return s, nil
}

// writeMetaProperty marshalls MetaProperty
func writeMetaProperty(s MetaProperty, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("failed to write Uid field: " + err.Error())
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
	}
	if err := basic.WriteString(s.Signature, w); err != nil {
		return fmt.Errorf("failed to write Signature field: " + err.Error())
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
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[uint32]MetaMethod, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := readMetaMethod(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Methods field: " + err.Error())
	}
	if s.Signals, err = func() (m map[uint32]MetaSignal, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[uint32]MetaSignal, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := readMetaSignal(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Signals field: " + err.Error())
	}
	if s.Properties, err = func() (m map[uint32]MetaProperty, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[uint32]MetaProperty, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := readMetaProperty(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Properties field: " + err.Error())
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Description field: " + err.Error())
	}
	return s, nil
}

// writeMetaObject marshalls MetaObject
func writeMetaObject(s MetaObject, w io.Writer) (err error) {
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Methods)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.Methods {
			err = basic.WriteUint32(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = writeMetaMethod(v, w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Methods field: " + err.Error())
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Signals)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.Signals {
			err = basic.WriteUint32(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = writeMetaSignal(v, w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Signals field: " + err.Error())
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Properties)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.Properties {
			err = basic.WriteUint32(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = writeMetaProperty(v, w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Properties field: " + err.Error())
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("failed to write Description field: " + err.Error())
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
		return s, fmt.Errorf("failed to read MinValue field: " + err.Error())
	}
	if s.MaxValue, err = basic.ReadFloat32(r); err != nil {
		return s, fmt.Errorf("failed to read MaxValue field: " + err.Error())
	}
	if s.CumulatedValue, err = basic.ReadFloat32(r); err != nil {
		return s, fmt.Errorf("failed to read CumulatedValue field: " + err.Error())
	}
	return s, nil
}

// writeMinMaxSum marshalls MinMaxSum
func writeMinMaxSum(s MinMaxSum, w io.Writer) (err error) {
	if err := basic.WriteFloat32(s.MinValue, w); err != nil {
		return fmt.Errorf("failed to write MinValue field: " + err.Error())
	}
	if err := basic.WriteFloat32(s.MaxValue, w); err != nil {
		return fmt.Errorf("failed to write MaxValue field: " + err.Error())
	}
	if err := basic.WriteFloat32(s.CumulatedValue, w); err != nil {
		return fmt.Errorf("failed to write CumulatedValue field: " + err.Error())
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
		return s, fmt.Errorf("failed to read Count field: " + err.Error())
	}
	if s.Wall, err = readMinMaxSum(r); err != nil {
		return s, fmt.Errorf("failed to read Wall field: " + err.Error())
	}
	if s.User, err = readMinMaxSum(r); err != nil {
		return s, fmt.Errorf("failed to read User field: " + err.Error())
	}
	if s.System, err = readMinMaxSum(r); err != nil {
		return s, fmt.Errorf("failed to read System field: " + err.Error())
	}
	return s, nil
}

// writeMethodStatistics marshalls MethodStatistics
func writeMethodStatistics(s MethodStatistics, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Count, w); err != nil {
		return fmt.Errorf("failed to write Count field: " + err.Error())
	}
	if err := writeMinMaxSum(s.Wall, w); err != nil {
		return fmt.Errorf("failed to write Wall field: " + err.Error())
	}
	if err := writeMinMaxSum(s.User, w); err != nil {
		return fmt.Errorf("failed to write User field: " + err.Error())
	}
	if err := writeMinMaxSum(s.System, w); err != nil {
		return fmt.Errorf("failed to write System field: " + err.Error())
	}
	return nil
}

// Timeval is serializable
type Timeval struct {
	Tvsec  int64
	Tvusec int64
}

// readTimeval unmarshalls Timeval
func readTimeval(r io.Reader) (s Timeval, err error) {
	if s.Tvsec, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read Tvsec field: " + err.Error())
	}
	if s.Tvusec, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read Tvusec field: " + err.Error())
	}
	return s, nil
}

// writeTimeval marshalls Timeval
func writeTimeval(s Timeval, w io.Writer) (err error) {
	if err := basic.WriteInt64(s.Tvsec, w); err != nil {
		return fmt.Errorf("failed to write Tvsec field: " + err.Error())
	}
	if err := basic.WriteInt64(s.Tvusec, w); err != nil {
		return fmt.Errorf("failed to write Tvusec field: " + err.Error())
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
		return s, fmt.Errorf("failed to read Id field: " + err.Error())
	}
	if s.Kind, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("failed to read Kind field: " + err.Error())
	}
	if s.SlotId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read SlotId field: " + err.Error())
	}
	if s.Arguments, err = value.NewValue(r); err != nil {
		return s, fmt.Errorf("failed to read Arguments field: " + err.Error())
	}
	if s.Timestamp, err = readTimeval(r); err != nil {
		return s, fmt.Errorf("failed to read Timestamp field: " + err.Error())
	}
	if s.UserUsTime, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read UserUsTime field: " + err.Error())
	}
	if s.SystemUsTime, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read SystemUsTime field: " + err.Error())
	}
	if s.CallerContext, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read CallerContext field: " + err.Error())
	}
	if s.CalleeContext, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read CalleeContext field: " + err.Error())
	}
	return s, nil
}

// writeEventTrace marshalls EventTrace
func writeEventTrace(s EventTrace, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Id, w); err != nil {
		return fmt.Errorf("failed to write Id field: " + err.Error())
	}
	if err := basic.WriteInt32(s.Kind, w); err != nil {
		return fmt.Errorf("failed to write Kind field: " + err.Error())
	}
	if err := basic.WriteUint32(s.SlotId, w); err != nil {
		return fmt.Errorf("failed to write SlotId field: " + err.Error())
	}
	if err := s.Arguments.Write(w); err != nil {
		return fmt.Errorf("failed to write Arguments field: " + err.Error())
	}
	if err := writeTimeval(s.Timestamp, w); err != nil {
		return fmt.Errorf("failed to write Timestamp field: " + err.Error())
	}
	if err := basic.WriteInt64(s.UserUsTime, w); err != nil {
		return fmt.Errorf("failed to write UserUsTime field: " + err.Error())
	}
	if err := basic.WriteInt64(s.SystemUsTime, w); err != nil {
		return fmt.Errorf("failed to write SystemUsTime field: " + err.Error())
	}
	if err := basic.WriteUint32(s.CallerContext, w); err != nil {
		return fmt.Errorf("failed to write CallerContext field: " + err.Error())
	}
	if err := basic.WriteUint32(s.CalleeContext, w); err != nil {
		return fmt.Errorf("failed to write CalleeContext field: " + err.Error())
	}
	return nil
}
