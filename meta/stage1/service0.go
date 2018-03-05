package stage1

import (
	"github.com/lugu/qiloop/meta/signature"
)

var MetaService0 MetaObject = MetaObject{
	Methods: map[uint32]MetaMethod{
		2: MetaMethod{
			Uid:                 2,
			ReturnSignature:     signature.MetaObjectSignature,
			Name:                "metaObject",
			ParametersSignature: "(I)",
			Description:         "request self description",
			Parameters: []MetaMethodParameter{
				MetaMethodParameter{
					"", "",
				},
			},
			ReturnDescription: "",
		},
		8: MetaMethod{
			Uid:                 8,
			ReturnSignature:     "{sm}",
			Name:                "authenticate",
			ParametersSignature: "({sm})",
			Description:         "",
			Parameters: []MetaMethodParameter{
				MetaMethodParameter{
					"", "",
				},
			},
			ReturnDescription: "",
		},
	},
	Description: "Server",
}
