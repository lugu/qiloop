package directory_test

import (
	"bytes"
	proxy "github.com/lugu/qiloop/bus/client/services"
	dir "github.com/lugu/qiloop/bus/server/directory"
	sess "github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
	"io"
	"testing"
)

func TestNewServer(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := dir.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	session, err := sess.NewSession(addr)
	if err != nil {
		t.Error(err)
	}
	services := proxy.Services(session)
	directory, err := services.ServiceDirectory()
	if err != nil {
		t.Fatalf("failed to create directory: %s", err)
	}
	machineID, err := directory.MachineId()
	if err != nil {
		t.Error(err)
	}
	if machineID == "" {
		t.Fatalf("empty machine id")
	}
}

type mockServiceDirectorySignalHelper struct {
	t            *testing.T
	expectedName string
	expectedUid  uint32
}

func newServiceDirectorySignalHelper(t *testing.T, name string, uid uint32) dir.ServiceDirectorySignalHelper {
	return &mockServiceDirectorySignalHelper{
		t:            t,
		expectedName: name,
		expectedUid:  uid,
	}
}
func (h *mockServiceDirectorySignalHelper) SignalServiceAdded(P0 uint32, P1 string) error {
	return nil
}

func (h *mockServiceDirectorySignalHelper) SignalServiceRemoved(P0 uint32, P1 string) error {
	return nil
}

func newObjectRef(serviceID uint32) object.ObjectReference {
	return object.ObjectReference{
		Boolean:    true,
		MetaObject: object.ObjectMetaObject,
		ParentID:   0,
		ServiceID:  serviceID,
		ObjectID:   1,
	}
}

func newClientInfo(name string) proxy.ServiceInfo {
	return proxy.ServiceInfo{
		Name:      name,
		MachineId: util.MachineID(),
		ProcessId: util.ProcessID(),
		Endpoints: []string{
			util.NewUnixAddr(),
		},
	}
}

func newInfo(name string) dir.ServiceInfo {
	return dir.ServiceInfo{
		Name:      name,
		MachineId: util.MachineID(),
		ProcessId: util.ProcessID(),
		Endpoints: []string{
			util.NewUnixAddr(),
		},
	}
}

func compareInfo(t *testing.T, observed, expected dir.ServiceInfo) {
	if observed.Name != expected.Name {
		t.Errorf("unexpected expected: %#v, expecting %#v", observed, expected)
	}
	if observed.ServiceId != expected.ServiceId {
		t.Errorf("unexpected expected: %#v, expecting %#v", observed, expected)
	}
	if observed.MachineId != expected.MachineId {
		t.Errorf("unexpected expected: %#v, expecting %#v", observed, expected)
	}
	if observed.ProcessId != expected.ProcessId {
		t.Errorf("unexpected expected: %#v, expecting %#v", observed, expected)
	}
	if len(observed.Endpoints) != len(expected.Endpoints) {
		t.Errorf("unexpected expected: %#v, expecting %#v", observed, expected)
	}
	if observed.Endpoints[0] != expected.Endpoints[0] {
		t.Errorf("unexpected expected: %#v, expecting %#v", observed, expected)
	}
}

