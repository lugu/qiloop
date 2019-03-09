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
	ref, err := AddBomb(f.service, f.serviceID)
	if err != nil {
		panic(err)
	}
	proxy, err := f.session.Object(*ref)
	if err != nil {
		panic(err)
	}
	bomb := MakeBomb(f.session, proxy)
	return bomb, nil
}

func (f *spacecraftImpl) Ammo(b BombProxy) error {
	return nil
}

type bombImpl struct{}

func (f *bombImpl) Activate(activation server.Activation,
	helper BombSignalHelper) error {
	return nil
}

func (f *bombImpl) OnTerminate() {
}

func NewBombObject() server.ServerObject {
	return BombObject(&bombImpl{})
}

// Move into the stub
func AddBomb(service server.Service, serviceID uint32) (
	*object.ObjectReference, error) {

	var stb stubBomb
	stb.impl = &bombImpl{}
	stb.obj = generic.NewObject(stb.metaObject())

	objectID, err := service.Add(&stb)
	if err != nil {
		return nil, err
	}
	// FIXME: at this stage, stb has been activated and received
	// the service and object ids. Update the stub to save those
	// informations.

	return &object.ObjectReference{
		true, // with meta object
		object.FullMetaObject(stb.metaObject()),
		0,
		serviceID,
		objectID,
	}, nil
}
