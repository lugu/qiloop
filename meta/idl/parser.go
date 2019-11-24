package idl

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"

	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	parsec "github.com/prataprc/goparsec"
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
		parsec.Atom("any", ""),
		parsec.Atom("unknown", ""))
}

// Context catures the current state of the parser.
type Context struct {
	scope      Scope
	typeParser parsec.Parser
}

// NewContext creates a new context.
func NewContext() *Context {
	var ctx Context
	parser := func(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
		return ctx.typeParser(s)
	}
	ctx.scope = NewScope()
	ctx.typeParser = parser
	ctx.typeParser = typeParser(&ctx)
	return &ctx
}

func mapType(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifyMap,
		parsec.Atom("Map<", "Map<"),
		ctx.typeParser,
		parsec.Atom(",", ","),
		ctx.typeParser,
		parsec.Atom(">", ">"))
}

func tupleType(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifyTuple,
		parsec.Atom("Tuple<", "Tuple<"),
		parsec.Many(
			nodifyList,
			ctx.typeParser,
			parsec.Atom(",", ","),
		),
		parsec.Atom(">", ">"))
}

func vecType(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifyVec,
		parsec.Atom("Vec<", "Vec<"),
		ctx.typeParser,
		parsec.Atom(">", ">"))
}

func typeParser(ctx *Context) parsec.Parser {
	return parsec.OrdChoice(
		nodifyType,
		basicType(),
		mapType(ctx),
		tupleType(ctx),
		vecType(ctx),
		referenceType(ctx),
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

func returns(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifyReturns,
		parsec.Maybe(
			nodifyMaybeReturns,
			parsec.And(
				nodifyReturnsType,
				parsec.Atom("->", "->"),
				ctx.typeParser,
			),
		),
	)
}

func parameter(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifyParam,
		ident(),
		parsec.Atom(":", ":"),
		ctx.typeParser,
	)
}

func parameters(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifyAndParams,
		parsec.Maybe(
			nodifyMaybeParams,
			parsec.Many(
				nodifyParams,
				parameter(ctx),
				parsec.Atom(",", ","),
			),
		),
	)
}

func ident() parsec.Parser {
	return parsec.Token(`[_A-Za-z][0-9a-zA-Z_]*`, "IDENT")
}

func typeIdent() parsec.Parser {
	return parsec.OrdTokens(
		[]string{
			`[_A-Za-z][0-9a-zA-Z_]*<[0-9a-zA-Z_]*>`,
			`[_A-Za-z][0-9a-zA-Z_]*`,
		},
		[]string{
			"TYPE_IDENT_CPP",
			"TYPE_IDENT",
		},
	)
}

func method(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifyMethod,
		parsec.Atom("fn", "fn"),
		ident(),
		parsec.Atom("(", "("),
		parameters(ctx),
		parsec.Atom(")", ")"),
		returns(ctx),
		comments(),
	)
}

func signal(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifySignal,
		parsec.Atom("sig", "sig"),
		ident(),
		parsec.Atom("(", "("),
		parameters(ctx),
		parsec.Atom(")", ")"),
		comments(),
	)
}

func property(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifyProperty,
		parsec.Atom("prop", "prop"),
		ident(),
		parsec.Atom("(", "("),
		parameters(ctx),
		parsec.Atom(")", ")"),
		comments(),
	)
}

func action(ctx *Context) parsec.Parser {
	return parsec.OrdChoice(
		nodifyAction,
		method(ctx),
		signal(ctx),
		property(ctx),
	)
}

func interfaceParser(ctx *Context) parsec.Parser {
	return parsec.And(
		makeNodifyInterface(ctx.scope),
		parsec.Atom("interface", "interface"),
		ident(),
		comments(),
		parsec.Kleene(nodifyActionList, action(ctx)),
		parsec.Atom("end", "end"),
		comments(),
	)
}

func referenceType(ctx *Context) parsec.Parser {
	return parsec.And(
		makeNodifyTypeReference(ctx.scope),
		typeIdent(),
	)
}

func member(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifyMember,
		ident(),
		parsec.Atom(":", ":"),
		ctx.typeParser,
		comments(),
	)
}

func constValue() parsec.Parser {
	return parsec.Int()
}

