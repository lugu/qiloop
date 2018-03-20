package idl

import (
	"fmt"
	. "github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	parsec "github.com/prataprc/goparsec"
	"io"
	"io/ioutil"
	"reflect"
)

func basicType() parsec.Parser {
	return parsec.OrdChoice(nodifyBasicType,
		parsec.Atom("int32", ""),
		parsec.Atom("uint32", ""),
		parsec.Atom("int64", ""),
		parsec.Atom("uint64", ""),
		parsec.Atom("float32", ""),
		parsec.Atom("float64", ""),
		parsec.Atom("int64", ""),
		parsec.Atom("uint64", ""),
		parsec.Atom("bool", ""),
		parsec.Atom("str", ""),
		parsec.Atom("obj", ""),
		parsec.Atom("any", ""))
}

// required to break the recursive definition of typeParser()
var mapParser parsec.Parser
var vecParser parsec.Parser
var referenceParser parsec.Parser

func init() {
	mapParser = mapType()
	vecParser = vecType()
	referenceParser = referenceType()
}

func mapType() parsec.Parser {
	return parsec.And(
		nodifyMap,
		parsec.Atom("Map<", "Map<"),
		typeParser(),
		parsec.Atom(",", ","),
		typeParser(),
		parsec.Atom(">", ">"),
	)
}

func vecType() parsec.Parser {
	return parsec.And(
		nodifyVec,
		parsec.Atom("Vec<", "Vec<"),
		typeParser(),
		parsec.Atom(">", ">"),
	)
}

func typeParser() parsec.Parser {
	return parsec.OrdChoice(
		nodifyType,
		basicType(),
		&mapParser,
		&vecParser,
		&referenceParser,
	)
}

func comments() parsec.Parser {
	return parsec.And(
		nodifyComment,
		parsec.Maybe(
			nodifyMaybeComment,
			parsec.And(
				nodifyCommentContent,
				parsec.Atom("//", "//"),
				parsec.Token(`.*`, "comment"),
			),
		),
	)
}

func returns() parsec.Parser {
	return parsec.And(
		nodifyReturns,
		parsec.Maybe(
			nodifyMaybeReturns,
			parsec.And(
				nodifyReturnsType,
				parsec.Atom("->", "->"),
				typeParser(),
			),
		),
	)
}

func parameter() parsec.Parser {
	return parsec.And(
		nodifyParam,
		Ident(),
		parsec.Atom(":", ":"),
		typeParser(),
	)
}

func parameters() parsec.Parser {
	return parsec.And(
		nodifyAndParams,
		parsec.Maybe(
			nodifyMaybeParams,
			parsec.Many(
				nodifyParams,
				parameter(),
				parsec.Atom(",", ","),
			),
		),
	)
}

func Ident() parsec.Parser {
	return parsec.Token(`[_A-Za-z][0-9a-zA-Z_]*`, "IDENT")
}

func method() parsec.Parser {
	return parsec.And(
		nodifyMethod,
		parsec.Atom("fn", "fn"),
		Ident(),
		parsec.Atom("(", "("),
		parameters(),
		parsec.Atom(")", ")"),
		returns(),
		comments(),
	)
}

func signal() parsec.Parser {
	return parsec.And(
		nodifySignal,
		parsec.Atom("sig", "sig"),
		Ident(),
		parsec.Atom("(", "("),
		parameters(),
		parsec.Atom(")", ")"),
		comments(),
	)
}

func action() parsec.Parser {
	return parsec.OrdChoice(
		nodifyAction,
		method(),
		signal(),
	)
}

func interfaceParser() parsec.Parser {
	return parsec.And(
		nodifyInterface,
		parsec.Atom("interface", "interface"),
		Ident(),
		comments(),
		parsec.Kleene(nodifyActionList, action()),
		parsec.Atom("end", "end"),
		comments(),
	)
}

func referenceType() parsec.Parser {
	return parsec.And(
		nodifyTypeReference,
		Ident(),
	)
}

func member() parsec.Parser {
	return parsec.And(
		nodifyMember,
		Ident(),
		parsec.Atom(":", ":"),
		typeParser(),
		comments(),
	)
}

