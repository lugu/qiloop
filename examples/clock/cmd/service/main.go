package main

import (
	"flag"
	"fmt"
	"log"
	"syscall"

	"github.com/lugu/qiloop"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/examples/clock"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	var sessionURL = flag.String("session-url", "tcp://localhost:9559", "Service directory URL")
	var listenURL = flag.String("listen-url", "", "URL to listen to")
	var userName = flag.String("user", "", "user name")
	var userToken string
	flag.Parse()

	if *listenURL == "" {
		*listenURL = util.NewUnixAddr()
	}

	if *userName != "" {
		fmt.Println("Your token: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("Token input error: %s", err)
		}
		userToken = string(bytePassword)
	}

	// Create a new session.
	session, err := qiloop.NewSession(
		*sessionURL, *userName, userToken)
	if err != nil {
		log.Fatalf("Failed to connect %s: %s", *sessionURL, err)
	}

	auth := bus.Yes{}
	server, err := services.NewServer(session, *listenURL, auth)
	if err != nil {
		log.Fatalf("Failed to start server at %s: %s", *listenURL, err)
	}
	defer server.Terminate()

	service, err := server.NewService("Timestamp", clock.NewTimestampObject())
	if err != nil {
		log.Fatalf("Failed to register service %s: %s", "Timestamp", err)
	}
	defer service.Terminate()

	log.Printf("Service Timestamp registered")

	err = <-server.WaitTerminate()
	if err != nil {
		log.Fatalf("Terminate server: %s", err)
	}
}
