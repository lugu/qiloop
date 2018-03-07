package main

import (
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/meta/stage1"
	"github.com/lugu/qiloop/meta/stage2"
	"io"
	"log"
	"os"
)

func main() {
	objects := make([]stage1.MetaObject, 2)

	objects[0] = stage2.MetaService0
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

	directory, err := stage1.ReadMetaObject(input)
	if err != nil {
		log.Fatalf("failed to parse MetaObject: %s", err)
	}
	objects[1] = directory
	objects[1].Description = "Directory"

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
