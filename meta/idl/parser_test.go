package idl

import (
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	parsec "github.com/prataprc/goparsec"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func helpParseMethod(t *testing.T, label, input string, expected object.MetaMethod) {
	root, _ := method()(parsec.NewScanner([]byte(input)))
	if root == nil {
		t.Fatalf("%s: cannot parse input:\n%s", label, input)
	}
	if err, ok := root.(error); ok {
		t.Fatalf("%s: parsing error: %v", label, err)
	}
	if method, ok := root.(*object.MetaMethod); !ok {
		t.Fatalf("%s; type error %+v: %+v", label, reflect.TypeOf(root), root)
	} else if !reflect.DeepEqual(*method, expected) {
		t.Fatalf("%s: expected %#v, got %#v", label, expected, *method)
	}
}

func helpParseType(t *testing.T, lang, sign string) {
	label := lang
	root, _ := typeParser()(parsec.NewScanner([]byte(lang)))
	if root == nil {
		t.Fatalf("%s: cannot parse input:\n%s", label, lang)
	}
	if err, ok := root.(error); ok {
		t.Fatalf("%s: parsing error: %v", label, err)
	}
	if typ, ok := root.(signature.Type); !ok {
		t.Fatalf("%s; type error %+v: %+v", label, reflect.TypeOf(root), root)
	} else if typ.Signature() != sign {
		t.Fatalf("%s: expected %#v, got %#v", label, sign, typ.Signature())
	}
}

func TestParseBasicType(t *testing.T) {
	helpParseType(t, "int32", "i")
	helpParseType(t, "uint32", "I")
	helpParseType(t, "int64", "l")
	helpParseType(t, "uint64", "L")
	helpParseType(t, "float32", "f")
	helpParseType(t, "float64", "d")
	helpParseType(t, "bool", "b")
	helpParseType(t, "str", "s")
	helpParseType(t, "any", "m")
}

func TestParseCompoundType(t *testing.T) {
	helpParseType(t, "Map<int32,uint64>", "{iL}")
	helpParseType(t, "Map<str,bool>", "{sb}")
	helpParseType(t, "Map<float32,any>", "{fm}")
}

func TestParseMethod0(t *testing.T) {
	input := `fn methodName()`
	expected := object.MetaMethod{
		Name:                "methodName",
		ReturnSignature:     "v",
		ParametersSignature: "()",
	}
	helpParseMethod(t, "TestParseMethod0", input, expected)
}

func TestParseMethod0bis(t *testing.T) {
	input := `fn methodName() //uid:100`
	expected := object.MetaMethod{
		Uid:                 100,
		Name:                "methodName",
		ReturnSignature:     "v",
		ParametersSignature: "()",
	}
	helpParseMethod(t, "TestParseMethod0bis", input, expected)
}

func TestParseMethod1(t *testing.T) {
	input := `fn methodName() -> int32`
	expected := object.MetaMethod{
		Name:                "methodName",
		ParametersSignature: "()",
		ReturnSignature:     "i",
	}
	helpParseMethod(t, "TestParseMethod1", input, expected)
}

func TestParseMethod1ter(t *testing.T) {
	input := `fn methodName() -> float64`
	expected := object.MetaMethod{
		Name:                "methodName",
		ParametersSignature: "()",
		ReturnSignature:     "d",
	}
	helpParseMethod(t, "TestParseMethod1ter", input, expected)
}

func TestParseMethod1bis(t *testing.T) {
	input := `fn methodName() -> bool`
	expected := object.MetaMethod{
		Name:                "methodName",
		ParametersSignature: "()",
		ReturnSignature:     "b",
	}
	helpParseMethod(t, "TestParseMethod1bis", input, expected)
}

func TestParseMethod2(t *testing.T) {
	input := `fn methodName(param1: int32, param2: float64)`
	expected := object.MetaMethod{
		Name:                "methodName",
		ParametersSignature: "(id)",
		ReturnSignature:     "v",
		Parameters: []object.MetaMethodParameter{
			object.MetaMethodParameter{
				Name: "param1",
			},
			object.MetaMethodParameter{
				Name: "param2",
			},
		},
	}
	helpParseMethod(t, "TestParseMethod2", input, expected)
}
func TestParseMethod3(t *testing.T) {
	input := `fn methodName(param1: int32, param2: float64) -> bool`
	expected := object.MetaMethod{
		Name:                "methodName",
		ParametersSignature: "(id)",
		ReturnSignature:     "b",
		Parameters: []object.MetaMethodParameter{
			object.MetaMethodParameter{
				Name: "param1",
			},
			object.MetaMethodParameter{
				Name: "param2",
			},
		},
	}
	helpParseMethod(t, "TestParseMethod3", input, expected)
}
func TestParseMethod3bis(t *testing.T) {
	input := `fn methodName(param1: int32, param2: float64) -> bool //uid:10`
	expected := object.MetaMethod{
		Uid:                 10,
		Name:                "methodName",
		ParametersSignature: "(id)",
		ReturnSignature:     "b",
		Parameters: []object.MetaMethodParameter{
			object.MetaMethodParameter{
				Name: "param1",
			},
			object.MetaMethodParameter{
				Name: "param2",
			},
		},
	}
	helpParseMethod(t, "TestParseMethod3bis", input, expected)
}

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
		t.Fatalf("%s:\nexpected: %#v\nobserved: %#v", label, *expectedMetaObj, metaList[0])
	}
}

