## random notes on refactoring

mirror design

type Object interface {
        MetaObject(objectID uint32) (MetaObject, error)
        Terminate(objectID uint32) error
        RegisterEvent(objectID uint32, actionID uint32, handler uint64) (
                uint64, error)
        [...]
}

type Proxy interface {
        Call(action string, payload []byte) ([]byte, error)
        CallID(action uint32, payload []byte) ([]byte, error)
        Subscribe(action string) (func(), chan []byte, error)
        [...]
}

type ProxyMirror interface {
    Proxy
    Mirror
}

type ObjectMirror interface {
    Object
    Mirror
}

type Mirror interface {
    object() ObjectMirror
    proxy() ProxyMirror
}

type Meta struct {
    bus.Proxy
    object.Object
}

func (m Meta) Mirror() Mirror {
    return m
}


### Goals
    - uniform Type handling of proxy and stub
        - TODO: split custom proxy interface from bus.Proxy and object.Object
        - TODO: rename ObjectProxy into ProxyObject
        - TODO: introduce Mirror pattern into the proxy
        - TODO: rename interface server must implement: FileImplementor

        - TODO: activation helper to add objects (ServerObject)
        - TODO: when a impl method needs to create and return an
          object it needs to: (1) create a ServerObject and register
          it and (2) return a ObjectProxy

        - object Object:

            type Object interface {
                    MetaObject(objectID uint32) (MetaObject, error)
                    Terminate(objectID uint32) error
                    RegisterEvent(objectID uint32, actionID uint32, handler uint64) (
                            uint64, error)
                    [...]
            }

        - bus Proxy:

            type Proxy interface {
                    Call(action string, payload []byte) ([]byte, error)
                    CallID(action uint32, payload []byte) ([]byte, error)
                    Subscribe(action string) (func(), chan []byte, error)
                    [...]
            }

        - proxy Interface:

            type File interface {
                    object.Object
                    bus.Proxy
                    // Custom methods
                    Read(size int32) ([]int32, error)
            }

        - server object:

            type ServerObject interface {
                    Receive(m *net.Message, from *Context) error
                    Activate(activation Activation) error
                    OnTerminate()
            }

    - clarify the various Proxy and Object definitions and concepts
    - flatten package hirarchy for a simpler public API
    - unique qiloop command (scan, proxy, stub, info)
    - clean-up (stage 1, multiple sessions/client/proxy)
    - standalone server (log and directory)
    - compatible with lugu/audit package
    - reconsilliation: can pass local objects as parameter

### Actions
    + implement proxy generation as part of interface Type
    - implement stub generation as a method
    - local object shared by a server can be passed as argument

### New hierarchy

    ./meta/type.go // type and typeset definition
    ./meta/signature.go // parser
    ./meta/basic.go // basic types and composites
    ./meta/interface.go // object type
    ./meta/ref.go // scope and references
    ./meta/proxy.go: create proxy
    ./meta/stub.go: create stub
