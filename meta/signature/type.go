package signature

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"reflect"

	"github.com/dave/jennifer/jen"
)

// TypeSet is a container which contains exactly one instance of each
// Type currently generated. It is used to generate the
// type declaration only once.
type TypeSet struct {
	Names []string
	Types []Type
}

// Declare writes all the registered Type into the jen.File.
func (s *TypeSet) Declare(f *jen.File) {
	for _, v := range s.Types {
		v.TypeDeclaration(f)
	}
}

// Search returns a Type if a name is associated.
func (s *TypeSet) Search(name string) Type {
	for i, n := range s.Names {
		if n == name {
			return s.Types[i]
		}
	}
	return nil
}

// ResolveCollision returns a type name not already present in the
// set. If name if not present, returns name, else it returns a new
// string derived from name.
func (s *TypeSet) ResolveCollision(originalName, signature string) string {
	name := originalName
	// loop 100 times to avoid name collision
	for i := 0; i < 100; i++ {
		ok := true // can use the name
		for i, n := range s.Names {
			if n == name {
				if s.Types[i].Signature() == signature {
					// already registered
					return name
				}
				ok = false
				break
			}
		}
		if ok {
			return name
		}
		name = fmt.Sprintf("%s_%d", originalName, i)
	}
	return "can_not_register_name_" + originalName
}

// NewTypeSet construct a new TypeSet.
func NewTypeSet() *TypeSet {
	typ := make([]Type, 0)
	nam := make([]string, 0)
	return &TypeSet{nam, typ}
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
	Reader() TypeReader
}

type typeConstructor struct {
	signature    string
	signatureIDL string
	typeName     *Statement
	marshal      func(id string, writer string) *Statement // returns an error
	unmarshal    func(reader string) *Statement            // returns (type, err)
	reader       TypeReader
}

func (t *typeConstructor) Signature() string {
	return t.signature
}
func (t *typeConstructor) SignatureIDL() string {
	return t.signatureIDL
}
func (t *typeConstructor) TypeName() *Statement {
	return t.typeName
}
func (t *typeConstructor) TypeDeclaration(f *jen.File) {
	return
}
func (t *typeConstructor) RegisterTo(s *TypeSet) {
	return
}
func (t *typeConstructor) Marshal(id string, writer string) *Statement {
	return t.marshal(id, writer)
}
func (t *typeConstructor) Unmarshal(reader string) *Statement {
	return t.unmarshal(reader)
}
func (t *typeConstructor) Reader() TypeReader {
	return t.reader
}

// Print render the type into a string. It is only used for testing.
func Print(v Type) string {
	buf := bytes.NewBufferString("")
	v.TypeName().Render(buf)
	return buf.String()
}

// NewInt8Type is a contructor for the representation of a uint64.
func NewInt8Type() Type {
	return &typeConstructor{
		signature:    "c",
		signatureIDL: "int8",
		typeName:     jen.Int8(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"WriteInt8").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"ReadInt8").Call(jen.Id(reader))
		},
		reader: constReader(1),
	}
}

// NewUint8Type is a contructor for the representation of a uint64.
func NewUint8Type() Type {
	return &typeConstructor{
		signature:    "C",
		signatureIDL: "uint8",
		typeName:     jen.Uint8(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"WriteUint8").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"ReadUint8").Call(jen.Id(reader))
		},
		reader: constReader(1),
	}
}

// NewInt16Type is a contructor for the representation of a uint64.
func NewInt16Type() Type {
	return &typeConstructor{
		signature:    "w",
		signatureIDL: "int16",
		typeName:     jen.Int16(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"WriteInt16").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"ReadInt16").Call(jen.Id(reader))
		},
		reader: constReader(2),
	}
}

// NewUint16Type is a contructor for the representation of a uint64.
func NewUint16Type() Type {
	return &typeConstructor{
		signature:    "W",
		signatureIDL: "uint16",
		typeName:     jen.Uint16(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"WriteUint16").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"ReadUint16").Call(jen.Id(reader))
		},
		reader: constReader(2),
	}
}

