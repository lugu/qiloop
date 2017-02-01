package main

import (
	"flag"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/proxy"
	"log"
	"os"
)

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
	proxies, err := os.Create("implementations.go")
	if err != nil {
		log.Fatalf("failed to create implementations.go: %s", err)
	}
	defer proxies.Close()
	interfaces, err := os.Create("interfaces.go")
	if err != nil {
		log.Fatalf("failed to create interfaces.go: %s", err)
	}
	defer interfaces.Close()
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("failed to open %s: %s", file, err)
		}
		metas, err := idl.ParseIDL(f)
		f.Close()
		proxy.GenerateProxys(metas, goBasePackageName, proxies)
		proxy.GenerateInterfaces(metas, goBasePackageName, interfaces)
	}
}
