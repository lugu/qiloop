package main

import (
	"bytes"
	"github.com/lugu/qiloop/meta/proxy"
	object "github.com/lugu/qiloop/meta/stage1"
	server "github.com/lugu/qiloop/meta/stage2"
	directory "github.com/lugu/qiloop/meta/stage3"
	"github.com/lugu/qiloop/net"
	"github.com/lugu/qiloop/value"
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

	endpoint := ":9559"
	conn, err := net.NewClient(endpoint)
	if err != nil {
		log.Fatalf("failed to connect %s: %s", endpoint, err)
	}
	server0 := server.Server{net.NewProxy(conn, 0, 0)}
	permissions := map[string]value.Value{
		"ClientServerSocket":    value.Bool(true),
		"MessageFlags":          value.Bool(true),
		"MetaObjectCache":       value.Bool(true),
		"RemoteCancelableCalls": value.Bool(true),
	}
	permissions, err = server0.Authenticate(permissions)
	if err != nil {
		log.Fatalf("authentication failed: %s", err)
	}

	directory := directory.Directory{net.NewProxy(conn, 1, 1)}
	serviceInfoList, err := directory.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}

	objects := make([]object.MetaObject, len(serviceInfoList)+1)

	objects[0] = server.MetaService0

	for i, s := range serviceInfoList {
		service := server.Server{net.NewProxy(conn, s.ServiceId, 1)}
		metaObj, err := service.MetaObject(1)
		if err != nil {
			log.Printf("failed to query MetaObject of %s: %s", s.Name, err)
		}
		metaObj.Description = s.Name

		// Go type system: can not convert type stage2.MetaObject into
		// stage1.MetaObject the trick is to serialize and deserialize
		// the object to convert it.

		buf := bytes.NewBuffer(make([]byte, 0))
		if err = server.WriteMetaObject(metaObj, buf); err != nil {
			log.Fatalf("failed to serialize stage2.MetaObject: %s", err)
		}

		if objects[i+1], err = object.ReadMetaObject(buf); err != nil {
			log.Fatalf("failed to deserialize stage1.MetaObject: %s", err)
		}

	}
	proxy.GenerateProxys(objects, "services", output)
}
