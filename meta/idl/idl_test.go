package idl

import (
	"github.com/lugu/qiloop/type/object"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func helpTestGenerate(t *testing.T, idlFileName, serviceName string,
	metaObj object.MetaObject) {

	path := filepath.Join("testdata", idlFileName)
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	idl, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	expected := string(idl)

	var w strings.Builder
	if err := GenerateIDL(&w, serviceName, metaObj); err != nil {
		t.Errorf("failed to parse server: %s", err)
	}
	if w.String() != expected {
		t.Errorf("Got:\n%s\nExpecting:\n%s\n", w.String(), expected)
	}
}

func TestServiceServer(t *testing.T) {
	helpTestGenerate(t, "service0.idl", "Server", object.MetaService0)
}

func TestObject(t *testing.T) {
	helpTestGenerate(t, "object.idl", "Object", object.ObjectMetaObject)
}

func TestServiceDirectory(t *testing.T) {
	path := filepath.Join("testdata", "meta-object.bin")
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	metaObj, err := object.ReadMetaObject(file)
	if err != nil {
		panic(err)
	}
	helpTestGenerate(t, "service1.idl", "ServiceDirectory", metaObj)
}
