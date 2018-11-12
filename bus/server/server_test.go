package server_test

import (
	"bytes"
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
	"testing"
)

func TestNewServer(t *testing.T) {

	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}

	var object server.ObjectDispatcher
	handler := func(d []byte) ([]byte, error) {
		return []byte{0xab, 0xcd}, nil
	}
	object.Wrap(3, handler)

	service := server.NewService(&object)
	router := server.NewRouter(server.ServiceAuthenticate(server.Yes{}))
	router.Add(1, service)
	srv := server.StandAloneServer(listener, router)
	go srv.Run()

	clt, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}

	err = client.Authenticate(clt)
	if err != nil {
		panic(err)
	}

	h := net.NewHeader(net.Call, 1, 1, 3, 4)
	mSent := net.NewMessage(h, make([]byte, 0))

	// client is prepared to receive a message
	received := make(chan *net.Message)
	go func() {
		msg, err := clt.ReceiveAny()
		if err != nil {
			t.Errorf("failed to receive net. %s", err)
		}
		received <- msg
	}()

	// client send a message
	if err := clt.Send(mSent); err != nil {
		t.Errorf("failed to send paquet: %s", err)
	}

	// server replied
	mReceived := <-received

	srv.Stop()

	if mReceived.Header.Type == net.Error {
		buf := bytes.NewBuffer(mReceived.Payload)
		errV, err := value.NewValue(buf)
		if err != nil {
			t.Errorf("invalid error value: %v", mReceived.Payload)
		}
		if str, ok := errV.(value.StringValue); ok {
			t.Errorf("error: %s", string(str))
		} else {
			t.Errorf("invalid error: %v", mReceived.Payload)
		}
	}

	if mSent.Header.ID != mReceived.Header.ID {
		t.Errorf("invalid message id: %d", mReceived.Header.ID)
	}
	if mReceived.Header.Type != net.Reply {
		t.Errorf("invalid message type: %d", mReceived.Header.Type)
	}
	if mSent.Header.Service != mReceived.Header.Service {
		t.Errorf("invalid message service: %d", mReceived.Header.Service)
	}
	if mSent.Header.Object != mReceived.Header.Object {
		t.Errorf("invalid message object: %d", mReceived.Header.Object)
	}
	if mSent.Header.Action != mReceived.Header.Action {
		t.Errorf("invalid message action: %d", mReceived.Header.Action)
	}
	if mReceived.Payload[0] != 0xab || mReceived.Payload[1] != 0xcd {
		t.Errorf("invalid message payload: %d", mReceived.Header.Type)
	}
}

func TestServerReturnError(t *testing.T) {
	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	// Empty meta object:
	obj := server.NewObject(object.MetaObject{
		Description: "test",
		Methods:     make(map[uint32]object.MetaMethod),
		Signals:     make(map[uint32]object.MetaSignal),
	})

	router := server.NewRouter(server.ServiceAuthenticate(server.Yes{}))

	router.Add(1, server.NewService(obj))

	srv := server.StandAloneServer(listener, router)
	go srv.Run()

	cltNet, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}

	err = client.Authenticate(cltNet)
	if err != nil {
		panic(err)
	}

	clt := client.NewClient(cltNet)

	serviceID := uint32(0x1)
	objectID := uint32(0x1)
	actionID := uint32(0x0)
	invalidServiceID := uint32(0x2)
	invalidObjectID := uint32(0x2)
	invalidActionID := uint32(0x100)

	testInvalid := func(t *testing.T, expected error, serviceID, objectID,
		actionID uint32) {
		_, err := clt.Call(serviceID, objectID, actionID, make([]byte, 0))
		if err == nil && expected == nil {
			// ok
		} else if err != nil && expected == nil {
			t.Errorf("unexpected error:\n%s\nexpecting:\nnil", err)
		} else if err == nil && expected != nil {
			t.Errorf("unexpected error:\nnil\nexpecting:\n%s", expected)
		} else if err.Error() != expected.Error() {
			t.Errorf("unexpected error:\n%s\nexpecting:\n%s", err, expected)
		}
	}
	testInvalid(t, server.ServiceNotFound, invalidServiceID, objectID, actionID)
	testInvalid(t, server.ObjectNotFound, serviceID, invalidObjectID, actionID)
	testInvalid(t, server.ActionNotFound, serviceID, objectID, invalidActionID)

	testInvalid(t, server.ObjectNotFound, serviceID, invalidObjectID, invalidActionID)
	testInvalid(t, server.ServiceNotFound, invalidServiceID, objectID, invalidActionID)
	testInvalid(t, server.ServiceNotFound, invalidServiceID, invalidObjectID, actionID)
	srv.Stop()
}

func TestStandAloneInit(t *testing.T) {

	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	// Empty meta object:
	obj := server.NewObject(object.MetaObject{
		Description: "test",
		Methods:     make(map[uint32]object.MetaMethod),
		Signals:     make(map[uint32]object.MetaSignal),
	})
	router := server.NewRouter(server.ServiceAuthenticate(server.Yes{}))
	router.Add(1, server.NewService(obj))

	srv := server.StandAloneServer(listener, router)
	go srv.Run()

	clt, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}

	// construct a proxy object
	serviceID := uint32(0x1)
	objectID := uint32(0x1)
	proxy := client.NewProxy(client.NewClient(clt),
		object.ObjectMetaObject, serviceID, objectID)

	// register to signal
	id, err := proxy.MethodUid("registerEvent")
	if err != nil {
		t.Errorf("proxy get register event action id: %s", err)
	}
	if id != 0 {
		t.Errorf("invalid action id")
	}
	srv.Stop()
}
