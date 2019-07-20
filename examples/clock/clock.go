package clock

import (
	"time"

	"github.com/lugu/qiloop/bus"
)

type timestampService time.Time

func (t timestampService) Activate(activation bus.Activation,
	helper TimestampSignalHelper) error {
	return nil
}
func (t timestampService) OnTerminate() {
}
func (t timestampService) Nanoseconds() (int64, error) {
	return time.Since(time.Time(t)).Nanoseconds(), nil
}

// NewTimestampObject creates a timestamp object which can be
// registered to a bus.Server.
func NewTimestampObject() bus.Actor {
	return TimestampObject(timestampService{})
}
