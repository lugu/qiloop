package main

import (
	"flag"
	"log"

	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/examples/clock"
)

func main() {
	flag.Parse()

	server, err := app.ServerFromFlag("Timestamp", clock.NewTimestampObject())
	if err != nil {
		log.Fatalf("Failed to register service %s: %s", "Timestamp", err)
	}

	println("Timestamp service running...")

	err = <-server.WaitTerminate()
	if err != nil {
		log.Fatalf("Terminate server: %s", err)
	}
}
