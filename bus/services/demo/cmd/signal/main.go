package main

import (
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
)

func main() {
	sess, err := session.NewSession(":9559")
	if err != nil {
		panic(err)
	}

	proxies := services.Services(sess)

	directory, err := proxies.ServiceDirectory()
	if err != nil {
		panic(err)
	}

	unsubscribe, channel, err := directory.SubscribeServiceAdded()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		e := <-channel
		println("service " + e.Name)
	}
	unsubscribe()
}