func structure() parsec.Parser {
	return parsec.And(
		nodifyStructure,
		parsec.Atom("struct", "struct"),
		Ident(),
		comments(),
		parsec.Kleene(nodifyMemberList, member()),
		parsec.Atom("end", "end"),
		comments(),
	)
}

func declaration() parsec.Parser {
	return parsec.OrdChoice(
		nodifyDeclaration,
		structure(),
		interfaceParser(),
	)
}

func declarations() parsec.Parser {
	return parsec.Many(
		nodifyDeclarationList,
		declaration(),
	)
}

func parse(input []byte) (*Declarations, error) {
	root, scanner := declarations()(parsec.NewScanner(input).TrackLineno())
	_, scanner = scanner.SkipWS()
	if !scanner.Endof() {
		return nil, fmt.Errorf("parsing error at line: %d", scanner.Lineno())
	}
	if root == nil {
		return nil, fmt.Errorf("cannot parse input:\n%s", input)
	}

	definitions, ok := root.(*Declarations)
	if !ok {
		if err, ok := root.(error); ok {
			return nil, err
		}
		return nil, fmt.Errorf("cannot parse IDL: %+v", reflect.TypeOf(root))
	}
	return definitions, nil
}

// ParseIDL read an IDL definition from a reader and returns the
// MetaObject associated with the IDL.
func ParseIDL(reader io.Reader) ([]object.MetaObject, error) {
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read input: %s", err)
	}

	unregsterTypeNames()
	// first pass needed to generate the TypeStruct correctly.
	definitions, err := parse(input)
	if err != nil {
		return nil, err
	}

	registerTypeNames(definitions)
	// second pass needed to generate the MetaObject correctly.
	definitions, err = parse(input)
	if err != nil {
		return nil, err
	}

	return definitions.Interfaces, nil
}

func nodifyDeclaration(nodes []Node) Node {
	return nodes[0]
}

func nodifyMaybeComment(nodes []Node) Node {
	return nodes[0]
}

func nodifyReturnsType(nodes []Node) Node {
	return nodes[1]
}

func nodifyAction(nodes []Node) Node {
	return nodes[0]
}

func nodifyType(nodes []Node) Node {
	return nodes[0]
}

func nodifyMaybeReturns(nodes []Node) Node {
	return nodes[0]
}

func nodifyMaybeParams(nodes []Node) Node {
	return nodes[0]
}

func nodifyComment(nodes []Node) Node {
	if _, ok := nodes[0].(parsec.MaybeNone); ok {
		return ""
	} else if comment, ok := nodes[0].(string); ok {
		return comment
	} else if uid, ok := nodes[0].(uint32); ok {
		return uid
	} else {
		return nodes[0]
	}
}

func nodifyCommentContent(nodes []Node) Node {
	comment := nodes[1].(*parsec.Terminal).GetValue()
	var uid uint32
	_, err := fmt.Sscanf(comment, "uid:%d", &uid)
	if err == nil {
		return uid
	}
	return comment
}

func nodifyReturns(nodes []Node) Node {
	if _, ok := nodes[0].(parsec.MaybeNone); ok {
		return NewVoidType()
	} else if ret, ok := nodes[0].(Type); ok {
		return ret
	} else {
		return fmt.Errorf("unexpected return value (%s): %v",
			reflect.TypeOf(nodes[0]), nodes[0])
	}
}

// Declarations represent the result of the IDL parsing. It contains a
// list of MetaObject and a list of StructType. The embedded
// MetaObjects and StructTypes can not be use: their signatures need to
// be re-unified.
type Declarations struct {
	Interfaces []object.MetaObject
	Struct     []StructType
}

// nodifyDeclarationList returns a Declarations structure holding a
// list of MetaObject and StructType.
func nodifyDeclarationList(nodes []Node) Node {
	interfaces := make([]object.MetaObject, 0)
	struc := make([]StructType, 0)
	for _, node := range nodes {
		if err, ok := node.(error); ok {
			return err
		}
		if metaObj, ok := node.(*object.MetaObject); ok {
			interfaces = append(interfaces, *metaObj)
		} else if s, ok := node.(*StructType); ok {
			struc = append(struc, *s)
		} else {
			return fmt.Errorf("Expecting MetaObject, got %+v: %+v", reflect.TypeOf(node), node)
		}
	}
	return &Declarations{
		Interfaces: interfaces,
		Struct:     struc,
	}
}

