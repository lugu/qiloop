package idl

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
)

// Statement is imported from jennifer code generator.
type Statement = jen.Statement

func objName(name string) string {
	return name + "Proxy"
}

func proxyName(name string) string {
	return "proxy" + name
}

func generateInterface(itf *InterfaceType, file *jen.File) error {
	serviceName := signature.CleanName(itf.Name)
	err := generateObjectInterface(itf, serviceName, file)
	if err != nil {
		return fmt.Errorf("declare interface %s: %s",
			itf.Name, err)
	}

	err = generateProxyObject(itf, serviceName, file)
	if err != nil {
		return fmt.Errorf("declare proxy %s: %s",
			itf.Name, err)
	}
	return nil
}

func skipActionInterface(serviceName string, action uint32) bool {
	if serviceName == "ServiceZero" {
		return false
	}
	if serviceName == "Object" && action > uint32(10) {
		return false
	}
	if action >= object.MinUserActionID {
		return false
	}
	return true
}

func generateObjectInterface(itf *InterfaceType, serviceName string,
	file *jen.File) error {

	definitions := make([]jen.Code, 0)

	method := func(m object.MetaMethod, methodName string) error {
		method := itf.Methods[m.Uid]
		methodName = signature.CleanMethodName(methodName)
		if serviceName == "Object" {
			methodName = signature.CleanName(methodName)
		}
		if skipActionInterface(serviceName, m.Uid) {
			return nil
		}
		err := generateMethodDef(file, serviceName, method,
			methodName, &definitions)
		if err != nil {
			return fmt.Errorf(
				"render method definition %s of %s: %s",
				m.Name, serviceName, err)
		}
		return nil
	}
	signal := func(s object.MetaSignal, signalName string) error {
		signal := itf.Signals[s.Uid]
		signalName = signature.CleanName("Subscribe" + signalName)
		if skipActionInterface(serviceName, s.Uid) {
			return nil
		}
		err := generateSignalDef(file, serviceName, signal,
			signalName, &definitions)
		if err != nil {
			return fmt.Errorf(
				"render signal %s of %s: %s",
				s.Name, serviceName, err)
		}
		return nil
	}
	property := func(p object.MetaProperty, propertyName string) error {
		property := itf.Properties[p.Uid]
		getMethodName := "Get" + propertyName
		setMethodName := "Set" + propertyName
		subscribeMethodName := "Subscribe" + propertyName
		if skipActionInterface(serviceName, p.Uid) {
			return nil
		}
		err := generatePropertyDef(file, serviceName, property,
			getMethodName, setMethodName, subscribeMethodName,
			&definitions)
		if err != nil {
			return fmt.Errorf(
				"render property %s of %s: %s",
				p.Name, serviceName, err)
		}
		return nil
	}

	metaObj := itf.MetaObject()
	if err := metaObj.ForEachMethodAndSignal(method, signal, property); err != nil {
		return fmt.Errorf("generate interface object %s: %s",
			serviceName, err)
	}

	if serviceName == "Object" {
		definitions = append(definitions,
			jen.Qual("github.com/lugu/qiloop/type/object", "Object"))
	} else if serviceName != "ServiceZero" {
		comment := jen.Comment(objName("Generic methods shared by all objects"))
		definitions = append(definitions, comment)
		definitions = append(definitions,
			jen.Qual("github.com/lugu/qiloop/bus", objName("Object")))
	}

	if serviceName != "ServiceZero" && serviceName != "Object" {
		comment := jen.Comment("WithContext can be used cancellation and timeout")
		definitions = append(definitions, comment)

		def := jen.Id("WithContext").Params(
			jen.Id("ctx").Qual("context", "Context"),
		).Params(jen.Id(objName(serviceName)))
		definitions = append(definitions, def)
	} else {
		comment := jen.Comment("Proxy returns a proxy.")
		definitions = append(definitions, comment)

		def := jen.Id("Proxy").Params().Params(
			jen.Qual("github.com/lugu/qiloop/bus", "Proxy"),
		)
		definitions = append(definitions, def)
	}

	file.Comment(objName(serviceName) + " represents a proxy object to the service")
	file.Type().Id(objName(serviceName)).Interface(
		definitions...,
	)

	return nil
}

