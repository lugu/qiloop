package client_test

import (
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/client/services"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/server/directory"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
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

func TestLookup(t *testing.T) {
	addr := util.NewUnixAddr()
	server, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	cache, err := client.NewCachedSession(addr)
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
	defer server.Terminate()
	_, err = client.NewCachedSession(addr)
	if err == nil {
		panic("expecting an error")
	}
}
