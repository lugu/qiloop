package idl

import (
	. "github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	parsec "github.com/prataprc/goparsec"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func helpParseMethod(t *testing.T, label, input string, expected Method) {
	ctx := NewContext()
	root, _ := method(ctx)(parsec.NewScanner([]byte(input)))
	if root == nil {
		t.Fatalf("%s: cannot parse input:\n%s", label, input)
	}
	if err, ok := root.(error); ok {
		t.Fatalf("%s: parsing error: %v", label, err)
	}
	function, ok := root.(Method)
	if !ok {
		t.Fatalf("%s; type error %+v: %+v", label, reflect.TypeOf(root),
			root)
	}
	if function.Name != expected.Name {
		t.Fatalf("%s: wrong name %s, got %s", label, expected.Name,
			function.Name)
	}
	if function.Return.Signature() != expected.Return.Signature() {
		t.Fatalf("%s: wrong return %s, got %s", label,
			expected.Return.Signature(), function.Return.Signature())
	}
	if len(function.Params) != len(expected.Params) {
		t.Fatalf("%s: wrong number of parameter", label)
	}
	for i, p := range function.Params {
		if p.Type.Signature() != expected.Params[i].Type.Signature() {
			t.Fatalf("%s: wrong param %d: %s, got %s", label, i,
				p.Type.Signature(),
				expected.Params[i].Type.Signature())
		}
	}
}

