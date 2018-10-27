package idl

import (
	"fmt"
	. "github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	parsec "github.com/prataprc/goparsec"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
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

type Context struct {
	scope      Scope
	typeParser parsec.Parser
}

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
		parsec.Atom(">", ">"),
	)
}

func vecType(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifyVec,
		parsec.Atom("Vec<", "Vec<"),
		ctx.typeParser,
		parsec.Atom(">", ">"),
	)
}

func typeParser(ctx *Context) parsec.Parser {
	return parsec.OrdChoice(
		nodifyType,
		basicType(),
		mapType(ctx),
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
		Ident(),
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

func Ident() parsec.Parser {
	return parsec.Token(`[_A-Za-z][0-9a-zA-Z_]*`, "IDENT")
}

func method(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifyMethod,
		parsec.Atom("fn", "fn"),
		Ident(),
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
		Ident(),
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
		Ident(),
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
		Ident(),
		comments(),
		parsec.Kleene(nodifyActionList, action(ctx)),
		parsec.Atom("end", "end"),
		comments(),
	)
}

func referenceType(ctx *Context) parsec.Parser {
	return parsec.And(
		makeNodifyTypeReference(ctx.scope),
		Ident(),
	)
}

func member(ctx *Context) parsec.Parser {
	return parsec.And(
		nodifyMember,
		Ident(),
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
		Ident(),
		parsec.Atom("=", "="),
		constValue(),
		comments(),
	)
}

func enum() parsec.Parser {
	return parsec.And(
		nodifyEnum,
		parsec.Atom("enum", "enum"),
		Ident(),
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
		Ident(),
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

type PackageDeclaration struct {
	Name  string
	Types []Type
}

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

// returns a signature.Type or an error
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

// nodifyDeclarationList returns a list of signature.Type.
func nodifyDeclarationList(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
		return err
	}
	typeList := make([]Type, 0)
	for i, node := range nodes {
		if err, ok := checkError([]Node{node}); ok {
			return fmt.Errorf("failed to parse parameter %d: %s", i, err)
		} else if typ, ok := node.(Type); ok {
			typeList = append(typeList, typ)
		} else {
			return fmt.Errorf("failed to parse type %d: got %+v: %+v",
				i, reflect.TypeOf(node), node)
		}
	}
	return typeList
}

func nodifyPackageNameAnd(nodes []Node) Node {
	return nodes[1].(*parsec.Terminal).GetValue()
}

func nodifyPackageNameMaybe(nodes []Node) Node {
	return nodes[0]
}

func nodifyPackageName(nodes []Node) Node {
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
func nodifyPackage(nodes []Node) Node {
	var packageName = ""
	packageNode := nodes[0]
	definitions := nodes[1]
	packageName = packageNode.(string)
	if typeList, ok := definitions.([]Type); !ok {
		return fmt.Errorf("Expecting type list, got %+v: %+v",
			reflect.TypeOf(definitions), definitions)
	} else {
		return &PackageDeclaration{
			Name:  packageName,
			Types: typeList,
		}
	}
}

func makeNodifyTypeReference(sc Scope) func([]Node) Node {
	return func(nodes []Node) Node {
		typeNode := nodes[0]
		typeName := typeNode.(*parsec.Terminal).GetValue()
		return NewRefType(typeName, sc)
	}
}

func nodifyEnum(nodes []Node) Node {
	nameNode := nodes[1]
	enumNode := nodes[3]
	if enum, ok := enumNode.(*EnumType); ok {
		enum.Name = nameNode.(*parsec.Terminal).GetValue()
		return enum
	} else if err, ok := enumNode.(error); ok {
		return err
	} else {
		return fmt.Errorf("unexpected enum value: %v", enumNode)
	}
}

func nodifyEnumMembers(nodes []Node) Node {
	var enum EnumType
	enum.Values = make(map[string]int)
	for _, node := range nodes {
		if member, ok := node.(EnumMember); ok {
			enum.Values[member.Const] = member.Value
		} else if err, ok := node.(error); ok {
			return err
		} else {
			return fmt.Errorf("unexpected member node: %v", node)
		}
	}
	return &enum
}

func nodifyEnumConst(nodes []Node) Node {
	nameNode := nodes[0]
	constNode := nodes[2]
	valueStr := constNode.(*parsec.Terminal).GetValue()
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fmt.Errorf("Integer expected in enum definition: %s", valueStr)
	}
	return EnumMember{
		Const: nameNode.(*parsec.Terminal).GetValue(),
		Value: value,
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

// makeNodifyStructure returns a *StructType of an error.
// the structure is added to the scope
func makeNodifyStructure(sc Scope) func([]Node) Node {
	return func(nodes []Node) Node {
		nameNode := nodes[1]
		structNode := nodes[3]
		structType, ok := structNode.(*StructType)
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
func makeNodifyInterface(sc Scope) func([]Node) Node {
	return func(nodes []Node) Node {
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

// nodifyMethod returns a Method or an error.
func nodifyMethod(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
		return fmt.Errorf("failed to parse method: %s", err)
	}
	retNode := nodes[5]
	paramNode := nodes[3]
	commentNode := nodes[6]
	params, ok := paramNode.([]Parameter)
	if !ok {
		return fmt.Errorf("method: failed to convert param type (%s): %v",
			reflect.TypeOf(paramNode), paramNode)
	}
	retType, ok := retNode.(Type)
	if !ok {
		return fmt.Errorf("method: failed to convert return type (%s): %v",
			reflect.TypeOf(retNode), retNode)
	}

	var method Method
	method.Name = nodes[1].(*parsec.Terminal).GetValue()
	method.Params = params
	method.Return = retType

	if uid, ok := commentNode.(uint32); ok {
		method.Id = uid
	}
	return method
}

// nodifySignal returns either a Signal or an error.
func nodifySignal(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
		return fmt.Errorf("failed to parse method: %s", err)
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
		signal.Id = uid
	}
	return signal
}

// nodifyProperty returns either a Property or an error.
func nodifyProperty(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
		return fmt.Errorf("failed to parse method: %s", err)
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
		prop.Id = uid
	}
	return prop
}

// nodifyActionList returns an *InterfaceType or an error
// the InterfaceType has no name or namespce.
func nodifyActionList(nodes []Node) Node {
	var itf InterfaceType

	itf.Methods = make(map[uint32]Method)
	itf.Signals = make(map[uint32]Signal)
	itf.Properties = make(map[uint32]Property)

	var customAction uint32 = 100
	for _, node := range nodes {
		if err, ok := node.(error); ok {
			return err
		}
		if method, ok := node.(Method); ok {
			if method.Id == 0 && method.Name != "registerEvent" {
				method.Id = customAction
				customAction++
			}
			itf.Methods[method.Id] = method
		} else if signal, ok := node.(Signal); ok {
			if signal.Id == 0 {
				signal.Id = customAction
				customAction++
			}
			itf.Signals[signal.Id] = signal
		} else if property, ok := node.(Property); ok {
			if property.Id == 0 {
				property.Id = customAction
				customAction++
			}
			itf.Properties[property.Id] = property
		} else {
			return fmt.Errorf("Expecting action, got %+v: %+v",
				reflect.TypeOf(node), node)
		}
	}
	return &itf
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
		return Parameter{
			Name: nodes[0].(*parsec.Terminal).GetValue(),
			Type: typ,
		}
	} else {
		return fmt.Errorf("failed to parse param: %s", nodes[2])
	}
}

// nodifyParams returns an slice of Parameter
func nodifyParams(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
		return err
	}
	params := make([]Parameter, len(nodes))
	for i, node := range nodes {
		if err, ok := checkError([]Node{node}); ok {
			return fmt.Errorf("failed to parse parameter %d: %s", i, err)
		} else if member, ok := node.(Parameter); ok {
			params[i] = member
		} else {
			return fmt.Errorf("failed to parse parameter %d: %s", i, node)
		}
	}
	return params
}

// return a slice of Parameter
func nodifyAndParams(nodes []Node) Node {
	if err, ok := checkError(nodes); ok {
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