// NewIntType is a contructor for the representation of an int32.
func NewIntType() Type {
	return &typeConstructor{
		signature:    "i",
		signatureIDL: "int32",
		typeName:     jen.Int32(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"WriteInt32").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"ReadInt32").Call(jen.Id(reader))
		},
		reader: constReader(4),
	}
}

// NewUIntType is a contructor for the representation of an uint32.
func NewUintType() Type {
	return &typeConstructor{
		signature:    "I",
		signatureIDL: "uint32",
		typeName:     jen.Uint32(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"WriteUint32").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"ReadUint32").Call(jen.Id(reader))
		},
		reader: constReader(4),
	}
}

// NewLongType is a contructor for the representation of a uint64.
func NewLongType() Type {
	return &typeConstructor{
		signature:    "l",
		signatureIDL: "int64",
		typeName:     jen.Int64(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"WriteInt64").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"ReadInt64").Call(jen.Id(reader))
		},
		reader: constReader(8),
	}
}

// NewULongType is a contructor for the representation of a uint64.
func NewULongType() Type {
	return &typeConstructor{
		signature:    "L",
		signatureIDL: "uint64",
		typeName:     jen.Uint64(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"WriteUint64").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"ReadUint64").Call(jen.Id(reader))
		},
		reader: constReader(8),
	}
}

// NewFloatType is a contructor for the representation of a float32.
func NewFloatType() Type {
	return &typeConstructor{
		signature:    "f",
		signatureIDL: "float32",
		typeName:     jen.Float32(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"WriteFloat32").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"ReadFloat32").Call(jen.Id(reader))
		},
		reader: constReader(4),
	}
}

// NewDoubleType is a contructor for the representation of a float32.
func NewDoubleType() Type {
	return &typeConstructor{
		signature:    "d",
		signatureIDL: "float64",
		typeName:     jen.Float64(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"WriteFloat64").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"ReadFloat64").Call(jen.Id(reader))
		},
		reader: constReader(8),
	}
}

// NewStringType is a contructor for the representation of a string.
func NewStringType() Type {
	return &typeConstructor{
		signature:    "s",
		signatureIDL: "str",
		typeName:     jen.String(),
		marshal: func(id string, writer string) *Statement {
			return jen.Id("basic.WriteString").Call(jen.Id(id),
				jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic",
				"ReadString").Call(jen.Id(reader))
		},
		reader: stringReader{},
	}
}

// NewVoidType is a contructor for the representation of the
// absence of a return type. Only used in the context of a returned
// type.
func NewVoidType() Type {
	return &typeConstructor{
		signature:    "v",
		signatureIDL: "nothing",
		typeName:     jen.Empty(),
		marshal: func(id string, writer string) *Statement {
			return jen.Nil()
		},
		unmarshal: func(reader string) *Statement {
			return jen.Empty()
		},
		reader: constReader(0),
	}
}

// NewValueType is a contructor for the representation of a Value.
func NewValueType() Type {
	return &typeConstructor{
		signature:    "m",
		signatureIDL: "any",
		typeName:     jen.Qual("github.com/lugu/qiloop/type/value", "Value"),
		marshal: func(id string, writer string) *Statement {
			return jen.Id(id).Dot("Write").Call(jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/value", "NewValue").Call(jen.Id(reader))
		},
		reader: valueReader{},
	}
}

// NewBoolType is a contructor for the representation of a bool.
func NewBoolType() Type {
	return &typeConstructor{
		signature:    "b",
		signatureIDL: "bool",
		typeName:     jen.Bool(),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic", "WriteBool").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/basic", "ReadBool").Call(jen.Id(reader))
		},
		reader: constReader(1),
	}
}

