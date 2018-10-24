package main

import (
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
	"log"
)

func main() {
	sess, err := session.NewSession(":9559")
	if err != nil {
		log.Fatalf("failed to connect: %s\n", err)
	}

	srv := services.Services(sess)
	directory, err := srv.ServiceDirectory()
	if err != nil {
		log.Fatalf("failed to connect service manager: %s\n", err)
	}

	serviceList, err := directory.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s\n", err)
	}

	log.Printf("%d services found.\n", len(serviceList))
}
