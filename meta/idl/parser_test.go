package idl_test

import (
	. "github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/type/object"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func helpParserTest(t *testing.T, label, idlFileName string, expectedMetaObj *object.MetaObject) {
	path := filepath.Join("testdata", idlFileName)
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	metaList, err := Parse(file)
	if err != nil {
		t.Fatalf("%s: failed to parse idl: %s", label, err)
	}
	if len(metaList) < 1 {
		t.Fatalf("%s: no interface found", label)
	}
	if len(metaList) > 1 {
		t.Fatalf("%s: too many interfaces: %d", label, len(metaList))
	}
	if !reflect.DeepEqual(metaList[0], *expectedMetaObj) {
		t.Fatalf("%s: expected %#v, got %#v", label, expectedMetaObj, metaList[0])
	}
}

func TestParseEmptyService(t *testing.T) {
	helpParserTest(t, "Empty interface", "empty.idl", new(object.MetaObject))
}

func TestParseService0(t *testing.T) {
	helpParserTest(t, "Service 0", "service0.idl", &object.MetaService0)
}

func TestParseService1(t *testing.T) {
	path := filepath.Join("testdata", "meta-object.bin")
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	metaObj, err := object.ReadMetaObject(file)
	helpParserTest(t, "Service 1", "service1.idl", &metaObj)
}

func TestParseObject(t *testing.T) {
	helpParserTest(t, "Object", "object.idl", &object.ObjectMetaObject)
}
