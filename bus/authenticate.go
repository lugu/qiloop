package bus

import (
	"bytes"
	"errors"
	"io"

	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
)

const (
	capabilityMapSizeMax = 4096
)

var (
	// ErrCapabilityTooLong is returned when a capability map.
	ErrCapabilityTooLong = errors.New("capability map too long")
)

// Authenticator decides if a user/token tuple is valid. It is used to
// construct the service server (i.e. service zero).
type Authenticator interface {
	Authenticate(user, token string) bool
}

// Dictionary is an Authenticator which reads its permission from a
// dictionnary.
func Dictionary(passwords map[string]string) Authenticator {
	return dictionary(passwords)
}

// Yes is an Authenticator which accepts anything.
type Yes struct{}

// Authenticate returns true
func (y Yes) Authenticate(user, token string) bool {
	return true
}

// No is an Authenticator which refuses anything.
type No struct{}

// Authenticate returns false
func (n No) Authenticate(user, token string) bool {
	return false
}

type dictionary map[string]string

func (d dictionary) Authenticate(user, token string) bool {
	pwd, ok := d[user]
	return ok && pwd == token
}

// WriteCapabilityMap marshals the capability map.
func WriteCapabilityMap(m CapabilityMap, out io.Writer) error {
	return write_capabilityMap(_capabilityMap{
		Properties: m,
	}, out)
}

// ReadCapabilityMap unmarshals the capability map.
func ReadCapabilityMap(in io.Reader) (m CapabilityMap, err error) {
	caps, err := read_capabilityMap(in)
	return caps.Properties, err
}

type serviceAuthenticate struct {
	auth Authenticator
}

func (s *serviceAuthenticate) Receive(m *net.Message, from Channel) error {
	if m.Header.Action != object.AuthenticateActionID {
		return from.SendError(m, ErrActionNotFound)
	}
	response, err := s.wrapAuthenticate(from, m.Payload)
	if err != nil {
		return from.SendError(m, err)
	}
	return from.SendReply(m, response)
}

func (s *serviceAuthenticate) Activate(activation Activation) error {
	return nil
}

func (s *serviceAuthenticate) OnTerminate() {
}

func (s *serviceAuthenticate) wrapAuthenticate(from Channel, payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	m, err := ReadCapabilityMap(buf)
	if err != nil {
		return nil, err
	}
	ret := s.Authenticate(from, m)
	var out bytes.Buffer
	err = WriteCapabilityMap(ret, &out)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (s *serviceAuthenticate) capError() CapabilityMap {
	return CapabilityMap{
		KeyState: value.Uint(StateError),
	}
}

func (s *serviceAuthenticate) Authenticate(from Channel, cap CapabilityMap) CapabilityMap {
	var user, token string
	if userValue, ok := cap[KeyUser]; ok {
		if userStr, ok := userValue.(value.StringValue); ok {
			user = userStr.Value()
		} else {
			return s.capError()
		}
	}
	if tokenValue, ok := cap[KeyToken]; ok {
		if tokenStr, ok := tokenValue.(value.StringValue); ok {
			token = tokenStr.Value()
		} else {
			return s.capError()
		}
	}
	if s.auth.Authenticate(user, token) {
		from.SetAuthenticated()
		return from.Cap()
	}
	return s.capError()
}

// ServiceAuthenticate represents the servie server (serivce zero)
// used to authenticate a new connection.
func ServiceAuthenticate(auth Authenticator) Actor {
	return &serviceAuthenticate{
		auth: auth,
	}
}
