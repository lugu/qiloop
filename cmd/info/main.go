package main

import (
	"encoding/json"
	"fmt"
	"github.com/lugu/qiloop/services"
	"github.com/lugu/qiloop/session"
	"github.com/lugu/qiloop/value"
	"log"
)

func main() {
	endpoint := ":9559"
	conn, err := session.NewClient(endpoint)
	if err != nil {
		log.Fatalf("failed to connect %s: %s", endpoint, err)
	}
	server := services.Server{session.NewProxy(conn, 0, 0)}
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

	directory := services.ServiceDirectory{session.NewProxy(conn, 1, 1)}
	services, err := directory.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}

	json, err := json.MarshalIndent(services, "", "    ")
	if err != nil {
		log.Fatalf("json encoding failed: %s", err)
	}
	fmt.Println(string(json))
}
