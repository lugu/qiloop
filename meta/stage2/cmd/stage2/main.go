package main

import (
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/object"
	"io"
	"log"
	"os"
)

func main() {
	objects := make([]object.MetaObject, 2)

	objects[0] = object.MetaService0
	objects[0].Description = "Server"

	var input io.Reader
	if len(os.Args) > 1 {
		filename := os.Args[1]

		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
			return
		}
		input = file
		defer file.Close()
	} else {
		input = os.Stdin
	}

	directory, err := object.ReadMetaObject(input)
	if err != nil {
		log.Fatalf("failed to parse MetaObject: %s", err)
	}
	objects[1] = directory
	objects[1].Description = "ServiceDirectory"

	var output io.Writer
	if len(os.Args) > 2 {
		filename := os.Args[2]

		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
			return
		}
		output = file
		defer file.Close()
	} else {
		output = os.Stdout
	}

	err = proxy.GenerateProxys(objects, "stage2", output)
	if err != nil {
		log.Fatalf("proxy generation failed: %s\n", err)
	}

	return
}
