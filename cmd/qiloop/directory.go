package main

import (
	"log"

	"github.com/lugu/qiloop/bus"
	dir "github.com/lugu/qiloop/bus/directory"
)

func directory(serverURL string) {
	server, err := dir.NewServer(serverURL, bus.Yes{})
	if err != nil {
		log.Fatalf("Failed to listen at %s: %s", serverURL, err)
	}
	defer server.Terminate()

	log.Printf("Listening at %s", serverURL)

	err = <-server.WaitTerminate()
	if err != nil {
		log.Fatalf("Server terminate: %s", err)
	}
}
