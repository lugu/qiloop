package bus_test

import (
	gonet "net"
	"testing"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/directory"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
)

func TestCache(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	cache := bus.NewCache(endpoint)
	defer cache.Terminate()
	cache.AddService("Server", 0, object.MetaService0)

	err = cache.Lookup("ServiceDirectory", 1)
	if err == nil {
		panic("expecting an authentication error")
	}
}

func TestServerProxy(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	cache := bus.NewCache(endpoint)
	defer cache.Terminate()

	_, err = bus.ServiceServer(cache)
	if err == nil {
		panic("shall not create service not in the cache")
	}

	cache.AddService("ServiceZero", 0, object.MetaService0)

	_, err = bus.ServiceServer(cache)
	if err != nil {
		panic("expecting an authentication error")
	}

	_, err = bus.ServiceZero(cache)
	if err != nil {
		panic("expecting an authentication error")
	}
}

func TestLookup(t *testing.T) {
	addr := util.NewUnixAddr()
	server, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	cache, err := bus.NewCachedSession(addr)
	if err != nil {
		panic(err)
	}
	defer cache.Terminate()
	cache.AddService("Server", 0, object.MetaService0)

	_, err = services.ServiceDirectory(cache)
	if err == nil {
		panic("ServiceDirectory not yet cached")
	}

	err = cache.Lookup("ServiceDirectory", 1)
	if err != nil {
		panic(err)
	}

	directory, err := services.ServiceDirectory(cache)
	if err != nil {
		panic(err)
	}
	list, err := directory.Services()
	if err != nil {
		panic(err)
	}
	if len(list) == 0 {
		panic("service list empty")
	}
}

func TestDialError(t *testing.T) {
	_, err := bus.NewCachedSession("nope://")
	if err == nil {
		panic("expecting an error")
	}
}
func TestCacheAuthError(t *testing.T) {
	addr := util.NewUnixAddr()

	listener, err := net.Listen(addr)
	if err != nil {
		panic(err)
	}

	server, err := bus.StandAloneServer(listener, bus.No{},
		bus.PrivateNamespace())
	if err != nil {
		panic(err)
	}
	defer server.Terminate()
	_, err = bus.NewCachedSession(addr)
	if err == nil {
		panic("expecting an error")
	}
}

func TestServerError(t *testing.T) {
	clt, srv := gonet.Pipe()

	buf := make([]byte, 10)

	go func() {
		_, err := srv.Read(buf)
		if err != nil {
			panic(err)
		}
		srv.Close()
	}()

	cache := bus.NewCache(net.ConnEndPoint(clt))
	defer cache.Terminate()
	cache.AddService("ServiceZero", 0, object.MetaService0)

	server0, err := bus.ServiceZero(cache)
	if err != nil {
		panic("expecting an authentication error")
	}

	capabilityMap := map[string]value.Value{
		bus.KeyUser:  value.String("a"),
		bus.KeyToken: value.String("b"),
	}
	_, err = server0.Authenticate(capabilityMap)
	if err == nil {
		panic("shall not accept ...")
	}
}
