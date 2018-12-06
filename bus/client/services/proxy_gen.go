// Package services contains a generated proxy
// File generated. DO NOT EDIT.
package services

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	object1 "github.com/lugu/qiloop/bus/client/object"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	"io"
	"log"
)

// ServicesConstructor gives access to remote services
type ServicesConstructor struct {
	session bus.Session
}

// Services gives access to the services constructor
func Services(s bus.Session) ServicesConstructor {
	return ServicesConstructor{session: s}
}

// ServiceDirectory is a proxy object to the remote service
type ServiceDirectory interface {
	object.Object
	bus.Proxy
	// Service calls the remote procedure
	Service(P0 string) (ServiceInfo, error)
	// Services calls the remote procedure
	Services() ([]ServiceInfo, error)
	// RegisterService calls the remote procedure
	RegisterService(P0 ServiceInfo) (uint32, error)
	// UnregisterService calls the remote procedure
	UnregisterService(P0 uint32) error
	// ServiceReady calls the remote procedure
	ServiceReady(P0 uint32) error
	// UpdateServiceInfo calls the remote procedure
	UpdateServiceInfo(P0 ServiceInfo) error
	// MachineId calls the remote procedure
	MachineId() (string, error)
	// SocketOfService calls the remote procedure
	SocketOfService(P0 uint32) (object.ObjectReference, error)
	// SignalServiceAdded subscribe to a remote signal
	SignalServiceAdded(cancel chan int) (chan struct {
		P0 uint32
		P1 string
	}, error)
	// SignalServiceRemoved subscribe to a remote signal
	SignalServiceRemoved(cancel chan int) (chan struct {
		P0 uint32
		P1 string
	}, error)
}

// LogManager is a proxy object to the remote service
type LogManager interface {
	object.Object
	bus.Proxy
	// Log calls the remote procedure
	Log(P0 []LogMessage) error
	// CreateListener calls the remote procedure
	CreateListener() (object.ObjectReference, error)
	// GetListener calls the remote procedure
	GetListener() (object.ObjectReference, error)
	// AddProvider calls the remote procedure
	AddProvider(P0 object.ObjectReference) (int32, error)
	// RemoveProvider calls the remote procedure
	RemoveProvider(P0 int32) error
}

// ServiceDirectoryProxy implements ServiceDirectory
type ServiceDirectoryProxy struct {
	object1.ObjectProxy
}

