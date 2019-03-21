package logger

import (
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/server/generic"
	"github.com/lugu/qiloop/type/object"
)

type logListenerImpl struct {
	cancel      chan struct{}
	logs        chan *LogMessage
	activation  server.Activation
	helper      LogListenerSignalHelper
	onTerminate func()
}

func CreateLogListener(session bus.Session, service server.Service,
	serviceID uint32, producer chan *LogMessage, onTerminate func()) (
	LogListenerProxy, error) {

	impl := &logListenerImpl{
		logs:        producer,
		cancel:      make(chan struct{}),
		onTerminate: onTerminate,
	}
	var stb stubLogListener
	stb.impl = impl
	stb.obj = generic.NewObject(stb.metaObject(), stb.onPropertyChange)

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
	return MakeLogListener(session, proxy), nil
}

func (l *logListenerImpl) Activate(activation server.Activation,
	helper LogListenerSignalHelper) error {

	l.helper = helper
	l.activation = activation

	if err := helper.UpdateVerbosity(LogLevelInfo); err != nil {
		return err
	}
	if err := helper.UpdateFilters(make(map[string]int32)); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-l.cancel:
				break
			case msg, ok := <-l.logs:
				if !ok {
					l.activation.Terminate()
				}
				l.helper.SignalOnLogMessage(*msg)
			}
		}
	}()
	return nil
}

func (l *logListenerImpl) OnTerminate() {
	close(l.cancel)
	l.onTerminate()
}

func (l *logListenerImpl) SetCategory(category string, level LogLevel) error {
	panic("not yet implemented")
}

func (l *logListenerImpl) ClearFilters() error {
	panic("not yet implemented")
}

func (l *logListenerImpl) OnVerbosityChange(level LogLevel) error {
	panic("not yet implemented")
}

func (l *logListenerImpl) OnFiltersChange(filters map[string]int32) error {
	panic("not yet implemented")
}
