package stub

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	"io"
	"strings"
)

func stubName(name string) string {
	return "stub" + name
}

// GeneratePackage generate the stub for the package declaration.
func GeneratePackage(w io.Writer, pkg *idl.PackageDeclaration) error {
	if pkg.Name == "" {
		return fmt.Errorf("empty package name")
	}
	file := jen.NewFile(pkg.Name)
	msg := "Package " + pkg.Name + " contains a generated stub"
	file.PackageComment(msg)
	file.PackageComment("File generated. DO NOT EDIT.")

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
	if err := generateStubObject(f, itf); err != nil {
		return err
	}
	if err := generateStubMethods(f, itf); err != nil {
		return err
	}
	if err := generateStubMetaObject(f, itf); err != nil {
		return err
	}
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
	}
	return jen.Id(methodName).Add(tuple.Params()).Params(ret.TypeName(),
		jen.Error()), nil
}

func generateSignalDef(itf *idl.InterfaceType, set *signature.TypeSet,
	signal idl.Signal, signalName string) (jen.Code, error) {

	tuple := signal.Tuple()
	tuple.RegisterTo(set)

	return jen.Id(signalName).Add(tuple.Params()).Error(), nil
}

func methodBodyBlock(itf *idl.InterfaceType, method idl.Method,
	methodName string) (*jen.Statement, error) {

	writing := make([]jen.Code, 0)
	params := make([]jen.Code, 0)
	code := jen.Id("buf").Op(":=").Qual(
		"bytes", "NewBuffer",
	).Call(jen.Id("payload"))
	if len(method.Params) != 0 {
		writing = append(writing, code)
	}

	for _, param := range method.Params {
		params = append(params, jen.Id(param.Name))
		code = jen.List(jen.Id(param.Name), jen.Err()).Op(":=").Add(
			param.Type.Unmarshal("buf"),
		)
		writing = append(writing, code)
		code = jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("cannot read "+param.Name+": %s"),
					jen.Err(),
				),
			),
		)
		writing = append(writing, code)
	}
	// if has not return value
	if method.Return.Signature() == "v" {
		code = jen.Id("callErr := s.impl").Dot(methodName).Call(params...)
	} else {
		code = jen.Id("ret, callErr := s.impl").Dot(methodName).Call(params...)
	}
	writing = append(writing, code)
	code = jen.Id(`if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer`)
	writing = append(writing, code)
	if method.Return.Signature() != "v" {
		code = jen.Id("errOut").Op(":=").Add(method.Return.Marshal("ret", "&out"))
		writing = append(writing, code)

		code = jen.Id(`if errOut != nil {
		    return nil, fmt.Errorf("cannot write response: %s", errOut)
	        }`)
		writing = append(writing, code)
	}
	code = jen.Id(`return out.Bytes(), nil`)
	writing = append(writing, code)

	return jen.Block(
		writing...,
	), nil

}

func generateMethodMarshal(file *jen.File, itf *idl.InterfaceType,
	method idl.Method, methodName string) error {

	body, err := methodBodyBlock(itf, method, methodName)
	if err != nil {
		return fmt.Errorf("failed to create method body: %s", err)
	}

	file.Func().Params(jen.Id("s").Op("*").Id(stubName(itf.Name))).Id(methodName).Params(
		jen.Id("payload []byte"),
	).Params(
		jen.Id("[]byte, error"),
	).Add(body)
	return nil
}

func signalBodyBlock(itf *idl.InterfaceType, signal idl.Signal,
	signalName string) (*jen.Statement, error) {

	writing := make([]jen.Code, 0)
	code := jen.Var().Id("buf").Qual("bytes", "Buffer")
	writing = append(writing, code)

	for _, param := range signal.Params {
		code = jen.If(jen.Err().Op(":=").Add(
			param.Type.Marshal(param.Name, "&buf"),
		).Op(";").Err().Op("!=").Nil()).Block(
			jen.Id(`return fmt.Errorf("failed to serialize ` +
				param.Name + `: %s", err)`),
		)
		writing = append(writing, code)
	}
	// if has not return value
	code = jen.Id("err := s.obj.UpdateSignal").Call(
		jen.Lit(signal.ID),
		jen.Id("buf.Bytes()"),
	)
	writing = append(writing, code)
	code = jen.Id(`
	if err != nil {
	    return fmt.Errorf("failed to update ` +
		signalName + `: %s", err)
	}
	return nil`)
	writing = append(writing, code)
	return jen.Block(
		writing...,
	), nil
}

