package logger

import (
	"fmt"
	"log"
	"regexp"
	"sync"

	"github.com/lugu/qiloop/bus"
)

// logListenerImpl implements LogListenerImplementor
type logListenerImpl struct {
	filters      map[string]LogLevel
	filtersReg   map[string]*regexp.Regexp
	defaultLevel LogLevel
	filtersMutex sync.RWMutex
	manager      *logManager

	activation bus.Activation
	helper     LogListenerSignalHelper
}

func (l *logListenerImpl) filter(msg *LogMessage) bool {
	if msg.Level.Level == LogLevelNone.Level {
		return false
	}
	l.filtersMutex.RLock()
	defer l.filtersMutex.RUnlock()

	for pattern, reg := range l.filtersReg {
		if reg.Match([]byte(msg.Category)) {
			if msg.Level.Level <= l.filters[pattern].Level {
				return true
			}
		}
	}
	if msg.Level.Level <= l.defaultLevel.Level {
		return true
	}
	return false
}

func (l *logListenerImpl) Activate(activation bus.Activation,
	helper LogListenerSignalHelper) error {

	l.helper = helper
	l.activation = activation

	if err := helper.UpdateLogLevel(LogLevelInfo); err != nil {
		return err
	}

	return nil
}

func (l *logListenerImpl) Messages(messages []LogMessage) error {
	kept := []LogMessage{}
	for _, msg := range messages {
		if l.filter(&msg) {
			l.helper.SignalOnLogMessage(msg)
			kept = append(kept, msg)
		}
	}
	return l.helper.SignalOnLogMessages(kept)
}

func (l *logListenerImpl) OnTerminate() {
	err := l.manager.terminateListener(l)
	if err != nil {
		log.Printf("failed to remove listener: %s", err)
	}
}

func validateLevel(level LogLevel) error {
	if level.Level < LogLevelNone.Level ||
		level.Level > LogLevelDebug.Level {
		return fmt.Errorf("invalid level (%d)", level.Level)
	}
	return nil
}

func (l *logListenerImpl) AddFilter(category string, level LogLevel) error {
	if err := validateLevel(level); err != nil {
		return err
	}
	reg, err := regexp.Compile(category)
	if err != nil {
		return fmt.Errorf("invalid regexp (%s): %s", category, err)
	}
	l.filtersMutex.Lock()
	defer l.filtersMutex.Unlock()
	l.filters[category] = level
	l.filtersReg[category] = reg
	l.manager.UpdateFilters()
	return nil
}

func (l *logListenerImpl) ClearFilters() error {
	l.filtersMutex.Lock()
	l.filters = make(map[string]LogLevel)
	l.filtersReg = make(map[string]*regexp.Regexp)
	l.filtersMutex.Unlock()
	l.manager.UpdateFilters()
	return nil
}

func (l *logListenerImpl) OnLogLevelChange(level LogLevel) error {
	if err := validateLevel(level); err != nil {
		return err
	}
	l.filtersMutex.Lock()
	l.defaultLevel = level
	l.filtersMutex.Unlock()
	l.manager.UpdateVerbosity()
	return nil
}
func (l *logListenerImpl) SetLevel(level LogLevel) error {
	if err := validateLevel(level); err != nil {
		return err
	}
	l.filtersMutex.Lock()
	l.defaultLevel = level
	l.filtersMutex.Unlock()
	l.manager.UpdateVerbosity()
	return nil
}
