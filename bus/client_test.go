package bus_test

import (
	"sync"
	"testing"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/directory"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
)

func TestClientCall(t *testing.T) {

	serviceEndpoint, clientEndpoint := net.Pipe()
	defer serviceEndpoint.Close()
	defer clientEndpoint.Close()

	msgChan, err := serviceEndpoint.ReceiveAny()
	if err != nil {
		t.Error(err)
	}

	// accept a single connection
	go func() {
		m, ok := <-msgChan
		if !ok {
			t.Fatalf("connection closed")
		}
		m.Header.Type = net.Reply
		err := serviceEndpoint.Send(*m)
		if err != nil {
			t.Errorf("send meesage: %s", err)
		}
	}()

	// client connection
	c := bus.NewClient(bus.NewContext(clientEndpoint))
	_, err = c.Call(nil, 1, 2, 3, []byte{0xab, 0xcd})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestClientAlreadyCancelled(t *testing.T) {

	serviceEndpoint, clientEndpoint := net.Pipe()
	defer serviceEndpoint.Close()
	defer clientEndpoint.Close()

	cancel := make(chan struct{})
	close(cancel)
	c := bus.NewClient(bus.NewContext(clientEndpoint))
	_, err := c.Call(cancel, 1, 2, 3, []byte{0xab, 0xcd})
	if err != bus.ErrCancelled {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestClientCancel(t *testing.T) {

	serviceEndpoint, clientEndpoint := net.Pipe()
	defer serviceEndpoint.Close()
	defer clientEndpoint.Close()

	msgChan, err := serviceEndpoint.ReceiveAny()
	if err != nil {
		t.Error(err)
	}

	// testSync: make sure we cancel after the call is sent
	testSync := make(chan struct{})

	// accept a single connection
	go func() {
		m, ok := <-msgChan
		if !ok {
			t.Fatalf("connection closed")
		}
		if m.Header.Type != net.Call {
			t.Errorf("not a call")
		}
		msgChan, err = serviceEndpoint.ReceiveAny()
		if err != nil {
			t.Error(err)
		}
		close(testSync)
		m, ok = <-msgChan
		if !ok {
			t.Fatalf("connection closed")
		}
		if m.Header.Type != net.Cancel {
			t.Errorf("not a cancel")
		}
		m.Header.Type = net.Cancelled
		serviceEndpoint.Send(net.NewMessage(m.Header, []byte{}))
	}()

	// client connection
	c := bus.NewClient(bus.NewContext(clientEndpoint))

	cancel := make(chan struct{})
	go func() {
		// wait until the service has received the call
		<-testSync
		close(cancel)
	}()
	_, err = c.Call(cancel, 1, 2, 3, []byte{0xab, 0xcd})
	if err != bus.ErrCancelled {
		t.Errorf("proxy failed to call service: %s", err)
	}
}

func TestSelectEndPoint(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	// shall connect to unix socket
	naddr, channel, err := bus.SelectEndPoint([]string{
		"tcp://198.18.0.1:12",
		addr,
		"tcps://192.168.0.1:12",
	}, "", "")
	if err != nil {
		t.Error(err)
	}
	if naddr != addr {
		t.Error("unexpected address")
	}
	endpoint := channel.EndPoint()
	defer endpoint.Close()
	// shall refuse to connect
	naddr, channel, err = bus.SelectEndPoint([]string{
		"tcp://198.18.1.0",
		"tcps://192.168.0.0",
	}, "", "")
	if err == nil {
		t.Error("shall not be able to connect")
	}
	if channel != nil {
		t.Error("non empty channel")
	}
	if naddr != "" {
		t.Error("non empty address")
	}
	// shall refuse to connect to empty list
	_, channel, err = bus.SelectEndPoint(make([]string, 0), "", "")
	if err == nil {
		t.Fatalf("empty list")
	}
	if channel != nil {
		t.Error("non empty endpoint")
	}
	if naddr != "" {
		t.Error("non empty address")
	}
}

func TestClientDisconnectionError(t *testing.T) {

	serviceEndpoint, clientEndpoint := net.Pipe()
	defer clientEndpoint.Close()

	// client connection
	c := bus.NewClient(bus.NewContext(clientEndpoint))

	var wait sync.WaitGroup
	var disconnectError error

	wait.Add(1)
	cont := func(err error) {
		disconnectError = err
		wait.Done()
	}
	c.OnDisconnect(cont)

	serviceEndpoint.Close()
	wait.Wait()
	if disconnectError == nil {
		t.Error("expecting a disconnection error")
	}
}

func TestClientState(t *testing.T) {

	_, clientEndpoint := net.Pipe()
	defer clientEndpoint.Close()

	// client connection
	c := bus.NewClient(bus.NewContext(clientEndpoint))

	if c.State("boom", 1) != 1 {
		t.Error("incorrect value")
	}

	if c.State("boom", -1) != 0 {
		t.Error("incorrect value")
	}

	if c.State("boom", -1) != -1 {
		t.Error("incorrect value")
	}

	if c.State("boom", 1) != 0 {
		t.Error("incorrect value")
	}

	if c.State("bam", 0) != 0 {
		t.Error("incorrect value")
	}
}

func TestClientDisconnectionSuccess(t *testing.T) {

	serviceEndpoint, clientEndpoint := net.Pipe()
	defer serviceEndpoint.Close()

	// client connection
	c := bus.NewClient(bus.NewContext(clientEndpoint))

	var wait sync.WaitGroup
	var err error

	wait.Add(1)
	cont := func(disconnectError error) {
		err = disconnectError
		wait.Done()
	}
	c.OnDisconnect(cont)

	clientEndpoint.Close()
	wait.Wait()
	if err != nil {
		t.Errorf("not expecting a disconnection error: %s", err)
	}
}