// NewMetaObjectType is a contructor for the representation of an
// object.
func NewMetaObjectType() Type {
	reader, err := MakeReader(MetaObjectSignature)
	if err != nil {
		panic(fmt.Errorf("invalid MetaObjectSignature: %s",
			MetaObjectSignature))
	}
	return &typeConstructor{
		signature:    MetaObjectSignature,
		signatureIDL: "MetaObject",
		typeName:     jen.Qual("github.com/lugu/qiloop/type/object", "MetaObject"),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/object", "WriteMetaObject").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/object", "ReadMetaObject").Call(jen.Id(reader))
		},
		reader: reader,
	}
}

// NewObjectType is a contructor for the representation of a Value.
func NewObjectType() Type {
	sig := "(({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>II)<ObjectReference,metaObject,serviceID,objectID>"

	reader, _ := MakeReader(sig)
	return &typeConstructor{
		signature:    "o",
		signatureIDL: "obj",
		typeName:     jen.Qual("github.com/lugu/qiloop/type/object", "ObjectReference"),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/object", "WriteObjectReference").Call(jen.Id(id), jen.Id(writer))
		},
		unmarshal: func(reader string) *Statement {
			return jen.Qual("github.com/lugu/qiloop/type/object", "ReadObjectReference").Call(jen.Id(reader))
		},
		reader: reader,
	}
}

func TypeIsObjectReference(t Type) bool {
	if t.Signature() != "o" {
		return false
	}
	name := t.TypeName()
	if reflect.DeepEqual(name,
		jen.Qual("github.com/lugu/qiloop/type/object",
		"ObjectReference"),
	) {
		return true
	}
	return false
}

// NewUnknownType is a contructor for an unknown type.
func NewUnknownType() Type {
	return &typeConstructor{
		signature:    "X",
		signatureIDL: "unknown",
		typeName:     jen.Id("interface{}"),
		marshal: func(id string, writer string) *Statement {
			return jen.Qual("fmt", "Errorf").Call(
				jen.Lit("marshal unknown type: %v"),
				jen.Id(id),
			)
		},
		unmarshal: func(reader string) *Statement {
			return jen.List(jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("unmarshal unknown type: %v"),
					jen.Id(reader),
				),
			)
		},
		reader: UnknownReader("X"),
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

// TypeDeclaration writes the type declaration into file.
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
            return fmt.Errorf("write slice size: %s", err)
        }`),
		jen.For(
			jen.Id("_, v := range "+listID),
		).Block(
			jen.Err().Op("=").Add(l.value.Marshal("v", writer)),
			jen.Id(`if (err != nil) {
                return fmt.Errorf("write slice value: %s", err)
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
			jen.Return(jen.Id("b"), jen.Qual("fmt", "Errorf").Call(jen.Id(`"read slice size: %s", err`)))),
		jen.Id("b").Op("=").Id("make").Call(l.TypeName(), jen.Id("size")),
		jen.For(
			jen.Id("i := 0; i < int(size); i++"),
		).Block(
			jen.Id("b[i], err =").Add(l.value.Unmarshal(reader)),
			jen.Id(`if (err != nil) {
                return b, fmt.Errorf("read slice value: %s", err)
            }`),
		),
		jen.Return(jen.Id("b"), jen.Nil()),
	).Call()
}

