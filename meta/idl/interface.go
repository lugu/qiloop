package idl

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
	name       string
	functions  map[uint32]Function
	signals    map[uint32]Signal
	properties map[uint32]Property
	scope      Scope
	namespace  Namespace
}
