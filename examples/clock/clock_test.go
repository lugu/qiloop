package clock

import (
	"testing"
	"time"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
)

func TestTimestampImplementation(t *testing.T) {
	var timestamp timestampService
	beforeA := time.Now()
	a, err := timestamp.Nanoseconds()
	if err != nil {
		t.Fatal(err)
	}
	afterA := time.Now()
	b, err := timestamp.Nanoseconds()
	if err != nil {
		t.Fatal(err)
	}
	afterB := time.Now()
	c, err := timestamp.Nanoseconds()
	if err != nil {
		t.Fatal(err)
	}
	if a >= b || b >= c {
		t.Error("non monotonic")
	}
	if b-a >= afterB.Sub(beforeA).Nanoseconds() {
		t.Error("too fast")
	}
	if c-a <= afterB.Sub(afterA).Nanoseconds() {
		t.Error("too slow")
	}
}

func TestTimestampService(t *testing.T) {

	// 1. create a server which listen to a socket file.
	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	ns := bus.PrivateNamespace()
	srv, err := bus.StandAloneServer(listener, bus.Yes{}, ns)
	if err != nil {
		t.Error(err)
	}
	defer srv.Terminate()

	// 2. register the service
	obj := NewTimestampObject()
	_, err = srv.NewService("Timestamp", obj)
	if err != nil {
		t.Error(err)
	}

	// 3. client connects to the service
	session := srv.Session()
	proxies := Services(session)

	timestamp, err := proxies.Timestamp()
	if err != nil {
		t.Fatal(err)
	}

	// 4. client request a timestamp
	nano, err := timestamp.Nanoseconds()
	if err != nil {
		t.Fatal(err)
	}
	if nano == 0 {
		t.Fatal("zero timestamp")
	}
}
