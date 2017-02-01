package session_test

import (
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/session"
	"testing"
)

func TestNewService(t *testing.T) {
	server, client := net.NewPipe()
	defer server.Close()
	defer client.Close()

	wrapper := map[uint32]bus.ActionWrapper{
		8: func(s bus.Service, d []byte) ([]byte, error) {
			return []byte{1, 2, 3, 4}, nil
		},
	}
	_ = session.NewService(0, 0, server, wrapper)
}
