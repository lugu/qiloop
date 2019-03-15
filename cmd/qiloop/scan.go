package main

import (
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/client/services"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/meta/idl"
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

func scan(serverURL, serviceName, idlFile string) {

	output := open(idlFile)
	defer output.Close()

	sess, err := session.NewSession(serverURL)
	if err != nil {
		log.Fatalf("failed to create session: %s", err)
	}

	// service direction id = 1
	dir, err := services.Services(sess).ServiceDirectory()
	if err != nil {
		log.Fatalf("failed to create proxy: %s", err)
	}

	serviceInfoList, err := dir.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}

	objects := make([]object.MetaObject, 0)
	objects = append(objects, object.MetaService0)
	objects = append(objects, object.ObjectMetaObject)

	for _, s := range serviceInfoList {

		if serviceName != "" && serviceName != s.Name {
			continue
		}

		// sort the addresses based on their value
		for _, addr := range s.Endpoints {
			// do not connect the test range.
			if strings.Contains(addr, "198.18.0") {
				continue
			}
			cache, err := client.NewCachedSession(addr)
			if err != nil {
				continue
			}
			err = cache.Lookup(s.Name, s.ServiceId)
			if err != nil {
				continue
			}

			meta := cache.Services[s.ServiceId]
			meta.Description = s.Name
			objects = append(objects, meta)
			err = idl.GenerateIDL(output, s.Name, meta)
			if err != nil {
				log.Printf("failed to generate IDL of %s: %s",
					s.Name, err)
			}
			break
		}
	}
}
