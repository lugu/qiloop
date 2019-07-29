package main

import (
	"log"

	"github.com/lugu/qiloop/bus"
	dir "github.com/lugu/qiloop/bus/directory"
	qilog "github.com/lugu/qiloop/bus/logger"
	"github.com/lugu/qiloop/bus/session/token"
)

func server(serverURL string) {
	user, token := token.GetUserToken()
	server, err := dir.NewServer(serverURL, bus.Dictionary(
		map[string]string{
			user: token,
		}))
	defer server.Terminate()

	_, err = server.NewService("LogManager", qilog.NewLogManager())
	if err != nil {
		log.Fatalf("Failed to start log manager: %s", err)
	}

	log.Printf("Listening at %s", serverURL)

	err = <-server.WaitTerminate()
	if err != nil {
		log.Fatalf("Server terminate: %s", err)
	}
}
