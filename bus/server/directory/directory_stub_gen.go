// Package directory contains a generated stub
// File generated. DO NOT EDIT.
package directory

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	object1 "github.com/lugu/qiloop/bus/client/object"
	net "github.com/lugu/qiloop/bus/net"
	server "github.com/lugu/qiloop/bus/server"
	generic "github.com/lugu/qiloop/bus/server/generic"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	"io"
	"log"
)

// ServiceDirectoryImplementor interface of the service implementation
type ServiceDirectoryImplementor interface {
	// Activate is called before any other method.
	// It shall be used to initialize the interface.
	// activation provides runtime informations.
	// activation.Terminate() unregisters the object.
	// activation.Session can access other services.
	// helper enables signals an properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation server.Activation, helper ServiceDirectorySignalHelper) error
	OnTerminate()
	Service(name string) (ServiceInfo, error)
	Services() ([]ServiceInfo, error)
	RegisterService(info ServiceInfo) (uint32, error)
	UnregisterService(serviceID uint32) error
	ServiceReady(serviceID uint32) error
	UpdateServiceInfo(info ServiceInfo) error
	MachineId() (string, error)
	_socketOfService(serviceID uint32) (object.ObjectReference, error)
}

// ServiceDirectorySignalHelper provided to ServiceDirectory a companion object
type ServiceDirectorySignalHelper interface {
	SignalServiceAdded(serviceID uint32, name string) error
	SignalServiceRemoved(serviceID uint32, name string) error
}

// stubServiceDirectory implements server.ServerObject.
type stubServiceDirectory struct {
	obj     generic.Object
	impl    ServiceDirectoryImplementor
	session bus.Session
}

