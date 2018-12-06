package signature

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	parsec "github.com/prataprc/goparsec"
	"io"
	"reflect"
)

// MetaObjectSignature is the signature of MetaObject. It is used to
// generate the MetaObject struct which is used to generate the
// services. This step is referred as stage 1.
const MetaObjectSignature string = "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>"

// ObjectSignature is the signature of ObjectReference.
var ObjectSignature = fmt.Sprintf("(b%sIII)<ObjectReference,boolean,metaObject,parentID,serviceID,objectID>",
	MetaObjectSignature)

// Statement is a short version for jen.Statement.
type Statement = jen.Statement

func basicType() parsec.Parser {
	return parsec.OrdChoice(nodifyBasicType,
		parsec.Atom("I", "uint32"),
		parsec.Atom("i", "int32"),
		parsec.Atom("s", "string"),
		parsec.Atom("L", "uint64"),
		parsec.Atom("l", "int64"),
		parsec.Atom("b", "bool"),
		parsec.Atom("f", "float32"),
		parsec.Atom("d", "float64"),
		parsec.Atom("m", "value"),
		parsec.Atom("o", "github.com/lugu/qiloop/type/object.Object"),
		parsec.Atom("X", "interface{}"),
		parsec.Atom("v", "void"),
		parsec.Atom("c", "int8"),
		parsec.Atom("C", "uint8"),
		parsec.Atom("w", "int16"),
		parsec.Atom("W", "uint16"),
	)
}

func typeName() parsec.Parser {
	return parsec.Ident()
}

func structName() parsec.Parser {
	patterns := []string{
		// Allow C++ style name (ex: List<double>)
		`[A-Za-z][0-9a-zA-Z_]*\<[A-Za-z][0-9a-zA-Z_]*>`,
		`[A-Za-z][0-9a-zA-Z_]*`,
	}
	names := []string{
		"structTemplateName",
		"structName",
	}
	return parsec.OrdTokens(patterns, names)
}

// Node is an alias to parsec.ParsecNode
type Node = parsec.ParsecNode

func nodifyBasicType(nodes []Node) Node {
	if len(nodes) != 1 {
		return fmt.Errorf("wrong basic arguments %+v", nodes)
	}
	signature := nodes[0].(*parsec.Terminal).GetValue()
	switch signature {
	case "i":
		return NewIntType()
	case "I":
		return NewUIntType()
	case "l":
		return NewLongType()
	case "L":
		return NewULongType()
	case "s":
		return NewStringType()
	case "b":
		return NewBoolType()
	case "f":
		return NewFloatType()
	case "d":
		return NewDoubleType()
	case "v":
		return NewVoidType()
	case "m":
		return NewValueType()
	case "o":
		return NewObjectType()
	case "X":
		return NewUnknownType()
	case "c":
		return NewInt8Type()
	case "C":
		return NewUint8Type()
	case "w":
		return NewInt16Type()
	case "W":
		return NewUint16Type()
	default:
		return fmt.Errorf("wrong signature %s", signature)
	}
}

func extractValue(object interface{}) (Type, error) {
	nodes, ok := object.([]Node)
	if !ok {
		return nil, fmt.Errorf("extraction failed: %+v", reflect.TypeOf(object))
	}
	value, ok := nodes[0].(Type)
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
		return fmt.Errorf("key conversion failed: %s", err)
	}
	value, err := extractValue(nodes[2])
	if err != nil {
		return fmt.Errorf("value conversion failed: %s", err)
	}
	return NewMapType(key, value)
}

func nodifyArrayType(nodes []Node) Node {
	if len(nodes) != 3 {
		return fmt.Errorf("wrong arguments %+v", nodes)
	}
	value, err := extractValue(nodes[1])
	if err != nil {
		return fmt.Errorf("value conversion failed: %s", err)
	}
	return NewListType(value)
}

func extractMembersTypes(node Node) ([]Type, error) {
	typeList, ok := node.([]Node)
	if !ok {
		return nil, fmt.Errorf("member type list is not a list: %s", reflect.TypeOf(node))
	}
	types := make([]Type, len(typeList))
	for i := range typeList {
		memberType, err := extractValue(typeList[i])
		if err != nil {
			return nil, fmt.Errorf("member type value conversion failed: %s", err)
		}
		types[i] = memberType
	}
	return types, nil
}

