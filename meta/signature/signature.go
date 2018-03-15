package signature

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	parsec "github.com/prataprc/goparsec"
	"io"
	"log"
	"reflect"
)

// MetaObjectSignature is the signature of MetaObject. It is used to
// generate the MetaObject struct which is used to generate the
// services. This step is referred as stage 1.
const MetaObjectSignature string = "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>"

var ObjectSignature string = fmt.Sprintf("(b%sIII)<ObjectReference,boolean,metaObject,parentID,serviceID,objectID>",
	MetaObjectSignature)

// Statement is a short version for jen.Statement.
type Statement = jen.Statement

func basicType() parsec.Parser {
	return parsec.OrdChoice(nodifyBasicType,
		parsec.Atom("I", "uint32"),
		parsec.Atom("i", "uint32"),
		parsec.Atom("s", "string"),
		parsec.Atom("L", "uint64"),
		parsec.Atom("l", "uint64"),
		parsec.Atom("b", "bool"),
		parsec.Atom("f", "float32"),
		parsec.Atom("d", "float64"),
		parsec.Atom("m", "value"),
		parsec.Atom("o", "github.com/lugu/qiloop/type/object.Object"),
		parsec.Atom("X", "interface{}"),
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
	default:
		log.Panicf("wrong signature %s", signature)
	}
	return nil
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
		log.Panicf("key conversion failed: %s", err)
	}
	value, err := extractValue(nodes[2])
	if err != nil {
		log.Panicf("value conversion failed: %s", err)
	}
	return NewMapType(key, value)
}

func nodifyArrayType(nodes []Node) Node {
	if len(nodes) != 3 {
		log.Panicf("wrong arguments %+v", nodes)
	}
	value, err := extractValue(nodes[1])
	if err != nil {
		log.Panicf("value conversion failed: %s", err)
	}
	return NewListType(value)
}

func extractMembersTypes(node Node) []Type {
	typeList, ok := node.([]Node)
	if !ok {
		log.Panicf("member type list is not a list: %s", reflect.TypeOf(node))
	}
	types := make([]Type, len(typeList))
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

func extractMembers(typesNode, namesNode Node) []MemberType {

	types := extractMembersTypes(typesNode)
	names := extractMembersName(namesNode)

	if len(types) != len(names) {
		log.Panicf("member types and names of different size: %+v, %+v", types, names)
	}

	members := make([]MemberType, len(names))

	for i := range types {
		members[i] = NewMemberType(names[i], types[i])
	}
	return members
}

func nodifyTupleType(nodes []Node) Node {

	types := extractMembersTypes(nodes[1])

	return NewTupleType(types)
}

func nodifyTypeDefinition(nodes []Node) Node {

	terminal, ok := nodes[4].(*parsec.Terminal)
	if !ok {
		log.Panicf("wrong name %s", reflect.TypeOf(nodes[4]))
	}
	name := terminal.GetValue()
	members := extractMembers(nodes[1], nodes[6])

	return NewStrucType(name, members)
}

// Parse reads a signature contained in a string and constructs its
// type representation.
func Parse(input string) (Type, error) {
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
		return nil, fmt.Errorf("failed to parse signature: %s", input)
	}
	types, ok := root.([]Node)
	if !ok {
		return nil, fmt.Errorf("failed to convert array: %+v", reflect.TypeOf(root))
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