func enumConst() parsec.Parser {
	return parsec.And(
		nodifyEnumConst,
		ident(),
		parsec.Atom("=", "="),
		constValue(),
		comments(),
	)
}

func enum() parsec.Parser {
	return parsec.And(
		nodifyEnum,
		parsec.Atom("enum", "enum"),
		ident(),
		comments(),
		parsec.Kleene(nodifyEnumMembers, enumConst()),
		parsec.Atom("end", "end"),
		comments(),
	)
}

func structure(ctx *Context) parsec.Parser {
	return parsec.And(
		makeNodifyStructure(ctx.scope),
		parsec.Atom("struct", "struct"),
		typeIdent(),
		comments(),
		parsec.Kleene(nodifyMemberList, member(ctx)),
		parsec.Atom("end", "end"),
		comments(),
	)
}

func declaration(ctx *Context) parsec.Parser {
	return parsec.OrdChoice(
		nodifyDeclaration,
		structure(ctx),
		enum(),
		interfaceParser(ctx),
	)
}

func declarationsList(ctx *Context) parsec.Parser {
	return parsec.Kleene(
		nodifyDeclarationList,
		declaration(ctx),
	)
}

func packageName() parsec.Parser {
	return parsec.And(nodifyPackageName,
		parsec.Maybe(nodifyPackageNameMaybe,
			parsec.And(nodifyPackageNameAnd,
				parsec.Atom("package", "package"),
				parsec.Token(`[_A-Za-z][0-9a-zA-Z-._]*`, "IDENT"),
				comments(),
			),
		),
	)
}

func packageParser(ctx *Context) parsec.Parser {
	return parsec.And(nodifyPackage,
		packageName(),
		// imports(),
		declarationsList(ctx),
	)
}

// PackageDeclaration is the result of the parsing of an IDL file. It
// represents a package declaration containing a package name and a
// list of declared types.
type PackageDeclaration struct {
	Name  string
	Types []signature.Type
}

// ParsePackage takes an IDL file as input, parse it and returns a
// PackageDeclaration.
func ParsePackage(input []byte) (*PackageDeclaration, error) {
	context := NewContext()
	parser := packageParser(context)
	root, scanner := parser(parsec.NewScanner(input).TrackLineno())
	_, scanner = scanner.SkipWS()
	if !scanner.Endof() {
		return nil, fmt.Errorf("parsing error at line: %d", scanner.Lineno())
	}
	if root == nil {
		return nil, fmt.Errorf("cannot parse input:\n%s", input)
	}
	definitions, ok := root.(*PackageDeclaration)
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
	pkg, err := ParsePackage(input)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %s", err)
	}

	metas := make([]object.MetaObject, 0)
	for _, decl := range pkg.Types {
		itf, ok := decl.(*InterfaceType)
		if ok {
			meta := itf.MetaObject()
			metas = append(metas, meta)
		}
	}
	return metas, nil
}

func nodifyDeclaration(nodes []signature.Node) signature.Node {
	return nodes[0]
}

func nodifyMaybeComment(nodes []signature.Node) signature.Node {
	return nodes[0]
}

func nodifyReturnsType(nodes []signature.Node) signature.Node {
	return nodes[1]
}

func nodifyAction(nodes []signature.Node) signature.Node {
	return nodes[0]
}

func nodifyType(nodes []signature.Node) signature.Node {
	return nodes[0]
}

func nodifyMaybeReturns(nodes []signature.Node) signature.Node {
	return nodes[0]
}

func nodifyMaybeParams(nodes []signature.Node) signature.Node {
	return nodes[0]
}

