package clock_test

import (
	"testing"
	"time"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/examples/clock"
)

func TestTimestampImplementation(t *testing.T) {
	timestamp := clock.Timestamper(time.Now())
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
		t.Errorf("non monotonic: %d, %d, %d", a, b, c)
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
	obj := clock.NewTimestampObject()
	_, err = srv.NewService("Timestamp", obj)
	if err != nil {
		t.Error(err)
	}

	// 3. client connects to the service
	session := srv.Session()
	proxies := clock.Services(session)

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

func TestSynchronizedTimestamp(t *testing.T) {

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
	obj := clock.NewTimestampObject()
	_, err = srv.NewService("Timestamp", obj)
	if err != nil {
		t.Error(err)
	}

	// 3. get a client connection to the service
	session := srv.Session()
	proxies := clock.Services(session)

	timestampService, err := proxies.Timestamp()
	if err != nil {
		t.Fatal(err)
	}

	// 4. get a synchronized timestamper
	timestamper, err := clock.SynchronizedTimestamper(session)

	nano1, _ := timestamper.Nanoseconds()
	nano2, err := timestampService.Nanoseconds()
	nano3, _ := timestamper.Nanoseconds()

	if err != nil {
		t.Errorf("reference timestamp: %s", err)
	}

	if nano2 <= nano1 || nano3 <= nano2 {
		t.Errorf("Synchronization failed: %d, %d, %d",
			nano1, nano2, nano3)
	}
}
