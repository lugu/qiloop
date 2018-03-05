package main

import (
	"io"
	"log"
	"os"
	"github.com/lugu/qiloop/meta/signature"
)

func main() {
	var output io.Writer

	if len(os.Args) > 1 {
		filename := os.Args[1]

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

	typ, err := signature.Parse(signature.MetaObjectSignature)
	if err != nil {
		log.Fatalf("parsing error: %s\n", err)
	}
	if err = signature.GenerateType(typ, "stage1", output); err != nil {
		log.Fatalf("code generation failed: %s\n", err)
	}
}