// Reader returns a list TypeReader.
func (l *ListType) Reader() TypeReader {
	return varReader{
		reader: l.value.Reader(),
	}
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

// TypeDeclaration writes the type declaration into file.
func (m *MapType) TypeDeclaration(file *jen.File) {
	return
}

// Reader returns a map TypeReader.
func (m *MapType) Reader() TypeReader {
	return varReader{
		reader: tupleReader([]memberReader{
			memberReader{
				name:   "key",
				reader: m.key.Reader(),
			},
			memberReader{
				name:   "value",
				reader: m.value.Reader(),
			},
		}),
	}
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
            return fmt.Errorf("write map size: %s", err)
        }`),
		jen.For(
			jen.Id("k, v := range "+mapID),
		).Block(
			jen.Err().Op("=").Add(m.key.Marshal("k", writer)),
			jen.Id(`if (err != nil) {
                return fmt.Errorf("write map key: %s", err)
            }`),
			jen.Err().Op("=").Add(m.value.Marshal("v", writer)),
			jen.Id(`if (err != nil) {
                return fmt.Errorf("write map value: %s", err)
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
			jen.Return(jen.Id("m"), jen.Qual("fmt", "Errorf").Call(jen.Id(`"read map size: %s", err`)))),
		jen.Id("m").Op("=").Id("make").Call(m.TypeName(), jen.Id("size")),
		jen.For(
			jen.Id("i := 0; i < int(size); i++"),
		).Block(
			jen.Id("k, err :=").Add(m.key.Unmarshal(reader)),
			jen.Id(`if (err != nil) {
                return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
            }`),
			jen.Id("v, err :=").Add(m.value.Unmarshal(reader)),
			jen.Id(`if (err != nil) {
                return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
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
	Name string
	Type Type
}

// Title is the public name of the field.
func (m MemberType) Title() string {
	return CleanName(m.Name)
}

// NewTupleType is a contructor for the representation of a series of
// types. Used to describe a method parameters list.
func NewTupleType(values []Type) *TupleType {
	var tuple TupleType
	tuple.Members = make([]MemberType, 0)
	for i, v := range values {
		tuple.Members = append(tuple.Members,
			MemberType{fmt.Sprintf("P%d", i), v})
	}
	return &tuple
}

// TupleType a list of a parameter of a method.
type TupleType struct {
	Members []MemberType
}

// Signature returns "(<signature 1><signature 2>...)" where
// <signature X> is the signature of the elements.
func (t *TupleType) Signature() string {
	sig := "("
	for _, m := range t.Members {
		sig += m.Type.Signature()
	}
	sig += ")"
	return sig
}

// ParamIDL returns "name1 signature1, name2 signature2" where
// signatureX is the signature of the elements.
func (t *TupleType) ParamIDL() string {
	var sig string
	for i, typ := range t.Members {
		sig += typ.Name + ": " + typ.Type.SignatureIDL()
		if i != len(t.Members)-1 {
			sig += ", "
		}
	}
	return sig
}

// SignatureIDL returns "Tuple<signature1,signature2>" where
// signatureX is the signature of the elements.
func (t *TupleType) SignatureIDL() string {
	var sig string
	for i, typ := range t.Members {
		sig += typ.Type.SignatureIDL()
		if i != len(t.Members)-1 {
			sig += ","
		}
	}
	return fmt.Sprintf("Tuple<%s>", sig)
}

// Params returns a statement representing the list of parameter of
// a method.
func (t *TupleType) Params() *Statement {
	arguments := make([]jen.Code, len(t.Members))
	for i, m := range t.Members {
		arguments[i] = jen.Id(m.Name).Add(m.Type.TypeName())
	}
	return jen.Params(arguments...)
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (t *TupleType) TypeName() *Statement {
	params := make([]jen.Code, 0)
	for _, typ := range t.Members {
		params = append(params, jen.Id(strings.Title(typ.Name)).Add(typ.Type.TypeName()))
	}
	return jen.Struct(params...)
}

// RegisterTo adds the type to the TypeSet.
func (t *TupleType) RegisterTo(s *TypeSet) {
	for _, m := range t.Members {
		m.Type.RegisterTo(s)
	}
	return
}

// TypeDeclaration writes the type declaration into file.
func (t *TupleType) TypeDeclaration(*jen.File) {
	return
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (t *TupleType) Marshal(tupleID string, writer string) *Statement {
	statements := make([]jen.Code, 0)
	for _, typ := range t.Members {
		s1 := jen.Err().Op("=").Add(typ.Type.Marshal(tupleID+"."+strings.Title(typ.Name), writer))
		s2 := jen.Id(`if (err != nil) {
			return fmt.Errorf("write tuple member: %s", err)
		}`)
		statements = append(statements, s1)
		statements = append(statements, s2)
	}
	statements = append(statements, jen.Return(jen.Nil()))
	return jen.Func().Params().Params(jen.Error()).Block(
		statements...,
	).Call()
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (t *TupleType) Unmarshal(reader string) *Statement {
	statements := make([]jen.Code, 0)
	for _, typ := range t.Members {
		s1 := jen.List(jen.Id("s."+strings.Title(typ.Name)), jen.Err()).Op("=").Add(typ.Type.Unmarshal(reader))
		s2 := jen.Id(`if (err != nil) {
			return s, fmt.Errorf("read tuple member: %s", err)
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

// Reader returns a map TypeReader.
func (t *TupleType) Reader() TypeReader {
	readers := make([]memberReader, len(t.Members))
	for i, m := range t.Members {
		readers[i] = memberReader{
			name:   m.Name,
			reader: m.Type.Reader(),
		}
	}
	return tupleReader(readers)
}

// ConvertMetaObjects replace any element type which has the same
// signature as MetaObject with an element of the type
// object.MetaObject. This is required to generate proxy services
// which implements the object.Object interface and avoid a circular
// dependency.
func (t *TupleType) ConvertMetaObjects() {
	for i, m := range t.Members {
		if m.Type.Signature() == MetaObjectSignature {
			t.Members[i].Type = NewMetaObjectType()
		}
	}
}

// NewStructType is a contructor for the representation of a struct.
func NewStructType(name string, members []MemberType) *StructType {
	return &StructType{name, members}
}

// StructType represents a struct.
type StructType struct {
	Name    string
	Members []MemberType
}

// Signature returns the signature of the struct.
func (s *StructType) Signature() string {
	if len(s.Members) == 0 {
		return fmt.Sprintf("()<%s>", s.Name)
	}
	types := ""
	names := make([]string, 0, len(s.Members))
	for _, v := range s.Members {
		names = append(names, v.Name)
		types += v.Type.Signature()
	}
	return fmt.Sprintf("(%s)<%s,%s>", types,
		s.Name, strings.Join(names, ","))
}

func (s *StructType) name() string {
	return CleanName(s.Name)
}

// SignatureIDL returns the idl signature of the struct.
func (s *StructType) SignatureIDL() string {
	return s.Name
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (s *StructType) TypeName() *Statement {
	return jen.Id(s.name())
}

// RegisterTo adds the type to the TypeSet.
// Structure name is updated if needed  after collision
// resolution. Register search if a type is in the TypeSet under a
// given name. If the same name and signature is already present it
// does nothing. otherwise it adds the type and search for a new which
// does not conflict with the names already present.
func (s *StructType) RegisterTo(set *TypeSet) {
	for _, v := range s.Members {
		v.Type.RegisterTo(set)
	}
	s.Name = set.ResolveCollision(s.Name, s.Signature())
	if set.Search(s.Name) == nil {
		set.Types = append(set.Types, s)
		set.Names = append(set.Names, s.Name)
	}
}

// TypeDeclaration writes the type declaration into file.
func (s *StructType) TypeDeclaration(file *jen.File) {
	fields := make([]jen.Code, len(s.Members))
	for i, v := range s.Members {
		fields[i] = jen.Id(v.Title()).Add(v.Type.TypeName())
	}
	file.Commentf("%s is serializable", s.name())
	file.Type().Id(s.name()).Struct(fields...)

	readFields := make([]jen.Code, len(s.Members)+1)
	writeFields := make([]jen.Code, len(s.Members)+1)
	for i, v := range s.Members {
		readFields[i] = jen.If(
			jen.Id("s."+v.Title()+", err =").Add(v.Type.Unmarshal("r")),
			jen.Id("err != nil")).Block(
			jen.Return(jen.Id("s"),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit(`read `+v.Title()+` field: %s`),
					jen.Id("err"),
				)),
		)
		writeFields[i] = jen.If(
			jen.Id("err :=").Add(v.Type.Marshal("s."+v.Title(), "w")),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Qual("fmt", "Errorf").Call(
				jen.Lit(`write `+v.Title()+` field: %s`),
				jen.Id("err"),
			)),
		)
	}
	readFields[len(s.Members)] = jen.Return(jen.Id("s"), jen.Nil())
	writeFields[len(s.Members)] = jen.Return(jen.Nil())

	file.Commentf("read%s unmarshalls %s", s.name(), s.name())
	file.Func().Id("read"+s.name()).Params(
		jen.Id("r").Id("io.Reader"),
	).Params(
		jen.Id("s").Id(s.name()), jen.Err().Error(),
	).Block(readFields...)
	file.Commentf("write%s marshalls %s", s.name(), s.name())
	file.Func().Id("write"+s.name()).Params(
		jen.Id("s").Id(s.name()),
		jen.Id("w").Qual("io", "Writer"),
	).Params(jen.Err().Error()).Block(writeFields...)
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (s *StructType) Marshal(structID string, writer string) *Statement {
	return jen.Id("write"+s.name()).Call(jen.Id(structID), jen.Id(writer))
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (s *StructType) Unmarshal(reader string) *Statement {
	return jen.Id("read" + s.name()).Call(jen.Id(reader))
}

// Reader returns a struct TypeReader.
func (s *StructType) Reader() TypeReader {
	readers := make([]memberReader, len(s.Members))
	for i, m := range s.Members {
		readers[i] = memberReader{
			name:   m.Name,
			reader: m.Type.Reader(),
		}
	}
	return tupleReader(readers)
}

// EnumType represents a const.
type EnumType struct {
	Name   string
	Values map[string]int
}

// EnumMember is used during Enum parsing
type EnumMember struct {
	Const string
	Value int
}

// NewEnumType returns an enum type
func NewEnumType(name string, values map[string]int) Type {
	return &EnumType{
		Name:   name,
		Values: values,
	}
}

// Signature returns "i" since enum are integer constant.
func (e *EnumType) Signature() string {
	return "i"
}

// SignatureIDL returns the name of the enum.
func (e *EnumType) SignatureIDL() string {
	return e.Name
}

// TypeName returns a statement to be inserted when the type is to be
// declared.
func (e *EnumType) TypeName() *Statement {
	return jen.Id(e.Name)
}

// RegisterTo adds the enum to the type set.
func (e *EnumType) RegisterTo(set *TypeSet) {
	// do not register anonymous enum
	if e.Name == "" {
		return
	}
	for i, name := range set.Names {
		if name == e.Name {
			if set.Types[i].Signature() != e.Signature() {
				log.Printf("type set collision %s: %s vs %s",
					name, set.Types[i].Signature(),
					e.Signature())
			}
			return
		}
	}
	set.Names = append(set.Names, e.Name)
	set.Types = append(set.Types, e)
}

// TypeDeclaration writes the type declaration into file.
func (e *EnumType) TypeDeclaration(file *jen.File) {
	file.Type().Id(e.Name).Int()
	var defs = make([]jen.Code, 0)
	for i, v := range e.Values {
		defs = append(defs, jen.Id(CleanName(i)).Op("=").Lit(v))
	}
	file.Const().Defs(defs...)
}

// Marshal returns a statement which represent the code needed to put
// the variable "id" into the io.Writer "writer" while returning an
// error.
func (e *EnumType) Marshal(id string, writer string) *Statement {
	return NewIntType().Marshal(id, writer)
}

// Unmarshal returns a statement which represent the code needed to read
// from a reader "reader" of type io.Reader and returns both the value
// read and an error.
func (e *EnumType) Unmarshal(reader string) *Statement {
	return NewIntType().Unmarshal(reader)
}

// Reader returns an enum TypeReader.
func (e *EnumType) Reader() TypeReader {
	return constReader(4)
}
