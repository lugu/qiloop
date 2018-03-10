package main

import (
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/object"
	"io"
	"log"
	"os"
)

func main() {

	var err error
	var input io.Reader = os.Stdin
	var output io.Writer = os.Stdout

	if len(os.Args) > 1 {
		filename := os.Args[1]
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
		}
		input = file
		defer file.Close()
	}

	if len(os.Args) > 2 {
		filename := os.Args[2]
		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
		}
		output = file
		defer file.Close()
	}

	objects := make([]object.MetaObject, 3)
	objects[0] = object.MetaService0
	objects[1] = object.ObjectMetaObject
	objects[2], err = object.ReadMetaObject(input)
	if err != nil {
		log.Fatalf("failed to parse MetaObject: %s", err)
	}
	objects[2].Description = "ServiceDirectory"

	err = proxy.GenerateProxys(objects, "stage2", output)
	if err != nil {
		log.Fatalf("proxy generation failed: %s\n", err)
	}
}
