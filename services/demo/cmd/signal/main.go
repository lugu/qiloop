package main

import (
	"github.com/lugu/qiloop/services"
	"github.com/lugu/qiloop/session/dummy"
	"log"
)

func main() {
	sess, err := dummy.NewSession(":9559")
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	directory, err := services.NewServiceDirectory(sess, 1)
	if err != nil {
		log.Fatalf("failed to connect log manager: %s", err)
	}

	id, err := directory.RegisterEvent(1, 107, 0x0000006b00000001)
	if err != nil {
		log.Fatalf("failed to register event: %s", err)
	}

	cancel := make(chan int)

	channel, err := directory.SignalServiceRemoved(cancel)
	if err != nil {
		log.Fatalf("failed to get remote signal channel: %s", err)
	}

	for e := range channel {
		log.Printf("service removed: %s (%d)", e.P1, e.P0)
	}

	err = directory.UnregisterEvent(1, 107, id)
	if err != nil {
		log.Fatalf("failed to register event: %s", err)
	}
}
