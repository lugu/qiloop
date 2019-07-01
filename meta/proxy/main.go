package proxy

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/lugu/qiloop/meta/idl"
)

// GenerateProxy writes the service stub from an IDL file
func GenerateProxy(idlFileName, proxyFileName, packageName string) {

	file, err := os.Open(idlFileName)
	if err != nil {
		log.Fatalf("cannot open %s: %s", idlFileName, err)
	}

	input, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		log.Fatalf("cannot read %s: %s", idlFileName, err)
	}

	output := os.Stdout
	if proxyFileName != "-" {
		output, err = os.Create(proxyFileName)
		if err != nil {
			log.Fatalf("failed to create %s: %s", proxyFileName, err)
		}
		defer output.Close()
	}

	pkg, err := idl.ParsePackage([]byte(input))
	if err != nil {
		log.Fatalf("failed to parse %s: %s", idlFileName, err)
	}

	if err := GeneratePackage(output, packageName, pkg); err != nil {
		log.Fatalf("failed to generate output: %s", err)
	}
}
