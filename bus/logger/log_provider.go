package logger

import (
	"fmt"
	"os"
	"runtime"
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
	id             int32
	location       string
	category       string
	logger         Logger
	verbosity      LogLevel
	verbosityMutex sync.RWMutex
	filters        map[string]LogLevel
	filtersMutex   sync.RWMutex
}

// Logger sends a log message to the log manager if a log listener
// request it, else it does nothing, except Fatal which calls
// os.Exit(1).
type Logger interface {
	Fatal(format string, v ...interface{})
	Error(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Info(format string, v ...interface{})
	Verbose(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Terminate()
}

func newLogMessage(level LogLevel, location, category,
	msg string) LogMessage {
	now := time.Now().UnixNano()

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	return LogMessage{
		Source:     fmt.Sprintf("%s:%d", file, line),
		Level:      level,
		Category:   category,
		Location:   location,
		Message:    msg,
		Id:         0,
		Date:       TimePoint{uint64(now)},
		SystemDate: TimePoint{uint64(now)},
	}
}

func newLogProviderImpl(category string) (LogProviderImplementor, *logProvider) {
	location := fmt.Sprintf("%s:%d", util.MachineID(), util.ProcessID())
	impl := &logProvider{
		location:  location,
		category:  category,
		verbosity: LogLevelInfo,
		filters:   make(map[string]LogLevel),
	}
	return impl, impl
}

func (l *logProvider) logf(level LogLevel, format string, v ...interface{}) {
	filterLevel, ok := l.filters[l.category]
	if ok && level.Level > filterLevel.Level {
		return
	}
	l.verbosityMutex.RLock()
	if level.Level > l.verbosity.Level {
		l.verbosityMutex.RUnlock()
		return
	}
	l.verbosityMutex.RUnlock()

	msg := fmt.Sprintf(format, v...)
	logMsg := newLogMessage(level, l.location, l.category, msg)
	l.manager.Log([]LogMessage{logMsg})
}

func (l *logProvider) Fatal(format string, v ...interface{}) {
	l.logf(LogLevelFatal, format, v...)
	os.Exit(1)
}

func (l *logProvider) Error(format string, v ...interface{}) {
	l.logf(LogLevelError, format, v...)
}

func (l *logProvider) Warning(format string, v ...interface{}) {
	l.logf(LogLevelWarning, format, v...)
}

func (l *logProvider) Info(format string, v ...interface{}) {
	l.logf(LogLevelInfo, format, v...)
}

func (l *logProvider) Verbose(format string, v ...interface{}) {
	l.logf(LogLevelVerbose, format, v...)
}

func (l *logProvider) Debug(format string, v ...interface{}) {
	l.logf(LogLevelDebug, format, v...)
}

func (l *logProvider) Terminate() {
	l.manager.RemoveProvider(l.id)
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
	l.SetVerbosity(LogLevelNone)
}
func (l *logProvider) SetVerbosity(level LogLevel) error {
	l.verbosityMutex.Lock()
	l.verbosity = level
	l.verbosityMutex.Unlock()
	return nil
}
func (l *logProvider) SetCategory(category string, level LogLevel) error {
	if category == l.category {
		l.SetVerbosity(level)
	}
	return nil
}
func (l *logProvider) ClearAndSet(filters map[string]LogLevel) error {
	for filter, level := range filters {
		l.SetCategory(filter, level)
	}
	return nil
}

// NewLogger returns a new Logger. NewLogger creates a new log
// provider and associate it with log manager.
func NewLogger(session bus.Session, category string) (Logger, error) {
	constructor := Services(session)
	logManager, err := constructor.LogManager()
	if err != nil {
		return nil, err
	}
	service := logManager.FIXMEProxy().ProxyService(session)
	impl, logger := newLogProviderImpl(category)
	proxy, err := constructor.NewLogProvider(service, impl)
	logger.id, err = logManager.AddProvider(proxy)
	if err != nil {
		return nil, err
	}
	return logger, err
}
