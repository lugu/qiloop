package logger

import (
	"sync"

	"github.com/lugu/qiloop/bus"
)

var (
	// LogLevelNone does not output logs
	LogLevelNone LogLevel = LogLevel{Level: 0}
	// LogLevelFatal only shows catastrophic errors
	LogLevelFatal LogLevel = LogLevel{Level: 1}
	// LogLevelError shows all errors
	LogLevelError LogLevel = LogLevel{Level: 2}
	// LogLevelWarning shows errors and warning
	LogLevelWarning LogLevel = LogLevel{Level: 3}
	// LogLevelInfo shows errors, warning and info
	LogLevelInfo LogLevel = LogLevel{Level: 4}
	// LogLevelVerbose shows an excessive amount of log
	LogLevelVerbose LogLevel = LogLevel{Level: 5}
	// LogLevelDebug shows log to debug
	LogLevelDebug LogLevel = LogLevel{Level: 6}
)

type logClient struct {
	send chan []LogMessage
	dest LogListenerProxy
}

type logManager struct {
	clientsMutex sync.RWMutex
	clients      map[uint32]logClient
	clientsNext  uint32
	activation   bus.Activation
}

// NewLogManager returns an object to be registered as the LogManager.
func NewLogManager() bus.Actor {
	return LogManagerObject(&logManager{
		clients: make(map[uint32]logClient),
	})
}

func (l *logManager) Activate(activation bus.Activation,
	helper LogManagerSignalHelper) error {
	l.activation = activation
	return nil
}
func (l *logManager) OnTerminate() {
	l.clientsMutex.RLock()
	defer l.clientsMutex.RUnlock()
	for _, client := range l.clients {
		close(client.send)
	}
}

func (l *logManager) Log(messages []LogMessage) error {
	l.clientsMutex.RLock()
	defer l.clientsMutex.RUnlock()
	for _, client := range l.clients {
		client.send <- messages
	}
	return nil
}

func (l *logManager) CreateListener() (LogListenerProxy, error) {
	l.clientsMutex.Lock()
	index := l.clientsNext
	l.clientsNext++
	l.clientsMutex.Unlock()

	onTerminate := func() {
		l.clientsMutex.Lock()
		delete(l.clients, index)
		l.clientsMutex.Unlock()
	}

	send := make(chan []LogMessage)
	proxy, err := CreateLogListener(l.activation.Session,
		l.activation.Service, send, onTerminate)

	if err != nil {
		return nil, err
	}

	l.clientsMutex.Lock()
	l.clients[index] = logClient{
		send: send,
		dest: proxy,
	}
	l.clientsMutex.Unlock()
	return proxy, nil
}

func (l *logManager) GetListener() (LogListenerProxy, error) {
	return l.CreateListener()
}
func (l *logManager) AddProvider(source LogProviderProxy) (int32, error) {
	panic("not yet implemented")
}
func (l *logManager) RemoveProvider(providerID int32) error {
	panic("not yet implemented")
}
