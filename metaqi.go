package metaqi

import (
    "fmt"
    parsec "github.com/prataprc/goparsec"
)

type (
    MetaMethodParameter struct {
        Name string
        Description string
    }

    MetaMethod struct {
        Uuid uint32
        ReturnSignature string
        Name string
        ParametersSignature string
        Description string
        Parameters MetaMethodParameter
        ReturnDescription string
    }

    MetaSignal struct {
        Uuid uint32
        Name string
        Signature string
    }

    MetaProperty struct {
        Uuid uint32
        Name string
        Signature string
    }

    MetaObject struct {
        Methods map[int] MetaMethod
        Signals map[int] MetaSignal
        Properties map[int] MetaProperty
        Description string
    }
)

func String() parsec.Parser {
    return parsec.Atom("s", "string")
}
func Int() parsec.Parser {
    return parsec.Atom("I", "uint32")
}
func BasicType(ast *parsec.AST) parsec.Parser {
    return ast.OrdChoice("BasicType", nil, Int(), String())
}
func TypeName(ast *parsec.AST) parsec.Parser {
    return parsec.Ident()
}

func parse(input string) string {
    text := []byte(input)
    ast := parsec.NewAST("signature", 100)

    var embeddedType parsec.Parser
    var mapType parsec.Parser
    var typeDefinition parsec.Parser

    var declarationType = ast.OrdChoice("DeclarationType", nil,
        BasicType(ast), &mapType, &embeddedType, &typeDefinition)

    embeddedType = ast.And("EmbeddedType", nil,
        parsec.Atom("[", "MapStart"),
        &declarationType,
        parsec.Atom("]", "MapClose"))

    var listType = ast.Kleene("ListType", nil, &declarationType)

    var typeMemberList = ast.Maybe("MaybeMemberList", nil,
        ast.Many("TypeMemberList", nil,
        TypeName(ast), parsec.Atom(",", "Separator")))

    typeDefinition = ast.And("TypeDefinition", nil,
        parsec.Atom("(", "TypeParameterStart"),
        &listType,
        parsec.Atom(")", "TypeParameterClose"),
        parsec.Atom("<", "TypeDefinitionStart"),
        TypeName(ast),
        parsec.Atom(",", "TypeName"),
        &typeMemberList,
        parsec.Atom(">", "TypeDefinitionClose"))

    mapType = ast.And("MapType", nil,
        parsec.Atom("{", "MapStart"),
        &declarationType, &declarationType,
        parsec.Atom("}", "MapClose"))

    var typeSignature = declarationType

    root, scanner := ast.Parsewith(typeSignature, parsec.NewScanner(text))
    if (root != nil) {
        return root.GetName()
    }
    fmt.Println(scanner.GetCursor())
    return "not recognized"
}
