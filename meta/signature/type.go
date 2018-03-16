package signature

import (
	"bytes"
	"fmt"
	"github.com/dave/jennifer/jen"
	"strings"
)

// TypeSet is a container which contains exactly one instance of each
// Type currently generated. It is used to generate the
// type declaration only once.
type TypeSet struct {
	Signatures map[string]string // maps type name with signatures
	Types      []Type
}

// Declare writes all the registered Type into the jen.File.
func (s *TypeSet) Declare(f *jen.File) {
	for _, v := range s.Types {
		v.TypeDeclaration(f)
	}
}

// NewTypeSet construct a new TypeSet.
func NewTypeSet() *TypeSet {
	sig := make(map[string]string)
	typ := make([]Type, 0)
	return &TypeSet{sig, typ}
}

// Type represents a type of a signature or a type
// embedded inside a signature. Type represents types for
// primitive types (int, long, float, string), vectors of a type,
// associative maps and structures.
type Type interface {
	Signature() string
	SignatureIDL() string
	TypeName() *Statement
	TypeDeclaration(*jen.File)
	RegisterTo(s *TypeSet)
	Marshal(id string, writer string) *Statement // returns an error
	Unmarshal(reader string) *Statement          // returns (type, err)
}

type TypeConstructor struct {
	signature    string
	signatureIDL string
	typeName     *Statement
	marshal      func(id string, writer string) *Statement // returns an error
	unmarshal    func(reader string) *Statement            // returns (type, err)
}

func (t *TypeConstructor) Signature() string {
	return t.signature
}
func (t *TypeConstructor) SignatureIDL() string {
	return t.signatureIDL
}
func (t *TypeConstructor) TypeName() *Statement {
	return t.typeName
}
func (t *TypeConstructor) TypeDeclaration(f *jen.File) {
	return
}
func (t *TypeConstructor) RegisterTo(s *TypeSet) {
	return
}
func (t *TypeConstructor) Marshal(id string, writer string) *Statement {
	return t.marshal(id, writer)
}
func (t *TypeConstructor) Unmarshal(reader string) *Statement {
	return t.unmarshal(reader)
}

// Print render the type into a string. It is only used for testing.
func Print(v Type) string {
	buf := bytes.NewBufferString("")
	v.TypeName().Render(buf)
	return buf.String()
}

// NewLongType is a contructor for the representation of a uint64.
func NewLongType() Type {
	return &TypeConstructor{
		signature:    "l",
		signatureIDL: "int64",
		typeName:     jen.Int64(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteInt64").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Id("basic.ReadInt64").Call(jen.Id(reader))
		},
	}
}

// NewULongType is a contructor for the representation of a uint64.
func NewULongType() Type {
	return &TypeConstructor{
		signature:    "L",
		signatureIDL: "uint64",
		typeName:     jen.Uint64(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteUint64").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Id("basic.ReadUint64").Call(jen.Id(reader))
		},
	}
}

// NewFloatType is a contructor for the representation of a float32.
func NewFloatType() Type {
	return &TypeConstructor{
		signature:    "f",
		signatureIDL: "float32",
		typeName:     jen.Float32(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteFloat32").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Id("basic.ReadFloat32").Call(jen.Id(reader))
		},
	}
}

// NewDoubleType is a contructor for the representation of a float32.
func NewDoubleType() Type {
	return &TypeConstructor{
		signature:    "d",
		signatureIDL: "float64",
		typeName:     jen.Float64(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteFloat64").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic", "ReadFloat64").Call(jen.Id(reader))
		},
	}
}

// NewIntType is a contructor for the representation of an int32.
func NewIntType() Type {
	return &TypeConstructor{
		signature:    "i",
		signatureIDL: "int32",
		typeName:     jen.Int32(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteInt32").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Id("basic.ReadInt32").Call(jen.Id(reader))
		},
	}
}

// NewUIntType is a contructor for the representation of an uint32.
func NewUIntType() Type {
	return &TypeConstructor{
		signature:    "I",
		signatureIDL: "uint32",
		typeName:     jen.Uint32(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteUint32").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Id("basic.ReadUint32").Call(jen.Id(reader))
		},
	}
}