// nodifyTypeReference returns a Type or an error. The type is a
// RefType.
func nodifyTypeReference(nodes []Node) Node {
	typeNode := nodes[0]
	typeName := typeNode.(*parsec.Terminal).GetValue()
	if ref, err := NewRefType(typeName); err != nil {
		return err
	} else {
		return ref
	}
}

// nodifyMember returns a MemberType or an error.
func nodifyMember(nodes []Node) Node {
	nameNode := nodes[0]
	typeNode := nodes[2]
	var member MemberType
	var ok bool
	member.Name = nameNode.(*parsec.Terminal).GetValue()
	member.Value, ok = typeNode.(Type)
	if !ok {
		return fmt.Errorf("Expecting Type, got %+v: %+v", reflect.TypeOf(typeNode), typeNode)
	}
	return member
}

// nodifyMemberList returns a StructType or an error.
func nodifyMemberList(nodes []Node) Node {
	members := make([]MemberType, len(nodes))
	for i, node := range nodes {
		if member, ok := node.(MemberType); ok {
			members[i] = member
		} else {
			return fmt.Errorf("MemberList: unexpected type, got %+v: %+v", reflect.TypeOf(node), node)
		}
	}
	return NewStructType("parameters", members)
}

// nodifyStructure returns a StructType of an error.
func nodifyStructure(nodes []Node) Node {
	nameNode := nodes[1]
	structNode := nodes[3]
	structType, ok := structNode.(*StructType)
	if !ok {
		return fmt.Errorf("unexpected non struct type(%s): %v",
			reflect.TypeOf(structNode), structNode)
	}
	structType.Name = nameNode.(*parsec.Terminal).GetValue()
	return structType
}

func nodifyInterface(nodes []Node) Node {
	nameNode := nodes[1]
	objNode := nodes[3]
	metaObj, ok := objNode.(*object.MetaObject)
	if !ok {
		return fmt.Errorf("Expecting MetaObject, got %+v: %+v", reflect.TypeOf(objNode), objNode)
	}
	metaObj.Description = nameNode.(*parsec.Terminal).GetValue()
	return metaObj
}

func nodifyVec(nodes []Node) Node {

	elementNode := nodes[1]
	elementType, ok := elementNode.(Type)
	if !ok {
		return fmt.Errorf("invalid vector element: %v", elementNode)
	}
	return NewListType(elementType)
}

func nodifyMap(nodes []Node) Node {

	keyNode := nodes[1]
	valueNode := nodes[3]
	keyType, ok := keyNode.(Type)
	if !ok {
		return fmt.Errorf("invalid map key: %v", keyNode)
	}
	valueType, ok := valueNode.(Type)
	if !ok {
		return fmt.Errorf("invalid map value: %v", valueNode)
	}
	return NewMapType(keyType, valueType)
}

func nodifyBasicType(nodes []Node) Node {
	if len(nodes) != 1 {
		return fmt.Errorf("basic type array size: %d: %s", len(nodes), nodes)
	}
	sig := nodes[0].(*parsec.Terminal).GetValue()
	switch sig {
	case "int32":
		return NewIntType()
	case "uint32":
		return NewUIntType()
	case "int64":
		return NewLongType()
	case "uint64":
		return NewULongType()
	case "float32":
		return NewFloatType()
	case "float64":
		return NewDoubleType()
	case "str":
		return NewStringType()
	case "bool":
		return NewBoolType()
	case "any":
		return NewValueType()
	case "obj":
		return NewObjectType()
	default:
		return fmt.Errorf("unknown type: %s", sig)
	}
}

