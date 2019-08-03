// Package main illustrate how to use the pong package to create a
// service "PingPong" registered to the service directory.
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/examples/pong"
)

func main() {
	flag.Parse()

	// the bus.Actor object used by the server.
	service := pong.PingPongObject(pong.PingPongImpl())

	// create a server with a ping pong service.
	server, err := app.ServerFromFlag("PingPong", service)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Terminate()

	log.Print("PingPong service running...")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// wait until the server fails or is interrupted.
	select {
	case err = <-server.WaitTerminate():
		if err != nil {
			log.Fatal(err)
		}
	case <-interrupt:
		log.Print("interrupt, quitting.")
	}
}
