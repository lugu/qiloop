package net_test

import (
	"fmt"
	"github.com/lugu/qiloop/net"
	gonet "net"
	"reflect"
	"testing"
)

func TestPingPong(t *testing.T) {

	var p int
	var err error
	var ln gonet.Listener
	var addr string

	// 1. establish server
	for p = 1024; p < 66535; p++ {
		addr = fmt.Sprintf(":%d", p)
		ln, err = gonet.Listen("tcp", addr)
		if err == nil {
			defer ln.Close()
			break
		}
	}
	if ln == nil {
		t.Errorf("cannot find available port")
	}

	// 2. accept a single connection
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Errorf("failed to accept: %s", err)
			return
		}
		defer conn.Close()

		endpoint := net.NewEndPoint(conn)
		m, err := endpoint.ReceiveAny()
		if err != nil {
			t.Errorf("failed to receive meesage: %s", err)
		}
		err = endpoint.Send(*m)
		if err != nil {
			t.Errorf("failed to send meesage: %s", err)
		}
	}()

	// 3. client estable connection
	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		t.Errorf("dial failed: %s", err)
	}

	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	mSent := net.NewMessage(h, []byte{0xab, 0xcd})

	// 4. prepare to receive a message
	received := make(chan *net.Message)
	go func() {
		msg, err := endpoint.ReceiveAny()
		if err != nil {
			t.Errorf("failed to receive net. %s", err)
		}
		received <- msg
	}()

	// 5. client send a message
	if err = endpoint.Send(mSent); err != nil {
		t.Errorf("failed to send paquet: %s", err)
	}

	// 6. server replied
	mReceived := <-received

	// 7. check packet integrity
	if !reflect.DeepEqual(mSent, *mReceived) {
		t.Errorf("expected %#v, go %#v", mSent, mReceived)
	}
}
