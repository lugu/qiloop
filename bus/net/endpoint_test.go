package net_test

import (
	"github.com/lugu/qiloop/bus/net"
	"io/ioutil"
	"log"
	gonet "net"
	"os"
	"reflect"
	"sync"
	"testing"
)

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

func TestSend(t *testing.T) {
	a, b := gonet.Pipe()
	defer a.Close()
	defer b.Close()

	e := net.NewEndPoint(a)

	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	m := net.NewMessage(h, []byte{0xab, 0xcd})

	var wait sync.WaitGroup

	wait.Add(1)
	go func() {
		err := e.Send(m)
		if err != nil {
			panic(err)
		}
		wait.Done()
	}()

	buf := make([]byte, net.HeaderSize+2)
	count, err := b.Read(buf)
	if err != nil {
		panic(err)
	}
	if count != len(buf) {
		t.Errorf("read less than expected: %d", count)
	}
	wait.Wait()
}

func TestReceiveAny(t *testing.T) {
	a, b := gonet.Pipe()
	defer a.Close()
	defer b.Close()

	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	m := net.NewMessage(h, []byte{0xab, 0xcd})

	var wait sync.WaitGroup

	wait.Add(1)
	go func() {
		err := m.Write(a)
		if err != nil {
			panic(err)
		}
		wait.Done()
	}()

	e := net.NewEndPoint(b)
	msg, err := e.ReceiveAny()
	if err != nil {
		panic(err)
	}
	if len(msg.Payload) != 2 {
		t.Errorf("read less than expected: %d", len(msg.Payload))
	}
	wait.Wait()
}

func TestPingPong(t *testing.T) {

	server, client := net.NewPipe()
	defer server.Close()
	defer client.Close()

	var wait sync.WaitGroup

	// reply to a single message
	wait.Add(1)
	go func() {
		m, err := server.ReceiveAny()
		if err != nil {
			t.Errorf("failed to receive meesage: %s", err)
		}
		err = server.Send(*m)
		if err != nil {
			t.Errorf("failed to send meesage: %s", err)
		}
		wait.Done()
	}()

	var msg *net.Message
	var err error

	wait.Add(1)
	go func() {
		// server replied
		msg, err = client.ReceiveAny()
		if err != nil {
			t.Errorf("failed to receive net. %s", err)
		}
		wait.Done()
	}()

	// client send a message
	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	mSent := net.NewMessage(h, []byte{0xab, 0xcd})
	if err := client.Send(mSent); err != nil {
		t.Errorf("failed to send paquet: %s", err)
	}

	wait.Wait()

	// check packet integrity
	if !reflect.DeepEqual(mSent, *msg) {
		t.Errorf("expected %#v, go %#v", mSent, *msg)
	}
}
