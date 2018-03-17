package idl

import (
	"fmt"
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	"io"
)

// generateMethod writes the method declaration. Does not use
// methodName, because is the Go method name generated to avoid
// conflicts. QiMessage do not have such constraint and thus we don't
// use this name when creating IDL files.
func generateMethod(writer io.Writer, set *signature.TypeSet, m object.MetaMethod, methodName string) error {
	paramType, err := signature.Parse(m.ParametersSignature)
	if err != nil {
		return fmt.Errorf("failed to parse parms of %s: %s", m.Name, err)
	}
	retType, err := signature.Parse(m.ReturnSignature)
	if err != nil {
		return fmt.Errorf("failed to parse return of %s: %s", m.Name, err)
	}
	if retType.Signature() == "v" {
		fmt.Fprintf(writer, "\tfn %s(%s)\n", m.Name, paramType.SignatureIDL())
	} else {
		fmt.Fprintf(writer, "\tfn %s(%s) -> %s\n", m.Name, paramType.SignatureIDL(), retType.SignatureIDL())
	}
	paramType.RegisterTo(set)
	retType.RegisterTo(set)
	return nil
}

// generateSignal writes the signal declaration. Does not use
// methodName, because is the Go method name generated to avoid
// conflicts. QiMessage do not have such constraint and thus we don't
// use this name when creating IDL files.
func generateSignal(writer io.Writer, set *signature.TypeSet, s object.MetaSignal, methodName string) error {
	signalType, err := signature.Parse(s.Signature)
	if err != nil {
		return fmt.Errorf("failed to parse signal of %s: %s", s.Name, err)
	}
	signalType.RegisterTo(set)
	fmt.Fprintf(writer, "\tsig %s(%s)\n", s.Name, signalType.SignatureIDL())
	return nil
}

func generateStructure(writer io.Writer, s *signature.StructType) error {
	fmt.Fprintf(writer, "struct %s\n", s.Name)
	for _, mem := range s.Members {
		fmt.Fprintf(writer, "\t%s: %s\n", mem.Name, mem.Value.SignatureIDL())
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

// GenerateIDL writes the IDL definition of a MetaObject into a
// writer. This IDL definition can be used to re-create the MetaObject
// with the method ParseIDL.
func GenerateIDL(writer io.Writer, serviceName string, metaObj object.MetaObject) error {
	set := signature.NewTypeSet()

	fmt.Fprintf(writer, "interface %s\n", serviceName)

	methodCall := func(m object.MetaMethod, methodName string) error {
		return generateMethod(writer, set, m, methodName)
	}
	signalCall := func(s object.MetaSignal, methodName string) error {
		return generateSignal(writer, set, s, methodName)
	}

	if err := metaObj.ForEachMethodAndSignal(methodCall, signalCall); err != nil {
		return fmt.Errorf("failed to generate proxy object %s: %s", serviceName, err)
	}
	fmt.Fprintf(writer, "end\n")

	if err := generateStructures(writer, set); err != nil {
		return fmt.Errorf("failed to generate structures: %s", err)
	}
	return nil
}