// nodifyMethod returns either a MetaMethod or an error.
func nodifyMethod(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
		return fmt.Errorf("failed to parse method: %s", err)
	}
	retNode := nodes[5]
	paramNode := nodes[3]
	commentNode := nodes[6]
	structType, ok := paramNode.(*StructType)
	if !ok {
		return fmt.Errorf("failed to convert param type (%s): %v", reflect.TypeOf(paramNode), paramNode)
	}
	retType, ok := retNode.(Type)
	if !ok {
		return fmt.Errorf("failed to convert return type (%s): %v", reflect.TypeOf(retNode), retNode)
	}

	params := make([]Type, len(structType.Members))
	for i, member := range structType.Members {
		params[i] = member.Value
	}
	tupleType := NewTupleType(params)

	method := new(object.MetaMethod)
	method.Name = nodes[1].(*parsec.Terminal).GetValue()
	method.ParametersSignature = tupleType.Signature()
	if len(structType.Members) != 0 {
		method.Parameters = make([]object.MetaMethodParameter, len(structType.Members))

		for i, member := range structType.Members {
			method.Parameters[i].Name = member.Name
		}
	}
	if uid, ok := commentNode.(uint32); ok {
		method.Uid = uid
	}
	method.ReturnSignature = retType.Signature()
	return method
}

// nodifySignal returns either a MetaSignal or an error.
func nodifySignal(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
		return fmt.Errorf("failed to parse method: %s", err)
	}
	commentNode := nodes[5]
	paramNode := nodes[3]
	structType, ok := paramNode.(*StructType)
	if !ok {
		return fmt.Errorf("failed to convert param type (%s): %v", reflect.TypeOf(paramNode), paramNode)
	}

	params := make([]Type, len(structType.Members))
	for i, member := range structType.Members {
		params[i] = member.Value
	}
	tupleType := NewTupleType(params)

	signal := new(object.MetaSignal)
	signal.Name = nodes[1].(*parsec.Terminal).GetValue()
	signal.Signature = tupleType.Signature()
	if uid, ok := commentNode.(uint32); ok {
		signal.Uid = uid
	}
	return signal
}

// nodifyActionList returns either a MetaObject or an error.
func nodifyActionList(nodes []Node) Node {
	metaObj := new(object.MetaObject)
	methods := make(map[uint32]object.MetaMethod)
	signals := make(map[uint32]object.MetaSignal)
	var customAction uint32 = 100
	for _, node := range nodes {
		if err, ok := node.(error); ok {
			return err
		}
		if method, ok := node.(*object.MetaMethod); ok {
			if method.Uid == 0 && method.Name != "registerEvent" {
				method.Uid = customAction
				customAction++
			}
			methods[method.Uid] = *method
		} else if signal, ok := node.(*object.MetaSignal); ok {
			if signal.Uid == 0 {
				signal.Uid = customAction
				customAction++
			}
			signals[signal.Uid] = *signal
		} else {
			return fmt.Errorf("Expecting MetaMethod or MetaSignal, got %+v: %+v", reflect.TypeOf(node), node)
		}
	}
	if len(methods) != 0 {
		metaObj.Methods = methods
	}
	if len(signals) != 0 {
		metaObj.Signals = signals
	}
	return metaObj
}

func checkError(nodes []Node) (error, bool) {
	if len(nodes) == 1 {
		if err, ok := nodes[0].(error); ok {
			return err, true
		}
	}
	return nil, false
}

func nodifyParam(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
		return err
	}
	if typ, ok := nodes[2].(Type); ok {
		return MemberType{
			Name:  nodes[0].(*parsec.Terminal).GetValue(),
			Value: typ,
		}
	} else {
		return fmt.Errorf("failed to parse param: %s", nodes[2])
	}
}

func nodifyParams(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
		return err
	}
	members := make([]MemberType, len(nodes))
	for i, node := range nodes {
		if err, ok := checkError([]Node{node}); ok {
			return fmt.Errorf("failed to parse parameter %d: %s", i, err)
		} else if member, ok := node.(MemberType); ok {
			members[i] = member
		} else {
			return fmt.Errorf("failed to parse parameter %d: %s", i, node)
		}
	}
	return NewStructType("parameters", members)
}

func nodifyAndParams(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
		return err
	}
	if _, ok := nodes[0].(parsec.MaybeNone); ok {
		return NewStructType("parameters", make([]MemberType, 0))
	} else if ret, ok := nodes[0].(*StructType); ok {
		return ret
	} else {
		return fmt.Errorf("unexpected non struct type(%s): %v",
			reflect.TypeOf(nodes[0]), nodes[0])
	}
}
