package metaqi

import (
    "bytes"
    "fmt"
    "reflect"
    "io"
    "log"
    "metaqi/basic"
    parsec "github.com/prataprc/goparsec"
    "github.com/dave/jennifer/jen"
    "strings"
)

/*
Bootstrap stages:

    1. Extract the MetaObject signature using wireshark
    2. Generate the code of MetaObject type and MetaObject constructor:
        type MetaObject struct { ... }
        func NewMetaObject(io.Reader) MetaObject { ... }
    3. Extract the MetaObject data using wireshark
    4. Parse the MetaObject binary data and generate the ServiceDirectory proxy
        type ServiceDirectory { ... }
        func NewServiceDirectory(...) ServiceDirectory
    5. Construct the MetaObject of each services declared in the ServiceDirectory
    6. Parse the MetaObject of each service and generate the associated proxy
        type ServiceXXX { ... }
        func NewServiceXXX(...) ServiceXXX

*/

type Statement = jen.Statement

type ValueConstructor interface {
    Signature() string
    TypeName() *Statement
    TypeDeclaration(*jen.File)
    ReadFrom(r io.Reader) error
    WriteTo(w io.Writer) error
    Marshal(id string, writer string) *Statement // returns an error
    Unmarshal(writer string) *Statement // returns (type, err)
}

func Print(v ValueConstructor) string {
    buf := bytes.NewBufferString("")
    v.TypeName().Render(buf)
    return buf.String()
}

func NewIntValue() *IntValue {
    return &IntValue{ 0 }
}

func NewStringValue() *StringValue {
    return &StringValue{ "" }
}

func NewMapValue(key, value ValueConstructor) *MapValue {
    return &MapValue{ key, value, nil}
}

func NewMemberValue(name string, value ValueConstructor) MemberValue {
    return MemberValue{ name, value }
}

func NewStrucValue(name string, members []MemberValue) *StructValue {
    return &StructValue{ name, members }
}

type IntValue struct {
    value uint32
}

func (i *IntValue) Signature() string {
    return "I"
}

func (i *IntValue) TypeName() *Statement {
    return jen.Uint32()
}

func (i *IntValue) TypeDeclaration(file *jen.File) {
    return
}

func (i *IntValue) Marshal(id string, writer string) *Statement {
    return jen.Qual("metaqi/basic", "WriteUint32").Call(jen.Id(id), jen.Id(writer))
}

func (i *IntValue) Unmarshal(writer string) *Statement {
    return jen.Id("basic.ReadUint32").Call(jen.Id(writer))
}

func (i *IntValue) ReadFrom(r io.Reader) (err error) {
    i.value, err = basic.ReadUint32(r)
    return err
}

func (i *IntValue) WriteTo(w io.Writer) error {
    return basic.WriteUint32(i.value, w)
}

type StringValue struct {
    value string
}

func (i *StringValue) Signature() string {
    return "s"
}

func (i *StringValue) TypeName() *Statement {
    return jen.String()
}

func (i *StringValue) TypeDeclaration(file *jen.File) {
    return
}

func (i *StringValue) Marshal(id string, writer string) *Statement {
    return jen.Id("basic.WriteString").Call(jen.Id(id), jen.Id(writer))
}

func (i *StringValue) Unmarshal(writer string) *Statement {
    return jen.Id("basic.ReadString").Call(jen.Id(writer))
}

func (i *StringValue) ReadFrom(r io.Reader) (err error) {
    i.value, err = basic.ReadString(r)
    return err
}

func (i *StringValue) WriteTo(w io.Writer) error {
    return basic.WriteString(i.value, w)
}

type MapValue struct {
    key ValueConstructor
    value ValueConstructor
    values map[ValueConstructor]ValueConstructor
}

func (m *MapValue) Signature() string {
    return fmt.Sprintf("{%s%s}", m.key.Signature(), m.value.Signature())
}

func (m *MapValue) TypeName() *Statement {
    return jen.Map(m.key.TypeName()).Add(m.value.TypeName())
}

func (m *MapValue) TypeDeclaration(file *jen.File) {
    m.key.TypeDeclaration(file)
    m.value.TypeDeclaration(file)
}

func (m *MapValue) ReadFrom(r io.Reader) error {
    size, err := basic.ReadUint32(r)
    if err != nil {
        return fmt.Errorf("failed to read map size: %s", err)
    }
    m.values = make(map[ValueConstructor]ValueConstructor, size)
    for i := 0; i < int(size); i++ {
        if err := m.key.ReadFrom(r); err != nil {
            return fmt.Errorf("failed to read key %d: %s", i, err)
        }
        if err := m.value.ReadFrom(r); err != nil {
            return fmt.Errorf("failed to value key %d: %s", i, err)
        }
        m.values[m.key] = m.value
    }
    return nil
}

func (m *MapValue) WriteTo(w io.Writer) error {
    err := basic.WriteUint32(uint32(len(m.values)), w)
    if err != nil {
        return fmt.Errorf("failed to write map size: %s", err)
    }
    for k,v := range m.values {
        if err := k.WriteTo(w); err != nil {
            return fmt.Errorf("failed to write key: %s", err)
        }
        if err := v.WriteTo(w); err != nil {
            return fmt.Errorf("failed to write value: %s", err)
        }
    }
    return nil
}