func generateSignalDef(file *jen.File, serviceName string,
	signal Signal, signalName string, definitions *[]jen.Code) error {

	signalType := signal.Type()
	retType := jen.Params(
		jen.Id("unsubscribe").Func().Params(),
		jen.Id("updates").Chan().Add(signalType.TypeName()),
		jen.Id("err").Error(),
	)
	def := jen.Id(signalName).Params().Add(retType)
	*definitions = append(*definitions, def)
	return nil
}

// generatePropertyDef declares the setter, getter and subscribe methods
func generatePropertyDef(file *jen.File, serviceName string,
	property Property, getMethodName, setMethodName,
	subscribeMethodName string, definitions *[]jen.Code) error {

	propertyType := property.Type()

	retType := jen.Params(propertyType.TypeName(), jen.Error())
	getMethod := jen.Id(getMethodName).Params().Add(retType)
	*definitions = append(*definitions, getMethod)

	paramType := jen.Params(propertyType.TypeName())
	setMethod := jen.Id(setMethodName).Add(paramType).Error()
	*definitions = append(*definitions, setMethod)

	retType = jen.Params(
		jen.Id("unsubscribe").Func().Params(),
		jen.Id("updates").Chan().Add(propertyType.TypeName()),
		jen.Id("err").Error(),
	)
	subscribeMethod := jen.Id(subscribeMethodName).Params().Add(retType)
	*definitions = append(*definitions, subscribeMethod)
	return nil
}

func generateMethodDef(file *jen.File, serviceName string,
	method Method, methodName string, definitions *[]jen.Code) error {

	paramType := method.Tuple()
	paramType.ConvertMetaObjects()

	returnType := method.Return

	// implementation detail: use object.MetaObject when generating
	// proxy in order for proxy to implement the object.Object
	// interface.
	if returnType.Signature() == signature.MetaObjectSignature {
		returnType = signature.NewMetaObjectType()
	}

	def := jen.Id(methodName).Add(
		paramType.Params(),
	).Params(returnType.TypeName(), jen.Error())

	if returnType.Signature() == "v" {
		def = jen.Id(methodName).Add(paramType.Params()).Error()
	}

	*definitions = append(*definitions, def)
	return nil
}

func generateProxyObject(itf *InterfaceType, serviceName string,
	file *jen.File) error {

	proxyName := proxyName(serviceName)
	generateProxyType(file, serviceName, proxyName, itf)

	method := func(m object.MetaMethod, methodName string) error {
		method := itf.Methods[m.Uid]
		if serviceName != "Object" && serviceName != "ServiceZero" &&
			m.Uid < object.MinUserActionID {
			return nil
		}
		return generateProxyMethod(file, proxyName, method, methodName)
	}
	signal := func(s object.MetaSignal, signalName string) error {
		signal := itf.Signals[s.Uid]
		signalName = signature.CleanName("Subscribe" + signalName)
		if serviceName != "Object" && serviceName != "ServiceZero" &&
			s.Uid < object.MinUserActionID {
			return nil
		}
		return generateProxySignal(file, proxyName, signal, signalName)
	}
	property := func(p object.MetaProperty, propertyName string) error {
		property := itf.Properties[p.Uid]
		getMethodName := "Get" + propertyName
		setMethodName := "Set" + propertyName
		subscribeMethodName := "Subscribe" + propertyName
		if serviceName != "Object" && serviceName != "ServiceZero" &&
			p.Uid < object.MinUserActionID {
			return nil
		}
		return generateProxyProperty(file, proxyName, property,
			getMethodName, setMethodName, subscribeMethodName)
	}

	metaObj := itf.MetaObject()
	if err := metaObj.ForEachMethodAndSignal(method, signal, property); err != nil {
		return fmt.Errorf("generate proxy object %s: %s",
			serviceName, err)
	}
	return nil
}

