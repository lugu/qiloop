package dummy

import (
	"fmt"
	"github.com/lugu/qiloop/net"
	"github.com/lugu/qiloop/object"
	"github.com/lugu/qiloop/services"
	"github.com/lugu/qiloop/session"
	"github.com/lugu/qiloop/value"
	"strings"
)

func Authenticate(endpoint net.EndPoint) error {

	const serviceID = 0
	const objectID = 0
	const messageID = 0

	client0 := &blockingClient{endpoint, messageID}
	proxy0 := NewProxy(client0, object.MetaService0, serviceID, objectID)
	server0 := services.Server{proxy0}

	permissions := map[string]value.Value{
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

// staticSession implements the Session interface. It is a dummy
// implementation of Session. It does not update the list of services
// and returns dummy blockingClients.
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
	proxy, err := metaProxy(endpoint, ref.ServiceID, ref.ObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get object meta object (%s): %s", info.Name, err)
	}
	return &services.Object{proxy}, nil
}

func newService(info services.ServiceInfo, objectID uint32) (p session.Proxy, err error) {
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

func (d *staticSession) Proxy(name string, objectID uint32) (p session.Proxy, err error) {

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

// metaProxy is to create proxies to the directory and server
// services needed for a session.
func metaProxy(e net.EndPoint, serviceID, objectID uint32) (p session.Proxy, err error) {
	client := &blockingClient{e, 3}
	meta, err := session.MetaObject(client, serviceID, objectID)
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
	directory := services.ServiceDirectory{proxy}
	s = new(staticSession)
	s.services, err = directory.Services()
	if err != nil {
		return s, fmt.Errorf("failed to list services: %s", err)
	}
	return s, nil
}
