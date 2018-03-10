package main

import (
	"github.com/lugu/qiloop/meta/signature"
	"io"
	"log"
	"os"
)

func main() {
	var output io.Writer

	if len(os.Args) == 0 {
		log.Fatalf("missing argument: require at least a package name")
	}
	packageName := os.Args[1]

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

	typ, err := signature.Parse(signature.ObjectSignature)
	if err != nil {
		log.Fatalf("parsing error: %s\n", err)
	}
	if err = signature.GenerateType(typ, packageName, output); err != nil {
		log.Fatalf("code generation failed: %s\n", err)
	}
}
