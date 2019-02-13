package generic

import (
	"errors"
	"fmt"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
)

// ErrWrongObjectID is returned when a method argument is given the
// wrong object ID.
var ErrWrongObjectID = errors.New("Wrong object ID")

func (s *stubGeneric) UpdateSignal(signal uint32, value []byte) error {
	return s.obj.UpdateSignal(signal, value)
}

func (s *stubGeneric) Wrap(id uint32, fn server.ActionWrapper) {
	s.obj.Wrap(id, fn)
}

type objectImpl struct {
	obj       *BasicObject
	meta      object.MetaObject
	serviceID uint32
	objectID  uint32
	signal    GenericSignalHelper
	terminate server.Terminator
}

// NewObject returns an Object which implements ServerObject. It
// handles all the generic methods and signals common to all objects.
// Services implementation user this Object and fill it with the
// extra actions they wish to handle using the Wrap method. See
// type/object.Object for a list of the default methods.
func NewObject(meta object.MetaObject) Object {
	obj := GenericObject(genericObject(meta))
	stub := obj.(*stubGeneric)
	return stub
}

func genericObject(meta object.MetaObject) Generic {
	impl := objectImpl{
		obj:  NewBasicObject(),
		meta: object.FullMetaObject(meta),
	}
	return &impl
}

func (o *objectImpl) Activate(activation server.Activation,
	signal GenericSignalHelper) error {

	o.signal = signal
	o.serviceID = activation.ServiceID
	o.objectID = activation.ObjectID
	o.terminate = activation.Terminate
	return nil
}

func (o *objectImpl) OnTerminate() {
	o.obj.OnTerminate()
}

func (o *objectImpl) RegisterEvent(objectID uint32, actionID uint32,
	handler uint64) (uint64, error) {
	return 0, fmt.Errorf("Not implemented: see BasicObject")
}

func (o *objectImpl) UnregisterEvent(objectID uint32, actionID uint32,
	handler uint64) error {
	return fmt.Errorf("Not implemented: see BasicObject")
}

func (o *objectImpl) MetaObject(objectID uint32) (object.MetaObject, error) {
	if objectID != o.objectID {
		return o.meta, ErrWrongObjectID
	}
	return o.meta, nil
}

func (o *objectImpl) Terminate(objectID uint32) error {
	if objectID != o.objectID {
		return ErrWrongObjectID
	}
	o.terminate()
	return nil
}

func (o *objectImpl) Property(name value.Value) (value.Value, error) {
	panic("Not yet implemented")
}

func (o *objectImpl) SetProperty(name value.Value, value value.Value) error {
	panic("Not yet implemented")
}

func (o *objectImpl) Properties() ([]string, error) {
	panic("Not yet implemented")
}

func (o *objectImpl) RegisterEventWithSignature(objectID uint32,
	actionID uint32, handler uint64, P3 string) (uint64, error) {
	panic("Not yet implemented")
}

func (o *objectImpl) IsStatsEnabled() (bool, error) {
	panic("Not yet implemented")
}

func (o *objectImpl) EnableStats(enabled bool) error {
	panic("Not yet implemented")
}

func (o *objectImpl) Stats() (map[uint32]MethodStatistics, error) {
	panic("Not yet implemented")
}

func (o *objectImpl) ClearStats() error {
	panic("Not yet implemented")
}

func (o *objectImpl) IsTraceEnabled() (bool, error) {
	panic("Not yet implemented")
}

func (o *objectImpl) EnableTrace(traced bool) error {
	panic("Not yet implemented")
}
