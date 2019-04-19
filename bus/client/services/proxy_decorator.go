package services

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/client/object"
	"time"
)

// Proxy returns a proxy to the service described by info.
func (info ServiceInfo) Proxy(sess bus.Session) (object.ObjectProxy, error) {
	proxy, err := sess.Proxy(info.Name, 1)
	if err != nil {
		return nil, fmt.Errorf("cannot connect %s: %s", info.Name, err)
	}
	return object.MakeObject(proxy), nil
}

var (
	LogLevelNone    LogLevel = LogLevel{Level: 0}
	LogLevelFatal   LogLevel = LogLevel{Level: 1}
	LogLevelError   LogLevel = LogLevel{Level: 2}
	LogLevelWarning LogLevel = LogLevel{Level: 3}
	LogLevelInfo    LogLevel = LogLevel{Level: 4}
	LogLevelVerbose LogLevel = LogLevel{Level: 5}
	LogLevelDebug   LogLevel = LogLevel{Level: 6}
)

func Now() TimePoint {
	return TimePoint{
		uint64(time.Now().UnixNano()),
	}
}
