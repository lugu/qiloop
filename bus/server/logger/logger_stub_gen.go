// Package logger contains a generated stub
// File generated. DO NOT EDIT.

package logger

import (
	bytes "bytes"
	fmt "fmt"
	bus "github.com/lugu/qiloop/bus"
	net "github.com/lugu/qiloop/bus/net"
	server "github.com/lugu/qiloop/bus/server"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
	io "io"
	log "log"
)

// LogProviderImplementor interface of the service implementation
type LogProviderImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals an properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation server.Activation, helper LogProviderSignalHelper) error
	OnTerminate()
	SetVerbosity(level LogLevel) error
	SetCategory(category string, level LogLevel) error
	ClearAndSet(filters map[string]LogLevel) error
}

// LogProviderSignalHelper provided to LogProvider a companion object
type LogProviderSignalHelper interface{}

// stubLogProvider implements server.ServerObject.
type stubLogProvider struct {
	obj     server.Object
	impl    LogProviderImplementor
	session bus.Session
}

// LogProviderObject returns an object using LogProviderImplementor
func LogProviderObject(impl LogProviderImplementor) server.ServerObject {
	var stb stubLogProvider
	stb.impl = impl
	stb.obj = server.NewObject(stb.metaObject(), stb.onPropertyChange)
	stb.obj.Wrap(uint32(0x64), stb.SetVerbosity)
	stb.obj.Wrap(uint32(0x65), stb.SetCategory)
	stb.obj.Wrap(uint32(0x66), stb.ClearAndSet)
	return &stb
}
func (p *stubLogProvider) Activate(activation server.Activation) error {
	p.session = activation.Session
	p.obj.Activate(activation)
	return p.impl.Activate(activation, p)
}
func (p *stubLogProvider) OnTerminate() {
	p.impl.OnTerminate()
	p.obj.OnTerminate()
}
func (p *stubLogProvider) Receive(msg *net.Message, from *server.Context) error {
	return p.obj.Receive(msg, from)
}
func (p *stubLogProvider) onPropertyChange(name string, data []byte) error {
	switch name {
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}
func (p *stubLogProvider) SetVerbosity(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	level, err := ReadLogLevel(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read level: %s", err)
	}
	callErr := p.impl.SetVerbosity(level)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubLogProvider) SetCategory(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	category, err := basic.ReadString(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read category: %s", err)
	}
	level, err := ReadLogLevel(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read level: %s", err)
	}
	callErr := p.impl.SetCategory(category, level)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubLogProvider) ClearAndSet(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	filters, err := func() (m map[string]LogLevel, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]LogLevel, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := ReadLogLevel(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return nil, fmt.Errorf("cannot read filters: %s", err)
	}
	callErr := p.impl.ClearAndSet(filters)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubLogProvider) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "LogProvider",
		Methods: map[uint32]object.MetaMethod{
			uint32(0x64): {
				Name:                "setVerbosity",
				ParametersSignature: "((i)<LogLevel,level>)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x64),
			},
			uint32(0x65): {
				Name:                "setCategory",
				ParametersSignature: "(s(i)<LogLevel,level>)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x65),
			},
			uint32(0x66): {
				Name:                "clearAndSet",
				ParametersSignature: "({s(i)<LogLevel,level>})",
				ReturnSignature:     "v",
				Uid:                 uint32(0x66),
			},
		},
		Properties: map[uint32]object.MetaProperty{},
		Signals:    map[uint32]object.MetaSignal{},
	}
}

// LogListenerImplementor interface of the service implementation
type LogListenerImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals an properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation server.Activation, helper LogListenerSignalHelper) error
	OnTerminate()
	SetCategory(category string, level LogLevel) error
	ClearFilters() error
	// OnVerbosityChange is called when the property is updated.
	// Returns an error if the property value is not allowed
	OnVerbosityChange(level LogLevel) error
	// OnFiltersChange is called when the property is updated.
	// Returns an error if the property value is not allowed
	OnFiltersChange(filters map[string]int32) error
}

// LogListenerSignalHelper provided to LogListener a companion object
type LogListenerSignalHelper interface {
	SignalOnLogMessage(msg LogMessage) error
	UpdateVerbosity(level LogLevel) error
	UpdateFilters(filters map[string]int32) error
}

