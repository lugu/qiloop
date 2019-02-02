package main

import (
	"github.com/lugu/qiloop/bus/client/services"
	"github.com/lugu/qiloop/bus/session"
	"log"
)

func main() {
	sess, err := session.NewSession(":9559")
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	srv := services.Services(sess)
	directory, err := srv.ServiceDirectory()
	if err != nil {
		log.Fatalf("failed to connect log manager: %s", err)
	}

	cancel, channel, err := directory.SubscribeServiceAdded()
	if err != nil {
		log.Fatalf("failed to get remote signal channel: %s", err)
	}

	for i := 0; i < 10; i++ {
		e := <-channel
		log.Printf("service added: %s (%d)", e.P1, e.P0)
	}
	cancel()
}
