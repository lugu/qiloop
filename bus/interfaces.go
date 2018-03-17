package bus

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
)

// Service represents a service which can answer call and emit signal
// events.
type Service interface {
	Reply(objectID uint32, actionID uint32, payload []byte) ([]byte, error)
	Emit(objectID, signalID uint32, cancel chan int) (chan []byte, error)
}

// Client represents a client connection to a service.
type Client interface {
	Call(serviceID uint32, objectID uint32, actionID uint32, payload []byte) ([]byte, error)
	Subscribe(serviceID, objectID, signalID uint32, cancel chan int) (chan []byte, error)
}

type Proxy interface {
	Call(action string, payload []byte) ([]byte, error)
	CallID(action uint32, payload []byte) ([]byte, error)

	// SignalSubscribe returns a channel with the values of a signal
	Subscribe(signal string, cancel chan int) (chan []byte, error)
	SubscribeID(signal uint32, cancel chan int) (chan []byte, error)

	MethodUid(name string) (uint32, error)
	SignalUid(name string) (uint32, error)

	// ServiceID returns the related service identifier
	ServiceID() uint32
	// ServiceID returns object identifier within the service
	// namespace.
	ObjectID() uint32
}

type Session interface {
	Proxy(name string, objectID uint32) (Proxy, error)
	Object(ref object.ObjectReference) (object.Object, error)
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
