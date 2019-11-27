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

	// 4. send a message
	proxy := bus.NewProxy(c, object.MetaService0, 1, 2)
	_, err = proxy.CallID(3, []byte{0xab, 0xcd})
	if err != nil {
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
	services := services.Services(session)
	directory, err := services.ServiceDirectory(nil)
	if err != nil {
		t.Error(err)
	}
	signalID, err := directory.SignalID("serviceAdded")
	if err != nil {
		t.Error(err)
	}
	if signalID != 106 {
		t.Fatalf("wrong signal id")
	}
	methodID, err := directory.MethodID("services")
	if err != nil {
		t.Error(err)
	}
	if methodID != 101 {
		t.Fatalf("wrong method id")
	}
	if directory.ObjectID() != 1 {
		t.Fatalf("wrong object id")
	}
	if directory.ServiceID() != 1 {
		t.Fatalf("wrong service id")
	}
	cancel, _, err := directory.SubscribeID(signalID)
	if err != nil {
		t.Error(err)
	}
	cancel()
	cancel, _, err = directory.Subscribe("serviceAdded")
	if err != nil {
		t.Error(err)
	}
	cancel()
	_, _, err = directory.Subscribe("unknownSignal")
	if err == nil {
		t.Fatalf("must fail")
	}
	_, _, err = directory.SubscribeID(12345)
	if err == nil {
		// TODO
	}
	_, err = directory.Call("unknown service", make([]byte, 0))
	if err == nil {
		t.Fatalf("must fail")
	}
	resp, err := directory.Call("services", make([]byte, 0))
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
