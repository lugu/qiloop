package clock

import (
	"fmt"
	"time"

	"github.com/lugu/qiloop/bus"
)

// Timestamper creates monotonic timestamps.
type Timestamper time.Time

// Nanoseconds returns the timestamp.
func (t Timestamper) Nanoseconds() (int64, error) {
	return time.Since(time.Time(t)).Nanoseconds(), nil
}

// timestampService implements TimestampImplementor.
type timestampService struct {
	Timestamper // inherits the Nanoseconds method.
}

// Activate is called once the service is online. It provides the
// implementation important runtime informations.
func (t timestampService) Activate(activation bus.Activation,
	helper TimestampSignalHelper) error {
	return nil
}

// OnTerminate is called when the service is termninated.
func (t timestampService) OnTerminate() {
}

// NewTimestampObject creates a timestamp object which can be
// registered to a bus.Server.
func NewTimestampObject() bus.Actor {
	return TimestampObject(timestampService{
		Timestamper(time.Now()),
	})
}

// SynchronizedTimestamper returns a locally generated timestamp
// source synchronized is the remote Timestamp service.
func SynchronizedTimestamper(session bus.Session) (Timestamper, error) {

	ref := time.Now()

	constructor := Services(session)
	timestampProxy, err := constructor.Timestamp()
	if err != nil {
		return Timestamper(ref),
			fmt.Errorf("reference timestamp: %s", err)
	}

	delta1 := time.Since(ref)
	ts, err := timestampProxy.Nanoseconds()
	delta2 := time.Since(ref)

	if err != nil {
		return Timestamper(ref),
			fmt.Errorf("reference timestamp: %s", err)
	}

	offset := ((delta1 + delta2) / 2) - time.Duration(ts)
	return Timestamper(ref.Add(offset)), nil
}
