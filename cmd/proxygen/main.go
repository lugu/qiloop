package main

import (
	"flag"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/type/object"
	"log"
	"os"
)

func main() {
	var idlFile = flag.String("idl", "", "idl file")
	var packageName = flag.String("package", "", "package name")
	var proxyFile = flag.String("proxy", "proxy_gen.go", "file to generate")
	flag.Parse()

	metas := make([]object.MetaObject, 0)

	f, err := os.Open(*idlFile)
	if err != nil {
		log.Fatalf("failed to open %s: %s", *idlFile, err)
	}
	objects, err := idl.ParseIDL(f)
	f.Close()
	if err != nil {
		log.Fatalf("failed to parse %s: %s", *idlFile, err)
	}
	for _, meta := range objects {
		metas = append(metas, meta)
	}

	proxies, err := os.Create(*proxyFile)
	defer proxies.Close()
	if err != nil {
		log.Fatalf("failed to create proxies.go: %s", err)
	}
	if err := proxy.Generate(metas, *packageName, proxies); err != nil {
		log.Fatalf("failed to generate proxies: %s", err)
	}
}
