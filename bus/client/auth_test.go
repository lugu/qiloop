package client_test

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
	"io"
	gonet "net"
	"strings"
	"testing"
)

func TestAuth(t *testing.T) {
	addr := util.NewUnixAddr()

	listener, err := gonet.Listen("unix", strings.TrimPrefix(addr,
		"unix://"))
	if err != nil {
		panic(err)
	}

	auth := server.Dictionary(map[string]string{
		"user1": "aaa",
		"user2": "bbb",
	})
	server, err := server.StandAloneServer(listener, auth,
		server.PrivateNamespace())
	defer server.Terminate()
	if err != nil {
		panic(err)
	}

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	err = client.Authenticate(endpoint)
	if err == nil {
		panic("must fail")
	}
	endpoint, err = net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	err = client.AuthenticateUser(endpoint, "user1", "bbb")
	if err == nil {
		panic("must fail")
	}
	endpoint, err = net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	err = client.AuthenticateUser(endpoint, "user1", "aaa")
	if err != nil {
		panic(err)
	}
	err = client.AuthenticateUser(endpoint, "user1", "aaa")
	if err != nil {
		panic(err)
	}
	endpoint.Close()
	endpoint, err = net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	err = client.AuthenticateUser(endpoint, "user2", "aaa")
	if err == nil {
		panic("must fail")
	}
	endpoint, err = net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	err = client.AuthenticateUser(endpoint, "user2", "bbb")
	if err != nil {
		panic(err)
	}
	err = client.AuthenticateUser(endpoint, "user3", "")
	if err == nil {
		panic(err)
	}
	endpoint.Close()
}

type ServerMock struct {
	EndPoint      net.EndPoint
	ExpectedUser  string
	ExpectedToken string
	Response      func(user, token string) client.CapabilityMap
}

func NewServer(response func(user, token string) client.CapabilityMap) net.EndPoint {
	client, server := net.Pipe()
	s := &ServerMock{
		EndPoint: server,
		Response: response,
	}
	filter := func(hdr *net.Header) (bool, bool) {
		return true, true
	}
	consumer := func(msg *net.Message) error {
		return s.Receive(msg)
	}
	closer := func(err error) {
		if err != io.EOF {
			panic("unexpected error: " + err.Error())
		}
	}
	s.EndPoint.AddHandler(filter, consumer, closer)
	return client
}

func (s *ServerMock) Receive(m *net.Message) error {
	if m.Header.Action != object.AuthenticateActionID {
		return util.ReplyError(s.EndPoint, m, server.ErrActionNotFound)
	}

	if s.Response == nil {
		return util.ReplyError(s.EndPoint, m,
			fmt.Errorf("service unavailable"))
	}

	response, err := s.wrapAuthenticate(m.Payload)
	if err != nil {
		return util.ReplyError(s.EndPoint, m, err)
	}
	hdr := net.NewHeader(net.Reply, 0, 0, object.AuthenticateActionID,
		m.Header.ID)
	reply := net.NewMessage(hdr, response)
	return s.EndPoint.Send(reply)
}

