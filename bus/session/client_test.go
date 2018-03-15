package session_test

import (
	"fmt"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/type/object"
	gonet "net"
	"testing"
)

func TestProxyCall(t *testing.T) {

	var p int
	var err error
	var ln gonet.Listener

	// 1. establish server
	for p = 1024; p < 66535; p++ {
		ln, err = gonet.Listen("tcp", fmt.Sprintf(":%d", p))
		if err == nil {
			break
		}
	}
	defer ln.Close()

	// 2. accept a single connection
	go func() {
		var err error
		conn, err := ln.Accept()
		if err != nil {
			t.Errorf("failed to accept: %s", err)
			return
		}
		defer conn.Close()
		endpoint := net.NewEndPoint(conn)

		for i := 0; i < 2; i++ {
			m, err := endpoint.ReceiveAny()
			if err != nil {
				t.Errorf("failed to receive meesage: %s", err)
			}
			m.Header.Type = net.Reply
			err = endpoint.Send(*m)
			if err != nil {
				t.Errorf("failed to send meesage: %s", err)
			}
		}
	}()

	// 3. client estable connection
	client, err := session.NewClient(fmt.Sprintf(":%d", p))
	if err != nil {
		t.Errorf("failed to create client failed: %s", err)
	}

	// 4. create proxy
	proxy := session.NewProxy(client, object.MetaService0, 1, 2)
	// 4. client send a message
	_, err = proxy.CallID(3, []byte{0xab, 0xcd})
	if err != nil {
		t.Errorf("proxy failed to call service: %s", err)
	}
}
