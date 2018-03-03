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
    // type ObjectProxy struct {
    //     Proxy
    // }
    file.Type().Id(typ).Struct(jen.Id("Proxy"))
}

func GenerateProxy(metaObj object.MetaObject, packageName, serviceName string, w io.Writer) error {

    var file *jen.File = jen.NewFile(packageName)
    generateProxyType(file, serviceName, metaObj)
    for k, m := range metaObj.Methods {
        if err := generateMethod(file, k, serviceName, m); err != nil {
            return fmt.Errorf("failed to render method %s of %s: %s", m.Name, serviceName, err)
        }
    }
    if err := file.Render(w); err != nil {
        return fmt.Errorf("failed to render %s: %s", serviceName, err)
    }
    return nil
}

func generateMethodArguments(m object.MetaMethod) (*Statement, error) {
    _, err := signature.Parse(m.ParametersSignature)
    if err != nil {
        fmt.Printf("failed to parse argument signature %s: %s\n", m.ParametersSignature, err)
        // return nil, fmt.Errorf("failed to parse signature %s: %s", m.ParametersSignature, err)
    }
    return jen.Empty(), nil
}

func generateMethodReturnType(m object.MetaMethod) (*Statement, error) {
    _, err := signature.Parse(m.ReturnSignature)
    if err != nil {
        fmt.Printf("failed to parse return signature %s: %s\n", m.ReturnSignature, err)
        // return nil, fmt.Errorf("failed to parse signature %s: %s", m.ParametersSignature, err)
    }
    return jen.Error(), nil
}

func generateMethodBody(m object.MetaMethod) (*Statement, error) {
    return jen.Empty(), nil
}

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
    file.Func().Params(jen.Id("p").Op("*").Id(typ)).Id(strings.Title(m.Name)).Params(
        arguments,
    ).Params(
        returnedType,
    ).Block(
        body,
    )
    return nil
}


// warning: always add an error to the signature
// TODO:
// - create a type which inherited from Object
// - for each method, construct the method like:
//      - generate the method signature
//      - generate the parameter serialisation: func(a, b uint) []byte { }
//      - generate the Proxy call
//      - generate the deserialization func(p []byte) uint32 { }

// func (p *object.Proxy) RegisterEvent(a, b uint32, c uint64) (uint64, error) {
//     var parameterBytes []byte // TODO: serialized params
//     _, err := p.Call(0, parameterBytes)
//     if err != nil {
//         return 0, fmt.Errorf("call to Register event failed: %s", err)
//     }
//     var returned uint64 // TODO unserialize response
//     return returned, nil
// }
