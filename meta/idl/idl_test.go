package idl

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lugu/qiloop/type/object"
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
	if err := GenerateIDL(&w, "test", map[string]object.MetaObject{
		serviceName: metaObj,
	}); err != nil {
		t.Errorf("parse server: %s", err)
	}
	if w.String() != expected {
		t.Errorf("Got:\n->%s<-\nExpecting:\n->%s<-\n", w.String(), expected)
	}
}

func TestServiceServer(t *testing.T) {
	helpTestGenerate(t, "service0.idl", "Server", object.MetaService0)
}

func TestObject(t *testing.T) {
	helpTestGenerate(t, "object.idl", "Object", object.ObjectMetaObject)
}

func TestProperties(t *testing.T) {
	helpTestGenerate(t, "property.idl", "Test",
		object.MetaObject{
			Description: "Test",
			Methods:     map[uint32]object.MetaMethod{},
			Signals:     map[uint32]object.MetaSignal{},
			Properties: map[uint32]object.MetaProperty{
				0x64: {
					Uid:       0x64,
					Signature: "I",
					Name:      "a",
				},
				0x65: {
					Uid:       0x65,
					Signature: "s",
					Name:      "b",
				},
			},
		})
}

func TestNames(t *testing.T) {
	params := []string{
		"", "for", "len", "func", "range", "some happy name ", "c",
	}
	expected := `package test
interface a
	fn b(P0: str,for_1: str,len: str,func_3: str,range_4: str,somehappyname: str,c: str) -> str //uid:1
end
`
	f := func(names []string) []object.MetaMethodParameter {

		params := []object.MetaMethodParameter{}
		for _, name := range names {
			params = append(params, object.MetaMethodParameter{
				Name:        name,
				Description: name,
			})
		}
		return params
	}
	var w strings.Builder
	if err := GenerateIDL(&w, "test", map[string]object.MetaObject{
		"a": object.MetaObject{
			Description: "Test",
			Methods: map[uint32]object.MetaMethod{
				1: object.MetaMethod{
					Uid:                 1,
					Name:                "b",
					ParametersSignature: "(sssssss)",
					Parameters:          f(params),
					ReturnSignature:     "s",
				},
			},
		},
	}); err != nil {
		t.Errorf("parse server: %s", err)
	}
	if w.String() != expected {
		t.Errorf("Got:\n->%s<-\nExpecting:\n->%s<-\n", w.String(), expected)
	}
}

func TestTupleRet(t *testing.T) {
	expected := `package test
interface a
	fn b() -> Tuple<str,str> //uid:1
end
`
	var w strings.Builder
	if err := GenerateIDL(&w, "test", map[string]object.MetaObject{
		"a": object.MetaObject{
			Description: "Test",
			Methods: map[uint32]object.MetaMethod{
				1: object.MetaMethod{
					Uid:                 1,
					Name:                "b",
					ParametersSignature: "()",
					ReturnSignature:     "(ss)",
				},
			},
		},
	}); err != nil {
		t.Errorf("parse server: %s", err)
	}
	if w.String() != expected {
		t.Errorf("Got:\n->%s<-\nExpecting:\n->%s<-\n", w.String(), expected)
	}
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
