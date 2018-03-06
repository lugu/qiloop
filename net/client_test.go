package net_test

import (
	"fmt"
	"github.com/lugu/qiloop/message"
	qinet "github.com/lugu/qiloop/net"
	"net"
	"testing"
)

func TestProxyCall(t *testing.T) {

	var p int
	var err error
	var conn net.Conn
	var ln net.Listener

	// 1. establish server
	for p = 1024; p < 66535; p++ {
		ln, err = net.Listen("tcp", fmt.Sprintf(":%d", p))
		if err == nil {
			break
		}
	}

	// 2. accept a single connection
	go func() {
		conn, err = ln.Accept()
		if err != nil {
			t.Errorf("failed to accept: %s", err)
			return
		}
		endpoint := qinet.AcceptedEndPoint(conn)
		m, err := endpoint.Receive()
		if err != nil {
			t.Errorf("failed to receive meesage: %s", err)
		}
		m.Header.Type = message.Reply
		err = endpoint.Send(m)
		if err != nil {
			t.Errorf("failed to send meesage: %s", err)
		}
		conn.Close()
		ln.Close()
	}()

	// 3. client estable connection
	client, err := qinet.NewClient(fmt.Sprintf(":%d", p))
	if err != nil {
		t.Errorf("failed to create client failed: %s", err)
	}

	// 4. create proxy
	proxy := qinet.NewProxy(client, 1, 2)
	// 4. client send a message
	_, err = proxy.Call(3, []byte{0xab, 0xcd})
	if err != nil {
		t.Errorf("proxy failed to call service: %s", err)
	}
}
