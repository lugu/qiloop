package logger

import (
	"bytes"
	"context"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	net "github.com/lugu/qiloop/bus/net"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
	"io"
	"log"
)

// LogProviderImplementor interface of the service implementation
type LogProviderImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals and properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation bus.Activation, helper LogProviderSignalHelper) error
	OnTerminate()
	SetVerbosity(level LogLevel) error
	SetCategory(category string, level LogLevel) error
	ClearAndSet(filters map[string]LogLevel) error
}

// LogProviderSignalHelper provided to LogProvider a companion object
type LogProviderSignalHelper interface{}

// stubLogProvider implements server.Actor.
type stubLogProvider struct {
	impl      LogProviderImplementor
	session   bus.Session
	service   bus.Service
	serviceID uint32
	signal    bus.SignalHandler
}

// LogProviderObject returns an object using LogProviderImplementor
func LogProviderObject(impl LogProviderImplementor) bus.Actor {
	var stb stubLogProvider
	stb.impl = impl
	obj := bus.NewBasicObject(&stb, stb.metaObject(), stb.onPropertyChange)
	stb.signal = obj
	return obj
}

// CreateLogProvider registers a new object to a service
// and returns a proxy to the newly created object
func CreateLogProvider(session bus.Session, service bus.Service, impl LogProviderImplementor) (LogProviderProxy, error) {
	obj := LogProviderObject(impl)
	objectID, err := service.Add(obj)
	if err != nil {
		return nil, err
	}
	stb := &stubLogProvider{}
	meta := object.FullMetaObject(stb.metaObject())
	client := bus.DirectClient(obj)
	proxy := bus.NewProxy(client, meta, service.ServiceID(), objectID)
	return MakeLogProvider(session, proxy), nil
}
func (p *stubLogProvider) Activate(activation bus.Activation) error {
	p.session = activation.Session
	p.service = activation.Service
	p.serviceID = activation.ServiceID
	return p.impl.Activate(activation, p)
}
func (p *stubLogProvider) OnTerminate() {
	p.impl.OnTerminate()
}
func (p *stubLogProvider) Receive(msg *net.Message, from bus.Channel) error {
	// action dispatch
	switch msg.Header.Action {
	case 100:
		return p.SetVerbosity(msg, from)
	case 101:
		return p.SetCategory(msg, from)
	case 102:
		return p.ClearAndSet(msg, from)
	default:
		return from.SendError(msg, bus.ErrActionNotFound)
	}
}
func (p *stubLogProvider) onPropertyChange(name string, data []byte) error {
	switch name {
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}
func (p *stubLogProvider) SetVerbosity(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	level, err := readLogLevel(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read level: %s", err))
	}
	callErr := p.impl.SetVerbosity(level)

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
func (p *stubLogProvider) SetCategory(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	category, err := basic.ReadString(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read category: %s", err))
	}
	level, err := readLogLevel(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read level: %s", err))
	}
	callErr := p.impl.SetCategory(category, level)

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
func (p *stubLogProvider) ClearAndSet(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	filters, err := func() (m map[string]LogLevel, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return m, fmt.Errorf("read map size: %s", err)
		}
		m = make(map[string]LogLevel, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(buf)
			if err != nil {
				return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
			}
			v, err := readLogLevel(buf)
			if err != nil {
				return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read filters: %s", err))
	}
	callErr := p.impl.ClearAndSet(filters)

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
func (p *stubLogProvider) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "LogProvider",
		Methods: map[uint32]object.MetaMethod{
			100: {
				Name:                "setVerbosity",
				ParametersSignature: "((i)<LogLevel,level>)",
				ReturnSignature:     "v",
				Uid:                 100,
			},
			101: {
				Name:                "setCategory",
				ParametersSignature: "(s(i)<LogLevel,level>)",
				ReturnSignature:     "v",
				Uid:                 101,
			},
			102: {
				Name:                "clearAndSet",
				ParametersSignature: "({s(i)<LogLevel,level>})",
				ReturnSignature:     "v",
				Uid:                 102,
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
	// helper enables signals and properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation bus.Activation, helper LogListenerSignalHelper) error
	OnTerminate()
	SetLevel(level LogLevel) error
	AddFilter(category string, level LogLevel) error
	ClearFilters() error
	// OnLogLevelChange is called when the property is updated.
	// Returns an error if the property value is not allowed
	OnLogLevelChange(level LogLevel) error
}

// LogListenerSignalHelper provided to LogListener a companion object
type LogListenerSignalHelper interface {
	SignalOnLogMessage(log LogMessage) error
	SignalOnLogMessages(logs []LogMessage) error
	SignalOnLogMessagesWithBacklog(logs []LogMessage) error
	UpdateLogLevel(level LogLevel) error
}

// stubLogListener implements server.Actor.
type stubLogListener struct {
	impl      LogListenerImplementor
	session   bus.Session
	service   bus.Service
	serviceID uint32
	signal    bus.SignalHandler
}

// LogListenerObject returns an object using LogListenerImplementor
func LogListenerObject(impl LogListenerImplementor) bus.Actor {
	var stb stubLogListener
	stb.impl = impl
	obj := bus.NewBasicObject(&stb, stb.metaObject(), stb.onPropertyChange)
	stb.signal = obj
	return obj
}

// CreateLogListener registers a new object to a service
// and returns a proxy to the newly created object
func CreateLogListener(session bus.Session, service bus.Service, impl LogListenerImplementor) (LogListenerProxy, error) {
	obj := LogListenerObject(impl)
	objectID, err := service.Add(obj)
	if err != nil {
		return nil, err
	}
	stb := &stubLogListener{}
	meta := object.FullMetaObject(stb.metaObject())
	client := bus.DirectClient(obj)
	proxy := bus.NewProxy(client, meta, service.ServiceID(), objectID)
	return MakeLogListener(session, proxy), nil
}
func (p *stubLogListener) Activate(activation bus.Activation) error {
	p.session = activation.Session
	p.service = activation.Service
	p.serviceID = activation.ServiceID
	return p.impl.Activate(activation, p)
}
func (p *stubLogListener) OnTerminate() {
	p.impl.OnTerminate()
}
func (p *stubLogListener) Receive(msg *net.Message, from bus.Channel) error {
	// action dispatch
	switch msg.Header.Action {
	case 100:
		return p.SetLevel(msg, from)
	case 101:
		return p.AddFilter(msg, from)
	case 102:
		return p.ClearFilters(msg, from)
	default:
		return from.SendError(msg, bus.ErrActionNotFound)
	}
}
func (p *stubLogListener) onPropertyChange(name string, data []byte) error {
	switch name {
	case "logLevel":
		buf := bytes.NewBuffer(data)
		prop, err := readLogLevel(buf)
		if err != nil {
			return fmt.Errorf("cannot read LogLevel: %s", err)
		}
		return p.impl.OnLogLevelChange(prop)
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}
func (p *stubLogListener) SetLevel(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	level, err := readLogLevel(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read level: %s", err))
	}
	callErr := p.impl.SetLevel(level)

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
func (p *stubLogListener) AddFilter(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	category, err := basic.ReadString(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read category: %s", err))
	}
	level, err := readLogLevel(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read level: %s", err))
	}
	callErr := p.impl.AddFilter(category, level)

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
func (p *stubLogListener) ClearFilters(msg *net.Message, c bus.Channel) error {
	callErr := p.impl.ClearFilters()

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
func (p *stubLogListener) SignalOnLogMessage(log LogMessage) error {
	var buf bytes.Buffer
	if err := writeLogMessage(log, &buf); err != nil {
		return fmt.Errorf("serialize log: %s", err)
	}
	err := p.signal.UpdateSignal(103, buf.Bytes())

	if err != nil {
		return fmt.Errorf("update SignalOnLogMessage: %s", err)
	}
	return nil
}
func (p *stubLogListener) SignalOnLogMessages(logs []LogMessage) error {
	var buf bytes.Buffer
	if err := func() error {
		err := basic.WriteUint32(uint32(len(logs)), &buf)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range logs {
			err = writeLogMessage(v, &buf)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize logs: %s", err)
	}
	err := p.signal.UpdateSignal(104, buf.Bytes())

	if err != nil {
		return fmt.Errorf("update SignalOnLogMessages: %s", err)
	}
	return nil
}
func (p *stubLogListener) SignalOnLogMessagesWithBacklog(logs []LogMessage) error {
	var buf bytes.Buffer
	if err := func() error {
		err := basic.WriteUint32(uint32(len(logs)), &buf)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range logs {
			err = writeLogMessage(v, &buf)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize logs: %s", err)
	}
	err := p.signal.UpdateSignal(105, buf.Bytes())

	if err != nil {
		return fmt.Errorf("update SignalOnLogMessagesWithBacklog: %s", err)
	}
	return nil
}
func (p *stubLogListener) UpdateLogLevel(level LogLevel) error {
	var buf bytes.Buffer
	if err := writeLogLevel(level, &buf); err != nil {
		return fmt.Errorf("serialize level: %s", err)
	}
	err := p.signal.UpdateProperty(106, "(i)<LogLevel,level>", buf.Bytes())

	if err != nil {
		return fmt.Errorf("update UpdateLogLevel: %s", err)
	}
	return nil
}
func (p *stubLogListener) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "LogListener",
		Methods: map[uint32]object.MetaMethod{
			100: {
				Name:                "setLevel",
				ParametersSignature: "((i)<LogLevel,level>)",
				ReturnSignature:     "v",
				Uid:                 100,
			},
			101: {
				Name:                "addFilter",
				ParametersSignature: "(s(i)<LogLevel,level>)",
				ReturnSignature:     "v",
				Uid:                 101,
			},
			102: {
				Name:                "clearFilters",
				ParametersSignature: "()",
				ReturnSignature:     "v",
				Uid:                 102,
			},
		},
		Properties: map[uint32]object.MetaProperty{106: {
			Name:      "logLevel",
			Signature: "(i)<LogLevel,level>",
			Uid:       106,
		}},
		Signals: map[uint32]object.MetaSignal{
			103: {
				Name:      "onLogMessage",
				Signature: "(s(i)<LogLevel,level>sssI(L)<TimePoint,ns>(L)<TimePoint,ns>)<LogMessage,source,level,category,location,message,id,date,systemDate>",
				Uid:       103,
			},
			104: {
				Name:      "onLogMessages",
				Signature: "[(s(i)<LogLevel,level>sssI(L)<TimePoint,ns>(L)<TimePoint,ns>)<LogMessage,source,level,category,location,message,id,date,systemDate>]",
				Uid:       104,
			},
			105: {
				Name:      "onLogMessagesWithBacklog",
				Signature: "[(s(i)<LogLevel,level>sssI(L)<TimePoint,ns>(L)<TimePoint,ns>)<LogMessage,source,level,category,location,message,id,date,systemDate>]",
				Uid:       105,
			},
		},
	}
}

// LogManagerImplementor interface of the service implementation
type LogManagerImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals and properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation bus.Activation, helper LogManagerSignalHelper) error
	OnTerminate()
	Log(messages []LogMessage) error
	CreateListener() (LogListenerProxy, error)
	GetListener() (LogListenerProxy, error)
	AddProvider(source LogProviderProxy) (int32, error)
	RemoveProvider(sourceID int32) error
}

// LogManagerSignalHelper provided to LogManager a companion object
type LogManagerSignalHelper interface{}

// stubLogManager implements server.Actor.
type stubLogManager struct {
	impl      LogManagerImplementor
	session   bus.Session
	service   bus.Service
	serviceID uint32
	signal    bus.SignalHandler
}

// LogManagerObject returns an object using LogManagerImplementor
func LogManagerObject(impl LogManagerImplementor) bus.Actor {
	var stb stubLogManager
	stb.impl = impl
	obj := bus.NewBasicObject(&stb, stb.metaObject(), stb.onPropertyChange)
	stb.signal = obj
	return obj
}

// CreateLogManager registers a new object to a service
// and returns a proxy to the newly created object
func CreateLogManager(session bus.Session, service bus.Service, impl LogManagerImplementor) (LogManagerProxy, error) {
	obj := LogManagerObject(impl)
	objectID, err := service.Add(obj)
	if err != nil {
		return nil, err
	}
	stb := &stubLogManager{}
	meta := object.FullMetaObject(stb.metaObject())
	client := bus.DirectClient(obj)
	proxy := bus.NewProxy(client, meta, service.ServiceID(), objectID)
	return MakeLogManager(session, proxy), nil
}
func (p *stubLogManager) Activate(activation bus.Activation) error {
	p.session = activation.Session
	p.service = activation.Service
	p.serviceID = activation.ServiceID
	return p.impl.Activate(activation, p)
}
func (p *stubLogManager) OnTerminate() {
	p.impl.OnTerminate()
}
func (p *stubLogManager) Receive(msg *net.Message, from bus.Channel) error {
	// action dispatch
	switch msg.Header.Action {
	case 100:
		return p.Log(msg, from)
	case 101:
		return p.CreateListener(msg, from)
	case 102:
		return p.GetListener(msg, from)
	case 103:
		return p.AddProvider(msg, from)
	case 104:
		return p.RemoveProvider(msg, from)
	default:
		return from.SendError(msg, bus.ErrActionNotFound)
	}
}
func (p *stubLogManager) onPropertyChange(name string, data []byte) error {
	switch name {
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}
func (p *stubLogManager) Log(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	messages, err := func() (b []LogMessage, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]LogMessage, size)
		for i := 0; i < int(size); i++ {
			b[i], err = readLogMessage(buf)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read messages: %s", err))
	}
	callErr := p.impl.Log(messages)

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
func (p *stubLogManager) CreateListener(msg *net.Message, c bus.Channel) error {
	ret, callErr := p.impl.CreateListener()

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
func (p *stubLogManager) GetListener(msg *net.Message, c bus.Channel) error {
	ret, callErr := p.impl.GetListener()

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
func (p *stubLogManager) AddProvider(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	source, err := func() (LogProviderProxy, error) {
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
		return MakeLogProvider(p.session, proxy), nil
	}()
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read source: %s", err))
	}
	ret, callErr := p.impl.AddProvider(source)

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := basic.WriteInt32(ret, &out)
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubLogManager) RemoveProvider(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	sourceID, err := basic.ReadInt32(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read sourceID: %s", err))
	}
	callErr := p.impl.RemoveProvider(sourceID)

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
func (p *stubLogManager) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "LogManager",
		Methods: map[uint32]object.MetaMethod{
			100: {
				Name:                "log",
				ParametersSignature: "([(s(i)<LogLevel,level>sssI(L)<TimePoint,ns>(L)<TimePoint,ns>)<LogMessage,source,level,category,location,message,id,date,systemDate>])",
				ReturnSignature:     "v",
				Uid:                 100,
			},
			101: {
				Name:                "createListener",
				ParametersSignature: "()",
				ReturnSignature:     "o",
				Uid:                 101,
			},
			102: {
				Name:                "getListener",
				ParametersSignature: "()",
				ReturnSignature:     "o",
				Uid:                 102,
			},
			103: {
				Name:                "addProvider",
				ParametersSignature: "(o)",
				ReturnSignature:     "i",
				Uid:                 103,
			},
			104: {
				Name:                "removeProvider",
				ParametersSignature: "(i)",
				ReturnSignature:     "v",
				Uid:                 104,
			},
		},
		Properties: map[uint32]object.MetaProperty{},
		Signals:    map[uint32]object.MetaSignal{},
	}
}

// LogLevel is serializable
type LogLevel struct {
	Level int32
}

// readLogLevel unmarshalls LogLevel
func readLogLevel(r io.Reader) (s LogLevel, err error) {
	if s.Level, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("read Level field: %s", err)
	}
	return s, nil
}

// writeLogLevel marshalls LogLevel
func writeLogLevel(s LogLevel, w io.Writer) (err error) {
	if err := basic.WriteInt32(s.Level, w); err != nil {
		return fmt.Errorf("write Level field: %s", err)
	}
	return nil
}

// TimePoint is serializable
type TimePoint struct {
	Ns uint64
}

// readTimePoint unmarshalls TimePoint
func readTimePoint(r io.Reader) (s TimePoint, err error) {
	if s.Ns, err = basic.ReadUint64(r); err != nil {
		return s, fmt.Errorf("read Ns field: %s", err)
	}
	return s, nil
}

// writeTimePoint marshalls TimePoint
func writeTimePoint(s TimePoint, w io.Writer) (err error) {
	if err := basic.WriteUint64(s.Ns, w); err != nil {
		return fmt.Errorf("write Ns field: %s", err)
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

// readLogMessage unmarshalls LogMessage
func readLogMessage(r io.Reader) (s LogMessage, err error) {
	if s.Source, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Source field: %s", err)
	}
	if s.Level, err = readLogLevel(r); err != nil {
		return s, fmt.Errorf("read Level field: %s", err)
	}
	if s.Category, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Category field: %s", err)
	}
	if s.Location, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Location field: %s", err)
	}
	if s.Message, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Message field: %s", err)
	}
	if s.Id, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read Id field: %s", err)
	}
	if s.Date, err = readTimePoint(r); err != nil {
		return s, fmt.Errorf("read Date field: %s", err)
	}
	if s.SystemDate, err = readTimePoint(r); err != nil {
		return s, fmt.Errorf("read SystemDate field: %s", err)
	}
	return s, nil
}

// writeLogMessage marshalls LogMessage
func writeLogMessage(s LogMessage, w io.Writer) (err error) {
	if err := basic.WriteString(s.Source, w); err != nil {
		return fmt.Errorf("write Source field: %s", err)
	}
	if err := writeLogLevel(s.Level, w); err != nil {
		return fmt.Errorf("write Level field: %s", err)
	}
	if err := basic.WriteString(s.Category, w); err != nil {
		return fmt.Errorf("write Category field: %s", err)
	}
	if err := basic.WriteString(s.Location, w); err != nil {
		return fmt.Errorf("write Location field: %s", err)
	}
	if err := basic.WriteString(s.Message, w); err != nil {
		return fmt.Errorf("write Message field: %s", err)
	}
	if err := basic.WriteUint32(s.Id, w); err != nil {
		return fmt.Errorf("write Id field: %s", err)
	}
	if err := writeTimePoint(s.Date, w); err != nil {
		return fmt.Errorf("write Date field: %s", err)
	}
	if err := writeTimePoint(s.SystemDate, w); err != nil {
		return fmt.Errorf("write SystemDate field: %s", err)
	}
	return nil
}

// LogProviderProxy represents a proxy object to the service
type LogProviderProxy interface {
	SetVerbosity(level LogLevel) error
	SetCategory(category string, level LogLevel) error
	ClearAndSet(filters map[string]LogLevel) error
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
	WithContext(ctx context.Context) LogProviderProxy
}

// proxyLogProvider implements LogProviderProxy
type proxyLogProvider struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeLogProvider returns a specialized proxy.
func MakeLogProvider(sess bus.Session, proxy bus.Proxy) LogProviderProxy {
	return &proxyLogProvider{bus.MakeObject(proxy), sess}
}

// LogProvider returns a proxy to a remote service
func LogProvider(session bus.Session) (LogProviderProxy, error) {
	proxy, err := session.Proxy("LogProvider", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeLogProvider(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyLogProvider) WithContext(ctx context.Context) LogProviderProxy {
	return MakeLogProvider(p.session, p.Proxy().WithContext(ctx))
}

// SetVerbosity calls the remote procedure
func (p *proxyLogProvider) SetVerbosity(level LogLevel) error {
	var err error
	var buf bytes.Buffer
	if err = writeLogLevel(level, &buf); err != nil {
		return fmt.Errorf("serialize level: %s", err)
	}
	_, err = p.Proxy().Call("setVerbosity", buf.Bytes())
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
		return fmt.Errorf("serialize category: %s", err)
	}
	if err = writeLogLevel(level, &buf); err != nil {
		return fmt.Errorf("serialize level: %s", err)
	}
	_, err = p.Proxy().Call("setCategory", buf.Bytes())
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
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range filters {
			err = basic.WriteString(k, &buf)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = writeLogLevel(v, &buf)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize filters: %s", err)
	}
	_, err = p.Proxy().Call("clearAndSet", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearAndSet failed: %s", err)
	}
	return nil
}

// LogListenerProxy represents a proxy object to the service
type LogListenerProxy interface {
	SetLevel(level LogLevel) error
	AddFilter(category string, level LogLevel) error
	ClearFilters() error
	SubscribeOnLogMessage() (unsubscribe func(), updates chan LogMessage, err error)
	SubscribeOnLogMessages() (unsubscribe func(), updates chan []LogMessage, err error)
	SubscribeOnLogMessagesWithBacklog() (unsubscribe func(), updates chan []LogMessage, err error)
	GetLogLevel() (LogLevel, error)
	SetLogLevel(LogLevel) error
	SubscribeLogLevel() (unsubscribe func(), updates chan LogLevel, err error)
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
	WithContext(ctx context.Context) LogListenerProxy
}

// proxyLogListener implements LogListenerProxy
type proxyLogListener struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeLogListener returns a specialized proxy.
func MakeLogListener(sess bus.Session, proxy bus.Proxy) LogListenerProxy {
	return &proxyLogListener{bus.MakeObject(proxy), sess}
}

// LogListener returns a proxy to a remote service
func LogListener(session bus.Session) (LogListenerProxy, error) {
	proxy, err := session.Proxy("LogListener", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeLogListener(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyLogListener) WithContext(ctx context.Context) LogListenerProxy {
	return MakeLogListener(p.session, p.Proxy().WithContext(ctx))
}

// SetLevel calls the remote procedure
func (p *proxyLogListener) SetLevel(level LogLevel) error {
	var err error
	var buf bytes.Buffer
	if err = writeLogLevel(level, &buf); err != nil {
		return fmt.Errorf("serialize level: %s", err)
	}
	_, err = p.Proxy().Call("setLevel", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setLevel failed: %s", err)
	}
	return nil
}

// AddFilter calls the remote procedure
func (p *proxyLogListener) AddFilter(category string, level LogLevel) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(category, &buf); err != nil {
		return fmt.Errorf("serialize category: %s", err)
	}
	if err = writeLogLevel(level, &buf); err != nil {
		return fmt.Errorf("serialize level: %s", err)
	}
	_, err = p.Proxy().Call("addFilter", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call addFilter failed: %s", err)
	}
	return nil
}

// ClearFilters calls the remote procedure
func (p *proxyLogListener) ClearFilters() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Proxy().Call("clearFilters", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearFilters failed: %s", err)
	}
	return nil
}

// SubscribeOnLogMessage subscribe to a remote property
func (p *proxyLogListener) SubscribeOnLogMessage() (func(), chan LogMessage, error) {
	propertyID, err := p.Proxy().SignalID("onLogMessage")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "onLogMessage", err)
	}
	ch := make(chan LogMessage)
	cancel, chPay, err := p.Proxy().SubscribeID(propertyID)
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
			e, err := readLogMessage(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// SubscribeOnLogMessages subscribe to a remote property
func (p *proxyLogListener) SubscribeOnLogMessages() (func(), chan []LogMessage, error) {
	propertyID, err := p.Proxy().SignalID("onLogMessages")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "onLogMessages", err)
	}
	ch := make(chan []LogMessage)
	cancel, chPay, err := p.Proxy().SubscribeID(propertyID)
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
			e, err := func() (b []LogMessage, err error) {
				size, err := basic.ReadUint32(buf)
				if err != nil {
					return b, fmt.Errorf("read slice size: %s", err)
				}
				b = make([]LogMessage, size)
				for i := 0; i < int(size); i++ {
					b[i], err = readLogMessage(buf)
					if err != nil {
						return b, fmt.Errorf("read slice value: %s", err)
					}
				}
				return b, nil
			}()
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// SubscribeOnLogMessagesWithBacklog subscribe to a remote property
func (p *proxyLogListener) SubscribeOnLogMessagesWithBacklog() (func(), chan []LogMessage, error) {
	propertyID, err := p.Proxy().SignalID("onLogMessagesWithBacklog")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "onLogMessagesWithBacklog", err)
	}
	ch := make(chan []LogMessage)
	cancel, chPay, err := p.Proxy().SubscribeID(propertyID)
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
			e, err := func() (b []LogMessage, err error) {
				size, err := basic.ReadUint32(buf)
				if err != nil {
					return b, fmt.Errorf("read slice size: %s", err)
				}
				b = make([]LogMessage, size)
				for i := 0; i < int(size); i++ {
					b[i], err = readLogMessage(buf)
					if err != nil {
						return b, fmt.Errorf("read slice value: %s", err)
					}
				}
				return b, nil
			}()
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// GetLogLevel updates the property value
func (p *proxyLogListener) GetLogLevel() (ret LogLevel, err error) {
	name := value.String("logLevel")
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
	ret, err = readLogLevel(&buf)
	return ret, err
}

// SetLogLevel updates the property value
func (p *proxyLogListener) SetLogLevel(update LogLevel) error {
	name := value.String("logLevel")
	var buf bytes.Buffer
	err := writeLogLevel(update, &buf)
	if err != nil {
		return fmt.Errorf("marshall error: %s", err)
	}
	val := value.Opaque("(i)<LogLevel,level>", buf.Bytes())
	return p.SetProperty(name, val)
}

// SubscribeLogLevel subscribe to a remote property
func (p *proxyLogListener) SubscribeLogLevel() (func(), chan LogLevel, error) {
	propertyID, err := p.Proxy().PropertyID("logLevel")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "logLevel", err)
	}
	ch := make(chan LogLevel)
	cancel, chPay, err := p.Proxy().SubscribeID(propertyID)
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
			e, err := readLogLevel(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// LogManagerProxy represents a proxy object to the service
type LogManagerProxy interface {
	Log(messages []LogMessage) error
	CreateListener() (LogListenerProxy, error)
	GetListener() (LogListenerProxy, error)
	AddProvider(source LogProviderProxy) (int32, error)
	RemoveProvider(sourceID int32) error
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
	WithContext(ctx context.Context) LogManagerProxy
}

// proxyLogManager implements LogManagerProxy
type proxyLogManager struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeLogManager returns a specialized proxy.
func MakeLogManager(sess bus.Session, proxy bus.Proxy) LogManagerProxy {
	return &proxyLogManager{bus.MakeObject(proxy), sess}
}

// LogManager returns a proxy to a remote service
func LogManager(session bus.Session) (LogManagerProxy, error) {
	proxy, err := session.Proxy("LogManager", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeLogManager(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyLogManager) WithContext(ctx context.Context) LogManagerProxy {
	return MakeLogManager(p.session, p.Proxy().WithContext(ctx))
}

// Log calls the remote procedure
func (p *proxyLogManager) Log(messages []LogMessage) error {
	var err error
	var buf bytes.Buffer
	if err = func() error {
		err := basic.WriteUint32(uint32(len(messages)), &buf)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range messages {
			err = writeLogMessage(v, &buf)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize messages: %s", err)
	}
	_, err = p.Proxy().Call("log", buf.Bytes())
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
	response, err := p.Proxy().Call("createListener", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call createListener failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (LogListenerProxy, error) {
		ref, err := object.ReadObjectReference(resp)
		if err != nil {
			return nil, fmt.Errorf("get meta: %s", err)
		}
		proxy, err := p.session.Object(ref)
		if err != nil {
			return nil, fmt.Errorf("get proxy: %s", err)
		}
		return MakeLogListener(p.session, proxy), nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse createListener response: %s", err)
	}
	return ret, nil
}

// GetListener calls the remote procedure
func (p *proxyLogManager) GetListener() (LogListenerProxy, error) {
	var err error
	var ret LogListenerProxy
	var buf bytes.Buffer
	response, err := p.Proxy().Call("getListener", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getListener failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (LogListenerProxy, error) {
		ref, err := object.ReadObjectReference(resp)
		if err != nil {
			return nil, fmt.Errorf("get meta: %s", err)
		}
		proxy, err := p.session.Object(ref)
		if err != nil {
			return nil, fmt.Errorf("get proxy: %s", err)
		}
		return MakeLogListener(p.session, proxy), nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getListener response: %s", err)
	}
	return ret, nil
}

// AddProvider calls the remote procedure
func (p *proxyLogManager) AddProvider(source LogProviderProxy) (int32, error) {
	var err error
	var ret int32
	var buf bytes.Buffer
	if err = func() error {
		meta, err := source.MetaObject(source.Proxy().ObjectID())
		if err != nil {
			return fmt.Errorf("get meta: %s", err)
		}
		ref := object.ObjectReference{
			MetaObject: meta,
			ServiceID:  source.Proxy().ServiceID(),
			ObjectID:   source.Proxy().ObjectID(),
		}
		return object.WriteObjectReference(ref, &buf)
	}(); err != nil {
		return ret, fmt.Errorf("serialize source: %s", err)
	}
	response, err := p.Proxy().Call("addProvider", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call addProvider failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadInt32(resp)
	if err != nil {
		return ret, fmt.Errorf("parse addProvider response: %s", err)
	}
	return ret, nil
}

// RemoveProvider calls the remote procedure
func (p *proxyLogManager) RemoveProvider(sourceID int32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteInt32(sourceID, &buf); err != nil {
		return fmt.Errorf("serialize sourceID: %s", err)
	}
	_, err = p.Proxy().Call("removeProvider", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call removeProvider failed: %s", err)
	}
	return nil
}
