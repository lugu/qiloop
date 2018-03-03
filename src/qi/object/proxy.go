package object

import (
    "qi/net"
    "fmt"
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

type ObjectProxy struct {
    Proxy
}

// warning: always add an error to the signature
// TODO:
// - create a type which inherited from Object
// - for each method, construct the method like:
//      - generate the method signature
//      - generate the parameter serialisation: func(a, b uint) []byte { }
//      - generate the Proxy call
//      - generate the deserialization func(p []byte) uint32 { }
func (p *ObjectProxy) RegisterEvent(a, b uint32, c uint64) (uint64, error) {
    var parameterBytes []byte // TODO: serialized params
    _, err := p.Call(0, parameterBytes)
    if err != nil {
        return 0, fmt.Errorf("call to Register event failed: %s", err)
    }
    var returned uint64 // TODO unserialize response
    return returned, nil
}
