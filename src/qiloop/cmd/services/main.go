package main

import (
	"log"
	"os"
	"path/filepath"
	"qiloop/net"
	"qiloop/services"
	"qiloop/value"
)

func main() {
	endpoint := ":9559"
	conn, err := net.NewClient(endpoint)
	if err != nil {
		log.Fatalf("failed to connect %s: %s", endpoint, err)
	}
	server := services.Server{net.NewProxy(conn, 0, 0)}
	permissions := map[string]value.Value{
		"ClientServerSocket":    value.Bool(true),
		"MessageFlags":          value.Bool(true),
		"MetaObjectCache":       value.Bool(true),
		"RemoteCancelableCalls": value.Bool(true),
	}
	permissions, err = server.Authenticate(permissions)
	if err != nil {
		log.Fatalf("authentication failed: %s", err)
	}

	directory := services.Directory{net.NewProxy(conn, 1, 1)}
	serviceInfoList, err := directory.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}

	for _, s := range serviceInfoList {
		service := services.Directory{net.NewProxy(conn, s.ServiceId, 1)}
		_, err := service.MetaObject(1)
		if err != nil {
			log.Printf("failed to query MetaObject of %s: %s", s.Name, err)
		}

		filename := filepath.Join("/tmp", s.Name+".go")
		file, err := os.Create(filename)
		if err != nil {
			log.Printf("failed to open %s: %s", filename, err)
			return
		}
		defer file.Close()

		// err = proxy.GenerateProxy(metaObj.(net.MetaObject), "services", "Server", file)
		if err != nil {
			log.Printf("proxy generation failed: %s\n", err)
		}
	}
}
