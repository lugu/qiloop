package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/examples/clock"
)

func main() {
	flag.Parse()

	server, err := app.ServerFromFlag("Timestamp", clock.NewTimestampObject())
	if err != nil {
		log.Fatal(err)
	}
	defer server.Terminate()

	log.Print("Timestamp service running...")

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
