package main

import (
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/stage2"
	"github.com/lugu/qiloop/type/object"
	"io"
	"log"
	"os"
)

var ServicesActionID = uint32(101)

var metaServiceDirectory object.MetaObject = object.MetaObject{
	Description: "ServiceDirectory",
	Methods: map[uint32]object.MetaMethod{
		ServicesActionID: object.MetaMethod{
			Uid:  ServicesActionID,
			Name: "services",
		},
	},
}

func NewServiceDirectory(addr string, serviceID uint32) (d stage2.ServiceDirectory, err error) {

	cache, err := client.NewCachedSession(addr)
	if err != nil {
		log.Fatalf("authentication failed: %s", err)
	}

	err = cache.Lookup("ServiceDirectory", serviceID)
	if err != nil {
		log.Fatalf("failed to query meta object of ServiceDirectory: %s", err)
	}

	return stage2.NewServiceDirectory(cache.Session(), serviceID)
}

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

	directoryServiceID := uint32(1)
	dir, err := NewServiceDirectory(addr, directoryServiceID)
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
		cache, err := client.NewCachedSession(s.Endpoints[0])
		if err != nil {
			log.Fatalf("authentication failed: %s", err)
		}
		err = cache.Lookup(s.Name, s.ServiceId)
		if err != nil {
			log.Fatalf("failed to query meta object: %s", err)
		}

		meta := cache.Services[s.ServiceId]
		if err := idl.GenerateIDL(output, s.Name, meta); err != nil {
			log.Printf("failed to generate IDL of %s: %s", s.Name, err)
			continue
		}
	}
}
