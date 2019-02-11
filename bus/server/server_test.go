package server_test

import (
	"bytes"
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/client/services"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/server/directory"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
	"testing"
)

// ObjectDispatcher implements ServerObject
type ObjectDispatcher struct {
	wrapper server.Wrapper
}

func (o *ObjectDispatcher) UpdateSignal(signal uint32, value []byte) error {
	panic("not available")
}

func (o *ObjectDispatcher) Wrap(id uint32, fn server.ActionWrapper) {
	if o.wrapper == nil {
		o.wrapper = make(map[uint32]server.ActionWrapper)
	}
	o.wrapper[id] = fn
}

func (o *ObjectDispatcher) Activate(activation server.Activation) error {
	return nil
}

func (o *ObjectDispatcher) OnTerminate() {
}

func (o *ObjectDispatcher) Receive(m *net.Message, from *server.Context) error {
	if o.wrapper == nil {
		return util.ReplyError(from.EndPoint, m, server.ErrActionNotFound)
	}
	a, ok := o.wrapper[m.Header.Action]
	if !ok {
		return util.ReplyError(from.EndPoint, m, server.ErrActionNotFound)
	}
	response, err := a(m.Payload)

	if err != nil {
		return util.ReplyError(from.EndPoint, m, err)
	}
	reply := net.NewMessage(m.Header, response)
	reply.Header.Type = net.Reply
	return from.EndPoint.Send(reply)
}

func newObject() server.ServerObject {
	var object ObjectDispatcher
	handler := func(d []byte) ([]byte, error) {
		return []byte{0xab, 0xcd}, nil
	}
	object.Wrap(3, handler)
	return &object
}

