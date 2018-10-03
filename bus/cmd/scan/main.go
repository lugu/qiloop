package main

import (
	"flag"
	"github.com/lugu/qiloop/bus/session/basic"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/type/object"
	"io"
	"log"
	"os"
	"strings"
)

func open(filename string) io.WriteCloser {
	switch filename {
	case "":
		file, err := os.Create(os.DevNull)
		if err != nil {
			log.Fatalf("failed to open %s: %s", os.DevNull, err)
		}
		return file
	case "-":
		return os.Stdout
	default:
		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
		}
		return file
	}
}

func main() {
	var serverURL = flag.String("qi-url", "tcp://localhost:9559", "server URL")
	var interfaceFile = flag.String("interface", "", "File to write interface (default none)")
	var implemFile = flag.String("proxy", "", "File to write proxy (default none)")
	var idlFile = flag.String("idl", "-", "File to write IDL definition (default stdout)")
	var serviceName = flag.String("service", "", "Name of the service (default all)")

	flag.Parse()

	outputInterface := open(*interfaceFile)
	defer outputInterface.Close()

	outputImplementation := open(*implemFile)
	defer outputImplementation.Close()

	outputIDL := open(*idlFile)
	defer outputIDL.Close()

	// service direction id = 1
	// object id = 1
	dir, err := basic.NewServiceDirectory(*serverURL, 1, 1, 101)
	if err != nil {
		log.Fatalf("failed to create directory: %s", err)
	}

	serviceInfoList, err := dir.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}

	objects := make([]object.MetaObject, 0)
	objects = append(objects, object.MetaService0)
	objects = append(objects, object.ObjectMetaObject)

	for _, s := range serviceInfoList {

		if *serviceName != "" && *serviceName != s.Name {
			continue
		}

		// sort the addresses based on their value
		for _, ep := range s.Endpoints {
			// do not connect the test range.
			// FIXME: unless a local interface has such IP
			// address.
			if strings.Contains(ep, "198.18.0") {
				continue
			}
			obj, err := basic.NewObject(ep, s.ServiceId, 1, 2)
			if err != nil {
				// log.Printf("failed to create service of %s: %s", s.Name, err)
				continue
			}
			meta, err := obj.MetaObject(1)
			if err != nil {
				log.Printf("failed to query MetaObject of %s: %s", s.Name, err)
				break
			}
			meta.Description = s.Name
			objects = append(objects, meta)
			if err := idl.GenerateIDL(outputIDL, s.Name, meta); err != nil {
				log.Printf("failed to generate IDL of %s: %s", s.Name, err)
			}
			break
		}
	}
	proxy.GenerateInterfaces(objects, "services", outputInterface)
	proxy.GenerateProxys(objects, "services", outputImplementation)
}
