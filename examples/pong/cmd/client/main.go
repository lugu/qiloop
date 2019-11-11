// Package main illustrate how to interract with the PingPong service.
// It calls ten times the "ping" method and waits for ten signals
// "pong".
package main

import (
	"flag"

	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/examples/pong"
)

func main() {
	flag.Parse()

	// session represents a connection to the service directory.
	session, err := app.SessionFromFlag()
	if err != nil {
		panic(err)
	}

	// proxies is an helper to access the specialized proxy.
	proxies := pong.Services(session)

	// obtain a representation of the ping pong service
	client, err := proxies.PingPong(nil)
	if err != nil {
		panic(err)
	}

	// subscribe to the signal "pong" of the ping pong service.
	cancel, pong, err := client.SubscribePong()
	if err != nil {
		panic(err)
	}
	defer cancel()

	for i := 0; i < 10; i++ {
		err = client.Ping("hello")
		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < 10; i++ {
		answer := <-pong
		println(answer)
	}
}
