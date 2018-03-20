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
	ref  Type
}

// NewRefType is a contructor for the representation of a type reference to be
// resolved with a TypeSet.
func NewRefType(name string, ref Type) (*RefType, error) {
	refType := &RefType{name, ref}
	if ref == nil && types != nil {
		if err := refType.resolve(types); err != nil {
			return nil, fmt.Errorf("unknown reference %s: %s", name, err)
		}
	}
	return refType, nil
}

func (r *RefType) Signature() string {
	if r.ref == nil {
		return NewUnknownType().Signature()
	}
	return r.ref.Signature()
}

func (r *RefType) SignatureIDL() string {
	if r.ref == nil {
		return NewUnknownType().SignatureIDL()
	}
	return r.ref.SignatureIDL()
}

func (r *RefType) TypeName() *Statement {
	return jen.Id(r.Name)
}

func (r *RefType) RegisterTo(set *TypeSet) {
	if r.ref != nil {
		r.ref.RegisterTo(set)
	}
}

func (r *RefType) TypeDeclaration(file *jen.File) {
	if r.ref != nil {
		r.ref.TypeDeclaration(file)
	}
}

func (r *RefType) Marshal(id string, writer string) *Statement {
	if r.ref != nil {
		return r.ref.Marshal(id, writer)
	}
	return jen.Qual("fmt", "Errorf").Call(
		jen.Lit("reference type serialization not implemented: %v"), jen.Id(id),
	)
}

func (r *RefType) Unmarshal(reader string) *Statement {
	if r.ref != nil {
		return r.ref.Unmarshal(reader)
	}
	return jen.Return(
		jen.Nil(),
		jen.Qual("fmt", "Errorf").Call(jen.Lit("reference type not implemented")),
	)
}

func (r *RefType) resolve(set *TypeSet) error {
	if sig, ok := set.Signatures[r.Name]; ok {
		for _, typ := range set.Types {
			if typ.Signature() == sig {
				fmt.Printf("resolved: %s => %s\n", r.Name, typ.Signature())
				r.ref = typ
				return nil
			}
		}
		for _, typ := range set.Types {
			fmt.Printf("signature %s\n", typ.Signature())
		}
		return fmt.Errorf("%s: failed to find signature %s", r.Name, sig)
	}
	return fmt.Errorf("%s: ref not found in typeset (size: %d).", r.Name, len(set.Types))
}

// Register associates the StructType of declarations into a global
// variable which will be used to resolve RefType during the second
// pass.
//
// TODO: register MetaObject as well in order to allow interfaces to
// refer to other interfaces.
func registerTypeNames(declarations *Declarations) {
	types = NewTypeSet()
	for _, struc := range declarations.Struct {
		struc.RegisterTo(types)
	}
	for i, _ := range types.Types {
		fmt.Printf("register 3: %s\n", types.Types[i].Signature())
	}
	for k, v := range types.Signatures {
		fmt.Printf("register 4: %s => %s\n", k, v)
	}
}

// Unregster erase all the previously registered StrucType.
func unregsterTypeNames() {
	types = nil
}