// stubLogListener implements server.ServerObject.
type stubLogListener struct {
	obj     server.Object
	impl    LogListenerImplementor
	session bus.Session
}

// LogListenerObject returns an object using LogListenerImplementor
func LogListenerObject(impl LogListenerImplementor) server.ServerObject {
	var stb stubLogListener
	stb.impl = impl
	stb.obj = server.NewObject(stb.metaObject(), stb.onPropertyChange)
	stb.obj.Wrap(uint32(0x65), stb.SetCategory)
	stb.obj.Wrap(uint32(0x66), stb.ClearFilters)
	return &stb
}
func (p *stubLogListener) Activate(activation server.Activation) error {
	p.session = activation.Session
	p.obj.Activate(activation)
	return p.impl.Activate(activation, p)
}
func (p *stubLogListener) OnTerminate() {
	p.impl.OnTerminate()
	p.obj.OnTerminate()
}
func (p *stubLogListener) Receive(msg *net.Message, from *server.Context) error {
	return p.obj.Receive(msg, from)
}
func (p *stubLogListener) onPropertyChange(name string, data []byte) error {
	switch name {
	case "verbosity":
		buf := bytes.NewBuffer(data)
		prop, err := ReadLogLevel(buf)
		if err != nil {
			return fmt.Errorf("cannot read Verbosity: %s", err)
		}
		return p.impl.OnVerbosityChange(prop)
	case "filters":
		buf := bytes.NewBuffer(data)
		prop, err := func() (m map[string]int32, err error) {
			size, err := basic.ReadUint32(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map size: %s", err)
			}
			m = make(map[string]int32, size)
			for i := 0; i < int(size); i++ {
				k, err := basic.ReadString(buf)
				if err != nil {
					return m, fmt.Errorf("failed to read map key: %s", err)
				}
				v, err := basic.ReadInt32(buf)
				if err != nil {
					return m, fmt.Errorf("failed to read map value: %s", err)
				}
				m[k] = v
			}
			return m, nil
		}()
		if err != nil {
			return fmt.Errorf("cannot read Filters: %s", err)
		}
		return p.impl.OnFiltersChange(prop)
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}
func (p *stubLogListener) SetCategory(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	category, err := basic.ReadString(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read category: %s", err)
	}
	level, err := ReadLogLevel(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read level: %s", err)
	}
	callErr := p.impl.SetCategory(category, level)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubLogListener) ClearFilters(payload []byte) ([]byte, error) {
	callErr := p.impl.ClearFilters()
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubLogListener) SignalOnLogMessage(msg LogMessage) error {
	var buf bytes.Buffer
	if err := WriteLogMessage(msg, &buf); err != nil {
		return fmt.Errorf("failed to serialize msg: %s", err)
	}
	err := p.obj.UpdateSignal(uint32(0x67), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalOnLogMessage: %s", err)
	}
	return nil
}
func (p *stubLogListener) UpdateVerbosity(level LogLevel) error {
	var buf bytes.Buffer
	if err := WriteLogLevel(level, &buf); err != nil {
		return fmt.Errorf("failed to serialize level: %s", err)
	}
	err := p.obj.UpdateProperty(uint32(0x68), "(i)<LogLevel,level>", buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update UpdateVerbosity: %s", err)
	}
	return nil
}
func (p *stubLogListener) UpdateFilters(filters map[string]int32) error {
	var buf bytes.Buffer
	if err := func() error {
		err := basic.WriteUint32(uint32(len(filters)), &buf)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range filters {
			err = basic.WriteString(k, &buf)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = basic.WriteInt32(v, &buf)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to serialize filters: %s", err)
	}
	err := p.obj.UpdateProperty(uint32(0x69), "{si}", buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update UpdateFilters: %s", err)
	}
	return nil
}
func (p *stubLogListener) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "LogListener",
		Methods: map[uint32]object.MetaMethod{
			uint32(0x65): {
				Name:                "setCategory",
				ParametersSignature: "(s(i)<LogLevel,level>)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x65),
			},
			uint32(0x66): {
				Name:                "clearFilters",
				ParametersSignature: "()",
				ReturnSignature:     "v",
				Uid:                 uint32(0x66),
			},
		},
		Properties: map[uint32]object.MetaProperty{
			uint32(0x68): {
				Name:      "verbosity",
				Signature: "(i)<LogLevel,level>",
				Uid:       uint32(0x68),
			},
			uint32(0x69): {
				Name:      "filters",
				Signature: "{si}",
				Uid:       uint32(0x69),
			},
		},
		Signals: map[uint32]object.MetaSignal{uint32(0x67): {
			Name:      "onLogMessage",
			Signature: "(s(i)<LogLevel,level>sssI(L)<TimePoint,ns>(L)<TimePoint,ns>)<LogMessage,source,level,category,location,message,id,date,systemDate>",
			Uid:       uint32(0x67),
		}},
	}
}

