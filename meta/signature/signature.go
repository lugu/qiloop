package signature

import (
	"bytes"
	"fmt"
	"github.com/dave/jennifer/jen"
	parsec "github.com/prataprc/goparsec"
	"io"
	"log"
	"reflect"
	"strings"
)

// MetaObjectSignature is the signature of MetaObject. It is used to
// generate the MetaObject struct which is used to generate the
// services. This step is referred as stage 1.
const MetaObjectSignature string = "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>"

// TypeSet is a container which contains exactly one instance of each
// ValueConstructor currently generated. It is used to generate the
// type declaration only once.
type TypeSet struct {
	signatures map[string]bool
	types      []ValueConstructor
}

// Register allows a ValueConstructor to declare itself for lated code
// generation.
func (s *TypeSet) Register(v ValueConstructor) {
	if _, ok := s.signatures[v.Signature()]; !ok {
		s.signatures[v.Signature()] = true
		s.types = append(s.types, v)
	}
}

// Declare writes all the registered ValueConstructor into the jen.File.
func (s *TypeSet) Declare(f *jen.File) {
	for _, v := range s.types {
		v.typeDeclaration(f)
	}
}

// NewTypeSet construct a new TypeSet.
func NewTypeSet() *TypeSet {
	sig := make(map[string]bool)
	typ := make([]ValueConstructor, 0)
	return &TypeSet{sig, typ}
}

// Statement is a short version for jen.Statement.
type Statement = jen.Statement

// ValueConstructor represents a type of a signature or a type
// embedded inside a signature. ValueConstructor represents types for
// primitive types (int, long, float, string), vectors of a type,
// associative maps and structures.
type ValueConstructor interface {
	Signature() string
	TypeName() *Statement
	typeDeclaration(*jen.File)
	RegisterTo(s *TypeSet)
	Marshal(id string, writer string) *Statement // returns an error
	Unmarshal(reader string) *Statement          // returns (type, err)
}

// Print render the type into a string. It is only used for testing.
func Print(v ValueConstructor) string {
	buf := bytes.NewBufferString("")
	v.TypeName().Render(buf)
	return buf.String()
}

// NewLongValue is a contructor for the representation of a uint64.
func NewLongValue() LongValue {
	return LongValue{}
}

// NewFloatValue is a contructor for the representation of a float32.
func NewFloatValue() FloatValue {
	return FloatValue{}
}

// NewIntValue is a contructor for the representation of an uint32.
func NewIntValue() IntValue {
	return IntValue{}
}

// NewStringValue is a contructor for the representation of a string.
func NewStringValue() StringValue {
	return StringValue{}
}

// NewVoidValue is a contructor for the representation of the
// absence of a return type. Only used in the context of a returned
// type.
func NewVoidValue() VoidValue {
	return VoidValue{}
}

// NewValueValue is a contructor for the representation of a Value.
func NewValueValue() ValueValue {
	return ValueValue{}
}

// NewBoolValue is a contructor for the representation of a bool.
func NewBoolValue() BoolValue {
	return BoolValue{}
}

// NewListValue is a contructor for the representation of a slice.
func NewListValue(value ValueConstructor) *ListValue {
	return &ListValue{value}
}

// NewMapValue is a contructor for the representation of a map.
func NewMapValue(key, value ValueConstructor) *MapValue {
	return &MapValue{key, value}
}

// NewMemberValue is a contructor for the representation of a field in
// a struct.
func NewMemberValue(name string, value ValueConstructor) MemberValue {
	return MemberValue{name, value}
}

// NewStrucValue is a contructor for the representation of a struct.
func NewStrucValue(name string, members []MemberValue) *StructValue {
	return &StructValue{name, members}
}

// NewTupleValue is a contructor for the representation of a series of
// types. Used to describe a method parameters list.
func NewTupleValue(values []ValueConstructor) *TupleValue {
	return &TupleValue{values}
}

// NewMetaObjectValue is a contructor for the representation of an
// object.
func NewMetaObjectValue() MetaObjectValue {
	return MetaObjectValue{}
}

// IntValue represents an integer.
type IntValue struct {
}

// Signature returns "I". "i" is also accepted as an integer
// signature.
func (i IntValue) Signature() string {
	return "I"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (i IntValue) TypeName() *Statement {
	return jen.Uint32()
}

// RegisterTo adds the type to the TypeSet.
func (i IntValue) RegisterTo(s *TypeSet) {
	return
}

func (i IntValue) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (i IntValue) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/basic", "WriteUint32").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (i IntValue) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadUint32").Call(jen.Id(reader))
}

