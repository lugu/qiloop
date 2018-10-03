package main

import (
	"flag"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/meta/stub"
	"io"
	"log"
	"os"
)

func open(filename string) io.WriteCloser {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed to open %s: %s", filename, err)
	}
	return file
}

func main() {

	var directory string
	var goBasePackageName string
	flag.StringVar(&directory, "directory", ".", "directory containing the IDL files")
	flag.StringVar(&goBasePackageName, "package", "", "go base package name")

	dir, err := os.Open(directory)
	if err != nil {
		log.Fatalf("failed to open directory: %s", directory)
	}
	defer dir.Close()
	files, err := dir.Readdirnames(-1)
	if err != nil {
		log.Fatalf("failed to open %s: %s", directory, err)
	}

	proxies := open("implementations.go")
	defer proxies.Close()

	interfaces := open("interfaces.go")
	defer interfaces.Close()

	stubs := open("stubs.go")
	defer stubs.Close()

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("failed to open %s: %s", file, err)
		}
		metas, err := idl.ParseIDL(f)
		f.Close()
		proxy.GenerateProxys(metas, goBasePackageName, proxies)
		proxy.GenerateInterfaces(metas, goBasePackageName, interfaces)
		stub.GenerateStubs(metas, goBasePackageName, stubs)
	}
}
