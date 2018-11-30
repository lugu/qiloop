package client_test

import (
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/client/services"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/server/directory"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
	"testing"
)

func TestProxyCall(t *testing.T) {

	serviceEndpoint, clientEndpoint := net.NewPipe()
	defer serviceEndpoint.Close()
	defer clientEndpoint.Close()

	msgChan, err := serviceEndpoint.ReceiveAny()
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

func TestProxy(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	session := server.Session()
	services := services.Services(session)
	directory, err := services.ServiceDirectory()
	if err != nil {
		panic(err)
	}
	signalID, err := directory.SignalUid("serviceAdded")
	if err != nil {
		panic(err)
	}
	if signalID != 106 {
		panic("wrong signal id")
	}
	methodID, err := directory.MethodUid("services")
	if err != nil {
		panic(err)
	}
	if methodID != 101 {
		panic("wrong method id")
	}
	if directory.ObjectID() != 1 {
		panic("wrong object id")
	}
	if directory.ServiceID() != 1 {
		panic("wrong service id")
	}
	cancel := make(chan int)
	_, err = directory.SubscribeID(signalID, cancel)
	if err != nil {
		panic(err)
	}
	cancel <- 1
	_, err = directory.SubscribeSignal("serviceAdded", cancel)
	if err != nil {
		panic(err)
	}
	cancel <- 1
	_, err = directory.SubscribeSignal("unknownSignal", cancel)
	if err == nil {
		panic("must fail")
	}
	_, err = directory.SubscribeID(12345, cancel)
	if err == nil {
		// TODO
	}
	_, err = directory.Call("unknown service", []byte{})
	if err == nil {
		panic("must fail")
	}
	resp, err := directory.Call("services", []byte{})
	if err != nil {
		panic(err)
	}
	if resp == nil {
		panic(err)
	}
	if len(resp) == 0 {
		panic(err)
	}
	directory.Disconnect()
}

func TestSelectEndPoint(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	// shall connect to unix socket
	endpoint, err := client.SelectEndPoint([]string{
		"tcp://198.18.0.1:12",
		addr,
		"tcps://192.168.0.1:12",
	})
	if err != nil {
		panic(err)
	}
	defer endpoint.Close()
	// shall refuse to connect
	_, err = client.SelectEndPoint([]string{
		"tcp://198.18.1.0",
		"tcps://192.168.0.0",
	})
	if err == nil {
		panic("shall not be able to connect")
	}
}

func TestSelectError(t *testing.T) {

	addr := util.NewUnixAddr()

	listener, err := net.Listen(addr)
	if err != nil {
		panic(err)
	}
	server, err := server.StandAloneServer(listener, server.No{},
		server.PrivateNamespace())
	if err != nil {
		panic(err)
	}
	defer server.Terminate()

	// shall fail to authenticate
	endpoint, err := client.SelectEndPoint([]string{
		"tcp://198.18.0.1:12",
		addr,
		"tcps://192.168.0.1:12",
	})
	defer endpoint.Close()
	if err == nil {
		panic("shall fail to authenticate")
	}
	// shall refuse to connect to empty list
	_, err = client.SelectEndPoint([]string{})
	if err == nil {
		panic("empty list")
	}
}
