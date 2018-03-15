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

	directory, err := services.NewServiceDirectory(sess, 1)
	if err != nil {
		log.Fatalf("failed to connect service manager: %s\n", err)
	}

	serviceList, err := directory.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s\n", err)
	}

	service := serviceList[len(serviceList)-1]
	proxy, err := sess.Proxy(service.Name, 1)
	if err != nil {
		log.Fatalf("failed get proxy of %s: %s\n", service.Name, err)
	}
	_, err = proxy.Call("exit", make([]byte, 0))
	if err != nil {
		log.Fatalf("failed call exit on %s: %s\n", service.Name, err)
	}
	log.Printf("Just killed: %s\n", service.Name)
}
