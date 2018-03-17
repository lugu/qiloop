package idl

import (
	"fmt"
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	parsec "github.com/prataprc/goparsec"
	"io"
	"io/ioutil"
	"reflect"
)

// Node is an alias to parsec.ParsecNode
type Node = parsec.ParsecNode

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
		parsec.Atom("any", ""))
}

// required to break the recursive definition of typeParser()
var mapParser parsec.Parser
var vecParser parsec.Parser

func init() {
	mapParser = mapType()
	vecParser = vecType()
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
		parsec.Ident(),
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

func method() parsec.Parser {
	return parsec.And(
		nodifyMethod,
		parsec.Atom("fn", "fn"),
		parsec.Ident(),
		parsec.Atom("(", "("),
		parameters(),
		parsec.Atom(")", ")"),
		returns(),
	)
}

func interfaceParser() parsec.Parser {
	return parsec.And(
		nodifyInterface,
		parsec.Atom("interface", "interface"),
		parsec.Ident(),
		parsec.Kleene(nodifyMethodList, method()),
		parsec.Atom("end", "end"),
	)
}

// ParseIDL read an IDL definition from a reader and returns the
// MetaObject associated with the IDL.
func Parse(reader io.Reader) ([]object.MetaObject, error) {
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read input: %s", err)
	}

	interfaceList := parsec.Many(nodifyInterfaceList, interfaceParser())

	root, _ := interfaceList(parsec.NewScanner(input))
	if root == nil {
		return nil, fmt.Errorf("cannot parse input:\n%s", input)
	}

	metas, ok := root.([]object.MetaObject)
	if !ok {
		err, ok := root.(error)
		if !ok {
			return nil, fmt.Errorf("cannot parse IDL: %+v", reflect.TypeOf(root))
		}
		return nil, err
	}
	return metas, nil
}

func nodifyReturnsType(nodes []Node) Node {
	return nodes[1]
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

func nodifyReturns(nodes []Node) Node {
	if _, ok := nodes[0].(parsec.MaybeNone); ok {
		return signature.NewVoidType()
	} else if ret, ok := nodes[0].(signature.Type); ok {
		return ret
	} else {
		return fmt.Errorf("unexpected return value (%s): %v",
			reflect.TypeOf(nodes[0]), nodes[0])
	}
}

func nodifyInterfaceList(nodes []Node) Node {
	interfaces := make([]object.MetaObject, 0, len(nodes))
	for _, node := range nodes {
		if err, ok := node.(error); ok {
			return err
		}
		if metaObj, ok := node.(*object.MetaObject); ok {
			interfaces = append(interfaces, *metaObj)
		} else {
			return fmt.Errorf("Expecting MetaObject, got %+v: %#v", reflect.TypeOf(node), node)
		}
	}
	return interfaces
}

func nodifyInterface(nodes []Node) Node {
	metaObj := new(object.MetaObject)
	metaObj.Description = nodes[1].(*parsec.Terminal).GetValue()
	return metaObj
}

func nodifyVec(nodes []Node) Node {

	elementNode := nodes[1]
	elementType, ok := elementNode.(signature.Type)
	if !ok {
		return fmt.Errorf("invalid vector element: %v", elementNode)
	}
	return signature.NewListType(elementType)
}

func nodifyMap(nodes []Node) Node {

	keyNode := nodes[1]
	valueNode := nodes[3]
	keyType, ok := keyNode.(signature.Type)
	if !ok {
		return fmt.Errorf("invalid map key: %v", keyNode)
	}
	valueType, ok := valueNode.(signature.Type)
	if !ok {
		return fmt.Errorf("invalid map value: %v", valueNode)
	}
	return signature.NewMapType(keyType, valueType)
}

func nodifyBasicType(nodes []Node) Node {
	if len(nodes) != 1 {
		return fmt.Errorf("basic type array size: %d: %s", len(nodes), nodes)
	}
	sig := nodes[0].(*parsec.Terminal).GetValue()
	switch sig {
	case "int32":
		return signature.NewIntType()
	case "uint32":
		return signature.NewUIntType()
	case "int64":
		return signature.NewLongType()
	case "uint64":
		return signature.NewULongType()
	case "float32":
		return signature.NewFloatType()
	case "float64":
		return signature.NewDoubleType()
	case "str":
		return signature.NewStringType()
	case "bool":
		return signature.NewBoolType()
	case "any":
		return signature.NewValueType()
	case "obj":
		return signature.NewValueType()
	default:
		return fmt.Errorf("unknown type: %s", sig)
	}
}

func nodifyMethod(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
		return fmt.Errorf("failed to parse method: %s", err)
	}
	retNode := nodes[5]
	paramNode := nodes[3]
	structType, ok := paramNode.(*signature.StructType)
	if !ok {
		return fmt.Errorf("failed to convert param type (%s): %v", reflect.TypeOf(paramNode), paramNode)
	}
	retType, ok := retNode.(signature.Type)
	if !ok {
		return fmt.Errorf("failed to convert return type (%s): %v", reflect.TypeOf(retNode), retNode)
	}

	params := make([]signature.Type, len(structType.Members))
	for i, member := range structType.Members {
		params[i] = member.Value
	}
	tupleType := signature.NewTupleType(params)

	method := new(object.MetaMethod)
	method.Name = nodes[1].(*parsec.Terminal).GetValue()
	method.ParametersSignature = tupleType.Signature()
	if len(structType.Members) != 0 {
		method.Parameters = make([]object.MetaMethodParameter, len(structType.Members))
		for i, member := range structType.Members {
			method.Parameters[i].Name = member.Name
		}
	}
	method.ReturnSignature = retType.Signature()
	return method
}

func nodifyMethodList(nodes []Node) Node {
	methods := make([]object.MetaMethod, 0, len(nodes))
	for _, node := range nodes {
		if err, ok := node.(error); ok {
			return err
		}
		if method, ok := node.(object.MetaMethod); ok {
			methods = append(methods, method)
		} else {
			return fmt.Errorf("Expecting MetaMethod, got %+v: %#v", reflect.TypeOf(node), node)
		}
	}
	return methods
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
	if typ, ok := nodes[2].(signature.Type); ok {
		return signature.MemberType{
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
	members := make([]signature.MemberType, len(nodes))
	for i, node := range nodes {
		if err, ok := checkError([]Node{node}); ok {
			return fmt.Errorf("failed to parse parameter %d: %s", i, err)
		} else if member, ok := node.(signature.MemberType); ok {
			members[i] = member
		} else {
			return fmt.Errorf("failed to parse parameter %d: %s", i, node)
		}
	}
	return signature.NewStrucType("parameters", members)
}

func nodifyAndParams(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
		return err
	}
	if _, ok := nodes[0].(parsec.MaybeNone); ok {
		return signature.NewStrucType("parameters", make([]signature.MemberType, 0))
	} else if ret, ok := nodes[0].(*signature.StructType); ok {
		return ret
	} else {
		return fmt.Errorf("unexpected non struct type(%s): %v",
			reflect.TypeOf(nodes[0]), nodes[0])
	}
}
