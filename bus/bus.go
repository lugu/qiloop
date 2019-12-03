package bus

import (
	"bytes"
	"fmt"

	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
)

// GetMetaObject queries the meta object of the remote object.
func GetMetaObject(client Client, serviceID uint32, objectID uint32) (m object.MetaObject, err error) {
	var buf bytes.Buffer
	err = basic.WriteUint32(objectID, &buf)
	if err != nil {
		return m, fmt.Errorf("Can not serialize objectID: %s", err)
	}
	response, err := client.Call(nil, serviceID, objectID, object.MetaObjectMethodID, buf.Bytes())
	if err != nil {
		return m, fmt.Errorf("Can not call MetaObject: %s", err)
	}
	ret := bytes.NewBuffer(response)
	m, err = object.ReadMetaObject(ret)
	if err != nil {
		return m, fmt.Errorf("parse metaObject response: %s", err)
	}
	return m, nil
}

// MakeObject converts a Proxy into an ObjectProxy.
func MakeObject(proxy Proxy) ObjectProxy {
	return &proxyObject{proxy}
}

// ServiceServer retruns a proxy to the authenticating service (ID 0)
func (s Constructor) ServiceServer() (ServiceZeroProxy, error) {
	proxy, err := s.session.Proxy("ServiceZero", 0)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}
	return &proxyServiceZero{proxy}, nil
}
