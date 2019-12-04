package bus_test

import (
	"sync"
	"testing"

	dir "github.com/lugu/qiloop/bus/directory"
	"github.com/lugu/qiloop/bus/services"
	sess "github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/bus/util"
)

func TestTwoSubscribersDontOverlap(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := dir.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	session, err := sess.NewSession(addr)
	if err != nil {
		t.Error(err)
	}
	directory, err := services.ServiceDirectory(session)
	if err != nil {
		t.Fatalf("create directory: %s", err)
	}

	unsubscribe1, updates1, err := directory.SubscribeServiceRemoved()
	if err != nil {
		t.Fatal(err)
	}
	unsubscribe2, updates2, err := directory.SubscribeServiceRemoved()
	if err != nil {
		t.Fatal(err)
	}

	event1Count := 0
	event2Count := 0
	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		defer wait.Done()
		for {
			_, ok := <-updates1
			if !ok {
				return
			}
			event1Count++
		}
	}()
	go func() {
		defer wait.Done()
		for {
			_, ok := <-updates2
			if !ok {
				return
			}
			event2Count++
		}
	}()
	directory.UnregisterService(1)
	unsubscribe1()
	unsubscribe2()

	wait.Wait()
	if event1Count != 1 {
		t.Errorf("too many messages in 1: %d", event1Count)
	}
	if event2Count != 1 {
		t.Errorf("too many messages in 2: %d", event2Count)
	}
}