func nodifyComment(nodes []signature.Node) signature.Node {
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

func nodifyCommentContent(nodes []signature.Node) signature.Node {
	comment := nodes[1].(*parsec.Terminal).GetValue()
	var uid uint32
	_, err := fmt.Sscanf(comment, "uid:%d", &uid)
	if err == nil {
		return uid
	}
	return comment
}

// returns a signature.Type or an error
func nodifyReturns(nodes []signature.Node) signature.Node {
	if _, ok := nodes[0].(parsec.MaybeNone); ok {
		return signature.NewVoidType()
	} else if ret, ok := nodes[0].(signature.Type); ok {
		return ret
	} else {
		return fmt.Errorf("unexpected return value (%s): %v",
			reflect.TypeOf(nodes[0]), nodes[0])
	}
}

// nodifyDeclarationList returns a list of signature.Type.
func nodifyDeclarationList(nodes []signature.Node) signature.Node {
	if ok, err := checkError(nodes); ok {
		return err
	}
	typeList := make([]signature.Type, 0)
	for i, node := range nodes {
		if ok, err := checkError([]signature.Node{node}); ok {
			return fmt.Errorf("parse parameter %d: %s", i, err)
		} else if typ, ok := node.(signature.Type); ok {
			typeList = append(typeList, typ)
		} else {
			return fmt.Errorf("parse type %d: got %+v: %+v",
				i, reflect.TypeOf(node), node)
		}
	}
	return typeList
}

func nodifyPackageNameAnd(nodes []signature.Node) signature.Node {
	return nodes[1].(*parsec.Terminal).GetValue()
}

func nodifyPackageNameMaybe(nodes []signature.Node) signature.Node {
	return nodes[0]
}

func nodifyPackageName(nodes []signature.Node) signature.Node {
	if _, ok := nodes[0].(parsec.MaybeNone); ok {
		return ""
	}
	if name, ok := nodes[0].(string); ok {
		return name
	}
	return fmt.Errorf("Expecting package name, got %+v: %+v",
		reflect.TypeOf(nodes[0]), nodes[0])
}

// nodifyPackage returns a package structure.
func nodifyPackage(nodes []signature.Node) signature.Node {
	packageNode := nodes[0]
	definitions := nodes[1]
	packageName := packageNode.(string)
	typeList, ok := definitions.([]signature.Type)
	if !ok {
		return fmt.Errorf("Expecting type list, got %+v: %+v",
			reflect.TypeOf(definitions), definitions)
	}
	return &PackageDeclaration{
		Name:  packageName,
		Types: typeList,
	}
}

func makeNodifyTypeReference(sc Scope) func([]signature.Node) signature.Node {
	return func(nodes []signature.Node) signature.Node {
		typeNode := nodes[0]
		typeName := typeNode.(*parsec.Terminal).GetValue()
		return NewRefType(typeName, sc)
	}
}

func nodifyEnum(nodes []signature.Node) signature.Node {
	nameNode := nodes[1]
	enumNode := nodes[3]
	if enum, ok := enumNode.(*signature.EnumType); ok {
		enum.Name = nameNode.(*parsec.Terminal).GetValue()
		return enum
	} else if err, ok := enumNode.(error); ok {
		return err
	} else {
		return fmt.Errorf("unexpected enum value: %v", enumNode)
	}
}

func nodifyEnumMembers(nodes []signature.Node) signature.Node {
	var enum signature.EnumType
	enum.Values = make(map[string]int)
	for _, node := range nodes {
		if member, ok := node.(signature.EnumMember); ok {
			enum.Values[member.Const] = member.Value
		} else if err, ok := node.(error); ok {
			return err
		} else {
			return fmt.Errorf("unexpected member node: %v", node)
		}
	}
	return &enum
}

func nodifyEnumConst(nodes []signature.Node) signature.Node {
	nameNode := nodes[0]
	constNode := nodes[2]
	valueStr := constNode.(*parsec.Terminal).GetValue()
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fmt.Errorf("Integer expected in enum definition: %s", valueStr)
	}
	return signature.EnumMember{
		Const: nameNode.(*parsec.Terminal).GetValue(),
		Value: value,
	}
}

// nodifyMember returns a MemberType or an error.
func nodifyMember(nodes []signature.Node) signature.Node {
	nameNode := nodes[0]
	typeNode := nodes[2]
	var member signature.MemberType
	var ok bool
	member.Name = nameNode.(*parsec.Terminal).GetValue()
	member.Type, ok = typeNode.(signature.Type)
	if !ok {
		return fmt.Errorf("Expecting Type, got %+v: %+v", reflect.TypeOf(typeNode), typeNode)
	}
	return member
}