// LogManagerImplementor interface of the service implementation
type LogManagerImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals an properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation server.Activation, helper LogManagerSignalHelper) error
	OnTerminate()
	Log(messages []LogMessage) error
	CreateListener() (LogListenerProxy, error)
	GetListener() (LogListenerProxy, error)
	AddProvider(source LogProviderProxy) (int32, error)
	RemoveProvider(sourceID int32) error
}

// LogManagerSignalHelper provided to LogManager a companion object
type LogManagerSignalHelper interface{}

// stubLogManager implements server.ServerObject.
type stubLogManager struct {
	obj     server.Object
	impl    LogManagerImplementor
	session bus.Session
}

// LogManagerObject returns an object using LogManagerImplementor
func LogManagerObject(impl LogManagerImplementor) server.ServerObject {
	var stb stubLogManager
	stb.impl = impl
	stb.obj = server.NewObject(stb.metaObject(), stb.onPropertyChange)
	stb.obj.Wrap(uint32(0x64), stb.Log)
	stb.obj.Wrap(uint32(0x65), stb.CreateListener)
	stb.obj.Wrap(uint32(0x66), stb.GetListener)
	stb.obj.Wrap(uint32(0x67), stb.AddProvider)
	stb.obj.Wrap(uint32(0x68), stb.RemoveProvider)
	return &stb
}
func (p *stubLogManager) Activate(activation server.Activation) error {
	p.session = activation.Session
	p.obj.Activate(activation)
	return p.impl.Activate(activation, p)
}
func (p *stubLogManager) OnTerminate() {
	p.impl.OnTerminate()
	p.obj.OnTerminate()
}
func (p *stubLogManager) Receive(msg *net.Message, from *server.Context) error {
	return p.obj.Receive(msg, from)
}
func (p *stubLogManager) onPropertyChange(name string, data []byte) error {
	switch name {
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}
func (p *stubLogManager) Log(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	messages, err := func() (b []LogMessage, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]LogMessage, size)
		for i := 0; i < int(size); i++ {
			b[i], err = ReadLogMessage(buf)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return nil, fmt.Errorf("cannot read messages: %s", err)
	}
	callErr := p.impl.Log(messages)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubLogManager) CreateListener(payload []byte) ([]byte, error) {
	ret, callErr := p.impl.CreateListener()
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := func() error {
		meta, err := ret.MetaObject(ret.ObjectID())
		if err != nil {
			return fmt.Errorf("failed to get meta: %s", err)
		}
		ref := object.ObjectReference{
			true,
			meta,
			0,
			ret.ServiceID(),
			ret.ObjectID(),
		}
		return object.WriteObjectReference(ref, &out)
	}()
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubLogManager) GetListener(payload []byte) ([]byte, error) {
	ret, callErr := p.impl.GetListener()
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := func() error {
		meta, err := ret.MetaObject(ret.ObjectID())
		if err != nil {
			return fmt.Errorf("failed to get meta: %s", err)
		}
		ref := object.ObjectReference{
			true,
			meta,
			0,
			ret.ServiceID(),
			ret.ObjectID(),
		}
		return object.WriteObjectReference(ref, &out)
	}()
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubLogManager) AddProvider(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	source, err := func() (LogProviderProxy, error) {
		ref, err := object.ReadObjectReference(buf)
		if err != nil {
			return nil, fmt.Errorf("failed to get meta: %s", err)
		}
		proxy, err := p.session.Object(ref)
		if err != nil {
			return nil, fmt.Errorf("failed to get proxy: %s", err)
		}
		return MakeLogProvider(p.session, proxy), nil
	}()
	if err != nil {
		return nil, fmt.Errorf("cannot read source: %s", err)
	}
	ret, callErr := p.impl.AddProvider(source)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteInt32(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubLogManager) RemoveProvider(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	sourceID, err := basic.ReadInt32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read sourceID: %s", err)
	}
	callErr := p.impl.RemoveProvider(sourceID)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubLogManager) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "LogManager",
		Methods: map[uint32]object.MetaMethod{
			uint32(0x64): {
				Name:                "log",
				ParametersSignature: "([(s(i)<LogLevel,level>sssI(L)<TimePoint,ns>(L)<TimePoint,ns>)<LogMessage,source,level,category,location,message,id,date,systemDate>])",
				ReturnSignature:     "v",
				Uid:                 uint32(0x64),
			},
			uint32(0x65): {
				Name:                "createListener",
				ParametersSignature: "()",
				ReturnSignature:     "o",
				Uid:                 uint32(0x65),
			},
			uint32(0x66): {
				Name:                "getListener",
				ParametersSignature: "()",
				ReturnSignature:     "o",
				Uid:                 uint32(0x66),
			},
			uint32(0x67): {
				Name:                "addProvider",
				ParametersSignature: "(o)",
				ReturnSignature:     "i",
				Uid:                 uint32(0x67),
			},
			uint32(0x68): {
				Name:                "removeProvider",
				ParametersSignature: "(i)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x68),
			},
		},
		Properties: map[uint32]object.MetaProperty{},
		Signals:    map[uint32]object.MetaSignal{},
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

