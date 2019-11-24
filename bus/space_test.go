package bus_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/examples/space"
)

var (
	// Hook let the test program record what is happening
	Hook = make(chan string)
)

// NewSpacecraftObject creates a new server side Spacecraft object.
func NewSpacecraftObject() bus.Actor {
	return SpacecraftObject(&spacecraftImpl{})
}

// spacecraftImpl implements SpacecraftImplementor.
type spacecraftImpl struct {
	session   bus.Session
	terminate func()
	service   bus.Service
	ammo      BombProxy
}

func (f *spacecraftImpl) Activate(activation bus.Activation,
	helper SpacecraftSignalHelper) error {
	f.session = activation.Session
	f.terminate = activation.Terminate
	f.service = activation.Service
	ammo, err := CreateBomb(f.session, f.service)
	f.ammo = ammo
	return err
}

func (f *spacecraftImpl) OnTerminate() {
	select {
	case Hook <- "SpaceCraft.OnTerminate()":
	default:
	}
}

func (f *spacecraftImpl) Shoot() (BombProxy, error) {
	return f.ammo, nil
}

func (f *spacecraftImpl) Ammo(b BombProxy) error {
	f.ammo = b
	return nil
}

// bombImpl implements BombImplementor.
type bombImpl struct{}

func (f *bombImpl) Activate(activation bus.Activation,
	helper BombSignalHelper) error {

	err := helper.UpdateDelay(10)
	if err != nil {
		return err
	}
	return nil
}

func (f *bombImpl) OnTerminate() {
	select {
	case Hook <- "Bomb.OnTerminate()":
	default:
	}

}

func (f *bombImpl) OnDelayChange(duration int32) error {
	if duration < 0 {
		return fmt.Errorf("duration cannot be negative (%d)", duration)
	}
	return nil
}

// NewBombObject returns the server side implementation of a Bomb
// object.
func NewBombObject() bus.Actor {
	return BombObject(&bombImpl{})
}

// CreateBomb creates a new bomb and add it the to service.
func CreateBomb(session bus.Session, service bus.Service) (BombProxy, error) {
	constructor := Services(session)
	return constructor.NewBomb(service, &bombImpl{})
}
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

	spacecraft, err := proxies.Spacecraft(nil)
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

	props, err := ammoProxy.Properties()
	if err != nil {
		t.Error(err)
	}
	if len(props) != 1 || props[0] != "delay" {
		t.Error(err)
	}

	cancel, tik, err := ammoProxy.SubscribeDelay()
	if err != nil {
		t.Error(err)
	}
	defer cancel()

	// pass the client object to the service.
	err = spacecraft.Ammo(ammoProxy)
	if err != nil {
		t.Error(err)
	}

	err = ammoProxy.SetDelay(10)
	if err != nil {
		t.Error(err)
	}

	timer := time.NewTimer(time.Second)
	select {
	case <-timer.C:
		t.Error("missing delay")
	case <-tik:
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

	spacecraft, err := proxies.Spacecraft(nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = spacecraft.Shoot()
	if err != nil {
		t.Fatal(err)
	}

	proxyService := spacecraft.ProxyService(session)

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

	spacecraft, err := proxies.Spacecraft(nil)
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

	err = bomb.Terminate(bomb.ObjectID())
	if err != nil {
		t.Fatal(err)
	}
	wait.Wait()
}
