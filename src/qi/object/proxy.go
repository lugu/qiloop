package object

import (
    "qi/net"
)

type Proxy struct {
	client      net.Client
	service     uint32
	object      uint32
	description string
}

func (p Proxy) Call(action uint32, payload []byte) ([]byte, error) {
	return p.client.Call(p.service, p.object, action, payload)
}

func NewProxy(c net.Client, service, object uint32, description string) Proxy {
	return Proxy{c, service, object, description}
}

type Value interface {
    Signature() string
    Value() interface{}
}

type Object interface {
    RegisterEvent(uint32,uint32,uint64) (uint64, error)
    UnregisterEvent(uint32,uint32,uint64) error
    MetaObject(uint32) (MetaObject, error)
    Terminate(uint32) error
    Property(Value) (Value, error)
    SetProperty(Value,Value) error
    Properties() ([]string, error)
    RegisterEventWithSignature(uint32,uint32,uint64,string) (uint64, error)
}

