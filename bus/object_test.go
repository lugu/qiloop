package bus

import (
	"testing"

	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
)

func newObject() BasicObject {
	return NewObject(object.MetaObject{
		Description: "",
		Methods:     make(map[uint32]object.MetaMethod),
		Signals:     make(map[uint32]object.MetaSignal),
		Properties:  make(map[uint32]object.MetaProperty),
	}, func(string, []byte) error { return nil })
}

func TestMethodStatistics(t *testing.T) {
	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	srv, err := StandAloneServer(listener, Yes{},
		PrivateNamespace())
	if err != nil {
		t.Error(err)
	}
	srv.NewService("serviceA", newObject())

	sess := srv.Session()
	proxy, err := sess.Proxy("serviceA", 1)
	if err != nil {
		t.Error(err)
	}
	remoteObj := MakeObject(proxy)
	enabled, err := remoteObj.IsStatsEnabled()
	if err != nil {
		t.Error(err)
	}
	if enabled {
		t.Errorf("Stats shall not be enabled")
	}
	// FIXME
	t.Skip("Statistics not working")
	err = remoteObj.EnableStats(true)
	if err != nil {
		t.Error(err)
	}
	enabled, err = remoteObj.IsStatsEnabled()
	if err != nil {
		t.Error(err)
	}
	if !enabled {
		t.Errorf("Stats shall be enabled")
	}
	stats, err := remoteObj.Stats()
	if err != nil {
		t.Error(err)
	}
	actionID := uint32(0x50) // isStatsEnabled
	methodStats, ok := stats[actionID]
	if !ok {
		t.Fatalf("missing action from stats")
	}
	if methodStats.Count != 1 {
		t.Errorf("cound not 1 (%d)", methodStats.Count)
	}
	if methodStats.Wall.CumulatedValue == 0 {
		t.Errorf("null cumulative value")
	}
	err = remoteObj.EnableStats(false)
	if err != nil {
		t.Error(err)
	}
	enabled, err = remoteObj.IsStatsEnabled()
	if err != nil {
		t.Error(err)
	}
	if enabled {
		t.Errorf("Stats shall not be enabled")
	}
}

func TestTraceEvent(t *testing.T) {
	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	srv, err := StandAloneServer(listener, Yes{},
		PrivateNamespace())
	if err != nil {
		t.Error(err)
	}
	srv.NewService("serviceA", newObject())

	sess := srv.Session()
	proxy, err := sess.Proxy("serviceA", 1)
	if err != nil {
		t.Error(err)
	}
	remoteObj := MakeObject(proxy)
	enabled, err := remoteObj.IsTraceEnabled()
	if err != nil {
		t.Fatal(err)
	}
	if enabled {
		t.Errorf("Trace shall not be enabled")
	}
	cancel, traces, err := remoteObj.SubscribeTraceObject()
	if err != nil {
		t.Fatal(err)
	}
	err = remoteObj.EnableTrace(true)
	if err != nil {
		t.Fatal(err)
	}
	enabled, err = remoteObj.IsTraceEnabled()
	if err != nil {
		t.Fatal(err)
	}
	if !enabled {
		t.Errorf("Trace shall be enabled")
	}
	// FIXME
	t.Skip("Traces not working")
	trace := <-traces
	if trace.Id < 84 || trace.Id > 86 { // tracing actions
		t.Errorf("unexpected action %#v", trace)
	}
	cancel()
	err = remoteObj.EnableTrace(false)
	if err != nil {
		t.Error(err)
	}
	enabled, err = remoteObj.IsTraceEnabled()
	if err != nil {
		t.Error(err)
	}
	if enabled {
		t.Errorf("Trace shall not be enabled")
	}
}
