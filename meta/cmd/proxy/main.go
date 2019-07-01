package main

import (
	"flag"
	"log"

	"github.com/lugu/qiloop/meta/proxy"
)

func main() {
	var idlFileName = flag.String("idl", "", "IDL file name")
	var proxyFileName = flag.String("output", "", "generated proxy file")
	var packageName = flag.String("path", "", "package name")
	flag.Parse()

	if *idlFileName == "" {
		log.Fatalf("missing idl file")
	}
	if *proxyFileName == "" {
		log.Fatalf("missing proxy file")
	}
	proxy.GenerateProxy(*idlFileName, *proxyFileName, *packageName)
}
