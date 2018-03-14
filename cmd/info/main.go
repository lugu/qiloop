package main

import (
	"encoding/json"
	"fmt"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/object"
	"github.com/lugu/qiloop/services"
	"log"
	"os"
)

func Print(i interface{}) {
	json, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		log.Fatalf("json encoding failed: %s", err)
	}
	fmt.Println(string(json))
}

func main() {
	sess, err := session.NewSession(":9559")
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	if len(os.Args) > 1 {
		serviceName := os.Args[1]
		proxy, err := sess.Proxy(serviceName, 1)
		if err != nil {
			log.Fatalf("failed to connect service (%s): %s", serviceName, err)
		}
		var obj object.Object = &services.ObjectProxy{proxy}
		meta, err := obj.MetaObject(1)
		if err != nil {
			log.Fatalf("failed to get metaobject (%s): %s", serviceName, err)
		}
		Print(meta)
	} else {
		directory, err := services.NewServiceDirectory(sess, 1)
		if err != nil {
			log.Fatalf("directory creation failed: %s", err)
		}
		services, err := directory.Services()
		if err != nil {
			log.Fatalf("failed to list services: %s", err)
		}
		Print(services)
	}
}