func (s *ServerMock) wrapAuthenticate(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	m, err := server.ReadCapabilityMap(buf)
	if err != nil {
		return nil, err
	}
	ret := s.Authenticate(m)
	var out bytes.Buffer
	err = server.WriteCapabilityMap(ret, &out)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (s *ServerMock) capError() client.CapabilityMap {
	return client.CapabilityMap{
		client.KeyState: value.Uint(client.StateError),
	}
}

func (s *ServerMock) Authenticate(cap client.CapabilityMap) client.CapabilityMap {
	var user, token string
	if userValue, ok := cap[client.KeyUser]; ok {
		if userStr, ok := userValue.(value.StringValue); ok {
			user = userStr.Value()
		} else {
			return s.capError()
		}
	}
	if tokenValue, ok := cap[client.KeyToken]; ok {
		if tokenStr, ok := tokenValue.(value.StringValue); ok {
			token = tokenStr.Value()
		} else {
			return s.capError()
		}
	}
	return s.Response(user, token)
}

func helpTest(t *testing.T, user, token string, status uint32) {
	response := func(gotUser, gotToken string) client.CapabilityMap {
		if user != gotUser {
			panic("not expecting user " + gotUser)
		}
		if token != gotToken {
			panic("not expecting token " + gotToken)
		}
		if status == client.StateContinue {
			status = client.StateDone
			return client.CapabilityMap{
				client.KeyState:    value.Uint(client.StateContinue),
				client.KeyNewToken: value.String(token),
			}
		}
		return client.CapabilityMap{
			client.KeyState: value.Uint(status),
		}
	}
	endpoint := NewServer(response)
	defer endpoint.Close()
	err := client.AuthenticateUser(endpoint, user, token)
	switch status {
	case client.StateDone:
		if err != nil {
			panic("expecting a success, got " + err.Error())
		}
	case client.StateError:
		if err == nil {
			panic("expecting an error")
		}
	case client.StateContinue:
		if err != nil {
			panic("expecting a success, got " + err.Error())
		}
	default:
		if err == nil {
			panic("expecting an error")
		}
	}
}

func TestAuthContinue(t *testing.T) {
	helpTest(t, "userA", "correct passwd", client.StateDone)
	helpTest(t, "userA", "incorrect passwd", client.StateError)
	helpTest(t, "userA", "negotiation", client.StateContinue)
	helpTest(t, "userA", "invalid state", uint32(0xdead))
}

func helpAuthError(t *testing.T, user, token string,
	response func(user, token string) client.CapabilityMap) {

	endpoint := NewServer(response)
	err := client.AuthenticateUser(endpoint, user, token)
	if err == nil {
		panic("must fail")
	}
	endpoint.Close()
}

func TestAuthError(t *testing.T) {
	// missing state
	helpAuthError(t, "a", "b", func(user, token string) client.CapabilityMap {
		if user != "a" {
			panic("unexpected user " + user)
		}
		if token != "b" {
			panic("unexpected token " + token)
		}
		return client.CapabilityMap{}
	})

	// service not available
	helpAuthError(t, "a", "b", nil)

	// wrong state type
	helpAuthError(t, "a", "b", func(user, token string) client.CapabilityMap {
		return client.CapabilityMap{
			client.KeyState: value.String("not correct"),
		}
	})

	// missing new token
	helpAuthError(t, "a", "b", func(user, token string) client.CapabilityMap {
		return client.CapabilityMap{
			client.KeyState: value.Uint(client.StateContinue),
		}
	})

	// wrong new token type
	helpAuthError(t, "a", "b", func(user, token string) client.CapabilityMap {
		return client.CapabilityMap{
			client.KeyState:    value.Uint(client.StateContinue),
			client.KeyNewToken: value.Uint(12),
		}
	})

	// continue reply missing state
	state := 1
	helpAuthError(t, "a", "b", func(user, token string) client.CapabilityMap {
		if state == 1 {
			state = 2
			return client.CapabilityMap{
				client.KeyState:    value.Uint(client.StateContinue),
				client.KeyNewToken: value.String("bb"),
			}
		}
		if user != "a" {
			panic("unexpected user " + user)
		}
		if token != "bb" {
			panic("unexpected token " + token)
		}
		return client.CapabilityMap{}
	})

	// continue reply wrong state type
	state = 1
	helpAuthError(t, "a", "b", func(user, token string) client.CapabilityMap {
		if state == 1 {
			state = 2
			return client.CapabilityMap{
				client.KeyState:    value.Uint(client.StateContinue),
				client.KeyNewToken: value.String("bb"),
			}
		}
		return client.CapabilityMap{
			client.KeyState: value.String("bb"),
		}
	})

	// continue reply wrong state value
	state = 1
	helpAuthError(t, "a", "b", func(user, token string) client.CapabilityMap {
		if state == 1 {
			state = 2
			return client.CapabilityMap{
				client.KeyState:    value.Uint(client.StateContinue),
				client.KeyNewToken: value.String("bb"),
			}
		}
		return client.CapabilityMap{
			client.KeyState: value.Uint(1234),
		}
	})

	// continue reply state continue
	state = 1
	helpAuthError(t, "a", "b", func(user, token string) client.CapabilityMap {
		if state == 1 {
			state = 2
			return client.CapabilityMap{
				client.KeyState:    value.Uint(client.StateContinue),
				client.KeyNewToken: value.String("bb"),
			}
		}
		return client.CapabilityMap{
			client.KeyState: value.Uint(client.StateContinue),
		}
	})

	// continue reply state error
	state = 1
	helpAuthError(t, "a", "b", func(user, token string) client.CapabilityMap {
		if state == 1 {
			state = 2
			return client.CapabilityMap{
				client.KeyState:    value.Uint(client.StateContinue),
				client.KeyNewToken: value.String("bb"),
			}
		}
		return client.CapabilityMap{
			client.KeyState: value.Uint(client.StateError),
		}
	})

	// continue with call error
	state = 1
	var endpoint net.EndPoint
	response := func(user, token string) client.CapabilityMap {
		if state == 1 {
			state = 2
			return client.CapabilityMap{
				client.KeyState:    value.Uint(client.StateContinue),
				client.KeyNewToken: value.String("bb"),
			}
		}
		if token != "bb" {
			panic("expecting token bb")
		}
		endpoint.Close()
		return client.CapabilityMap{}
	}

	endpoint = NewServer(response)
	err := client.AuthenticateUser(endpoint, "a", "b")
	if err == nil {
		panic("must fail")
	}
	endpoint.Close()
}