func TestServerDirectory(t *testing.T) {
	helper := newServiceDirectorySignalHelper(t, "test", 2)
	impl := dir.NewServiceDirectory()
	impl.Activate(nil, 1, 1, helper)
	info := newInfo("test")
	uid, err := impl.RegisterService(info)
	if err != nil {
		t.Error(err)
	}
	info.ServiceId = uid
	err = impl.ServiceReady(uid)
	if err != nil {
		t.Error(err)
	}
	// shall not be able to register twice with the same name
	_, err = impl.RegisterService(info)
	if err == nil {
		t.Error(err)
	}
	services, err := impl.Services()
	if err != nil {
		t.Error(err)
	}
	if len(services) != 1 {
		t.Errorf("wrong service number: %d", len(services))
	}
	compareInfo(t, services[0], info)
	info2 := newInfo("test2")
	uid, err = impl.RegisterService(info2)
	if err != nil {
		t.Error(err)
	}
	info2.ServiceId = uid
	err = impl.ServiceReady(uid)
	if err != nil {
		t.Error(err)
	}
	// shall not be able to register twice with the same name
	_, err = impl.RegisterService(info2)
	if err == nil {
		t.Error(err)
	}
	services, err = impl.Services()
	if err != nil {
		t.Error(err)
	}
	if len(services) != 2 {
		t.Errorf("wrong service number: %d", len(services))
	}
	compareInfo(t, services[0], info)
	compareInfo(t, services[1], info2)

	err = impl.UnregisterService(info2.ServiceId + 1)
	if err == nil {
		t.Fatalf("shall fail")
	}
	err = impl.UnregisterService(info.ServiceId)
	if err != nil {
		t.Error(err)
	}
	err = impl.UnregisterService(info.ServiceId)
	if err == nil {
		t.Error(err)
	}
	services, err = impl.Services()
	if err != nil {
		t.Error(err)
	}
	if len(services) != 1 {
		t.Errorf("wrong service number: %d", len(services))
	}
	compareInfo(t, services[0], info2)

}

