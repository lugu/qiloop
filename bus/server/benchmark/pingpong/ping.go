package pingpong

import (
	"github.com/lugu/qiloop/bus"
)

type impl struct {
	signal PingPongSignalHelper
}

func NewPingPong() PingPong {
	return new(impl)
}

func (p *impl) Activate(sess bus.Session, serviceID, objectID uint32,
	signal PingPongSignalHelper) error {
	p.signal = signal
	return nil
}

func (p *impl) OnTerminate() {
}

func (p *impl) Hello(a string) (string, error) {
	return "Hello, World!", nil
}

func (p *impl) Ping(a string) error {
	return p.signal.SignalPong(a)
}
