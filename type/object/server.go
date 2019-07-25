package object

import (
	"github.com/lugu/qiloop/meta/signature"
)

// AuthenticateActionID represents the action id of the authentication
// process.
const AuthenticateActionID = uint32(8)

// MetaService0 represents the service id "0" used during
// authentication. It is a special service with only one method.
var MetaService0 MetaObject = MetaObject{
	Description: "Server",
	Methods: map[uint32]MetaMethod{
		AuthenticateActionID: {
			Uid:                 AuthenticateActionID,
			ReturnSignature:     "{sm}",
			Name:                "authenticate",
			ParametersSignature: "({sm})",
			Parameters: []MetaMethodParameter{
				{
					Name: "capability",
				},
			},
		},
	},
	Signals:    make(map[uint32]MetaSignal),
	Properties: make(map[uint32]MetaProperty),
}

// ObjectMetaObject represents the generic actions all services (and
// objects) have (except service zero, see MetaService0).
var ObjectMetaObject MetaObject = MetaObject{
	Description: "Object",
	Methods: map[uint32]MetaMethod{
		0x0: {
			Uid:                 0x0,
			ReturnSignature:     "L",
			Name:                "registerEvent",
			ParametersSignature: "(IIL)",
		},
		0x1: {
			Uid:                 0x1,
			ReturnSignature:     "v",
			Name:                "unregisterEvent",
			ParametersSignature: "(IIL)",
		},
		0x2: {
			Uid:                 0x2,
			ReturnSignature:     signature.MetaObjectSignature,
			Name:                "metaObject",
			ParametersSignature: "(I)",
		},
		0x3: {
			Uid:                 0x3,
			ReturnSignature:     "v",
			Name:                "terminate",
			ParametersSignature: "(I)",
		},
		0x5: {
			Uid:                 0x5,
			ReturnSignature:     "m",
			Name:                "property",
			ParametersSignature: "(m)",
		},
		0x6: {
			Uid:                 0x6,
			ReturnSignature:     "v",
			Name:                "setProperty",
			ParametersSignature: "(mm)",
		},
		0x7: {
			Uid:                 0x7,
			ReturnSignature:     "[s]",
			Name:                "properties",
			ParametersSignature: "()",
		},
		0x8: {
			Uid:                 0x8,
			ReturnSignature:     "L",
			Name:                "registerEventWithSignature",
			ParametersSignature: "(IILs)",
		},
		uint32(0x50): {
			Name:                "isStatsEnabled",
			ParametersSignature: "()",
			ReturnSignature:     "b",
			Uid:                 uint32(0x50),
		},
		uint32(0x51): {
			Name:                "enableStats",
			ParametersSignature: "(b)",
			ReturnSignature:     "v",
			Uid:                 uint32(0x51),
		},
		uint32(0x52): {
			Name:                "stats",
			ParametersSignature: "()",
			ReturnSignature:     "{I(I(fff)<MinMaxSum,minValue,maxValue,cumulatedValue>(fff)<MinMaxSum,minValue,maxValue,cumulatedValue>(fff)<MinMaxSum,minValue,maxValue,cumulatedValue>)<MethodStatistics,count,wall,user,system>}",
			Uid:                 uint32(0x52),
		},
		uint32(0x53): {
			Name:                "clearStats",
			ParametersSignature: "()",
			ReturnSignature:     "v",
			Uid:                 uint32(0x53),
		},
		uint32(0x54): {
			Name:                "isTraceEnabled",
			ParametersSignature: "()",
			ReturnSignature:     "b",
			Uid:                 uint32(0x54),
		},
		uint32(0x55): {
			Name:                "enableTrace",
			ParametersSignature: "(b)",
			ReturnSignature:     "v",
			Uid:                 uint32(0x55),
		},
	},
	Signals: map[uint32]MetaSignal{
		0x56: {
			Uid:       0x56,
			Name:      "traceObject",
			Signature: "((IiIm(ll)<timeval,tv_sec,tv_usec>llII)<EventTrace,id,kind,slotId,arguments,timestamp,userUsTime,systemUsTime,callerContext,calleeContext>)",
		},
	},
	Properties: make(map[uint32]MetaProperty),
}
