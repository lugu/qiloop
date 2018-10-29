package stub

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	"io"
)

func GeneratePackage(w io.Writer, pkg *idl.PackageDeclaration) error {
	if pkg.Name == "" {
		return fmt.Errorf("empty package name")
	}
	file := jen.NewFile(pkg.Name)
	file.PackageComment("file generated. DO NOT EDIT.")

	set := signature.NewTypeSet()
	for _, typ := range pkg.Types {
		err := generateType(file, set, typ)
		if err != nil {
			return err
		}
	}

	return file.Render(w)
}

func generateInterface(f *jen.File, set *signature.TypeSet, itf *idl.InterfaceType) error {
	if err := generateObjectInterface(f, set, itf); err != nil {
		return err
	}
	if err := generateStub(f, itf); err != nil {
		return err
	}
	return nil
}

func generateStub(f *jen.File, itf *idl.InterfaceType) error {
	if err := generateStubType(f, itf); err != nil {
		return err
	}
	if err := generateStubConstructor(f, itf); err != nil {
		return err
	}
	// TODO: generate methods serializer
	return nil
}

func generateType(f *jen.File, set *signature.TypeSet, typ signature.Type) error {

	itf, ok := typ.(*idl.InterfaceType)
	if ok {
		return generateInterface(f, set, itf)
	}
	typ.TypeDeclaration(f)
	return nil
}

func generateMethodDef(itf *idl.InterfaceType, set *signature.TypeSet,
	method idl.Method, methodName string) (jen.Code, error) {

	tuple := method.Tuple()
	ret := method.Return

	tuple.RegisterTo(set)
	method.Return.RegisterTo(set)

	if ret.Signature() == "v" {
		return jen.Id(methodName).Add(tuple.Params()).Params(jen.Error()), nil
	} else {
		return jen.Id(methodName).Add(tuple.Params()).Params(ret.TypeName(),
			jen.Error()), nil
	}
}

func generateSignalDef(itf *idl.InterfaceType, set *signature.TypeSet,
	signal idl.Signal, signalName string) (jen.Code, error) {

	tuple := signal.Tuple()
	tuple.RegisterTo(set)

	retType := jen.Params(jen.Chan().Add(tuple.TypeName()), jen.Error())
	return jen.Id(signalName).Params(jen.Id("cancel").Chan().Int()).Add(retType), nil
}

func generateStubConstructor(file *jen.File, itf *idl.InterfaceType) error {
	file.Func().Id("New"+itf.Name).Params(
		jen.Id("impl").Id(itf.Name),
	).Qual(
		"github.com/lugu/qiloop/bus/session", "Object",
	).Block(
		jen.Var().Id("stb").Id(itf.Name+"Stub"),
		jen.Id("sbt").Dot("Wrapper = bus.Wrapper(make(map[uint32]bus.ActionWrapper))"),
		// TODO: iterate over the wrapper
		jen.Return().Op("&").Id("stb"),
	)
	return nil
}

func generateStubType(file *jen.File, itf *idl.InterfaceType) error {
	file.Type().Id(itf.Name+"Stub").Struct(
		jen.Qual("github.com/lugu/qiloop/bus/session", "ObjectDispather"),
		jen.Id("impl").Id(itf.Name),
	)
	return nil
}

func generateObjectInterface(file *jen.File, set *signature.TypeSet,
	itf *idl.InterfaceType) error {

	// Proxy and stub shall generate the name method name: reuse
	// the MetaObject method ForEachMethodAndSignal to get an
	// ordered list of the method with uniq name.
	meta := itf.MetaObject()
	definitions := make([]jen.Code, 0)

	methodCall := func(m object.MetaMethod, methodName string) error {
		method := itf.Methods[m.Uid]
		def, err := generateMethodDef(itf, set, method, methodName)
		if err != nil {
			return fmt.Errorf("failed to render method definition %s of %s: %s",
				method.Name, itf.Name, err)
		}
		definitions = append(definitions, def)
		return nil
	}
	signalCall := func(s object.MetaSignal, signalName string) error {
		signal := itf.Signals[s.Uid]
		def, err := generateSignalDef(itf, set, signal, signalName)
		if err != nil {
			return fmt.Errorf("failed to render signal %s of %s: %s", s.Name,
				itf.Name, err)
		}
		definitions = append(definitions, def)
		return nil
	}

	if err := meta.ForEachMethodAndSignal(methodCall, signalCall); err != nil {
		return fmt.Errorf("failed to generate interface object %s: %s",
			itf.Name, err)
	}

	file.Type().Id(itf.Name).Interface(
		definitions...,
	)
	return nil
}
