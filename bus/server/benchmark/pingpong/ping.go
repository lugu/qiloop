package pingpong

import (
	"github.com/lugu/qiloop/bus/session"
)

type impl struct {
	signal PingPongSignalHelper
}

func NewPingPong() PingPong {
	return &impl{}
}

func (p *impl) Activate(sess bus.Session, serviceID, objectID uint32,
	signal PingPongSignalHelper) {
	p.signal = signal
}

func (p *impl) Ping(a string) error {
	return p.signal.SignalPong(a)
}
