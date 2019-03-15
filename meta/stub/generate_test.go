// Package stub_test contains a generated stub
// File generated. DO NOT EDIT.

package stub_test

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	net "github.com/lugu/qiloop/bus/net"
	server "github.com/lugu/qiloop/bus/server"
	generic "github.com/lugu/qiloop/bus/server/generic"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
	"io"
	"log"
)

// ObjectImplementor interface of the service implementation
type ObjectImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals an properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation server.Activation, helper ObjectSignalHelper) error
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

// stubObject implements server.ServerObject.
type stubObject struct {
	obj     generic.Object
	impl    ObjectImplementor
	session bus.Session
}

// ObjectObject returns an object using ObjectImplementor
func ObjectObject(impl ObjectImplementor) server.ServerObject {
	var stb stubObject
	stb.impl = impl
	stb.obj = generic.NewObject(stb.metaObject())
	stb.obj.Wrap(uint32(0x0), stb.RegisterEvent)
	stb.obj.Wrap(uint32(0x1), stb.UnregisterEvent)
	stb.obj.Wrap(uint32(0x2), stb.MetaObject)
	stb.obj.Wrap(uint32(0x3), stb.Terminate)
	stb.obj.Wrap(uint32(0x5), stb.Property)
	stb.obj.Wrap(uint32(0x6), stb.SetProperty)
	stb.obj.Wrap(uint32(0x7), stb.Properties)
	stb.obj.Wrap(uint32(0x8), stb.RegisterEventWithSignature)
	stb.obj.Wrap(uint32(0x50), stb.IsStatsEnabled)
	stb.obj.Wrap(uint32(0x51), stb.EnableStats)
	stb.obj.Wrap(uint32(0x52), stb.Stats)
	stb.obj.Wrap(uint32(0x53), stb.ClearStats)
	stb.obj.Wrap(uint32(0x54), stb.IsTraceEnabled)
	stb.obj.Wrap(uint32(0x55), stb.EnableTrace)
	return &stb
}
func (p *stubObject) Activate(activation server.Activation) error {
	p.session = activation.Session
	p.obj.Activate(activation)
	return p.impl.Activate(activation, p)
}
func (p *stubObject) OnTerminate() {
	p.impl.OnTerminate()
	p.obj.OnTerminate()
}
func (p *stubObject) Receive(msg *net.Message, from *server.Context) error {
	return p.obj.Receive(msg, from)
}
func (p *stubObject) RegisterEvent(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read objectID: %s", err)
	}
	actionID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read actionID: %s", err)
	}
	handler, err := basic.ReadUint64(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read handler: %s", err)
	}
	ret, callErr := p.impl.RegisterEvent(objectID, actionID, handler)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteUint64(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubObject) UnregisterEvent(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read objectID: %s", err)
	}
	actionID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read actionID: %s", err)
	}
	handler, err := basic.ReadUint64(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read handler: %s", err)
	}
	callErr := p.impl.UnregisterEvent(objectID, actionID, handler)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubObject) MetaObject(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read objectID: %s", err)
	}
	ret, callErr := p.impl.MetaObject(objectID)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := object.WriteMetaObject(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubObject) Terminate(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read objectID: %s", err)
	}
	callErr := p.impl.Terminate(objectID)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubObject) Property(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	name, err := value.NewValue(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read name: %s", err)
	}
	ret, callErr := p.impl.Property(name)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := ret.Write(&out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubObject) SetProperty(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	name, err := value.NewValue(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read name: %s", err)
	}
	value, err := value.NewValue(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read value: %s", err)
	}
	callErr := p.impl.SetProperty(name, value)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubObject) Properties(payload []byte) ([]byte, error) {
	ret, callErr := p.impl.Properties()
	if callErr != nil {
		return nil, callErr
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
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubObject) RegisterEventWithSignature(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read objectID: %s", err)
	}
	actionID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read actionID: %s", err)
	}
	handler, err := basic.ReadUint64(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read handler: %s", err)
	}
	P3, err := basic.ReadString(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P3: %s", err)
	}
	ret, callErr := p.impl.RegisterEventWithSignature(objectID, actionID, handler, P3)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteUint64(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubObject) IsStatsEnabled(payload []byte) ([]byte, error) {
	ret, callErr := p.impl.IsStatsEnabled()
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteBool(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubObject) EnableStats(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	enabled, err := basic.ReadBool(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read enabled: %s", err)
	}
	callErr := p.impl.EnableStats(enabled)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubObject) Stats(payload []byte) ([]byte, error) {
	ret, callErr := p.impl.Stats()
	if callErr != nil {
		return nil, callErr
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
			err = WriteMethodStatistics(v, &out)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}()
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubObject) ClearStats(payload []byte) ([]byte, error) {
	callErr := p.impl.ClearStats()
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubObject) IsTraceEnabled(payload []byte) ([]byte, error) {
	ret, callErr := p.impl.IsTraceEnabled()
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteBool(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubObject) EnableTrace(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	traced, err := basic.ReadBool(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read traced: %s", err)
	}
	callErr := p.impl.EnableTrace(traced)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubObject) SignalTraceObject(event EventTrace) error {
	var buf bytes.Buffer
	if err := WriteEventTrace(event, &buf); err != nil {
		return fmt.Errorf("failed to serialize event: %s", err)
	}
	err := p.obj.UpdateSignal(uint32(0x56), buf.Bytes())

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
		Signals: map[uint32]object.MetaSignal{uint32(0x56): {
			Name:      "traceObject",
			Signature: "((IiIm(ll)<timeval,tv_sec,tv_usec>llII)<EventTrace,id,kind,slotId,arguments,timestamp,userUsTime,systemUsTime,callerContext,calleeContext>)",
			Uid:       uint32(0x56),
		}},
	}
}

// Constructor gives access to remote services
type Constructor struct {
	session bus.Session
}

// Services gives access to the services constructor
func Services(s bus.Session) Constructor {
	return Constructor{session: s}
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

// Object represents a proxy object to the service
type ObjectProxy interface {
	object.Object
	bus.Proxy
	Object
}

// proxyObject implements ObjectProxy
type proxyObject struct {
	bus.Proxy
}

// Object retruns a proxy to a remote service
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
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(objectID, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize objectID: %s", err)
	}
	if err = basic.WriteUint32(actionID, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize actionID: %s", err)
	}
	if err = basic.WriteUint64(handler, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize handler: %s", err)
	}
	response, err := p.Call("registerEvent", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerEvent failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse registerEvent response: %s", err)
	}
	return ret, nil
}

// UnregisterEvent calls the remote procedure
func (p *proxyObject) UnregisterEvent(objectID uint32, actionID uint32, handler uint64) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(objectID, buf); err != nil {
		return fmt.Errorf("failed to serialize objectID: %s", err)
	}
	if err = basic.WriteUint32(actionID, buf); err != nil {
		return fmt.Errorf("failed to serialize actionID: %s", err)
	}
	if err = basic.WriteUint64(handler, buf); err != nil {
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
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(objectID, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize objectID: %s", err)
	}
	response, err := p.Call("metaObject", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call metaObject failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = object.ReadMetaObject(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse metaObject response: %s", err)
	}
	return ret, nil
}

// Terminate calls the remote procedure
func (p *proxyObject) Terminate(objectID uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(objectID, buf); err != nil {
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
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = name.Write(buf); err != nil {
		return ret, fmt.Errorf("failed to serialize name: %s", err)
	}
	response, err := p.Call("property", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call property failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = value.NewValue(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse property response: %s", err)
	}
	return ret, nil
}

// SetProperty calls the remote procedure
func (p *proxyObject) SetProperty(name value.Value, value value.Value) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = name.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize name: %s", err)
	}
	if err = value.Write(buf); err != nil {
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
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("properties", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call properties failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(buf)
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
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(objectID, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize objectID: %s", err)
	}
	if err = basic.WriteUint32(actionID, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize actionID: %s", err)
	}
	if err = basic.WriteUint64(handler, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize handler: %s", err)
	}
	if err = basic.WriteString(P3, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P3: %s", err)
	}
	response, err := p.Call("registerEventWithSignature", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerEventWithSignature failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse registerEventWithSignature response: %s", err)
	}
	return ret, nil
}

// IsStatsEnabled calls the remote procedure
func (p *proxyObject) IsStatsEnabled() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("isStatsEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isStatsEnabled failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse isStatsEnabled response: %s", err)
	}
	return ret, nil
}

// EnableStats calls the remote procedure
func (p *proxyObject) EnableStats(enabled bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(enabled, buf); err != nil {
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
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("stats", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call stats failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (m map[uint32]MethodStatistics, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[uint32]MethodStatistics, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := ReadMethodStatistics(buf)
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
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
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
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("isTraceEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isTraceEnabled failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse isTraceEnabled response: %s", err)
	}
	return ret, nil
}

// EnableTrace calls the remote procedure
func (p *proxyObject) EnableTrace(traced bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(traced, buf); err != nil {
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
			e, err := ReadEventTrace(buf)
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

// ReadMetaMethodParameter unmarshalls MetaMethodParameter
func ReadMetaMethodParameter(r io.Reader) (s MetaMethodParameter, err error) {
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Description field: " + err.Error())
	}
	return s, nil
}

// WriteMetaMethodParameter marshalls MetaMethodParameter
func WriteMetaMethodParameter(s MetaMethodParameter, w io.Writer) (err error) {
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

// ReadMetaMethod unmarshalls MetaMethod
func ReadMetaMethod(r io.Reader) (s MetaMethod, err error) {
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
			b[i], err = ReadMetaMethodParameter(r)
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

// WriteMetaMethod marshalls MetaMethod
func WriteMetaMethod(s MetaMethod, w io.Writer) (err error) {
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
			err = WriteMetaMethodParameter(v, w)
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

// ReadMetaSignal unmarshalls MetaSignal
func ReadMetaSignal(r io.Reader) (s MetaSignal, err error) {
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

// WriteMetaSignal marshalls MetaSignal
func WriteMetaSignal(s MetaSignal, w io.Writer) (err error) {
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

// ReadMetaProperty unmarshalls MetaProperty
func ReadMetaProperty(r io.Reader) (s MetaProperty, err error) {
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

// WriteMetaProperty marshalls MetaProperty
func WriteMetaProperty(s MetaProperty, w io.Writer) (err error) {
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

// ReadMetaObject unmarshalls MetaObject
func ReadMetaObject(r io.Reader) (s MetaObject, err error) {
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
			v, err := ReadMetaMethod(r)
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
			v, err := ReadMetaSignal(r)
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
			v, err := ReadMetaProperty(r)
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

// WriteMetaObject marshalls MetaObject
func WriteMetaObject(s MetaObject, w io.Writer) (err error) {
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
			err = WriteMetaMethod(v, w)
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
			err = WriteMetaSignal(v, w)
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
			err = WriteMetaProperty(v, w)
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

// ReadMinMaxSum unmarshalls MinMaxSum
func ReadMinMaxSum(r io.Reader) (s MinMaxSum, err error) {
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

// WriteMinMaxSum marshalls MinMaxSum
func WriteMinMaxSum(s MinMaxSum, w io.Writer) (err error) {
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

// ReadMethodStatistics unmarshalls MethodStatistics
func ReadMethodStatistics(r io.Reader) (s MethodStatistics, err error) {
	if s.Count, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read Count field: " + err.Error())
	}
	if s.Wall, err = ReadMinMaxSum(r); err != nil {
		return s, fmt.Errorf("failed to read Wall field: " + err.Error())
	}
	if s.User, err = ReadMinMaxSum(r); err != nil {
		return s, fmt.Errorf("failed to read User field: " + err.Error())
	}
	if s.System, err = ReadMinMaxSum(r); err != nil {
		return s, fmt.Errorf("failed to read System field: " + err.Error())
	}
	return s, nil
}

// WriteMethodStatistics marshalls MethodStatistics
func WriteMethodStatistics(s MethodStatistics, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Count, w); err != nil {
		return fmt.Errorf("failed to write Count field: " + err.Error())
	}
	if err := WriteMinMaxSum(s.Wall, w); err != nil {
		return fmt.Errorf("failed to write Wall field: " + err.Error())
	}
	if err := WriteMinMaxSum(s.User, w); err != nil {
		return fmt.Errorf("failed to write User field: " + err.Error())
	}
	if err := WriteMinMaxSum(s.System, w); err != nil {
		return fmt.Errorf("failed to write System field: " + err.Error())
	}
	return nil
}

// Timeval is serializable
type Timeval struct {
	Tvsec  int64
	Tvusec int64
}

// ReadTimeval unmarshalls Timeval
func ReadTimeval(r io.Reader) (s Timeval, err error) {
	if s.Tvsec, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read Tvsec field: " + err.Error())
	}
	if s.Tvusec, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read Tvusec field: " + err.Error())
	}
	return s, nil
}

// WriteTimeval marshalls Timeval
func WriteTimeval(s Timeval, w io.Writer) (err error) {
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

// ReadEventTrace unmarshalls EventTrace
func ReadEventTrace(r io.Reader) (s EventTrace, err error) {
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
	if s.Timestamp, err = ReadTimeval(r); err != nil {
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

// WriteEventTrace marshalls EventTrace
func WriteEventTrace(s EventTrace, w io.Writer) (err error) {
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
	if err := WriteTimeval(s.Timestamp, w); err != nil {
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
