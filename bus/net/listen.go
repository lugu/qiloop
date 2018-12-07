package net

import (
	"crypto/tls"
	"fmt"
	"github.com/lugu/qiloop/bus/net/cert"
	"log"
	gonet "net"
	"net/url"
	"strings"
)

func listenTCP(addr string) (gonet.Listener, error) {
	return gonet.Listen("tcp", addr)
}

func listenTLS(addr string) (gonet.Listener, error) {
	var err1, err2 error
	cer, err1 := cert.Certificate()
	if err1 != nil {
		log.Printf("Failed to read x509 certificate: %s", err1)
		cer, err2 = cert.GenerateCertificate()
		if err2 != nil {
			log.Printf("Failed to create x509 certificate: %s", err2)
			return nil, fmt.Errorf("no certificate available (%s, %s)",
				err1, err2)
		}
	}

	conf := &tls.Config{
		Certificates:       []tls.Certificate{cer},
		InsecureSkipVerify: true,
	}

	return tls.Listen("tcp", addr, conf)
}

func listenUNIX(name string) (gonet.Listener, error) {
	return gonet.Listen("unix", name)
}

// Listen reads the transport of addr and listen at the address. addr
// can be of the form: unix://, tcp:// or tcps://.
func Listen(addr string) (gonet.Listener, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("listen: invalid address: %s", err)
	}
	switch u.Scheme {
	case "tcp":
		return listenTCP(u.Host)
	case "tcps":
		return listenTLS(u.Host)
	case "unix":
		return listenUNIX(strings.TrimPrefix(addr, "unix://"))
	default:
		return nil, fmt.Errorf("unknown URL scheme: %s", addr)
	}
}
