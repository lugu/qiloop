// file generated. DO NOT EDIT.
package client

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	basic "github.com/lugu/qiloop/type/basic"
	value "github.com/lugu/qiloop/type/value"
)

type NewServices struct {
	session bus.Session
}

func Services(s bus.Session) NewServices {
	return NewServices{session: s}
}

type Server interface {
	bus.Proxy
	Authenticate(P0 map[string]value.Value) (map[string]value.Value, error)
}
type ServerProxy struct {
	bus.Proxy
}

func NewServer(ses bus.Session, obj uint32) (Server, error) {
	proxy, err := ses.Proxy("Server", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &ServerProxy{proxy}, nil
}
func (s NewServices) Server() (Server, error) {
	return NewServer(s.session, 1)
}
func (p *ServerProxy) Authenticate(P0 map[string]value.Value) (map[string]value.Value, error) {
	var err error
	var ret map[string]value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = func() error {
		err := basic.WriteUint32(uint32(len(P0)), buf)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range P0 {
			err = basic.WriteString(k, buf)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = v.Write(buf)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("authenticate", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call authenticate failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (m map[string]value.Value, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]value.Value, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := value.NewValue(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse authenticate response: %s", err)
	}
	return ret, nil
}
