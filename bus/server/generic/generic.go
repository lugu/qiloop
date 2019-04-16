package generic

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
	"log"
	"sync"
	"time"
)

// ErrWrongObjectID is returned when a method argument is given the
// wrong object ID.
var ErrWrongObjectID = errors.New("Wrong object ID")

// ErrNotYetImplemented is returned when a feature is not yet
// implemented.
var ErrNotYetImplemented = errors.New("Not supported")

func (s *stubGeneric) UpdateSignal(signal uint32, data []byte) error {
	return s.obj.UpdateSignal(signal, data)
}

func (s *stubGeneric) UpdateProperty(id uint32, sig string, data []byte) error {
	objImpl, ok := (s.impl).(*objectImpl)
	if !ok {
		return fmt.Errorf("unexpected implementation")
	}
	prop, ok := objImpl.meta.Properties[id]
	if !ok {
		return fmt.Errorf("missing property (%d), %#v", id,
			objImpl.meta)
	}
	err := objImpl.onPropertyChange(prop.Name, data)
	if err != nil {
		return err
	}
	newValue := value.Opaque(sig, data)
	err = objImpl.saveProperty(prop.Name, newValue)
	if err != nil {
		return err
	}
	return s.obj.UpdateProperty(id, sig, data)
}

func (s *stubGeneric) Wrap(id uint32, fn server.ActionWrapper) {
	s.obj.Wrap(id, fn)
}

type objectImpl struct {
	obj               *BasicObject
	meta              object.MetaObject
	onPropertyChange  func(string, []byte) error
	objectID          uint32
	signal            GenericSignalHelper
	properties        map[string]value.Value
	propertiesMutex   sync.RWMutex
	terminate         func()
	stats             map[uint32]MethodStatistics
	statsLock         sync.RWMutex
	observableWrapper server.Wrapper
	traceEnabled      bool
	traceMutex        sync.RWMutex
}

// NewObject returns an Object which implements ServerObject. It
// handles all the generic methods and signals common to all objects.
// Services implementation user this Object and fill it with the
// extra actions they wish to handle using the Wrap method. See
// type/object.Object for a list of the default methods.
// onPropertyChange is called each time a property is udpated.
func NewObject(meta object.MetaObject,
	onPropertyChange func(string, []byte) error) Object {

	impl := &objectImpl{
		meta:              object.FullMetaObject(meta),
		onPropertyChange:  onPropertyChange,
		stats:             nil,
		properties:        make(map[string]value.Value),
		observableWrapper: make(map[uint32]server.ActionWrapper),
	}
	obj := GenericObject(impl)
	stub := obj.(*stubGeneric)
	impl.basicObject(stub.obj.(*BasicObject))
	return stub
}

func (o *objectImpl) basicObject(obj *BasicObject) {
	o.obj = obj
	obj.tracer = o.tracer()
}

func (o *objectImpl) Activate(activation server.Activation,
	signal GenericSignalHelper) error {

	o.signal = signal
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
	stringValue, ok := name.(value.StringValue)
	if !ok {
		return nil, fmt.Errorf("property name must be a string value")
	}
	nameStr := stringValue.Value()
	o.propertiesMutex.RLock()
	defer o.propertiesMutex.RUnlock()
	val, ok := o.properties[nameStr]
	if !ok {
		return nil, fmt.Errorf("property unknown: %s, %#v", nameStr,
			o.properties)
	}
	return val, nil
}

func (o *objectImpl) SetProperty(name value.Value, newValue value.Value) error {
	var nameStr string
	stringValue, ok := name.(value.StringValue)
	if ok {
		nameStr = stringValue.Value()
	} else {
		idValue, ok := name.(value.UintValue)
		if !ok {
			return fmt.Errorf("incorrect name type")
		}
		property, ok := o.meta.Properties[idValue.Value()]
		if !ok {
			return fmt.Errorf(
				"incorrect property id value, got %d",
				idValue.Value())
		}
		nameStr = property.Name
	}
	var buf bytes.Buffer
	err := newValue.Write(&buf)
	if err != nil {
		return fmt.Errorf("cannot write value: %s", err)
	}
	sig, err := basic.ReadString(&buf)
	if err != nil {
		return fmt.Errorf("invalid signature: %s", err)
	}
	data := buf.Bytes()
	err = o.onPropertyChange(nameStr, data)
	if err != nil {
		return err
	}
	err = o.saveProperty(nameStr, newValue)
	if err != nil {
		return err
	}
	id, err := o.meta.PropertyID(nameStr)
	if err != nil {
		return fmt.Errorf("cannot set property: %s", err)
	}
	return o.obj.UpdateProperty(id, sig, data)
}

func (o *objectImpl) saveProperty(name string, newValue value.Value) error {
	o.propertiesMutex.Lock()
	defer o.propertiesMutex.Unlock()
	o.properties[name] = newValue
	return nil
}

func (o *objectImpl) Properties() ([]string, error) {
	properties := make([]string, 0)
	o.propertiesMutex.RLock()
	defer o.propertiesMutex.RUnlock()
	for property, _ := range o.properties {
		properties = append(properties, property)
	}
	return properties, nil
}

func (o *objectImpl) RegisterEventWithSignature(objectID uint32,
	actionID uint32, handler uint64, P3 string) (uint64, error) {
	return 0, fmt.Errorf("Not yet implemented")
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
	o.traceMutex.RLock()
	defer o.traceMutex.RUnlock()
	return o.traceEnabled, nil
}

func (o *objectImpl) EnableTrace(enable bool) error {
	o.traceMutex.RLock()
	defer o.traceMutex.RUnlock()
	o.traceEnabled = enable
	return nil
}

func (o *objectImpl) tracer() func(msg *net.Message) {
	return func(msg *net.Message) {
		o.trace(msg)
	}
}

func (o *objectImpl) trace(msg *net.Message) {
	o.traceMutex.RLock()
	enabled := o.traceEnabled
	o.traceMutex.RUnlock()
	if !enabled {
		return
	}
	// do not trace traceObject signal
	if msg.Header.Action == 86 {
		return
	}
	now := time.Now()
	timeval := Timeval{
		Tvsec:  int64(now.Second()),
		Tvusec: int64(now.Nanosecond() / 1000),
	}
	event := EventTrace{
		Id:        msg.Header.Action,
		Kind:      int32(msg.Header.Type),
		SlotId:    msg.Header.ID,
		Arguments: value.Void(),
		Timestamp: timeval,
	}
	err := o.signal.SignalTraceObject(event)
	if err != nil {
		log.Printf("trace error: %s", err)
	}
}
