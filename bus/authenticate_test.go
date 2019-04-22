package bus_test

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/value"
	"io"
	"testing"
)

func helpAuth(t *testing.T, creds map[string]string, user, token string, ok bool) {
	addr := util.NewUnixAddr()

	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}

	auth := bus.Dictionary(creds)
	srv, err := bus.StandAloneServer(listener, auth, bus.PrivateNamespace())
	if err != nil {
		panic(err)
	}

	ep, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	err = bus.AuthenticateUser(ep, user, token)
	if ok && err != nil {
		t.Errorf("user: %s, token: %s, error: %s", user, token, err)
	}
	if !ok && err == nil {
		t.Errorf("shall not pass: %s, %s", user, token)
	}
	srv.Terminate()
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

func LimitedReader(c bus.CapabilityMap, size int) io.Reader {
	var buf bytes.Buffer
	bus.WriteCapabilityMap(c, &buf)
	return &io.LimitedReader{
		R: &buf,
		N: int64(size),
	}
}

type LimitedWriter struct {
	size int
}

func (b *LimitedWriter) Write(buf []byte) (int, error) {
	if len(buf) <= b.size {
		b.size -= len(buf)
		return len(buf), nil
	}
	old_size := b.size
	b.size = 0
	return old_size, io.EOF
}

func NewLimitedWriter(size int) io.Writer {
	return &LimitedWriter{
		size: size,
	}
}

func TestWriterCapabilityMapError(t *testing.T) {
	c := bus.CapabilityMap{
		bus.KeyState: value.Uint(bus.StateDone),
	}
	var buf bytes.Buffer
	err := bus.WriteCapabilityMap(c, &buf)
	if err != nil {
		panic(err)
	}
	max := len(buf.Bytes())

	for i := 0; i < max-1; i++ {
		w := NewLimitedWriter(i)
		err := bus.WriteCapabilityMap(c, w)
		if err == nil {
			panic(fmt.Errorf("not expecting a success at %d", i))
		}
	}
	w := NewLimitedWriter(max)
	err = bus.WriteCapabilityMap(c, w)
	if err != nil {
		panic(err)
	}
}

func TestReadHeaderError(t *testing.T) {
	c := bus.CapabilityMap{
		bus.KeyState: value.Uint(bus.StateDone),
	}
	var buf bytes.Buffer
	err := bus.WriteCapabilityMap(c, &buf)
	if err != nil {
		panic(err)
	}
	max := len(buf.Bytes())

	for i := 0; i < max; i++ {
		r := LimitedReader(c, i)
		_, err := bus.ReadCapabilityMap(r)
		if err == nil {
			panic(fmt.Errorf("not expecting a success at %d", i))
		}
	}
	r := LimitedReader(c, max)
	_, err = bus.ReadCapabilityMap(r)
	if err != nil {
		panic(err)
	}
}

func TestAuthenticateYesNo(t *testing.T) {
	var auth bus.Authenticator
	auth = bus.No{}
	if auth.Authenticate("", "") {
		panic("error")
	}
	auth = bus.Yes{}
	if !auth.Authenticate("", "") {
		panic("error")
	}
}
