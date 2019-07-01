package main

import (
	p "github.com/lugu/qiloop/meta/proxy"
)

// GenerateProxy write the specialized proxy from an IDL file
func proxy(idlFileName, proxyFileName, packageName string) {
	p.GenerateProxy(idlFileName, proxyFileName, packageName)
}
