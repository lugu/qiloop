package main

import (
	"io"
	"log"
	"os"
	"qiloop/meta/signature"
)

/*
Bootstrap stages:

    1. Extract the MetaObject signature using wireshark
    2. Generate the code of MetaObject type and MetaObject constructor:
        type MetaObject struct { ... }
        func NewMetaObject(io.Reader) MetaObject { ... }
    3. Extract the MetaObject data using wireshark
    4. Parse the MetaObject binary data and generate the ServiceDirectory proxy
        type ServiceDirectory { ... }
        func NewServiceDirectory(...) ServiceDirectory
    5. Construct the MetaObject of each services declared in the ServiceDirectory
    6. Parse the MetaObject of each service and generate the associated proxy
        type ServiceXXX { ... }
        func NewServiceXXX(...) ServiceXXX

*/

const MetaObjectSignature string = "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>"

func main() {
	var output io.Writer

	if len(os.Args) > 1 {
		filename := os.Args[1]

		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
			return
		}
		output = file
		defer file.Close()
	} else {
		output = os.Stdout
	}

	typ, err := signature.Parse(MetaObjectSignature)
	if err != nil {
		log.Fatalf("parsing error: %s\n", err)
	}
	if err = signature.GenerateType(typ, "object", output); err != nil {
		log.Fatalf("code generation failed: %s\n", err)
	}
}
