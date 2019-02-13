package generic

import (
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"testing"
)

func TestBasicObjectWrap(t *testing.T) {
	obj := NewBasicObject()
	passed := false

	obj.Wrap(123, func(payload []byte) ([]byte, error) {
		passed = true
		return nil, nil
	})
	in, out := net.Pipe()

	channel, err := out.ReceiveAny()
	if err != nil {
		t.Error(err)
	}

	ctx := server.NewContext(in)
	hdr := net.NewHeader(net.Call, 0, 0, 123, 0)
	msg := net.NewMessage(hdr, nil)

	err = obj.Receive(&msg, ctx)
	if err != nil {
		t.Error(err)
	}
	if passed == false {
		t.Errorf("failed to pass")
	}
	reply := <-channel
	if reply.Header.Type != net.Reply {
		t.Errorf("type is %d", reply.Header.Type)
	}

	channel, err = out.ReceiveAny()
	if err != nil {
		t.Error(err)
	}

	hdr = net.NewHeader(net.Call, 0, 0, 124, 0)
	msg = net.NewMessage(hdr, nil)

	err = obj.Receive(&msg, ctx)
	if err != nil {
		t.Error(err)
	}
	reply = <-channel
	if reply.Header.Type != net.Error {
		t.Errorf("type is %d", reply.Header.Type)
	}
}
