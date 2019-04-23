package session

import (
	"github.com/lugu/qiloop/bus"
	dir "github.com/lugu/qiloop/bus/directory"
	"github.com/lugu/qiloop/bus/util"
	"testing"
)

func TestNewSession(t *testing.T) {
	addr := util.NewUnixAddr()
	server, err := dir.NewServer(addr, bus.Yes{})
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	sess, err := NewSession(addr)
	if err != nil {
		t.Fatal(err)
	}
	sess.Destroy()
}

func TestNewSessionError(t *testing.T) {
	addr := util.NewUnixAddr()
	server, err := dir.NewServer(addr, bus.No{})
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	_, err = NewSession(addr)
	if err == nil {
		t.Fatal("expecting an error")
	}
}
