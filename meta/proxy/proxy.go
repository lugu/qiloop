package proxy

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	"io"
	"strings"
)

type Statement = jen.Statement

func newFileAndSet(packageName string) (*jen.File, *signature.TypeSet) {
	file := jen.NewFile(packageName)
	file.PackageComment("file generated. DO NOT EDIT.")
	return file, signature.NewTypeSet()
}

func generateProxyType(file *jen.File, serviceName, proxyName string, metaObj object.MetaObject) {

	if proxyName == "ObjectProxy" || proxyName == "ServerProxy" {
		file.Type().Id(proxyName).Struct(jen.Qual("github.com/lugu/qiloop/bus", "Proxy"))
	} else {
		file.Type().Id(proxyName).Struct(jen.Qual(
			"github.com/lugu/qiloop/bus/client/object",
			"ObjectProxy",
		))
	}
	blockContructor := jen.Return(

		jen.Op("&").Id(proxyName).Values(jen.Qual(
			"github.com/lugu/qiloop/bus/client/object",
			"ObjectProxy",
		).Values(jen.Id("proxy"))),
		jen.Nil(),
	)
	if proxyName == "ObjectProxy" || proxyName == "ServerProxy" {
		blockContructor = jen.Id(`return &` + proxyName + `{ proxy }, nil`)
	}
	file.Func().Id("New"+serviceName).Params(
		jen.Id("ses").Qual("github.com/lugu/qiloop/bus", "Session"),
		jen.Id("obj").Uint32(),
	).Params(
		jen.Id(serviceName),
		jen.Error(),
	).Block(
		jen.List(jen.Id("proxy"), jen.Err()).Op(":=").Id("ses").Dot("Proxy").Call(
			jen.Lit(serviceName),
			jen.Id("obj"),
		),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("failed to contact service: %s"),
					jen.Err(),
				),
			),
		),
		blockContructor,
	)
	file.Func().Params(
		jen.Id("s").Id("NewServices"),
	).Id(
		strings.Title(serviceName),
	).Params().Params(
		jen.Id(serviceName), jen.Error(),
	).Block(
		jen.Id(`return New` + serviceName + `(s.session, 1)`),
	)
}

func generateMethodDef(file *jen.File, set *signature.TypeSet, serviceName string, m object.MetaMethod, methodName string) (jen.Code, error) {

	paramType, err := signature.Parse(m.ParametersSignature)
	if err != nil {
		return nil, fmt.Errorf("failed to parse parameters %s: %s", m.ParametersSignature, err)
	}
	paramType.RegisterTo(set)
	tuple, ok := paramType.(*signature.TupleType)
	if !ok {
		tuple = signature.NewTupleType([]signature.Type{paramType})
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
		ret = signature.NewMetaObjectType()
	}
	ret.RegisterTo(set)
	if ret.Signature() == "v" {
		return jen.Id(methodName).Add(tuple.Params()).Error(), nil
	} else {
		return jen.Id(methodName).Add(tuple.Params()).Params(ret.TypeName(), jen.Error()), nil
	}
}

func generateObjectInterface(metaObj object.MetaObject, serviceName string, set *signature.TypeSet, file *jen.File) error {

	definitions := make([]jen.Code, 0)
	if serviceName != "Server" {
		definitions = append(definitions, jen.Qual("github.com/lugu/qiloop/type/object", "Object"))
	}
	definitions = append(definitions, jen.Qual("github.com/lugu/qiloop/bus", "Proxy"))

	methodCall := func(m object.MetaMethod, methodName string) error {
		if serviceName != "Server" && m.Uid < object.MinUserActionID {
			return nil
		}
		def, err := generateMethodDef(file, set, serviceName, m, methodName)
		if err != nil {
			return fmt.Errorf("failed to render method definition %s of %s: %s", m.Name, serviceName, err)
		}
		definitions = append(definitions, def)
		return nil
	}
	signalCall := func(s object.MetaSignal, methodName string) error {
		if serviceName != "Server" && s.Uid < object.MinUserActionID {
			return nil
		}
		def, err := generateSignalDef(file, set, serviceName, s, methodName)
		if err != nil {
			return fmt.Errorf("failed to render signal %s of %s: %s", s.Name, serviceName, err)
		}
		definitions = append(definitions, def)
		return nil
	}

	if err := metaObj.ForEachMethodAndSignal(methodCall, signalCall); err != nil {
		return fmt.Errorf("failed to generate interface object %s: %s", serviceName, err)
	}

	file.Type().Id(serviceName).Interface(
		definitions...,
	)
	return nil
}