// ServiceDirectoryObject returns an object using ServiceDirectoryImplementor
func ServiceDirectoryObject(impl ServiceDirectoryImplementor) server.ServerObject {
	var stb stubServiceDirectory
	stb.impl = impl
	stb.obj = generic.NewObject(stb.metaObject())
	stb.obj.Wrap(uint32(0x64), stb.Service)
	stb.obj.Wrap(uint32(0x65), stb.Services)
	stb.obj.Wrap(uint32(0x66), stb.RegisterService)
	stb.obj.Wrap(uint32(0x67), stb.UnregisterService)
	stb.obj.Wrap(uint32(0x68), stb.ServiceReady)
	stb.obj.Wrap(uint32(0x69), stb.UpdateServiceInfo)
	stb.obj.Wrap(uint32(0x6c), stb.MachineId)
	stb.obj.Wrap(uint32(0x6d), stb._socketOfService)
	return &stb
}
func (p *stubServiceDirectory) Activate(activation server.Activation) error {
	p.session = activation.Session
	p.obj.Activate(activation)
	return p.impl.Activate(activation, p)
}
func (p *stubServiceDirectory) OnTerminate() {
	p.impl.OnTerminate()
	p.obj.OnTerminate()
}
func (p *stubServiceDirectory) Receive(msg *net.Message, from *server.Context) error {
	return p.obj.Receive(msg, from)
}
func (p *stubServiceDirectory) Service(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	name, err := basic.ReadString(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read name: %s", err)
	}
	ret, callErr := p.impl.Service(name)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := WriteServiceInfo(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubServiceDirectory) Services(payload []byte) ([]byte, error) {
	ret, callErr := p.impl.Services()
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
			err = WriteServiceInfo(v, &out)
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
func (p *stubServiceDirectory) RegisterService(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	info, err := ReadServiceInfo(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read info: %s", err)
	}
	ret, callErr := p.impl.RegisterService(info)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteUint32(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubServiceDirectory) UnregisterService(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	serviceID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read serviceID: %s", err)
	}
	callErr := p.impl.UnregisterService(serviceID)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubServiceDirectory) ServiceReady(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	serviceID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read serviceID: %s", err)
	}
	callErr := p.impl.ServiceReady(serviceID)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubServiceDirectory) UpdateServiceInfo(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	info, err := ReadServiceInfo(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read info: %s", err)
	}
	callErr := p.impl.UpdateServiceInfo(info)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (p *stubServiceDirectory) MachineId(payload []byte) ([]byte, error) {
	ret, callErr := p.impl.MachineId()
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteString(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubServiceDirectory) _socketOfService(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	serviceID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read serviceID: %s", err)
	}
	ret, callErr := p.impl._socketOfService(serviceID)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := object.WriteObjectReference(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (p *stubServiceDirectory) SignalServiceAdded(serviceID uint32, name string) error {
	var buf bytes.Buffer
	if err := basic.WriteUint32(serviceID, &buf); err != nil {
		return fmt.Errorf("failed to serialize serviceID: %s", err)
	}
	if err := basic.WriteString(name, &buf); err != nil {
		return fmt.Errorf("failed to serialize name: %s", err)
	}
	err := p.obj.UpdateSignal(uint32(0x6a), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalServiceAdded: %s", err)
	}
	return nil
}
func (p *stubServiceDirectory) SignalServiceRemoved(serviceID uint32, name string) error {
	var buf bytes.Buffer
	if err := basic.WriteUint32(serviceID, &buf); err != nil {
		return fmt.Errorf("failed to serialize serviceID: %s", err)
	}
	if err := basic.WriteString(name, &buf); err != nil {
		return fmt.Errorf("failed to serialize name: %s", err)
	}
	err := p.obj.UpdateSignal(uint32(0x6b), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalServiceRemoved: %s", err)
	}
	return nil
}
func (p *stubServiceDirectory) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "ServiceDirectory",
		Methods: map[uint32]object.MetaMethod{
			uint32(0x64): {
				Name:                "service",
				ParametersSignature: "(s)",
				ReturnSignature:     "(sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>",
				Uid:                 uint32(0x64),
			},
			uint32(0x65): {
				Name:                "services",
				ParametersSignature: "()",
				ReturnSignature:     "[(sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>]",
				Uid:                 uint32(0x65),
			},
			uint32(0x66): {
				Name:                "registerService",
				ParametersSignature: "((sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>)",
				ReturnSignature:     "I",
				Uid:                 uint32(0x66),
			},
			uint32(0x67): {
				Name:                "unregisterService",
				ParametersSignature: "(I)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x67),
			},
			uint32(0x68): {
				Name:                "serviceReady",
				ParametersSignature: "(I)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x68),
			},
			uint32(0x69): {
				Name:                "updateServiceInfo",
				ParametersSignature: "((sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x69),
			},
			uint32(0x6c): {
				Name:                "machineId",
				ParametersSignature: "()",
				ReturnSignature:     "s",
				Uid:                 uint32(0x6c),
			},
			uint32(0x6d): {
				Name:                "_socketOfService",
				ParametersSignature: "(I)",
				ReturnSignature:     "o",
				Uid:                 uint32(0x6d),
			},
		},
		Signals: map[uint32]object.MetaSignal{
			uint32(0x6a): {
				Name:      "serviceAdded",
				Signature: "(Is)",
				Uid:       uint32(0x6a),
			},
			uint32(0x6b): {
				Name:      "serviceRemoved",
				Signature: "(Is)",
				Uid:       uint32(0x6b),
			},
		},
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

// ServiceAdded is serializable
type ServiceAdded struct {
	ServiceID uint32
	Name      string
}

// ReadServiceAdded unmarshalls ServiceAdded
func ReadServiceAdded(r io.Reader) (s ServiceAdded, err error) {
	if s.ServiceID, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read ServiceID field: " + err.Error())
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	return s, nil
}

// WriteServiceAdded marshalls ServiceAdded
func WriteServiceAdded(s ServiceAdded, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.ServiceID, w); err != nil {
		return fmt.Errorf("failed to write ServiceID field: " + err.Error())
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
	}
	return nil
}

// ServiceRemoved is serializable
type ServiceRemoved struct {
	ServiceID uint32
	Name      string
}

// ReadServiceRemoved unmarshalls ServiceRemoved
func ReadServiceRemoved(r io.Reader) (s ServiceRemoved, err error) {
	if s.ServiceID, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read ServiceID field: " + err.Error())
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	return s, nil
}

// WriteServiceRemoved marshalls ServiceRemoved
func WriteServiceRemoved(s ServiceRemoved, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.ServiceID, w); err != nil {
		return fmt.Errorf("failed to write ServiceID field: " + err.Error())
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
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
	// SocketOfService calls the remote procedure
	SocketOfService(serviceID uint32) (object.ObjectReference, error)
	// SubscribeServiceAdded subscribe to a remote signal
	SubscribeServiceAdded() (unsubscribe func(), updates chan ServiceAdded, err error)
	// SubscribeServiceRemoved subscribe to a remote signal
	SubscribeServiceRemoved() (unsubscribe func(), updates chan ServiceRemoved, err error)
}

// ServiceDirectory represents a proxy object to the service
type ServiceDirectoryProxy interface {
	object.Object
	bus.Proxy
	ServiceDirectory
}

// proxyServiceDirectory implements ServiceDirectoryProxy
type proxyServiceDirectory struct {
	object1.ObjectProxy
	session bus.Session
}

// MakeServiceDirectory constructs ServiceDirectoryProxy
func MakeServiceDirectory(sess bus.Session, proxy bus.Proxy) ServiceDirectoryProxy {
	return &proxyServiceDirectory{object1.MakeObject(proxy), sess}
}

// NewServiceDirectory constructs ServiceDirectoryProxy
func NewServiceDirectory(sess bus.Session, obj uint32) (ServiceDirectoryProxy, error) {
	proxy, err := sess.Proxy("ServiceDirectory", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return MakeServiceDirectory(sess, proxy), nil
}

// ServiceDirectory retruns a proxy to a remote service
func (s Constructor) ServiceDirectory() (ServiceDirectoryProxy, error) {
	return NewServiceDirectory(s.session, 1)
}

// Service calls the remote procedure
func (p *proxyServiceDirectory) Service(name string) (ServiceInfo, error) {
	var err error
	var ret ServiceInfo
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(name, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize name: %s", err)
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
func (p *proxyServiceDirectory) Services() ([]ServiceInfo, error) {
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
func (p *proxyServiceDirectory) RegisterService(info ServiceInfo) (uint32, error) {
	var err error
	var ret uint32
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = WriteServiceInfo(info, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize info: %s", err)
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
func (p *proxyServiceDirectory) UnregisterService(serviceID uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(serviceID, buf); err != nil {
		return fmt.Errorf("failed to serialize serviceID: %s", err)
	}
	_, err = p.Call("unregisterService", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterService failed: %s", err)
	}
	return nil
}

// ServiceReady calls the remote procedure
func (p *proxyServiceDirectory) ServiceReady(serviceID uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(serviceID, buf); err != nil {
		return fmt.Errorf("failed to serialize serviceID: %s", err)
	}
	_, err = p.Call("serviceReady", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call serviceReady failed: %s", err)
	}
	return nil
}

// UpdateServiceInfo calls the remote procedure
func (p *proxyServiceDirectory) UpdateServiceInfo(info ServiceInfo) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = WriteServiceInfo(info, buf); err != nil {
		return fmt.Errorf("failed to serialize info: %s", err)
	}
	_, err = p.Call("updateServiceInfo", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call updateServiceInfo failed: %s", err)
	}
	return nil
}

// MachineId calls the remote procedure
func (p *proxyServiceDirectory) MachineId() (string, error) {
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
func (p *proxyServiceDirectory) SocketOfService(serviceID uint32) (object.ObjectReference, error) {
	var err error
	var ret object.ObjectReference
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(serviceID, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize serviceID: %s", err)
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

// SubscribeServiceAdded subscribe to a remote property
func (p *proxyServiceDirectory) SubscribeServiceAdded() (func(), chan ServiceAdded, error) {
	propertyID, err := p.SignalID("serviceAdded")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "serviceAdded", err)
	}

	handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register event for %s: %s", "serviceAdded", err)
	}
	ch := make(chan ServiceAdded)
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
			e, err := ReadServiceAdded(buf)
			if err != nil {
				log.Printf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return cancel, ch, nil
}

// SubscribeServiceRemoved subscribe to a remote property
func (p *proxyServiceDirectory) SubscribeServiceRemoved() (func(), chan ServiceRemoved, error) {
	propertyID, err := p.SignalID("serviceRemoved")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "serviceRemoved", err)
	}

	handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register event for %s: %s", "serviceRemoved", err)
	}
	ch := make(chan ServiceRemoved)
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
			e, err := ReadServiceRemoved(buf)
			if err != nil {
				log.Printf("failed to unmarshall tuple: %s", err)
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

// ReadServiceInfo unmarshalls ServiceInfo
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

// WriteServiceInfo marshalls ServiceInfo
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
