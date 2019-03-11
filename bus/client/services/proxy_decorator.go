package services

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/client/object"
)

// Proxy returns a proxy to the service described by info.
func (info ServiceInfo) Proxy(sess bus.Session) (object.ObjectProxy, error) {
	proxy, err := sess.Proxy(info.Name, 1)
	if err != nil {
		return nil, fmt.Errorf("cannot connect %s: %s", info.Name, err)
	}
	return object.MakeObject(proxy), nil
}
