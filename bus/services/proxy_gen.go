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

// Constructor gives access to remote services
type Constructor struct {
	session bus.Session
}

// Services gives access to the services constructor
func Services(s bus.Session) Constructor {
	return Constructor{session: s}
}

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

// ServiceDirectory is the abstract interface of the service
type ServiceDirectory interface {
	// Service calls the remote procedure
	Service(name string) (ServiceInfo, error)
	// Services calls the remote procedure
	Services() ([]ServiceInfo, error)
	// RegisterService calls the remote procedure
	RegisterService(info ServiceInfo) (uint32, error)
	// UnregisterService calls the remote procedure
	UnregisterService(serviceID uint32) error
	// ServiceReady calls the remote procedure
	ServiceReady(serviceID uint32) error
	// UpdateServiceInfo calls the remote procedure
	UpdateServiceInfo(info ServiceInfo) error
	// MachineId calls the remote procedure
	MachineId() (string, error)
	// _socketOfService calls the remote procedure
	_socketOfService(serviceID uint32) (object.ObjectReference, error)
	// SubscribeServiceAdded subscribe to a remote signal
	SubscribeServiceAdded() (unsubscribe func(), updates chan ServiceAdded, err error)
	// SubscribeServiceRemoved subscribe to a remote signal
	SubscribeServiceRemoved() (unsubscribe func(), updates chan ServiceRemoved, err error)
}

