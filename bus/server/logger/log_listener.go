package logger

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/server/generic"
	"github.com/lugu/qiloop/type/object"
	"regexp"
	"sync"
)

type logListenerImpl struct {
	filters      map[string]int32
	filtersReg   map[string]*regexp.Regexp
	level        int32
	filtersMutex sync.RWMutex

	cancel      chan struct{}
	logs        chan []LogMessage
	activation  server.Activation
	helper      LogListenerSignalHelper
	onTerminate func()
}

func CreateLogListener(session bus.Session, service server.Service,
	producer chan []LogMessage, onTerminate func()) (
	LogListenerProxy, error) {

	impl := &logListenerImpl{
		filters:     make(map[string]int32),
		filtersReg:  make(map[string]*regexp.Regexp),
		level:       LogLevelInfo.Level,
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
		service.ServiceID(),
		objectID,
	}
	proxy, err := session.Object(ref)
	if err != nil {
		return nil, err
	}
	return MakeLogListener(session, proxy), nil
}

func (l *logListenerImpl) filter(msg *LogMessage) bool {
	if msg.Level.Level == LogLevelNone.Level {
		return false
	}
	l.filtersMutex.RLock()
	defer l.filtersMutex.RUnlock()

	for pattern, reg := range l.filtersReg {
		if reg.FindString(msg.Category) != "" {
			if msg.Level.Level <= l.filters[pattern] {
				return true
			}
		}
	}
	if msg.Level.Level <= l.level {
		return true
	}
	return false
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
				return
			case messages, ok := <-l.logs:
				if !ok {
					l.activation.Terminate()
					return
				}
				for _, msg := range messages {
					if l.filter(&msg) {
						l.helper.SignalOnLogMessage(msg)
					}
				}
			}
		}
	}()
	return nil
}

func (l *logListenerImpl) OnTerminate() {
	close(l.cancel)
	l.onTerminate()
}

func validateLevel(level LogLevel) error {
	if level.Level < LogLevelNone.Level ||
		level.Level > LogLevelDebug.Level {
		return fmt.Errorf("invalid level (%d)", level.Level)
	}
	return nil
}

func (l *logListenerImpl) SetCategory(category string, level LogLevel) error {
	if err := validateLevel(level); err != nil {
		return err
	}
	reg, err := regexp.Compile(category)
	if err != nil {
		return fmt.Errorf("invalid regexp (%s): %s", category, err)
	}
	l.filtersMutex.Lock()
	defer l.filtersMutex.Unlock()
	l.filters[category] = level.Level
	l.filtersReg[category] = reg
	return nil
}

func (l *logListenerImpl) ClearFilters() error {
	l.filtersMutex.Lock()
	defer l.filtersMutex.Unlock()
	l.filters = make(map[string]int32)
	l.filtersReg = make(map[string]*regexp.Regexp)
	return nil
}

func (l *logListenerImpl) OnVerbosityChange(level LogLevel) error {
	if err := validateLevel(level); err != nil {
		return err
	}
	l.filtersMutex.Lock()
	defer l.filtersMutex.Unlock()
	l.level = level.Level
	return nil
}

func (l *logListenerImpl) OnFiltersChange(filters map[string]int32) error {

	newFilters := make(map[string]int32)
	newFiltersReg := make(map[string]*regexp.Regexp)
	for pattern, level := range filters {
		if err := validateLevel(LogLevel{level}); err != nil {
			return err
		}
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regexp (%s): %s", pattern, err)
		}
		newFilters[pattern] = level
		newFiltersReg[pattern] = reg
	}

	l.filtersMutex.Lock()
	defer l.filtersMutex.Unlock()
	l.filters = newFilters
	l.filtersReg = newFiltersReg
	return nil
}
