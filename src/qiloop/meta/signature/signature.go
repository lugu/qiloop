package signature

import (
	"bytes"
	"fmt"
	"github.com/dave/jennifer/jen"
	"io"
	"log"
	"reflect"
	"strings"
	parsec "github.com/prataprc/goparsec"
)

type TypeSet struct {
    signatures map[string]bool
    types []ValueConstructor
}

func (s *TypeSet) Register(v ValueConstructor) {
    if _, ok := s.signatures[v.Signature()]; !ok {
        s.signatures[v.Signature()] = true
        s.types = append(s.types, v)
    }
}

func (s *TypeSet) Declare(f *jen.File) {
    for _, v := range s.types {
        v.typeDeclaration(f)
    }
}

func NewTypeSet() *TypeSet {
    sig := make(map[string]bool)
    typ := make([]ValueConstructor,0)
    return &TypeSet { sig, typ }
}

type Statement = jen.Statement

type ValueConstructor interface {
	Signature() string
	TypeName() *Statement
	typeDeclaration(*jen.File)
	RegisterTo(s *TypeSet)
	Marshal(id string, writer string) *Statement // returns an error
	Unmarshal(reader string) *Statement          // returns (type, err)
}

func Print(v ValueConstructor) string {
	buf := bytes.NewBufferString("")
	v.TypeName().Render(buf)
	return buf.String()
}

func NewLongValue() LongValue {
	return LongValue{}
}

func NewFloatValue() FloatValue {
	return FloatValue{}
}

func NewIntValue() IntValue {
	return IntValue{}
}

func NewStringValue() StringValue {
	return StringValue{}
}

func NewVoidValue() VoidValue {
	return VoidValue{}
}

func NewValueValue() ValueValue {
	return ValueValue{}
}

func NewBoolValue() BoolValue {
	return BoolValue{}
}

func NewListValue(value ValueConstructor) *ListValue {
	return &ListValue{value}
}

func NewMapValue(key, value ValueConstructor) *MapValue {
	return &MapValue{key, value}
}

func NewMemberValue(name string, value ValueConstructor) MemberValue {
	return MemberValue{name, value}
}

func NewStrucValue(name string, members []MemberValue) *StructValue {
	return &StructValue{name, members}
}

func NewTupleValue(values []ValueConstructor) *TupleValue {
	return &TupleValue{values}
}

type IntValue struct {
}

func (i IntValue) Signature() string {
	return "I"
}

func (i IntValue) TypeName() *Statement {
	return jen.Uint32()
}

func (i IntValue) RegisterTo(s *TypeSet) {
	return
}

func (i IntValue) typeDeclaration(file *jen.File) {
	return
}

func (i IntValue) Marshal(id string, writer string) *Statement {
	return jen.Qual("qiloop/basic", "WriteUint32").Call(jen.Id(id), jen.Id(writer))
}

func (i IntValue) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadUint32").Call(jen.Id(reader))
}

type LongValue struct {
}

func (i LongValue) Signature() string {
	return "L"
}

func (i LongValue) TypeName() *Statement {
	return jen.Uint64()
}

func (i LongValue) RegisterTo(s *TypeSet) {
	return
}

func (i LongValue) typeDeclaration(file *jen.File) {
	return
}

func (i LongValue) Marshal(id string, writer string) *Statement {
	return jen.Qual("qiloop/basic", "WriteUint64").Call(jen.Id(id), jen.Id(writer))
}

func (i LongValue) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadUint64").Call(jen.Id(reader))
}

type FloatValue struct {
}

func (f FloatValue) Signature() string {
	return "f"
}

func (f FloatValue) TypeName() *Statement {
	return jen.Float32()
}

func (f FloatValue) RegisterTo(s *TypeSet) {
	return
}

func (f FloatValue) typeDeclaration(file *jen.File) {
	return
}

func (f FloatValue) Marshal(id string, writer string) *Statement {
	return jen.Qual("qiloop/basic", "WriteFloat32").Call(jen.Id(id), jen.Id(writer))
}

func (f FloatValue) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadFloat32").Call(jen.Id(reader))
}

type BoolValue struct {
}

func (b BoolValue) Signature() string {
	return "b"
}

