package main

import (
    "fmt"
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
        fmt.Printf(typeDescription.TypeDeclaration())
    } else {
        fmt.Printf("missing signature argument\n")
    }
}
