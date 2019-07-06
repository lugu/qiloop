package logger

import (
	"fmt"
	"sync"
	"time"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
)

type logProvider struct {
	manager        LogManager
	location       string
	source         string
	category       string
	logger         Logger
	verbosity      LogLevel
	verbosityMutex sync.RWMutex
	filters        map[string]LogLevel
	filtersMutex   sync.RWMutex
}

type Logger interface {
	Log(level LogLevel, message string)
}

func NewLogMessage(level LogLevel, location, source, category,
	msg string) LogMessage {
	now := time.Now().UnixNano()

	return LogMessage{
		Source:     source,
		Level:      level,
		Category:   category,
		Location:   location,
		Message:    msg,
		Id:         0,
		Date:       TimePoint{uint64(now)},
		SystemDate: TimePoint{uint64(now)},
	}
}

func NewLogProviderImpl(source, category string) (LogProviderImplementor, Logger) {
	location := fmt.Sprintf("%s:%d", util.MachineID(), util.ProcessID())
	impl := &logProvider{
		source:    source,
		location:  location,
		category:  category,
		verbosity: LogLevelInfo,
		filters:   make(map[string]LogLevel),
	}
	return impl, impl
}

func (l *logProvider) Log(level LogLevel, msg string) {
	filterLevel, ok := l.filters[l.category]
	if ok && level.Level > filterLevel.Level {
		return
	}
	if level.Level > l.verbosity.Level {
		return
	}
	logMsg := NewLogMessage(level, l.location, l.source, l.category, msg)
	l.manager.Log([]LogMessage{logMsg})
}

func (l *logProvider) Activate(activation bus.Activation,
	helper LogProviderSignalHelper) (err error) {
	services := Services(activation.Session)
	l.manager, err = services.LogManager()
	if err != nil {
		return fmt.Errorf("Cannot create LogProvider: %s", err)
	}
	return nil
}
func (l *logProvider) OnTerminate() {
}
func (l *logProvider) SetVerbosity(level LogLevel) error {
	panic("not yet implemented")
}
func (l *logProvider) SetCategory(category string, level LogLevel) error {
	panic("not yet implemented")
}
func (l *logProvider) ClearAndSet(filters map[string]LogLevel) error {
	panic("not yet implemented")
}

func CreateLogProvider(session bus.Session, service bus.Service,
	impl LogProviderImplementor) (
	LogProviderProxy, error) {

	var stb stubLogProvider
	stb.impl = impl
	obj := bus.NewBasicObject(&stb, stb.metaObject(), stb.onPropertyChange)
	stb.signal = obj

	objectID, err := service.Add(obj)
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
	return MakeLogProvider(session, proxy), nil
}