func (b BoolValue) TypeName() *Statement {
	return jen.Bool()
}

func (b BoolValue) RegisterTo(s *TypeSet) {
	return
}

func (b BoolValue) typeDeclaration(file *jen.File) {
	return
}

func (b BoolValue) Marshal(id string, writer string) *Statement {
	return jen.Qual("qiloop/basic", "WriteBool").Call(jen.Id(id), jen.Id(writer))
}

func (b BoolValue) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadBool").Call(jen.Id(reader))
}

type ValueValue struct {
}

func (b ValueValue) Signature() string {
	return "m"
}

func (b ValueValue) TypeName() *Statement {
	return jen.Qual("qiloop/value", "Value")
}

func (b ValueValue) RegisterTo(s *TypeSet) {
	return
}

func (b ValueValue) typeDeclaration(file *jen.File) {
	return
}

func (b ValueValue) Marshal(id string, writer string) *Statement {
	return jen.Id(id).Dot("Write").Call(jen.Id(writer))
}

func (b ValueValue) Unmarshal(reader string) *Statement {
    return jen.Qual("qiloop/value", "NewValue").Call(jen.Id(reader))
}

type VoidValue struct {
}

func (v VoidValue) Signature() string {
	return "v"
}

func (v VoidValue) TypeName() *Statement {
	return jen.Empty()
}

func (v VoidValue) RegisterTo(s *TypeSet) {
	return
}

func (v VoidValue) typeDeclaration(file *jen.File) {
	return
}

func (v VoidValue) Marshal(id string, writer string) *Statement {
	return jen.Nil()
}

func (v VoidValue) Unmarshal(reader string) *Statement {
	return jen.Empty()
}

type StringValue struct {
}

func (s StringValue) Signature() string {
	return "s"
}

func (s StringValue) TypeName() *Statement {
	return jen.String()
}

func (s StringValue) RegisterTo(set *TypeSet) {
	return
}

func (s StringValue) typeDeclaration(file *jen.File) {
	return
}

func (s StringValue) Marshal(id string, writer string) *Statement {
	return jen.Id("basic.WriteString").Call(jen.Id(id), jen.Id(writer))
}

func (s StringValue) Unmarshal(reader string) *Statement {
	return jen.Id("basic.ReadString").Call(jen.Id(reader))
}

type ListValue struct {
	value ValueConstructor
}

func (l *ListValue) Signature() string {
	return fmt.Sprintf("[%s]", l.value.Signature())
}

func (l *ListValue) TypeName() *Statement {
	return jen.Index().Add(l.value.TypeName())
}

func (l *ListValue) RegisterTo(s *TypeSet) {
    l.value.RegisterTo(s)
	return
}

func (l *ListValue) typeDeclaration(file *jen.File) {
    return
}

func (l *ListValue) Marshal(listId string, writer string) *Statement {
	return jen.Func().Params().Params(jen.Error()).Block(
        jen.Err().Op(":=").Qual("qiloop/basic", "WriteUint32").Call(jen.Id("uint32").Call(
			jen.Id("len").Call(jen.Id(listId))),
			jen.Id(writer)),
		jen.Id(`if (err != nil) {
            return fmt.Errorf("failed to write slice size: %s", err)
        }`),
		jen.For(
			jen.Id("_, v := range "+listId),
		).Block(
			jen.Err().Op("=").Add(l.value.Marshal("v", writer)),
			jen.Id(`if (err != nil) {
                return fmt.Errorf("failed to write slice value: %s", err)
            }`),
		),
		jen.Return(jen.Nil()),
	).Call()
}

func (l *ListValue) Unmarshal(reader string) *Statement {
	return jen.Func().Params().Params(
		jen.Id("b").Index().Add(l.value.TypeName()),
		jen.Err().Error(),
	).Block(
		jen.Id("size, err := basic.ReadUint32").Call(jen.Id(reader)),
		jen.If(jen.Id("err != nil")).Block(
			jen.Return(jen.Id("b"), jen.Qual("fmt", "Errorf").Call(jen.Id(`"failed to read slice size: %s", err`)))),
		jen.Id("b").Op("=").Id("make").Call(l.TypeName(), jen.Id("size")),
		jen.For(
			jen.Id("i := 0; i < int(size); i++"),
		).Block(
			jen.Id("b[i], err =").Add(l.value.Unmarshal(reader)),
			jen.Id(`if (err != nil) {
                return b, fmt.Errorf("failed to read slice value: %s", err)
            }`),
		),
		jen.Return(jen.Id("b"), jen.Nil()),
	).Call()
}