// LogLevel is serializable
type LogLevel struct {
	Level int32
}

// ReadLogLevel unmarshalls LogLevel
func ReadLogLevel(r io.Reader) (s LogLevel, err error) {
	if s.Level, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("failed to read Level field: " + err.Error())
	}
	return s, nil
}

// WriteLogLevel marshalls LogLevel
func WriteLogLevel(s LogLevel, w io.Writer) (err error) {
	if err := basic.WriteInt32(s.Level, w); err != nil {
		return fmt.Errorf("failed to write Level field: " + err.Error())
	}
	return nil
}

// TimePoint is serializable
type TimePoint struct {
	Ns uint64
}

// ReadTimePoint unmarshalls TimePoint
func ReadTimePoint(r io.Reader) (s TimePoint, err error) {
	if s.Ns, err = basic.ReadUint64(r); err != nil {
		return s, fmt.Errorf("failed to read Ns field: " + err.Error())
	}
	return s, nil
}

// WriteTimePoint marshalls TimePoint
func WriteTimePoint(s TimePoint, w io.Writer) (err error) {
	if err := basic.WriteUint64(s.Ns, w); err != nil {
		return fmt.Errorf("failed to write Ns field: " + err.Error())
	}
	return nil
}

// LogMessage is serializable
type LogMessage struct {
	Source     string
	Level      LogLevel
	Category   string
	Location   string
	Message    string
	Id         uint32
	Date       TimePoint
	SystemDate TimePoint
}

// ReadLogMessage unmarshalls LogMessage
func ReadLogMessage(r io.Reader) (s LogMessage, err error) {
	if s.Source, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Source field: " + err.Error())
	}
	if s.Level, err = ReadLogLevel(r); err != nil {
		return s, fmt.Errorf("failed to read Level field: " + err.Error())
	}
	if s.Category, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Category field: " + err.Error())
	}
	if s.Location, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Location field: " + err.Error())
	}
	if s.Message, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Message field: " + err.Error())
	}
	if s.Id, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read Id field: " + err.Error())
	}
	if s.Date, err = ReadTimePoint(r); err != nil {
		return s, fmt.Errorf("failed to read Date field: " + err.Error())
	}
	if s.SystemDate, err = ReadTimePoint(r); err != nil {
		return s, fmt.Errorf("failed to read SystemDate field: " + err.Error())
	}
	return s, nil
}

// WriteLogMessage marshalls LogMessage
func WriteLogMessage(s LogMessage, w io.Writer) (err error) {
	if err := basic.WriteString(s.Source, w); err != nil {
		return fmt.Errorf("failed to write Source field: " + err.Error())
	}
	if err := WriteLogLevel(s.Level, w); err != nil {
		return fmt.Errorf("failed to write Level field: " + err.Error())
	}
	if err := basic.WriteString(s.Category, w); err != nil {
		return fmt.Errorf("failed to write Category field: " + err.Error())
	}
	if err := basic.WriteString(s.Location, w); err != nil {
		return fmt.Errorf("failed to write Location field: " + err.Error())
	}
	if err := basic.WriteString(s.Message, w); err != nil {
		return fmt.Errorf("failed to write Message field: " + err.Error())
	}
	if err := basic.WriteUint32(s.Id, w); err != nil {
		return fmt.Errorf("failed to write Id field: " + err.Error())
	}
	if err := WriteTimePoint(s.Date, w); err != nil {
		return fmt.Errorf("failed to write Date field: " + err.Error())
	}
	if err := WriteTimePoint(s.SystemDate, w); err != nil {
		return fmt.Errorf("failed to write SystemDate field: " + err.Error())
	}
	return nil
}

