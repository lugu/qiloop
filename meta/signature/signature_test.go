package signature

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	parsec "github.com/prataprc/goparsec"
)

func testUtil(t *testing.T, input string, expected Type) {
	result, err := Parse(input)
	if err != nil {
		t.Error(err)
	} else if result == nil {
		t.Error("wrong return")
	} else if strings.ToLower(result.Signature()) != strings.ToLower(expected.Signature()) {
		buf := bytes.NewBufferString("")
		result.TypeName().Render(buf)
		t.Errorf("invalid type: %s (%s)", buf.String(), result.Signature())
	}
}

func testSignature(t *testing.T, signature string) {
	result, err := Parse(signature)
	if err != nil {
		t.Error(err)
	} else if result == nil {
		t.Error("wrong return")
	} else if result.Signature() != signature {
		t.Errorf("invalid signature: %s for %s",
			result.Signature(), signature)
	}
}

func helpStructName(t *testing.T, input string) {
	text := []byte(input)
	root, _ := structName()(parsec.NewScanner(text))
	if root == nil {
		t.Errorf("parse signature: %s", input)
	}
	terminal, ok := root.(*parsec.Terminal)
	if !ok {
		t.Errorf("parse signature: %s: %+v",
			input, reflect.TypeOf(root))
	}
	if input != terminal.GetValue() {
		t.Errorf("parse signature name: %s instead of %s",
			terminal.GetValue(), input)
	}
}

func TestStructName(t *testing.T) {
	helpStructName(t, "random_name")
	helpStructName(t, "RandomName2")
	helpStructName(t, "RandomName3<float>")
	helpStructName(t, "a<a>")
}

func TestParseBasics(t *testing.T) {
	testUtil(t, "v", NewVoidType())
	testUtil(t, "c", NewInt8Type())
	testUtil(t, "C", NewUint8Type())
	testUtil(t, "w", NewInt16Type())
	testUtil(t, "W", NewUint16Type())
	testUtil(t, "i", NewIntType())
	testUtil(t, "I", NewUIntType())
	testUtil(t, "s", NewStringType())
	testUtil(t, "L", NewULongType())
	testUtil(t, "l", NewLongType())
	testUtil(t, "b", NewBoolType())
	testUtil(t, "f", NewFloatType())
	testUtil(t, "d", NewDoubleType())
	testUtil(t, "m", NewValueType())
	testUtil(t, "X", NewUnknownType())
}

func TestParseMultipleString(t *testing.T) {
	testUtil(t, "ss", NewStringType())
}

func TestParseEmpty(t *testing.T) {
	t.SkipNow()
	testUtil(t, "", nil)
}

func TestParseMap(t *testing.T) {
	testUtil(t, "{cw}", NewMapType(NewInt8Type(), NewInt16Type()))
	testUtil(t, "{Wc}", NewMapType(NewUint16Type(), NewUint8Type()))
	testUtil(t, "{ss}", NewMapType(NewStringType(), NewStringType()))
	testUtil(t, "{sI}", NewMapType(NewStringType(), NewUIntType()))
	testUtil(t, "{is}", NewMapType(NewIntType(), NewStringType()))
	testUtil(t, "{iI}", NewMapType(NewIntType(), NewUIntType()))
	testUtil(t, "{Li}", NewMapType(NewULongType(), NewIntType()))
	testUtil(t, "{sl}", NewMapType(NewStringType(), NewLongType()))
}

func TestParseList(t *testing.T) {
	testUtil(t, "[s]", NewListType(NewStringType()))
	testUtil(t, "[i]", NewListType(NewIntType()))
	testUtil(t, "[b]", NewListType(NewBoolType()))
	testUtil(t, "[{bi}]", NewListType(NewMapType(NewBoolType(), NewIntType())))
	testUtil(t, "{b[i]}", NewMapType(NewBoolType(), NewListType(NewIntType())))
}

func TestParseTuple(t *testing.T) {
	testUtil(t, "(s)", NewTupleType([]Type{NewStringType()}))
	testUtil(t, "(i)", NewTupleType([]Type{NewIntType()}))
	testUtil(t, "(ii)", NewTupleType([]Type{NewIntType(), NewIntType()}))
	testUtil(t, "(fbd)", NewTupleType([]Type{NewFloatType(), NewBoolType(), NewDoubleType()}))
}

