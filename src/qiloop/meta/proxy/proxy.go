package proxy

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"io"
	"qiloop/meta/signature"
	"qiloop/object"
	"strings"
)

type Statement = jen.Statement

func generateProxyType(file *jen.File, typ string, metaObj object.MetaObject) {
	file.Type().Id(typ).Struct(jen.Qual("qiloop/object", "Proxy"))
	file.Func().Id("New"+typ).Params(
		jen.Id("endpoint").String(),
		jen.Id("service").Uint32(),
		jen.Id("obj").Uint32(),
	).Params(
		jen.Op("*").Id(typ),
		jen.Error(),
	).Block(
		jen.If(
			jen.List(jen.Id("conn"), jen.Err()).Op(":=").Qual("qiloop/net", "NewClient").Call(
				jen.Id("endpoint"),
			).Op(";").Err().Op("!=").Nil()).Block(
			jen.Id(`return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)`),
		).Else().Block(
			jen.Id(`proxy := object.NewProxy(conn, service, obj)`),
			jen.Id(`return &`+typ+`{ proxy }, nil`),
		),
	)
}

func GenerateProxy(metaObj object.MetaObject, packageName, serviceName string, w io.Writer) error {

	file := jen.NewFile(packageName)
	set := signature.NewTypeSet()
	generateProxyType(file, serviceName, metaObj)
	for k, m := range metaObj.Methods {
		if err := generateMethod(file, set, k, serviceName, m); err != nil {
			// FIXME: uncomment
			// return fmt.Errorf("failed to render method %s of %s: %s", m.Name, serviceName, err)
			fmt.Printf("failed to render method %s of %s: %s\n", m.Name, serviceName, err)
		}
	}
	set.Declare(file)
	if err := file.Render(w); err != nil {
		return fmt.Errorf("failed to render %s: %s", serviceName, err)
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
		writing[i] = jen.Id(fmt.Sprintf(`response, err := p.Call(%d, buf.Bytes())`, m.Uid))
		i++
	} else {
		writing[i] = jen.Id(fmt.Sprintf(`_, err = p.Call(%d, buf.Bytes())`, m.Uid))
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

func generateMethod(file *jen.File, set *signature.TypeSet, id uint32, typ string, m object.MetaMethod) error {
	tupleType, err := signature.Parse(m.ParametersSignature)
	if err != nil {
		return fmt.Errorf("failed to parse signature %s: %s", m.ParametersSignature, err)
	}
	params, ok := tupleType.(*signature.TupleValue)
	if !ok {
		return fmt.Errorf("failed to parse method parameters: expected a tuple, got %#v", tupleType)
	}
	params.RegisterTo(set)

	ret, err := signature.Parse(m.ReturnSignature)
	if err != nil {
		return fmt.Errorf("failed to parse return signature %s: %s", m.ReturnSignature, err)
	}
	ret.RegisterTo(set)

	body, err := methodBodyBlock(m, params, ret)
	if err != nil {
		return fmt.Errorf("failed to generate body: %s", err)
	}

	retType := jen.Params(ret.TypeName(), jen.Error())
	if _, ok := ret.(signature.VoidValue); ok {
		retType = jen.Error()
	}

	file.Func().Params(jen.Id("p").Op("*").Id(typ)).Id(strings.Title(m.Name)).Add(
		params.Params(),
	).Add(
		retType,
	).Add(
		body,
	)
	return nil
}