// NewStringType is a contructor for the representation of a string.
func NewStringType() Type {
	return &TypeConstructor{
		signature:    "s",
		signatureIDL: "str",
		typeName:     jen.String(),
		marshal: func(id string, writer string) *Statement {
			return jen.Id("basic.WriteString").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Id("basic.ReadString").Call(jen.Id(reader))
		},
	}
}

// NewVoidType is a contructor for the representation of the
// absence of a return type. Only used in the context of a returned
// type.
func NewVoidType() Type {
	return &TypeConstructor{
		signature:    "v",
		signatureIDL: "nothing",
		typeName:     jen.Empty(),
		marshal: func(id string, writer string) *Statement {
			return jen.Nil()
		},
		unmarshal: func(reader string) *Statement {
			return jen.Empty()
		},
	}
}

// NewValueType is a contructor for the representation of a Value.
func NewValueType() Type {
	return &TypeConstructor{
		signature:    "m",
		signatureIDL: "any",
		typeName:     jen.Qual("github.com/lugu/qiloop/type/value", "Value"),
		marshal: func(id string, writer string) *Statement {
			return jen.Id(id).Dot("Write").Call(jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/value", "NewValue").Call(jen.Id(reader))
		},
	}
}

// NewBoolType is a contructor for the representation of a bool.
func NewBoolType() Type {
	return &TypeConstructor{
		signature:    "b",
		signatureIDL: "bool",
		typeName:     jen.Bool(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteBool").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Id("basic.ReadBool").Call(jen.Id(reader))
		},
	}
}

// NewMetaObjectType is a contructor for the representation of an
// object.
func NewMetaObjectType() Type {
	return &TypeConstructor{
		signature:    MetaObjectSignature,
		signatureIDL: "MetaObject",
		typeName:     jen.Qual("github.com/lugu/qiloop/type/object", "MetaObject"),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/object", "WriteMetaObject").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/object", "ReadMetaObject").Call(jen.Id(reader))
		},
	}
}

// NewObjectType is a contructor for the representation of a Value.
func NewObjectType() Type {
	return &TypeConstructor{
		signature:    "o",
		signatureIDL: "obj",
		typeName:     jen.Qual("github.com/lugu/qiloop/type/object", "ObjectReference"),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/object", "WriteObjectReference").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/object", "ReadObjectReference").Call(jen.Id(reader))
		},
	}
}

// NewUnknownType is a contructor for an unkown type.
func NewUnknownType() Type {
	return &TypeConstructor{
		signature:    "X",
		signatureIDL: "unknown",
		typeName:     jen.Id("interface{}"),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("fmt", "Errorf").Call(jen.Lit("unknown type serialization not supported: %v"), jen.Id(id))
		},
		unmarshal: func(reader string) *Statement {
			return jen.List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("unknown type deserialization not supported")))
		},
	}
}

// NewListType is a contructor for the representation of a slice.
func NewListType(value Type) *ListType {
	return &ListType{value}
}

// ListType represents a slice.
type ListType struct {
	value Type
}

// Signature returns "[<signature>]" where <signature> is the
// signature of the type of the list.
func (l *ListType) Signature() string {
	return fmt.Sprintf("[%s]", l.value.Signature())
}

