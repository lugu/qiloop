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
	var outputInterfaces io.Writer = os.Stdout
	var outputImplementation io.Writer = os.Stdout

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
		outputInterfaces = file
		defer file.Close()
	}

	if len(os.Args) > 3 {
		filename := os.Args[3]
		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
		}
		outputImplementation = file
		defer file.Close()
	}

	objects := make([]object.MetaObject, 0)
	objects = append(objects, object.MetaService0)
	objects = append(objects, object.ObjectMetaObject)
	directory, err := object.ReadMetaObject(input)
	if err != nil {
		log.Fatalf("failed to parse MetaObject: %s", err)
	}
	directory.Description = "ServiceDirectory"
	objects = append(objects, directory)

	err = proxy.GenerateProxys(objects, "stage2", outputImplementation)
	if err != nil {
		log.Fatalf("proxy generation failed: %s\n", err)
	}
	err = proxy.GenerateInterfaces(objects, "stage2", outputInterfaces)
	if err != nil {
		log.Fatalf("interface generation failed: %s\n", err)
	}
}
