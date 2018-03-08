package main

import (
	"encoding/json"
	"fmt"
	"github.com/lugu/qiloop/services"
	"github.com/lugu/qiloop/session/dummy"
	"log"
)

func main() {
	sess, err := dummy.NewSession(":9559")
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
	directory, err := services.NewServiceDirectory(sess, 1)
	if err != nil {
		log.Fatalf("directory creation failed: %s", err)
	}

	services, err := directory.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}

	json, err := json.MarshalIndent(services, "", "    ")
	if err != nil {
		log.Fatalf("json encoding failed: %s", err)
	}
	fmt.Println(string(json))
}