type MapValue struct {
	key    ValueConstructor
	value  ValueConstructor
}

func (m *MapValue) Signature() string {
	return fmt.Sprintf("{%s%s}", m.key.Signature(), m.value.Signature())
}

func (m *MapValue) TypeName() *Statement {
	return jen.Map(m.key.TypeName()).Add(m.value.TypeName())
}

func (m *MapValue) RegisterTo(s *TypeSet) {
    m.key.RegisterTo(s)
    m.value.RegisterTo(s)
	return
}

func (m *MapValue) typeDeclaration(file *jen.File) {
    return
}

func (m *MapValue) Marshal(mapId string, writer string) *Statement {
	return jen.Func().Params().Params(jen.Error()).Block(
        jen.Err().Op(":=").Qual("qiloop/basic", "WriteUint32").Call(jen.Id("uint32").Call(
			jen.Id("len").Call(jen.Id(mapId))),
			jen.Id(writer)),
		jen.Id(`if (err != nil) {
            return fmt.Errorf("failed to write map size: %s", err)
        }`),
		jen.For(
			jen.Id("k, v := range "+mapId),
		).Block(
			jen.Err().Op("=").Add(m.key.Marshal("k", writer)),
			jen.Id(`if (err != nil) {
                return fmt.Errorf("failed to write map key: %s", err)
            }`),
			jen.Err().Op("=").Add(m.value.Marshal("v", writer)),
			jen.Id(`if (err != nil) {
                return fmt.Errorf("failed to write map value: %s", err)
            }`),
		),
		jen.Return(jen.Nil()),
	).Call()
}

func (m *MapValue) Unmarshal(reader string) *Statement {
	return jen.Func().Params().Params(
		jen.Id("m").Map(m.key.TypeName()).Add(m.value.TypeName()),
		jen.Err().Error(),
	).Block(
		jen.Id("size, err := basic.ReadUint32").Call(jen.Id(reader)),
		jen.If(jen.Id("err != nil")).Block(
			jen.Return(jen.Id("m"), jen.Qual("fmt", "Errorf").Call(jen.Id(`"failed to read map size: %s", err`)))),
		jen.Id("m").Op("=").Id("make").Call(m.TypeName(), jen.Id("size")),
		jen.For(
			jen.Id("i := 0; i < int(size); i++"),
		).Block(
			jen.Id("k, err :=").Add(m.key.Unmarshal(reader)),
			jen.Id(`if (err != nil) {
                return m, fmt.Errorf("failed to read map key: %s", err)
            }`),
			jen.Id("v, err :=").Add(m.value.Unmarshal(reader)),
			jen.Id(`if (err != nil) {
                return m, fmt.Errorf("failed to read map value: %s", err)
            }`),
			jen.Id("m[k] = v"),
		),
		jen.Return(jen.Id("m"), jen.Nil()),
	).Call()
}

type MemberValue struct {
	Name  string
	Value ValueConstructor
}

func (m MemberValue) Title() string {
    return strings.Title(m.Name)
}

type TupleValue struct {
	values []ValueConstructor
}

func (t *TupleValue) Signature() string {
	sig := "("
	for _, v := range t.values {
        sig += v.Signature()
    }
	sig += ")"
    return sig
}

func (t *TupleValue) Members() []MemberValue {
	members := make([]MemberValue, len(t.values))
	for i, v := range t.values {
        members[i] = MemberValue{ fmt.Sprintf("p%d", i), v }
	}
    return members
}

func (t *TupleValue) Params() *Statement {
	arguments := make([]jen.Code, len(t.values))
	for i, v := range t.values {
        arguments[i] = jen.Id(fmt.Sprintf("p%d", i)).Add(v.TypeName())
	}
	return jen.Params(arguments...)
}

