package session

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
	"strings"
)

type CapabilityMap map[string]value.Value

func Authenticate(endpoint net.EndPoint) error {

	const serviceID = 0
	const objectID = 0

	client0 := newClient(endpoint)
	proxy0 := NewProxy(client0, object.MetaService0, serviceID, objectID)
	server0 := services.ServerProxy{proxy0}

	permissions := CapabilityMap{
		"ClientServerSocket":    value.Bool(true),
		"MessageFlags":          value.Bool(true),
		"MetaObjectCache":       value.Bool(true),
		"RemoteCancelableCalls": value.Bool(true),
	}
	if _, err := server0.Authenticate(permissions); err != nil {
		fmt.Errorf("authentication failed: %s", err)
	}
	return nil
}

// staticSession implements the Session interface. It is an
// implementation of Session. It does not update the list of services
// and returns clients.

type staticSession struct {
	services []services.ServiceInfo
}

func newEndPoint(info services.ServiceInfo) (endpoint net.EndPoint, err error) {
	if len(info.Endpoints) == 0 {
		return endpoint, fmt.Errorf("missing address for service %s", info.Name)
	}
	addr := strings.TrimPrefix(info.Endpoints[0], "tcp://")
	endpoint, err = net.DialEndPoint(addr)
	if err != nil {
		return endpoint, fmt.Errorf("%s: connection failed (%s) : %s", info.Name, addr, err)
	}
	if err = Authenticate(endpoint); err != nil {
		return endpoint, fmt.Errorf("authentication error (%s): %s", info.Name, err)
	}
	return endpoint, nil
}

func newObject(info services.ServiceInfo, ref object.ObjectReference) (object.Object, error) {
	endpoint, err := newEndPoint(info)
	if err != nil {
		return nil, fmt.Errorf("object connection error (%s): %s", info.Name, err)
	}
	proxy := newProxy(endpoint, ref.MetaObject, ref.ServiceID, ref.ObjectID)
	return &services.ObjectProxy{proxy}, nil
}

func newService(info services.ServiceInfo, objectID uint32) (p bus.Proxy, err error) {
	endpoint, err := newEndPoint(info)
	if err != nil {
		return nil, fmt.Errorf("service connection error (%s): %s", info.Name, err)
	}
	proxy, err := metaProxy(endpoint, info.ServiceId, objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service meta object (%s): %s", info.Name, err)
	}
	return proxy, nil
}

// FIXME: objectID does not seems needed: it can be deduce from the
// name
func (d *staticSession) Proxy(name string, objectID uint32) (p bus.Proxy, err error) {

	for _, service := range d.services {
		if service.Name == name {
			return newService(service, objectID)
		}
	}
	return p, fmt.Errorf("service not found: %s", name)
}

func (d *staticSession) Object(ref object.ObjectReference) (o object.Object, err error) {
	for _, service := range d.services {
		if service.ServiceId == ref.ServiceID {
			return newObject(service, ref)
		}
	}
	return o, fmt.Errorf("Not yet implemented")
}

func (d *staticSession) Register(name string, service bus.Service) error {
	return fmt.Errorf("not yet implemented")
}

func newProxy(e net.EndPoint, meta object.MetaObject, serviceID, objectID uint32) bus.Proxy {
	return NewProxy(newClient(e), meta, serviceID, objectID)
}

// metaProxy is to create proxies to the directory and server
// services needed for a session.
func metaProxy(e net.EndPoint, serviceID, objectID uint32) (p bus.Proxy, err error) {
	client := newClient(e)
	meta, err := bus.MetaObject(client, serviceID, objectID)
	if err != nil {
		return p, fmt.Errorf("Can not reach metaObject: %s", err)
	}
	return NewProxy(client, meta, serviceID, objectID), nil
}

func NewSession(addr string) (s *staticSession, err error) {

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		return s, fmt.Errorf("failed to contact %s: %s", addr, err)
	}
	if err = Authenticate(endpoint); err != nil {
		return s, fmt.Errorf("authenitcation failed: %s", err)
	}

	proxy, err := metaProxy(endpoint, 1, 1)
	if err != nil {
		return s, fmt.Errorf("failed to get directory meta object: %s", err)
	}
	directory := services.ServiceDirectoryProxy{proxy}
	s = new(staticSession)
	s.services, err = directory.Services()
	if err != nil {
		return s, fmt.Errorf("failed to list services: %s", err)
	}
	return s, nil
}