func TestParseDefinition(t *testing.T) {
	testUtil(t, "()<test>", NewStructType("test", make([]MemberType, 0)))
	testUtil(t, "(s)<test,a>", NewStructType("test", []MemberType{NewMemberType("a", NewStringType())}))
	testUtil(t, "(ss)<test,a,a>", NewStructType("test", []MemberType{
		NewMemberType("a", NewStringType()),
		NewMemberType("a", NewStringType()),
	}))
	testUtil(t, "(sss)<test,a,a,a>", NewStructType("test", []MemberType{
		NewMemberType("a", NewStringType()),
		NewMemberType("a", NewStringType()),
		NewMemberType("a", NewStringType()),
	}))
}

func TestParseEmbeddedDefinition(t *testing.T) {
	testUtil(t, "([s])<test,a>", NewStructType("test", []MemberType{
		NewMemberType("a", NewListType(NewStringType()))}))
	testUtil(t, "({si})<test,a>", NewStructType("test", []MemberType{
		NewMemberType("a", NewMapType(NewStringType(), NewIntType()))}))
}

func TestParseMapMap(t *testing.T) {
	testSignature(t, "{{ii}i}")
	testSignature(t, "{i{ii}}")
	testSignature(t, "{{ss}{ii}}")
	testSignature(t, "{{{si}s}{ii}}")
}

func TestParseDefinitionSignature(t *testing.T) {
	testSignature(t, "(s)<test,a>")
	testSignature(t, "(si)<test,a,b>")
	testSignature(t, "(iii)<test,a,b,c>")
	testSignature(t, "(s{ii})<test,a,b>")
	testSignature(t, "({ss})<test,a>")
}

func TestParseEmbeddedDefinitionSignature(t *testing.T) {
	testSignature(t, "([(s)<test2,b>])<test,a>")
	testSignature(t, "(s[(s)<test2,b>])<test,a,b>")
	testSignature(t, "([(s)<test2,b>]s)<test,a,b>")
	testSignature(t, "(b(s)<b,one>)<a,one,two>")
	testSignature(t, "(b(o(iW)<c,one,two>)<b,one,two>m)<a,one,two,three>")
}

func TestParseMetaSignal(t *testing.T) {
	testSignature(t, "(Iss)<MetaSignal,uid,name,signature>")
}
func TestParseMetaProperty(t *testing.T) {
	testSignature(t, "(Iss)<MetaProperty,uid,name,signature>")
}
func TestParseMetaMethodParameter(t *testing.T) {
	testSignature(t, "(ss)<MetaMethodParameter,name,description>")
}
func TestParseMetaMethod(t *testing.T) {
	testSignature(t, "(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>")
}
func TestParseMetaSignalMap2(t *testing.T) {
	testSignature(t, "{(Iss)<MetaSignal,uid,name,signature>I}")
}
func TestParseMetaSignalMap(t *testing.T) {
	testSignature(t, "{I(Iss)<MetaSignal,uid,name,signature>}")
}
func TestParseMetaPropertyMap(t *testing.T) {
	testSignature(t, "{I(Iss)<MetaProperty,uid,name,signature>}")
}
func TestParseMetaMethodMap(t *testing.T) {
	testSignature(t, "{I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}")
}
func TestParseMetaObject(t *testing.T) {
	testSignature(t, "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>")
}
func TestParseServiceInfo(t *testing.T) {
	testSignature(t, "[(sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>]")
}

func TestParseTextProcessingContext(t *testing.T) {
	testSignature(t, "{sc}")
	testSignature(t, "({sc}fs)<AgentGrd,concepts,confidence,userId>")
}

func TestParseRobotFullState(t *testing.T) {
	testSignature(t, "(ff)<ValueConfidence<float>,value,confidence>")
}

func TestGenerateType(t *testing.T) {
	typ := NewBoolType()
	var buf bytes.Buffer
	GenerateType(typ, "package1", &buf)
}

func TestParseError(t *testing.T) {
	check := func(s string) {
		typ, err := Parse(s)
		if err == nil {
			t.Errorf("unexpecting parse success: %s, %#v",
				s, typ)
		}
	}
	_, err := Parse("o")
	if err != nil {
		t.Error(err)
	}
	check("")
	check("#c")
	check("#c")
	check("{")
	check("{}")
	check("{c}")
	check("{sss}")
	check("[")
	check("[]")
	check("[sb]")
}

func TestInternalFunc(t *testing.T) {
	err := nodifyBasicType(make([]Node, 0))
	if err == nil {
		t.Errorf("unexpected")
	}
}
