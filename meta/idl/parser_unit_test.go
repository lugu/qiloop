package idl

import (
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	parsec "github.com/prataprc/goparsec"
	"reflect"
	"testing"
)

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
