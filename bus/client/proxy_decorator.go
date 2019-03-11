package client

import "fmt"

// ServiceServer retruns a proxy to the authenticating service (ID 0)
func (s Constructor) ServiceServer() (ServerProxy, error) {
	proxy, err := s.session.Proxy("Server", 0)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &proxyServer{proxy}, nil
}
