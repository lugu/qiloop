package proxy

import (
    "io"
    "strings"
    "fmt"
    "qi/object"
    "qi/meta/signature"
	"github.com/dave/jennifer/jen"
)

type Statement = jen.Statement

func generateProxyType(file *jen.File, typ string, metaObj object.MetaObject) {
    file.Type().Id(typ).Struct(jen.Id("Proxy"))
}

func GenerateProxy(metaObj object.MetaObject, packageName, serviceName string, w io.Writer) error {

    var file *jen.File = jen.NewFile(packageName)
    generateProxyType(file, serviceName, metaObj)
    for k, m := range metaObj.Methods {
        if err := generateMethod(file, k, serviceName, m); err != nil {
            // FIXME
            // return fmt.Errorf("failed to render method %s of %s: %s", m.Name, serviceName, err)
            fmt.Printf("failed to render method %s of %s: %s\n", m.Name, serviceName, err)
        }
    }
    if err := file.Render(w); err != nil {
        return fmt.Errorf("failed to render %s: %s", serviceName, err)
    }
    return nil
}

func generateMethodArguments(m object.MetaMethod) (*Statement, error) {
    argument, err := signature.Parse(m.ParametersSignature)
    if err != nil {
        return nil, fmt.Errorf("failed to parse signature %s: %s", m.ParametersSignature, err)
    }
    tupleType, ok := argument.(*signature.TupleValue)
    if !ok {
        return nil, fmt.Errorf("failed to parse method parameters: expected a tuple, got %#v", argument)
    }
    return tupleType.Params(), nil
}

func generateMethodReturnType(m object.MetaMethod) (*Statement, error) {
    typ, err := signature.Parse(m.ReturnSignature)
    if err != nil {
        return nil, fmt.Errorf("failed to parse signature %s: %s", m.ParametersSignature, err)
    }
    return jen.Params(typ.TypeName(), jen.Error()), nil
}

func generateMethodBody(m object.MetaMethod) (*Statement, error) {
    return jen.Block(
        jen.Var().Err().Error(),
        jen.Var().Id("ret").Error(),
        jen.Return(jen.Id("ret"), jen.Err()),
    ), nil
}

// func (p *object.Proxy) RegisterEvent(a, b uint32, c uint64) (uint64, error) {
//     var parameterBytes []byte // TODO: serialized params
//     _, err := p.Call(0, parameterBytes)
//     if err != nil {
//         return 0, fmt.Errorf("call to Register event failed: %s", err)
//     }
//     var returned uint64 // TODO unserialize response
//     return returned, nil
// }

func generateMethod(file *jen.File, id uint32, typ string, m object.MetaMethod) error {
    arguments, err := generateMethodArguments(m)
    if err != nil {
        return fmt.Errorf("failed to parse arguments: %s", err)
    }
    returnedType, err := generateMethodReturnType(m)
    if err != nil {
        return fmt.Errorf("failed to parse arguments: %s", err)
    }
    body, err := generateMethodBody(m)
    if err != nil {
        return fmt.Errorf("failed to generate body: %s", err)
    }
    file.Func().Params(jen.Id("p").Op("*").Id(typ)).Id(strings.Title(m.Name)).Add(
        arguments,
    ).Add(
        returnedType,
    ).Add(
        body,
    )
    return nil
}

