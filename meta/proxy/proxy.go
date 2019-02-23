package proxy

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
	"io"
)

// Statement is imported from jennifer code generator.
type Statement = jen.Statement

// GeneratePackage generate the proxy for the package declaration.
func GeneratePackage(w io.Writer, pkg *idl.PackageDeclaration) error {

	if pkg.Name == "" {
		return fmt.Errorf("empty package name")
	}
	file := jen.NewFile(pkg.Name)
	msg := "Package " + pkg.Name + " contains a generated proxy"
	file.PackageComment(msg)
	file.PackageComment("File generated. DO NOT EDIT.")

	generateNewServices(file)

	set := signature.NewTypeSet()
	for _, typ := range pkg.Types {
		err := generateType(file, set, typ)
		if err != nil {
			return err
		}
	}

	set.Declare(file)
	return file.Render(w)
}

func generateType(file *jen.File, set *signature.TypeSet, typ signature.Type) error {
	itf, ok := typ.(*idl.InterfaceType)
	if ok {
		serviceName := util.CleanName(itf.Name)
		err := generateObjectInterface(itf, serviceName, set, file)
		if err != nil {
			return fmt.Errorf("failed to render interface %s: %s",
				itf.Name, err)
		}

		err = generateProxyObject(itf, serviceName, set, file)
		if err != nil {
			return fmt.Errorf("failed to render proxy %s: %s",
				itf.Name, err)
		}
	} else {
		typ.RegisterTo(set)
	}
	return nil
}

func generateNewServices(file *jen.File) {
	file.Comment("Constructor gives access to remote services")
	file.Type().Id(
		"Constructor",
	).Struct(
		jen.Id("session").Qual("github.com/lugu/qiloop/bus", "Session"),
	)
	file.Comment("Services gives access to the services constructor")
	file.Func().Id(
		"Services",
	).Params(
		jen.Id("s").Qual("github.com/lugu/qiloop/bus", "Session"),
	).Id(
		"Constructor",
	).Block(
		jen.Id(`return Constructor{ session: s, }`),
	)
}

func generateObjectInterface(itf *idl.InterfaceType, serviceName string,
	set *signature.TypeSet, file *jen.File) error {

	definitions := make([]jen.Code, 0)
	if serviceName != "Server" {
		definitions = append(definitions,
			jen.Qual("github.com/lugu/qiloop/type/object", "Object"))
	}
	definitions = append(definitions,
		jen.Qual("github.com/lugu/qiloop/bus", "Proxy"))

	method := func(m object.MetaMethod, methodName string) error {
		method := itf.Methods[m.Uid]
		methodName = util.CleanName(methodName)
		if serviceName != "Server" && m.Uid < object.MinUserActionID {
			return nil
		}
		err := generateMethodDef(file, set, serviceName, method,
			methodName, &definitions)
		if err != nil {
			return fmt.Errorf(
				"failed to render method definition %s of %s: %s",
				m.Name, serviceName, err)
		}
		return nil
	}
	signal := func(s object.MetaSignal, signalName string) error {
		signal := itf.Signals[s.Uid]
		signalName = util.CleanName("Subscribe" + signalName)
		if serviceName != "Server" && s.Uid < object.MinUserActionID {
			return nil
		}
		err := generateSignalDef(file, set, serviceName, signal,
			signalName, &definitions)
		if err != nil {
			return fmt.Errorf(
				"failed to render signal %s of %s: %s",
				s.Name, serviceName, err)
		}
		return nil
	}
	property := func(p object.MetaProperty, propertyName string) error {
		property := itf.Properties[p.Uid]
		getMethodName := "Get" + propertyName
		setMethodName := "Set" + propertyName
		subscribeMethodName := "Subscribe" + propertyName
		if serviceName != "Server" && p.Uid < object.MinUserActionID {
			return nil
		}
		err := generatePropertyDef(file, set, serviceName, property,
			getMethodName, setMethodName, subscribeMethodName,
			&definitions)
		if err != nil {
			return fmt.Errorf(
				"failed to render property %s of %s: %s",
				p.Name, serviceName, err)
		}
		return nil
	}

	metaObj := itf.MetaObject()
	if err := metaObj.ForEachMethodAndSignal(method, signal, property); err != nil {
		return fmt.Errorf("failed to generate interface object %s: %s",
			serviceName, err)
	}

	file.Comment(serviceName + " is a proxy object to the remote service")
	file.Type().Id(serviceName).Interface(
		definitions...,
	)
	return nil
}

func generateSignalDef(file *jen.File, set *signature.TypeSet, serviceName string,
	signal idl.Signal, signalName string, definitions *[]jen.Code) error {

	signalType := signal.Type()
	signalType.RegisterTo(set)
	retType := jen.Params(
		jen.Id("unsubscribe").Func().Params(),
		jen.Id("updates").Chan().Add(signalType.TypeName()),
		jen.Id("err").Error(),
	)
	comment := jen.Comment(signalName + " subscribe to a remote signal")
	*definitions = append(*definitions, comment)
	def := jen.Id(signalName).Params().Add(retType)
	*definitions = append(*definitions, def)
	return nil
}

