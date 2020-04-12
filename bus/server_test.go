package bus_test

import (
	"testing"

	sess "github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/examples/pong"
	"github.com/lugu/qiloop/bus/directory"
)

func helpConnectClient(t *testing.T, addr string) {

	session, err := sess.NewSession(addr)
	if err != nil {
		t.Error(err)
	}
	defer session.Terminate()

	// obtain a representation of the ping pong service
	client, err := pong.PingPong(session)
	if err != nil {
		t.Error(err)
	}

	// subscribe to the signal "pong" of the ping pong service.
	cancel, pong, err := client.SubscribePong()
	if err != nil {
		t.Error(err)
	}
	defer cancel()

	for i := 0; i < 10; i++ {
		err = client.Ping("hello")
		if err != nil {
			t.Error(err)
		}
	}

	for i := 0; i < 10; i++ {
		<-pong
	}
}

func TestServerConnection(t *testing.T) {

	addr := util.NewUnixAddr()
	server, err := directory.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()


	// the bus.Actor object used by the server.
	actor := pong.PingPongObject(pong.PingPongImpl())
	service, err := server.NewService("PingPong", actor)
	defer service.Terminate()

	// create a server with a ping pong service.
	helpConnectClient(t, addr)
	helpConnectClient(t, addr)
	helpConnectClient(t, addr)
}

