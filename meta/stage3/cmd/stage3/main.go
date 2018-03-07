package main

import (
	"bytes"
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/meta/stage1"
	"github.com/lugu/qiloop/meta/stage2"
	"github.com/lugu/qiloop/net"
	"io"
	"log"
	"os"
)

func main() {
	var output io.Writer

	if len(os.Args) > 1 {
		filename := os.Args[1]

		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
			return
		}
		output = file
		defer file.Close()
	} else {
		output = os.Stdout
	}

	addr := ":9559"
	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		log.Fatalf("failed to connect: %d", err)
	}
	if err = authenticate(endpoint); err != nil {
		log.Fatalf("failed to authenticate: %d", err)
	}

	dir := stage2.Directory{manualProxy(endpoint, 1, 1)}
	serviceInfoList, err := dir.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}

	objects := make([]stage1.MetaObject, 0)
	objects = append(objects, stage2.MetaService0)

	for i, s := range serviceInfoList {

		if i > 1 {
			continue
		}

		service := stage2.Server{manualProxy(endpoint, s.ServiceId, 1)}
		metaObj, err := service.MetaObject(1)
		if err != nil {
			log.Printf("failed to query MetaObject of %s: %s", s.Name, err)
		}
		metaObj.Description = s.Name

		// Go type system: can not convert type stage2.MetaObject into
		// stage1.MetaObject the trick is to serialize and deserialize
		// the object to convert it.

		buf := bytes.NewBuffer(make([]byte, 0))
		if err = stage2.WriteMetaObject(metaObj, buf); err != nil {
			log.Fatalf("failed to serialize stage2.MetaObject: %s", err)
		}

		if object, err := stage1.ReadMetaObject(buf); err != nil {
			log.Fatalf("failed to deserialize stage1.MetaObject: %s", err)
		} else {
			objects = append(objects, object)
		}
	}
	proxy.GenerateProxys(objects, "services", output)
}
