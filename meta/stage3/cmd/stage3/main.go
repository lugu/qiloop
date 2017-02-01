package main

import (
	"github.com/lugu/qiloop/bus/session/basic"
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/type/object"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	var outputImplementation io.Writer
	var outputInterface io.Writer
	var addr string

	if len(os.Args) > 1 {
		filename := os.Args[1]

		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
			return
		}
		outputInterface = file
		defer file.Close()
	} else {
		outputInterface = os.Stdout
	}

	if len(os.Args) > 2 {
		filename := os.Args[2]

		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
			return
		}
		outputImplementation = file
		defer file.Close()
	} else {
		outputImplementation = os.Stdout
	}

	if len(os.Args) > 3 {
		addr = os.Args[3]
	} else {
		addr = ":9559"
	}

	// directoryServiceID := 1
	// directoryObjectID := 1
	dir, err := basic.NewServiceDirectory(addr, 1, 1, 101)
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

		addr := strings.TrimPrefix(s.Endpoints[0], "tcp://")
		obj, err := basic.NewObject(addr, s.ServiceId, 1, 2)
		if err != nil {
			log.Printf("failed to create servinceof %s: %s", s.Name, err)
			continue
		}
		meta, err := obj.MetaObject(1)
		if err != nil {
			log.Printf("failed to query MetaObject of %s: %s", s.Name, err)
			continue
		}
		meta.Description = s.Name
		objects = append(objects, meta)
	}
	proxy.GenerateInterfaces(objects, "services", outputInterface)
	proxy.GenerateProxys(objects, "services", outputImplementation)
}