// NewServiceDirectory constructs ServiceDirectory
func NewServiceDirectory(ses bus.Session, obj uint32) (ServiceDirectory, error) {
	proxy, err := ses.Proxy("ServiceDirectory", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &ServiceDirectoryProxy{object1.ObjectProxy{proxy}}, nil
}

// ServiceDirectory retruns a proxy to a remote service
func (s ServicesConstructor) ServiceDirectory() (ServiceDirectory, error) {
	return NewServiceDirectory(s.session, 1)
}

// Service calls the remote procedure
func (p *ServiceDirectoryProxy) Service(P0 string) (ServiceInfo, error) {
	var err error
	var ret ServiceInfo
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("service", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call service failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = ReadServiceInfo(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse service response: %s", err)
	}
	return ret, nil
}

// Services calls the remote procedure
func (p *ServiceDirectoryProxy) Services() ([]ServiceInfo, error) {
	var err error
	var ret []ServiceInfo
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("services", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call services failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (b []ServiceInfo, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]ServiceInfo, size)
		for i := 0; i < int(size); i++ {
			b[i], err = ReadServiceInfo(buf)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse services response: %s", err)
	}
	return ret, nil
}

// RegisterService calls the remote procedure
func (p *ServiceDirectoryProxy) RegisterService(P0 ServiceInfo) (uint32, error) {
	var err error
	var ret uint32
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = WriteServiceInfo(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("registerService", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerService failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadUint32(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse registerService response: %s", err)
	}
	return ret, nil
}

// UnregisterService calls the remote procedure
func (p *ServiceDirectoryProxy) UnregisterService(P0 uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("unregisterService", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterService failed: %s", err)
	}
	return nil
}

// ServiceReady calls the remote procedure
func (p *ServiceDirectoryProxy) ServiceReady(P0 uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("serviceReady", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call serviceReady failed: %s", err)
	}
	return nil
}

// UpdateServiceInfo calls the remote procedure
func (p *ServiceDirectoryProxy) UpdateServiceInfo(P0 ServiceInfo) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = WriteServiceInfo(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("updateServiceInfo", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call updateServiceInfo failed: %s", err)
	}
	return nil
}

// MachineId calls the remote procedure
func (p *ServiceDirectoryProxy) MachineId() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("machineId", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call machineId failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse machineId response: %s", err)
	}
	return ret, nil
}

// SocketOfService calls the remote procedure
func (p *ServiceDirectoryProxy) SocketOfService(P0 uint32) (object.ObjectReference, error) {
	var err error
	var ret object.ObjectReference
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("_socketOfService", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _socketOfService failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = object.ReadObjectReference(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _socketOfService response: %s", err)
	}
	return ret, nil
}

// SignalServiceAdded subscribe to a remote signal
func (p *ServiceDirectoryProxy) SignalServiceAdded(cancel chan int) (chan struct {
	P0 uint32
	P1 string
}, error) {
	signalID, err := p.SignalUid("serviceAdded")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "serviceAdded", err)
	}

	handlerID := uint64(signalID)<<32 + 1 // FIXME: read it from proxy
	_, err = p.RegisterEvent(p.ObjectID(), signalID, handlerID)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "serviceAdded", err)
	}
	ch := make(chan struct {
		P0 uint32
		P1 string
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				p.UnregisterEvent(p.ObjectID(), signalID, handlerID)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 uint32
				P1 string
			}, err error) {
				s.P0, err = basic.ReadUint32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				s.P1, err = basic.ReadString(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				log.Printf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}

// SignalServiceRemoved subscribe to a remote signal
func (p *ServiceDirectoryProxy) SignalServiceRemoved(cancel chan int) (chan struct {
	P0 uint32
	P1 string
}, error) {
	signalID, err := p.SignalUid("serviceRemoved")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "serviceRemoved", err)
	}

	handlerID := uint64(signalID)<<32 + 1 // FIXME: read it from proxy
	_, err = p.RegisterEvent(p.ObjectID(), signalID, handlerID)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "serviceRemoved", err)
	}
	ch := make(chan struct {
		P0 uint32
		P1 string
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				// connection lost or cancellation.
				close(ch)
				p.UnregisterEvent(p.ObjectID(), signalID, handlerID)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 uint32
				P1 string
			}, err error) {
				s.P0, err = basic.ReadUint32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				s.P1, err = basic.ReadString(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				log.Printf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}

// LogManagerProxy implements LogManager
type LogManagerProxy struct {
	object1.ObjectProxy
}

// NewLogManager constructs LogManager
func NewLogManager(ses bus.Session, obj uint32) (LogManager, error) {
	proxy, err := ses.Proxy("LogManager", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &LogManagerProxy{object1.ObjectProxy{proxy}}, nil
}

// LogManager retruns a proxy to a remote service
func (s ServicesConstructor) LogManager() (LogManager, error) {
	return NewLogManager(s.session, 1)
}

// Log calls the remote procedure
func (p *LogManagerProxy) Log(P0 []LogMessage) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = func() error {
		err := basic.WriteUint32(uint32(len(P0)), buf)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range P0 {
			err = WriteLogMessage(v, buf)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("log", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call log failed: %s", err)
	}
	return nil
}

// CreateListener calls the remote procedure
func (p *LogManagerProxy) CreateListener() (object.ObjectReference, error) {
	var err error
	var ret object.ObjectReference
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("createListener", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call createListener failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = object.ReadObjectReference(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse createListener response: %s", err)
	}
	return ret, nil
}

// GetListener calls the remote procedure
func (p *LogManagerProxy) GetListener() (object.ObjectReference, error) {
	var err error
	var ret object.ObjectReference
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("getListener", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getListener failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = object.ReadObjectReference(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse getListener response: %s", err)
	}
	return ret, nil
}

// AddProvider calls the remote procedure
func (p *LogManagerProxy) AddProvider(P0 object.ObjectReference) (int32, error) {
	var err error
	var ret int32
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = object.WriteObjectReference(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("addProvider", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call addProvider failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadInt32(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse addProvider response: %s", err)
	}
	return ret, nil
}

// RemoveProvider calls the remote procedure
func (p *LogManagerProxy) RemoveProvider(P0 int32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteInt32(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("removeProvider", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call removeProvider failed: %s", err)
	}
	return nil
}

type ServiceInfo struct {
	Name      string
	ServiceId uint32
	MachineId string
	ProcessId uint32
	Endpoints []string
	SessionId string
}

func ReadServiceInfo(r io.Reader) (s ServiceInfo, err error) {
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	if s.ServiceId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read ServiceId field: " + err.Error())
	}
	if s.MachineId, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read MachineId field: " + err.Error())
	}
	if s.ProcessId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read ProcessId field: " + err.Error())
	}
	if s.Endpoints, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(r)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Endpoints field: " + err.Error())
	}
	if s.SessionId, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read SessionId field: " + err.Error())
	}
	return s, nil
}
func WriteServiceInfo(s ServiceInfo, w io.Writer) (err error) {
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
	}
	if err := basic.WriteUint32(s.ServiceId, w); err != nil {
		return fmt.Errorf("failed to write ServiceId field: " + err.Error())
	}
	if err := basic.WriteString(s.MachineId, w); err != nil {
		return fmt.Errorf("failed to write MachineId field: " + err.Error())
	}
	if err := basic.WriteUint32(s.ProcessId, w); err != nil {
		return fmt.Errorf("failed to write ProcessId field: " + err.Error())
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Endpoints)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.Endpoints {
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Endpoints field: " + err.Error())
	}
	if err := basic.WriteString(s.SessionId, w); err != nil {
		return fmt.Errorf("failed to write SessionId field: " + err.Error())
	}
	return nil
}

type LogMessage struct {
	Source     string
	Level      int32
	Category   string
	Location   string
	Message    string
	Id         uint32
	Date       uint64
	SystemDate uint64
}

func ReadLogMessage(r io.Reader) (s LogMessage, err error) {
	if s.Source, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Source field: " + err.Error())
	}
	if s.Level, err = basic.ReadInt32(r); err != nil {
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
	if s.Date, err = basic.ReadUint64(r); err != nil {
		return s, fmt.Errorf("failed to read Date field: " + err.Error())
	}
	if s.SystemDate, err = basic.ReadUint64(r); err != nil {
		return s, fmt.Errorf("failed to read SystemDate field: " + err.Error())
	}
	return s, nil
}
func WriteLogMessage(s LogMessage, w io.Writer) (err error) {
	if err := basic.WriteString(s.Source, w); err != nil {
		return fmt.Errorf("failed to write Source field: " + err.Error())
	}
	if err := basic.WriteInt32(s.Level, w); err != nil {
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
	if err := basic.WriteUint64(s.Date, w); err != nil {
		return fmt.Errorf("failed to write Date field: " + err.Error())
	}
	if err := basic.WriteUint64(s.SystemDate, w); err != nil {
		return fmt.Errorf("failed to write SystemDate field: " + err.Error())
	}
	return nil
}
