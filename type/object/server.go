package object

import (
	"github.com/lugu/qiloop/meta/signature"
)

const AuthenticateActionID = 8

var MetaService0 MetaObject = MetaObject{
	Description: "Server",
	Methods: map[uint32]MetaMethod{
		AuthenticateActionID: MetaMethod{
			Uid:                 AuthenticateActionID,
			ReturnSignature:     "{sm}",
			Name:                "authenticate",
			ParametersSignature: "({sm})",
			Parameters: []MetaMethodParameter{
				MetaMethodParameter{
					Name: "capability",
				},
			},
		},
	},
}

var ObjectMetaObject MetaObject = MetaObject{
	Description: "Object",
	Methods: map[uint32]MetaMethod{
		0x0: MetaMethod{
			Uid:                 0x0,
			ReturnSignature:     "L",
			Name:                "registerEvent",
			ParametersSignature: "(IIL)",
		},
		0x1: MetaMethod{
			Uid:                 0x1,
			ReturnSignature:     "v",
			Name:                "unregisterEvent",
			ParametersSignature: "(IIL)",
		},
		0x2: MetaMethod{
			Uid:                 0x2,
			ReturnSignature:     signature.MetaObjectSignature,
			Name:                "metaObject",
			ParametersSignature: "(I)",
		},
		0x3: MetaMethod{
			Uid:                 0x3,
			ReturnSignature:     "v",
			Name:                "terminate",
			ParametersSignature: "(I)",
		},
		0x5: MetaMethod{
			Uid:                 0x5,
			ReturnSignature:     "m",
			Name:                "property",
			ParametersSignature: "(m)",
		},
		0x6: MetaMethod{
			Uid:                 0x6,
			ReturnSignature:     "v",
			Name:                "setProperty",
			ParametersSignature: "(mm)",
		},
		0x7: MetaMethod{
			Uid:                 0x7,
			ReturnSignature:     "[s]",
			Name:                "properties",
			ParametersSignature: "()",
		},
		0x8: MetaMethod{
			Uid:                 0x8,
			ReturnSignature:     "L",
			Name:                "registerEventWithSignature",
			ParametersSignature: "(IILs)",
		},
	},
	Signals: map[uint32]MetaSignal{
		0x56: MetaSignal{
			Uid:       0x56,
			Name:      "traceObject",
			Signature: "((IiIm(ll)<timeval,tv_sec,tv_usec>llII)<EventTrace,id,kind,slotId,arguments,timestamp,userUsTime,systemUsTime,callerContext,calleeContext>)",
		},
	},
}
