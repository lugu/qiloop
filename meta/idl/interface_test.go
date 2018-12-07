package idl

import (
	"github.com/lugu/qiloop/meta/signature"
	parsec "github.com/prataprc/goparsec"
	"reflect"
	"testing"
)

func newInterfaceType(t *testing.T, name string) *InterfaceType {
	input := `interface ` + name + `
	  	fn method1(a: int32)
		prop property1(b: str)
		sig signal1(c: bool)
	end`
	ctx := NewContext()
	root, _ := interfaceParser(ctx)(parsec.NewScanner([]byte(input)))
	if root == nil {
		t.Errorf("cannot parse input:\n%s", input)
	}
	if err, ok := root.(error); ok {
		t.Errorf("parsing error: %v", err)
	}
	itf, ok := root.(*InterfaceType)
	if !ok {
		t.Errorf("type error %+v: %+v", reflect.TypeOf(root), root)
	}
	return itf
}

func TestInterface(t *testing.T) {
	itf := newInterfaceType(t, "Test")
	if itf.Name != "Test" {
		t.Errorf("expected Test: got: %s", itf.Name)
		return
	}
	if itf.Signature() != "o" {
		t.Errorf("unexpected signature: %s", itf.Signature())
	}
	if itf.SignatureIDL() != "obj" {
		t.Errorf("unexpected signature: %s", itf.SignatureIDL())
	}
	if itf.TypeName() == nil {
		t.Error("unexpected nil type name")
	}
	if len(itf.Methods) != 1 {
		t.Errorf("shall have only one method: %d found",
			len(itf.Methods))
	}
	if len(itf.Signals) != 1 {
		t.Errorf("shall have only one signal: %d found",
			len(itf.Signals))
	}
	if len(itf.Properties) != 1 {
		t.Errorf("shall have only one properties: %d found",
			len(itf.Properties))
	}
	set := signature.NewTypeSet()
	itf.RegisterTo(set)
	for _, method := range itf.Methods {
		method.Tuple()
		method.Meta(0)
	}
	for _, signal := range itf.Signals {
		signal.Tuple()
		signal.Meta(0)
	}
	for _, property := range itf.Properties {
		property.Meta(0)
	}
}

func TestTypeRef(t *testing.T) {
	scope := NewScope()
	scope.Add("Test1", signature.NewBoolType())

	ref := NewRefType("Test1", scope)
	if ref.Signature() != "b" {
		t.Errorf("unexpected signature: %s", ref.Signature())
	}
	if ref.SignatureIDL() != "bool" {
		t.Errorf("unexpected signature: %s", ref.SignatureIDL())
	}
	if ref.TypeName() == nil {
		t.Error("unexpected nil type name")
	}
	set := signature.NewTypeSet()
	ref.RegisterTo(set)

	ref.Marshal("a", "b")
	ref.Unmarshal("a")

	ref = NewRefType("Unknown", scope)
	ref.Signature()
	ref.SignatureIDL()
	ref.TypeName()
	ref.Marshal("a", "b")
	ref.Unmarshal("a")
}
