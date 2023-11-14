package bus_test

import (
	"testing"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/object"
)

func TestCancelledProxyCall(t *testing.T) {
	serviceEndpoint, clientEndpoint := net.Pipe()
	defer serviceEndpoint.Close()
	defer clientEndpoint.Close()

	msgChan, err := serviceEndpoint.ReceiveAny()
	if err != nil {
		t.Error(err)
	}

	go func() {
		m, ok := <-msgChan
		if !ok {
			panic("closed channel")
		}
		m.Header.Type = net.Cancelled
		m.Header.Size = 0
		m.Payload = []byte{}
		err := serviceEndpoint.Send(*m)
		if err != nil {
			panic(err)
		}
	}()

	c := bus.NewClient(bus.NewContext(clientEndpoint))
	proxy := bus.NewProxy(c, object.MetaService0, 1, 2)
	_, err = proxy.CallID(3, []byte{0xab, 0xcd})
	if err != bus.ErrCancelled {
		t.Errorf("wrong error type: %s", err)
	}
}