func helpParseType(t *testing.T, lang, sign string) {
	label := lang
	ctx := NewContext()
	root, _ := ctx.typeParser(parsec.NewScanner([]byte(lang)))
	if root == nil {
		t.Fatalf("%s: cannot parse input:\n%s", label, lang)
	}
	if err, ok := root.(error); ok {
		t.Fatalf("%s: parsing error: %v", label, err)
	}
	if typ, ok := root.(Type); !ok {
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
	expected := Method{
		Name:   "methodName",
		Id:     0,
		Return: NewVoidType(),
		Params: []Parameter{},
	}
	helpParseMethod(t, "TestParseMethod0", input, expected)
}

func TestParseMethod0bis(t *testing.T) {
	input := `fn methodName() //uid:200`
	expected := Method{
		Name:   "methodName",
		Id:     200,
		Return: NewVoidType(),
		Params: []Parameter{},
	}
	helpParseMethod(t, "TestParseMethod0bis", input, expected)
}

func TestParseMethod1(t *testing.T) {
	input := `fn methodName() -> int32`
	expected := Method{
		Name:   "methodName",
		Id:     0,
		Return: NewIntType(),
		Params: []Parameter{},
	}
	helpParseMethod(t, "TestParseMethod1", input, expected)
}

func TestParseMethod1ter(t *testing.T) {
	input := `fn methodName() -> float64`
	expected := Method{
		Name:   "methodName",
		Id:     0,
		Return: NewDoubleType(),
		Params: []Parameter{},
	}
	helpParseMethod(t, "TestParseMethod1ter", input, expected)
}

func TestParseMethod1bis(t *testing.T) {
	input := `fn methodName() -> bool`
	expected := Method{
		Name:   "methodName",
		Id:     0,
		Return: NewBoolType(),
		Params: []Parameter{},
	}
	helpParseMethod(t, "TestParseMethod1bis", input, expected)
}

func TestParseMethod2(t *testing.T) {
	input := `fn methodName(param1: int32, param2: float64)`
	expected := Method{
		Name:   "methodName",
		Id:     0,
		Return: NewVoidType(),
		Params: []Parameter{
			Parameter{
				"param1",
				NewIntType(),
			},
			Parameter{
				"param2",
				NewDoubleType(),
			},
		},
	}
	helpParseMethod(t, "TestParseMethod2", input, expected)
}
func TestParseMethod3(t *testing.T) {
	input := `fn methodName(param1: int32, param2: float64) -> bool`
	expected := Method{
		Name:   "methodName",
		Id:     0,
		Return: NewBoolType(),
		Params: []Parameter{
			Parameter{
				"param1",
				NewIntType(),
			},
			Parameter{
				"param2",
				NewDoubleType(),
			},
		},
	}
	helpParseMethod(t, "TestParseMethod3", input, expected)
}
func TestParseMethod3bis(t *testing.T) {
	input := `fn methodName(param1: int32, param2: float64) -> bool //uid:10`
	expected := Method{
		Name:   "methodName",
		Id:     10,
		Return: NewBoolType(),
		Params: []Parameter{
			Parameter{
				"param1",
				NewIntType(),
			},
			Parameter{
				"param2",
				NewDoubleType(),
			},
		},
	}
	helpParseMethod(t, "TestParseMethod3bis", input, expected)
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

func helpParserTest(t *testing.T, label, idlFileName string, expected *object.MetaObject) {
	path := filepath.Join("testdata", idlFileName)
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	metaList, err := ParseIDL(file)
	if err != nil {
		t.Fatalf("%s: failed to parse idl: %s", label, err)
	}
	if len(metaList) < 1 {
		t.Fatalf("%s: no interface found", label)
	}
	if len(metaList) > 1 {
		t.Fatalf("%s: too many interfaces: %d", label, len(metaList))
	}
	result := metaList[0]
	for i, _ := range result.Methods {
		if result.Methods[i].ParametersSignature != expected.Methods[i].ParametersSignature {
			t.Errorf("%s: wrong parameter signature %d", label, i)
		}
		if result.Methods[i].ReturnSignature != expected.Methods[i].ReturnSignature {
			t.Errorf("%s: wrong return signature %d", label, i)
			t.Errorf("%s: observed: %s", label, result.Methods[i].ReturnSignature)
			t.Errorf("%s: expected: %s", label, expected.Methods[i].ReturnSignature)
		}
	}
	for i, _ := range result.Signals {
		if result.Signals[i].Signature != expected.Signals[i].Signature {
			t.Errorf("%s: wrong signal signature %d", label, i)
		}
	}
}

func TestParseEmptyService(t *testing.T) {
	empty := object.MetaObject{
		Description: "Empty",
	}
	helpParserTest(t, "Empty interface", "empty.idl", &empty)
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

func helpParseInterface(t *testing.T, label, input, name string) {
	ctx := NewContext()
	root, _ := interfaceParser(ctx)(parsec.NewScanner([]byte(input)))
	if root == nil {
		t.Errorf("%s: cannot parse input:\n%s", label, input)
		return
	}
	if err, ok := root.(error); ok {
		t.Errorf("%s: parsing error: %v", label, err)
		return
	}
	if itf, ok := root.(*InterfaceType); !ok {
		t.Errorf("%s; type error %+v: %+v", label, reflect.TypeOf(root), root)
		return
	} else if itf.Name != name {
		t.Errorf("%s; expected :%s: got: %s", label, name, itf.Name)
		return
	}
}

func TestParseInterfaces(t *testing.T) {
	idl := `interface Test1
	end`
	helpParseInterface(t, "1", idl, "Test1")
	idl = `interface Test2
	// this is a comment.
	end`
	helpParseInterface(t, "2", idl, "Test2")
	idl = `interface Test3 // this is a comment.
	end`
	helpParseInterface(t, "3", idl, "Test3")
}

func TestParseReturns(t *testing.T) {
	input := "-> int32"
	expected := NewIntType()
	parser := returns(NewContext())
	root, _ := parser(parsec.NewScanner([]byte(input)))
	if root == nil {
		t.Fatalf("error parsing returns:\n%s", input)
	}
	if err, ok := root.(error); ok {
		t.Fatalf("cannot parse returns: %v", err)
	}
	if ret, ok := root.(Type); !ok {
		t.Fatalf("return type error: %+v", root)
	} else if ret.Signature() != expected.Signature() {
		t.Fatalf("cannot generate signature: %+v", ret)
	}
}

func helpParseStruct(t *testing.T, label, input, expected string) {
	root, _ := structure(NewContext())(parsec.NewScanner([]byte(input)))
	if root == nil {
		t.Errorf("%s: error parsing struuture:\n%s", label, input)
		return
	}
	if err, ok := root.(error); ok {
		t.Errorf("%s: cannot parse returns: %v", label, err)
		return
	}
	if structType, ok := root.(*StructType); !ok {
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
	end`, "()<Test>")
	helpParseStruct(t, "2", `struct Test
	a: int32
	b: str
	end`, "(is)<Test,a,b>")
	helpParseStruct(t, "3", `struct Test
	c: float32 // test
	d: bool
	end`, "(fb)<Test,c,d>")
}

func newDeclaration(t *testing.T) []*StructType {
	scope := NewScope()
	basic1 := NewRefType("basic1", scope)
	basic2 := NewRefType("basic2", scope)
	complex1 := NewRefType("complex1", scope)
	complex2 := NewRefType("complex2", scope)

	structs := []*StructType{
		&StructType{
			Name: "basic1",
			Members: []MemberType{
				MemberType{
					Name: "a",
					Type: NewIntType(),
				},
				MemberType{
					Name: "b",
					Type: NewIntType(),
				},
			},
		},
		&StructType{
			Name: "basic2",
			Members: []MemberType{
				MemberType{
					Name: "a",
					Type: NewIntType(),
				},
				MemberType{
					Name: "b",
					Type: NewIntType(),
				},
			},
		},
		&StructType{
			Name: "complex1",
			Members: []MemberType{
				MemberType{
					Name: "simple1",
					Type: basic1,
				},
				MemberType{
					Name: "b",
					Type: NewIntType(),
				},
			},
		},
		&StructType{
			Name: "complex3",
			Members: []MemberType{
				MemberType{
					Name: "notSimple2",
					Type: complex2,
				},
			},
		},
		&StructType{
			Name: "complex2",
			Members: []MemberType{
				MemberType{
					Name: "notSimple1",
					Type: complex1,
				},
				MemberType{
					Name: "simple2",
					Type: basic2,
				},
				MemberType{
					Name: "b",
					Type: NewIntType(),
				},
			},
		},
	}
	for _, str := range structs {
		scope.Add(str.Name, str)
	}
	return structs
}

func TestSimpleTest(t *testing.T) {
	types := newDeclaration(t)

	expected := []string{
		"(ii)<basic1,a,b>",
		"(ii)<basic2,a,b>",
		"((ii)<basic1,a,b>i)<complex1,simple1,b>",
		"((((ii)<basic1,a,b>i)<complex1,simple1,b>(ii)<basic2,a,b>i)<complex2,notSimple1,simple2,b>)<complex3,notSimple2>",
		"(((ii)<basic1,a,b>i)<complex1,simple1,b>(ii)<basic2,a,b>i)<complex2,notSimple1,simple2,b>",
	}
	for i, s := range types {
		if s.Signature() != expected[i] {
			t.Errorf("%d: invalid: \nobserved: %s, \nexpected: %s",
				i, s.Signature(), expected[i])
		}
	}
}

func TestParseSimpleIDL(t *testing.T) {
	expected := object.MetaObject{
		Description: "I",
		Methods: map[uint32]object.MetaMethod{
			100: object.MetaMethod{
				Uid:                 100,
				ReturnSignature:     "v",
				Name:                "c",
				ParametersSignature: "(((if)<A,a,b>f)<B,a,b>)",
				Parameters: []object.MetaMethodParameter{
					object.MetaMethodParameter{
						Name:        "d",
						Description: "B",
					},
				},
				ReturnDescription: "nothing",
			},
		},
		Signals:    map[uint32]object.MetaSignal{},
		Properties: map[uint32]object.MetaProperty{},
	}
	idl := `
	struct B
		a: A
		b: float32
	end
	struct A
		a: int32
		b: float32
	end
	interface I
		fn c(d: B)
	end
	`
	metaList, err := ParseIDL(strings.NewReader(idl))
	if err != nil {
		t.Fatalf("Failed to parse input: %s", err)
	}
	if len(metaList) != 1 {
		t.Fatalf("Failed to parse I (%d)", len(metaList))
	}
	metaObj := metaList[0]
	if !reflect.DeepEqual(metaObj, expected) {
		t.Fatalf("\nexpected: %#v\nobserved: %#v", expected, metaObj)
	}
}

func TestParseEnum(t *testing.T) {
	input := `
	enum EnumType0
		valueA = 1
		valueB = 2
	end
	`
	declarations, err := ParsePackage([]byte(input))
	if err != nil {
		t.Fatalf("Failed to parse input: %s", err)
	}
	typeList := declarations.Types
	if len(typeList) != 1 {
		t.Fatalf("missing enum (%d)", len(typeList))
	}
	if enum, ok := typeList[0].(*EnumType); ok {
		if enum.Name != "EnumType0" {
			t.Errorf("invalide name %s", enum.Name)
		}
		if val, ok := enum.Values["valueB"]; ok {
			if val != 2 {
				t.Errorf("incorrect value %d", val)
			}
		} else {
			t.Errorf("missing const %s", "valueB")
		}
		if val, ok := enum.Values["valueA"]; ok {
			if val != 1 {
				t.Errorf("incorrect value %d", val)
			}
		} else {
			t.Errorf("missing const %s", "valueA")
		}
	} else {
		t.Fatalf("invalide type %#v", typeList[0])
	}
}

func TestPackageName(t *testing.T) {
	input := `package bla_bla.te-st // comment test `
	declarations, err := ParsePackage([]byte(input))
	if err != nil {
		t.Fatalf("Failed to parse input: %s", err)
	}
	if declarations.Name != "bla_bla.te-st" {
		t.Fatalf("Package name invalid: %s", declarations.Name)
	}
}