// LogProvider is the abstract interface of the service
type LogProvider interface {
	// SetVerbosity calls the remote procedure
	SetVerbosity(level LogLevel) error
	// SetCategory calls the remote procedure
	SetCategory(category string, level LogLevel) error
	// ClearAndSet calls the remote procedure
	ClearAndSet(filters map[string]LogLevel) error
}

// LogProvider represents a proxy object to the service
type LogProviderProxy interface {
	object.Object
	bus.Proxy
	LogProvider
}

// proxyLogProvider implements LogProviderProxy
type proxyLogProvider struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeLogProvider constructs LogProviderProxy
func MakeLogProvider(sess bus.Session, proxy bus.Proxy) LogProviderProxy {
	return &proxyLogProvider{bus.MakeObject(proxy), sess}
}

// LogProvider retruns a proxy to a remote service
func (s Constructor) LogProvider() (LogProviderProxy, error) {
	proxy, err := s.session.Proxy("LogProvider", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return MakeLogProvider(s.session, proxy), nil
}

// SetVerbosity calls the remote procedure
func (p *proxyLogProvider) SetVerbosity(level LogLevel) error {
	var err error
	var buf bytes.Buffer
	if err = WriteLogLevel(level, &buf); err != nil {
		return fmt.Errorf("failed to serialize level: %s", err)
	}
	_, err = p.Call("setVerbosity", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setVerbosity failed: %s", err)
	}
	return nil
}

// SetCategory calls the remote procedure
func (p *proxyLogProvider) SetCategory(category string, level LogLevel) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(category, &buf); err != nil {
		return fmt.Errorf("failed to serialize category: %s", err)
	}
	if err = WriteLogLevel(level, &buf); err != nil {
		return fmt.Errorf("failed to serialize level: %s", err)
	}
	_, err = p.Call("setCategory", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setCategory failed: %s", err)
	}
	return nil
}

// ClearAndSet calls the remote procedure
func (p *proxyLogProvider) ClearAndSet(filters map[string]LogLevel) error {
	var err error
	var buf bytes.Buffer
	if err = func() error {
		err := basic.WriteUint32(uint32(len(filters)), &buf)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range filters {
			err = basic.WriteString(k, &buf)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = WriteLogLevel(v, &buf)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to serialize filters: %s", err)
	}
	_, err = p.Call("clearAndSet", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearAndSet failed: %s", err)
	}
	return nil
}

// LogListener is the abstract interface of the service
type LogListener interface {
	// SetCategory calls the remote procedure
	SetCategory(category string, level LogLevel) error
	// ClearFilters calls the remote procedure
	ClearFilters() error
	// SubscribeOnLogMessage subscribe to a remote signal
	SubscribeOnLogMessage() (unsubscribe func(), updates chan LogMessage, err error)
	// GetVerbosity returns the property value
	GetVerbosity() (LogLevel, error)
	// SetVerbosity sets the property value
	SetVerbosity(LogLevel) error
	// SubscribeVerbosity regusters to a property
	SubscribeVerbosity() (unsubscribe func(), updates chan LogLevel, err error)
	// GetFilters returns the property value
	GetFilters() (map[string]int32, error)
	// SetFilters sets the property value
	SetFilters(map[string]int32) error
	// SubscribeFilters regusters to a property
	SubscribeFilters() (unsubscribe func(), updates chan map[string]int32, err error)
}

// LogListener represents a proxy object to the service
type LogListenerProxy interface {
	object.Object
	bus.Proxy
	LogListener
}

// proxyLogListener implements LogListenerProxy
type proxyLogListener struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeLogListener constructs LogListenerProxy
func MakeLogListener(sess bus.Session, proxy bus.Proxy) LogListenerProxy {
	return &proxyLogListener{bus.MakeObject(proxy), sess}
}

// LogListener retruns a proxy to a remote service
func (s Constructor) LogListener() (LogListenerProxy, error) {
	proxy, err := s.session.Proxy("LogListener", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return MakeLogListener(s.session, proxy), nil
}

// SetCategory calls the remote procedure
func (p *proxyLogListener) SetCategory(category string, level LogLevel) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(category, &buf); err != nil {
		return fmt.Errorf("failed to serialize category: %s", err)
	}
	if err = WriteLogLevel(level, &buf); err != nil {
		return fmt.Errorf("failed to serialize level: %s", err)
	}
	_, err = p.Call("setCategory", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setCategory failed: %s", err)
	}
	return nil
}

// ClearFilters calls the remote procedure
func (p *proxyLogListener) ClearFilters() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("clearFilters", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearFilters failed: %s", err)
	}
	return nil
}