func generateProxyType(file *jen.File, serviceName, ProxyName string,
	itf *InterfaceType) {

	file.Comment(ProxyName + " implements " + objName(serviceName))
	if ProxyName == proxyName("Object") || ProxyName == proxyName("ServiceZero") {
		file.Type().Id(ProxyName).Struct(
			jen.Id("proxy").Qual("github.com/lugu/qiloop/bus", "Proxy"),
		)
	} else {
		file.Type().Id(ProxyName).Struct(
			jen.Qual(
				"github.com/lugu/qiloop/bus",
				"ObjectProxy",
			),
			jen.Id("session").Qual(
				"github.com/lugu/qiloop/bus",
				"Session",
			),
		)
	}
	if ProxyName != proxyName("Object") && ProxyName != proxyName("ServiceZero") {
		file.Comment("Make" + serviceName + " returns a specialized proxy.")
		file.Func().Id("Make"+serviceName).Params(
			jen.Id("sess").Qual("github.com/lugu/qiloop/bus", "Session"),
			jen.Id("proxy").Qual("github.com/lugu/qiloop/bus", "Proxy"),
		).Id(objName(serviceName)).Block(
			jen.Return(
				jen.Op("&").Id(ProxyName).Values(
					jen.Qual("github.com/lugu/qiloop/bus", "MakeObject").Call(jen.Id("proxy")),
					jen.Id("sess"),
				),
			),
		)
	}
	blockContructor := jen.Id(
		`return Make` + serviceName + `(session, proxy), nil`,
	)
	if ProxyName == proxyName("Object") || ProxyName == proxyName("ServiceZero") {
		blockContructor = jen.Id(`return &` + ProxyName + `{ proxy }, nil`)
	}
	file.Comment(serviceName + " returns a proxy to a remote service")
	file.Func().Id(
		signature.CleanName(serviceName),
	).Params(
		jen.Id("session").Qual("github.com/lugu/qiloop/bus", "Session"),
	).Params(
		jen.Id(objName(serviceName)), jen.Error(),
	).Block(
		jen.List(jen.Id("proxy"), jen.Err()).Op(":=").Id("session").Dot("Proxy").Call(
			jen.Lit(serviceName),
			jen.Lit(1),
		),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("contact service: %s"),
					jen.Err(),
				),
			),
		),
		blockContructor,
	)

	if ProxyName != proxyName("Object") && ProxyName != proxyName("ServiceZero") {
		file.Comment("WithContext bound future calls to the context deadline and cancellation")
		file.Func().Params(
			jen.Id("p").Op("*").Id(proxyName(serviceName)),
		).Id("WithContext").Params(
			jen.Id("ctx").Qual("context", "Context"),
		).Params(jen.Id(objName(serviceName))).Block(
			jen.Id(`return Make` + serviceName + `(p.session, p.Proxy().WithContext(ctx))`),
		)
	} else {
		file.Comment("Proxy returns a proxy.")
		file.Func().Params(
			jen.Id("p").Op("*").Id(proxyName(serviceName)),
		).Id("Proxy").Params().Params(
			jen.Qual("github.com/lugu/qiloop/bus", "Proxy"),
		).Block(
			jen.Id(`return p.proxy`),
		)
	}
}

