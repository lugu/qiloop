package main

import (
	"flag"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/type/object"
	"log"
	"os"
)

func PrintService(sess bus.Session, serviceName string) {
	proxy, err := sess.Proxy(serviceName, 1)
	if err != nil {
		log.Fatalf("failed to connect service (%s): %s", serviceName, err)
	}
	var obj object.Object = &services.ObjectProxy{proxy}
	meta, err := obj.MetaObject(1)
	if err != nil {
		log.Fatalf("failed to get metaobject (%s): %s", serviceName, err)
	}
	if err := idl.GenerateIDL(os.Stdout, serviceName, meta); err != nil {
		log.Printf("failed to generate IDL of %s: %s", serviceName, err)
	}
}

func main() {
	var serverURL = flag.String("qi-url", "tcp://127.0.0.1:9559", "server URL")
	flag.Parse()

	sess, err := session.NewSession(*serverURL)
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	if len(flag.Args()) > 1 {
		PrintService(sess, flag.Args()[1])
	} else {
		directory, err := services.NewServiceDirectory(sess, 1)
		if err != nil {
			log.Fatalf("directory creation failed: %s", err)
		}
		services, err := directory.Services()
		if err != nil {
			log.Fatalf("failed to list services: %s", err)
		}
		for _, s := range services {
			PrintService(sess, s.Name)
		}
	}
}
