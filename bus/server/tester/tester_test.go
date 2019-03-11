package tester_test

import (
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/server/tester"
	"github.com/lugu/qiloop/bus/util"
	"testing"
)

func TestAddRemoveObject(t *testing.T) {

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

	obj := tester.NewSpacecraftObject()
	service, err := srv.NewService("Spacecraft", obj)
	if err != nil {
		t.Error(err)
	}

	session := srv.Session()
	proxies := tester.Services(session)

	spacecraft, err := proxies.Spacecraft()
	if err != nil {
		t.Fatal(err)
	}

	bomb, err := spacecraft.Shoot()
	if err != nil {
		t.Fatal(err)
	}

	err = spacecraft.Ammo(bomb)
	if err != nil {
		t.Error(err)
	}

	ammo := tester.NewBombObject()
	id, err := service.Add(ammo)
	if err != nil {
		t.Error(err)
	}

	proxy, err := session.Proxy("Spacecraft", id)
	if err != nil {
		t.Error(err)
	}

	ammoProxy := tester.MakeBomb(session, proxy)
	err = spacecraft.Ammo(ammoProxy)
	if err != nil {
		t.Error(err)
	}

	service.Terminate()
	srv.Terminate()

}
