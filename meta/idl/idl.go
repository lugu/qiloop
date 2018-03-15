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
	fmt.Fprintf(writer, "\tfn %s\n", m.Name)
	return nil
}

// generateSignal writes the signal declaration. Does not use
// methodName, because is the Go method name generated to avoid
// conflicts. QiMessage do not have such constraint and thus we don't
// use this name when creating IDL files.
func generateSignal(writer io.Writer, set *signature.TypeSet, s object.MetaSignal, methodName string) error {
	fmt.Fprintf(writer, "\tsig %s\n", s.Name)
	return nil
}

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

	return nil
}

func ParseIDL(reader io.Reader) (*object.MetaObject, error) {
	return nil, fmt.Errorf("Not yet implemented")
}