// LongValue represents a long.
type LongValue struct {
}

// Signature returns "L".
func (i LongValue) Signature() string {
	return "L"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (i LongValue) TypeName() *Statement {
	return jen.Uint64()
}

// RegisterTo adds the type to the TypeSet.
func (i LongValue) RegisterTo(s *TypeSet) {
	return
}

func (i LongValue) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (i LongValue) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/basic", "WriteUint64").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (i LongValue) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadUint64").Call(jen.Id(reader))
}

// FloatValue represents a float.
type FloatValue struct {
}

// Signature returns "f".
func (f FloatValue) Signature() string {
	return "f"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (f FloatValue) TypeName() *Statement {
	return jen.Float32()
}

// RegisterTo adds the type to the TypeSet.
func (f FloatValue) RegisterTo(s *TypeSet) {
	return
}

func (f FloatValue) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (f FloatValue) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/basic", "WriteFloat32").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (f FloatValue) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadFloat32").Call(jen.Id(reader))
}

// BoolValue represents a bool.
type BoolValue struct {
}

// Signature returns "b".
func (b BoolValue) Signature() string {
	return "b"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (b BoolValue) TypeName() *Statement {
	return jen.Bool()
}

// RegisterTo adds the type to the TypeSet.
func (b BoolValue) RegisterTo(s *TypeSet) {
	return
}

func (b BoolValue) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (b BoolValue) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/basic", "WriteBool").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (b BoolValue) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadBool").Call(jen.Id(reader))
}

// ValueValue represents a Value.
type ValueValue struct {
}

