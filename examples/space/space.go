// Package space contains the implementation of the Spacecraft and Bomb
// objects.
//
// The file space_stub_gen.go is generated with:
//
// 	$ qiloop stub --idl space.qi.idl --output space_stub_gen.go
//
package space

import (
	"fmt"

	"github.com/lugu/qiloop/bus"
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
	ammo, err := NewBomb(f.session, f.service)
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

// NewBomb creates a new bomb and add it the to service.
func NewBomb(session bus.Session, service bus.Service) (BombProxy, error) {
	return CreateBomb(session, service, &bombImpl{})
}
