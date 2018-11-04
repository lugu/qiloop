package directory_test

import (
	proxy "github.com/lugu/qiloop/bus/client/services"
	"github.com/lugu/qiloop/bus/net"
	srv "github.com/lugu/qiloop/bus/server"
	dir "github.com/lugu/qiloop/bus/server/directory"
	sess "github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/bus/util"
	"log"
	"testing"
)

func TestNewServer(t *testing.T) {
	name := util.MakeTempFileName()

	listener, err := net.Listen("unix://" + name)
	if err != nil {
		t.Fatal(err)
	}
	object := dir.ServiceDirectoryObject(dir.NewServiceDirectory())
	router := srv.NewRouter()
	router.Add(srv.NewService(object))
	server := srv.StandAloneServer(listener, router)

	go server.Run()

	session, err := sess.NewSession("unix://" + name)
	if err != nil {
		panic(err)
	}
	services := proxy.Services(session)
	directory, err := services.ServiceDirectory()
	if err != nil {
		log.Fatalf("failed to create directory: %s", err)
	}
	machineID, err := directory.MachineId()
	if err != nil {
		panic(err)
	}
	if machineID == "" {
		panic("empty machine id")
	}
}
