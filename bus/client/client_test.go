package client_test

import (
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/object"
	"io"
	"testing"
)

func TestProxyCall(t *testing.T) {

	serviceEndpoint, clientEndpoint := net.NewPipe()
	defer serviceEndpoint.Close()
	defer clientEndpoint.Close()

	// accept a single connection
	go func() {
		for i := 0; i < 2; i++ {
			m, err := serviceEndpoint.ReceiveAny()
			if err == io.EOF {
				break
			} else if err != nil {
				t.Errorf("failed to receive meesage: %s", err)
			}
			m.Header.Type = net.Reply
			err = serviceEndpoint.Send(*m)
			if err != nil {
				t.Errorf("failed to send meesage: %s", err)
			}
		}
	}()

	// client connection
	c := client.NewClient(clientEndpoint)

	// 4. send a message
	proxy := client.NewProxy(c, object.MetaService0, 1, 2)
	_, err := proxy.CallID(3, []byte{0xab, 0xcd})
	if err != nil {
		t.Errorf("proxy failed to call service: %s", err)
	}
}