func (t *TupleValue) TypeName() *Statement {
	return jen.Id("...interface{}")
}

func (t *TupleValue) RegisterTo(s *TypeSet) {
    for _, v := range t.values {
        v.RegisterTo(s)
    }
    return
}

func (s *TupleValue) typeDeclaration(*jen.File) {
    return
}

func (s *TupleValue) Marshal(variadicIdentifier string, writer string) *Statement {
    // TODO: shall returns an error
    return jen.Empty()
}

func (s *TupleValue) Unmarshal(reader string) *Statement {
    // TODO: shall returns (type, err)
    return jen.Empty()
}


type StructValue struct {
	name    string
	members []MemberValue
}

func (s *StructValue) Signature() string {
	types := ""
	names := make([]string, 0, len(s.members))
	for _, v := range s.members {
		names = append(names, v.Name)
		if s, ok := v.Value.(*StructValue); ok {
			types += "[" + s.Signature() + "]"
		} else {
			types += v.Value.Signature()
		}
	}
	return fmt.Sprintf("(%s)<%s,%s>", types,
		s.name, strings.Join(names, ","))
}

func (s *StructValue) TypeName() *Statement {
	return jen.Id(s.name)
}

func (t *StructValue) RegisterTo(s *TypeSet) {
    for _, v := range t.members {
        v.Value.RegisterTo(s)
    }
    s.Register(t)
    return
}

func (s *StructValue) typeDeclaration(file *jen.File) {
	fields := make([]jen.Code, len(s.members))
	for i, v := range s.members {
		fields[i] = jen.Id(v.Title()).Add(v.Value.TypeName())
	}
	file.Type().Id(s.name).Struct(fields...)

	readFields := make([]jen.Code, len(s.members)+1)
	writeFields := make([]jen.Code, len(s.members)+1)
	for i, v := range s.members {
		readFields[i] = jen.If(
			jen.Id("s." + v.Title() + ", err =").Add(v.Value.Unmarshal("r")),
			jen.Id("err != nil")).Block(
			jen.Id(`return s, fmt.Errorf("failed to read ` + v.Title() + ` field: %s", err)`),
		)
		writeFields[i] = jen.If(
			jen.Id("err :=").Add(v.Value.Marshal("s."+ v.Title(), "w")),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Id(`return fmt.Errorf("failed to write ` + v.Title() + ` field: %s", err)`),
		)
	}
	readFields[len(s.members)] = jen.Return(jen.Id("s"), jen.Nil())
	writeFields[len(s.members)] = jen.Return(jen.Nil())

	file.Func().Id("Read"+s.name).Params(
		jen.Id("r").Id("io.Reader"),
	).Params(
		jen.Id("s").Id(s.name), jen.Err().Error(),
	).Block(readFields...)
	file.Func().Id("Write"+s.name).Params(
		jen.Id("s").Id(s.name),
		jen.Id("w").Qual("io", "Writer"),
	).Params(jen.Err().Error()).Block(writeFields...)
}

func (s *StructValue) Marshal(strucId string, writer string) *Statement {
	return jen.Id("Write"+s.name).Call(jen.Id(strucId), jen.Id(writer))
}

func (s *StructValue) Unmarshal(reader string) *Statement {
	return jen.Id("Read" + s.name).Call(jen.Id(reader))
}

func BasicType() parsec.Parser {
	return parsec.OrdChoice(nodifyBasicType,
		parsec.Atom("I", "uint32"),
		parsec.Atom("s", "string"),
		parsec.Atom("L", "uint64"),
		parsec.Atom("b", "bool"),
		parsec.Atom("f", "float32"),
		parsec.Atom("m", "value"),
		parsec.Atom("v", "void"))
}

func TypeName() parsec.Parser {
	return parsec.Ident()
}

type Node = parsec.ParsecNode

func nodifyBasicType(nodes []Node) Node {
	if len(nodes) != 1 {
		log.Panicf("wrong basic arguments %+v\n", nodes)
	}
	signature := nodes[0].(*parsec.Terminal).GetValue()
	switch signature {
	case "I":
		return NewIntValue()
	case "L":
		return NewLongValue()
	case "s":
		return NewStringValue()
	case "b":
		return NewBoolValue()
	case "f":
		return NewFloatValue()
	case "v":
		return NewVoidValue()
	case "m":
		return NewValueValue()
	default:
		log.Panicf("wrong signature %s", signature)
	}
	return nil
}