// SignatureIDL returns "Vec<signature>" where "signature" is the IDL
// signature of the type of the list.
func (l *ListType) SignatureIDL() string {
	return fmt.Sprintf("Vec<%s>", l.value.SignatureIDL())
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (l *ListType) TypeName() *Statement {
	return jen.Index().Add(l.value.TypeName())
}

// RegisterTo adds the type to the TypeSet.
func (l *ListType) RegisterTo(s *TypeSet) {
	l.value.RegisterTo(s)
	return
}

func (l *ListType) TypeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (l *ListType) Marshal(listID string, writer string) *Statement {
	return jen.Func().Params().Params(jen.Error()).Block(
		jen.Err().Op(":=").Qual("github.com/lugu/qiloop/type/basic", "WriteUint32").Call(jen.Id("uint32").Call(
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
func (l *ListType) Unmarshal(reader string) *Statement {
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

// NewMapType is a contructor for the representation of a map.
func NewMapType(key, value Type) *MapType {
	return &MapType{key, value}
}

// MapType represents a map.
type MapType struct {
	key   Type
	value Type
}

// Signature returns "{<signature key><signature value>}" where
// <signature key> is the signature of the key and <signature value>
// the signature of the value.
func (m *MapType) Signature() string {
	return fmt.Sprintf("{%s%s}", m.key.Signature(), m.value.Signature())
}

// SignatureIDL returns "map<key><value>" where "key" is the IDL
// signature of the key and "value" the IDL signature of the value.
func (m *MapType) SignatureIDL() string {
	return fmt.Sprintf("Map<%s,%s>", m.key.SignatureIDL(), m.value.SignatureIDL())
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (m *MapType) TypeName() *Statement {
	return jen.Map(m.key.TypeName()).Add(m.value.TypeName())
}

// RegisterTo adds the type to the TypeSet.
func (m *MapType) RegisterTo(s *TypeSet) {
	m.key.RegisterTo(s)
	m.value.RegisterTo(s)
	return
}

func (m *MapType) TypeDeclaration(file *jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (m *MapType) Marshal(mapID string, writer string) *Statement {
	return jen.Func().Params().Params(jen.Error()).Block(
		jen.Err().Op(":=").Qual("github.com/lugu/qiloop/type/basic", "WriteUint32").Call(jen.Id("uint32").Call(
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
func (m *MapType) Unmarshal(reader string) *Statement {
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

// NewMemberType is a contructor for the representation of a field in
// a struct.
func NewMemberType(name string, value Type) MemberType {
	return MemberType{name, value}
}

// MemberType a field in a struct.
type MemberType struct {
	Name  string
	Value Type
}

// Title is the public name of the field.
func (m MemberType) Title() string {
	return strings.Title(m.Name)
}

// NewTupleType is a contructor for the representation of a series of
// types. Used to describe a method parameters list.
func NewTupleType(values []Type) *TupleType {
	return &TupleType{values}
}

// TupleType a list of a parameter of a method.
type TupleType struct {
	values []Type
}

// Signature returns "(<signature 1><signature 2>...)" where
// <signature X> is the signature of the elements.
func (t *TupleType) Signature() string {
	sig := "("
	for _, v := range t.values {
		sig += v.Signature()
	}
	sig += ")"
	return sig
}

// SignatureIDL returns "name1 signature1, name2 signature2" where
// signatureX is the signature of the elements.
func (t *TupleType) SignatureIDL() string {
	var i int = -1
	var typ MemberType
	sig := ""
	for i, typ = range t.Members() {
		if i == len(t.values)-1 {
			break
		}
		sig += typ.Name + ": " + typ.Value.SignatureIDL() + ", "
	}
	if len(t.values) > 0 {
		sig += t.Members()[i].Name + ": " + t.Members()[i].Value.SignatureIDL()
	}
	return sig
}

// Members returns the list of the types composing the TupleType.
func (t *TupleType) Members() []MemberType {
	members := make([]MemberType, len(t.values))
	for i, v := range t.values {
		members[i] = MemberType{fmt.Sprintf("P%d", i), v}
	}
	return members
}

// Params returns a statement representing the list of parameter of
// a method.
func (t *TupleType) Params() *Statement {
	arguments := make([]jen.Code, len(t.values))
	for i, v := range t.values {
		arguments[i] = jen.Id(fmt.Sprintf("P%d", i)).Add(v.TypeName())
	}
	return jen.Params(arguments...)
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (t *TupleType) TypeName() *Statement {
	params := make([]jen.Code, 0)
	for _, typ := range t.Members() {
		params = append(params, jen.Id(typ.Name).Add(typ.Value.TypeName()))
	}
	return jen.Struct(params...)
}

// RegisterTo adds the type to the TypeSet.
func (t *TupleType) RegisterTo(s *TypeSet) {
	for _, v := range t.values {
		v.RegisterTo(s)
	}
	return
}

func (t *TupleType) TypeDeclaration(*jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (t *TupleType) Marshal(tupleID string, writer string) *Statement {
	statements := make([]jen.Code, 0)
	for _, typ := range t.Members() {
		s1 := jen.Err().Op("=").Add(typ.Value.Marshal(tupleID+"."+typ.Name, writer))
		s2 := jen.Id(`if (err != nil) {
			return fmt.Errorf("failed to write tuple member: %s", err)
		}`)
		statements = append(statements, s1)
		statements = append(statements, s2)
	}
	statements = append(statements, jen.Return(jen.Nil()))
	return jen.Qual("fmt", "Errorf").Call(jen.Lit("unknown type serialization not implemented: %v"), jen.Id(tupleID))
	return jen.Func().Params().Params(jen.Error()).Block(
		statements...,
	).Call()
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (t *TupleType) Unmarshal(reader string) *Statement {
	statements := make([]jen.Code, 0)
	for _, typ := range t.Members() {
		s1 := jen.List(jen.Id("s."+typ.Name), jen.Err()).Op("=").Add(typ.Value.Unmarshal(reader))
		s2 := jen.Id(`if (err != nil) {
			return s, fmt.Errorf("failed to read tuple member: %s", err)
		}`)
		statements = append(statements, s1)
		statements = append(statements, s2)
	}
	statements = append(statements, jen.Return(jen.Id("s"), jen.Nil()))
	return jen.Func().Params().Params(
		jen.Id("s").Add(t.TypeName()),
		jen.Err().Error(),
	).Block(
		statements...,
	).Call()
}

// ConvertMetaObjects replace any element type which has the same
// signature as MetaObject with an element of the type
// object.MetaObject. This is required to generate proxy services
// which implements the object.Object interface and avoid a circular
// dependancy.
func (t *TupleType) ConvertMetaObjects() {
	for i, member := range t.values {
		if member.Signature() == MetaObjectSignature {
			t.values[i] = NewMetaObjectType()
		}
	}
}

// NewStrucType is a contructor for the representation of a struct.
func NewStrucType(name string, members []MemberType) *StructType {
	return &StructType{name, members}
}

// StructType represents a struct.
type StructType struct {
	Name    string
	Members []MemberType
}

// Signature returns the signature of the struct.
func (s *StructType) Signature() string {
	types := ""
	names := make([]string, 0, len(s.Members))
	for _, v := range s.Members {
		names = append(names, v.Name)
		if s, ok := v.Value.(*StructType); ok {
			types += "[" + s.Signature() + "]"
		} else {
			types += v.Value.Signature()
		}
	}
	return fmt.Sprintf("(%s)<%s,%s>", types,
		s.Name, strings.Join(names, ","))
}

// SignatureIDL returns the idl signature of the struct.
func (s *StructType) SignatureIDL() string {
	return s.Name
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (s *StructType) TypeName() *Statement {
	return jen.Id(s.Name)
}

// RegisterTo adds the type to the TypeSet.
func (s *StructType) RegisterTo(set *TypeSet) {
	for _, v := range s.Members {
		v.Value.RegisterTo(set)
	}

	// register the name of the struct. if the name is used by a
	// struct of a different signature, change the name.
	for i := 0; i < 100; i++ {
		if sgn, ok := set.Signatures[s.Name]; !ok {
			// name not yet used
			set.Signatures[s.Name] = s.Signature()
			set.Types = append(set.Types, s)
			break
		} else if sgn == s.Signature() {
			// same name and same signatures
			break
		} else {
			s.Name = fmt.Sprintf("%s_%d", s.Name, i)
		}
	}
	return
}

func (s *StructType) TypeDeclaration(file *jen.File) {
	fields := make([]jen.Code, len(s.Members))
	for i, v := range s.Members {
		fields[i] = jen.Id(v.Title()).Add(v.Value.TypeName())
	}
	file.Type().Id(s.Name).Struct(fields...)

	readFields := make([]jen.Code, len(s.Members)+1)
	writeFields := make([]jen.Code, len(s.Members)+1)
	for i, v := range s.Members {
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
	readFields[len(s.Members)] = jen.Return(jen.Id("s"), jen.Nil())
	writeFields[len(s.Members)] = jen.Return(jen.Nil())

	file.Func().Id("Read"+s.Name).Params(
		jen.Id("r").Id("io.Reader"),
	).Params(
		jen.Id("s").Id(s.Name), jen.Err().Error(),
	).Block(readFields...)
	file.Func().Id("Write"+s.Name).Params(
		jen.Id("s").Id(s.Name),
		jen.Id("w").Qual("io", "Writer"),
	).Params(jen.Err().Error()).Block(writeFields...)
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (s *StructType) Marshal(structID string, writer string) *Statement {
	return jen.Id("Write"+s.Name).Call(jen.Id(structID), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (s *StructType) Unmarshal(reader string) *Statement {
	return jen.Id("Read" + s.Name).Call(jen.Id(reader))
}
