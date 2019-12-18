package bus_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
)

func TestAuth(t *testing.T) {
	addr := util.NewUnixAddr()

	listener, err := net.Listen(addr)
	if err != nil {
		panic(err)
	}

	auth := bus.Dictionary(map[string]string{
		"user1": "aaa",
		"user2": "bbb",
	})
	server, err := bus.StandAloneServer(listener, auth,
		bus.PrivateNamespace())
	defer server.Terminate()
	if err != nil {
		panic(err)
	}

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	err = bus.Authenticate(endpoint)
	if err == nil {
		panic("must fail")
	}
	endpoint, err = net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	err = bus.AuthenticateUser(endpoint, "user1", "bbb")
	if err == nil {
		panic("must fail")
	}
	endpoint, err = net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	err = bus.AuthenticateUser(endpoint, "user1", "aaa")
	if err != nil {
		panic(err)
	}
	err = bus.AuthenticateUser(endpoint, "user1", "aaa")
	if err != nil {
		panic(err)
	}
	endpoint.Close()
	endpoint, err = net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	err = bus.AuthenticateUser(endpoint, "user2", "aaa")
	if err == nil {
		panic("must fail")
	}
	endpoint, err = net.DialEndPoint(addr)
	if err != nil {
		panic(err)
	}
	err = bus.AuthenticateUser(endpoint, "user2", "bbb")
	if err != nil {
		panic(err)
	}
	err = bus.AuthenticateUser(endpoint, "user3", "")
	if err == nil {
		panic(err)
	}
	endpoint.Close()
}

type ServerMock struct {
	EndPoint      net.EndPoint
	ExpectedUser  string
	ExpectedToken string
	Response      func(user, token string) bus.CapabilityMap
}

func NewServer(response func(user, token string) bus.CapabilityMap) net.EndPoint {
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
	channel := bus.NewContext(s.EndPoint)
	if m.Header.Action != object.AuthenticateActionID {
		return channel.SendError(m, bus.ErrActionNotFound)
	}

	if s.Response == nil {
		return channel.SendError(m,
			fmt.Errorf("service unavailable"))
	}

	response, err := s.wrapAuthenticate(m.Payload)
	if err != nil {
		return channel.SendError(m, err)
	}
	hdr := net.NewHeader(net.Reply, 0, 0, object.AuthenticateActionID,
		m.Header.ID)
	reply := net.NewMessage(hdr, response)
	return s.EndPoint.Send(reply)
}

