package directory_test

import (
	proxy "github.com/lugu/qiloop/bus/client/services"
	dir "github.com/lugu/qiloop/bus/server/directory"
	sess "github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/bus/util"
	"log"
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
		panic(err)
	}
	services := proxy.Services(session)
	directory, err := services.ServiceDirectory()
	if err != nil {
		log.Fatalf("failed to create directory: %s", err)
	}
	machineID, err := directory.MachineId()
	if err != nil {
		panic(err)
	}
	if machineID == "" {
		panic("empty machine id")
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
		panic(err)
	}
	info.ServiceId = uid
	err = impl.ServiceReady(uid)
	if err != nil {
		panic(err)
	}
	// shall not be able to register twice with the same name
	_, err = impl.RegisterService(info)
	if err == nil {
		panic(err)
	}
	services, err := impl.Services()
	if err != nil {
		panic(err)
	}
	if len(services) != 1 {
		t.Errorf("wrong service number: %d", len(services))
	}
	compareInfo(t, services[0], info)
	info2 := newInfo("test2")
	uid, err = impl.RegisterService(info2)
	if err != nil {
		panic(err)
	}
	info2.ServiceId = uid
	err = impl.ServiceReady(uid)
	if err != nil {
		panic(err)
	}
	// shall not be able to register twice with the same name
	_, err = impl.RegisterService(info2)
	if err == nil {
		panic(err)
	}
	services, err = impl.Services()
	if err != nil {
		panic(err)
	}
	if len(services) != 2 {
		t.Errorf("wrong service number: %d", len(services))
	}
	compareInfo(t, services[0], info)
	compareInfo(t, services[1], info2)

	err = impl.UnregisterService(info2.ServiceId + 1)
	if err == nil {
		panic("shall fail")
	}
	err = impl.UnregisterService(info.ServiceId)
	if err != nil {
		panic(err)
	}
	err = impl.UnregisterService(info.ServiceId)
	if err == nil {
		panic(err)
	}
	services, err = impl.Services()
	if err != nil {
		panic(err)
	}
	if len(services) != 1 {
		t.Errorf("wrong service number: %d", len(services))
	}
	compareInfo(t, services[0], info2)
	info2.Name = "test3"
}
