package metaqi

import (
    "encoding/binary"
    "fmt"
    "io"
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

func NewIntValue() IntValue {
    return IntValue{ 0 }
}

func NewStringValue() StringValue {
    return StringValue{ "" }
}

func MapStringValue(key, value ValueConstructor) MapValue {
    return MapValue{ key, value, nil}
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
    return fmt.Sprintf("(%s)<%s,%s", s.membersType(),
        s.TypeName(), s.membersName())
}

func (s *StructValue) membersType() string {
    types := ""
    for _, v := range s.members {
        types += v.value.TypeName()
    }
    return types
}

func (s *StructValue) membersName() string {
    names := make([]string, 0, len(s.members))
    for _, v := range s.members {
        names = append(names, v.value.TypeName())
    }
    return strings.Join(names, ", ")
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

func String() parsec.Parser {
    return parsec.Atom("s", "string")
}
func Int() parsec.Parser {
    return parsec.Atom("I", "uint32")
}
func BasicType(ast *parsec.AST) parsec.Parser {
    return ast.OrdChoice("BasicType", nil, Int(), String())
}
func TypeName(ast *parsec.AST) parsec.Parser {
    return parsec.Ident()
}

func parse(input string) string {
    text := []byte(input)
    ast := parsec.NewAST("signature", 100)

    var embeddedType parsec.Parser
    var mapType parsec.Parser
    var typeDefinition parsec.Parser

    var declarationType = ast.OrdChoice("DeclarationType", nil,
        BasicType(ast), &mapType, &embeddedType, &typeDefinition)

    embeddedType = ast.And("EmbeddedType", nil,
        parsec.Atom("[", "MapStart"),
        &declarationType,
        parsec.Atom("]", "MapClose"))

    var listType = ast.Kleene("ListType", nil, &declarationType)

    var typeMemberList = ast.Maybe("MaybeMemberList", nil,
        ast.Many("TypeMemberList", nil,
        TypeName(ast), parsec.Atom(",", "Separator")))

    typeDefinition = ast.And("TypeDefinition", nil,
        parsec.Atom("(", "TypeParameterStart"),
        &listType,
        parsec.Atom(")", "TypeParameterClose"),
        parsec.Atom("<", "TypeDefinitionStart"),
        TypeName(ast),
        parsec.Atom(",", "TypeName"),
        &typeMemberList,
        parsec.Atom(">", "TypeDefinitionClose"))

    mapType = ast.And("MapType", nil,
        parsec.Atom("{", "MapStart"),
        &declarationType, &declarationType,
        parsec.Atom("}", "MapClose"))

    var typeSignature = declarationType

    root, scanner := ast.Parsewith(typeSignature, parsec.NewScanner(text))
    if (root != nil) {
        return root.GetName()
    }
    fmt.Println(scanner.GetCursor())
    return "not recognized"
}
