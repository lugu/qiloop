package clock

import (
	"time"

	"github.com/lugu/qiloop/bus"
)

type timestampService struct{}

func (t timestampService) Activate(activation bus.Activation,
	helper TimestampSignalHelper) error {
	return nil
}
func (t timestampService) OnTerminate() {
}
func (t timestampService) Nanoseconds() (int64, error) {
	return time.Now().UnixNano(), nil
}

// NewTimestampObject returns an object to be used as a service.
func NewTimestampObject() bus.Actor {
	return TimestampObject(timestampService{})
}
