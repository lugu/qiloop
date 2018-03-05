package stage2

import (
	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/meta/stage1"
)

var MetaService0 stage1.MetaObject = stage1.MetaObject{
	Methods: map[uint32]stage1.MetaMethod{
		2: {
			Uid:                 2,
			ReturnSignature:     signature.MetaObjectSignature,
			Name:                "metaObject",
			ParametersSignature: "(I)",
			Description:         "request self description",
			Parameters: []stage1.MetaMethodParameter{
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
			Parameters: []stage1.MetaMethodParameter{
				{
					"", "",
				},
			},
			ReturnDescription: "",
		},
	},
	Description: "Server",
}
