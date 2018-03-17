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

// ParseIDL read an IDL definition from a reader and returns the
// MetaObject associated with the IDL.
func Parse(reader io.Reader) (*object.MetaObject, error) {
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read input: %s", err)
	}

	interfaceSignature := parsec.Atom("not implmemented", "TypeParameterClose")

	root, _ := interfaceSignature(parsec.NewScanner(input))
	if root == nil {
		return nil, fmt.Errorf("failed to parse signature: %s", input)
	}
	types, ok := root.([]Node)
	if !ok {
		return nil, fmt.Errorf("failed to convert array: %+v", reflect.TypeOf(root))
	}
	if len(types) != 1 {
		return nil, fmt.Errorf("did not parse only one IDL: %+v", root)
	}
	metaObj, ok := types[0].(object.MetaObject)
	if !ok {
		return nil, fmt.Errorf("failed to convert value: %+v", reflect.TypeOf(types[0]))
	}
	return &metaObj, nil
}