func (s *ServerMock) wrapAuthenticate(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	m, err := bus.ReadCapabilityMap(buf)
	if err != nil {
		return nil, err
	}
	ret := s.Authenticate(m)
	var out bytes.Buffer
	err = bus.WriteCapabilityMap(ret, &out)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (s *ServerMock) capError() bus.CapabilityMap {
	return bus.CapabilityMap{
		bus.KeyState: value.Uint(bus.StateError),
	}
}

func (s *ServerMock) Authenticate(cap bus.CapabilityMap) bus.CapabilityMap {
	var user, token string
	if userValue, ok := cap[bus.KeyUser]; ok {
		if userStr, ok := userValue.(value.StringValue); ok {
			user = userStr.Value()
		} else {
			return s.capError()
		}
	}
	if tokenValue, ok := cap[bus.KeyToken]; ok {
		if tokenStr, ok := tokenValue.(value.StringValue); ok {
			token = tokenStr.Value()
		} else {
			return s.capError()
		}
	}
	return s.Response(user, token)
}

func helpTest(t *testing.T, user, token string, status uint32) {
	response := func(gotUser, gotToken string) bus.CapabilityMap {
		if user != gotUser {
			panic("not expecting user " + gotUser)
		}
		if token != gotToken {
			panic("not expecting token " + gotToken)
		}
		if status == bus.StateContinue {
			status = bus.StateDone
			return bus.CapabilityMap{
				bus.KeyState:    value.Uint(bus.StateContinue),
				bus.KeyNewToken: value.String(token),
			}
		}
		return bus.CapabilityMap{
			bus.KeyState: value.Uint(status),
		}
	}
	endpoint := NewServer(response)
	defer endpoint.Close()
	err := bus.AuthenticateUser(endpoint, user, token)
	switch status {
	case bus.StateDone:
		if err != nil {
			panic("expecting a success, got " + err.Error())
		}
	case bus.StateError:
		if err == nil {
			panic("expecting an error")
		}
	case bus.StateContinue:
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
	helpTest(t, "userA", "correct passwd", bus.StateDone)
	helpTest(t, "userA", "incorrect passwd", bus.StateError)
	helpTest(t, "userA", "negotiation", bus.StateContinue)
	helpTest(t, "userA", "invalid state", uint32(0xdead))
}

func helpAuthError(t *testing.T, user, token string,
	response func(user, token string) bus.CapabilityMap) {

	endpoint := NewServer(response)
	err := bus.AuthenticateUser(endpoint, user, token)
	if err == nil {
		panic("must fail")
	}
	endpoint.Close()
}

func TestAuthError(t *testing.T) {
	// missing state
	helpAuthError(t, "a", "b", func(user, token string) bus.CapabilityMap {
		if user != "a" {
			panic("unexpected user " + user)
		}
		if token != "b" {
			panic("unexpected token " + token)
		}
		return bus.CapabilityMap{}
	})

	// service not available
	helpAuthError(t, "a", "b", nil)

	// wrong state type
	helpAuthError(t, "a", "b", func(user, token string) bus.CapabilityMap {
		return bus.CapabilityMap{
			bus.KeyState: value.String("not correct"),
		}
	})

	// missing new token
	helpAuthError(t, "a", "b", func(user, token string) bus.CapabilityMap {
		return bus.CapabilityMap{
			bus.KeyState: value.Uint(bus.StateContinue),
		}
	})

	// wrong new token type
	helpAuthError(t, "a", "b", func(user, token string) bus.CapabilityMap {
		return bus.CapabilityMap{
			bus.KeyState:    value.Uint(bus.StateContinue),
			bus.KeyNewToken: value.Uint(12),
		}
	})

	// continue reply missing state
	state := 1
	helpAuthError(t, "a", "b", func(user, token string) bus.CapabilityMap {
		if state == 1 {
			state = 2
			return bus.CapabilityMap{
				bus.KeyState:    value.Uint(bus.StateContinue),
				bus.KeyNewToken: value.String("bb"),
			}
		}
		if user != "a" {
			panic("unexpected user " + user)
		}
		if token != "bb" {
			panic("unexpected token " + token)
		}
		return bus.CapabilityMap{}
	})

	// continue reply wrong state type
	state = 1
	helpAuthError(t, "a", "b", func(user, token string) bus.CapabilityMap {
		if state == 1 {
			state = 2
			return bus.CapabilityMap{
				bus.KeyState:    value.Uint(bus.StateContinue),
				bus.KeyNewToken: value.String("bb"),
			}
		}
		return bus.CapabilityMap{
			bus.KeyState: value.String("bb"),
		}
	})

	// continue reply wrong state value
	state = 1
	helpAuthError(t, "a", "b", func(user, token string) bus.CapabilityMap {
		if state == 1 {
			state = 2
			return bus.CapabilityMap{
				bus.KeyState:    value.Uint(bus.StateContinue),
				bus.KeyNewToken: value.String("bb"),
			}
		}
		return bus.CapabilityMap{
			bus.KeyState: value.Uint(1234),
		}
	})

	// continue reply state continue
	state = 1
	helpAuthError(t, "a", "b", func(user, token string) bus.CapabilityMap {
		if state == 1 {
			state = 2
			return bus.CapabilityMap{
				bus.KeyState:    value.Uint(bus.StateContinue),
				bus.KeyNewToken: value.String("bb"),
			}
		}
		return bus.CapabilityMap{
			bus.KeyState: value.Uint(bus.StateContinue),
		}
	})

	// continue reply state error
	state = 1
	helpAuthError(t, "a", "b", func(user, token string) bus.CapabilityMap {
		if state == 1 {
			state = 2
			return bus.CapabilityMap{
				bus.KeyState:    value.Uint(bus.StateContinue),
				bus.KeyNewToken: value.String("bb"),
			}
		}
		return bus.CapabilityMap{
			bus.KeyState: value.Uint(bus.StateError),
		}
	})

	// continue with call error
	state = 1
	var endpoint net.EndPoint
	response := func(user, token string) bus.CapabilityMap {
		if state == 1 {
			state = 2
			return bus.CapabilityMap{
				bus.KeyState:    value.Uint(bus.StateContinue),
				bus.KeyNewToken: value.String("bb"),
			}
		}
		if token != "bb" {
			panic("expecting token bb")
		}
		return bus.CapabilityMap{}
	}

	endpoint = NewServer(response)
	err := bus.AuthenticateUser(endpoint, "a", "b")
	if err == nil {
		panic("must fail")
	}
	endpoint.Close()
}
