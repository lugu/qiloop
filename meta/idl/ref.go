package idl

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	. "github.com/lugu/qiloop/meta/signature"
)

var types *TypeSet = nil

// RefType represents a struct.
type RefType struct {
	Name string
}

// NewRefType is a contructor for the representation of a type reference to be
// resolved with a TypeSet.
func NewRefType(name string) (Type, error) {
	ref := &RefType{name}
	if types != nil {
		if typ, err := ref.resolve(types); err != nil {
			return nil, fmt.Errorf("unknown reference %s: %s", name, err)
		} else {
			return typ, nil
		}
	}
	return ref, nil
}

func (r *RefType) Signature() string {
	if types == nil {
		return NewStructType(r.Name+"_reference", nil).Signature()
	}
	typ, err := r.resolve(types)
	if err != nil {
		return NewStructType(err.Error(), nil).Signature()
	}
	return typ.Signature()
}

func (r *RefType) SignatureIDL() string {
	if types == nil {
		return NewStructType(r.Name+"_reference", nil).SignatureIDL()
	}
	typ, err := r.resolve(types)
	if err != nil {
		return NewStructType(err.Error(), nil).SignatureIDL()
	}
	return typ.Signature()
}

func (r *RefType) TypeName() *Statement {
	return jen.Id(r.Name)
}

func (r *RefType) RegisterTo(set *TypeSet) {
}

func (r *RefType) TypeDeclaration(file *jen.File) {
}

func (r *RefType) Marshal(id string, writer string) *Statement {
	return jen.Qual("fmt", "Errorf").Call(
		jen.Lit("reference type serialization not implemented: %v"), jen.Id(id),
	)
}

func (r *RefType) Unmarshal(reader string) *Statement {
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
