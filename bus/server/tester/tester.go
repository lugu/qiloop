package tester

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/type/object"
)

var (
	Hook = func(event string) {}
)

func NewSpacecraftObject() server.ServerObject {
	return SpacecraftObject(&spacecraftImpl{})
}

type spacecraftImpl struct {
	session   bus.Session
	terminate func()
	service   server.Service
}

func (f *spacecraftImpl) Activate(activation server.Activation,
	helper SpacecraftSignalHelper) error {
	f.session = activation.Session
	f.terminate = activation.Terminate
	f.service = activation.Service
	return nil
}

func (f *spacecraftImpl) OnTerminate() {
	Hook("SpaceCraft.OnTerminate()")
}

func (f *spacecraftImpl) Shoot() (BombProxy, error) {
	return CreateBomb(f.session, f.service, &bombImpl{})
}

func (f *spacecraftImpl) Ammo(b BombProxy) error {
	return nil
}

type bombImpl struct{}

func (f *bombImpl) Activate(activation server.Activation,
	helper BombSignalHelper) error {

	err := helper.UpdateDelay(10)
	if err != nil {
		return err
	}
	return nil
}

func (f *bombImpl) OnTerminate() {
	Hook("Bomb.OnTerminate()")

}

func (f *bombImpl) OnDelayChange(duration int32) error {
	if duration < 0 {
		return fmt.Errorf("duration cannot be negative (%d)", duration)
	}
	return nil
}

func NewBombObject() server.ServerObject {
	return BombObject(&bombImpl{})
}

// Not entirely satisfying: need to allow for client side object
// generation... Here comes the ObjectID question..
func CreateBomb(session bus.Session, service server.Service,
	impl BombImplementor) (BombProxy, error) {

	var stb stubBomb
	stb.impl = impl
	stb.obj = server.NewObject(stb.metaObject(), stb.onPropertyChange)

	objectID, err := service.Add(&stb)
	if err != nil {
		return nil, err
	}

	ref := object.ObjectReference{
		true, // with meta object
		object.FullMetaObject(stb.metaObject()),
		0,
		service.ServiceID(),
		objectID,
	}
	proxy, err := session.Object(ref)
	if err != nil {
		return nil, err
	}
	return MakeBomb(session, proxy), nil
}