func generateSignalHelper(file *jen.File, itf *idl.InterfaceType,
	signal idl.Signal, signalName string) error {
	tuple := signal.Tuple()

	body, err := signalBodyBlock(itf, signal, signalName)
	if err != nil {
		return fmt.Errorf("failed to create signal helper body: %s", err)
	}
	file.Func().Params(
		jen.Id("s").Op("*").Id(stubName(itf.Name)),
	).Id(signalName).Add(tuple.Params()).Error().Add(body)
	return nil
}

func generateStubMethods(file *jen.File, itf *idl.InterfaceType) error {

	method := func(m object.MetaMethod, methodName string) error {
		method := itf.Methods[m.Uid]
		err := generateMethodMarshal(file, itf, method, methodName)
		if err != nil {
			return fmt.Errorf("failed to create method marshall %s of %s: %s",
				method.Name, itf.Name, err)
		}
		return nil
	}
	signal := func(s object.MetaSignal, signalName string) error {
		signal := itf.Signals[s.Uid]
		signalName = strings.Replace(signalName, "Subscribe", "Signal", 1)
		err := generateSignalHelper(file, itf, signal, signalName)
		if err != nil {
			return fmt.Errorf("failed to create signal marshall %s of %s: %s",
				signal.Name, itf.Name, err)
		}
		return nil
	}
	property := func(p object.MetaProperty, getMethodName, setMethodName,
		subscribeMethodName string) error {
		// TODO: add property methods
		return nil
	}

	meta := itf.MetaObject()
	if err := meta.ForEachMethodAndSignal(method, signal, property); err != nil {
		return fmt.Errorf("failed to generate interface object %s: %s",
			itf.Name, err)
	}
	return nil
}
func generateStubMetaObject(file *jen.File, itf *idl.InterfaceType) error {
	metaMethods := func(d jen.Dict) {
		for _, method := range itf.Methods {
			d[jen.Lit(method.ID)] = jen.Values(jen.Dict{
				jen.Id("Uid"): jen.Lit(method.ID),
				jen.Id("ReturnSignature"): jen.Lit(
					method.Return.Signature(),
				),
				jen.Id("Name"): jen.Lit(method.Name),
				jen.Id("ParametersSignature"): jen.Lit(
					method.Tuple().Signature(),
				),
			})
		}
	}
	metaSignals := func(d jen.Dict) {
		for _, signal := range itf.Signals {
			d[jen.Lit(signal.ID)] = jen.Values(jen.Dict{
				jen.Id("Uid"):  jen.Lit(signal.ID),
				jen.Id("Name"): jen.Lit(signal.Name),
				jen.Id("Signature"): jen.Lit(
					signal.Tuple().Signature(),
				),
			})
		}
	}
	file.Func().Params(
		jen.Id("s").Op("*").Id(stubName(itf.Name)),
	).Id("metaObject").Params().Params(
		jen.Qual("github.com/lugu/qiloop/type/object", "MetaObject"),
	).Block(
		jen.Return().Qual(
			"github.com/lugu/qiloop/type/object", "MetaObject",
		).Values(
			jen.Dict{
				jen.Id("Description"): jen.Lit(itf.Name),
				jen.Id("Methods"): jen.Map(jen.Uint32()).Qual(
					"github.com/lugu/qiloop/type/object", "MetaMethod",
				).Values(
					jen.DictFunc(metaMethods),
				),
				jen.Id("Signals"): jen.Map(jen.Uint32()).Qual(
					"github.com/lugu/qiloop/type/object", "MetaSignal",
				).Values(
					jen.DictFunc(metaSignals),
				),
			},
		),
	)
	return nil
}

