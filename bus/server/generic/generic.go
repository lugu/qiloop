package generic

import (
	"errors"
	"fmt"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
	"sync"
	"time"
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
	obj               *BasicObject
	meta              object.MetaObject
	serviceID         uint32
	objectID          uint32
	signal            GenericSignalHelper
	terminate         server.Terminator
	stats             map[uint32]MethodStatistics
	statsLock         sync.RWMutex
	observableWrapper server.Wrapper
}

// NewObject returns an Object which implements ServerObject. It
// handles all the generic methods and signals common to all objects.
// Services implementation user this Object and fill it with the
// extra actions they wish to handle using the Wrap method. See
// type/object.Object for a list of the default methods.
func NewObject(meta object.MetaObject) Object {
	impl := &objectImpl{
		meta:              object.FullMetaObject(meta),
		stats:             nil,
		observableWrapper: make(map[uint32]server.ActionWrapper),
	}
	obj := GenericObject(impl)
	stub := obj.(*stubGeneric)
	impl.obj = stub.obj.(*BasicObject)
	return stub
}

func (o *objectImpl) Activate(activation server.Activation,
	signal GenericSignalHelper) error {

	o.signal = signal
	o.serviceID = activation.ServiceID
	o.objectID = activation.ObjectID
	o.terminate = activation.Terminate

	// During activation, all the action wrapper have been
	// registered. Dupplicate them to enable method statistics.
	for id, fn := range o.obj.wrapper {
		o.observableWrapper[id] = o.observer(id, fn)
	}
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
	o.statsLock.RLock()
	defer o.statsLock.RUnlock()
	return o.stats != nil, nil
}

func (m MethodStatistics) updateWith(t time.Duration) MethodStatistics {
	m.Count++
	duration := float32(t)
	m.Wall.CumulatedValue += duration
	if duration < m.Wall.MinValue {
		m.Wall.MinValue = duration
	}
	if duration > m.Wall.MaxValue {
		m.Wall.MaxValue = duration
	}
	return m
}

// observer returns an ActionWrapper based on fn which records statistics
func (o *objectImpl) observer(id uint32, fn server.ActionWrapper) server.ActionWrapper {
	return func(payload []byte) ([]byte, error) {
		start := time.Now()
		ret, err := fn(payload)
		duration := time.Since(start)
		o.statsLock.Lock()
		defer o.statsLock.Unlock()
		if o.stats != nil {
			o.stats[id] = o.stats[id].updateWith(duration)
		}
		return ret, err
	}
}

func (o *objectImpl) swipeWrapper() {
	o.observableWrapper, o.obj.wrapper = o.obj.wrapper, o.observableWrapper
}

func (o *objectImpl) EnableStats(enabled bool) error {
	o.statsLock.Lock()
	defer o.statsLock.Unlock()
	if enabled && o.stats == nil {
		o.stats = make(map[uint32]MethodStatistics)
		o.swipeWrapper()
	} else if !enabled && o.stats != nil {
		o.stats = nil
		o.swipeWrapper()
	}
	return nil
}

func (o *objectImpl) Stats() (map[uint32]MethodStatistics, error) {
	stats := make(map[uint32]MethodStatistics)
	o.statsLock.RLock()
	defer o.statsLock.RUnlock()
	if o.stats != nil {
		for id, stat := range o.stats {
			stats[id] = stat
		}
	}
	return stats, nil
}

func (o *objectImpl) ClearStats() error {
	o.statsLock.Lock()
	defer o.statsLock.Unlock()
	if o.stats != nil {
		o.stats = make(map[uint32]MethodStatistics)
	}
	return nil
}

func (o *objectImpl) IsTraceEnabled() (bool, error) {
	panic("Not yet implemented")
}

func (o *objectImpl) EnableTrace(traced bool) error {
	panic("Not yet implemented")
}
