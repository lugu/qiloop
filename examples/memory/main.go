// Package main illustrates how to make a method call to a remote
// object. It uses the specialized proxy of the memory service.
package main

import (
	"flag"
	"fmt"

	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/examples/memory/proxy"
	"github.com/lugu/qiloop/type/value"
)

func main() {
	flag.Parse()
	// session represents a connection to the service directory.
	session, err := app.SessionFromFlag()
	if err != nil {
		panic(err)
	}
	defer session.Terminate()

	// Access the specialized proxy constructor.
	proxies := proxy.Services(session)

	// Obtain a proxy to the service
	tts, err := proxies.ALTextToSpeech(nil)
	if err != nil {
		panic(err)
	}

	// Obtain a proxy to the service
	alMemory, err := proxies.ALMemory(nil)
	if err != nil {
		panic(err)
	}

	events, err := alMemory.GetEventList()
	if err != nil {
		panic(err)
	}

	if len(events) == 0 {
		panic("no events available")
	}

	fmt.Printf("robot can subscribe to %d events \n", len(events))

	var unsubscribe func()
	var channel chan value.Value

	alMemorySub, err := alMemory.Subscriber("FrontTactilTouched")
	if err != nil {
		panic(err)
	}

	// subscribe to the signal "FrontTactilTouched" of the memory service.
	unsubscribe, channel, err = alMemorySub.SubscribeSignal()
	if err != nil {
		panic(err)
	}

	// wait until touch signal is received.
	event := <-channel
	if event != nil {
		tts.Say("ouch")
	}

	unsubscribe()
}
