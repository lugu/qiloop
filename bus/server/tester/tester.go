package tester

import (
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/server/generic"
	"github.com/lugu/qiloop/type/object"
)

func NewSpacecraftObject() server.ServerObject {
	return SpacecraftObject(&spacecraftImpl{})
}

type spacecraftImpl struct {
	session   bus.Session
	terminate server.Terminator
	service   server.Service
	serviceID uint32
}

func (f *spacecraftImpl) Activate(activation server.Activation,
	helper SpacecraftSignalHelper) error {
	f.session = activation.Session
	f.terminate = activation.Terminate
	f.service = activation.Service
	f.serviceID = activation.ServiceID
	return nil
}

func (f *spacecraftImpl) OnTerminate() {
}

func (f *spacecraftImpl) Shoot() (BombProxy, error) {
	return CreateBomb(f.session, f.service, f.serviceID, &bombImpl{})
}

func (f *spacecraftImpl) Ammo(b BombProxy) error {
	return nil
}

type bombImpl struct{}

func (f *bombImpl) Activate(activation server.Activation,
	helper BombSignalHelper) error {

	helper.UpdateDelay(10)
	return nil
}

func (f *bombImpl) OnTerminate() {
}

func NewBombObject() server.ServerObject {
	return BombObject(&bombImpl{})
}

// Not entirely satisfying: need to allow for client side object
// generation... Here comes the ObjectID question..
func CreateBomb(session bus.Session, service server.Service, serviceID uint32,
	impl BombImplementor) (BombProxy, error) {

	var stb stubBomb
	stb.impl = impl
	stb.obj = generic.NewObject(stb.metaObject())

	objectID, err := service.Add(&stb)
	if err != nil {
		return nil, err
	}

	ref := object.ObjectReference{
		true, // with meta object
		object.FullMetaObject(stb.metaObject()),
		0,
		serviceID,
		objectID,
	}
	proxy, err := session.Object(ref)
	if err != nil {
		return nil, err
	}
	return MakeBomb(session, proxy), nil
}
