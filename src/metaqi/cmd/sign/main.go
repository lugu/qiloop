package main

import (
    "fmt"
    "github.com/dave/jennifer/jen"
    "metaqi"
    "os"
)

func main() {
    if len(os.Args) > 1 {
        signature := os.Args[1]
        typeDescription, err := metaqi.Parse(signature)
        if err != nil {
            fmt.Printf("parsing error: %s\n", err)
        }
        var file *jen.File = jen.NewFile("object")
        typeDescription.TypeDeclaration(file)
        file.Render(os.Stdout)
    } else {
        fmt.Printf("missing signature argument\n")
    }
}
