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
	return new(object.MetaObject)
}

// ParseIDL read an IDL definition from a reader and returns the
// MetaObject associated with the IDL.
func Parse(reader io.Reader) ([]object.MetaObject, error) {
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read input: %s", err)
	}

	interfaceSignature := parsec.And(
		nodifyInterface,
		parsec.Atom("interface", "interface"),
		parsec.Ident(),
		parsec.Atom("end", "end"),
	)

	interfaceList := parsec.Many(nodifyInterfaceList, &interfaceSignature)

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
