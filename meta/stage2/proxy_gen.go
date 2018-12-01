// file generated. DO NOT EDIT.
package stage2

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	object1 "github.com/lugu/qiloop/bus/client/object"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
	"io"
	"log"
)

type NewServices struct {
	session bus.Session
}

func Services(s bus.Session) NewServices {
	return NewServices{session: s}
}

type Server interface {
	bus.Proxy
	Authenticate(P0 map[string]value.Value) (map[string]value.Value, error)
}
type Object interface {
	object.Object
	bus.Proxy
}
type ServiceDirectory interface {
	object.Object
	bus.Proxy
	Service(P0 string) (ServiceInfo, error)
	Services() ([]ServiceInfo, error)
	RegisterService(P0 ServiceInfo) (uint32, error)
	UnregisterService(P0 uint32) error
	ServiceReady(P0 uint32) error
	UpdateServiceInfo(P0 ServiceInfo) error
	MachineId() (string, error)
	_socketOfService(P0 uint32) (object.ObjectReference, error)
	SignalServiceAdded(cancel chan int) (chan struct {
		P0 uint32
		P1 string
	}, error)
	SignalServiceRemoved(cancel chan int) (chan struct {
		P0 uint32
		P1 string
	}, error)
}
type ServerProxy struct {
	bus.Proxy
}

func NewServer(ses bus.Session, obj uint32) (Server, error) {
	proxy, err := ses.Proxy("Server", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &ServerProxy{proxy}, nil
}
func (s NewServices) Server() (Server, error) {
	return NewServer(s.session, 1)
}
func (p *ServerProxy) Authenticate(P0 map[string]value.Value) (map[string]value.Value, error) {
	var err error
	var ret map[string]value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = func() error {
		err := basic.WriteUint32(uint32(len(P0)), buf)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range P0 {
			err = basic.WriteString(k, buf)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = v.Write(buf)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("authenticate", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call authenticate failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (m map[string]value.Value, err error) {
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
		return ret, fmt.Errorf("failed to parse authenticate response: %s", err)
	}
	return ret, nil
}

type ObjectProxy struct {
	bus.Proxy
}

func NewObject(ses bus.Session, obj uint32) (Object, error) {
	proxy, err := ses.Proxy("Object", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &ObjectProxy{proxy}, nil
}
func (s NewServices) Object() (Object, error) {
	return NewObject(s.session, 1)
}
func (p *ObjectProxy) RegisterEvent(P0 uint32, P1 uint32, P2 uint64) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteUint32(P1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteUint64(P2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P2: %s", err)
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
func (p *ObjectProxy) UnregisterEvent(P0 uint32, P1 uint32, P2 uint64) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteUint32(P1, buf); err != nil {
		return fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteUint64(P2, buf); err != nil {
		return fmt.Errorf("failed to serialize P2: %s", err)
	}
	_, err = p.Call("unregisterEvent", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterEvent failed: %s", err)
	}
	return nil
}
func (p *ObjectProxy) MetaObject(P0 uint32) (object.MetaObject, error) {
	var err error
	var ret object.MetaObject
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
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
func (p *ObjectProxy) Terminate(P0 uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("terminate", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call terminate failed: %s", err)
	}
	return nil
}
func (p *ObjectProxy) Property(P0 value.Value) (value.Value, error) {
	var err error
	var ret value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = P0.Write(buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
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
func (p *ObjectProxy) SetProperty(P0 value.Value, P1 value.Value) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = P0.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = P1.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize P1: %s", err)
	}
	_, err = p.Call("setProperty", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setProperty failed: %s", err)
	}
	return nil
}
func (p *ObjectProxy) Properties() ([]string, error) {
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
func (p *ObjectProxy) RegisterEventWithSignature(P0 uint32, P1 uint32, P2 uint64, P3 string) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteUint32(P1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteUint64(P2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P2: %s", err)
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
func (p *ObjectProxy) SignalTraceObject(cancel chan int) (chan struct {
	P0 EventTrace
}, error) {
	signalID, err := p.SignalUid("traceObject")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "traceObject", err)
	}

	handlerID := uint64(signalID)<<32 + 1 // FIXME: read it from proxy
	_, err = p.RegisterEvent(p.ObjectID(), signalID, handlerID)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "traceObject", err)
	}
	ch := make(chan struct {
		P0 EventTrace
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
				p.UnregisterEvent(p.ObjectID(), signalID, handlerID)
				close(ch)
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 EventTrace
			}, err error) {
				s.P0, err = ReadEventTrace(buf)
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

type ServiceDirectoryProxy struct {
	object1.ObjectProxy
}

func NewServiceDirectory(ses bus.Session, obj uint32) (ServiceDirectory, error) {
	proxy, err := ses.Proxy("ServiceDirectory", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &ServiceDirectoryProxy{object1.ObjectProxy{proxy}}, nil
}
func (s NewServices) ServiceDirectory() (ServiceDirectory, error) {
	return NewServiceDirectory(s.session, 1)
}
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
func (p *ServiceDirectoryProxy) _socketOfService(P0 uint32) (object.ObjectReference, error) {
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
				p.UnregisterEvent(p.ObjectID(), signalID, handlerID)
				close(ch)
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
				p.UnregisterEvent(p.ObjectID(), signalID, handlerID)
				close(ch)
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

type timeval struct {
	Tv_sec  int64
	Tv_usec int64
}

func Readtimeval(r io.Reader) (s timeval, err error) {
	if s.Tv_sec, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read Tv_sec field: " + err.Error())
	}
	if s.Tv_usec, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read Tv_usec field: " + err.Error())
	}
	return s, nil
}
func Writetimeval(s timeval, w io.Writer) (err error) {
	if err := basic.WriteInt64(s.Tv_sec, w); err != nil {
		return fmt.Errorf("failed to write Tv_sec field: " + err.Error())
	}
	if err := basic.WriteInt64(s.Tv_usec, w); err != nil {
		return fmt.Errorf("failed to write Tv_usec field: " + err.Error())
	}
	return nil
}

type EventTrace struct {
	Id            uint32
	Kind          int32
	SlotId        uint32
	Arguments     value.Value
	Timestamp     timeval
	UserUsTime    int64
	SystemUsTime  int64
	CallerContext uint32
	CalleeContext uint32
}

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
	if s.Timestamp, err = Readtimeval(r); err != nil {
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
	if err := Writetimeval(s.Timestamp, w); err != nil {
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
