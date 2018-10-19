package idl

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	. "github.com/lugu/qiloop/meta/signature"
)

var types *TypeSet = nil

// RefType represents a struct.
type RefType struct {
	Scope Scope
	Name  string
}

// NewRefType is a contructor for the representation of a type reference to be
// resolved with a TypeSet.
func NewRefType(name string, scope Scope) Type {
	return &RefType{scope, name}
}

func (r *RefType) Signature() string {
	t, err := r.Scope.Search(r.Name)
	if err == nil {
		return t.Signature()
	}
	return NewStructType(err.Error(), nil).Signature()
}

func (r *RefType) SignatureIDL() string {
	t, err := r.Scope.Search(r.Name)
	if err == nil {
		return t.SignatureIDL()
	}
	return NewStructType(err.Error(), nil).SignatureIDL()
}

func (r *RefType) TypeName() *Statement {
	t, err := r.Scope.Search(r.Name)
	if err == nil {
		return t.TypeName()
	}
	return jen.Id(r.Name)
}

func (r *RefType) RegisterTo(set *TypeSet) {
}

func (r *RefType) TypeDeclaration(file *jen.File) {
}

func (r *RefType) Marshal(id string, writer string) *Statement {
	t, err := r.Scope.Search(r.Name)
	if err == nil {
		return t.Marshal(id, writer)
	}
	return jen.Qual("fmt", "Errorf").Call(
		jen.Lit("reference type serialization not implemented: %v"), jen.Id(id),
	)
}

func (r *RefType) Unmarshal(reader string) *Statement {
	t, err := r.Scope.Search(r.Name)
	if err == nil {
		return t.Unmarshal(reader)
	}
	return jen.Return(
		jen.Nil(),
		jen.Qual("fmt", "Errorf").Call(jen.Lit("reference type not implemented")),
	)
}

func (r *RefType) resolve(set *TypeSet) (Type, error) {
	typ := set.Search(r.Name)
	if typ == nil {
		return nil, fmt.Errorf("%s: not found in typeset (size: %d).", r.Name, len(set.Types))
	}
	return typ, nil
}

// Register associates the StructType of declarations into a global
// variable which will be used to resolve RefType during the second
// pass.
//
// TODO: register MetaObject as well in order to allow interfaces to
// refer to other interfaces.
func registerTypeNames(declarations *Declarations) {
	set := NewTypeSet()
	for _, struc := range declarations.Struct {
		struc.RegisterTo(set)
	}
	types = set
}

// Unregster erase all the previously registered StrucType.
func unregsterTypeNames() {
	types = nil
}

// EnumType represents a const.
type EnumType struct {
	Name   string
	Values map[string]int
}

type EnumMember struct {
	Const string
	Value int
}

func NewEnumType(name string, values map[string]int) Type {
	return &EnumType{
		Name:   name,
		Values: values,
	}
}

func (e *EnumType) Signature() string {
	return NewIntType().Signature()
}

func (e *EnumType) SignatureIDL() string {
	return NewIntType().Signature()
}

func (e *EnumType) TypeName() *Statement {
	return jen.Id(e.Name)
}

func (e *EnumType) RegisterTo(set *TypeSet) {
	// TODO: implement me
}

func (e *EnumType) TypeDeclaration(file *jen.File) {
	// TODO: implement me
	file.Id(`type ` + e.Name + ` int
	const (
		XX = YY
	)`)
}

func (e *EnumType) Marshal(id string, writer string) *Statement {
	return NewIntType().Marshal(id, writer)
}

func (e *EnumType) Unmarshal(reader string) *Statement {
	return NewIntType().Unmarshal(reader)
}