// ServiceDirectoryProxy represents a proxy object to the service
type ServiceDirectoryProxy interface {
	bus.ObjectProxy
	ServiceDirectory
	// WithContext returns a new proxy. Calls to this proxy can be
	// cancelled by the context
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
func (c Constructor) ServiceDirectory() (ServiceDirectoryProxy, error) {
	proxy, err := c.session.Proxy("ServiceDirectory", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeServiceDirectory(c.session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyServiceDirectory) WithContext(ctx context.Context) ServiceDirectoryProxy {
	return MakeServiceDirectory(p.session, bus.WithContext(p.FIXMEProxy(), ctx))
}

// Service calls the remote procedure
func (p *proxyServiceDirectory) Service(name string) (ServiceInfo, error) {
	var err error
	var ret ServiceInfo
	var buf bytes.Buffer
	if err = basic.WriteString(name, &buf); err != nil {
		return ret, fmt.Errorf("serialize name: %s", err)
	}
	response, err := p.FIXMEProxy().Call("service", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call service failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = readServiceInfo(resp)
	if err != nil {
		return ret, fmt.Errorf("parse service response: %s", err)
	}
	return ret, nil
}

// Services calls the remote procedure
func (p *proxyServiceDirectory) Services() ([]ServiceInfo, error) {
	var err error
	var ret []ServiceInfo
	var buf bytes.Buffer
	response, err := p.FIXMEProxy().Call("services", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call services failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []ServiceInfo, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]ServiceInfo, size)
		for i := 0; i < int(size); i++ {
			b[i], err = readServiceInfo(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse services response: %s", err)
	}
	return ret, nil
}

// RegisterService calls the remote procedure
func (p *proxyServiceDirectory) RegisterService(info ServiceInfo) (uint32, error) {
	var err error
	var ret uint32
	var buf bytes.Buffer
	if err = writeServiceInfo(info, &buf); err != nil {
		return ret, fmt.Errorf("serialize info: %s", err)
	}
	response, err := p.FIXMEProxy().Call("registerService", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerService failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadUint32(resp)
	if err != nil {
		return ret, fmt.Errorf("parse registerService response: %s", err)
	}
	return ret, nil
}

// UnregisterService calls the remote procedure
func (p *proxyServiceDirectory) UnregisterService(serviceID uint32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteUint32(serviceID, &buf); err != nil {
		return fmt.Errorf("serialize serviceID: %s", err)
	}
	_, err = p.FIXMEProxy().Call("unregisterService", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterService failed: %s", err)
	}
	return nil
}

// ServiceReady calls the remote procedure
func (p *proxyServiceDirectory) ServiceReady(serviceID uint32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteUint32(serviceID, &buf); err != nil {
		return fmt.Errorf("serialize serviceID: %s", err)
	}
	_, err = p.FIXMEProxy().Call("serviceReady", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call serviceReady failed: %s", err)
	}
	return nil
}

// UpdateServiceInfo calls the remote procedure
func (p *proxyServiceDirectory) UpdateServiceInfo(info ServiceInfo) error {
	var err error
	var buf bytes.Buffer
	if err = writeServiceInfo(info, &buf); err != nil {
		return fmt.Errorf("serialize info: %s", err)
	}
	_, err = p.FIXMEProxy().Call("updateServiceInfo", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call updateServiceInfo failed: %s", err)
	}
	return nil
}

// MachineId calls the remote procedure
func (p *proxyServiceDirectory) MachineId() (string, error) {
	var err error
	var ret string
	var buf bytes.Buffer
	response, err := p.FIXMEProxy().Call("machineId", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call machineId failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadString(resp)
	if err != nil {
		return ret, fmt.Errorf("parse machineId response: %s", err)
	}
	return ret, nil
}

// _socketOfService calls the remote procedure
func (p *proxyServiceDirectory) _socketOfService(serviceID uint32) (object.ObjectReference, error) {
	var err error
	var ret object.ObjectReference
	var buf bytes.Buffer
	if err = basic.WriteUint32(serviceID, &buf); err != nil {
		return ret, fmt.Errorf("serialize serviceID: %s", err)
	}
	response, err := p.FIXMEProxy().Call("_socketOfService", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _socketOfService failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = object.ReadObjectReference(resp)
	if err != nil {
		return ret, fmt.Errorf("parse _socketOfService response: %s", err)
	}
	return ret, nil
}

// SubscribeServiceAdded subscribe to a remote property
func (p *proxyServiceDirectory) SubscribeServiceAdded() (func(), chan ServiceAdded, error) {
	propertyID, err := p.FIXMEProxy().SignalID("serviceAdded")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "serviceAdded", err)
	}
	ch := make(chan ServiceAdded)
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
	propertyID, err := p.FIXMEProxy().SignalID("serviceRemoved")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "serviceRemoved", err)
	}
	ch := make(chan ServiceRemoved)
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

// LogProvider is the abstract interface of the service
type LogProvider interface {
	// SetVerbosity calls the remote procedure
	SetVerbosity(level LogLevel) error
	// SetCategory calls the remote procedure
	SetCategory(category string, level LogLevel) error
	// ClearAndSet calls the remote procedure
	ClearAndSet(filters map[string]int32) error
}

// LogProviderProxy represents a proxy object to the service
type LogProviderProxy interface {
	bus.ObjectProxy
	LogProvider
	// WithContext returns a new proxy. Calls to this proxy can be
	// cancelled by the context
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
func (c Constructor) LogProvider() (LogProviderProxy, error) {
	proxy, err := c.session.Proxy("LogProvider", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeLogProvider(c.session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyLogProvider) WithContext(ctx context.Context) LogProviderProxy {
	return MakeLogProvider(p.session, bus.WithContext(p.FIXMEProxy(), ctx))
}

// SetVerbosity calls the remote procedure
func (p *proxyLogProvider) SetVerbosity(level LogLevel) error {
	var err error
	var buf bytes.Buffer
	if err = writeLogLevel(level, &buf); err != nil {
		return fmt.Errorf("serialize level: %s", err)
	}
	_, err = p.FIXMEProxy().Call("setVerbosity", buf.Bytes())
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
	_, err = p.FIXMEProxy().Call("setCategory", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setCategory failed: %s", err)
	}
	return nil
}

// ClearAndSet calls the remote procedure
func (p *proxyLogProvider) ClearAndSet(filters map[string]int32) error {
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
			err = basic.WriteInt32(v, &buf)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize filters: %s", err)
	}
	_, err = p.FIXMEProxy().Call("clearAndSet", buf.Bytes())
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

// LogListenerProxy represents a proxy object to the service
type LogListenerProxy interface {
	bus.ObjectProxy
	LogListener
	// WithContext returns a new proxy. Calls to this proxy can be
	// cancelled by the context
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
func (c Constructor) LogListener() (LogListenerProxy, error) {
	proxy, err := c.session.Proxy("LogListener", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeLogListener(c.session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyLogListener) WithContext(ctx context.Context) LogListenerProxy {
	return MakeLogListener(p.session, bus.WithContext(p.FIXMEProxy(), ctx))
}

// SetCategory calls the remote procedure
func (p *proxyLogListener) SetCategory(category string, level LogLevel) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(category, &buf); err != nil {
		return fmt.Errorf("serialize category: %s", err)
	}
	if err = writeLogLevel(level, &buf); err != nil {
		return fmt.Errorf("serialize level: %s", err)
	}
	_, err = p.FIXMEProxy().Call("setCategory", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setCategory failed: %s", err)
	}
	return nil
}

// ClearFilters calls the remote procedure
func (p *proxyLogListener) ClearFilters() error {
	var err error
	var buf bytes.Buffer
	_, err = p.FIXMEProxy().Call("clearFilters", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearFilters failed: %s", err)
	}
	return nil
}

// SubscribeOnLogMessage subscribe to a remote property
func (p *proxyLogListener) SubscribeOnLogMessage() (func(), chan LogMessage, error) {
	propertyID, err := p.FIXMEProxy().SignalID("onLogMessage")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "onLogMessage", err)
	}
	ch := make(chan LogMessage)
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
	propertyID, err := p.FIXMEProxy().PropertyID("verbosity")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "verbosity", err)
	}
	ch := make(chan LogLevel)
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
	propertyID, err := p.FIXMEProxy().PropertyID("filters")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "filters", err)
	}
	ch := make(chan map[string]int32)
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
	RemoveProvider(providerID int32) error
}

// LogManagerProxy represents a proxy object to the service
type LogManagerProxy interface {
	bus.ObjectProxy
	LogManager
	// WithContext returns a new proxy. Calls to this proxy can be
	// cancelled by the context
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
func (c Constructor) LogManager() (LogManagerProxy, error) {
	proxy, err := c.session.Proxy("LogManager", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeLogManager(c.session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyLogManager) WithContext(ctx context.Context) LogManagerProxy {
	return MakeLogManager(p.session, bus.WithContext(p.FIXMEProxy(), ctx))
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
	_, err = p.FIXMEProxy().Call("log", buf.Bytes())
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
	response, err := p.FIXMEProxy().Call("createListener", buf.Bytes())
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
	response, err := p.FIXMEProxy().Call("getListener", buf.Bytes())
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
		meta, err := source.MetaObject(source.FIXMEProxy().ObjectID())
		if err != nil {
			return fmt.Errorf("get meta: %s", err)
		}
		ref := object.ObjectReference{
			MetaObject: meta,
			ServiceID:  source.FIXMEProxy().ServiceID(),
			ObjectID:   source.FIXMEProxy().ObjectID(),
		}
		return object.WriteObjectReference(ref, &buf)
	}(); err != nil {
		return ret, fmt.Errorf("serialize source: %s", err)
	}
	response, err := p.FIXMEProxy().Call("addProvider", buf.Bytes())
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
func (p *proxyLogManager) RemoveProvider(providerID int32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteInt32(providerID, &buf); err != nil {
		return fmt.Errorf("serialize providerID: %s", err)
	}
	_, err = p.FIXMEProxy().Call("removeProvider", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call removeProvider failed: %s", err)
	}
	return nil
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

// ALAnimatedSpeech is the abstract interface of the service
type ALAnimatedSpeech interface {
	// Say calls the remote procedure
	Say(text string) error
	// IsBodyTalkEnabled calls the remote procedure
	IsBodyTalkEnabled() (bool, error)
	// IsBodyLanguageEnabled calls the remote procedure
	IsBodyLanguageEnabled() (bool, error)
	// SetBodyTalkEnabled calls the remote procedure
	SetBodyTalkEnabled(enable bool) error
	// SetBodyLanguageEnabled calls the remote procedure
	SetBodyLanguageEnabled(enable bool) error
}

// ALAnimatedSpeechProxy represents a proxy object to the service
type ALAnimatedSpeechProxy interface {
	bus.ObjectProxy
	ALAnimatedSpeech
	// WithContext returns a new proxy. Calls to this proxy can be
	// cancelled by the context
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
func (c Constructor) ALAnimatedSpeech() (ALAnimatedSpeechProxy, error) {
	proxy, err := c.session.Proxy("ALAnimatedSpeech", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return MakeALAnimatedSpeech(c.session, proxy), nil
}

// WithContext bound future calls to the context deadline and cancellation
func (p *proxyALAnimatedSpeech) WithContext(ctx context.Context) ALAnimatedSpeechProxy {
	return MakeALAnimatedSpeech(p.session, bus.WithContext(p.FIXMEProxy(), ctx))
}

// Say calls the remote procedure
func (p *proxyALAnimatedSpeech) Say(text string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(text, &buf); err != nil {
		return fmt.Errorf("serialize text: %s", err)
	}
	_, err = p.FIXMEProxy().Call("say", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call say failed: %s", err)
	}
	return nil
}

// IsBodyTalkEnabled calls the remote procedure
func (p *proxyALAnimatedSpeech) IsBodyTalkEnabled() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.FIXMEProxy().Call("isBodyTalkEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isBodyTalkEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse isBodyTalkEnabled response: %s", err)
	}
	return ret, nil
}

// IsBodyLanguageEnabled calls the remote procedure
func (p *proxyALAnimatedSpeech) IsBodyLanguageEnabled() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.FIXMEProxy().Call("isBodyLanguageEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isBodyLanguageEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse isBodyLanguageEnabled response: %s", err)
	}
	return ret, nil
}

// SetBodyTalkEnabled calls the remote procedure
func (p *proxyALAnimatedSpeech) SetBodyTalkEnabled(enable bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(enable, &buf); err != nil {
		return fmt.Errorf("serialize enable: %s", err)
	}
	_, err = p.FIXMEProxy().Call("setBodyTalkEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setBodyTalkEnabled failed: %s", err)
	}
	return nil
}

// SetBodyLanguageEnabled calls the remote procedure
func (p *proxyALAnimatedSpeech) SetBodyLanguageEnabled(enable bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(enable, &buf); err != nil {
		return fmt.Errorf("serialize enable: %s", err)
	}
	_, err = p.FIXMEProxy().Call("setBodyLanguageEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setBodyLanguageEnabled failed: %s", err)
	}
	return nil
}