func extractMembersName(node Node) ([]string, error) {
	membersList, ok := node.([]Node)
	if !ok {
		return nil, fmt.Errorf("member name list is not a list: %s", reflect.TypeOf(node))
	}
	names := make([]string, len(membersList))
	for i, n := range membersList {
		memberName, ok := n.(*parsec.Terminal)
		if !ok {
			return nil, fmt.Errorf("failed to convert member names %s", reflect.TypeOf(n))
		}
		names[i] = memberName.GetValue()
	}
	return names, nil
}

func extractMembers(typesNode, namesNode Node) ([]MemberType, error) {

	types, err := extractMembersTypes(typesNode)
	if err != nil {
		return nil, fmt.Errorf("failed to extract members types: %s", err)
	}
	names, err := extractMembersName(namesNode)
	if err != nil {
		return nil, fmt.Errorf("failed to extract members name: %s", err)
	}

	if len(types) != len(names) {
		return nil, fmt.Errorf("member types and names of different size: %+v, %+v", types, names)
	}

	members := make([]MemberType, len(names))

	for i := range types {
		members[i] = NewMemberType(names[i], types[i])
	}
	return members, nil
}

func nodifyTupleType(nodes []Node) Node {

	types, err := extractMembersTypes(nodes[1])
	if err != nil {
		return fmt.Errorf("failed to extract tuple member type: %s", err)
	}

	return NewTupleType(types)
}

func nodifyTypeMember(nodes []Node) Node {
	return nodes[1]
}

func nodifyStrucType(nodes []Node) Node {

	terminal, ok := nodes[4].(*parsec.Terminal)
	if !ok {
		return fmt.Errorf("wrong name %s", reflect.TypeOf(nodes[4]))
	}
	name := terminal.GetValue()
	members, err := extractMembers(nodes[1], nodes[5])
	if err != nil {
		return fmt.Errorf("failed to extract type definition: %s", err)
	}

	return NewStructType(name, members)
}

// Parse reads a signature contained in a string and constructs its
// type representation.
func Parse(input string) (Type, error) {
	text := []byte(input)

	var arrayType parsec.Parser
	var mapType parsec.Parser
	var structType parsec.Parser
	var tupleType parsec.Parser

	var declarationType = parsec.OrdChoice(nil,
		basicType(), &mapType, &arrayType, &structType, &tupleType)

	arrayType = parsec.And(nodifyArrayType,
		parsec.Atom("[", "MapStart"),
		&declarationType,
		parsec.Atom("]", "MapClose"))

	var listType = parsec.Kleene(nil, &declarationType)

	var typeMemberList = parsec.Kleene(
		nil, parsec.And(
			nodifyTypeMember,
			parsec.Atom(",", "TypeComa"),
			typeName(),
		))

	tupleType = parsec.And(nodifyTupleType,
		parsec.Atom("(", "TypeParameterStart"),
		&listType,
		parsec.Atom(")", "TypeParameterClose"))

	structType = parsec.And(nodifyStrucType,
		parsec.Atom("(", "TypeParameterStart"),
		&listType,
		parsec.Atom(")", "TypeParameterClose"),
		parsec.Atom("<", "TypeDefinitionStart"),
		structName(),
		&typeMemberList,
		parsec.Atom(">", "TypeDefinitionClose"))

	mapType = parsec.And(nodifyMap,
		parsec.Atom("{", "MapStart"),
		&declarationType, &declarationType,
		parsec.Atom("}", "MapClose"))

	var typeSignature = declarationType

	root, _ := typeSignature(parsec.NewScanner(text))
	if root == nil {
		return nil, fmt.Errorf("failed to parse signature: %s", input)
	}
	types, ok := root.([]Node)
	if !ok {
		err, ok := root.(error)
		if !ok {
			return nil, fmt.Errorf("failed to convert array: %+v", reflect.TypeOf(root))
		}
		return nil, err
	}
	if len(types) != 1 {
		return nil, fmt.Errorf("did not parse only one type: %+v", root)
	}
	constructor, ok := types[0].(Type)
	if !ok {
		return nil, fmt.Errorf("failed to convert value: %+v", reflect.TypeOf(types[0]))
	}
	return constructor, nil
}

// GenerateType generate the code required to serialize the given
// type.
func GenerateType(v Type, packageName string, w io.Writer) error {
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
