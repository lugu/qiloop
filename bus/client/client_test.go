package client_test

import (
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/object"
	"testing"
)

func TestProxyCall(t *testing.T) {

	serviceEndpoint, clientEndpoint := net.NewPipe()
	defer serviceEndpoint.Close()
	defer clientEndpoint.Close()

	msgChan, err := serviceEndpoint.ReceiveOne()
	if err != nil {
		panic(err)
	}

	// accept a single connection
	go func() {
		m, ok := <-msgChan
		if !ok {
			panic("connection closed")
		}
		m.Header.Type = net.Reply
		err := serviceEndpoint.Send(*m)
		if err != nil {
			t.Errorf("failed to send meesage: %s", err)
		}
	}()

	// client connection
	c := client.NewClient(clientEndpoint)

	// 4. send a message
	proxy := client.NewProxy(c, object.MetaService0, 1, 2)
	_, err = proxy.CallID(3, []byte{0xab, 0xcd})
	if err != nil {
		t.Errorf("proxy failed to call service: %s", err)
	}
}