func generateProxyMethod(file *jen.File, serviceName string,
	method Method, methodName string) error {

	paramType := method.Tuple()
	paramType.ConvertMetaObjects()

	returnType := method.Return

	// use object.MetaObject when generating proxy in order for
	// proxy to implement the object.Object interface.
	if returnType.Signature() == signature.MetaObjectSignature {
		returnType = signature.NewMetaObjectType()
	}

	body, err := methodBodyBlock2(method, paramType, returnType)
	if err != nil {
		return fmt.Errorf("generate body: %s", err)
	}

	returnCode := jen.Params(returnType.TypeName(), jen.Error())
	if returnType.Signature() == "v" {
		returnCode = jen.Error()
	}

	goMethodName := signature.CleanMethodName(methodName)
	if serviceName == proxyName("Object") {
		goMethodName = signature.CleanName(methodName)
	}
	file.Comment(goMethodName + " calls the remote procedure")
	file.Func().Params(jen.Id("p").Op("*").Id(serviceName)).Id(goMethodName).Add(
		paramType.Params(),
	).Add(
		returnCode,
	).Add(
		body,
	)
	return nil
}

func generateSubscribe(file *jen.File, serviceName, actionName, methodName string,
	actionType signature.Type, isSignal bool) {

	retType := jen.Params(
		jen.Func().Params(),
		jen.Chan().Add(actionType.TypeName()),
		jen.Error(),
	)
	methodID := "PropertyID"
	if isSignal {
		methodID = "SignalID"
	}
	body := jen.Block(
		jen.Id(fmt.Sprintf(`signalID, err := p.Proxy().MetaObject().%s("%s", "%s")`,
			methodID, actionName, actionType.Signature())),
		jen.Id(`if err != nil {
			return nil, nil, fmt.Errorf("%s not available: %s", "`+actionName+`", err)
		}`),

		jen.Id("ch").Op(":=").Make(jen.Chan().Add(actionType.TypeName())),
		jen.List(
			jen.Id("cancel"),
			jen.Id("chPay"),
			jen.Err(),
		).Op(":=").Id("p.Proxy().SubscribeID").Call(jen.Id("signalID")),
		jen.Id(`if err != nil {
			return nil, nil, fmt.Errorf("request property: %s", err)
		}`),
		jen.Go().Func().Params().Block(jen.For().Block(
			jen.List(jen.Id("payload"), jen.Id("ok")).Op(":=").Op("<-").Id("chPay"),
			jen.Id(`if !ok {
					// connection lost or cancellation.
					close(ch)
					return
				}`),
			jen.Id("buf").Op(":=").Qual("bytes", "NewBuffer").
				Call(jen.Id("payload")),
			jen.Id("_ = buf // discard unused variable error"),
			jen.List(jen.Id("e"),
				jen.Err()).Op(":=").Add(actionType.Unmarshal("buf")),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Qual("log", "Printf").Call(
					jen.Lit("unmarshall tuple: %s"),
					jen.Id("err"),
				),
				jen.Continue(),
			),
			jen.Id(`ch<- e`),
		)).Call(),
		jen.Id(` return cancel, ch, nil`),
	)

	file.Comment(methodName + " subscribe to a remote property")
	file.Func().Params(
		jen.Id("p").Op("*").Id(serviceName),
	).Id(methodName).Params().Add(retType).Add(body)
}

func generateProxySignal(file *jen.File, serviceName string,
	signal Signal, methodName string) error {

	signalType := signal.Type()

	generateSubscribe(file, serviceName, signal.Name, methodName, signalType, true)

	return nil
}
func generateProxyProperty(file *jen.File, serviceName string,
	property Property, getMethodName, setMethodName, subscribeMethodName string) error {

	propertyType := property.Type()

	generatePropertyGet(file, serviceName, property, getMethodName)
	generatePropertySet(file, serviceName, property, setMethodName)
	generateSubscribe(file, serviceName, property.Name, subscribeMethodName, propertyType, false)

	return nil
}
func generatePropertyGet(file *jen.File, serviceName string,
	property Property, methodName string) error {

	retType := jen.Params(
		jen.Id("ret").Add(property.Type().TypeName()),
		jen.Err().Error(),
	)
	body := jen.Block(
		jen.Id("name").Op(":=").Qual(
			"github.com/lugu/qiloop/type/value",
			"String",
		).Params(jen.Lit(property.Name)),
		jen.Id(`value, err := p.Property(name)`),
		jen.Id(`if err != nil {
		    return ret, fmt.Errorf("get property: %s", err)
		}`),
		jen.Var().Id("buf").Qual("bytes", "Buffer"),
		jen.Id(`err = value.Write(&buf)`),
		jen.Id(`if err != nil {
		    return ret, fmt.Errorf("read response: %s", err)
		}`),
		jen.Id(`s, err := basic.ReadString(&buf)`),
		jen.Id(`if err != nil {
		    return ret, fmt.Errorf("read signature: %s", err)
		}`),
		jen.Id(`// check the signature`),
		jen.Id(`sig := "`+property.Type().Signature()+`"`),
		jen.Id(`if sig != s {
		    return ret, fmt.Errorf("unexpected signature: %s instead of %s",
			s, sig)
		}`),
		jen.Id("ret, err =").Add(property.Type().Unmarshal("&buf")),
		jen.Id("return ret, err"),
	)

	file.Comment(methodName + " updates the property value")
	file.Func().Params(
		jen.Id("p").Op("*").Id(serviceName),
	).Id(methodName).Params().Add(retType).Add(body)
	return nil
}

