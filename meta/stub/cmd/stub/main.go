package main

import (
	"flag"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/stub"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var filename = flag.String("idl", "", "IDL file")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("failed to create %s: %s", *filename, err)
	}
	input, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		log.Fatalf("cannot read %s: %s", *filename, err)
	}

	pkg, err := idl.ParsePackage([]byte(input))
	if err != nil {
		log.Fatalf("failed to parse %s: %s", *filename, err)
	}
	if len(pkg.Types) == 0 {
		log.Fatalf("parse error: missing type")
	}
	err = stub.GeneratePackage(os.Stdout, pkg)
	if err != nil {
		log.Fatalf("failed to generate stub: %s", err)
	}
}
