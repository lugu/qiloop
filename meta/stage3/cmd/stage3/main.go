package main

import (
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/meta/stage2"
	"github.com/lugu/qiloop/net"
	"github.com/lugu/qiloop/object"
	"io"
	"log"
	"os"
	"strings"
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

	objects := make([]object.MetaObject, 0)
	objects = append(objects, object.MetaService0)

	for _, s := range serviceInfoList {

		addr := strings.TrimPrefix(s.Endpoints[0], "tcp://")
		endpoint, err := net.DialEndPoint(addr)
		if err != nil {
			log.Fatalf("failed to connect: %d", err)
		}
		if err = authenticate(endpoint); err != nil {
			log.Fatalf("failed to authenticate: %d", err)
		}

		service := stage2.Directory{manualProxy(endpoint, s.ServiceId, 1)}
		metaObj, err := service.MetaObject(1)
		if err != nil {
			log.Printf("failed to query MetaObject of %s: %s", s.Name, err)
		}
		metaObj.Description = s.Name
		objects = append(objects, metaObj)
	}
	proxy.GenerateProxys(objects, "services", output)
}
