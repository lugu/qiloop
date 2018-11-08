package server_test

import (
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/util"
	gonet "net"
	"testing"
)

func helpAuth(t *testing.T, creds map[string]string, user, token string, ok bool) {
	name := util.MakeTempFileName()

	listener, err := gonet.Listen("unix", name)
	if err != nil {
		t.Fatal(err)
	}

	router := server.NewRouter(server.ServiceAuthenticate(server.Dictionary(creds)))
	srv := server.StandAloneServer(listener, router)
	go srv.Run()

	ep, err := net.DialEndPoint("unix://" + name)
	if err != nil {
		panic(err)
	}
	err = client.AuthenticateUser(ep, user, token)
	if ok && err != nil {
		t.Errorf("user: %s, token: %s, error: %s", user, token, err)
	}
	if !ok && err == nil {
		t.Errorf("shall not pass: %s, %s", user, token)
	}
}

func TestNewServiceAuthenticate(t *testing.T) {

	credentials := map[string]string{
		"foo":  "bar",
		"bazz": "bozz",
	}
	helpAuth(t, credentials, "foo", "bar", true)
	helpAuth(t, credentials, "foo", "", false)
	helpAuth(t, credentials, "", "bar", false)
	helpAuth(t, credentials, "", "", false)
	helpAuth(t, credentials, "user", "pass", false)
	helpAuth(t, credentials, "user", "", false)
	helpAuth(t, credentials, "", "pass", false)
	helpAuth(t, credentials, "bazz", "", false)
	helpAuth(t, credentials, "bazz", "bozz", true)
}
