package pingpong

import (
	"github.com/lugu/qiloop/bus/server"
)

type impl struct {
	signal PingPongSignalHelper
}

func NewPingPong() PingPong {
	return new(impl)
}

func (p *impl) Activate(activation server.Activation, helper PingPongSignalHelper) error {
	p.signal = helper
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
