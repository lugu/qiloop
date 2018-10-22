package net

import (
	"crypto/tls"
	"fmt"
	"github.com/lugu/qiloop/bus/net/cert"
	gonet "net"
	"net/url"
	"strings"
)

func listenTCP(addr string) (gonet.Listener, error) {
	return gonet.Listen("tcp", addr)
}

func listenTLS(addr string) (gonet.Listener, error) {
	certFile, keyFile := cert.GetCertKey()
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
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

func Listen(addr string) (gonet.Listener, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
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