// generatePropertyDef returns the set, get and subscribe methods
func generatePropertyDef(file *jen.File, set *signature.TypeSet, serviceName string,
	property idl.Property, getMethodName, setMethodName,
	subscribeMethodName string, definitions *[]jen.Code) error {

	propertyType := property.Type()
	propertyType.RegisterTo(set)

	retType := jen.Params(propertyType.TypeName(), jen.Error())
	getMethod := jen.Id(getMethodName).Params().Add(retType)
	comment := jen.Comment(getMethodName + " returns the property value")
	*definitions = append(*definitions, comment)
	*definitions = append(*definitions, getMethod)

	paramType := jen.Params(propertyType.TypeName())
	setMethod := jen.Id(setMethodName).Add(paramType).Error()
	comment = jen.Comment(setMethodName + " sets the property value")
	*definitions = append(*definitions, comment)
	*definitions = append(*definitions, setMethod)

	retType = jen.Params(
		jen.Id("unsubscribe").Func().Params(),
		jen.Id("updates").Chan().Add(propertyType.TypeName()),
		jen.Id("err").Error(),
	)
	subscribeMethod := jen.Id(subscribeMethodName).Params().Add(retType)
	comment = jen.Comment(subscribeMethodName + " regusters to a property")
	*definitions = append(*definitions, comment)
	*definitions = append(*definitions, subscribeMethod)
	return nil
}

func generateMethodDef(file *jen.File, set *signature.TypeSet, serviceName string,
	method idl.Method, methodName string, definitions *[]jen.Code) error {

	paramType := method.Tuple()
	paramType.RegisterTo(set)
	paramType.ConvertMetaObjects()

	returnType := method.Return

	// implementation detail: use object.MetaObject when generating
	// proxy in order for proxy to implement the object.Object
	// interface.
	if returnType.Signature() == signature.MetaObjectSignature {
		returnType = signature.NewMetaObjectType()
	}
	returnType.RegisterTo(set)

	comment := jen.Comment(methodName + " calls the remote procedure")
	*definitions = append(*definitions, comment)

	def := jen.Id(methodName).Add(
		paramType.Params(),
	).Params(returnType.TypeName(), jen.Error())

	if returnType.Signature() == "v" {
		def = jen.Id(methodName).Add(paramType.Params()).Error()
	}

	*definitions = append(*definitions, def)
	return nil
}

func generateProxyObject(itf *idl.InterfaceType, serviceName string,
	set *signature.TypeSet, file *jen.File) error {

	proxyName := serviceName + "Proxy"
	generateProxyType(file, serviceName, proxyName, itf)

	method := func(m object.MetaMethod, methodName string) error {
		method := itf.Methods[m.Uid]
		if serviceName != "Object" && serviceName != "Server" &&
			m.Uid < object.MinUserActionID {
			return nil
		}
		return generateMethod(file, set, proxyName, method, methodName)
	}
	signal := func(s object.MetaSignal, signalName string) error {
		signal := itf.Signals[s.Uid]
		signalName = util.CleanName("Subscribe" + signalName)
		if serviceName != "Object" && serviceName != "Server" &&
			s.Uid < object.MinUserActionID {
			return nil
		}
		return generateSignal(file, set, proxyName, signal, signalName)
	}
	property := func(p object.MetaProperty, propertyName string) error {
		property := itf.Properties[p.Uid]
		getMethodName := "Get" + propertyName
		setMethodName := "Set" + propertyName
		subscribeMethodName := "Subscribe" + propertyName
		if serviceName != "Object" && serviceName != "Server" &&
			p.Uid < object.MinUserActionID {
			return nil
		}
		return generateProperty(file, set, proxyName, property,
			getMethodName, setMethodName, subscribeMethodName)
	}

	metaObj := itf.MetaObject()
	if err := metaObj.ForEachMethodAndSignal(method, signal, property); err != nil {
		return fmt.Errorf("failed to generate proxy object %s: %s",
			serviceName, err)
	}
	return nil
}

func generateProxyType(file *jen.File, serviceName, proxyName string,
	itf *idl.InterfaceType) {

	file.Comment(proxyName + " implements " + serviceName)
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
	file.Comment("New" + serviceName + " constructs " + serviceName)
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
	file.Comment(serviceName + " retruns a proxy to a remote service")
	file.Func().Params(
		jen.Id("s").Id("Constructor"),
	).Id(
		util.CleanName(serviceName),
	).Params().Params(
		jen.Id(serviceName), jen.Error(),
	).Block(
		jen.Id(`return New` + serviceName + `(s.session, 1)`),
	)
}

