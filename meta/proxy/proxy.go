package proxy

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/object"
	"io"
	"sort"
	"strings"
)

type Statement = jen.Statement

func newFileAndSet(packageName string) (*jen.File, *signature.TypeSet) {
	file := jen.NewFile(packageName)
	file.PackageComment("file generated. DO NOT EDIT.")
	return file, signature.NewTypeSet()
}

func generateProxyType(file *jen.File, typ string, metaObj object.MetaObject) {

	file.Type().Id(typ).Struct(jen.Qual("github.com/lugu/qiloop/bus", "Proxy"))
	file.Func().Id("New"+typ).Params(
		jen.Id("ses").Qual("github.com/lugu/qiloop/bus", "Session"),
		jen.Id("obj").Uint32(),
	).Params(
		jen.Op("*").Id(typ),
		jen.Error(),
	).Block(
		jen.List(jen.Id("proxy"), jen.Err()).Op(":=").Id("ses").Dot("Proxy").Call(
			jen.Lit(typ),
			jen.Id("obj"),
		),
		jen.Id(`if err != nil {
			return nil, fmt.Errorf("failed to contact service: %s", err)
		}`),
		jen.Id(`return &`+typ+`{ proxy }, nil`),
	)
}

func generateMethodDef(file *jen.File, set *signature.TypeSet, serviceName string, m object.MetaMethod, methodName string) (jen.Code, error) {

	paramType, err := signature.Parse(m.ParametersSignature)
	if err != nil {
		return nil, fmt.Errorf("failed to parse parameters %s: %s", m.ParametersSignature, err)
	}
	paramType.RegisterTo(set)
	tuple, ok := paramType.(*signature.TupleValue)
	if !ok {
		tuple = signature.NewTupleValue([]signature.ValueConstructor{paramType})
	}
	tuple.ConvertMetaObjects()

	ret, err := signature.Parse(m.ReturnSignature)
	if err != nil {
		return nil, fmt.Errorf("failed to parse return: %s: %s", m.ReturnSignature, err)
	}

	// implementation detail: use object.MetaObject when generating
	// proxy in order for proxy to implement the object.Object
	// interface.
	if ret.Signature() == signature.MetaObjectSignature {
		ret = signature.NewMetaObjectValue()
	}
	if _, ok := ret.(signature.VoidValue); ok {
		return jen.Id(methodName).Add(tuple.Params()).Error(), nil
	} else {
		return jen.Id(methodName).Add(tuple.Params()).Params(ret.TypeName(), jen.Error()), nil
	}
}

