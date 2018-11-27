package net_test

import (
	"github.com/lugu/qiloop/bus/net"
	"io/ioutil"
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

func TestReceiveOne(t *testing.T) {
	a, b := gonet.Pipe()
	defer a.Close()
	defer b.Close()

	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	m := net.NewMessage(h, []byte{0xab, 0xcd})

	var wait sync.WaitGroup

	e := net.NewEndPoint(b)
	msgChan, err := e.ReceiveAny()
	if err != nil {
		panic(err)
	}
	err = m.Write(a)
	if err != nil {
		panic(err)
	}

	msg, ok := <-msgChan
	if !ok {
		panic("failed to receive")
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

	srvChan, err := server.ReceiveAny()
	// reply to a single message
	wait.Add(1)
	go func() {
		m, ok := <-srvChan
		if !ok {
			t.Fatalf("failed to receive meesage")
		}
		err = server.Send(*m)
		if err != nil {
			t.Errorf("failed to send meesage: %s", err)
		}
		wait.Done()
	}()

	var mRcv *net.Message
	var mSent net.Message

	// server replied
	cltChan, err := client.ReceiveAny()
	if err != nil {
		t.Fatalf("failed to receive net.")
	}

	// client send a message
	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	mSent = net.NewMessage(h, []byte{0xab, 0xcd})
	if err := client.Send(mSent); err != nil {
		t.Errorf("failed to send paquet: %s", err)
	}

	mRcv, ok := <-cltChan
	if !ok {
		panic("chanel closed")
	}
	wait.Wait()

	// check packet integrity
	if !reflect.DeepEqual(mSent, *mRcv) {
		t.Errorf("expected %#v, go %#v", mSent, *mRcv)
	}
}

func TestEndPointFinalizer(t *testing.T) {
	a, b := gonet.Pipe()
	defer a.Close()
	defer b.Close()
	go func() {
		msg := net.NewMessage(net.NewHeader(net.Call, 1, 1, 1, 1), []byte{})
		msg.Write(a)
	}()
	filter := func(hrd *net.Header) (bool, bool) {
		return true, false
	}
	wait := make(chan *net.Message)
	consumer := func(msg *net.Message) error {
		wait <- msg
		return nil
	}
	closer := func(err error) {
	}
	finalizer := func(e net.EndPoint) {
		e.AddHandler(filter, consumer, closer)
	}
	net.EndPointFinalizer(b, finalizer)
	msg, ok := <-wait
	if !ok {
		panic("expecting a message")
	}
	if msg == nil {
		panic("not expected nil")
	}
}
