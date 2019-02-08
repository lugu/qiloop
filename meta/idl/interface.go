package idl

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
)

// Proxy were generated from MetaObjects. This is convinient from a
// boostraping point of view. Since the data structures are well
// known, we introduce the InterfaceType which contains the
// information of the MetaObject with the possibility to resolve
// object references.

// Two stage parsing of the IDL:
// 1. Construct a set of Types associated with context.
// 2. Defer the resolution of the types to the proxy/stub generation.

// Namespace represents a set of packages extracted from IDL files.
// Each package is given a name and a set of types.
type Namespace map[string]Scope

// Parameter represents a method parameter. It is used to describe a
// method.
type Parameter struct {
	Name string
	Type signature.Type
}

// Method represents the signature of a method describe in an IDL
// file.
type Method struct {
	Name   string
	ID     uint32
	Return signature.Type
	Params []Parameter
}

// Meta translate the method signature into a MetaMethod use in a
// MetaObject. There is not a one-to-one mapping between the two
// structures: Method capture the reference to other interface
// (object) while MetaMethod treats all object with a generic
// reference.
func (m Method) Meta(id uint32) object.MetaMethod {
	var meta object.MetaMethod
	meta.Uid = id
	meta.Name = m.Name
	meta.ReturnSignature = m.Return.Signature()
	meta.ReturnDescription = m.Return.SignatureIDL()
	params := make([]signature.Type, 0)
	meta.Parameters = make([]object.MetaMethodParameter, 0)
	for _, p := range m.Params {
		var param object.MetaMethodParameter
		param.Name = p.Name
		param.Description = p.Type.SignatureIDL()
		meta.Parameters = append(meta.Parameters, param)
		params = append(params, p.Type)
	}
	meta.ParametersSignature = signature.NewTupleType(params).Signature()
	return meta
}

// Tuple returns a TupleType used to generate marshall/unmarshall
// operations.
func (m Method) Tuple() *signature.TupleType {
	var tuple signature.TupleType
	tuple.Members = make([]signature.MemberType, 0)
	for _, p := range m.Params {
		tuple.Members = append(tuple.Members,
			signature.MemberType{
				Name: p.Name,
				Type: p.Type,
			})
	}
	return &tuple
}

func (m Method) Type() signature.Type {
	if len(m.Params) == 1 {
		return m.Params[0].Type
	}
	return m.Tuple()
}

// Signal represent an interface signal
type Signal struct {
	Name   string
	ID     uint32
	Params []Parameter
}

// Tuple returns a TupleType used to generate marshall/unmarshall
// operations.
func (s Signal) Tuple() *signature.TupleType {
	var tuple signature.TupleType
	tuple.Members = make([]signature.MemberType, 0)
	for _, p := range s.Params {
		tuple.Members = append(tuple.Members,
			signature.MemberType{
				Name: p.Name,
				Type: p.Type,
			})
	}
	return &tuple
}

// Type returns a StructType used to generate marshall/unmarshall
// operations.
func (s Signal) Type() signature.Type {
	if len(s.Params) == 1 {
		return s.Params[0].Type
	}
	return signature.NewStructType(s.Name, s.Tuple().Members)
}

// Meta returns a MetaSignal.
func (s Signal) Meta(id uint32) object.MetaSignal {
	var meta object.MetaSignal
	meta.Uid = id
	meta.Name = s.Name
	types := make([]signature.Type, 0)
	for _, p := range s.Params {
		types = append(types, p.Type)
	}
	meta.Signature = signature.NewTupleType(types).Signature()
	return meta
}

// Property represents a property
type Property struct {
	Name   string
	ID     uint32
	Params []Parameter
}

// Tuple returns a TupleType used to generate marshall/unmarshall
// operations.
func (s Property) Tuple() *signature.TupleType {
	var tuple signature.TupleType
	tuple.Members = make([]signature.MemberType, 0)
	for _, p := range s.Params {
		tuple.Members = append(tuple.Members,
			signature.MemberType{
				Name: p.Name,
				Type: p.Type,
			})
	}
	return &tuple
}

// Type returns a StructType used to generate marshall/unmarshall
// operations.
func (p Property) Type() signature.Type {
	if len(p.Params) == 1 {
		return p.Params[0].Type
	}
	return signature.NewStructType(p.Name, p.Tuple().Members)
}

// Meta converts a property to a MetaProperty.
func (p Property) Meta(id uint32) object.MetaProperty {
	var meta object.MetaProperty
	meta.Uid = id
	meta.Name = p.Name
	types := make([]signature.Type, 0)
	for _, p := range p.Params {
		types = append(types, p.Type)
	}
	meta.Signature = signature.NewTupleType(types).Signature()
	return meta
}

// InterfaceType represents a parsed IDL interface. It implements
// signature.Type.
type InterfaceType struct {
	Name        string
	PackageName string
	Methods     map[uint32]Method
	Signals     map[uint32]Signal
	Properties  map[uint32]Property
	Scope       Scope
	Namespace   Namespace
}

// Signature returns "o".
func (s *InterfaceType) Signature() string {
	return "o"
}

// SignatureIDL returns "obj".
func (s *InterfaceType) SignatureIDL() string {
	return "obj"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (s *InterfaceType) TypeName() *jen.Statement {
	return jen.Qual(s.PackageName, s.Name)
}

// TypeDeclaration writes the type declaration into file.
// It generates the proxy for the interface.
func (s *InterfaceType) TypeDeclaration(f *jen.File) {
	// proxy.generateInterface(f, s)
	panic("not yet implemented")
}

// RegisterTo adds the type to the type set.
func (s *InterfaceType) RegisterTo(set *signature.TypeSet) {
	name := s.Name
	// loop 100 times to avoid name collision
	for i := 0; i < 100; i++ {
		ok := true // can use the name
		for i, n := range set.Names {
			if n == name {
				if set.Types[i].Signature() == s.Signature() {
					// already registered
					return
				}
				ok = false
				break
			}
		}
		if ok {
			// name is not taken
			set.Types = append(set.Types, s)
			set.Names = append(set.Names, name)
			return
		}
		name = fmt.Sprintf("%s_%d", s.Name, i)
	}
	panic("failed to register " + name)
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (s *InterfaceType) Marshal(id string, writer string) *jen.Statement {
	// TODO: InterfaceType instanciation are proxy from which one
	// can construct an ObjectReference. This ObjectReference can
	// be serialized.
	//
	// return jen.Qual("github.com/lugu/qiloop/type/object",
	// 	"WriteObjectReference").Call(jen.Id(id), jen.Id(writer))
	panic("not yet implemented")
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (s *InterfaceType) Unmarshal(reader string) *jen.Statement {
	// TODO: see Marshall
	panic("not yet implemented")
}

// MetaObject returs the MetaObject describing the interface.
func (s *InterfaceType) MetaObject() object.MetaObject {
	var meta object.MetaObject
	meta.Description = s.Name
	meta.Methods = make(map[uint32]object.MetaMethod)
	meta.Signals = make(map[uint32]object.MetaSignal)
	meta.Properties = make(map[uint32]object.MetaProperty)
	for id, m := range s.Methods {
		meta.Methods[id] = m.Meta(id)
	}
	for id, s := range s.Signals {
		meta.Signals[id] = s.Meta(id)
	}
	for id, p := range s.Properties {
		meta.Properties[id] = p.Meta(id)
	}
	return meta
}
