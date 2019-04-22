package directory

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
)

// NewServer starts a server listening on addr. If parameter auth is
// nil, the Yes authenticator is used.
func NewServer(addr string, auth bus.Authenticator) (bus.Server, error) {

	if auth == nil {
		auth = bus.Yes{}
	}

	listener, err := net.Listen(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to open socket %s: %s", addr, err)
	}

	sd := ServiceDirectoryImpl()
	namespace := sd.Namespace(addr)
	service1 := ServiceDirectoryObject(sd)

	s, err := bus.NewServer(listener, auth, namespace, service1)
	if err != nil {
		listener.Close()
		return nil, err
	}
	return s, nil
}
