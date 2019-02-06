package main

import (
	"flag"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/proxy"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var filename = flag.String("idl", "", "IDL file")
	var out = flag.String("output", "-", "go file to generate")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("cannot open %s: %s", *filename, err)
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

	output := os.Stdout
	if *out != "-" {
		output, err = os.Create(*out)
		if err != nil {
			log.Fatalf("failed to create %s: %s", *out, err)
		}
		defer output.Close()
	}

	if err := proxy.GeneratePackage(output, pkg); err != nil {
		log.Fatalf("failed to generate output: %s", err)
	}
}
