package session

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/basic"
	"github.com/lugu/qiloop/object"
)

// Client represents a client connection to a service.
type Client interface {
	Call(serviceID uint32, objectID uint32, actionID uint32, payload []byte) ([]byte, error)
}

type Proxy interface {
	Call(action string, payload []byte) ([]byte, error)
	CallID(action uint32, payload []byte) ([]byte, error)

	// ServiceID returns the service identifier. Allow services to
	// implement the object.Object interface.
	ServiceID() uint32
	// ObjectID returns the object identifier with the service. Allow
	// services to implement the object.Object interface.
	ObjectID() uint32
}

type Session interface {
	Proxy(name string, objectID uint32) (Proxy, error)
}

func MetaObject(client Client, serviceID uint32, objectID uint32) (m object.MetaObject, err error) {
	buf := bytes.NewBuffer(make([]byte, 4))
	basic.WriteUint32(objectID, buf)
	response, err := client.Call(serviceID, objectID, object.MetaObjectMethodID, buf.Bytes())
	if err != nil {
		return m, fmt.Errorf("Can not call MetaObject: %s", err)
	}
	buf = bytes.NewBuffer(response)
	m, err = object.ReadMetaObject(buf)
	if err != nil {
		return m, fmt.Errorf("failed to parse metaObject response: %s", err)
	}
	return m, nil
}
