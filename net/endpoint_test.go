package net_test

import (
	"fmt"
	"github.com/lugu/qiloop/message"
	qinet "github.com/lugu/qiloop/net"
	"net"
	"reflect"
	"testing"
)

func TestPingPong(t *testing.T) {

	var p int
	var err error
	var ln net.Listener

	// 1. establish server
	for p = 1024; p < 66535; p++ {
		ln, err = net.Listen("tcp", fmt.Sprintf(":%d", p))
		if err == nil {
			break
		}
	}
	defer ln.Close()

	// 2. accept a single connection
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Errorf("failed to accept: %s", err)
			return
		}
		defer conn.Close()
		endpoint := qinet.AcceptedEndPoint(conn)
		m, err := endpoint.Receive()
		if err != nil {
			t.Errorf("failed to receive meesage: %s", err)
		}
		err = endpoint.Send(m)
		if err != nil {
			t.Errorf("failed to send meesage: %s", err)
		}
	}()

	// 3. client estable connection
	endpoint, err := qinet.DialEndPoint(fmt.Sprintf(":%d", p))
	if err != nil {
		t.Errorf("dial failed: %s", err)
	}

	// 4. client send a message
	h := message.NewHeader(message.Call, 1, 2, 3, 4)
	mSent := message.NewMessage(h, []byte{0xab, 0xcd})

	if err = endpoint.Send(mSent); err != nil {
		t.Errorf("failed to send paquet: %s", err)
	}

	// 5. server reply
	mReceived, err := endpoint.Receive()
	if err != nil {
		t.Errorf("failed to receive message: %s", err)
	}

	// 6. check packet integrity
	if !reflect.DeepEqual(mSent, mReceived) {
		t.Errorf("expected %#v, go %#v", mSent, mReceived)
	}
}