// nodifyMemberList returns a StructType or an error.
func nodifyMemberList(nodes []signature.Node) signature.Node {
	members := make([]signature.MemberType, len(nodes))
	for i, node := range nodes {
		if member, ok := node.(signature.MemberType); ok {
			members[i] = member
		} else {
			return fmt.Errorf("MemberList: unexpected type, got %+v: %+v", reflect.TypeOf(node), node)
		}
	}
	return signature.NewStructType("parameters", members)
}

// makeNodifyStructure returns a *StructType of an error.
// the structure is added to the scope
func makeNodifyStructure(sc Scope) func([]signature.Node) signature.Node {
	return func(nodes []signature.Node) signature.Node {
		nameNode := nodes[1]
		structNode := nodes[3]
		structType, ok := structNode.(*signature.StructType)
		if !ok {
			return fmt.Errorf("unexpected non struct type(%s): %v",
				reflect.TypeOf(structNode), structNode)
		}
		structType.Name = nameNode.(*parsec.Terminal).GetValue()
		sc.Add(structType.Name, structType)
		return structType
	}
}

// makeNodifyInterface returns an *InterfaceType or an error
// the InterfaceType is added to the scope
func makeNodifyInterface(sc Scope) func([]signature.Node) signature.Node {
	return func(nodes []signature.Node) signature.Node {
		nameNode := nodes[1]
		itfNode := nodes[3]
		itf, ok := itfNode.(*InterfaceType)
		if !ok {
			return fmt.Errorf("Expecting InterfaceType, got %+v: %+v",
				reflect.TypeOf(itfNode), itfNode)
		}
		itf.Name = nameNode.(*parsec.Terminal).GetValue()
		itf.Scope = sc
		sc.Add(itf.Name, itf)
		return itf
	}
}

func nodifyVec(nodes []signature.Node) signature.Node {

	elementNode := nodes[1]
	elementType, ok := elementNode.(signature.Type)
	if !ok {
		return fmt.Errorf("invalid vector element: %v", elementNode)
	}
	return signature.NewListType(elementType)
}

func nodifyTuple(nodes []signature.Node) signature.Node {

	if params, ok := nodes[1].([]Parameter); ok {
		members := []signature.MemberType{}
		for _, p := range params {
			members = append(members, signature.MemberType{
				Name: p.Name,
				Type: p.Type,
			})
		}
		return &signature.TupleType{
			Members: members,
		}
	} else {
		return fmt.Errorf("unexpected non param type(%s): %v",
			reflect.TypeOf(nodes[0]), nodes[0])
	}
}

func nodifyMap(nodes []signature.Node) signature.Node {

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

func nodifyBasicType(nodes []signature.Node) signature.Node {
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
		return signature.NewObjectType()
	case "unknown":
		return signature.NewUnknownType()
	default:
		return fmt.Errorf("unknown type: %s", sig)
	}
}

// nodifyMethod returns a Method or an error.
func nodifyMethod(nodes []signature.Node) signature.Node {
	if ok, err := checkError(nodes); ok {
		return fmt.Errorf("parse method: %s", err)
	}
	retNode := nodes[5]
	paramNode := nodes[3]
	commentNode := nodes[6]
	params, ok := paramNode.([]Parameter)
	if !ok {
		return fmt.Errorf("method: failed to convert param type (%s): %v",
			reflect.TypeOf(paramNode), paramNode)
	}
	retType, ok := retNode.(signature.Type)
	if !ok {
		return fmt.Errorf("method: failed to convert return type (%s): %v",
			reflect.TypeOf(retNode), retNode)
	}

	var method Method
	method.Name = nodes[1].(*parsec.Terminal).GetValue()
	method.Params = params
	method.Return = retType

	if uid, ok := commentNode.(uint32); ok {
		method.ID = uid
	}
	return method
}

// nodifySignal returns either a Signal or an error.
func nodifySignal(nodes []signature.Node) signature.Node {
	if ok, err := checkError(nodes); ok {
		return fmt.Errorf("parse method: %s", err)
	}
	paramNode := nodes[3]
	commentNode := nodes[5]
	params, ok := paramNode.([]Parameter)
	if !ok {
		return fmt.Errorf("signal: failed to convert param type (%s): %v",
			reflect.TypeOf(paramNode), paramNode)
	}
	var signal Signal
	signal.Name = nodes[1].(*parsec.Terminal).GetValue()
	signal.Params = params
	if uid, ok := commentNode.(uint32); ok {
		signal.ID = uid
	}
	return signal
}