func generateProxyObject(metaObj object.MetaObject, serviceName string, set *signature.TypeSet, file *jen.File) error {

	proxyName := serviceName + "Proxy"
	generateProxyType(file, serviceName, proxyName, metaObj)

	methodCall := func(m object.MetaMethod, methodName string) error {
		if serviceName != "Object" && serviceName != "Server" &&
			m.Uid < object.MinUserActionID {
			return nil
		}
		return generateMethod(file, set, proxyName, m, methodName)
	}
	signalCall := func(s object.MetaSignal, methodName string) error {
		if serviceName != "Object" && serviceName != "Server" &&
			s.Uid < object.MinUserActionID {
			return nil
		}
		return generateSignal(file, set, proxyName, s, methodName)
	}

	if err := metaObj.ForEachMethodAndSignal(methodCall, signalCall); err != nil {
		return fmt.Errorf("failed to generate proxy object %s: %s", serviceName, err)
	}
	return nil
}

func generateNewServices(file *jen.File) {
	file.Type().Id(
		"NewServices",
	).Struct(
		jen.Id("session").Qual("github.com/lugu/qiloop/bus", "Session"),
	)
	file.Func().Id(
		"Services",
	).Params(
		jen.Id("s").Qual("github.com/lugu/qiloop/bus", "Session"),
	).Id(
		"NewServices",
	).Block(
		jen.Id(`return NewServices{ session: s, }`),
	)
}

func cleanMetaObject(meta *object.MetaObject) {
	meta.Description = util.CleanName(meta.Description)
	for i, m := range meta.Methods {
		m.Name = util.CleanName(m.Name)
		meta.Methods[i] = m
	}
	for i, s := range meta.Signals {
		s.Name = util.CleanName(s.Name)
		meta.Signals[i] = s
	}
	for i, p := range meta.Properties {
		p.Name = util.CleanName(p.Name)
		meta.Properties[i] = p
	}
}

func Generate(metaObjList []object.MetaObject, packageName string, w io.Writer) error {
	file, set := newFileAndSet(packageName)
	generateNewServices(file)

	for i := range metaObjList {
		cleanMetaObject(&metaObjList[i])
	}
	for _, metaObj := range metaObjList {
		generateObjectInterface(metaObj, metaObj.Description, set, file)
	}
	for i, metaObj := range metaObjList {
		err := generateProxyObject(metaObj, metaObj.Description, set, file)
		if err != nil {
			return fmt.Errorf("failed to render %s (%d): %s",
				metaObj.Description, i, err)
		}
	}
	set.Declare(file)

	if err := file.Render(w); err != nil {
		return fmt.Errorf("failed to render %s: %s", packageName, err)
	}
	return nil
}