func (m *MapValue) Marshal(mapId string, writer string) *Statement {
    return jen.Func().Params().Params(jen.Error()).Block(
        jen.Id("err := basic.WriteUint32").Call(jen.Id("uint32").Call(
            jen.Id("len").Call(jen.Id(mapId))),
            jen.Id(writer)),
        jen.For(
            jen.Id("k, v := range " + mapId),
        ).Block( // FIXME: check errors
            m.key.Marshal("k", writer),
            m.value.Marshal("v", writer),
        ),
        jen.Return(jen.Err()),
    ).Call()
}

func (m *MapValue) Unmarshal(reader string) *Statement {
    return jen.Func().Params().Params(
        jen.Map(m.key.TypeName()).Add(m.value.TypeName()),
        jen.Error(),
    ).Block(
        jen.Id("var").Id("m").Add(m.TypeName()),
        jen.Id("size, err := basic.ReadUint32").Call(jen.Id(reader)),
        jen.For(
            jen.Id("i := 0; i < int(size); i++"),
        ).Block( // FIXME: check errors
            jen.Id("k, _ :=").Add(m.key.Unmarshal(reader)),
            jen.Id("v, _ :=").Add(m.value.Unmarshal(reader)),
            jen.Id("m[k] = v"),
        ),
        jen.Return(jen.Id("m"),jen.Err()),
    ).Call()
}

type MemberValue struct {
    name string
    value ValueConstructor
}

type StructValue struct {
    name string
    members []MemberValue
}

func (s *StructValue) Signature() string {
    types := ""
    names := make([]string, 0, len(s.members))
    for _, v := range s.members {
        names = append(names, v.name)
        if s, ok := v.value.(*StructValue); ok {
            types += "[" + s.Signature() + "]"
        } else {
            types += v.value.Signature()
        }
    }
    return fmt.Sprintf("(%s)<%s,%s>", types,
        s.name, strings.Join(names, ","))
}

func (s *StructValue) TypeName() *Statement {
    return jen.Id(s.name)
}

func (s *StructValue) TypeDeclaration(file *jen.File) {
    fields := make([]jen.Code, len(s.members))
    for i, v := range s.members {
        v.value.TypeDeclaration(file)
        fields[i] = jen.Id(v.name).Add(v.value.TypeName())
    }
    file.Type().Id(s.name).Struct(fields...)


    readFields := make([]jen.Code, len(s.members) + 1)
    writeFields := make([]jen.Code, len(s.members) + 1)
    for i, v := range s.members {
        readFields[i] = jen.Id("s." + v.name + ", err =").Add(v.value.Unmarshal("r"))
        writeFields[i] = jen.Id("err =").Add(v.value.Marshal("s." + v.name, "w"))
    }
    readFields[len(s.members)] = jen.Return(jen.Id("s"), jen.Nil())
    writeFields[len(s.members)] = jen.Return(jen.Err())

    file.Func().Id("Read" + s.name).Params(
        jen.Id("r").Id("io.Reader"),
    ).Params(
        jen.Id("s").Id(s.name),jen.Err().Error(),
    ).Block(readFields...)
    file.Func().Id("Write" + s.name).Params(
        jen.Id("s").Id(s.name),
        jen.Id("w").Qual("io", "Writer"),
    ).Params(jen.Err().Error()).Block(writeFields...)
}

func (s *StructValue) ReadFrom(r io.Reader) error {
    return fmt.Errorf("not yet implemented")
}

func (s *StructValue) WriteTo(w io.Writer) error {
    return fmt.Errorf("not yet implemented")
}

func (s *StructValue) Marshal(strucId string, writer string) *Statement {
    return jen.Id("Write" + s.name).Call(jen.Id(strucId), jen.Id(writer))
}

func (s *StructValue) Unmarshal(reader string) *Statement {
    return jen.Id("Read" + s.name).Call(jen.Id(reader))
}

func BasicType() parsec.Parser {
    return parsec.OrdChoice(nodifyBasicType,
        parsec.Atom("I", "uint32"),
        parsec.Atom("s", "string"))
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
        case "I": return NewIntValue()
        case "s": return NewStringValue()
        default: log.Panicf("wrong signature %s", signature)
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

func nodifyEmbeddedType(nodes []Node) Node {
    if len(nodes) != 3 {
        log.Panicf("wrong arguments %+v", nodes)
    }
    subnode, ok := nodes[1].([]Node)
    if !ok {
        log.Panicf("embedded extraction failed %+v", reflect.TypeOf(nodes[1]))
    }
    return subnode[0]
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

    var embeddedType parsec.Parser
    var mapType parsec.Parser
    var typeDefinition parsec.Parser

    var declarationType = parsec.OrdChoice(nil,
        BasicType(), &mapType, &embeddedType, &typeDefinition)

    embeddedType = parsec.And(nodifyEmbeddedType,
        parsec.Atom("[", "MapStart"),
        &declarationType,
        parsec.Atom("]", "MapClose"))

    var listType = parsec.Kleene(nil, &declarationType)

    var typeMemberList = parsec.Maybe(nil,
        parsec.Many(nil, TypeName(), parsec.Atom(",", "Separator")))

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
    if (root == nil) {
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
