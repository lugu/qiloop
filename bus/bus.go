package bus

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
)

// MetaObject queries the meta object of the remote object.
func MetaObject(client Client, serviceID uint32, objectID uint32) (m object.MetaObject, err error) {
	var buf bytes.Buffer
	err = basic.WriteUint32(objectID, &buf)
	if err != nil {
		return m, fmt.Errorf("Can not serialize objectID: %s", err)
	}
	response, err := client.Call(serviceID, objectID, object.MetaObjectMethodID, buf.Bytes())
	if err != nil {
		return m, fmt.Errorf("Can not call MetaObject: %s", err)
	}
	ret := bytes.NewBuffer(response)
	m, err = object.ReadMetaObject(ret)
	if err != nil {
		return m, fmt.Errorf("failed to parse metaObject response: %s", err)
	}
	return m, nil
}
