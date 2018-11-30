package signature_test

import (
	"github.com/dave/jennifer/jen"
	. "github.com/lugu/qiloop/meta/signature"
	"testing"
)

func helpTestBasics(t *testing.T, typ Type, signature, idl string,
	typName *Statement) {

	if typ.Signature() != signature {
		panic("invalid signature: " + typ.Signature())
	}
	if typ.SignatureIDL() != idl {
		panic("invalid idl: " + typ.SignatureIDL())
	}
	// TODO: compare typName and typ.TypeName()
	if typ.Marshal("a", "b") == nil {
		panic("marshall is nil")
	}
	if typ.Unmarshal("a") == nil {
		panic("unmarshall is nil")
	}
	typ.RegisterTo(NewTypeSet())
}

func TestBasicTypes(t *testing.T) {
	helpTestBasics(t, NewInt8Type(), "c", "int8", jen.Int8())
	helpTestBasics(t, NewUint8Type(), "C", "uint8", jen.Uint8())
	helpTestBasics(t, NewInt16Type(), "w", "int16", jen.Int16())
	helpTestBasics(t, NewUint16Type(), "W", "uint16", jen.Uint16())
	helpTestBasics(t, NewIntType(), "i", "int32", jen.Int32())
	helpTestBasics(t, NewUIntType(), "I", "uint32", jen.Uint32())
	helpTestBasics(t, NewLongType(), "l", "int64", jen.Int64())
	helpTestBasics(t, NewULongType(), "L", "uint64", jen.Uint64())
	helpTestBasics(t, NewFloatType(), "f", "float32", jen.Float32())
	helpTestBasics(t, NewDoubleType(), "d", "float64", jen.Float64())
	helpTestBasics(t, NewStringType(), "s", "str", jen.String())
	helpTestBasics(t, NewVoidType(), "v", "nothing", jen.Empty())
	helpTestBasics(t, NewValueType(), "m", "any",
		jen.Qual("github.com/lugu/qiloop/type/value", "Value"))
	helpTestBasics(t, NewBoolType(), "b", "bool", jen.Bool())
	helpTestBasics(t, NewObjectType(), "o", "obj",
		jen.Qual("github.com/lugu/qiloop/type/object", "ObjectReference"))
	helpTestBasics(t, NewMetaObjectType(), MetaObjectSignature, "MetaObject",
		jen.Qual("github.com/lugu/qiloop/type/object", "MetaObject"))
	helpTestBasics(t, NewUnknownType(), "X", "unknown", jen.Id("interface{}"))
}

func TestListType(t *testing.T) {
	helpTestBasics(t, NewListType(NewStringType()), "[s]", "Vec<str>",
		jen.Index().Add(jen.String()))
	helpTestBasics(t, NewMapType(NewStringType(), NewBoolType()), "{sb}", "Map<str,bool>",
		jen.Map(jen.String()).Add(jen.Bool()))
	helpTestBasics(t, NewTupleType([]Type{NewStringType(), NewBoolType()}), "(sb)",
		"P0: str, P1: bool",
		jen.Struct(jen.Id("P0").Add(jen.String()), jen.Id("P1").Add(jen.Bool())))
	helpTestBasics(t, NewStructType("test", []MemberType{{"a", NewIntType()}}),
		"(i)<test,a>", "test", jen.Id("test"))
}
