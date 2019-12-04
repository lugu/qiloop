package directory

import (
	"bytes"
	"context"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	net "github.com/lugu/qiloop/bus/net"
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
	// helper enables signals and properties updates.
	// Properties must be initialized using helper,
	// during the Activate call.
	Activate(activation bus.Activation, helper ServiceDirectorySignalHelper) error
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

// stubServiceDirectory implements server.Actor.
type stubServiceDirectory struct {
	impl      ServiceDirectoryImplementor
	session   bus.Session
	service   bus.Service
	serviceID uint32
	signal    bus.SignalHandler
}

// ServiceDirectoryObject returns an object using ServiceDirectoryImplementor
func ServiceDirectoryObject(impl ServiceDirectoryImplementor) bus.Actor {
	var stb stubServiceDirectory
	stb.impl = impl
	obj := bus.NewBasicObject(&stb, stb.metaObject(), stb.onPropertyChange)
	stb.signal = obj
	return obj
}

// CreateServiceDirectory registers a new object to a service
// and returns a proxy to the newly created object
func CreateServiceDirectory(session bus.Session, service bus.Service, impl ServiceDirectoryImplementor) (ServiceDirectoryProxy, error) {
	obj := ServiceDirectoryObject(impl)
	objectID, err := service.Add(obj)
	if err != nil {
		return nil, err
	}
	stb := &stubServiceDirectory{}
	meta := object.FullMetaObject(stb.metaObject())
	client := bus.DirectClient(obj)
	proxy := bus.NewProxy(client, meta, service.ServiceID(), objectID)
	return MakeServiceDirectory(session, proxy), nil
}
func (p *stubServiceDirectory) Activate(activation bus.Activation) error {
	p.session = activation.Session
	p.service = activation.Service
	p.serviceID = activation.ServiceID
	return p.impl.Activate(activation, p)
}
func (p *stubServiceDirectory) OnTerminate() {
	p.impl.OnTerminate()
}
func (p *stubServiceDirectory) Receive(msg *net.Message, from bus.Channel) error {
	// action dispatch
	switch msg.Header.Action {
	case 100:
		return p.Service(msg, from)
	case 101:
		return p.Services(msg, from)
	case 102:
		return p.RegisterService(msg, from)
	case 103:
		return p.UnregisterService(msg, from)
	case 104:
		return p.ServiceReady(msg, from)
	case 105:
		return p.UpdateServiceInfo(msg, from)
	case 108:
		return p.MachineId(msg, from)
	case 109:
		return p._socketOfService(msg, from)
	default:
		return from.SendError(msg, bus.ErrActionNotFound)
	}
}
func (p *stubServiceDirectory) onPropertyChange(name string, data []byte) error {
	switch name {
	default:
		return fmt.Errorf("unknown property %s", name)
	}
}
func (p *stubServiceDirectory) Service(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	name, err := basic.ReadString(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read name: %s", err))
	}
	ret, callErr := p.impl.Service(name)

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := writeServiceInfo(ret, &out)
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubServiceDirectory) Services(msg *net.Message, c bus.Channel) error {
	ret, callErr := p.impl.Services()

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
			err = writeServiceInfo(v, &out)
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
func (p *stubServiceDirectory) RegisterService(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	info, err := readServiceInfo(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read info: %s", err))
	}
	ret, callErr := p.impl.RegisterService(info)

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := basic.WriteUint32(ret, &out)
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubServiceDirectory) UnregisterService(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	serviceID, err := basic.ReadUint32(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read serviceID: %s", err))
	}
	callErr := p.impl.UnregisterService(serviceID)

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
func (p *stubServiceDirectory) ServiceReady(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	serviceID, err := basic.ReadUint32(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read serviceID: %s", err))
	}
	callErr := p.impl.ServiceReady(serviceID)

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
func (p *stubServiceDirectory) UpdateServiceInfo(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	info, err := readServiceInfo(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read info: %s", err))
	}
	callErr := p.impl.UpdateServiceInfo(info)

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
func (p *stubServiceDirectory) MachineId(msg *net.Message, c bus.Channel) error {
	ret, callErr := p.impl.MachineId()

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := basic.WriteString(ret, &out)
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubServiceDirectory) _socketOfService(msg *net.Message, c bus.Channel) error {
	buf := bytes.NewBuffer(msg.Payload)
	serviceID, err := basic.ReadUint32(buf)
	if err != nil {
		return c.SendError(msg, fmt.Errorf("cannot read serviceID: %s", err))
	}
	ret, callErr := p.impl._socketOfService(serviceID)

	// do not respond to post messages.
	if msg.Header.Type == net.Post {
		return nil
	}
	if callErr != nil {
		return c.SendError(msg, callErr)
	}
	var out bytes.Buffer
	errOut := object.WriteObjectReference(ret, &out)
	if errOut != nil {
		return c.SendError(msg, fmt.Errorf("cannot write response: %s", errOut))
	}
	return c.SendReply(msg, out.Bytes())
}
func (p *stubServiceDirectory) SignalServiceAdded(serviceID uint32, name string) error {
	var buf bytes.Buffer
	if err := basic.WriteUint32(serviceID, &buf); err != nil {
		return fmt.Errorf("serialize serviceID: %s", err)
	}
	if err := basic.WriteString(name, &buf); err != nil {
		return fmt.Errorf("serialize name: %s", err)
	}
	err := p.signal.UpdateSignal(106, buf.Bytes())

	if err != nil {
		return fmt.Errorf("update SignalServiceAdded: %s", err)
	}
	return nil
}
func (p *stubServiceDirectory) SignalServiceRemoved(serviceID uint32, name string) error {
	var buf bytes.Buffer
	if err := basic.WriteUint32(serviceID, &buf); err != nil {
		return fmt.Errorf("serialize serviceID: %s", err)
	}
	if err := basic.WriteString(name, &buf); err != nil {
		return fmt.Errorf("serialize name: %s", err)
	}
	err := p.signal.UpdateSignal(107, buf.Bytes())

	if err != nil {
		return fmt.Errorf("update SignalServiceRemoved: %s", err)
	}
	return nil
}
func (p *stubServiceDirectory) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "ServiceDirectory",
		Methods: map[uint32]object.MetaMethod{
			100: {
				Name:                "service",
				ParametersSignature: "(s)",
				ReturnSignature:     "(sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>",
				Uid:                 100,
			},
			101: {
				Name:                "services",
				ParametersSignature: "()",
				ReturnSignature:     "[(sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>]",
				Uid:                 101,
			},
			102: {
				Name:                "registerService",
				ParametersSignature: "((sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>)",
				ReturnSignature:     "I",
				Uid:                 102,
			},
			103: {
				Name:                "unregisterService",
				ParametersSignature: "(I)",
				ReturnSignature:     "v",
				Uid:                 103,
			},
			104: {
				Name:                "serviceReady",
				ParametersSignature: "(I)",
				ReturnSignature:     "v",
				Uid:                 104,
			},
			105: {
				Name:                "updateServiceInfo",
				ParametersSignature: "((sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>)",
				ReturnSignature:     "v",
				Uid:                 105,
			},
			108: {
				Name:                "machineId",
				ParametersSignature: "()",
				ReturnSignature:     "s",
				Uid:                 108,
			},
			109: {
				Name:                "_socketOfService",
				ParametersSignature: "(I)",
				ReturnSignature:     "o",
				Uid:                 109,
			},
		},
		Properties: map[uint32]object.MetaProperty{},
		Signals: map[uint32]object.MetaSignal{
			106: {
				Name:      "serviceAdded",
				Signature: "(Is)<serviceAdded,serviceID,name>",
				Uid:       106,
			},
			107: {
				Name:      "serviceRemoved",
				Signature: "(Is)<serviceRemoved,serviceID,name>",
				Uid:       107,
			},
		},
	}
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
	return MakeServiceDirectory(p.session, bus.WithContext(p.Proxy(), ctx))
}

// Service calls the remote procedure
func (p *proxyServiceDirectory) Service(name string) (ServiceInfo, error) {
	var err error
	var ret ServiceInfo
	var buf bytes.Buffer
	if err = basic.WriteString(name, &buf); err != nil {
		return ret, fmt.Errorf("serialize name: %s", err)
	}
	response, err := p.Proxy().Call("service", buf.Bytes())
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
	response, err := p.Proxy().Call("services", buf.Bytes())
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
	response, err := p.Proxy().Call("registerService", buf.Bytes())
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
	_, err = p.Proxy().Call("unregisterService", buf.Bytes())
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
	_, err = p.Proxy().Call("serviceReady", buf.Bytes())
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
	_, err = p.Proxy().Call("updateServiceInfo", buf.Bytes())
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
	response, err := p.Proxy().Call("machineId", buf.Bytes())
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
	response, err := p.Proxy().Call("_socketOfService", buf.Bytes())
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
	propertyID, err := p.Proxy().SignalID("serviceAdded")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "serviceAdded", err)
	}
	ch := make(chan ServiceAdded)
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
	propertyID, err := p.Proxy().SignalID("serviceRemoved")
	if err != nil {
		return nil, nil, fmt.Errorf("property %s not available: %s", "serviceRemoved", err)
	}
	ch := make(chan ServiceRemoved)
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