func generateStubObject(file *jen.File, itf *idl.InterfaceType) error {
	file.Func().Params(
		jen.Id("s").Op("*").Id(stubName(itf.Name)),
	).Id("Activate").Params(
		jen.Id("activation").Qual(
			"github.com/lugu/qiloop/bus/server",
			"Activation",
		),
	).Params(
		jen.Error(),
	).Block(
		jen.Id(`s.obj.Activate(activation)`),
		jen.Id(`return s.impl.Activate(activation, s)`),
	)
	file.Func().Params(
		jen.Id("s").Op("*").Id(stubName(itf.Name)),
	).Id("OnTerminate").Params().Block(
		jen.Id(`s.impl.OnTerminate()`),
		jen.Id(`s.obj.OnTerminate()`),
	)
	file.Func().Params(
		jen.Id("s").Op("*").Id(stubName(itf.Name)),
	).Id("Receive").Params(
		jen.Id("msg").Op("*").Qual("github.com/lugu/qiloop/bus/net", "Message"),
		jen.Id("from").Op("*").Qual("github.com/lugu/qiloop/bus/server", "Context"),
	).Params(jen.Error()).Block(
		jen.Id(`return s.obj.Receive(msg, from)`),
	)
	return nil
}
func generateStubConstructor(file *jen.File, itf *idl.InterfaceType) error {
	writing := make([]jen.Code, 0)
	code := jen.Var().Id("stb").Id(stubName(itf.Name))
	writing = append(writing, code)
	code = jen.Id("stb.impl = impl")
	writing = append(writing, code)
	code = jen.Id("stb.obj").Op("=").Qual(
		"github.com/lugu/qiloop/bus/server",
		"NewObject",
	).Call(jen.Id("stb.metaObject()"))
	writing = append(writing, code)

	method := func(m object.MetaMethod, methodName string) error {
		method := itf.Methods[m.Uid]
		code = jen.Id("stb.obj.Wrap").Call(
			jen.Lit(method.ID),
			jen.Id("stb").Dot(methodName),
		)
		writing = append(writing, code)
		return nil
	}

	signal := func(m object.MetaSignal, signalName string) error {
		return nil
	}

	property := func(p object.MetaProperty, getMethodName, setMethodName,
		subscribeMethodName string) error {
		return nil
	}

	meta := itf.MetaObject()
	if err := meta.ForEachMethodAndSignal(method, signal, property); err != nil {
		return fmt.Errorf("failed to generate interface object %s: %s",
			itf.Name, err)
	}
	code = jen.Return().Op("&").Id("stb")
	writing = append(writing, code)

	file.Commentf("%s returns an object using %s", itf.Name+"Object",
		itf.Name)
	file.Func().Id(itf.Name+"Object").Params(
		jen.Id("impl").Id(itf.Name),
	).Qual(
		"github.com/lugu/qiloop/bus/server", "Object",
	).Block(writing...)
	return nil
}

func generateStubType(file *jen.File, itf *idl.InterfaceType) error {
	file.Commentf("%s implements server.Object.", stubName(itf.Name))
	file.Type().Id(stubName(itf.Name)).Struct(
		jen.Id("obj").Op("*").Qual(
			"github.com/lugu/qiloop/bus/server", "BasicObject",
		),
		jen.Id("impl").Id(itf.Name),
	)
	return nil
}

func generateObjectInterface(file *jen.File, set *signature.TypeSet,
	itf *idl.InterfaceType) error {

	// Proxy and stub shall generate the name method name: reuse
	// the MetaObject method ForEachMethodAndSignal to get an
	// ordered list of the method with uniq name.
	definitions := make([]jen.Code, 0)
	signalDefinitions := make([]jen.Code, 0)
	activate := jen.Id("Activate").Params(
		jen.Id("activation").Qual("github.com/lugu/qiloop/bus/server",
			"Activation"),
		jen.Id("helper").Id(itf.Name+"SignalHelper"),
	).Params(
		jen.Error(),
	)
	definitions = append(definitions, activate)
	terminate := jen.Id("OnTerminate()")
	definitions = append(definitions, terminate)

	method := func(m object.MetaMethod, methodName string) error {
		method := itf.Methods[m.Uid]
		def, err := generateMethodDef(itf, set, method, methodName)
		if err != nil {
			return fmt.Errorf("failed to render method definition %s of %s: %s",
				method.Name, itf.Name, err)
		}
		definitions = append(definitions, def)
		return nil
	}
	signal := func(s object.MetaSignal, signalName string) error {
		signal := itf.Signals[s.Uid]
		signalName = strings.Replace(signalName, "Subscribe", "Signal", 1)
		def, err := generateSignalDef(itf, set, signal, signalName)
		if err != nil {
			return fmt.Errorf("failed to render signal %s of %s: %s", s.Name,
				itf.Name, err)
		}
		signalDefinitions = append(signalDefinitions, def)
		return nil
	}
	property := func(p object.MetaProperty, getMethodName, setMethodName,
		subscribeMethodName string) error {
		// TODO: force an initial value to be passed during object
		// creation.
		return nil
	}

	meta := itf.MetaObject()
	if err := meta.ForEachMethodAndSignal(method, signal, property); err != nil {
		return fmt.Errorf("failed to generate interface object %s: %s",
			itf.Name, err)
	}

	file.Commentf("%s interface of the service implementation", itf.Name)
	file.Type().Id(itf.Name).Interface(
		definitions...,
	)

	file.Commentf("%s provided to %s a companion object",
		itf.Name+"SignalHelper", itf.Name)
	file.Type().Id(itf.Name + "SignalHelper").Interface(
		signalDefinitions...,
	)
	return nil
}
