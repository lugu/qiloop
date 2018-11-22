package pingpong_test

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/bus/server/benchmark/pingpong"
	"github.com/lugu/qiloop/bus/server/benchmark/pingpong/proxy"
	dir "github.com/lugu/qiloop/bus/server/directory"
	sess "github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/basic"
	"os"
	"strings"
	"testing"
)

func TestPingPong(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := dir.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	service := pingpong.PingPongObject(pingpong.NewPingPong())
	_, err = server.NewService("PingPong", service)
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
	pong, err := client.SignalPong(cancel)
	if err != nil {
		panic(err)
	}

	response, err := client.Hello("hello")
	if err != nil {
		panic(err)
	}
	if response != "Hello, World!" {
		t.Errorf("wrong reply: %s", response)
	}
	client.Ping("hello")
	answer := <-pong

	if answer.P0 != "hello" {
		panic(err)
	}
}

func testRemoteAddr(b *testing.B, addr string) {
	server, err := dir.NewServer(addr, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer server.Stop()

	service := pingpong.PingPongObject(pingpong.NewPingPong())
	_, err = server.NewService("PingPong", service)
	if err != nil {
		panic(err)
	}

	session, err := sess.NewSession(addr)
	if err != nil {
		panic(err)
	}
	services := proxy.Services(session)
	client, err := services.PingPong()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reply, err := client.Hello("hello")
		if err != nil {
			panic(err)
		}
		if reply != "Hello, World!" {
			panic(reply)
		}
	}
}

func BenchmarkPingPongUnix(b *testing.B) {
	testRemoteAddr(b, util.NewUnixAddr())
}

func BenchmarkPingPongTCP(b *testing.B) {
	testRemoteAddr(b, "tcp://localhost:12345")
}

// mkdir ~/.qiloop; cd ~/.qiloop
// openssl genrsa -out server.key 2048
// openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
// cat <<EOF> ~/.qi-cert.conf
// /home/user/.qiloop/server.crt
// /home/user/.qiloop/server.key
// EOF
func BenchmarkPingPongTLS(b *testing.B) {
	testRemoteAddr(b, "tcps://localhost:54321")
}

func BenchmarkPingPongLocal(b *testing.B) {

	addr := util.NewUnixAddr()
	server, err := dir.NewServer(addr, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer server.Stop()

	// remove socket file: so it can't be used
	filename := strings.TrimPrefix(addr, "unix://")
	os.Remove(filename)

	service := pingpong.PingPongObject(pingpong.NewPingPong())
	_, err = server.NewService("PingPong", service)
	if err != nil {
		panic(err)
	}

	session := server.Session()
	if err != nil {
		panic(err)
	}
	services := proxy.Services(session)
	client, err := services.PingPong()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reply, err := client.Hello("hello")
		if err != nil {
			panic(err)
		}
		if reply != "Hello, World!" {
			panic(reply)
		}
	}
}

func implHello(a string) (string, error) {
	return "Hello, World!", nil
}

func stubHello(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	a, err := basic.ReadString(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read a: %s", err)
	}
	ret, callErr := implHello(a)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteString(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}

func proxyHello(P0 string) (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := stubHello(buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call hello failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse hello response: %s", err)
	}
	return ret, nil
}

func BenchmarkSerialization(b *testing.B) {

	for i := 0; i < b.N; i++ {
		reply, err := proxyHello("hello")
		if err != nil {
			panic(err)
		}
		if reply != "Hello, World!" {
			panic(reply)
		}
	}
}

func BenchmarkImplementation(b *testing.B) {

	for i := 0; i < b.N; i++ {
		reply, err := implHello("hello")
		if err != nil {
			panic(err)
		}
		if reply != "Hello, World!" {
			panic(reply)
		}
	}
}
