package client

import (
	"github.com/lugu/qiloop/bus"
)

func MakeObject(proxy bus.Proxy) ObjectProxy {
	return &proxyObject{proxy}
}
