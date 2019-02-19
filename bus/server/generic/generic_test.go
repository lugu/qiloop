package generic

import (
	proxy "github.com/lugu/qiloop/bus/client/object"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
	"testing"
)

func TestBasicObjectWrap(t *testing.T) {
	obj := NewBasicObject()
	passed := false

	obj.Wrap(123, func(payload []byte) ([]byte, error) {
		passed = true
		return nil, nil
	})
	in, out := net.Pipe()

	channel, err := out.ReceiveAny()
	if err != nil {
		t.Error(err)
	}

	ctx := server.NewContext(in)
	hdr := net.NewHeader(net.Call, 0, 0, 123, 0)
	msg := net.NewMessage(hdr, nil)

	err = obj.Receive(&msg, ctx)
	if err != nil {
		t.Error(err)
	}
	if passed == false {
		t.Errorf("failed to pass")
	}
	reply := <-channel
	if reply.Header.Type != net.Reply {
		t.Errorf("type is %d", reply.Header.Type)
	}

	channel, err = out.ReceiveAny()
	if err != nil {
		t.Error(err)
	}

	hdr = net.NewHeader(net.Call, 0, 0, 124, 0)
	msg = net.NewMessage(hdr, nil)

	err = obj.Receive(&msg, ctx)
	if err != nil {
		t.Error(err)
	}
	reply = <-channel
	if reply.Header.Type != net.Error {
		t.Errorf("type is %d", reply.Header.Type)
	}
}

func newObject() Object {
	return NewObject(object.MetaObject{
		Description: "",
		Methods:     make(map[uint32]object.MetaMethod),
		Signals:     make(map[uint32]object.MetaSignal),
		Properties:  make(map[uint32]object.MetaProperty),
	})
}

func TestMethodStatistics(t *testing.T) {
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
	srv.NewService("serviceA", newObject())

	sess := srv.Session()
	client, err := sess.Proxy("serviceA", 1)
	if err != nil {
		t.Error(err)
	}
	remoteObj := &proxy.ObjectProxy{client}
	enabled, err := remoteObj.IsStatsEnabled()
	if err != nil {
		t.Error(err)
	}
	if enabled {
		t.Errorf("Stats shall not be enabled")
	}
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
	srv, err := server.StandAloneServer(listener, server.Yes{},
		server.PrivateNamespace())
	if err != nil {
		t.Error(err)
	}
	srv.NewService("serviceA", newObject())

	sess := srv.Session()
	client, err := sess.Proxy("serviceA", 1)
	if err != nil {
		t.Error(err)
	}
	remoteObj := &proxy.ObjectProxy{client}
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
