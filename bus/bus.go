package bus

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
)

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
