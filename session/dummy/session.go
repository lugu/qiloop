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

func newServiceProxy(info services.ServiceInfo, objectID uint32) (p session.Proxy, err error) {

	if len(info.Endpoints) == 0 {
		return p, fmt.Errorf("no known address for service %s", info.Name)
	}

	addr := strings.TrimPrefix(info.Endpoints[0], "tcp://")
	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		return p, fmt.Errorf("%s: %s", info.Name, err)
	}
	if err = Authenticate(endpoint); err != nil {
		return p, fmt.Errorf("%s: %s", info.Name, err)
	}

	// FIXME: object id do be defined
	return metaProxy(endpoint, info.ServiceId, objectID), nil
}

func (d *staticSession) Proxy(name string, objectID uint32) (p session.Proxy, err error) {

	for _, service := range d.services {
		if service.Name == name {
			return newServiceProxy(service, objectID)
		}
	}
	return p, fmt.Errorf("service not found: %s", name)
}

// metaProxy is to create proxies to the directory and server
// services needed for a session.
func metaProxy(e net.EndPoint, serviceID, objectID uint32) session.Proxy {
	client := &blockingClient{e, 3}
	// FIXME: request a meta object
	meta, err := session.MetaObject(client, serviceID, objectID)
	if err != nil {
	}
	return NewProxy(client, meta, serviceID, objectID)
}

func NewSession(addr string) (s *staticSession, err error) {

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		return s, err
	}
	if err = Authenticate(endpoint); err != nil {
		return s, fmt.Errorf("authenitcation failed: %s", err)
	}

	directory := services.ServiceDirectory{
		metaProxy(endpoint, 1, 1),
	}
	s = new(staticSession)
	s.services, err = directory.Services()
	if err != nil {
		return s, fmt.Errorf("failed to list services: %s", err)
	}
	return s, nil
}
