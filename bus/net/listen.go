package net

import (
	"crypto/tls"
	"fmt"
	"log"
	gonet "net"
	"net/url"
	"strings"

	"github.com/lugu/qiloop/bus/net/cert"
)

type connListener struct {
	l gonet.Listener
}

func (c connListener) Accept() (Stream, error) {
	conn, err := c.l.Accept()
	if err != nil {
		return nil, err
	}
	return ConnStream(conn), nil
}

func (c connListener) Close() error {
	return c.l.Close()
}

func listenTCP(addr string) (Listener, error) {
	conn, err := gonet.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return connListener{conn}, nil
}

func listenTLS(addr string) (Listener, error) {
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
		Certificates: []tls.Certificate{cer},
	}

	conn, err := tls.Listen("tcp", addr, conf)
	if err != nil {
		return nil, err
	}
	return connListener{conn}, nil
}

func listenUNIX(name string) (Listener, error) {
	conn, err := gonet.Listen("unix", name)
	if err != nil {
		return nil, err
	}
	return connListener{conn}, nil
}

type Listener interface {
	Accept() (Stream, error)
	Close() error
}

// Listen reads the transport of addr and listen at the address. addr
// can be of the form: unix://, tcp:// or tcps://.
func Listen(addr string) (Listener, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("listen: invalid address: %s", err)
	}
	switch u.Scheme {
	case "tcp":
		return listenTCP(u.Host)
	case "tcps":
		return listenTLS(u.Host)
	case "quic":
		return listenQUIC(u.Host)
	case "unix":
		return listenUNIX(strings.TrimPrefix(addr, "unix://"))
	default:
		return nil, fmt.Errorf("unknown URL scheme: %s", addr)
	}
}
