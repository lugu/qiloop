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
		parsec.Atom("any", ""))
}

func typeParser() parsec.Parser {
	return basicType()
}

func returns() parsec.Parser {
	return parsec.Maybe(
		nodifyMaybeReturns,
		parsec.And(
			nodifyReturns,
			parsec.Atom("->", "->"),
			typeParser(),
		),
	)
}

func parameters() parsec.Parser {
	return parsec.Kleene(
		nil,
		parsec.And(
			nil,
			parsec.Token(`[0-9a-zA-Z_]*:`, "return type"),
			typeParser(),
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

func nodifyReturns(nodes []Node) Node {
	return nodes[1]
}

func nodifyMaybeReturns(nodes []Node) Node {
	if _, ok := nodes[0].(parsec.MaybeNone); ok {
		return signature.NewVoidType()
	} else if ret, ok := nodes[0].(signature.Type); ok {
		return ret
	} else {
		return fmt.Errorf("unexpected return value: %v", nodes[0])
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
	method := new(object.MetaMethod)
	method.Name = nodes[1].(*parsec.Terminal).GetValue()
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
