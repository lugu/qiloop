package main

import (
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
	"log"
)

func main() {
	sess, err := session.NewSession("tcp://localhost:9559")
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	objectID := uint32(1)
	directory, err := services.NewServiceDirectory(sess, objectID)
	if err != nil {
		log.Fatalf("failed to create directory: %s", err)
	}

	serviceList, err := directory.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}

	for _, info := range serviceList {
		log.Printf("service %s, id: %d", info.Name, info.ServiceId)
	}
}