func TestServiceDirectoryInfo(t *testing.T) {
	helper := newServiceDirectorySignalHelper(t, "test", 2)
	impl := dir.NewServiceDirectory()
	impl.Activate(nil, 1, 1, helper)
	info := newInfo("test")
	uid, err := impl.RegisterService(info)
	if err != nil {
		t.Error(err)
	}
	info.ServiceId = uid
	_, err = impl.RegisterService(info)
	if err == nil {
		t.Fatalf("already registered")
	}
	err = impl.ServiceReady(uid)
	if err != nil {
		t.Error(err)
	}
	info2, err := impl.Service("test")
	if err != nil {
		t.Error(err)
	}
	compareInfo(t, info, info2)
	info.MachineId = "test"
	err = impl.UpdateServiceInfo(info)
	if err != nil {
		t.Error(err)
	}

	info2 = newInfo("test2")
	info2.ServiceId = uid
	err = impl.UpdateServiceInfo(info2)
	if err == nil {
		t.Fatalf("shall not accecpt update")
	}
	info2 = newInfo("test2")
	info2.ServiceId = uid + 1
	err = impl.UpdateServiceInfo(info2)
	if err == nil {
		t.Fatalf("shall not accecpt name update")
	}
	_, err = impl.Service("test2")
	if err == nil {
		t.Fatalf("invalid name")
	}
	err = impl.ServiceReady(uid + 1)
	if err == nil {
		t.Fatalf("invalid uid")
	}
	info.Name = ""
	_, err = impl.RegisterService(info)
	if err == nil {
		t.Fatalf("shall reject empty name")
	}
	info = newInfo("test")
	info.MachineId = ""
	_, err = impl.RegisterService(info)
	if err == nil {
		t.Fatalf("shall reject empty machine info")
	}
	info = newInfo("test")
	info.ProcessId = 0
	_, err = impl.RegisterService(info)
	if err == nil {
		t.Fatalf("shall reject empty process info")
	}
	info = newInfo("test")
	info.Endpoints = make([]string, 0)
	_, err = impl.RegisterService(info)
	if err == nil {
		t.Fatalf("shall reject empty endpoint info")
	}
	info.Endpoints = []string{""}
	_, err = impl.RegisterService(info)
	if err == nil {
		t.Fatalf("shall reject empty endpoint info")
	}
	info = newInfo("test2")
	uid, err = impl.RegisterService(info)
	if err != nil {
		t.Error(err)
	}
	err = impl.UnregisterService(uid)
	if err != nil {
		t.Error(err)
	}
}
func TestNamespace(t *testing.T) {
	helper := newServiceDirectorySignalHelper(t, "test", 2)
	impl := dir.NewServiceDirectory()
	impl.Activate(nil, 1, 1, helper)
	info := newInfo("test")
	uid, err := impl.RegisterService(info)
	if err != nil {
		t.Error(err)
	}
	ns := impl.Namespace("dummy")
	_, err = ns.Reserve("test")
	if err == nil {
		t.Fatalf("shall fail")
	}
	_, err = ns.Resolve("test2")
	if err == nil {
		t.Fatalf("shall fail")
	}
	err = ns.Enable(uid + 1)
	if err == nil {
		t.Fatalf("shall fail")
	}
	err = ns.Remove(uid + 1)
	if err == nil {
		t.Fatalf("shall fail")
	}
	uid, err = ns.Reserve("test2")
	if err != nil {
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

func TestStub(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := dir.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	session, err := sess.NewSession(addr)
	if err != nil {
		t.Error(err)
	}
	services := proxy.Services(session)
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

	_, added, err := directory.SubscribeServiceAdded()
	if err != nil {
		t.Error(err)
	}
	_, removed, err := directory.SubscribeServiceRemoved()
	if err != nil {
		t.Error(err)
	}
	info := newClientInfo("test")
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
	if info2.P1 != "test" {
		t.Fatalf(info.Name)
	}
	info2, ok = <-removed
	if !ok {
		t.Fatalf("unexpected")
	}
	if info2.P1 != "test" {
		t.Fatalf(info.Name)
	}
}

func TestSession(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := dir.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	session := server.Session()
	if err != nil {
		t.Error(err)
	}
	defer session.Destroy()
	services := proxy.Services(session)
	directory, err := services.ServiceDirectory()
	if err != nil {
		t.Fatalf("failed to create directory: %s", err)
	}
	machineID, err := directory.MachineId()
	if err != nil {
		t.Error(err)
	}
	if machineID == "" {
		t.Fatalf("empty machine id")
	}
	_, err = session.Proxy("test", 1)
	if err == nil {
		t.Fatalf("shall fail")
	}

	info := newClientInfo("test")
	info.Endpoints = []string{addr}
	uid, err := directory.RegisterService(info)
	if err != nil {
		t.Error(err)
	}
	err = directory.ServiceReady(uid)
	if err != nil {
		t.Error(err)
	}
	_, err = session.Proxy("test", 1)
	if err == nil {
		t.Fatalf("shall fail")
	}
}

func LimitedReader(s dir.ServiceInfo, size int) io.Reader {
	var buf bytes.Buffer
	err := dir.WriteServiceInfo(s, &buf)
	if err != nil {
		panic(err)
	}
	return &io.LimitedReader{
		R: &buf,
		N: int64(size),
	}
}

type LimitedWriter struct {
	size int
}

func (b *LimitedWriter) Write(buf []byte) (int, error) {
	if len(buf) <= b.size {
		b.size -= len(buf)
		return len(buf), nil
	}
	old_size := b.size
	b.size = 0
	return old_size, io.EOF
}

func NewLimitedWriter(size int) io.Writer {
	return &LimitedWriter{
		size: size,
	}
}

func TestWriterServiceInfo(t *testing.T) {
	info := newInfo("test")
	var buf bytes.Buffer
	err := dir.WriteServiceInfo(info, &buf)
	if err != nil {
		t.Error(err)
	}
	max := len(buf.Bytes())

	for i := 0; i < max-1; i++ {
		w := NewLimitedWriter(i)
		err := dir.WriteServiceInfo(info, w)
		if err == nil {
			t.Fatalf("not expecting a success at %d", i)
		}
	}
	w := NewLimitedWriter(max)
	err = dir.WriteServiceInfo(info, w)
	if err != nil {
		t.Error(err)
	}
}

func TestReadHeaderError(t *testing.T) {
	info := newInfo("test")
	var buf bytes.Buffer
	err := dir.WriteServiceInfo(info, &buf)
	if err != nil {
		t.Error(err)
	}
	max := len(buf.Bytes())

	for i := 0; i < max; i++ {
		r := LimitedReader(info, i)
		_, err := dir.ReadServiceInfo(r)
		if err == nil {
			t.Fatalf("not expecting a success at %d", i)
		}
	}
	r := LimitedReader(info, max)
	_, err = dir.ReadServiceInfo(r)
	if err != nil {
		t.Error(err)
	}
}

func TestStandalone(t *testing.T) {
	_, err := dir.NewServer("", nil)
	if err == nil {
		t.Fatalf("shall fail")
	}
}
