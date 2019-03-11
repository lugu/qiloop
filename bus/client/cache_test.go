package client_test

import (
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/client/services"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/server/directory"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
	gonet "net"
	"strings"
	"testing"
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
	cache := client.NewCache(endpoint)
	defer cache.Destroy()
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
	cache := client.NewCache(endpoint)
	defer cache.Destroy()

	_, err = client.Services(cache).ServiceServer()
	if err == nil {
		panic("shall not create service not in the cache")
	}

	cache.AddService("Server", 0, object.MetaService0)

	_, err = client.Services(cache).ServiceServer()
	if err != nil {
		panic("expecting an authentication error")
	}

	services := client.Services(cache)
	_, err = services.Server()
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

	cache, err := client.NewCachedSession(addr)
	if err != nil {
		panic(err)
	}
	defer cache.Destroy()
	cache.AddService("Server", 0, object.MetaService0)

	services := services.Services(cache)
	_, err = services.ServiceDirectory()
	if err == nil {
		panic("ServiceDirectory not yet cached")
	}

	err = cache.Lookup("ServiceDirectory", 1)
	if err != nil {
		panic(err)
	}

	directory, err := services.ServiceDirectory()
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
	_, err := client.NewCachedSession("nope://")
	if err == nil {
		panic("expecting an error")
	}
}
func TestCacheAuthError(t *testing.T) {
	addr := util.NewUnixAddr()

	listener, err := gonet.Listen("unix", strings.TrimPrefix(addr,
		"unix://"))
	if err != nil {
		panic(err)
	}

	server, err := server.StandAloneServer(listener, server.No{},
		server.PrivateNamespace())
	if err != nil {
		panic(err)
	}
	defer server.Terminate()
	_, err = client.NewCachedSession(addr)
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

	cache := client.NewCache(net.NewEndPoint(clt))
	defer cache.Destroy()
	cache.AddService("Server", 0, object.MetaService0)

	services := client.Services(cache)
	server0, err := services.Server()
	if err != nil {
		panic("expecting an authentication error")
	}

	capabilityMap := map[string]value.Value{
		client.KeyUser:  value.String("a"),
		client.KeyToken: value.String("b"),
	}
	_, err = server0.Authenticate(capabilityMap)
	if err == nil {
		panic("shall not accept ...")
	}
}