func extractValue(object interface{}) (ValueConstructor, error) {
	nodes, ok := object.([]Node)
	if !ok {
		return nil, fmt.Errorf("extraction failed: %+v", reflect.TypeOf(object))
	}
	value, ok := nodes[0].(ValueConstructor)
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
	return NewMapValue(key, value)
}

func nodifyArrayType(nodes []Node) Node {
	if len(nodes) != 3 {
		log.Panicf("wrong arguments %+v", nodes)
	}
	value, err := extractValue(nodes[1])
	if err != nil {
		log.Panicf("value conversion failed: %s", err)
	}
	return NewListValue(value)
}

func extractMembersTypes(node Node) []ValueConstructor {
	typeList, ok := node.([]Node)
	if !ok {
		log.Panicf("member type list is not a list: %s", reflect.TypeOf(node))
	}
	types := make([]ValueConstructor, len(typeList))
	for i, _ := range typeList {
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

func extractMembers(typesNode, namesNode Node) []MemberValue {

	types := extractMembersTypes(typesNode)
	names := extractMembersName(namesNode)

	if len(types) != len(names) {
		log.Panicf("member types and names of different size: %+v, %+v", types, names)
	}

	members := make([]MemberValue, len(names))

	for i := range types {
		members[i] = NewMemberValue(names[i], types[i])
	}
	return members
}

func nodifyTupleType(nodes []Node) Node {

	types := extractMembersTypes(nodes[1])

	return NewTupleValue(types)
}

func nodifyTypeDefinition(nodes []Node) Node {

	terminal, ok := nodes[4].(*parsec.Terminal)
	if !ok {
		log.Panicf("wrong name %s", reflect.TypeOf(nodes[4]))
	}
	name := terminal.GetValue()
	members := extractMembers(nodes[1], nodes[6])

	return NewStrucValue(name, members)
}

func Parse(input string) (ValueConstructor, error) {
	text := []byte(input)

	var arrayType parsec.Parser
	var mapType parsec.Parser
	var typeDefinition parsec.Parser
	var tupleType parsec.Parser

	var declarationType = parsec.OrdChoice(nil,
		BasicType(), &mapType, &arrayType, &typeDefinition, &tupleType)

	arrayType = parsec.And(nodifyArrayType,
		parsec.Atom("[", "MapStart"),
		&declarationType,
		parsec.Atom("]", "MapClose"))

	var listType = parsec.Kleene(nil, &declarationType)

	var typeMemberList = parsec.Maybe(nil,
		parsec.Many(nil, TypeName(), parsec.Atom(",", "Separator")))

	tupleType = parsec.And(nodifyTupleType,
		parsec.Atom("(", "TypeParameterStart"),
		&listType,
		parsec.Atom(")", "TypeParameterClose"))

	typeDefinition = parsec.And(nodifyTypeDefinition,
		parsec.Atom("(", "TypeParameterStart"),
		&listType,
		parsec.Atom(")", "TypeParameterClose"),
		parsec.Atom("<", "TypeDefinitionStart"),
		TypeName(),
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
		return nil, fmt.Errorf("failed to parse signature")
	}
	types, ok := root.([]Node)
	if !ok {
		return nil, fmt.Errorf("failed to convert array: %+v", reflect.TypeOf(root))
	}
	if len(types) != 1 {
		return nil, fmt.Errorf("did not parse only one type: %+v", root)
	}
	constructor, ok := types[0].(ValueConstructor)
	if !ok {
		return nil, fmt.Errorf("failed to convert value: %+v", reflect.TypeOf(types[0]))
	}
	return constructor, nil
}

func GenerateType(v ValueConstructor, packageName string, w io.Writer) error {
    var file *jen.File = jen.NewFile(packageName)
    set := NewTypeSet()
    v.RegisterTo(set)
    set.Declare(file)
    if err := file.Render(w); err != nil {
        return fmt.Errorf("failed to render %s: %s", v.Signature(), err)
    }
    return nil
}
