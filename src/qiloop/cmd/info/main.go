package main

import (
	"encoding/json"
	"log"
	"os"
	"qiloop/net"
	"qiloop/object"
	"qiloop/services"
	"qiloop/value"
)

func main() {
	endpoint := "localhost:9559"
	conn, err := net.NewClient(endpoint)
	if err != nil {
		log.Fatalf("failed to connect %s: %s", endpoint, err)
	}
	server := services.Server{object.NewProxy(conn, 0, 0)}
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

	directory := services.Directory{object.NewProxy(conn, 1, 1)}
	services, err := directory.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}
	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(services)
	if err != nil {
		log.Fatalf("json encoding failed: %s", err)
	}
}
