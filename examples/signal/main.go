// Package main illustrates how to subsribe to a signal in order to
// receive events. It uses the specialized proxy of the service directory.
package main

import (
	"flag"

	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/bus/services"
)

func main() {
	flag.Parse()
	// session represents a connection to the service directory.
	session, err := app.SessionFromFlag()
	if err != nil {
		panic(err)
	}

	// proxies is an helper to access the specialized proxy.
	proxies := services.Services(session)

	// obtain a representation of the service directory
	directory, err := proxies.ServiceDirectory(nil)
	if err != nil {
		panic(err)
	}

	var unsubscribe func()
	var channel chan services.ServiceAdded

	// subscribe to the signal "serviceAdded" of the service directory.
	unsubscribe, channel, err = directory.SubscribeServiceAdded()
	if err != nil {
		panic(err)
	}

	// wait until 10 services have been added.
	for i := 0; i < 10; i++ {
		e := <-channel
		println("service " + e.Name)
	}
	unsubscribe()
}
