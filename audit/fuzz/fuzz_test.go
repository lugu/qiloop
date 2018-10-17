package fuzz_test

import (
	"github.com/lugu/qiloop/bus/cmd/fuzz"
	"github.com/lugu/qiloop/bus/session"
	"io/ioutil"
	gonet "net"
	"path/filepath"
	"testing"
)

func TestFuzz(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("shall panic")
		}
	}()

	passwords := map[string]string{
		"nao": "nao",
	}
	object := session.NewServiceAuthenticate(passwords)
	ns := session.NewNamespace(object)
	router := session.NewRouter()
	router.Add(ns)
	listener, err := gonet.Listen("tcp", ":9559")
	if err != nil {
		panic(err)
	}
	server := session.NewServer(listener, router)
	go server.Run()

	fuzz.ServerURL = "tcp://localhost:9559"

	data, err := ioutil.ReadFile(filepath.Join("testdata", "cap-auth-failure.bin"))
	if err != nil {
		t.Errorf("cannot open test data %s", err)
	}
	fuzz.Fuzz(data)

	server.Stop()
}
