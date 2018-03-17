package idl

import (
	"fmt"
	"github.com/lugu/qiloop/type/object"
	parsec "github.com/prataprc/goparsec"
	"io"
	"io/ioutil"
	"reflect"
)

// Node is an alias to parsec.ParsecNode
type Node = parsec.ParsecNode

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

func methodParser() parsec.Parser {
	return parsec.And(
		nil,
		parsec.Atom("fn", "fn"),
		parsec.Ident(),
	)
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

func interfaceParser() parsec.Parser {
	return parsec.And(
		nodifyInterface,
		parsec.Atom("interface", "interface"),
		parsec.Ident(),
		parsec.Kleene(nodifyMethodList, methodParser()),
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
