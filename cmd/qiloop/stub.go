package main

import (
	"github.com/lugu/qiloop/meta/idl"
	s "github.com/lugu/qiloop/meta/stub"
	"io/ioutil"
	"log"
	"os"
)

func stub(idlFileName, stubFileName, packageName string) {

	file, err := os.Open(idlFileName)
	if err != nil {
		log.Fatalf("failed to create %s: %s", idlFileName, err)
	}
	input, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		log.Fatalf("cannot read %s: %s", idlFileName, err)
	}

	output := os.Stdout
	if stubFileName != "-" {
		output, err = os.Create(stubFileName)
		if err != nil {
			log.Fatalf("failed to create %s: %s", stubFileName, err)
		}
		defer output.Close()
	}

	pkg, err := idl.ParsePackage([]byte(input))
	if err != nil {
		log.Fatalf("failed to parse %s: %s", idlFileName, err)
	}
	if len(pkg.Types) == 0 {
		log.Fatalf("parse error: missing type")
	}
	err = s.GeneratePackage(output, packageName, pkg)
	if err != nil {
		log.Fatalf("failed to generate stub: %s", err)
	}
}
