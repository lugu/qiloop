package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/hashicorp/mdns"
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
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	defer server.Terminate()

	_, err = server.NewService("LogManager", qilog.NewLogManager())
	if err != nil {
		log.Fatalf("Failed to start log manager: %s", err)
	}

	log.Printf("Listening at %s", serverURL)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	url, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal(err)
	}
	info := []string{"RobotType=NAO", "Protocol=tcp"}
	service, err := mdns.NewMDNSService(url.Hostname(), "_naoqi._tcp", "", "", 9559, nil, info)
	if err != nil {
		log.Fatalf("MDNS Server: %s", err)
	}
	mdnsServer, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		log.Fatalf("MDNS Server: %s", err)
	}
	defer mdnsServer.Shutdown()

	select {
	case err = <-server.WaitTerminate():
		if err != nil {
			log.Fatalf("Server: %s", err)
		}
	case s := <-interrupt:
		log.Printf("%v: quitting.", s)
	}
}
