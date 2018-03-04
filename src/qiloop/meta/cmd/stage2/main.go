package main

import (
	"io"
	"log"
	"os"
	"qiloop/meta/proxy"
	"qiloop/meta/signature"
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

	err := proxy.GenerateProxy(signature.MetaService0, "stage2", "Server", output)
	if err != nil {
		log.Fatalf("proxy generation failed: %s\n", err)
	}
	return
}
