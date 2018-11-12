package pingpong_test

import (
	"github.com/lugu/qiloop/bus/server/benchmark/pingpong"
	"github.com/lugu/qiloop/bus/server/benchmark/pingpong/proxy"
	"github.com/lugu/qiloop/bus/server/directory"
)

func TestPingPong(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := dir.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	go server.Run()
	defer server.Stop()

	service := pingpong.PingPongObject(NewPingPong())
	_, err := server.NewService("PingPong", service)
	if err != nil {
		panic(err)
	}

	session, err := sess.NewSession(addr)
	if err != nil {
		panic(err)
	}
	services := proxy.Services(session)
	client, err := services.PingPong()

	cancel := make(chan int)
	pong := client.SignalPong(cancel)

	client.Ping("hello")
	answer := <-pong

	if answer.P0 != "hello" {
		panic(err)
	}
}
