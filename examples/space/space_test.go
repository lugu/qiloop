package space_test

import (
	"sync"
	"testing"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/examples/space"
)

func TestAddRemoveObject(t *testing.T) {

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

	obj := space.NewSpacecraftObject()
	service, err := srv.NewService("Spacecraft", obj)
	if err != nil {
		t.Error(err)
	}

	session := srv.Session()
	proxies := space.Services(session)

	spacecraft, err := proxies.Spacecraft()
	if err != nil {
		t.Fatal(err)
	}

	bomb, err := spacecraft.Shoot()
	if err != nil {
		t.Fatal(err)
	}

	// initial delay is 10 seconds
	delay, err := bomb.GetDelay()
	if err != nil {
		t.Error(err)
	} else if delay != 10 {
		t.Errorf("unexpected delay: %d", delay)
	}

	err = bomb.SetDelay(12)
	if err != nil {
		t.Error(err)
	}

	delay, err = bomb.GetDelay()
	if err != nil {
		t.Error(err)
	} else if delay != 12 {
		t.Errorf("unexpected delay: %d", delay)
	}

	err = bomb.SetDelay(-1)
	if err == nil {
		t.Error("error expected")
	}

	delay, err = bomb.GetDelay()
	if err != nil {
		t.Error(err)
	} else if delay != 12 {
		t.Errorf("unexpected delay: %d", delay)
	}

	err = spacecraft.Ammo(bomb)
	if err != nil {
		t.Error(err)
	}

	// create a new bomb and register it
	ammo := space.NewBombObject()
	id, err := service.Add(ammo)
	if err != nil {
		t.Error(err)
	}

	// get a proxy of the new bomb
	proxy, err := session.Proxy("Spacecraft", id)
	if err != nil {
		t.Error(err)
	}

	// get a specialized proxy of the new bomb
	ammoProxy := space.MakeBomb(session, proxy)

	// pass the client object to the service.
	err = spacecraft.Ammo(ammoProxy)
	if err != nil {
		t.Error(err)
	}

	service.Terminate()
	srv.Terminate()
}

func TestClientBomb(t *testing.T) {

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

	obj := space.NewSpacecraftObject()
	_, err = srv.NewService("Spacecraft", obj)
	if err != nil {
		t.Error(err)
	}

	session := srv.Session()
	defer session.Terminate()
	proxies := space.Services(session)

	spacecraft, err := proxies.Spacecraft()
	if err != nil {
		t.Fatal(err)
	}
	_, err = spacecraft.Shoot()
	if err != nil {
		t.Fatal(err)
	}

	proxyService := spacecraft.FIXMEProxy().ProxyService(session)

	bomb, err := space.CreateBomb(session, proxyService)
	if err != nil {
		t.Error(err)
	}

	err = spacecraft.Ammo(bomb)
	if err != nil {
		t.Error(err)
	}

	_, err = spacecraft.Shoot()
	if err != nil {
		t.Fatal(err)
	}

}

func TestOnTerminate(t *testing.T) {
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

	obj := space.NewSpacecraftObject()
	service, err := srv.NewService("Spacecraft", obj)
	if err != nil {
		t.Error(err)
	}
	defer service.Terminate()

	session := srv.Session()
	proxies := space.Services(session)

	spacecraft, err := proxies.Spacecraft()
	if err != nil {
		t.Fatal(err)
	}

	bomb, err := spacecraft.Shoot()
	if err != nil {
		t.Fatal(err)
	}

	var wait sync.WaitGroup
	var waiting sync.WaitGroup
	wait.Add(1)
	waiting.Add(1)
	go func() {
		for {
			waiting.Done()
			event := <-space.Hook
			if event == "Bomb.OnTerminate()" {
				wait.Done()
				return
			}
		}
	}()
	waiting.Wait()

	err = bomb.Terminate(bomb.FIXMEProxy().ObjectID())
	if err != nil {
		t.Fatal(err)
	}
	wait.Wait()
}
