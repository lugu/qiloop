package logger

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/lugu/qiloop/bus"
)

// ErrUnknownProvider can be returned by RemoveProvider.
var ErrUnknownProvider = fmt.Errorf("Unknown provider")

var (
	// LogLevelNone does not output logs
	LogLevelNone = LogLevel{Level: 0}
	// LogLevelFatal only shows catastrophic errors
	LogLevelFatal = LogLevel{Level: 1}
	// LogLevelError shows all errors
	LogLevelError = LogLevel{Level: 2}
	// LogLevelWarning shows errors and warning
	LogLevelWarning = LogLevel{Level: 3}
	// LogLevelInfo shows errors, warning and info
	LogLevelInfo = LogLevel{Level: 4}
	// LogLevelVerbose shows an excessive amount of log
	LogLevelVerbose = LogLevel{Level: 5}
	// LogLevelDebug shows log to debug
	LogLevelDebug = LogLevel{Level: 6}
)

type logManager struct {
	activation     bus.Activation
	providersMutex sync.RWMutex
	providers      map[int32]LogProviderProxy
	providersNext  int32
	listenersMutex sync.RWMutex
	listeners      map[uint32]*logListenerImpl
	listenersNext  uint32
	filters        map[string]LogLevel
	defaultLevel   LogLevel
}

// NewLogManager returns an object to be registered as the LogManager.
func NewLogManager() bus.Actor {
	return LogManagerObject(&logManager{
		providers: make(map[int32]LogProviderProxy),
		listeners: make(map[uint32]*logListenerImpl),
	})
}

func (l *logManager) Activate(activation bus.Activation,
	helper LogManagerSignalHelper) error {
	l.activation = activation
	return nil
}
func (l *logManager) OnTerminate() {
}

func (l *logManager) Log(messages []LogMessage) error {
	l.listenersMutex.RLock()
	defer l.listenersMutex.RUnlock()
	for _, listener := range l.listeners {
		listener.Messages(messages)
	}
	return nil
}

func (l *logManager) CreateListener() (LogListenerProxy, error) {
	l.listenersMutex.Lock()
	index := l.listenersNext
	l.listenersNext++
	l.listenersMutex.Unlock()

	listener := &logListenerImpl{
		filters:      make(map[string]LogLevel),
		filtersReg:   make(map[string]*regexp.Regexp),
		defaultLevel: LogLevelInfo,
		manager:      l,
	}
	service := l.activation.Service
	session := l.activation.Session
	proxy, err := Services(session).NewLogListener(service, listener)
	if err != nil {
		return nil, err
	}

	l.listenersMutex.Lock()
	l.listeners[index] = listener
	l.listenersMutex.Unlock()

	l.UpdateFilters()
	l.UpdateVerbosity()
	return proxy, nil
}

func (l *logManager) UpdateFilters() {
	filters := make(map[string]LogLevel)
	l.listenersMutex.RLock()
	for _, listener := range l.listeners {
		for cat, l := range listener.filters {
			previous, ok := filters[cat]
			if !ok || l.Level < previous.Level {
				filters[cat] = l
			}
		}
	}
	l.listenersMutex.RUnlock()
	l.providersMutex.RLock()
	for _, provider := range l.providers {
		provider.ClearAndSet(filters)
	}
	l.providersMutex.RUnlock()
}

func (l *logManager) UpdateVerbosity() {
	level := LogLevelNone
	l.listenersMutex.RLock()
	for _, listener := range l.listeners {
		if listener.defaultLevel.Level > level.Level {
			level = listener.defaultLevel
		}
	}
	l.listenersMutex.RUnlock()
	l.providersMutex.RLock()
	for _, provider := range l.providers {
		provider.SetVerbosity(level)
	}
	l.providersMutex.RUnlock()
}

func (l *logManager) GetListener() (LogListenerProxy, error) {
	return l.CreateListener()
}
func (l *logManager) AddProvider(provider LogProviderProxy) (int32, error) {
	l.providersMutex.Lock()
	index := l.providersNext
	l.providersNext++
	l.providers[index] = provider
	l.providersMutex.Unlock()
	l.UpdateFilters()
	l.UpdateVerbosity()
	return index, nil
}
func (l *logManager) RemoveProvider(providerID int32) error {
	l.providersMutex.Lock()
	defer l.providersMutex.Unlock()
	if _, ok := l.providers[providerID]; ok {
		delete(l.providers, providerID)
		return nil
	}
	return ErrUnknownProvider
}
