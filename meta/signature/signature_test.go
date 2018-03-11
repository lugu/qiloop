package signature

import (
	"bytes"
	"testing"
)

func testUtil(t *testing.T, input string, expected ValueConstructor) {
	result, err := Parse(input)
	if err != nil {
		t.Error(err)
	} else if result == nil {
		t.Error("wrong return")
	} else if Print(result) != Print(expected) {
		buf := bytes.NewBufferString("")
		result.TypeName().Render(buf)
		t.Error("invalid type: " + buf.String())
	}
}

func testSignature(t *testing.T, signature string) {
	result, err := Parse(signature)
	if err != nil {
		t.Error(err)
	} else if result == nil {
		t.Error("wrong return")
	} else if result.Signature() != signature {
		t.Error("invalid signature: " + result.Signature())
	}
}

func TestParseBasics(t *testing.T) {
	testUtil(t, "i", NewIntValue())
	testUtil(t, "I", NewIntValue())
	testUtil(t, "s", NewStringValue())
	testUtil(t, "L", NewLongValue())
	testUtil(t, "l", NewLongValue())
	testUtil(t, "b", NewBoolValue())
	testUtil(t, "f", NewFloatValue())
	testUtil(t, "m", NewValueValue())
	testUtil(t, "X", NewUnknownValue())
}

func TestParseMultipleString(t *testing.T) {
	testUtil(t, "ss", NewStringValue())
}

func TestParseEmpty(t *testing.T) {
	t.SkipNow()
	testUtil(t, "", nil)
}

func TestParseMap(t *testing.T) {
	testUtil(t, "{ss}", NewMapValue(NewStringValue(), NewStringValue()))
	testUtil(t, "{sI}", NewMapValue(NewStringValue(), NewIntValue()))
	testUtil(t, "{Is}", NewMapValue(NewIntValue(), NewStringValue()))
	testUtil(t, "{II}", NewMapValue(NewIntValue(), NewIntValue()))
	testUtil(t, "{LI}", NewMapValue(NewLongValue(), NewIntValue()))
	testUtil(t, "{sL}", NewMapValue(NewStringValue(), NewLongValue()))
}

func TestParseList(t *testing.T) {
	testUtil(t, "[s]", NewListValue(NewStringValue()))
	testUtil(t, "[I]", NewListValue(NewIntValue()))
	testUtil(t, "[b]", NewListValue(NewBoolValue()))
	testUtil(t, "[{bI}]", NewListValue(NewMapValue(NewBoolValue(), NewIntValue())))
	testUtil(t, "{b[I]}", NewMapValue(NewBoolValue(), NewListValue(NewIntValue())))
}

func TestParseDefinition(t *testing.T) {
	testUtil(t, "(s)<test,a>", NewStrucValue("test", []MemberValue{NewMemberValue("a", NewStringValue())}))
	testUtil(t, "(ss)<test,a,a>", NewStrucValue("test", []MemberValue{
		NewMemberValue("a", NewStringValue()),
		NewMemberValue("a", NewStringValue()),
	}))
	testUtil(t, "(sss)<test,a,a,a>", NewStrucValue("test", []MemberValue{
		NewMemberValue("a", NewStringValue()),
		NewMemberValue("a", NewStringValue()),
		NewMemberValue("a", NewStringValue()),
	}))
}

func TestParseEmbeddedDefinition(t *testing.T) {
	testUtil(t, "([s])<test,a>", NewStrucValue("test", []MemberValue{
		NewMemberValue("a", NewStringValue())}))
}

func TestParseMapMap(t *testing.T) {
	testSignature(t, "{{II}I}")
	testSignature(t, "{I{II}}")
	testSignature(t, "{{ss}{II}}")
	testSignature(t, "{{{sI}s}{II}}")
}

func TestParseDefinitionSignature(t *testing.T) {
	testSignature(t, "(s)<test,a>")
	testSignature(t, "(sI)<test,a,b>")
	testSignature(t, "(III)<test,a,b,c>")
	testSignature(t, "(s{II})<test,a,b>")
	testSignature(t, "({ss})<test,a>")
}

func TestParseEmbeddedDefinitionSignature(t *testing.T) {
	testSignature(t, "([(s)<test2,b>])<test,a>")
	testSignature(t, "(s[(s)<test2,b>])<test,a,b>")
	testSignature(t, "([(s)<test2,b>]s)<test,a,b>")
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
