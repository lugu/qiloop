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

type Namespace map[string]Scope

type Parameter struct {
	Name string
	Type signature.Type
}

type Method struct {
	Name   string
	Id     uint32
	Return signature.Type
	Params []Parameter
}

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

func (m Method) Tuple() *signature.TupleType {
	var tuple signature.TupleType
	tuple.Members = make([]signature.MemberType, 0)
	for _, p := range m.Params {
		tuple.Members = append(tuple.Members,
			signature.MemberType{p.Name, p.Type})
	}
	return &tuple
}

type Signal struct {
	Name   string
	Id     uint32
	Params []Parameter
}

func (s Signal) Tuple() *signature.TupleType {
	var tuple signature.TupleType
	tuple.Members = make([]signature.MemberType, 0)
	for _, p := range s.Params {
		tuple.Members = append(tuple.Members,
			signature.MemberType{p.Name, p.Type})
	}
	return &tuple
}

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

type Property struct {
	Name   string
	Id     uint32
	Params []Parameter
}

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

type InterfaceType struct {
	Name        string
	PackageName string
	Methods     map[uint32]Method
	Signals     map[uint32]Signal
	Properties  map[uint32]Property
	Scope       Scope
	Namespace   Namespace
}

func (s *InterfaceType) Signature() string {
	return "o"
}
func (s *InterfaceType) SignatureIDL() string {
	return "obj"
}

func (s *InterfaceType) TypeName() *jen.Statement {
	return jen.Qual(s.PackageName, s.Name)
}

// declare the interface as done by the proxy generation
func (s *InterfaceType) TypeDeclaration(f *jen.File) {
	// proxy.generateInterface(f, s)
	panic("not yet implemented")
}

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

func (s *InterfaceType) Marshal(id string, writer string) *jen.Statement {
	// TODO: InterfaceType instanciation are proxy from which one
	// can construct an ObjectReference. This ObjectReference can
	// be serialized.
	//
	// return jen.Qual("github.com/lugu/qiloop/type/object",
	// 	"WriteObjectReference").Call(jen.Id(id), jen.Id(writer))
	panic("not yet implemented")
	return nil
}

func (s *InterfaceType) Unmarshal(reader string) *jen.Statement {
	// TODO: see Marshall
	panic("not yet implemented")
	return nil
}
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
	for id, p := range s.Signals {
		meta.Signals[id] = p.Meta(id)
	}
	return meta
}
