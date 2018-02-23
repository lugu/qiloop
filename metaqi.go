package metaqi

import (
    "encoding/binary"
    "fmt"
    "reflect"
    "io"
    "log"
    parsec "github.com/prataprc/goparsec"
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

interface Type {
    func writeTo(io.Writer)
}

interface TypeConstructor {
    func readFrom(io.Reader) Type
}

interface TypeDeclarator {
    func typeName() string
    // Generates code for a Type and a TypeConstructor
    func declareType(io.Writer)
}

func parseSignature(string) TypeDeclaration
func readMetaObject(io.Reader) MetaObject

*/


type (
    MetaMethodParameter struct {
        Name string
        Description string
    }

    MetaMethod struct {
        Uuid uint32
        ReturnSignature string
        Name string
        ParametersSignature string
        Description string
        Parameters MetaMethodParameter
        ReturnDescription string
    }

    MetaSignal struct {
        Uuid uint32
        Name string
        Signature string
    }

    MetaProperty struct {
        Uuid uint32
        Name string
        Signature string
    }

    MetaObject struct {
        Methods map[int] MetaMethod
        Signals map[int] MetaSignal
        Properties map[int] MetaProperty
        Description string
    }
)

type ValueConstructor interface {
    Signature() string
    TypeName() string
    TypeDeclaration() string
    ReadFrom(r io.Reader) error
    WriteTo(w io.Writer) error
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

func (i *IntValue) TypeName() string {
    return "uint32"
}

func (i *IntValue) TypeDeclaration() string {
    return ""
}

func (i *IntValue) ReadFrom(r io.Reader) error {
        buf := []byte{0, 0, 0, 0}
        bytes, err := r.Read(buf)
        if (err != nil) {
            return err
        } else if (bytes != 4) {
            return fmt.Errorf("IntValue read %d instead of 4", bytes)
        }
        i.value = binary.BigEndian.Uint32(buf)
        return nil
}

func (i *IntValue) WriteTo(r io.Writer) error {
        buf := []byte{0, 0, 0, 0}
        binary.BigEndian.PutUint32(buf, i.value)
        bytes, err := r.Write(buf)
        if (err != nil) {
            return err
        } else if (bytes != 4) {
            return fmt.Errorf("IntValue write %d instead of 4", bytes)
        }
        return nil
}

type StringValue struct {
    value string
}

func (i *StringValue) Signature() string {
    return "s"
}

func (i *StringValue) TypeName() string {
    return "string"
}

func (i *StringValue) TypeDeclaration() string {
    return ""
}

func (i *StringValue) ReadFrom(r io.Reader) error {
        var size IntValue
        if err := size.ReadFrom(r); err != nil {
            return fmt.Errorf("StringValue failed to read size: %s", err)
        }
        buf := make([]byte, size.value, size.value + 1)
        bytes, err := r.Read(buf)
        if (err != nil) {
            return err
        } else if (bytes != int(size.value)) {
            return fmt.Errorf("StringValue read %d instead of %d", bytes, size.value)
        }
        buf[size.value] = 0
        i.value = string(buf)
        return nil
}

func (i *StringValue) WriteTo(r io.Writer) error {
        size := IntValue{ uint32(len(i.value)) }
        if err := size.WriteTo(r); err != nil {
            return fmt.Errorf("StringValue failed to write size: %s", err)
        }
        bytes, err := r.Write([]byte(i.value))
        if (err != nil) {
            return err
        } else if (bytes != 4) {
            return fmt.Errorf("StringValue write %d instead of %d", bytes, size.value)
        }
        return nil
}

type MapValue struct {
    key ValueConstructor
    value ValueConstructor
    values map[ValueConstructor]ValueConstructor
}

func (m *MapValue) Signature() string {
    return fmt.Sprintf("{%s%s}", m.key.Signature(), m.value.Signature())
}

func (m *MapValue) TypeName() string {
    return fmt.Sprintf("map[%s] %s", m.key.TypeName(), m.value.TypeName())
}

func (m *MapValue) TypeDeclaration() string {
    return fmt.Sprintf("%s\n%s\n",
        m.key.TypeDeclaration(), m.value.TypeDeclaration())
}

func (m *MapValue) ReadFrom(r io.Reader) error {
        var size IntValue
        if err := size.ReadFrom(r); err != nil {
            return err
        }
        m.values = make(map[ValueConstructor]ValueConstructor, size.value)
        for i := 0; i < int(size.value); i++ {
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
        size := IntValue{ uint32(len(m.values)) }
        if err := size.WriteTo(w); err != nil {
            return fmt.Errorf("StringValue failed to write size: %s", err)
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

type MemberValue struct {
    name string
    value ValueConstructor
}

func (m* MemberValue) fieldDeclaration() string {
    return fmt.Sprintf("%s %s\n", m.name, m.value.TypeName())
}

type StructValue struct {
    name string
    members []MemberValue
}

func (s *StructValue) Signature() string {
    return fmt.Sprintf("(%s)<%s,%s>", s.membersTypeSignature(),
        s.TypeName(), s.membersName())
}

func (s *StructValue) membersTypeSignature() string {
    types := ""
    for _, v := range s.members {
        if s, ok := v.value.(*StructValue); ok {
            types += "[" + s.Signature() + "]"
        } else {
            types += v.value.Signature()
        }
    }
    return types
}

func (s *StructValue) membersName() string {
    names := make([]string, 0, len(s.members))
    for _, v := range s.members {
        names = append(names, v.name)
    }
    return strings.Join(names, ",")
}

func (s *StructValue) TypeName() string {
    return s.name
}

func (s *StructValue) TypeDeclaration() string {
    fields := ""
    for _, v := range s.members {
        fields += v.fieldDeclaration()
    }
    return fmt.Sprintf("type %s struct {\n%s\n}",
        s.TypeName(),
        fields)
}

func (s *StructValue) ReadFrom(r io.Reader) error {
    return fmt.Errorf("not yet implemented")
}

func (s *StructValue) WriteTo(w io.Writer) error {
    return fmt.Errorf("not yet implemented")
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

func parse(input string) (ValueConstructor, error) {
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
