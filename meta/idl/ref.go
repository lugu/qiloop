package idl

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/meta/signature"
)

// RefType represents a reference to an IDL interface. It
// implements signature.Type.
type RefType struct {
	Scope Scope
	Name  string
}

// NewRefType is a contructor for the representation of a type reference to be
// resolved with a TypeSet.
func NewRefType(name string, scope Scope) signature.Type {
	return &RefType{scope, name}
}

// Signature returns the signature of the referenced type. If the
// reference can not be resolved, it returns an invalid struct type
// with a name describing the error.
func (r *RefType) Signature() string {
	t, err := r.Scope.Search(r.Name)
	if err == nil {
		return t.Signature()
	}
	return signature.NewStructType(err.Error(), nil).Signature()
}

// SignatureIDL returns the IDL signature of the referenced type. If the
// reference can not be resolved, it returns an invalid name
// describing the error.
func (r *RefType) SignatureIDL() string {
	t, err := r.Scope.Search(r.Name)
	if err == nil {
		return t.SignatureIDL()
	}
	return signature.NewStructType(err.Error(), nil).SignatureIDL()
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (r *RefType) TypeName() *signature.Statement {
	t, err := r.Scope.Search(r.Name)
	if err == nil {
		return t.TypeName()
	}
	return jen.Id(r.Name)
}

// RegisterTo adds the type to the type set. Does nothing on RefType.
func (r *RefType) RegisterTo(set *signature.TypeSet) {
}

// TypeDeclaration writes the type declaration into file.
func (r *RefType) TypeDeclaration(file *jen.File) {
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (r *RefType) Marshal(id string, writer string) *signature.Statement {
	t, err := r.Scope.Search(r.Name)
	if err == nil {
		return t.Marshal(id, writer)
	}
	return jen.Qual("fmt", "Errorf").Call(
		jen.Lit("reference type serialization not implemented: %v"), jen.Id(id),
	)
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (r *RefType) Unmarshal(reader string) *signature.Statement {
	t, err := r.Scope.Search(r.Name)
	if err == nil {
		return t.Unmarshal(reader)
	}
	return jen.Return(
		jen.Nil(),
		jen.Qual("fmt", "Errorf").Call(
			jen.Lit("reference type not implemented"),
		),
	)
}

func (r *RefType) resolve(set *signature.TypeSet) (signature.Type, error) {
	typ := set.Search(r.Name)
	if typ == nil {
		return nil, fmt.Errorf("%s: not found in typeset (size: %d)",
			r.Name, len(set.Types))
	}
	return typ, nil
}

// Reader returns a TypeReader of the referenced type.
func (r *RefType) Reader() signature.TypeReader {
	t, err := r.Scope.Search(r.Name)
	if err == nil {
		reader, err := signature.MakeReader(t.Signature())
		if err != nil {
			return signature.UnknownReader(r.Signature())
		}
		return reader
	}
	return t.Reader()
}
