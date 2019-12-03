package main

import (
	"io"
	"log"
	"os"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/type/object"
)

func open(filename string) io.WriteCloser {
	switch filename {
	case "":
		file, err := os.Create(os.DevNull)
		if err != nil {
			log.Fatalf("open %s: %s", os.DevNull, err)
		}
		return file
	case "-":
		return os.Stdout
	default:
		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("open %s: %s", filename, err)
		}
		return file
	}
}

func scan(serverURL, packageName, serviceName, idlFile string) {

	output := open(idlFile)
	defer output.Close()

	sess, err := session.NewSession(serverURL)
	if err != nil {
		log.Fatalf("create session: %s", err)
	}

	// service direction id = 1
	dir, err := services.Services(sess).ServiceDirectory()
	if err != nil {
		log.Fatalf("create proxy: %s", err)
	}

	serviceInfoList, err := dir.Services()
	if err != nil {
		log.Fatalf("list services: %s", err)
	}

	objects := map[string]object.MetaObject{}

	for _, info := range serviceInfoList {

		if serviceName != "" && serviceName != info.Name {
			continue
		}
		proxy, err := sess.Proxy(info.Name, 1)
		if err != nil {
			log.Printf("%s: %s", info.Name, err)
		}
		obj := bus.MakeObject(proxy)
		meta, err := obj.MetaObject(1)
		if err != nil {
			log.Printf("%s: %s", info.Name, err)
		}
		meta.Description = info.Name
		objects[info.Name] = meta
	}

	if packageName == "" {
		packageName = "unknown"
	}
	err = idl.GenerateIDL(output, packageName, objects)
	if err != nil {
		log.Printf("generate IDL: %s", err)
	}
}
