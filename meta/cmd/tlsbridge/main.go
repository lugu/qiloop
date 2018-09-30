package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"
)

func listenTLS(addr string) {
	// openssl genrsa -out server.key 2048
	// openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("%s", err)
		return
	}

	conf := &tls.Config{
		Certificates:       []tls.Certificate{cer},
		InsecureSkipVerify: true,
	}

	ln, err := tls.Listen("tcp", addr, conf)
	if err != nil {
		log.Fatalf("%s", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("%s", err)
			continue
		}
		go forwardConnection(conn, connectLocal())
	}
}

func listenRAW(addr, remoteAddr string) {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("%s", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("%s", err)
			continue
		}
		forwardConnection(conn, connectRemote(remoteAddr))
	}
}

func copyAndClose(reader io.Reader, writer io.WriteCloser) {
	if _, err := io.Copy(writer, bufio.NewReader(reader)); err == nil {
		writer.Close()
	}
}

func forwardConnection(client, server net.Conn) {
	go copyAndClose(client, server)
	go copyAndClose(server, client)
}

func connectLocal() net.Conn {
	conn, err := net.Dial("tcp", ":12345")
	if err != nil {
		log.Fatalf("%s", err)
	}
	return conn
}

func connectRemote(addr string) net.Conn {

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return conn
}

func main() {
	var remoteHost = flag.String("remote-host", "", "remote address (host:port)")
	var localPort = flag.String("listen-port", ":9503", "local TLS port")
	var forwardPort = flag.String("tcp-port", ":12345", "local TCP port loopback")

	flag.Parse()

	if *remoteHost == "" {
		flag.PrintDefaults()
		return
	}

	go listenRAW(*forwardPort, *remoteHost)
	listenTLS(*localPort)
}