func generateMethod(file *jen.File, set *signature.TypeSet, serviceName string,
	method idl.Method, methodName string) error {

	paramType := method.Tuple()
	paramType.RegisterTo(set)
	paramType.ConvertMetaObjects()

	returnType := method.Return

	// use object.MetaObject when generating proxy in order for
	// proxy to implement the object.Object interface.
	if returnType.Signature() == signature.MetaObjectSignature {
		returnType = signature.NewMetaObjectType()
	}
	returnType.RegisterTo(set)

	body, err := methodBodyBlock(method, paramType, returnType)
	if err != nil {
		return fmt.Errorf("failed to generate body: %s", err)
	}

	returnCode := jen.Params(returnType.TypeName(), jen.Error())
	if returnType.Signature() == "v" {
		returnCode = jen.Error()
	}

	goMethodName := util.CleanName(methodName)
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
		jen.Id(`propertyID, err := p.`+methodID+`("`+actionName+`")
		if err != nil {
			return nil, nil, fmt.Errorf("property %s not available: %s", "`+actionName+`", err)
		}

		handlerID, err := p.RegisterEvent(p.ObjectID(), propertyID, 0)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to register event for %s: %s", "`+actionName+`", err)
		}`),
		jen.Id("ch").Op(":=").Make(jen.Chan().Add(actionType.TypeName())),
		jen.List(
			jen.Id("cancel"),
			jen.Id("chPay"),
			jen.Err(),
		).Op(":=").Id("p.SubscribeID").Call(jen.Id("propertyID")),
		jen.Id(`if err != nil {
			return nil, nil, fmt.Errorf("failed to request property: %s", err)
		}`),
		jen.Go().Func().Params().Block(jen.For().Block(
			jen.List(jen.Id("payload"), jen.Id("ok")).Op(":=").Op("<-").Id("chPay"),
			jen.Id(`if !ok {
					// connection lost or cancellation.
					close(ch)
					p.UnregisterEvent(p.ObjectID(), propertyID, handlerID)
					return
				}`),
			jen.Id("buf := bytes.NewBuffer(payload)"),
			jen.Id("_ = buf // discard unused variable error"),
			jen.List(jen.Id("e"),
				jen.Err()).Op(":=").Add(actionType.Unmarshal("buf")),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Qual("log", "Printf").Call(
					jen.Lit("failed to unmarshall tuple: %s"),
					jen.Id("err"),
				),
				jen.Continue(),
			),
			jen.Id(`ch<- e`),
		)).Call(),
		jen.Return(jen.Id("cancel"), jen.Id("ch"), jen.Nil()),
	)

	file.Comment(methodName + " subscribe to a remote property")
	file.Func().Params(
		jen.Id("p").Op("*").Id(serviceName),
	).Id(methodName).Params().Add(retType).Add(body)
}

func generateSignal(file *jen.File, set *signature.TypeSet, serviceName string,
	signal idl.Signal, methodName string) error {

	signalType := signal.Type()
	signalType.RegisterTo(set)

	generateSubscribe(file, serviceName, signal.Name, methodName, signalType, true)

	return nil
}
func generateProperty(file *jen.File, set *signature.TypeSet, serviceName string,
	property idl.Property, getMethodName, setMethodName, subscribeMethodName string) error {

	propertyType := property.Type()
	propertyType.RegisterTo(set)

	generatePropertyGet(file, serviceName, property, getMethodName)
	generatePropertySet(file, serviceName, property, setMethodName)
	generateSubscribe(file, serviceName, property.Name, subscribeMethodName, propertyType, false)

	return nil
}
func generatePropertyGet(file *jen.File, serviceName string,
	property idl.Property, methodName string) error {

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
		jen.Id(`// discard the signature`),
		jen.Id(`_, err = basic.ReadString(&buf)`),
		jen.Id(`if err != nil {
		    return ret, fmt.Errorf("read signature: %s", err)
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
	property idl.Property, methodName string) error {

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

func methodBodyBlock(method idl.Method, params *signature.TupleType,
	ret signature.Type) (*Statement, error) {
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
		writing = append(writing, jen.Id(fmt.Sprintf(`response, err := p.Call("%s", buf.Bytes())`, method.Name)))
	} else {
		writing = append(writing, jen.Id(fmt.Sprintf(`_, err = p.Call("%s", buf.Bytes())`, method.Name)))
	}
	if ret.Signature() != "v" {
		writing = append(writing, jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id(fmt.Sprintf(`return ret, fmt.Errorf("call %s failed: %s", err)`, method.Name, "%s")),
		))
	} else {
		writing = append(writing, jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id(fmt.Sprintf(`return fmt.Errorf("call %s failed: %s", err)`, method.Name, "%s")),
		))
	}
	if ret.Signature() != "v" {
		writing = append(writing, jen.Id("buf = bytes.NewBuffer(response)"))
		writing = append(writing, jen.Id("ret, err =").Add(ret.Unmarshal("buf")))
		writing = append(writing, jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id(fmt.Sprintf(`return ret, fmt.Errorf("failed to parse %s response: %s", err)`, method.Name, "%s")),
		))
		writing = append(writing, jen.Return(jen.Id("ret"), jen.Nil()))
	} else {
		writing = append(writing, jen.Return(jen.Nil()))
	}

	return jen.Block(
		writing...,
	), nil
}
