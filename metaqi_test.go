package metaqi

import "testing"

func testUtil(t *testing.T, input, expected string) {
    result := parse(input)
    if (result != expected) {
        t.Error(result)
    }
}

func TestParseInt(t *testing.T) {
    testUtil(t, "I", "uint32")
}

func TestParseString(t *testing.T) {
    testUtil(t, "s", "string")
}

func TestParseMultipleString(t *testing.T) {
    testUtil(t, "ss", "string")
}

func TestParseEmpty(t *testing.T) {
    testUtil(t, "", "not recognized")
}

func TestParseMap(t *testing.T) {
    testUtil(t, "{ss}", "MapType")
    testUtil(t, "{sI}", "MapType")
    testUtil(t, "{Is}", "MapType")
    testUtil(t, "{II}", "MapType")
}

func TestParseEmbedded(t *testing.T) {
    testUtil(t, "[s]", "EmbeddedType")
    testUtil(t, "[I]", "EmbeddedType")
    testUtil(t, "[{ss}]", "EmbeddedType")
    testUtil(t, "[{Is}]", "EmbeddedType")
}

func TestParseDefinition(t *testing.T) {
    testUtil(t, "(s)<test,a>", "TypeDefinition")
    testUtil(t, "(s)<test,a>", "TypeDefinition")
    testUtil(t, "(ss)<test,a,a>", "TypeDefinition")
    testUtil(t, "(sss)<test,a,a,a>", "TypeDefinition")
}

func TestParseMetaSignal(t *testing.T) {
    testUtil(t, "(Iss)<MetaSignal,uid,name,signature>", "TypeDefinition")
}
func TestParseMetaProperty(t *testing.T) {
    testUtil(t, "(Iss)<MetaProperty,uid,name,signature>" , "TypeDefinition")
}
func TestParseMetaMethodParameter(t *testing.T) {
    testUtil(t, "(ss)<MetaMethodParameter,name,description>", "TypeDefinition")
}
func TestParseMetaMethod(t *testing.T) {
    testUtil(t, "(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>", "TypeDefinition")
}
func TestParseMetaSignalMap2(t *testing.T) {
    testUtil(t, "{(Iss)<MetaSignal,uid,name,signature>I}", "MapType")
}
func TestParseMetaSignalMap(t *testing.T) {
    testUtil(t, "{I(Iss)<MetaSignal,uid,name,signature>}", "MapType")
}
func TestParseMetaPropertyMap(t *testing.T) {
    testUtil(t, "{I(Iss)<MetaProperty,uid,name,signature>}" , "MapType")
}
func TestParseMetaMethodMap(t *testing.T) {
    testUtil(t, "{I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}", "MapType")
}
func TestParseMetaObject(t *testing.T) {
    testUtil(t, "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>", "TypeDefinition")
}