func helpParseComment(t *testing.T, label, input, expected string) {
	root, _ := comments()(parsec.NewScanner([]byte(input)))
	if root == nil {
		t.Errorf("%s: cannot parse input:\n%s", label, input)
		return
	}
	if err, ok := root.(error); ok {
		t.Errorf("%s: parsing error: %v", label, err)
		return
	}
	if comment, ok := root.(string); !ok {
		t.Errorf("%s; type error %+v: %+v", label, reflect.TypeOf(root), root)
		return
	} else if comment != expected {
		t.Errorf(`%s; expected: "%s", got: "%s"`, label, expected, comment)
		return
	}
}

func TestParseComments(t *testing.T) {
	helpParseComment(t, "1", "// test", "test")
	helpParseComment(t, "2", "//test", "test")
	helpParseComment(t, "3", "// this test is long", "this test is long")
	helpParseComment(t, "4", "// this test \n is long", "this test ")
}

func TestParseEmptyService(t *testing.T) {
	empty := object.MetaObject{
		Description: "Empty",
	}
	helpParserTest(t, "Empty interface", "empty.idl", &empty)
}

func helpParseInterface(t *testing.T, label, input, description string) {
	root, _ := interfaceParser()(parsec.NewScanner([]byte(input)))
	if root == nil {
		t.Errorf("%s: cannot parse input:\n%s", label, input)
		return
	}
	if err, ok := root.(error); ok {
		t.Errorf("%s: parsing error: %v", label, err)
		return
	}
	if metaObj, ok := root.(*object.MetaObject); !ok {
		t.Errorf("%s; type error %+v: %+v", label, reflect.TypeOf(root), root)
		return
	} else if metaObj.Description != description {
		t.Errorf("%s; expected :%s: got: %s", label, description, metaObj.Description)
		return
	}
}

func TestParseInterfaces(t *testing.T) {
	idl := `interface Test
	end`
	helpParseInterface(t, "1", idl, "Test")
	idl = `interface Test
	// this is a comment.
	end`
	helpParseInterface(t, "2", idl, "Test")
	idl = `interface Test // this is a comment.
	end`
	helpParseInterface(t, "3", idl, "Test")
}

func TestParseService0(t *testing.T) {
	helpParserTest(t, "Service 0", "service0.idl", &object.MetaService0)
}

func TestParseService1(t *testing.T) {
	t.Skip("not yet implemented")
	path := filepath.Join("testdata", "meta-object.bin")
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	metaObj, err := object.ReadMetaObject(file)
	helpParserTest(t, "Service 1", "service1.idl", &metaObj)
}

func TestParseObject(t *testing.T) {
	t.Skip("not yet implemented")
	helpParserTest(t, "Object", "object.idl", &object.ObjectMetaObject)
}

func TestParseReturns(t *testing.T) {
	input := "-> int32"
	expected := signature.NewIntType()
	parser := returns()
	root, _ := parser(parsec.NewScanner([]byte(input)))
	if root == nil {
		t.Fatalf("error parsing returns:\n%s", input)
	}
	if err, ok := root.(error); ok {
		t.Fatalf("cannot parse returns: %v", err)
	}
	if ret, ok := root.(signature.Type); !ok {
		t.Fatalf("return type error: %+v", root)
	} else if ret.Signature() != expected.Signature() {
		t.Fatalf("cannot generate signature: %+v", ret)
	}
}

func helpParseStruct(t *testing.T, label, input, expected string) {
	root, _ := structure()(parsec.NewScanner([]byte(input)))
	if root == nil {
		t.Errorf("%s: error parsing struuture:\n%s", label, input)
		return
	}
	if err, ok := root.(error); ok {
		t.Errorf("%s: cannot parse returns: %v", label, err)
		return
	}
	if structType, ok := root.(*signature.StructType); !ok {
		t.Errorf("%s: return type error: %+v", label, root)
		return
	} else if structType.Signature() != expected {
		t.Errorf("%s: invalid signature: expected: %s, got %s",
			label, expected, structType.Signature())
		return
	}
}

func TestStructureParser(t *testing.T) {
	helpParseStruct(t, "1", `struct Test
	end`, "()<Test,>")
	helpParseStruct(t, "2", `struct Test
	a: int32
	b: str
	end`, "(is)<Test,a,b>")
	helpParseStruct(t, "3", `struct Test
	c: float32 // test
	d: bool
	end`, "(fb)<Test,c,d>")
}