func generatePropertySet(file *jen.File, serviceName string,
	property Property, methodName string) error {

	propertyType := property.Type()
	paramType := jen.Params(jen.Id("update").Add(propertyType.TypeName()))
	// 1. serialize property type
	// 2. make it a value using its signature
	// 3. name the property with a value
	// 4. call setProperty
	body := jen.Block(
		jen.Id("name").Op(":=").Qual(
			"github.com/lugu/qiloop/type/value",
			"String",
		).Params(jen.Lit(property.Name)),
		jen.Var().Id("buf").Qual("bytes", "Buffer"),
		jen.Err().Op(":=").Add(propertyType.Marshal("update", "&buf")),
		jen.Id(`if err != nil {
		    return fmt.Errorf("marshall error: %s", err)
		}`),
		jen.Id("val").Op(":=").Qual(
			"github.com/lugu/qiloop/type/value",
			"Opaque",
		).Params(
			jen.Lit(propertyType.Signature()),
			jen.Id(`buf.Bytes()`),
		),
		jen.Id(`return p.SetProperty(name, val)`),
	)

	file.Comment(methodName + " updates the property value")
	file.Func().Params(
		jen.Id("p").Op("*").Id(serviceName),
	).Id(methodName).Add(paramType).Error().Add(body)
	return nil
}

// TODO: delete me
func methodBodyBlock(method Method, params *signature.TupleType,
	ret signature.Type) (*Statement, error) {
	writing := make([]jen.Code, 0)
	writing = append(writing, jen.Var().Err().Error())
	if ret.Signature() != "v" {
		writing = append(writing, jen.Var().Id("ret").Add(ret.TypeName()))
	}
	writing = append(writing, jen.Var().Id("buf").Qual("bytes", "Buffer"))
	for _, v := range params.Members {
		if ret.Signature() != "v" {
			writing = append(writing, jen.If(jen.Err().Op("=").Add(v.Type.Marshal(v.Name, "&buf")).Op(";").Err().Op("!=").Nil()).Block(
				jen.Id(`return ret, fmt.Errorf("serialize `+v.Name+`: %s", err)`),
			))
		} else {
			writing = append(writing, jen.If(jen.Err().Op("=").Add(v.Type.Marshal(v.Name, "&buf")).Op(";").Err().Op("!=").Nil()).Block(
				jen.Id(`return fmt.Errorf("serialize `+v.Name+`: %s", err)`),
			))
		}
	}
	writing = append(writing, jen.Id(fmt.Sprintf(`methodID, _, err := p.Proxy().MetaObject().MethodID("%s", "%s")`, method.Name, params.Signature())))
	if ret.Signature() != "v" {
		writing = append(writing, jen.Id(`if err != nil {
		return ret, err
		}`))
		writing = append(writing, jen.Id(`response, err := p.Proxy().CallID(methodID, buf.Bytes())`))
		writing = append(writing, jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id(fmt.Sprintf(`return ret, fmt.Errorf("call %s failed: %s", err)`, method.Name, "%s")),
		))
		writing = append(writing, jen.Id("resp := bytes.NewBuffer(response)"))
		writing = append(writing, jen.Id("ret, err =").Add(ret.Unmarshal("resp")))
		writing = append(writing, jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id(fmt.Sprintf(`return ret, fmt.Errorf("parse %s response: %s", err)`, method.Name, "%s")),
		))
		writing = append(writing, jen.Return(jen.Id("ret"), jen.Nil()))
	} else {
		writing = append(writing, jen.Id(`if err != nil {
		return err
		}`))
		writing = append(writing, jen.Id(`_, err = p.Proxy().CallID(methodID, buf.Bytes())`))
		writing = append(writing, jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id(fmt.Sprintf(`return fmt.Errorf("call %s failed: %s", err)`, method.Name, "%s")),
		))
		writing = append(writing, jen.Return(jen.Nil()))
	}

	return jen.Block(
		writing...,
	), nil
}

