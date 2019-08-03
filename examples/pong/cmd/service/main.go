// Package main illustrate how to use the pong package to create a
// service "PingPong" registered to the service directory.
package main

import (
	"flag"

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
		panic(err)
	}
	defer server.Terminate()

	println("service running...")

	// wait until the server fails.
	err = <-server.WaitTerminate()
	if err != nil {
		panic(err)
	}
}
