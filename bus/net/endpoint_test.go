package net_test

import (
	"github.com/lugu/qiloop/bus/net"
	"io/ioutil"
	gonet "net"
	"os"
	"reflect"
	"testing"
)

func TestPingPong(t *testing.T) {

	server, client := net.NewPipe()
	defer server.Close()
	defer client.Close()

	// reply to a single message
	go func() {
		m, err := server.ReceiveAny()
		if err != nil {
			t.Errorf("failed to receive meesage: %s", err)
		}
		err = server.Send(*m)
		if err != nil {
			t.Errorf("failed to send meesage: %s", err)
		}
	}()

	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	mSent := net.NewMessage(h, []byte{0xab, 0xcd})

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

	// check packet integrity
	if !reflect.DeepEqual(mSent, *mReceived) {
		t.Errorf("expected %#v, go %#v", mSent, mReceived)
	}
}

func TestConnectUnix(t *testing.T) {
	f, err := ioutil.TempFile("", "go-net-test")
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
	_, err = net.DialEndPoint("unix://" + name)
	if err != nil {
		t.Fatal(err)
	}
	listener.Close()
}
