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

	var waitSend sync.WaitGroup
	var waitClose sync.WaitGroup

	waitSend.Add(2)
	waitClose.Add(2)

	event1Count := 0
	event2Count := 0
	go func() {
		for {
			_, ok := <-updates1
			if !ok {
				waitClose.Done()
				return
			}
			event1Count++
			waitSend.Done()
		}
	}()
	go func() {
		for {
			_, ok := <-updates2
			if !ok {
				waitClose.Done()
				return
			}
			event2Count++
			waitSend.Done()
		}
	}()
	err = directory.UnregisterService(1)
	if err != nil {
		t.Fatal(err)
	}
	waitSend.Wait()

	unsubscribe1()
	unsubscribe2()

	waitClose.Wait()

	if event1Count != 1 {
		t.Errorf("too many messages in 1: %d", event1Count)
	}
	if event2Count != 1 {
		t.Errorf("too many messages in 2: %d", event2Count)
	}
}