// nodifyProperty returns either a Property or an error.
func nodifyProperty(nodes []signature.Node) signature.Node {
	if ok, err := checkError(nodes); ok {
		return fmt.Errorf("parse method: %s", err)
	}
	paramNode := nodes[3]
	commentNode := nodes[5]
	params, ok := paramNode.([]Parameter)
	if !ok {
		return fmt.Errorf("signal: failed to convert param type (%s): %v",
			reflect.TypeOf(paramNode), paramNode)
	}
	var prop Property
	prop.Name = nodes[1].(*parsec.Terminal).GetValue()
	prop.Params = params
	if uid, ok := commentNode.(uint32); ok {
		prop.ID = uid
	}
	return prop
}

// nodifyActionList returns an *InterfaceType or an error
// the InterfaceType has no name or namespce.
func nodifyActionList(nodes []signature.Node) signature.Node {
	var itf InterfaceType

	itf.Methods = make(map[uint32]Method)
	itf.Signals = make(map[uint32]Signal)
	itf.Properties = make(map[uint32]Property)

	var customAction = uint32(100)
	for _, node := range nodes {
		if err, ok := node.(error); ok {
			return err
		}
		if method, ok := node.(Method); ok {
			if method.ID == 0 && method.Name != "registerEvent" {
				method.ID = customAction
				customAction++
			}
			itf.Methods[method.ID] = method
		} else if signal, ok := node.(Signal); ok {
			if signal.ID == 0 {
				signal.ID = customAction
				customAction++
			}
			itf.Signals[signal.ID] = signal
		} else if property, ok := node.(Property); ok {
			if property.ID == 0 {
				property.ID = customAction
				customAction++
			}
			itf.Properties[property.ID] = property
		} else {
			return fmt.Errorf("Expecting action, got %+v: %+v",
				reflect.TypeOf(node), node)
		}
	}
	return &itf
}

// checkError tests if nodes is an error, if so it returns true with
// the error, else false and nil.
func checkError(nodes []signature.Node) (bool, error) {
	if len(nodes) == 1 {
		if err, ok := nodes[0].(error); ok {
			return true, err
		}
	}
	return false, nil
}

func nodifyList(nodes []signature.Node) signature.Node {
	params := make([]Parameter, len(nodes))
	for i, node := range nodes {
		if ok, err := checkError([]signature.Node{node}); ok {
			return fmt.Errorf("tuple list %d: %s", i, err)
		} else if typ, ok := node.(signature.Type); ok {
			params[i] = Parameter{
				Name: fmt.Sprintf("param%d", i),
				Type: typ,
			}
		} else {
			return fmt.Errorf("parse parameter %d: %s", i, node)
		}
	}
	return params
}

func nodifyParam(nodes []signature.Node) signature.Node {
	if ok, err := checkError(nodes); ok {
		return err
	}
	if typ, ok := nodes[2].(signature.Type); ok {
		return Parameter{
			Name: nodes[0].(*parsec.Terminal).GetValue(),
			Type: typ,
		}
	}
	return fmt.Errorf("parse param: %s", nodes[2])
}

// nodifyParams returns an slice of Parameter
func nodifyParams(nodes []signature.Node) signature.Node {
	if ok, err := checkError(nodes); ok {
		return err
	}
	params := make([]Parameter, len(nodes))
	for i, node := range nodes {
		if ok, err := checkError([]signature.Node{node}); ok {
			return fmt.Errorf("parse parameter %d: %s", i, err)
		} else if member, ok := node.(Parameter); ok {
			params[i] = member
		} else {
			return fmt.Errorf("parse parameter %d: %s", i, node)
		}
	}
	return params
}

// return a slice of Parameter
func nodifyAndParams(nodes []signature.Node) signature.Node {
	if ok, err := checkError(nodes); ok {
		return err
	}
	if _, ok := nodes[0].(parsec.MaybeNone); ok {
		return make([]Parameter, 0)
	} else if ret, ok := nodes[0].([]Parameter); ok {
		return ret
	} else {
		return fmt.Errorf("unexpected non struct type(%s): %v",
			reflect.TypeOf(nodes[0]), nodes[0])
	}
}
