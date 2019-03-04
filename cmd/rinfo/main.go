package main

import (
	"encoding/json"
	"flag"
	"fmt"
	objproxy "github.com/lugu/qiloop/bus/client/object"
	"github.com/lugu/qiloop/bus/client/services"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/type/object"
	"log"
)

func Print(i interface{}) {
	json, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		log.Fatalf("json encoding failed: %s", err)
	}
	fmt.Println(string(json))
}

func main() {
	var serverURL = flag.String("qi-url", "tcp://127.0.0.1:9559", "server URL")
	var serviceName = flag.String("service", "", "name of the service to lookup")
	flag.Parse()

	sess, err := session.NewSession(*serverURL)
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
	srv := services.Services(sess)

	if *serviceName == "" {
		directory, err := srv.ServiceDirectory()
		if err != nil {
			log.Fatalf("directory creation failed: %s", err)
		}
		services, err := directory.Services()
		if err != nil {
			log.Fatalf("failed to list services: %s", err)
		}
		Print(services)
	} else {
		proxy, err := sess.Proxy(*serviceName, 1)
		if err != nil {
			log.Fatalf("failed to connect service (%s): %s",
				*serviceName, err)
		}
		var obj object.Object = objproxy.MakeObject(proxy)
		meta, err := obj.MetaObject(1)
		if err != nil {
			log.Fatalf("failed to get metaobject (%s): %s",
				*serviceName, err)
		}
		Print(meta)
	}
}
