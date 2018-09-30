package main

import (
	"github.com/lugu/qiloop/bus/session"
	"log"
	gonet "net"
)

func main() {
	listener, err := gonet.Listen("tcp", ":9559")
	if err != nil {
		log.Fatalf("%s", err)
	}
	directory := session.NewDirectoryService()
	router := session.NewRouter()
	router.Add(directory)
	server := session.NewServer(listener, router)
	server.Run()
}
