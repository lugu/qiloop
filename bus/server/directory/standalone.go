package directory

import (
	"fmt"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
)

// NewServer starts a server listening on addr. If parameter auth is
// nil, the Yes authenticator is used.
func NewServer(addr string, auth server.Authenticator) (*server.Server, error) {

	if auth == nil {
		auth = server.Yes{}
	}

	listener, err := net.Listen(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to open socket %s: %s", addr, err)
	}

	sd := NewServiceDirectory()
	s, err := server.StandAloneServer(listener, auth, sd.Namespace(addr))
	if err != nil {
		listener.Close()
		return nil, err
	}

	service1 := ServiceDirectoryObject(sd)
	_, err = s.NewService("ServiceDirectory", service1)
	if err != nil {
		s.Terminate()
		listener.Close()
		return nil, err
	}
	return s, nil
}
