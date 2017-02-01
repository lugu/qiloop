package main

import (
	"github.com/lugu/qiloop/bus/session/basic"
	"github.com/lugu/qiloop/meta/idl"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	var output io.Writer
	var addr string

	if len(os.Args) > 1 {
		addr = os.Args[1]
	} else {
		addr = ":9559"
	}

	if len(os.Args) > 2 {
		filename := os.Args[2]
		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
			return
		}
		defer file.Close()
		output = file
	} else {
		output = os.Stdout
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

	for i, s := range serviceInfoList {

		// FIXME: change me to 100
		if i > 1 {
			continue
		}

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
		idl.GenerateIDL(output, s.Name, meta)
	}
}
