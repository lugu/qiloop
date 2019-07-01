package main

import (
	"flag"
	"log"

	"github.com/lugu/qiloop/meta/stub"
)

func main() {
	var idlFileName = flag.String("idl", "", "IDL file name")
	var stubFileName = flag.String("output", "", "generated stub file")
	var packageName = flag.String("path", "", "package name")
	flag.Parse()

	if *idlFileName == "" {
		log.Fatalf("missing idl file")
	}
	if *stubFileName == "" {
		log.Fatalf("missing stub file")
	}
	stub.GenerateStub(*idlFileName, *stubFileName, *packageName)
}
