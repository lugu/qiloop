package pong_test

import (
	"bytes"
	"fmt"
	"io"
	gonet "net"
	"os"
	"strings"
	"testing"

	dir "github.com/lugu/qiloop/bus/directory"
	sess "github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/examples/pong"
	"github.com/lugu/qiloop/type/basic"
)

func TestPingPong(t *testing.T) {
	addr := util.NewUnixAddr()

	server, err := dir.NewServer(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Terminate()

	service := pong.PingPongObject(pong.PingPongImpl())
	_, err = server.NewService("PingPong", service)
	if err != nil {
		panic(err)
	}

	session, err := sess.NewSession(addr)
	if err != nil {
		panic(err)
	}
	services := pong.Services(session)
	client, err := services.PingPong()
	if err != nil {
		panic(err)
	}

	cancel, pong, err := client.SubscribePong()
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

	if answer != "hello" {
		panic(err)
	}
	cancel()
}

func testRemoteAddr(b *testing.B, addr string) {
	server, err := dir.NewServer(addr, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer server.Terminate()

	service := pong.PingPongObject(pong.PingPongImpl())
	_, err = server.NewService("PingPong", service)
	if err != nil {
		panic(err)
	}

	session, err := sess.NewSession(addr)
	if err != nil {
		panic(err)
	}
	services := pong.Services(session)
	client, err := services.PingPong()
	if err != nil {
		panic(err)
	}

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

func BenchmarkPingPongQUIC(b *testing.B) {
	testRemoteAddr(b, "quic://localhost:54322")
}

func BenchmarkPingPongLocal(b *testing.B) {

	addr := util.NewUnixAddr()
	server, err := dir.NewServer(addr, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer server.Terminate()

	// remove socket file: so it can't be used
	filename := strings.TrimPrefix(addr, "unix://")
	os.Remove(filename)

	service := pong.PingPongObject(pong.PingPongImpl())
	_, err = server.NewService("PingPong", service)
	if err != nil {
		panic(err)
	}

	session := server.Session()
	if err != nil {
		panic(err)
	}
	services := pong.Services(session)
	client, err := services.PingPong()
	if err != nil {
		panic(err)
	}

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

func BenchmarkPingPongUnixRaw(b *testing.B) {
	addr := util.NewUnixAddr()
	filename := strings.TrimPrefix(addr, "unix://")
	listener, err := gonet.Listen("unix", filename)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go func(conn gonet.Conn) {
			defer conn.Close()
			buf := make([]byte, 100)
			for {
				_, err := conn.Read(buf)
				if err == io.EOF {
					return
				} else if err != nil {
					panic(err)
				}
				_, err = conn.Write([]byte("Hello, World!"))
				if err != nil {
					panic(err)
				}
			}
		}(conn)
	}()
	conn, err := gonet.Dial("unix", filename)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	buf := make([]byte, 100)
	for i := 0; i < b.N; i++ {
		count, err := conn.Write([]byte("hello"))
		if err != nil {
			panic(err)
		}
		if count == 0 {
			panic("too short")
		}
		_, err = conn.Read(buf)
		if err != nil {
			panic(err)
		}
	}
}
