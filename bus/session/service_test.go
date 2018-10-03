package session_test

import (
	"bytes"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/type/value"
	"io/ioutil"
	gonet "net"
	"os"
	"testing"
)

func TestNewServer(t *testing.T) {

	f, err := ioutil.TempFile("", "go-server-test")
	if err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	f.Close()
	os.Remove(name)

	listener, err := gonet.Listen("unix", name)
	if err != nil {
		t.Fatal(err)
	}

	wrapper := bus.Wrapper{
		3: func(d []byte) ([]byte, error) {
			return []byte{0xab, 0xcd}, nil
		},
	}
	object := &session.ObjectDispather{
		Wrapper: wrapper,
	}

	ns := session.NewNamespace(object)
	router := session.NewRouter()
	router.Add(ns)
	server := session.NewServer(listener, router)
	go server.Run()

	conn, err := gonet.Dial("unix", name)
	if err != nil {
		panic(err)
	}

	client := net.NewEndPoint(conn)

	h := net.NewHeader(net.Call, 0, 0, 3, 4)
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

	server.Stop()

	if mReceived.Header.Type == net.Error {
		buf := bytes.NewBuffer(mReceived.Payload)
		errV, err := value.NewValue(buf)
		if err != nil {
			t.Errorf("invalid error value: %v", mReceived.Payload)
		}
		if str, ok := errV.(value.StringValue); ok {
			t.Errorf("error: %s", string(str))
		} else {
			t.Errorf("invalid error: %v", mReceived.Payload)
		}
	}

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
