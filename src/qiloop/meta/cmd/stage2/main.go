package main

import (
	"io"
	"log"
	"os"
	"qiloop/meta/proxy"
	"qiloop/meta/signature"
	object "qiloop/meta/stage1"
)

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

	metaObj := object.MetaObject{
		Methods: map[uint32]object.MetaMethod{
			2: object.MetaMethod{
				Uid:                 2,
				ReturnSignature:     signature.MetaObjectSignature,
				Name:                "metaObject",
				ParametersSignature: "(I)",
				Description:         "request self description",
				Parameters: []object.MetaMethodParameter{
					object.MetaMethodParameter{
						"", "",
					},
				},
				ReturnDescription: "",
			},
			8: object.MetaMethod{
				Uid:                 8,
				ReturnSignature:     "{sm}",
				Name:                "authenticate",
				ParametersSignature: "({sm})",
				Description:         "",
				Parameters: []object.MetaMethodParameter{
					object.MetaMethodParameter{
						"", "",
					},
				},
				ReturnDescription: "",
			},
		},
	}

	err := proxy.GenerateProxy(metaObj, "stage2", "Server", output)
	if err != nil {
		log.Fatalf("proxy generation failed: %s\n", err)
	}
	return
}