func methodBodyBlock2(method Method, params *signature.TupleType,
	ret signature.Type) (*Statement, error) {
	writing := []jen.Code{}
	if ret.Signature() == "v" {
		writing = append(writing, jen.Var().Id("ret").Struct())
	} else {
		writing = append(writing, jen.Var().Id("ret").Add(ret.TypeName()))
	}
	args := []jen.Code{
		jen.Lit(params.Signature()),
	}
	for _, v := range params.Members {
		if v.Type.Signature() == "o" {
			args = append(args,
				jen.Qual("github.com/lugu/qiloop/bus",
					"ObjectReference").Call(
					jen.Id(v.Name).Op(".").Id("Proxy").Call(),
				))
		} else {
			args = append(args, jen.Id(v.Name))
		}
	}
	writing = append(writing, jen.Id("args").Op(":=").Qual(
		"github.com/lugu/qiloop/bus", "NewParams",
	).Call(args...))

	if ret.Signature() == "o" && !signature.TypeIsObjectReference(ret) {
		writing = append(writing, jen.Var().Id("retRef").Qual(
			"github.com/lugu/qiloop/type/object", "ObjectReference",
		))
		writing = append(writing, jen.Id("resp").Op(":=").Qual(
			"github.com/lugu/qiloop/bus", "NewResponse",
		).Call(jen.Lit(ret.Signature()), jen.Op("&").Id("retRef")))
	} else {
		writing = append(writing, jen.Id("resp").Op(":=").Qual(
			"github.com/lugu/qiloop/bus", "NewResponse",
		).Call(jen.Lit(ret.Signature()), jen.Op("&").Id("ret")))
	}
	writing = append(writing, jen.Id("err").Op(":=").Id("p").Op(".").
		Id("Proxy").Call().Op(".").Id("Call2").Call(
		jen.Lit(method.Name), jen.Id("args"), jen.Id("resp"),
	),
	)
	if ret.Signature() == "v" {
		writing = append(writing, jen.Id(`if err != nil {
		return fmt.Errorf("call `+method.Name+` failed: %s", err)
		}
		return nil`))
	} else {
		writing = append(writing, jen.Id(
			`if err != nil {
			return ret, fmt.Errorf("call `+method.Name+` failed: %s", err)
		}`))
		if ret.Signature() == "o" &&
			!signature.TypeIsObjectReference(ret) {

			var name string
			switch v := ret.(type) {
			case *RefType:
				name = v.Name
			case *InterfaceType:
				name = v.Name
			}
			writing = append(writing, jen.Id(
				`proxy, err := p.session.Object(retRef)
			if err != nil {
				return nil, fmt.Errorf("proxy: %s", err)
			}
			ret = Make`+name+`(p.session, proxy)`))
		}
		writing = append(writing, jen.Id(`return ret, nil`))
	}
	return jen.Block(
		writing...,
	), nil
}
