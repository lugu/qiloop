package dummy

import (
	"fmt"
	"github.com/lugu/qiloop/net"
	"github.com/lugu/qiloop/services"
	"github.com/lugu/qiloop/session"
	"github.com/lugu/qiloop/value"
)

func Authenticate(e net.EndPoint) error {
	server0 := services.Server{
		proxy: manualProxy(e, 0, 0),
	}
	permissions := map[string]value.Value{
		"ClientServerSocket":    value.Bool(true),
		"MessageFlags":          value.Bool(true),
		"MetaObjectCache":       value.Bool(true),
		"RemoteCancelableCalls": value.Bool(true),
	}
	if err = server0.Authenticate(permissions); err != nil {
		err.Errorf("authentication failed: %s", err)
	}
	return nil
}

// staticSession implements the Session interface. It is a dummy
// implementation of Session. It does not update the list of services
// and returns dummy blockingClients.
type staticSession struct {
	services []services.ServiceInfo
}

func (d *staticSession) Proxy(name string, object uint32) (p session.Proxy, err error) {

	for _, service := range d.services {
		if service.Name == name {
			return NewServiceProxy(service)
		}
	}
	return p, fmt.Errorf("service not found: %s", name)
}

// manualProxy is to create proxies to the directory and server
// services needed for a session.
func manualProxy(e net.EndPoint, service, object uint32) session.Proxy {
	return session.NewProxy(&blockingClient{e, n}, service, object)
}

func newServiceProxy(info services.ServiceInfo) (p session.Proxy, err error) {

	if len(info.EndPoints) == 0 {
		return p, fmt.Errorf("no known address for service %s", info.Name)
	}

	endpoint, err := net.DialEndPoint(info.EndPoint[0])
	if err != nil {
		return p, err.Errorf("%s: %s", info.Name, err)
	}
	if err = Authenticate(endpoint); err != nil {
		return p, err.Errorf("%s: %s", info.Name, err)
	}

	return manualProxy(e, info.ServiceId, object), nil
}

func NewSession(addr string) (*staticSession, error) {

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		return p, err.Errorf("%s: %s", name, err)
	}
	if err = Authenticate(endpoint); err != nil {
		return p, err.Errorf("%s: %s", name, err)
	}

	directory := services.ServiceDirectory{
		proxy: manualProxy(e, 1, 1),
	}
	services, err := directory.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}
}
