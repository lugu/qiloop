package bus_test

import (
	"sync"
	"testing"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/directory"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
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
	c := bus.NewClient(clientEndpoint)
	_, err = c.Call(nil, 1, 2, 3, []byte{0xab, 0xcd})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestProxyCall(t *testing.T) {

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
	c := bus.NewClient(clientEndpoint)
	proxy := bus.NewProxy(c, object.MetaService0, 1, 2)
	_, err = proxy.CallID(3, []byte{0xab, 0xcd})
	if err != nil {
		t.Errorf("proxy failed to call service: %s", err)
	}
}

func TestClientAlreadyCancelled(t *testing.T) {

	serviceEndpoint, clientEndpoint := net.Pipe()
	defer serviceEndpoint.Close()
	defer clientEndpoint.Close()

	cancel := make(chan struct{})
	close(cancel)
	c := bus.NewClient(clientEndpoint)
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
	c := bus.NewClient(clientEndpoint)

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

func TestProxy(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	session := server.Session()
	directory, err := services.ServiceDirectory(session)
	if err != nil {
		t.Error(err)
	}
	signalID, err := directory.Proxy().MetaObject().SignalID("serviceAdded",
		"(Is)<serviceAdded,serviceID,name>")
	if err != nil {
		t.Error(err)
	}
	if signalID != 106 {
		t.Fatalf("wrong signal id")
	}
	methodID, err := directory.Proxy().MetaObject().MethodID("services",
		"()", "[(sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>]")
	if err != nil {
		t.Error(err)
	}
	if methodID != 101 {
		t.Fatalf("wrong method id")
	}
	if directory.Proxy().ObjectID() != 1 {
		t.Fatalf("wrong object id")
	}
	if directory.Proxy().ServiceID() != 1 {
		t.Fatalf("wrong service id")
	}
	cancel, _, err := directory.Proxy().SubscribeID(signalID)
	if err != nil {
		t.Error(err)
	}
	cancel()
	signalID, err = directory.Proxy().MetaObject().SignalID("serviceAdded",
		"(Is)<serviceAdded,serviceID,name>")
	if err != nil {
		t.Error(err)
	}
	cancel, _, err = directory.Proxy().SubscribeID(signalID)
	if err != nil {
		t.Error(err)
	}
	cancel()
	_, _, err = directory.Proxy().SubscribeID(12345)
	if err == nil {
		// TODO
	}
	methodID, err = directory.Proxy().MetaObject().MethodID("services", "()", "[(sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>]")
	if err != nil {
		t.Error(err)
	}
	resp, err := directory.Proxy().CallID(methodID, make([]byte, 0))
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Error(err)
	}
	if len(resp) == 0 {
		t.Error(err)
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
	naddr, endpoint, err := bus.SelectEndPoint([]string{
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
	defer endpoint.Close()
	// shall refuse to connect
	naddr, endpoint, err = bus.SelectEndPoint([]string{
		"tcp://198.18.1.0",
		"tcps://192.168.0.0",
	}, "", "")
	if err == nil {
		t.Error("shall not be able to connect")
	}
	if endpoint != nil {
		t.Error("non empty endpoint")
	}
	if naddr != "" {
		t.Error("non empty address")
	}
	// shall refuse to connect to empty list
	_, _, err = bus.SelectEndPoint(make([]string, 0), "", "")
	if err == nil {
		t.Fatalf("empty list")
	}
	if endpoint != nil {
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
	c := bus.NewClient(clientEndpoint)

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
	c := bus.NewClient(clientEndpoint)

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
	c := bus.NewClient(clientEndpoint)

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
