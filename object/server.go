package object

import (
	"github.com/lugu/qiloop/meta/signature"
)

var MetaService0 MetaObject = MetaObject{
	Methods: map[uint32]MetaMethod{
		2: {
			Uid:                 2,
			ReturnSignature:     signature.MetaObjectSignature,
			Name:                "metaObject",
			ParametersSignature: "(I)",
			Description:         "request self description",
			Parameters: []MetaMethodParameter{
				{
					"", "",
				},
			},
			ReturnDescription: "",
		},
		8: {
			Uid:                 8,
			ReturnSignature:     "{sm}",
			Name:                "authenticate",
			ParametersSignature: "({sm})",
			Description:         "",
			Parameters: []MetaMethodParameter{
				{
					"", "",
				},
			},
			ReturnDescription: "",
		},
	},
	Description: "Server",
}