// SubscribeOnLogMessage subscribe to a remote property
func (p *proxyLogListener) SubscribeOnLogMessage() (func(), chan LogMessage, error) {
	propertyID, err := p.SignalID("onLogMessage")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "onLogMessage", err)
	}

	handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register event for %s: %s", "onLogMessage", err)
	}
	ch := make(chan LogMessage)
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
			e, err := ReadLogMessage(buf)
			if err != nil {
				log.Printf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// GetVerbosity updates the property value
func (p *proxyLogListener) GetVerbosity() (ret LogLevel, err error) {
	name := value.String("verbosity")
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
	sig := "(i)<LogLevel,level>"
	if sig != s {
		return ret, fmt.Errorf("unexpected signature: %s instead of %s",
			s, sig)
	}
	ret, err = ReadLogLevel(&buf)
	return ret, err
}

// SetVerbosity updates the property value
func (p *proxyLogListener) SetVerbosity(update LogLevel) error {
	name := value.String("verbosity")
	var buf bytes.Buffer
	err := WriteLogLevel(update, &buf)
	if err != nil {
		return fmt.Errorf("marshall error: %s", err)
	}
	val := value.Opaque("(i)<LogLevel,level>", buf.Bytes())
	return p.SetProperty(name, val)
}

// SubscribeVerbosity subscribe to a remote property
func (p *proxyLogListener) SubscribeVerbosity() (func(), chan LogLevel, error) {
	propertyID, err := p.PropertyID("verbosity")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "verbosity", err)
	}

	handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register event for %s: %s", "verbosity", err)
	}
	ch := make(chan LogLevel)
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
			e, err := ReadLogLevel(buf)
			if err != nil {
				log.Printf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// GetFilters updates the property value
func (p *proxyLogListener) GetFilters() (ret map[string]int32, err error) {
	name := value.String("filters")
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
	sig := "{si}"
	if sig != s {
		return ret, fmt.Errorf("unexpected signature: %s instead of %s",
			s, sig)
	}
	ret, err = func() (m map[string]int32, err error) {
		size, err := basic.ReadUint32(&buf)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]int32, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(&buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := basic.ReadInt32(&buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}()
	return ret, err
}

// SetFilters updates the property value
func (p *proxyLogListener) SetFilters(update map[string]int32) error {
	name := value.String("filters")
	var buf bytes.Buffer
	err := func() error {
		err := basic.WriteUint32(uint32(len(update)), &buf)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range update {
			err = basic.WriteString(k, &buf)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = basic.WriteInt32(v, &buf)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}()
	if err != nil {
		return fmt.Errorf("marshall error: %s", err)
	}
	val := value.Opaque("{si}", buf.Bytes())
	return p.SetProperty(name, val)
}

// SubscribeFilters subscribe to a remote property
func (p *proxyLogListener) SubscribeFilters() (func(), chan map[string]int32, error) {
	propertyID, err := p.PropertyID("filters")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "filters", err)
	}

	handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register event for %s: %s", "filters", err)
	}
	ch := make(chan map[string]int32)
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
			e, err := func() (m map[string]int32, err error) {
				size, err := basic.ReadUint32(buf)
				if err != nil {
					return m, fmt.Errorf("failed to read map size: %s", err)
				}
				m = make(map[string]int32, size)
				for i := 0; i < int(size); i++ {
					k, err := basic.ReadString(buf)
					if err != nil {
						return m, fmt.Errorf("failed to read map key: %s", err)
					}
					v, err := basic.ReadInt32(buf)
					if err != nil {
						return m, fmt.Errorf("failed to read map value: %s", err)
					}
					m[k] = v
				}
				return m, nil
			}()
			if err != nil {
				log.Printf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// LogManager is the abstract interface of the service
type LogManager interface {
	// Log calls the remote procedure
	Log(messages []LogMessage) error
	// CreateListener calls the remote procedure
	CreateListener() (LogListenerProxy, error)
	// GetListener calls the remote procedure
	GetListener() (LogListenerProxy, error)
	// AddProvider calls the remote procedure
	AddProvider(source LogProviderProxy) (int32, error)
	// RemoveProvider calls the remote procedure
	RemoveProvider(sourceID int32) error
}

// LogManager represents a proxy object to the service
type LogManagerProxy interface {
	object.Object
	bus.Proxy
	LogManager
}

// proxyLogManager implements LogManagerProxy
type proxyLogManager struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeLogManager constructs LogManagerProxy
func MakeLogManager(sess bus.Session, proxy bus.Proxy) LogManagerProxy {
	return &proxyLogManager{bus.MakeObject(proxy), sess}
}

// LogManager retruns a proxy to a remote service
func (s Constructor) LogManager() (LogManagerProxy, error) {
	proxy, err := s.session.Proxy("LogManager", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return MakeLogManager(s.session, proxy), nil
}

// Log calls the remote procedure
func (p *proxyLogManager) Log(messages []LogMessage) error {
	var err error
	var buf bytes.Buffer
	if err = func() error {
		err := basic.WriteUint32(uint32(len(messages)), &buf)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range messages {
			err = WriteLogMessage(v, &buf)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to serialize messages: %s", err)
	}
	_, err = p.Call("log", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call log failed: %s", err)
	}
	return nil
}

// CreateListener calls the remote procedure
func (p *proxyLogManager) CreateListener() (LogListenerProxy, error) {
	var err error
	var ret LogListenerProxy
	var buf bytes.Buffer
	response, err := p.Call("createListener", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call createListener failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (LogListenerProxy, error) {
		ref, err := object.ReadObjectReference(resp)
		if err != nil {
			return nil, fmt.Errorf("failed to get meta: %s", err)
		}
		proxy, err := p.session.Object(ref)
		if err != nil {
			return nil, fmt.Errorf("failed to get proxy: %s", err)
		}
		return MakeLogListener(p.session, proxy), nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse createListener response: %s", err)
	}
	return ret, nil
}

// GetListener calls the remote procedure
func (p *proxyLogManager) GetListener() (LogListenerProxy, error) {
	var err error
	var ret LogListenerProxy
	var buf bytes.Buffer
	response, err := p.Call("getListener", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getListener failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (LogListenerProxy, error) {
		ref, err := object.ReadObjectReference(resp)
		if err != nil {
			return nil, fmt.Errorf("failed to get meta: %s", err)
		}
		proxy, err := p.session.Object(ref)
		if err != nil {
			return nil, fmt.Errorf("failed to get proxy: %s", err)
		}
		return MakeLogListener(p.session, proxy), nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse getListener response: %s", err)
	}
	return ret, nil
}

// AddProvider calls the remote procedure
func (p *proxyLogManager) AddProvider(source LogProviderProxy) (int32, error) {
	var err error
	var ret int32
	var buf bytes.Buffer
	if err = func() error {
		meta, err := source.MetaObject(source.ObjectID())
		if err != nil {
			return fmt.Errorf("failed to get meta: %s", err)
		}
		ref := object.ObjectReference{
			true,
			meta,
			0,
			source.ServiceID(),
			source.ObjectID(),
		}
		return object.WriteObjectReference(ref, &buf)
	}(); err != nil {
		return ret, fmt.Errorf("failed to serialize source: %s", err)
	}
	response, err := p.Call("addProvider", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call addProvider failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadInt32(resp)
	if err != nil {
		return ret, fmt.Errorf("failed to parse addProvider response: %s", err)
	}
	return ret, nil
}

// RemoveProvider calls the remote procedure
func (p *proxyLogManager) RemoveProvider(sourceID int32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteInt32(sourceID, &buf); err != nil {
		return fmt.Errorf("failed to serialize sourceID: %s", err)
	}
	_, err = p.Call("removeProvider", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call removeProvider failed: %s", err)
	}
	return nil
}
