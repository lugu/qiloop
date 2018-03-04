package signature

import (
	object "qiloop/meta/stage1"
)

const MetaObjectSignature string = "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>"

var MetaService0 object.MetaObject = object.MetaObject{
	Methods: map[uint32]object.MetaMethod{
		2: object.MetaMethod{
			Uid:                 2,
			ReturnSignature:     MetaObjectSignature,
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
	Description: "Server",
}
