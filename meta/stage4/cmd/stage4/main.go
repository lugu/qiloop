package main

import (
	"flag"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/type/object"
	"log"
	"os"
	"path"
)

func main() {

	var directory string
	var packageName string
	flag.StringVar(&directory, "directory", ".", "directory with IDL files")
	flag.StringVar(&packageName, "package", "services", "package name")
	flag.Parse()
	metas := make([]object.MetaObject, 0)

	dir, err := os.Open(directory)
	if err != nil {
		log.Fatalf("failed to open directory: %s", directory)
	}
	defer dir.Close()
	files, err := dir.Readdirnames(-1)
	if err != nil {
		log.Fatalf("failed to open %s: %s", directory, err)
	}
	for _, entry := range files {
		file := path.Join(directory, entry)
		f, err := os.Open(file)
		if err != nil {
			log.Printf("failed to open %s: %s", file, err)
			continue
		}
		objects, err := idl.ParseIDL(f)
		f.Close()
		if err != nil {
			log.Printf("failed to parse %s: %s", file, err)
			continue
		}
		for _, meta := range objects {
			metas = append(metas, meta)
		}
	}

	proxies, err := os.Create("proxy_gen.go")
	defer proxies.Close()
	if err != nil {
		log.Fatalf("failed to create proxies.go: %s", err)
	}
	if err := proxy.Generate(metas, packageName, proxies); err != nil {
		log.Printf("failed to generate proxies: %s", err)
	}
}
