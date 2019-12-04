package idl_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/proxy"
)

func helpGenerateProxy(t *testing.T, idlFileName, packageName string) {

	path := filepath.Join("testdata", idlFileName)
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	input, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		t.Fatalf("cannot read %s: %s", idlFileName, err)
	}

	output := os.Stdout

	pkg, err := idl.ParsePackage([]byte(input))
	if err != nil {
		t.Fatalf("parse %s: %s", idlFileName, err)
	}

	if err := proxy.GeneratePackage(output, packageName, pkg); err != nil {
		t.Fatalf("generate output: %s", err)
	}
}

func TestProxy(t *testing.T) {
	helpGenerateProxy(t, "object.idl", "object")
	helpGenerateProxy(t, "service0.idl", "zero")
	helpGenerateProxy(t, "service1.idl", "one")
	helpGenerateProxy(t, "property.idl", "prop")
}
