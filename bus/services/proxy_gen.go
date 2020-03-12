// Package services contains a generated proxy
// .

package services

import (
	"bytes"
	"context"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
	"io"
	"log"
)

// ServiceAdded is serializable
type ServiceAdded struct {
	ServiceID uint32
	Name      string
}

// readServiceAdded unmarshalls ServiceAdded
func readServiceAdded(r io.Reader) (s ServiceAdded, err error) {
	if s.ServiceID, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read ServiceID field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	return s, nil
}

// writeServiceAdded marshalls ServiceAdded
func writeServiceAdded(s ServiceAdded, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.ServiceID, w); err != nil {
		return fmt.Errorf("write ServiceID field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	return nil
}

// ServiceRemoved is serializable
type ServiceRemoved struct {
	ServiceID uint32
	Name      string
}

// readServiceRemoved unmarshalls ServiceRemoved
func readServiceRemoved(r io.Reader) (s ServiceRemoved, err error) {
	if s.ServiceID, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read ServiceID field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	return s, nil
}

// writeServiceRemoved marshalls ServiceRemoved
func writeServiceRemoved(s ServiceRemoved, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.ServiceID, w); err != nil {
		return fmt.Errorf("write ServiceID field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	return nil
}

// ServiceDirectoryProxy represents a proxy object to the service
type ServiceDirectoryProxy interface {
	Service(name string) (ServiceInfo, error)
	Services() ([]ServiceInfo, error)
	RegisterService(info ServiceInfo) (uint32, error)
	UnregisterService(serviceID uint32) error
	ServiceReady(serviceID uint32) error
	UpdateServiceInfo(info ServiceInfo) error
	MachineId() (string, error)
	_socketOfService(serviceID uint32) (object.ObjectReference, error)
	SubscribeServiceAdded() (unsubscribe func(), updates chan ServiceAdded, err error)
	SubscribeServiceRemoved() (unsubscribe func(), updates chan ServiceRemoved, err error)
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
	WithContext(ctx context.Context) ServiceDirectoryProxy
}

// proxyServiceDirectory implements ServiceDirectoryProxy
type proxyServiceDirectory struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeServiceDirectory returns a specialized proxy.
func MakeServiceDirectory(sess bus.Session, proxy bus.Proxy) ServiceDirectoryProxy {
	return &proxyServiceDirectory{bus.MakeObject(proxy), sess}
}

// ServiceDirectory returns a proxy to a remote service
func ServiceDirectory(session bus.Session) (ServiceDirectoryProxy, error) {
	proxy, err := session.Proxy("ServiceDirectory", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeServiceDirectory(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyServiceDirectory) WithContext(ctx context.Context) ServiceDirectoryProxy {
	return MakeServiceDirectory(p.session, p.Proxy().WithContext(ctx))
}

// Service calls the remote procedure
func (p *proxyServiceDirectory) Service(name string) (ServiceInfo, error) {
	var ret ServiceInfo
	args := bus.NewParams("(s)", name)
	resp := bus.NewResponse("(sIsI[s]ss)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId,objectUid>", &ret)
	err := p.Proxy().Call2("service", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call service failed: %s", err)
	}
	return ret, nil
}

// Services calls the remote procedure
func (p *proxyServiceDirectory) Services() ([]ServiceInfo, error) {
	var ret []ServiceInfo
	args := bus.NewParams("()")
	resp := bus.NewResponse("[(sIsI[s]ss)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId,objectUid>]", &ret)
	err := p.Proxy().Call2("services", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call services failed: %s", err)
	}
	return ret, nil
}

// RegisterService calls the remote procedure
func (p *proxyServiceDirectory) RegisterService(info ServiceInfo) (uint32, error) {
	var ret uint32
	args := bus.NewParams("((sIsI[s]ss)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId,objectUid>)", info)
	resp := bus.NewResponse("I", &ret)
	err := p.Proxy().Call2("registerService", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call registerService failed: %s", err)
	}
	return ret, nil
}

// UnregisterService calls the remote procedure
func (p *proxyServiceDirectory) UnregisterService(serviceID uint32) error {
	var ret struct{}
	args := bus.NewParams("(I)", serviceID)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("unregisterService", args, resp)
	if err != nil {
		return fmt.Errorf("call unregisterService failed: %s", err)
	}
	return nil
}

// ServiceReady calls the remote procedure
func (p *proxyServiceDirectory) ServiceReady(serviceID uint32) error {
	var ret struct{}
	args := bus.NewParams("(I)", serviceID)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("serviceReady", args, resp)
	if err != nil {
		return fmt.Errorf("call serviceReady failed: %s", err)
	}
	return nil
}

// UpdateServiceInfo calls the remote procedure
func (p *proxyServiceDirectory) UpdateServiceInfo(info ServiceInfo) error {
	var ret struct{}
	args := bus.NewParams("((sIsI[s]ss)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId,objectUid>)", info)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("updateServiceInfo", args, resp)
	if err != nil {
		return fmt.Errorf("call updateServiceInfo failed: %s", err)
	}
	return nil
}

// MachineId calls the remote procedure
func (p *proxyServiceDirectory) MachineId() (string, error) {
	var ret string
	args := bus.NewParams("()")
	resp := bus.NewResponse("s", &ret)
	err := p.Proxy().Call2("machineId", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call machineId failed: %s", err)
	}
	return ret, nil
}

// _socketOfService calls the remote procedure
func (p *proxyServiceDirectory) _socketOfService(serviceID uint32) (object.ObjectReference, error) {
	var ret object.ObjectReference
	args := bus.NewParams("(I)", serviceID)
	resp := bus.NewResponse("o", &ret)
	err := p.Proxy().Call2("_socketOfService", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call _socketOfService failed: %s", err)
	}
	return ret, nil
}

// SubscribeServiceAdded subscribe to a remote property
func (p *proxyServiceDirectory) SubscribeServiceAdded() (func(), chan ServiceAdded, error) {
	signalID, err := p.Proxy().MetaObject().SignalID("serviceAdded", "(Is)<serviceAdded,serviceID,name>")
	if err != nil {
		return nil, nil, fmt.Errorf("%s not available: %s", "serviceAdded", err)
	}
	ch := make(chan ServiceAdded)
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
			e, err := readServiceAdded(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// SubscribeServiceRemoved subscribe to a remote property
func (p *proxyServiceDirectory) SubscribeServiceRemoved() (func(), chan ServiceRemoved, error) {
	signalID, err := p.Proxy().MetaObject().SignalID("serviceRemoved", "(Is)<serviceRemoved,serviceID,name>")
	if err != nil {
		return nil, nil, fmt.Errorf("%s not available: %s", "serviceRemoved", err)
	}
	ch := make(chan ServiceRemoved)
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
			e, err := readServiceRemoved(buf)
			if err != nil {
				log.Printf("unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// ServiceInfo is serializable
type ServiceInfo struct {
	Name      string
	ServiceId uint32
	MachineId string
	ProcessId uint32
	Endpoints []string
	SessionId string
	ObjectUid string
}

// readServiceInfo unmarshalls ServiceInfo
func readServiceInfo(r io.Reader) (s ServiceInfo, err error) {
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	if s.ServiceId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read ServiceId field: %s", err)
	}
	if s.MachineId, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read MachineId field: %s", err)
	}
	if s.ProcessId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read ProcessId field: %s", err)
	}
	if s.Endpoints, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(r)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("read Endpoints field: %s", err)
	}
	if s.SessionId, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read SessionId field: %s", err)
	}
	if s.ObjectUid, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read ObjectUid field: %s", err)
	}
	return s, nil
}

// writeServiceInfo marshalls ServiceInfo
func writeServiceInfo(s ServiceInfo, w io.Writer) (err error) {
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	if err := basic.WriteUint32(s.ServiceId, w); err != nil {
		return fmt.Errorf("write ServiceId field: %s", err)
	}
	if err := basic.WriteString(s.MachineId, w); err != nil {
		return fmt.Errorf("write MachineId field: %s", err)
	}
	if err := basic.WriteUint32(s.ProcessId, w); err != nil {
		return fmt.Errorf("write ProcessId field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Endpoints)), w)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range s.Endpoints {
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("write Endpoints field: %s", err)
	}
	if err := basic.WriteString(s.SessionId, w); err != nil {
		return fmt.Errorf("write SessionId field: %s", err)
	}
	if err := basic.WriteString(s.ObjectUid, w); err != nil {
		return fmt.Errorf("write ObjectUid field: %s", err)
	}
	return nil
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
	ClearAndSet(filters map[string]int32) error
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
	var ret struct{}
	args := bus.NewParams("((i)<LogLevel,level>)", level)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("setVerbosity", args, resp)
	if err != nil {
		return fmt.Errorf("call setVerbosity failed: %s", err)
	}
	return nil
}

// SetCategory calls the remote procedure
func (p *proxyLogProvider) SetCategory(category string, level LogLevel) error {
	var ret struct{}
	args := bus.NewParams("(s(i)<LogLevel,level>)", category, level)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("setCategory", args, resp)
	if err != nil {
		return fmt.Errorf("call setCategory failed: %s", err)
	}
	return nil
}

// ClearAndSet calls the remote procedure
func (p *proxyLogProvider) ClearAndSet(filters map[string]int32) error {
	var ret struct{}
	args := bus.NewParams("({si})", filters)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("clearAndSet", args, resp)
	if err != nil {
		return fmt.Errorf("call clearAndSet failed: %s", err)
	}
	return nil
}

// LogListenerProxy represents a proxy object to the service
type LogListenerProxy interface {
	SetCategory(category string, level LogLevel) error
	ClearFilters() error
	SubscribeOnLogMessage() (unsubscribe func(), updates chan LogMessage, err error)
	GetVerbosity() (LogLevel, error)
	SetVerbosity(LogLevel) error
	SubscribeVerbosity() (unsubscribe func(), updates chan LogLevel, err error)
	GetFilters() (map[string]int32, error)
	SetFilters(map[string]int32) error
	SubscribeFilters() (unsubscribe func(), updates chan map[string]int32, err error)
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

// SetCategory calls the remote procedure
func (p *proxyLogListener) SetCategory(category string, level LogLevel) error {
	var ret struct{}
	args := bus.NewParams("(s(i)<LogLevel,level>)", category, level)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("setCategory", args, resp)
	if err != nil {
		return fmt.Errorf("call setCategory failed: %s", err)
	}
	return nil
}

// ClearFilters calls the remote procedure
func (p *proxyLogListener) ClearFilters() error {
	var ret struct{}
	args := bus.NewParams("()")
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("clearFilters", args, resp)
	if err != nil {
		return fmt.Errorf("call clearFilters failed: %s", err)
	}
	return nil
}

// SubscribeOnLogMessage subscribe to a remote property
func (p *proxyLogListener) SubscribeOnLogMessage() (func(), chan LogMessage, error) {
	signalID, err := p.Proxy().MetaObject().SignalID("onLogMessage", "(s(i)<LogLevel,level>sssI(L)<TimePoint,ns>(L)<TimePoint,ns>)<LogMessage,source,level,category,location,message,id,date,systemDate>")
	if err != nil {
		return nil, nil, fmt.Errorf("%s not available: %s", "onLogMessage", err)
	}
	ch := make(chan LogMessage)
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
	ret, err = readLogLevel(&buf)
	return ret, err
}

// SetVerbosity updates the property value
func (p *proxyLogListener) SetVerbosity(update LogLevel) error {
	name := value.String("verbosity")
	var buf bytes.Buffer
	err := writeLogLevel(update, &buf)
	if err != nil {
		return fmt.Errorf("marshall error: %s", err)
	}
	val := value.Opaque("(i)<LogLevel,level>", buf.Bytes())
	return p.SetProperty(name, val)
}

// SubscribeVerbosity subscribe to a remote property
func (p *proxyLogListener) SubscribeVerbosity() (func(), chan LogLevel, error) {
	signalID, err := p.Proxy().MetaObject().PropertyID("verbosity", "(i)<LogLevel,level>")
	if err != nil {
		return nil, nil, fmt.Errorf("%s not available: %s", "verbosity", err)
	}
	ch := make(chan LogLevel)
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
			return m, fmt.Errorf("read map size: %s", err)
		}
		m = make(map[string]int32, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(&buf)
			if err != nil {
				return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
			}
			v, err := basic.ReadInt32(&buf)
			if err != nil {
				return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
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
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range update {
			err = basic.WriteString(k, &buf)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = basic.WriteInt32(v, &buf)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
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
	signalID, err := p.Proxy().MetaObject().PropertyID("filters", "{si}")
	if err != nil {
		return nil, nil, fmt.Errorf("%s not available: %s", "filters", err)
	}
	ch := make(chan map[string]int32)
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
			e, err := func() (m map[string]int32, err error) {
				size, err := basic.ReadUint32(buf)
				if err != nil {
					return m, fmt.Errorf("read map size: %s", err)
				}
				m = make(map[string]int32, size)
				for i := 0; i < int(size); i++ {
					k, err := basic.ReadString(buf)
					if err != nil {
						return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
					}
					v, err := basic.ReadInt32(buf)
					if err != nil {
						return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
					}
					m[k] = v
				}
				return m, nil
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

// LogManagerProxy represents a proxy object to the service
type LogManagerProxy interface {
	Log(messages []LogMessage) error
	CreateListener() (LogListenerProxy, error)
	GetListener() (LogListenerProxy, error)
	AddProvider(source LogProviderProxy) (int32, error)
	RemoveProvider(providerID int32) error
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
	var ret struct{}
	args := bus.NewParams("([(s(i)<LogLevel,level>sssI(L)<TimePoint,ns>(L)<TimePoint,ns>)<LogMessage,source,level,category,location,message,id,date,systemDate>])", messages)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("log", args, resp)
	if err != nil {
		return fmt.Errorf("call log failed: %s", err)
	}
	return nil
}

// CreateListener calls the remote procedure
func (p *proxyLogManager) CreateListener() (LogListenerProxy, error) {
	var ret LogListenerProxy
	args := bus.NewParams("()")
	var retRef object.ObjectReference
	resp := bus.NewResponse("o", &retRef)
	err := p.Proxy().Call2("createListener", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call createListener failed: %s", err)
	}
	proxy, err := p.session.Object(retRef)
	if err != nil {
		return nil, fmt.Errorf("proxy: %s", err)
	}
	ret = MakeLogListener(p.session, proxy)
	return ret, nil
}

// GetListener calls the remote procedure
func (p *proxyLogManager) GetListener() (LogListenerProxy, error) {
	var ret LogListenerProxy
	args := bus.NewParams("()")
	var retRef object.ObjectReference
	resp := bus.NewResponse("o", &retRef)
	err := p.Proxy().Call2("getListener", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call getListener failed: %s", err)
	}
	proxy, err := p.session.Object(retRef)
	if err != nil {
		return nil, fmt.Errorf("proxy: %s", err)
	}
	ret = MakeLogListener(p.session, proxy)
	return ret, nil
}

// AddProvider calls the remote procedure
func (p *proxyLogManager) AddProvider(source LogProviderProxy) (int32, error) {
	var ret int32
	args := bus.NewParams("(o)", bus.ObjectReference(source.Proxy()))
	resp := bus.NewResponse("i", &ret)
	err := p.Proxy().Call2("addProvider", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call addProvider failed: %s", err)
	}
	return ret, nil
}

// RemoveProvider calls the remote procedure
func (p *proxyLogManager) RemoveProvider(providerID int32) error {
	var ret struct{}
	args := bus.NewParams("(i)", providerID)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("removeProvider", args, resp)
	if err != nil {
		return fmt.Errorf("call removeProvider failed: %s", err)
	}
	return nil
}

// ALTextToSpeechProxy represents a proxy object to the service
type ALTextToSpeechProxy interface {
	Say(stringToSay string) error
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
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
func ALTextToSpeech(session bus.Session) (ALTextToSpeechProxy, error) {
	proxy, err := session.Proxy("ALTextToSpeech", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeALTextToSpeech(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyALTextToSpeech) WithContext(ctx context.Context) ALTextToSpeechProxy {
	return MakeALTextToSpeech(p.session, p.Proxy().WithContext(ctx))
}

// Say calls the remote procedure
func (p *proxyALTextToSpeech) Say(stringToSay string) error {
	var ret struct{}
	args := bus.NewParams("(s)", stringToSay)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("say", args, resp)
	if err != nil {
		return fmt.Errorf("call say failed: %s", err)
	}
	return nil
}

// ALAnimatedSpeechProxy represents a proxy object to the service
type ALAnimatedSpeechProxy interface {
	Say(text string) error
	IsBodyTalkEnabled() (bool, error)
	IsBodyLanguageEnabled() (bool, error)
	SetBodyTalkEnabled(enable bool) error
	SetBodyLanguageEnabled(enable bool) error
	// Generic methods shared by all objectsProxy
	bus.ObjectProxy
	// WithContext can be used cancellation and timeout
	WithContext(ctx context.Context) ALAnimatedSpeechProxy
}

// proxyALAnimatedSpeech implements ALAnimatedSpeechProxy
type proxyALAnimatedSpeech struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeALAnimatedSpeech returns a specialized proxy.
func MakeALAnimatedSpeech(sess bus.Session, proxy bus.Proxy) ALAnimatedSpeechProxy {
	return &proxyALAnimatedSpeech{bus.MakeObject(proxy), sess}
}

// ALAnimatedSpeech returns a proxy to a remote service
func ALAnimatedSpeech(session bus.Session) (ALAnimatedSpeechProxy, error) {
	proxy, err := session.Proxy("ALAnimatedSpeech", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeALAnimatedSpeech(session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyALAnimatedSpeech) WithContext(ctx context.Context) ALAnimatedSpeechProxy {
	return MakeALAnimatedSpeech(p.session, p.Proxy().WithContext(ctx))
}

// Say calls the remote procedure
func (p *proxyALAnimatedSpeech) Say(text string) error {
	var ret struct{}
	args := bus.NewParams("(s)", text)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("say", args, resp)
	if err != nil {
		return fmt.Errorf("call say failed: %s", err)
	}
	return nil
}

// IsBodyTalkEnabled calls the remote procedure
func (p *proxyALAnimatedSpeech) IsBodyTalkEnabled() (bool, error) {
	var ret bool
	args := bus.NewParams("()")
	resp := bus.NewResponse("b", &ret)
	err := p.Proxy().Call2("isBodyTalkEnabled", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call isBodyTalkEnabled failed: %s", err)
	}
	return ret, nil
}

// IsBodyLanguageEnabled calls the remote procedure
func (p *proxyALAnimatedSpeech) IsBodyLanguageEnabled() (bool, error) {
	var ret bool
	args := bus.NewParams("()")
	resp := bus.NewResponse("b", &ret)
	err := p.Proxy().Call2("isBodyLanguageEnabled", args, resp)
	if err != nil {
		return ret, fmt.Errorf("call isBodyLanguageEnabled failed: %s", err)
	}
	return ret, nil
}

// SetBodyTalkEnabled calls the remote procedure
func (p *proxyALAnimatedSpeech) SetBodyTalkEnabled(enable bool) error {
	var ret struct{}
	args := bus.NewParams("(b)", enable)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("setBodyTalkEnabled", args, resp)
	if err != nil {
		return fmt.Errorf("call setBodyTalkEnabled failed: %s", err)
	}
	return nil
}

// SetBodyLanguageEnabled calls the remote procedure
func (p *proxyALAnimatedSpeech) SetBodyLanguageEnabled(enable bool) error {
	var ret struct{}
	args := bus.NewParams("(b)", enable)
	resp := bus.NewResponse("v", &ret)
	err := p.Proxy().Call2("setBodyLanguageEnabled", args, resp)
	if err != nil {
		return fmt.Errorf("call setBodyLanguageEnabled failed: %s", err)
	}
	return nil
}
