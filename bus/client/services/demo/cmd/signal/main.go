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

	proxies := services.Services(sess)
	directory, err := proxies.ServiceDirectory()
	if err != nil {
		log.Fatalf("failed to connect log manager: %s", err)
	}

	unsubscribe, channel, err := directory.SubscribeServiceAdded()
	if err != nil {
		log.Fatalf("failed to get remote signal channel: %s", err)
	}

	for i := 0; i < 10; i++ {
		e := <-channel
		log.Printf("service added: %s (%d)", e.Name, e.ServiceID)
	}
	unsubscribe()
}