func generateObjectInterface(metaObj object.MetaObject, serviceName string, set *signature.TypeSet, file *jen.File) error {

	definitions := make([]jen.Code, 0)
	definitions = append(definitions, jen.Qual("github.com/lugu/qiloop/bus", "Proxy"))
	definitions = append(definitions, jen.Qual("github.com/lugu/qiloop/object", "Object"))

	methodNames := make(map[string]bool)
	keys := make([]int, 0)
	for k := range metaObj.Methods {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for id, i := range keys {
		if id < 10 {
			// those are already defined as part of object.Object
			continue
		}
		k := uint32(i)
		m := metaObj.Methods[k]

		methodName := registerName(strings.Title(m.Name), methodNames)
		def, err := generateMethodDef(file, set, serviceName, m, methodName)
		if err != nil {
			return fmt.Errorf("failed to render method definition %s of %s: %s", m.Name, serviceName, err)
		}
		definitions = append(definitions, def)
	}
	keys = make([]int, 0)
	for k := range metaObj.Signals {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, i := range keys {
		k := uint32(i)
		s := metaObj.Signals[k]
		methodName := registerName("Signal"+strings.Title(s.Name), methodNames)
		def, err := generateSignalDef(file, set, serviceName, s, methodName)
		if err != nil {
			return fmt.Errorf("failed to render signal %s of %s: %s", s.Name, serviceName, err)
		}
		definitions = append(definitions, def)
	}

	file.Type().Id("I" + serviceName).Interface(
		definitions...,
	)
	return nil
}

func generateProxyObject(metaObj object.MetaObject, serviceName string, set *signature.TypeSet, file *jen.File) error {

	generateProxyType(file, serviceName, metaObj)

	methodNames := make(map[string]bool)
	keys := make([]int, 0)
	for k := range metaObj.Methods {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, i := range keys {
		k := uint32(i)
		m := metaObj.Methods[k]

		// generate uniq name for the method
		methodName := registerName(strings.Title(m.Name), methodNames)
		if err := generateMethod(file, set, serviceName, m, methodName); err != nil {
			return fmt.Errorf("failed to render method %s of %s: %s", m.Name, serviceName, err)
		}
	}
	keys = make([]int, 0)
	for k := range metaObj.Signals {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, i := range keys {
		k := uint32(i)
		s := metaObj.Signals[k]
		methodName := registerName("Signal"+strings.Title(s.Name), methodNames)

		if err := generateSignal(file, set, serviceName, s, methodName); err != nil {
			return fmt.Errorf("failed to render signal %s of %s: %s", s.Name, serviceName, err)
		}
	}
	return nil
}

func GenerateInterfaces(metaObjList []object.MetaObject, packageName string, w io.Writer) error {
	file, set := newFileAndSet(packageName)

	for _, metaObj := range metaObjList {
		generateObjectInterface(metaObj, metaObj.Description, set, file)
	}

	if err := file.Render(w); err != nil {
		return fmt.Errorf("failed to render %s: %s", packageName, err)
	}
	return nil
}

func GenerateProxy(metaObj object.MetaObject, packageName, serviceName string, w io.Writer) error {
	file, set := newFileAndSet(packageName)

	generateProxyObject(metaObj, serviceName, set, file)
	set.Declare(file)

	if err := file.Render(w); err != nil {
		return fmt.Errorf("failed to render %s: %s", serviceName, err)
	}
	return nil
}

func GenerateProxys(metaObjList []object.MetaObject, packageName string, w io.Writer) error {
	file, set := newFileAndSet(packageName)

	for _, metaObj := range metaObjList {
		generateProxyObject(metaObj, metaObj.Description, set, file)
	}
	set.Declare(file)

	if err := file.Render(w); err != nil {
		return fmt.Errorf("failed to render %s: %s", packageName, err)
	}
	return nil
}

func methodBodyBlock(m object.MetaMethod, params *signature.TupleValue, ret signature.ValueConstructor) (*Statement, error) {
	i := 0
	length := 20 + len(params.Members())
	writing := make([]jen.Code, length)
	writing[i] = jen.Var().Err().Error()
	i++
	if _, ok := ret.(signature.VoidValue); !ok {
		writing[i] = jen.Var().Id("ret").Add(ret.TypeName())
		i++
	}
	writing[i] = jen.Var().Id("buf").Op("*").Qual("bytes", "Buffer")
	i++
	writing[i] = jen.Id("buf = bytes.NewBuffer(make([]byte, 0))")
	i++
	for _, v := range params.Members() {
		if _, ok := ret.(signature.VoidValue); !ok {
			writing[i] = jen.If(jen.Err().Op("=").Add(v.Value.Marshal(v.Name, "buf")).Op(";").Err().Op("!=").Nil()).Block(
				jen.Id(`return ret, fmt.Errorf("failed to serialize ` + v.Name + `: %s", err)`),
			)
			i++
		} else {
			writing[i] = jen.If(jen.Err().Op("=").Add(v.Value.Marshal(v.Name, "buf")).Op(";").Err().Op("!=").Nil()).Block(
				jen.Id(`return fmt.Errorf("failed to serialize ` + v.Name + `: %s", err)`),
			)
			i++
		}
	}
	if _, ok := ret.(signature.VoidValue); !ok {
		writing[i] = jen.Id(fmt.Sprintf(`response, err := p.Call("%s", buf.Bytes())`, m.Name))
		i++
	} else {
		writing[i] = jen.Id(fmt.Sprintf(`_, err = p.Call("%s", buf.Bytes())`, m.Name))
		i++
	}
	if _, ok := ret.(signature.VoidValue); !ok {
		writing[i] = jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id(fmt.Sprintf(`return ret, fmt.Errorf("call %s failed: %s", err)`, m.Name, "%s")),
		)
		i++
	} else {
		writing[i] = jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id(fmt.Sprintf(`return fmt.Errorf("call %s failed: %s", err)`, m.Name, "%s")),
		)
		i++
	}
	if _, ok := ret.(signature.VoidValue); !ok {
		writing[i] = jen.Id("buf = bytes.NewBuffer(response)")
		i++
		writing[i] = jen.Id("ret, err =").Add(ret.Unmarshal("buf"))
		i++
		writing[i] = jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id(fmt.Sprintf(`return ret, fmt.Errorf("failed to parse %s response: %s", err)`, m.Name, "%s")),
		)
		i++
		writing[i] = jen.Return(jen.Id("ret"), jen.Nil())
		i++
	} else {
		writing[i] = jen.Return(jen.Nil())
		i++
	}

	return jen.Block(
		writing...,
	), nil
}

func registerName(name string, names map[string]bool) string {
	newName := name
	for i := 0; i < 100; i++ {
		if _, ok := names[newName]; !ok {
			break
		}
		newName = fmt.Sprintf("%s_%d", name, i)
	}
	names[newName] = true
	return newName
}

func generateMethod(file *jen.File, set *signature.TypeSet, serviceName string, m object.MetaMethod, methodName string) error {

	paramType, err := signature.Parse(m.ParametersSignature)
	if err != nil {
		return fmt.Errorf("failed to parse parameters %s: %s", m.ParametersSignature, err)
	}
	paramType.RegisterTo(set)
	tuple, ok := paramType.(*signature.TupleValue)
	if !ok {
		tuple = signature.NewTupleValue([]signature.ValueConstructor{paramType})
	}

	tuple.ConvertMetaObjects()

	ret, err := signature.Parse(m.ReturnSignature)
	if err != nil {
		return fmt.Errorf("failed to parse return: %s: %s", m.ReturnSignature, err)
	}

	// implementation detail: use object.MetaObject when generating
	// proxy in order for proxy to implement the object.Object
	// interface.
	if ret.Signature() == signature.MetaObjectSignature {
		ret = signature.NewMetaObjectValue()
	}
	ret.RegisterTo(set)

	body, err := methodBodyBlock(m, tuple, ret)
	if err != nil {
		return fmt.Errorf("failed to generate body: %s", err)
	}

	retType := jen.Params(ret.TypeName(), jen.Error())
	if _, ok := ret.(signature.VoidValue); ok {
		retType = jen.Error()
	}

	file.Func().Params(jen.Id("p").Op("*").Id(serviceName)).Id(strings.Title(methodName)).Add(
		tuple.Params(),
	).Add(
		retType,
	).Add(
		body,
	)
	return nil
}

func generateSignal(file *jen.File, set *signature.TypeSet, serviceName string, s object.MetaSignal, methodName string) error {

	signalType, err := signature.Parse(s.Signature)
	if err != nil {
		return fmt.Errorf("failed to parse signal %s: %s", s.Signature, err)
	}

	signalType.RegisterTo(set)

	retType := jen.Params(jen.Chan().Add(signalType.TypeName()), jen.Error())
	body := jen.Block(
		jen.Id(`signalID, err := p.SignalUid("`+s.Name+`")
		if err != nil {
			return nil, fmt.Errorf("signal %s not available: %s", "`+s.Name+`", err)
		}

		id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
		if err != nil {
			return nil, fmt.Errorf("failed to register event for %s: %s", "`+s.Name+`", err)
		}`),
		jen.Id("ch").Op(":=").Make(jen.Chan().Add(signalType.TypeName())),
		jen.List(jen.Id("chPay"), jen.Err()).Op(":=").Id("p.SignalStreamID").Call(jen.Id("signalID"), jen.Id("cancel")),
		jen.Id(`if err != nil {
			return nil, fmt.Errorf("failed to request signal: %s", err)
		}`),
		jen.Go().Func().Params().Block(
			jen.For().Block(
				jen.List(jen.Id("payload"), jen.Id("ok")).Op(":=").Op("<-").Id("chPay"),
				jen.Id(`if !ok {
					close(ch) // upstream is closed.
					err = p.UnregisterEvent(p.ObjectID(), signalID, id)
					if err != nil {
						// FIXME: implement proper logging.
						fmt.Printf("failed to unregister event %s: %s", "`+s.Name+`", err)
					}
					return
				}`),
				jen.Id("buf := bytes.NewBuffer(payload)"),
				jen.Id("_ = buf // discard unused variable error"),
				jen.List(jen.Id("e"), jen.Err()).Op(":=").Add(signalType.Unmarshal("buf")),
				jen.Id(`if err != nil {
					fmt.Errorf("failed to unmarshall tuple: %s", err)
					continue
				}
				ch<- e`),
			),
		).Call(),
		jen.Return(jen.Id("ch"), jen.Nil()),
	)

	file.Func().Params(jen.Id("p").Op("*").Id(serviceName)).Id(methodName).Params(
		jen.Id("cancel").Chan().Int(),
	).Add(
		retType,
	).Add(
		body,
	)
	return nil
}

func generateSignalDef(file *jen.File, set *signature.TypeSet, serviceName string, s object.MetaSignal, methodName string) (jen.Code, error) {
	signalType, err := signature.Parse(s.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to parse signal %s: %s", s.Signature, err)
	}
	retType := jen.Params(jen.Chan().Add(signalType.TypeName()), jen.Error())
	return jen.Id(methodName).Params(jen.Id("cancel").Chan().Int()).Add(retType), nil
}
