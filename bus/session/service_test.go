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
		3: func(s bus.Service, d []byte) ([]byte, error) {
			return []byte{0xab, 0xcd}, nil
		},
	}
	service := session.NewService(1, 2, server, wrapper)

	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	mSent := net.NewMessage(h, make([]byte, 0))

	// client is prepared to receive a message
	received := make(chan *net.Message)
	go func() {
		msg, err := client.ReceiveAny()
		if err != nil {
			t.Errorf("failed to receive net. %s", err)
		}
		received <- msg
	}()

	// client send a message
	if err := client.Send(mSent); err != nil {
		t.Errorf("failed to send paquet: %s", err)
	}

	// server replied
	mReceived := <-received

	service.Unregister()

	if mSent.Header.ID != mReceived.Header.ID {
		t.Errorf("invalid message id: %d", mReceived.Header.ID)
	}
	if mReceived.Header.Type != net.Reply {
		t.Errorf("invalid message type: %d", mReceived.Header.Type)
	}
	if mSent.Header.Service != mReceived.Header.Service {
		t.Errorf("invalid message service: %d", mReceived.Header.Service)
	}
	if mSent.Header.Object != mReceived.Header.Object {
		t.Errorf("invalid message object: %d", mReceived.Header.Object)
	}
	if mSent.Header.Action != mReceived.Header.Action {
		t.Errorf("invalid message action: %d", mReceived.Header.Action)
	}
	if mReceived.Payload[0] != 0xab || mReceived.Payload[1] != 0xcd {
		t.Errorf("invalid message payload: %d", mReceived.Header.Type)
	}
}
