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

func generateProxyType(file *jen.File, typ string, metaObj object.MetaObject) {

	file.Type().Id(typ).Struct(jen.Qual("github.com/lugu/qiloop/session", "Proxy"))
	file.Func().Id("New"+typ).Params(
		jen.Id("ses").Qual("github.com/lugu/qiloop/session", "Session"),
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
		methodName := registerName(m.Name, methodNames)
		if err := generateMethod(file, set, k, serviceName, m, methodName); err != nil {
			// FIXME: uncomment
			// return fmt.Errorf("failed to render method %s of %s: %s", m.Name, serviceName, err)
			fmt.Printf("failed to render method %s of %s: %s\n", m.Name, serviceName, err)
		}
	}
	signalNames := make(map[string]bool)
	keys = make([]int, 0)
	for k := range metaObj.Signals {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, i := range keys {
		k := uint32(i)
		s := metaObj.Signals[k]
		signalName := registerName(s.Name, signalNames)
		if err := generateSignal(file, set, k, serviceName, s, signalName); err != nil {
			// FIXME: uncomment
			// return fmt.Errorf("failed to render signal %s of %s: %s", s.Name, serviceName, err)
			fmt.Printf("failed to render signal %s of %s: %s\n", s.Name, serviceName, err)
		}
	}
	return nil
}

func GenerateProxy(metaObj object.MetaObject, packageName, serviceName string, w io.Writer) error {

	file := jen.NewFile(packageName)
	file.PackageComment("file generated. DO NOT EDIT.")
	set := signature.NewTypeSet()

	generateProxyObject(metaObj, serviceName, set, file)

	set.Declare(file)

	if err := file.Render(w); err != nil {
		return fmt.Errorf("failed to render %s: %s", serviceName, err)
	}
	return nil
}

func GenerateProxys(metaObjList []object.MetaObject, packageName string, w io.Writer) error {

	file := jen.NewFile(packageName)
	file.PackageComment("file generated. DO NOT EDIT.")

	set := signature.NewTypeSet()

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

func generateMethod(file *jen.File, set *signature.TypeSet, id uint32, typ string, m object.MetaMethod, methodName string) error {

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

	file.Func().Params(jen.Id("p").Op("*").Id(typ)).Id(strings.Title(methodName)).Add(
		tuple.Params(),
	).Add(
		retType,
	).Add(
		body,
	)
	return nil
}

func generateSignal(file *jen.File, set *signature.TypeSet, id uint32, typ string, s object.MetaSignal, signalName string) error {

	signalType, err := signature.Parse(s.Signature)
	if err != nil {
		return fmt.Errorf("failed to parse signal %s: %s", s.Signature, err)
	}

	signalType.RegisterTo(set)

	retType := jen.Params(jen.Chan().Add(signalType.TypeName()), jen.Error())
	body := jen.Block(jen.Return(jen.Nil(), jen.Nil()))

	file.Func().Params(jen.Id("p").Op("*").Id(typ)).Id("Signal" + strings.Title(signalName)).Params(
		jen.Id("cancel").Chan().Int(),
	).Add(
		retType,
	).Add(
		body,
	)
	/*
		func (o Object) Signal%d(cancel chan int) (chan %T, error) {

			chanPayload, err := SignalStream(%d, cancel)
			if err != nil {
				return fmt.Errorf("failed to get signal stream: %s", err)
			}
			values := make(chan %T)
			go func() {
				for {
					val e %T
					select {
					case payload, ok :=<- chanPayload:
						if !ok {
							// stream closed
							close(values)
							return
						}
						e, err := %s
						if err != nil {
							fmt.Errorf("failed to unmarshall %T: %s", err)
							continue
						}
						values<- e
					}
				}
			}
			return nil
		}
	*/
	return nil
}
