package main

import (
	"github.com/lugu/qiloop/bus/services"
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

	cancel := make(chan int)

	channel, err := directory.SignalServiceRemoved(cancel)
	if err != nil {
		log.Fatalf("failed to get remote signal channel: %s", err)
	}

	for e := range channel {
		log.Printf("service removed: %s (%d)", e.P1, e.P0)
	}
}
