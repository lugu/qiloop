package server

import (
	"bytes"
	"fmt"
	. "github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
	"io"
)

type Authenticator struct {
	passwords map[string]string
}

func (a *Authenticator) Authenticate(cap CapabilityMap) CapabilityMap {
	if userValue, ok := cap[KeyUser]; ok {
		if userStr, ok := userValue.(value.StringValue); ok {
			if tokenValue, ok := cap[KeyToken]; ok {
				if tokenStr, ok := tokenValue.(value.StringValue); ok {
					if pwd, ok := a.passwords[string(userStr)]; ok {
						if pwd == string(tokenStr) {
							return CapabilityMap{
								KeyState: value.Uint(StateDone),
							}
						}
					}
				}
			}
		}
	}
	return CapabilityMap{
		KeyState: value.Uint(StateError),
	}
}

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

type ServiceAuthenticate struct {
	auth Authenticator
}

func (s *ServiceAuthenticate) Receive(m *net.Message, from *Context) error {
	if m.Header.Action != object.AuthenticateActionID {
		return ActionNotFound
	}
	response, err := s.wrapAuthenticate(m.Payload)

	if err != nil {
		hdr := net.NewHeader(net.Error, 0, 0,
			object.AuthenticateActionID, m.Header.ID)
		mError := net.NewMessage(hdr, util.ErrorPaylad(err))
		return from.EndPoint.Send(mError)
	}
	hdr := net.NewHeader(net.Reply, 0, 0, object.AuthenticateActionID,
		m.Header.ID)
	reply := net.NewMessage(hdr, response)
	return from.EndPoint.Send(reply)
}

func (s *ServiceAuthenticate) Activate(sess session.Session, serviceID, objectID uint32) {
}

func (s *ServiceAuthenticate) wrapAuthenticate(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	m, err := ReadCapabilityMap(buf)
	if err != nil {
		return nil, err
	}
	ret := s.auth.Authenticate(m)
	buf = bytes.NewBuffer(make([]byte, 0))
	err = WriteCapabilityMap(ret, buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func NewServiceAuthenticate(passwords map[string]string) Object {
	return &ServiceAuthenticate{
		auth: Authenticator{
			passwords,
		},
	}
}