func TestNewServer(t *testing.T) {

	addr := util.NewUnixAddr()

	srv, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Terminate()

	object := newObject()
	srv.NewService("test", object)

	clt, err := net.DialEndPoint(addr)
	if err != nil {
		t.Error(err)
	}

	err = client.Authenticate(clt)
	if err != nil {
		t.Error(err)
	}

	h := net.NewHeader(net.Call, 2, 1, 3, 4)
	mSent := net.NewMessage(h, make([]byte, 0))

	received, err := clt.ReceiveAny()
	if err != nil {
		t.Errorf("failed to receive net. %s", err)
	}

	// client send a message
	if err := clt.Send(mSent); err != nil {
		t.Errorf("failed to send paquet: %s", err)
	}

	// server replied
	mReceived, ok := <-received
	if !ok {
		t.Fatalf("connection closed")
	}

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

	srv, err := server.StandAloneServer(listener, server.Yes{},
		server.PrivateNamespace())
	if err != nil {
		t.Error(err)
	}
	srv.NewService("test", obj)

	cltNet, err := net.DialEndPoint(addr)
	if err != nil {
		t.Error(err)
	}

	err = client.Authenticate(cltNet)
	if err != nil {
		t.Error(err)
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
	testInvalid(t, server.ErrServiceNotFound, invalidServiceID, objectID, actionID)
	testInvalid(t, server.ErrObjectNotFound, serviceID, invalidObjectID, actionID)
	testInvalid(t, server.ErrActionNotFound, serviceID, objectID, invalidActionID)

	testInvalid(t, server.ErrObjectNotFound, serviceID, invalidObjectID, invalidActionID)
	testInvalid(t, server.ErrServiceNotFound, invalidServiceID, objectID, invalidActionID)
	testInvalid(t, server.ErrServiceNotFound, invalidServiceID, invalidObjectID, actionID)
	srv.Terminate()
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

	srv, err := server.StandAloneServer(listener, server.Yes{},
		server.PrivateNamespace())
	if err != nil {
		t.Error(err)
	}
	defer srv.Terminate()
	srv.NewService("test", obj)

	clt, err := net.DialEndPoint(addr)
	if err != nil {
		t.Error(err)
	}

	// construct a proxy object
	serviceID := uint32(0x1)
	objectID := uint32(0x1)
	proxy := client.NewProxy(client.NewClient(clt),
		object.ObjectMetaObject, serviceID, objectID)

	// register to signal
	id, err := proxy.MethodID("registerEvent")
	if err != nil {
		t.Errorf("proxy get register event action id: %s", err)
	}
	if id != 0 {
		t.Errorf("invalid action id")
	}
}

func TestSession(t *testing.T) {
	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	obj := server.NewObject(object.MetaObject{
		Description: "test",
		Methods:     make(map[uint32]object.MetaMethod),
		Signals:     make(map[uint32]object.MetaSignal),
	})

	srv, err := server.StandAloneServer(listener, server.Yes{},
		server.PrivateNamespace())
	if err != nil {
		t.Error(err)
	}
	defer srv.Terminate()
	srv.NewService("test", obj)

	sess := srv.Session()
	proxy, err := sess.Proxy("test", 1)
	if err != nil {
		t.Error(err)
	}

	metaBytes, err := proxy.Call("metaObject", []byte{1, 0, 0, 0})
	if err != nil {
		t.Error(err)
	}
	buf := bytes.NewBuffer(metaBytes)
	_, err = object.ReadMetaObject(buf)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoteServer(t *testing.T) {

	addr := util.NewUnixAddr()
	srv, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Terminate()
	obj1 := server.NewObject(object.MetaObject{
		Description: "service1",
		Methods:     make(map[uint32]object.MetaMethod),
		Signals:     make(map[uint32]object.MetaSignal),
	})
	_, err = srv.NewService("service1", obj1)
	if err != nil {
		t.Error(err)
	}
	sess := srv.Session()
	services := services.Services(sess)
	directory, err := services.ServiceDirectory()
	if err != nil {
		t.Fatalf("failed to connect directory: %s", err)
	}
	info, err := directory.Service("service1")
	if err != nil {
		t.Error(err)
	}
	if info.ServiceId != 2 {
		t.Errorf("unexpected service id: %d", info.ServiceId)
	}
	proxy1, err := sess.Proxy("service1", 1)
	if err != nil {
		t.Fatal(err)
	}

	metaBytes, err := proxy1.Call("metaObject", []byte{1, 0, 0, 0})
	if err != nil {
		t.Error(err)
	}
	buf := bytes.NewBuffer(metaBytes)
	_, err = object.ReadMetaObject(buf)
	if err != nil {
		t.Error(err)
	}

	addr2 := util.NewUnixAddr()
	srv2, err := server.NewServer(sess, addr2, server.Yes{})
	if err != nil {
		t.Error(err)
	}
	defer srv2.Terminate()
	sess2 := srv2.Session()

	obj2 := server.NewObject(object.MetaObject{
		Description: "service2",
		Methods:     make(map[uint32]object.MetaMethod),
		Signals:     make(map[uint32]object.MetaSignal),
	})
	_, err = srv2.NewService("service2", obj2)
	if err != nil {
		t.Error(err)
	}
	info, err = directory.Service("service2")
	if err != nil {
		t.Error(err)
	}
	if info.ServiceId != 3 {
		t.Errorf("unexpected service id: %d", info.ServiceId)
	}
	if info.Endpoints[0] != addr2 {
		t.Errorf("unexpected address: %s", info.Endpoints[0])
	}
	// access service 2 through the session of server 1
	proxy2, err := sess.Proxy("service2", 1)
	if err != nil {
		t.Fatal(err)
	}

	metaBytes, err = proxy2.Call("metaObject", []byte{1, 0, 0, 0})
	if err != nil {
		t.Error(err)
	}
	buf = bytes.NewBuffer(metaBytes)
	_, err = object.ReadMetaObject(buf)
	if err != nil {
		t.Error(err)
	}

	// access service 2 through the session of server 2
	proxy2, err = sess2.Proxy("service2", 1)
	if err != nil {
		t.Fatal(err)
	}

	metaBytes, err = proxy2.Call("metaObject", []byte{1, 0, 0, 0})
	if err != nil {
		t.Error(err)
	}
	buf = bytes.NewBuffer(metaBytes)
	_, err = object.ReadMetaObject(buf)
	if err != nil {
		t.Error(err)
	}
}

func TestNamespace(t *testing.T) {
	ns := server.PrivateNamespace()
	_, err := ns.Reserve("")
	if err == nil {
		t.Fatalf("shall fail")
	}
	_, err = ns.Resolve("test2")
	if err == nil {
		t.Fatalf("shall fail")
	}
	err = ns.Enable(1)
	if err == nil {
		t.Fatalf("shall fail")
	}
	err = ns.Remove(1)
	if err == nil {
		t.Fatalf("shall fail")
	}
	uid, err := ns.Reserve("test2")
	if err != nil {
		t.Error(err)
	}
	_, err = ns.Reserve("test2")
	if err == nil {
		t.Error(err)
	}
	err = ns.Enable(uid)
	if err != nil {
		t.Error(err)
	}
	uid2, err := ns.Resolve("test2")
	if err != nil {
		t.Error(err)
	}
	if uid != uid2 {
		t.Error(err)
	}
	err = ns.Remove(uid)
	if err != nil {
		t.Error(err)
	}
}

func TestLocalSession(t *testing.T) {
	addr := util.NewUnixAddr()
	srv, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Terminate()

	ns := server.PrivateNamespace()
	// activate service directory
	sess := ns.Session(srv)
	_, err = sess.Proxy("ServiceDirectory", 1)
	if err == nil {
		t.Errorf("service directory not registered")
	}
	uid, err := ns.Reserve("ServiceDirectory")
	if err != nil {
		t.Error(err)
	}
	if uid != 1 {
		t.Errorf("expecting first id, got %d", uid)
	}
	err = ns.Enable(uid)
	if err != nil {
		t.Error(err)
	}
	_, err = sess.Proxy("ServiceDirectory", 1)
	if err != nil {
		t.Error(err)
	}
	_, err = sess.Proxy("ServiceDirectory", 2)
	if err == nil {
		t.Errorf("shall not have an second object")
	}
	_, err = sess.Proxy("unknown", 1)
	if err == nil {
		t.Errorf("shall fail: wrong object id")
	}
	srv.NewService("known", newObject())
	_, err = sess.Proxy("known", 1)
	if err == nil {
		t.Errorf("shall fail: object does not implement metaObject()")
	}

	sess.Destroy()
}

func TestServer(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	session, err := session.NewSession(addr)
	if err != nil {
		t.Error(err)
	}
	info := services.ServiceInfo{
		Name:      "test",
		MachineId: util.MachineID(),
		ProcessId: util.ProcessID(),
		Endpoints: []string{
			util.NewUnixAddr(),
		},
	}
	services := services.Services(session)
	directory, err := services.ServiceDirectory()
	if err != nil {
		t.Fatalf("failed to create directory: %s", err)
	}
	list, err := directory.Services()
	if err != nil {
		t.Error(err)
	}
	if len(list) != 1 {
		t.Errorf("not expecting %d", len(list))
	}

	cancelAdded, added, err := directory.SubscribeServiceAdded()
	if err != nil {
		t.Error(err)
	}
	cancelRemoved, removed, err := directory.SubscribeServiceRemoved()
	if err != nil {
		t.Error(err)
	}
	info.ServiceId, err = directory.RegisterService(info)
	if err != nil {
		t.Error(err)
	}
	_, err = directory.RegisterService(info)
	if err == nil {
		t.Fatalf("shall fail")
	}
	err = directory.ServiceReady(info.ServiceId)
	if err != nil {
		t.Error(err)
	}
	err = directory.ServiceReady(info.ServiceId)
	if err == nil {
		t.Fatalf("shall fail")
	}
	_, err = directory.Service("test")
	if err != nil {
		t.Error(err)
	}
	_, err = directory.Service("test2")
	if err == nil {
		t.Fatalf("shall fail")
	}
	info.ProcessId = 2
	err = directory.UpdateServiceInfo(info)
	if err != nil {
		t.Error(err)
	}
	info.ProcessId = 0
	err = directory.UpdateServiceInfo(info)
	if err == nil {
		t.Fatalf("shall fail")
	}
	err = directory.UnregisterService(info.ServiceId)
	if err != nil {
		t.Error(err)
	}
	err = directory.UnregisterService(info.ServiceId)
	if err == nil {
		t.Fatalf("shall fail")
	}
	info2, ok := <-added
	if !ok {
		t.Fatalf("unexpected")
	}
	if info2.Name != "test" {
		t.Fatalf(info.Name)
	}
	info3, ok := <-removed
	if !ok {
		t.Fatalf("unexpected")
	}
	if info3.Name != "test" {
		t.Fatalf(info.Name)
	}
	cancelAdded()
	cancelRemoved()
	_, ok = <-added
	if ok {
		t.Fatalf("unexpected added ok")
	}
	_, ok = <-removed
	if ok {
		t.Fatalf("unexpected removed ok")
	}
}

func TestServiceImpl(t *testing.T) {
	service := server.NewService(newObject())
	uid, err := service.Add(newObject())
	if err != nil {
		t.Error(err)
	}
	err = service.Remove(uid)
	if err != nil {
		t.Error(err)
	}
	err = service.Remove(uid + 1)
	if err == nil {
		t.Fatalf("shall fail")
	}
}

func TestRouter(t *testing.T) {
	service0 := server.ServiceAuthenticate(server.Yes{})
	namespace := server.PrivateNamespace()
	router := server.NewRouter(service0, namespace)
	object := newObject()
	service := server.NewService(object)
	err := router.Add(0, service)
	if err == nil {
		t.Error(err)
	}
	err = router.Add(1, service)
	if err != nil {
		t.Error(err)
	}
	err = router.Remove(1)
	if err != nil {
		t.Error(err)
	}
	err = router.Remove(1)
	if err == nil {
		t.Errorf("shall fail")
	}
	err = router.Terminate()
	if err != nil {
		t.Error(err)
	}
}

func TestNewContext(t *testing.T) {
	endpoint, _ := net.Pipe()
	context := server.NewContext(endpoint)
	if context == nil {
		t.Errorf("nil context")
	}
	if context.Authenticated == true {
		t.Errorf("shall not")
	}
}

type ObjectTerminaison struct {
	server.BasicObject
	Terminated bool
}

// OnTerminate is called when the object is terminated.
func (o *ObjectTerminaison) OnTerminate() {
	o.Terminated = true
}

func TestObjectTerminaison(t *testing.T) {

	basicObj := server.NewObject(object.MetaObject{
		Description: "test",
		Methods:     make(map[uint32]object.MetaMethod),
		Signals:     make(map[uint32]object.MetaSignal),
	})
	obj := &ObjectTerminaison{
		*basicObj,
		false,
	}

	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	srv, err := server.StandAloneServer(listener, server.Yes{},
		server.PrivateNamespace())
	if err != nil {
		t.Error(err)
	}
	srv.NewService("service1", obj)
	srv.Terminate()
	if obj.Terminated == false {
		t.Error("obj shall have terminated")
	}
}

func TestServiceTerminaison(t *testing.T) {
	basicObj := server.NewObject(object.MetaObject{
		Description: "test",
		Methods:     make(map[uint32]object.MetaMethod),
		Signals:     make(map[uint32]object.MetaSignal),
	})
	obj := &ObjectTerminaison{
		*basicObj,
		false,
	}

	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	ns := server.PrivateNamespace()
	srv, err := server.StandAloneServer(listener, server.Yes{}, ns)
	if err != nil {
		t.Error(err)
	}
	service, err := srv.NewService("service1", obj)
	if err != nil {
		t.Error(err)
	}
	uid, err := ns.Resolve("service1")
	if err != nil {
		t.Error(err)
	}
	if uid != 1 {
		t.Errorf("unexpectped uid: %d", uid)
	}
	service.Terminate()
	if obj.Terminated == false {
		t.Error("obj shall have terminated")
	}
	_, err = ns.Resolve("service1")
	if err == nil {
		t.Error("shall have failed")
	}
}
