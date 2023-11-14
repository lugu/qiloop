package bus_test

import (
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
			panic("connection closed")
		}
		m.Header.Type = net.Reply
		err := serviceEndpoint.Send(*m)
		if err != nil {
			panic(err)
		}
	}()

	// client connection
	c := bus.NewClient(bus.NewContext(clientEndpoint))
	proxy := bus.NewProxy(c, object.MetaService0, 1, 2)
	_, err = proxy.CallID(3, []byte{0xab, 0xcd})
	if err != nil {
		t.Errorf("proxy failed to call service: %s", err)
	}
}

func TestProxySubscribe(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	session := server.Session()
	dir, err := services.ServiceDirectory(session)
	if err != nil {
		t.Error(err)
	}
	signalID, err := dir.Proxy().MetaObject().SignalID("serviceAdded",
		"(Is)<serviceAdded,serviceID,name>")
	if err != nil {
		t.Error(err)
	}
	if signalID != 106 {
		t.Fatalf("wrong signal id")
	}
	cancel, addChan, err := dir.Proxy().SubscribeID(signalID)
	if err != nil {
		t.Error(err)
	}
	cancel()
	_, ok := <-addChan
	if ok {
		t.Error("unexpected event")
	}

	cancel, _, err = dir.Proxy().SubscribeID(signalID)
	if err != nil {
		t.Error(err)
	}
	cancel()
	_, _, err = dir.Proxy().SubscribeID(12345)
	if err == nil {
		t.Error("unexpected success")
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
	dir, err := services.ServiceDirectory(session)
	if err != nil {
		t.Error(err)
	}
	methodID, _, err := dir.Proxy().MetaObject().MethodID("services", "()")
	if err != nil {
		t.Error(err)
	}
	if methodID != 101 {
		t.Fatalf("wrong method id")
	}
	if dir.Proxy().ObjectID() != 1 {
		t.Fatalf("wrong object id")
	}
	if dir.Proxy().ServiceID() != 1 {
		t.Fatalf("wrong service id")
	}
	methodID, _, err = dir.Proxy().MetaObject().MethodID("services", "()")
	if err != nil {
		t.Error(err)
	}
	resp, err := dir.Proxy().CallID(methodID, make([]byte, 0))
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
