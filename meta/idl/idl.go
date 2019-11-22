package idl

import (
	"fmt"
	"io"
	"strings"

	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
)

func cleanVarName(i int, name string) string {
	if name == "" {
		return fmt.Sprintf("P%d", i)
	}
	// remove space from parameter name
	name = strings.ReplaceAll(name, " ", "")
	for _, keyword := range []string{
		"break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct", "chan",
		"else", "goto", "package", "switch", "const",
		"fallthrough", "if", "range", "type", "continue",
		"for", "import", "return", "var",
	} {
		if name == keyword {
			name = fmt.Sprintf("%s_%d", name, i)
			break
		}
	}
	return name
}

// generateMethod writes the method declaration. Does not use
// methodName, because is the Go method name generated to avoid
// conflicts. QiMessage do not have such constraint and thus we don't
// use this name when creating IDL files.
func generateMethod(writer io.Writer, set *signature.TypeSet,
	m object.MetaMethod, methodName string) error {

	// paramType is a tuple it needs to be unified with
	// m.MetaMethodParameter.
	paramType, err := signature.Parse(m.ParametersSignature)
	if err != nil {
		return fmt.Errorf("parse parms of %s: %s", m.Name, err)
	}
	retType, err := signature.Parse(m.ReturnSignature)
	if err != nil {
		return fmt.Errorf("parse return of %s: %s", m.Name, err)
	}

	tupleType, ok := paramType.(*signature.TupleType)
	if !ok {
		// some buggy service don' t return tupples
		tupleType = signature.NewTupleType([]signature.Type{paramType})
	}

	paramSignature := ""
	if m.Parameters == nil || len(m.Parameters) != len(tupleType.Members) {
		paramSignature = tupleType.ParamIDL()
	} else {
		for i, p := range m.Parameters {
			if paramSignature != "" {
				paramSignature += ","
			}
			name := cleanVarName(i, p.Name)
			paramSignature += name + ": " + tupleType.Members[i].Type.SignatureIDL()
		}
	}

	returnSignature := "-> " + retType.SignatureIDL() + " "
	if retType.Signature() == "v" {
		returnSignature = ""
	}

	fmt.Fprintf(writer, "\tfn %s(%s) %s//uid:%d\n", m.Name, paramSignature, returnSignature, m.Uid)

	paramType.RegisterTo(set)
	retType.RegisterTo(set)
	return nil
}

// generateProperties writes the signal declaration. Does not use
// methodName, because is the Go method name generated to avoid
// conflicts. QiMessage do not have such constraint and thus we don't
// use this name when creating IDL files.
func generateProperty(writer io.Writer, set *signature.TypeSet,
	p object.MetaProperty, propertyName string) error {

	propertyType, err := signature.Parse(p.Signature)
	if err != nil {
		return fmt.Errorf("parse property of %s: %s", p.Name, err)
	}
	propertyType.RegisterTo(set)

	tupleType, ok := propertyType.(*signature.TupleType)
	if !ok {
		tupleType = &signature.TupleType{
			Members: []signature.MemberType{
				signature.MemberType{
					Name: "param",
					Type: propertyType,
				},
			},
		}
	}

	fmt.Fprintf(writer, "\tprop %s(%s) //uid:%d\n", p.Name,
		tupleType.ParamIDL(), p.Uid)
	return nil
}

// generateSignal writes the signal declaration. Does not use
// methodName, because is the Go method name generated to avoid
// conflicts. QiMessage do not have such constraint and thus we don't
// use this name when creating IDL files.
func generateSignal(writer io.Writer, set *signature.TypeSet, s object.MetaSignal, methodName string) error {
	signalType, err := signature.Parse(s.Signature)
	if err != nil {
		return fmt.Errorf("parse signal of %s: %s", s.Name, err)
	}
	signalType.RegisterTo(set)

	tupleType, ok := signalType.(*signature.TupleType)
	if !ok {
		tupleType = signature.NewTupleType([]signature.Type{signalType})
	}
	fmt.Fprintf(writer, "\tsig %s(%s) //uid:%d\n", s.Name, tupleType.ParamIDL(), s.Uid)
	return nil
}

func generateStructure(writer io.Writer, s *signature.StructType) error {
	fmt.Fprintf(writer, "struct %s\n", s.Name)
	for _, mem := range s.Members {
		fmt.Fprintf(writer, "\t%s: %s\n", mem.Name, mem.Type.SignatureIDL())
	}
	fmt.Fprintf(writer, "end\n")
	return nil
}

func generateStructures(writer io.Writer, set *signature.TypeSet) error {
	for _, typ := range set.Types {
		if structure, ok := typ.(*signature.StructType); ok {
			generateStructure(writer, structure)
		}
	}
	return nil
}

// GenerateIDL writes the IDL definition for the meta object in objs.  a MetaObject into a
// writer. This IDL definition can be used to re-create the MetaObject
// with the method ParseIDL.
func GenerateIDL(writer io.Writer, packageName string, objs map[string]object.MetaObject) error {
	set := signature.NewTypeSet()

	fmt.Fprintf(writer, "package %s\n", packageName)

	for name, meta := range objs {

		fmt.Fprintf(writer, "interface %s\n", name)

		method := func(m object.MetaMethod, methodName string) error {
			return generateMethod(writer, set, m, methodName)
		}
		signal := func(s object.MetaSignal, signalName string) error {
			return generateSignal(writer, set, s, "Subscribe"+signalName)
		}
		property := func(p object.MetaProperty, propertyName string) error {
			return generateProperty(writer, set, p, propertyName)
		}

		if err := meta.ForEachMethodAndSignal(method, signal, property); err != nil {
			return fmt.Errorf("generate proxy object %s: %s", name, err)
		}
		fmt.Fprintf(writer, "end\n")
	}

	if err := generateStructures(writer, set); err != nil {
		return fmt.Errorf("generate structures: %s", err)
	}
	return nil
}
