package bus

import (
	"bytes"
	"fmt"
	"io"

	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
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
	err := basic.WriteUint32(uint32(len(m)), out)
	if err != nil {
		return fmt.Errorf("failed to write map size: %s", err)
	}
	for k, v := range m {
		err = basic.WriteString(k, out)
		if err != nil {
			return fmt.Errorf("failed to write map key: %s", err)
		}
		err = v.Write(out)
		if err != nil {
			return fmt.Errorf("failed to write map value: %s", err)
		}
	}
	return nil
}

// ReadCapabilityMap unmarshals the capability map.
func ReadCapabilityMap(in io.Reader) (m CapabilityMap, err error) {

	size, err := basic.ReadUint32(in)
	if err != nil {
		return m, fmt.Errorf("failed to read map size: %s", err)
	}
	m = make(map[string]value.Value, size)
	for i := 0; i < int(size); i++ {
		k, err := basic.ReadString(in)
		if err != nil {
			return m, fmt.Errorf("failed to read map key: %s", err)
		}
		v, err := value.NewValue(in)
		if err != nil {
			return m, fmt.Errorf("failed to read map value: %s", err)
		}
		m[k] = v
	}
	return m, nil
}

type serviceAuthenticate struct {
	auth Authenticator
}

func (s *serviceAuthenticate) Receive(m *net.Message, from *Channel) error {
	if m.Header.Action != object.AuthenticateActionID {
		return util.ReplyError(from.EndPoint, m, ErrActionNotFound)
	}
	response, err := s.wrapAuthenticate(from, m.Payload)

	if err != nil {
		return util.ReplyError(from.EndPoint, m, err)
	}
	hdr := net.NewHeader(net.Reply, 0, 0, object.AuthenticateActionID,
		m.Header.ID)
	reply := net.NewMessage(hdr, response)
	return from.EndPoint.Send(reply)
}

func (s *serviceAuthenticate) Activate(activation Activation) error {
	return nil
}

func (s *serviceAuthenticate) OnTerminate() {
}

func (s *serviceAuthenticate) wrapAuthenticate(from *Channel, payload []byte) ([]byte, error) {
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

func (s *serviceAuthenticate) Authenticate(from *Channel, cap CapabilityMap) CapabilityMap {
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
		from.Authenticate()
		return from.Cap
	}
	return s.capError()
}

// ServiceAuthenticate represents the servie server (serivce zero)
// used to authenticate a new connection.
func ServiceAuthenticate(auth Authenticator) ServerObject {
	return &serviceAuthenticate{
		auth: auth,
	}
}
