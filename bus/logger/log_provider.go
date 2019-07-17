package logger

import (
	"fmt"
	"sync"
	"time"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/util"
)

// ErrNotImplemented is returned when a method is not supported by
// this implementation.
var ErrNotImplemented = fmt.Errorf("Not supported")

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
	return ErrNotImplemented
}
func (l *logProvider) SetCategory(category string, level LogLevel) error {
	return ErrNotImplemented
}
func (l *logProvider) ClearAndSet(filters map[string]LogLevel) error {
	return ErrNotImplemented
}

// CreateLogProvider create a LogProvider object and add it to the
// service. It returns a proxy of the provider.
func CreateLogProvider(session bus.Session, service bus.Service,
	source, category string) (
	LogProviderProxy, Logger, error) {

	impl, logger := NewLogProviderImpl(source, category)
	proxy, err := Services(session).NewLogProvider(service, impl)
	return proxy, logger, err
}
