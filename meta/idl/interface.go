package idl

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/meta/signature"
)

// Proxy were generated from MetaObjects. This is convinient from a
// boostraping point of view. Since the data structures are well
// known, we introduce the InterfaceType which contains the
// information of the MetaObject with the possibility to resolve
// object references.

// Two stage parsing of the IDL:
// 1. Construct a set of Types associated with context.
// 2. Defer the resolution of the types to the proxy/stub generation.

type Scope map[string]signature.Type
type Namespace map[string]Scope

type Function struct {
	Name     string
	Return   string
	ArgNames []string
	ArgTypes []string
}

type Signal struct {
	Name     string
	ArgNames []string
	ArgTypes []string
}

type Property struct {
	Name   string
	Return string
}

type InterfaceType struct {
	name        string
	packageName string
	functions   map[uint32]Function
	signals     map[uint32]Signal
	properties  map[uint32]Property
	scope       Scope
	namespace   Namespace
}

func (s *InterfaceType) Signature() string {
	return "o"
}
func (s *InterfaceType) SignatureIDL() string {
	return "obj"
}

func (s *InterfaceType) TypeName() *jen.Statement {
	return jen.Qual(s.packageName, s.name)
}

// declare the interface as done by the proxy generation
func (s *InterfaceType) TypeDeclaration(f *jen.File) {
	// proxy.generateInterface(f, s)
	panic("not yet implemented")
}

func (s *InterfaceType) RegisterTo(set *signature.TypeSet) {
	name := s.name
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
		name = fmt.Sprintf("%s_%d", s.name, i)
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
