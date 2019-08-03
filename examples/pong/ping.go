// Package pong contains the implementation of the PingPong service.
//
// The file ping_stub_gen.go is generated with:
//
// 	$ qiloop stub --idl ping.qi.idl --output ping_stub_gen.go
//
package pong

import (
	"github.com/lugu/qiloop/bus"
)

type impl struct {
	signal PingPongSignalHelper
}

// PingPongImpl returns an implementation of ping pong.
func PingPongImpl() PingPongImplementor {
	return new(impl)
}

func (p *impl) Activate(activation bus.Activation, helper PingPongSignalHelper) error {
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
