package session

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/object"
)

// Client represents a client connection to a service.
type Client interface {
	Call(serviceID uint32, objectID uint32, actionID uint32, payload []byte) ([]byte, error)
}

type Proxy interface {
	Call(action string, payload []byte) ([]byte, error)
	CallID(action uint32, payload []byte) ([]byte, error)
}

type Session interface {
	Proxy(name string, objectID uint32) (Proxy, error)
}

func MetaObject(client Client, serviceID uint32, objectID uint32) (m object.MetaObject, err error) {
	response, err := client.Call(serviceID, objectID, object.MetaObjectMethodID, nil)
	if err != nil {
		return m, fmt.Errorf("Can not call MetaObject: %s", err)
	}
	buf := bytes.NewBuffer(response)
	m, err = object.ReadMetaObject(buf)
	if err != nil {
		return m, fmt.Errorf("failed to parse metaObject response: %s", err)
	}
	return m, nil
}