// Signature returns "m".
func (b ValueValue) Signature() string {
	return "m"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (b ValueValue) TypeName() *Statement {
	return jen.Qual("github.com/lugu/qiloop/value", "Value")
}

// RegisterTo adds the type to the TypeSet.
func (b ValueValue) RegisterTo(s *TypeSet) {
	return
}

func (b ValueValue) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (b ValueValue) Marshal(id string, writer string) *Statement {
	return jen.Id(id).Dot("Write").Call(jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (b ValueValue) Unmarshal(reader string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/value", "NewValue").Call(jen.Id(reader))
}

// VoidValue represents the return type of a method.
type VoidValue struct {
}

// Signature returns "v".
func (v VoidValue) Signature() string {
	return "v"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (v VoidValue) TypeName() *Statement {
	return jen.Empty()
}

// RegisterTo adds the type to the TypeSet.
func (v VoidValue) RegisterTo(s *TypeSet) {
	return
}

func (v VoidValue) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (v VoidValue) Marshal(id string, writer string) *Statement {
	return jen.Nil()
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (v VoidValue) Unmarshal(reader string) *Statement {
	return jen.Empty()
}

// StringValue represents a string.
type StringValue struct {
}

// Signature returns "s".
func (s StringValue) Signature() string {
	return "s"
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (s StringValue) TypeName() *Statement {
	return jen.String()
}

// RegisterTo adds the type to the TypeSet.
func (s StringValue) RegisterTo(set *TypeSet) {
	return
}

func (s StringValue) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (s StringValue) Marshal(id string, writer string) *Statement {
	return jen.Id("basic.WriteString").Call(jen.Id(id), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (s StringValue) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadString").Call(jen.Id(reader))
}

// ListValue represents a slice.
type ListValue struct {
	value ValueConstructor
}

// Signature returns "[<signature>]" where <signature> is the
// signature of the type of the list.
func (l *ListValue) Signature() string {
	return fmt.Sprintf("[%s]", l.value.Signature())
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (l *ListValue) TypeName() *Statement {
	return jen.Index().Add(l.value.TypeName())
}

// RegisterTo adds the type to the TypeSet.
func (l *ListValue) RegisterTo(s *TypeSet) {
	l.value.RegisterTo(s)
	return
}

func (l *ListValue) typeDeclaration(file *jen.File) {
	return
}

type MetaObjectValue struct {
}

func (m MetaObjectValue) Signature() string {
	return MetaObjectSignature
}

func (m MetaObjectValue) TypeName() *Statement {
	return jen.Qual("github.com/lugu/qiloop/object", "MetaObject")
}

func (m MetaObjectValue) typeDeclaration(*jen.File) {
	return
}

func (m MetaObjectValue) RegisterTo(s *TypeSet) {
	return
}

func (m MetaObjectValue) Marshal(id string, writer string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/object", "WriteMetaObject").Call(jen.Id(id), jen.Id(writer))
}

func (m MetaObjectValue) Unmarshal(reader string) *Statement {
	return jen.Qual("github.com/lugu/qiloop/object", "ReadMetaObject").Call(jen.Id(reader))
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (l *ListValue) Marshal(listID string, writer string) *Statement {
	return jen.Func().Params().Params(jen.Error()).Block(
		jen.Err().Op(":=").Qual("github.com/lugu/qiloop/basic", "WriteUint32").Call(jen.Id("uint32").Call(
			jen.Id("len").Call(jen.Id(listID))),
			jen.Id(writer)),
		jen.Id(`if (err != nil) {
            return fmt.Errorf("failed to write slice size: %s", err)
        }`),
		jen.For(
			jen.Id("_, v := range "+listID),
		).Block(
			jen.Err().Op("=").Add(l.value.Marshal("v", writer)),
			jen.Id(`if (err != nil) {
                return fmt.Errorf("failed to write slice value: %s", err)
            }`),
		),
		jen.Return(jen.Nil()),
	).Call()
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (l *ListValue) Unmarshal(reader string) *Statement {
	return jen.Func().Params().Params(
		jen.Id("b").Index().Add(l.value.TypeName()),
		jen.Err().Error(),
	).Block(
		jen.Id("size, err := basic.ReadUint32").Call(jen.Id(reader)),
		jen.If(jen.Id("err != nil")).Block(
			jen.Return(jen.Id("b"), jen.Qual("fmt", "Errorf").Call(jen.Id(`"failed to read slice size: %s", err`)))),
		jen.Id("b").Op("=").Id("make").Call(l.TypeName(), jen.Id("size")),
		jen.For(
			jen.Id("i := 0; i < int(size); i++"),
		).Block(
			jen.Id("b[i], err =").Add(l.value.Unmarshal(reader)),
			jen.Id(`if (err != nil) {
                return b, fmt.Errorf("failed to read slice value: %s", err)
            }`),
		),
		jen.Return(jen.Id("b"), jen.Nil()),
	).Call()
}

// MapValue represents a map.
type MapValue struct {
	key   ValueConstructor
	value ValueConstructor
}

// Signature returns "{<signature key><signature value>}" where
// <signature key> is the signature of the key and <signature value>
// the signature of the value.
func (m *MapValue) Signature() string {
	return fmt.Sprintf("{%s%s}", m.key.Signature(), m.value.Signature())
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (m *MapValue) TypeName() *Statement {
	return jen.Map(m.key.TypeName()).Add(m.value.TypeName())
}

// RegisterTo adds the type to the TypeSet.
func (m *MapValue) RegisterTo(s *TypeSet) {
	m.key.RegisterTo(s)
	m.value.RegisterTo(s)
	return
}

func (m *MapValue) typeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (m *MapValue) Marshal(mapID string, writer string) *Statement {
	return jen.Func().Params().Params(jen.Error()).Block(
		jen.Err().Op(":=").Qual("github.com/lugu/qiloop/basic", "WriteUint32").Call(jen.Id("uint32").Call(
			jen.Id("len").Call(jen.Id(mapID))),
			jen.Id(writer)),
		jen.Id(`if (err != nil) {
            return fmt.Errorf("failed to write map size: %s", err)
        }`),
		jen.For(
			jen.Id("k, v := range "+mapID),
		).Block(
			jen.Err().Op("=").Add(m.key.Marshal("k", writer)),
			jen.Id(`if (err != nil) {
                return fmt.Errorf("failed to write map key: %s", err)
            }`),
			jen.Err().Op("=").Add(m.value.Marshal("v", writer)),
			jen.Id(`if (err != nil) {
                return fmt.Errorf("failed to write map value: %s", err)
            }`),
		),
		jen.Return(jen.Nil()),
	).Call()
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (m *MapValue) Unmarshal(reader string) *Statement {
	return jen.Func().Params().Params(
		jen.Id("m").Map(m.key.TypeName()).Add(m.value.TypeName()),
		jen.Err().Error(),
	).Block(
		jen.Id("size, err := basic.ReadUint32").Call(jen.Id(reader)),
		jen.If(jen.Id("err != nil")).Block(
			jen.Return(jen.Id("m"), jen.Qual("fmt", "Errorf").Call(jen.Id(`"failed to read map size: %s", err`)))),
		jen.Id("m").Op("=").Id("make").Call(m.TypeName(), jen.Id("size")),
		jen.For(
			jen.Id("i := 0; i < int(size); i++"),
		).Block(
			jen.Id("k, err :=").Add(m.key.Unmarshal(reader)),
			jen.Id(`if (err != nil) {
                return m, fmt.Errorf("failed to read map key: %s", err)
            }`),
			jen.Id("v, err :=").Add(m.value.Unmarshal(reader)),
			jen.Id(`if (err != nil) {
                return m, fmt.Errorf("failed to read map value: %s", err)
            }`),
			jen.Id("m[k] = v"),
		),
		jen.Return(jen.Id("m"), jen.Nil()),
	).Call()
}

// MemberValue a field in a struct.
type MemberValue struct {
	Name  string
	Value ValueConstructor
}

// Title is the public name of the field.
func (m MemberValue) Title() string {
	return strings.Title(m.Name)
}

// TupleValue a list of a parameter of a method.
type TupleValue struct {
	values []ValueConstructor
}

// Signature returns "(<signature 1><signature 2>...)" where
// <signature X> is the signature of the elements.
func (t *TupleValue) Signature() string {
	sig := "("
	for _, v := range t.values {
		sig += v.Signature()
	}
	sig += ")"
	return sig
}

// Members returns the list of the types composing the TupleValue.
func (t *TupleValue) Members() []MemberValue {
	members := make([]MemberValue, len(t.values))
	for i, v := range t.values {
		members[i] = MemberValue{fmt.Sprintf("p%d", i), v}
	}
	return members
}

// Params returns a statement representing the list of parameter of
// a method.
func (t *TupleValue) Params() *Statement {
	arguments := make([]jen.Code, len(t.values))
	for i, v := range t.values {
		arguments[i] = jen.Id(fmt.Sprintf("p%d", i)).Add(v.TypeName())
	}
	return jen.Params(arguments...)
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (t *TupleValue) TypeName() *Statement {
	return jen.Id("...interface{}")
}

// RegisterTo adds the type to the TypeSet.
func (t *TupleValue) RegisterTo(s *TypeSet) {
	for _, v := range t.values {
		v.RegisterTo(s)
	}
	return
}

func (t *TupleValue) typeDeclaration(*jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (t *TupleValue) Marshal(variadicIdentifier string, writer string) *Statement {
	// TODO: shall returns an error
	return jen.Empty()
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (t *TupleValue) Unmarshal(reader string) *Statement {
	// TODO: shall returns (type, err)
	return jen.Empty()
}

// ConvertMetaObjects replace any element type which has the same
// signature as MetaObject with an element of the type
// object.MetaObject. This is required to generate proxy services
// which implements the object.Object interface and avoid a circular
// dependancy.
func (t *TupleValue) ConvertMetaObjects() {
	for i, member := range t.values {
		if member.Signature() == MetaObjectSignature {
			t.values[i] = NewMetaObjectValue()
		}
	}
}

// StructValue represents a struct.
type StructValue struct {
	name    string
	members []MemberValue
}

// Signature returns the signature of the struct.
func (s *StructValue) Signature() string {
	types := ""
	names := make([]string, 0, len(s.members))
	for _, v := range s.members {
		names = append(names, v.Name)
		if s, ok := v.Value.(*StructValue); ok {
			types += "[" + s.Signature() + "]"
		} else {
			types += v.Value.Signature()
		}
	}
	return fmt.Sprintf("(%s)<%s,%s>", types,
		s.name, strings.Join(names, ","))
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (s *StructValue) TypeName() *Statement {
	return jen.Id(s.name)
}

// RegisterTo adds the type to the TypeSet.
func (s *StructValue) RegisterTo(set *TypeSet) {
	for _, v := range s.members {
		v.Value.RegisterTo(set)
	}
	set.Register(s)
	return
}

func (s *StructValue) typeDeclaration(file *jen.File) {
	fields := make([]jen.Code, len(s.members))
	for i, v := range s.members {
		fields[i] = jen.Id(v.Title()).Add(v.Value.TypeName())
	}
	file.Type().Id(s.name).Struct(fields...)

	readFields := make([]jen.Code, len(s.members)+1)
	writeFields := make([]jen.Code, len(s.members)+1)
	for i, v := range s.members {
		readFields[i] = jen.If(
			jen.Id("s."+v.Title()+", err =").Add(v.Value.Unmarshal("r")),
			jen.Id("err != nil")).Block(
			jen.Id(`return s, fmt.Errorf("failed to read ` + v.Title() + ` field: %s", err)`),
		)
		writeFields[i] = jen.If(
			jen.Id("err :=").Add(v.Value.Marshal("s."+v.Title(), "w")),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Id(`return fmt.Errorf("failed to write ` + v.Title() + ` field: %s", err)`),
		)
	}
	readFields[len(s.members)] = jen.Return(jen.Id("s"), jen.Nil())
	writeFields[len(s.members)] = jen.Return(jen.Nil())

	file.Func().Id("Read"+s.name).Params(
		jen.Id("r").Id("io.Reader"),
	).Params(
		jen.Id("s").Id(s.name), jen.Err().Error(),
	).Block(readFields...)
	file.Func().Id("Write"+s.name).Params(
		jen.Id("s").Id(s.name),
		jen.Id("w").Qual("io", "Writer"),
	).Params(jen.Err().Error()).Block(writeFields...)
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (s *StructValue) Marshal(structID string, writer string) *Statement {
	return jen.Id("Write"+s.name).Call(jen.Id(structID), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (s *StructValue) Unmarshal(reader string) *Statement {
	return jen.Id("Read" + s.name).Call(jen.Id(reader))
}

func basicType() parsec.Parser {
	return parsec.OrdChoice(nodifyBasicType,
		parsec.Atom("I", "uint32"),
		parsec.Atom("s", "string"),
		parsec.Atom("L", "uint64"),
		parsec.Atom("b", "bool"),
		parsec.Atom("f", "float32"),
		parsec.Atom("m", "value"),
		parsec.Atom("v", "void"))
}

func typeName() parsec.Parser {
	return parsec.Ident()
}

// Node is an alias to parsec.ParsecNode
type Node = parsec.ParsecNode

func nodifyBasicType(nodes []Node) Node {
	if len(nodes) != 1 {
		log.Panicf("wrong basic arguments %+v\n", nodes)
	}
	signature := nodes[0].(*parsec.Terminal).GetValue()
	switch signature {
	case "I":
		return NewIntValue()
	case "L":
		return NewLongValue()
	case "s":
		return NewStringValue()
	case "b":
		return NewBoolValue()
	case "f":
		return NewFloatValue()
	case "v":
		return NewVoidValue()
	case "m":
		return NewValueValue()
	default:
		log.Panicf("wrong signature %s", signature)
	}
	return nil
}

func extractValue(object interface{}) (ValueConstructor, error) {
	nodes, ok := object.([]Node)
	if !ok {
		return nil, fmt.Errorf("extraction failed: %+v", reflect.TypeOf(object))
	}
	value, ok := nodes[0].(ValueConstructor)
	if !ok {
		return nil, fmt.Errorf("extraction failed bis: %+v", reflect.TypeOf(nodes[0]))
	}
	return value, nil
}

func nodifyMap(nodes []Node) Node {
	if len(nodes) != 4 {
		fmt.Printf("wrong map arguments %+v\n", nodes)
	}
	key, err := extractValue(nodes[1])
	if err != nil {
		log.Panicf("key conversion failed: %s", err)
	}
	value, err := extractValue(nodes[2])
	if err != nil {
		log.Panicf("value conversion failed: %s", err)
	}
	return NewMapValue(key, value)
}

func nodifyArrayType(nodes []Node) Node {
	if len(nodes) != 3 {
		log.Panicf("wrong arguments %+v", nodes)
	}
	value, err := extractValue(nodes[1])
	if err != nil {
		log.Panicf("value conversion failed: %s", err)
	}
	return NewListValue(value)
}

func extractMembersTypes(node Node) []ValueConstructor {
	typeList, ok := node.([]Node)
	if !ok {
		log.Panicf("member type list is not a list: %s", reflect.TypeOf(node))
	}
	types := make([]ValueConstructor, len(typeList))
	for i := range typeList {
		memberType, err := extractValue(typeList[i])
		if err != nil {
			log.Panicf("member type value conversion failed: %s", err)
		}
		types[i] = memberType
	}
	return types
}

func extractMembersName(node Node) []string {
	baseList, ok := node.([]Node)
	if !ok {
		log.Panicf("member name list is not a list: %s", reflect.TypeOf(node))
	}
	membersList, ok := baseList[0].([]Node)
	if !ok {
		log.Panicf("member name list is not a list: %s", reflect.TypeOf(node))
	}
	names := make([]string, len(membersList))
	for i, n := range membersList {
		memberName, ok := n.(*parsec.Terminal)
		if !ok {
			log.Panicf("failed to convert member names %s", reflect.TypeOf(n))
		}
		names[i] = memberName.GetValue()
	}
	return names
}

func extractMembers(typesNode, namesNode Node) []MemberValue {

	types := extractMembersTypes(typesNode)
	names := extractMembersName(namesNode)

	if len(types) != len(names) {
		log.Panicf("member types and names of different size: %+v, %+v", types, names)
	}

	members := make([]MemberValue, len(names))

	for i := range types {
		members[i] = NewMemberValue(names[i], types[i])
	}
	return members
}

func nodifyTupleType(nodes []Node) Node {

	types := extractMembersTypes(nodes[1])

	return NewTupleValue(types)
}

func nodifyTypeDefinition(nodes []Node) Node {

	terminal, ok := nodes[4].(*parsec.Terminal)
	if !ok {
		log.Panicf("wrong name %s", reflect.TypeOf(nodes[4]))
	}
	name := terminal.GetValue()
	members := extractMembers(nodes[1], nodes[6])

	return NewStrucValue(name, members)
}

// Parse reads a signature contained in a string and constructs its
// type representation.
func Parse(input string) (ValueConstructor, error) {
	text := []byte(input)

	var arrayType parsec.Parser
	var mapType parsec.Parser
	var typeDefinition parsec.Parser
	var tupleType parsec.Parser

	var declarationType = parsec.OrdChoice(nil,
		basicType(), &mapType, &arrayType, &typeDefinition, &tupleType)

	arrayType = parsec.And(nodifyArrayType,
		parsec.Atom("[", "MapStart"),
		&declarationType,
		parsec.Atom("]", "MapClose"))

	var listType = parsec.Kleene(nil, &declarationType)

	var typeMemberList = parsec.Maybe(nil,
		parsec.Many(nil, typeName(), parsec.Atom(",", "Separator")))

	tupleType = parsec.And(nodifyTupleType,
		parsec.Atom("(", "TypeParameterStart"),
		&listType,
		parsec.Atom(")", "TypeParameterClose"))

	typeDefinition = parsec.And(nodifyTypeDefinition,
		parsec.Atom("(", "TypeParameterStart"),
		&listType,
		parsec.Atom(")", "TypeParameterClose"),
		parsec.Atom("<", "TypeDefinitionStart"),
		typeName(),
		parsec.Atom(",", "TypeName"),
		&typeMemberList,
		parsec.Atom(">", "TypeDefinitionClose"))

	mapType = parsec.And(nodifyMap,
		parsec.Atom("{", "MapStart"),
		&declarationType, &declarationType,
		parsec.Atom("}", "MapClose"))

	var typeSignature = declarationType

	root, _ := typeSignature(parsec.NewScanner(text))
	if root == nil {
		return nil, fmt.Errorf("failed to parse signature")
	}
	types, ok := root.([]Node)
	if !ok {
		return nil, fmt.Errorf("failed to convert array: %+v", reflect.TypeOf(root))
	}
	if len(types) != 1 {
		return nil, fmt.Errorf("did not parse only one type: %+v", root)
	}
	constructor, ok := types[0].(ValueConstructor)
	if !ok {
		return nil, fmt.Errorf("failed to convert value: %+v", reflect.TypeOf(types[0]))
	}
	return constructor, nil
}

// GenerateType generate the code required to serialize the given
// type.
func GenerateType(v ValueConstructor, packageName string, w io.Writer) error {
	file := jen.NewFile(packageName)
	file.PackageComment("file generated. DO NOT EDIT.")
	set := NewTypeSet()
	v.RegisterTo(set)
	set.Declare(file)
	if err := file.Render(w); err != nil {
		return fmt.Errorf("failed to render %s: %s", v.Signature(), err)
	}
	return nil
}