func methodBodyBlock(m object.MetaMethod, params *signature.TupleType, ret signature.Type) (*Statement, error) {
	writing := make([]jen.Code, 0)
	writing = append(writing, jen.Var().Err().Error())
	if ret.Signature() != "v" {
		writing = append(writing, jen.Var().Id("ret").Add(ret.TypeName()))
	}
	writing = append(writing, jen.Var().Id("buf").Op("*").Qual("bytes", "Buffer"))
	writing = append(writing, jen.Id("buf = bytes.NewBuffer(make([]byte, 0))"))
	for _, v := range params.Members {
		if ret.Signature() != "v" {
			writing = append(writing, jen.If(jen.Err().Op("=").Add(v.Type.Marshal(v.Name, "buf")).Op(";").Err().Op("!=").Nil()).Block(
				jen.Id(`return ret, fmt.Errorf("failed to serialize `+v.Name+`: %s", err)`),
			))
		} else {
			writing = append(writing, jen.If(jen.Err().Op("=").Add(v.Type.Marshal(v.Name, "buf")).Op(";").Err().Op("!=").Nil()).Block(
				jen.Id(`return fmt.Errorf("failed to serialize `+v.Name+`: %s", err)`),
			))
		}
	}
	if ret.Signature() != "v" {
		writing = append(writing, jen.Id(fmt.Sprintf(`response, err := p.Call("%s", buf.Bytes())`, m.Name)))
	} else {
		writing = append(writing, jen.Id(fmt.Sprintf(`_, err = p.Call("%s", buf.Bytes())`, m.Name)))
	}
	if ret.Signature() != "v" {
		writing = append(writing, jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id(fmt.Sprintf(`return ret, fmt.Errorf("call %s failed: %s", err)`, m.Name, "%s")),
		))
	} else {
		writing = append(writing, jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id(fmt.Sprintf(`return fmt.Errorf("call %s failed: %s", err)`, m.Name, "%s")),
		))
	}
	if ret.Signature() != "v" {
		writing = append(writing, jen.Id("buf = bytes.NewBuffer(response)"))
		writing = append(writing, jen.Id("ret, err =").Add(ret.Unmarshal("buf")))
		writing = append(writing, jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id(fmt.Sprintf(`return ret, fmt.Errorf("failed to parse %s response: %s", err)`, m.Name, "%s")),
		))
		writing = append(writing, jen.Return(jen.Id("ret"), jen.Nil()))
	} else {
		writing = append(writing, jen.Return(jen.Nil()))
	}

	return jen.Block(
		writing...,
	), nil
}

func generateMethod(file *jen.File, set *signature.TypeSet, serviceName string, m object.MetaMethod, methodName string) error {

	paramType, err := signature.Parse(m.ParametersSignature)
	if err != nil {
		return fmt.Errorf("failed to parse parameters %s: %s", m.ParametersSignature, err)
	}
	paramType.RegisterTo(set)
	tuple, ok := paramType.(*signature.TupleType)
	if !ok {
		tuple = signature.NewTupleType([]signature.Type{paramType})
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
		ret = signature.NewMetaObjectType()
	}
	ret.RegisterTo(set)

	body, err := methodBodyBlock(m, tuple, ret)
	if err != nil {
		return fmt.Errorf("failed to generate body: %s", err)
	}

	retType := jen.Params(ret.TypeName(), jen.Error())
	if ret.Signature() == "v" {
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

		_, err = p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
		if err != nil {
			return nil, fmt.Errorf("failed to register event for %s: %s", "`+s.Name+`", err)
		}`),
		jen.Id("ch").Op(":=").Make(jen.Chan().Add(signalType.TypeName())),
		jen.List(jen.Id("chPay"), jen.Err()).Op(":=").Id("p.SubscribeID").Call(jen.Id("signalID"), jen.Id("cancel")),
		jen.Id(`if err != nil {
			return nil, fmt.Errorf("failed to request signal: %s", err)
		}`),
		jen.Go().Func().Params().Block(
			jen.For().Block(
				jen.List(jen.Id("payload"), jen.Id("ok")).Op(":=").Op("<-").Id("chPay"),
				jen.Id(`if !ok {
					// connection lost.
					close(ch)
					return
				}`),
				jen.Id("buf := bytes.NewBuffer(payload)"),
				jen.Id("_ = buf // discard unused variable error"),
				jen.List(jen.Id("e"), jen.Err()).Op(":=").Add(signalType.Unmarshal("buf")),
				jen.If(jen.Id("err").Op("!=").Nil()).Block(
					jen.Qual("log", "Printf").Call(jen.Lit("failed to unmarshall tuple: %s"), jen.Id("err")),
					jen.Continue(),
				),
				jen.Id(`ch<- e`),
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
	signalType.RegisterTo(set)
	retType := jen.Params(jen.Chan().Add(signalType.TypeName()), jen.Error())
	return jen.Id(methodName).Params(jen.Id("cancel").Chan().Int()).Add(retType), nil
}
