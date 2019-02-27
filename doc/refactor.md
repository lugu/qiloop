## random notes on refactoring

### Goals
    - uniform Type handling of proxy and stub
    - flatten package hirarchy for a simpler public API
    - unique qiloop command (scan, proxy, stub, info)
    - clean-up (stage 1, multiple sessions/client/proxy)
    - standalone server (log and directory)
    - compatible with lugu/audit package
    - reconsilliation: can pass local objects as parameter

### Actions
    - implement proxy generation as part of interface Type
    - implement stub generation as a method
    - local object shared by a server can be passed as argument

#### Proxy changes

Problem statement:
- proxy shall take specialized objects and return specialized objects
- framework: proxies are interface Type:
        - marshall/unmarshall them
        - have a constructor

- concept: allow for client side stub and router implementation so
  clients can instancitate objects.

- specialized objects are
        type CustomService interface {
                object.Object
                bus.Proxy
                [... methods ...]
        }
- actually is:
        type ServiceDirectoryProxy struct {
                object1.ObjectProxy
        }
        type ObjectProxy struct {
                bus.Proxy
                [... object.Object methods ...]
        }

### Current status
- signature:
     - types from signature
     - generate a Type from a signature (move to stage1)
     - declare Type and types (can be moved)
     - depends on Type
- idl:
     - parse an IDL into a package declaration
     	- depends on Type
     - generate an IDL base on a meta object
     	- shall be rewriteen into (1) metaobject to interface
     	type and (2) interface type into IDL
     	- depends on Type
- proxy:
     - generate a proxy
- stub:
     - generate a stub
- stage1:
     - uses signature: generate a type

### New hierarchy

    ./meta/type.go // type and typeset definition
    ./meta/signature.go // parser
    ./meta/basic.go // basic types and composites
    ./meta/interface.go // object type
    ./meta/ref.go // scope and references
    ./meta/proxy.go: create proxy
    ./meta/stub.go: create stub


TODO: Should just call typ.RegisterTo(set) whatever the type is.
the set of type shall be rendered as without having to worry about
that type it is. Kill de dependancy to meta/idl package.

TODO: review the meta hierarchy: what are the types? shall all the
types be in meta/signature?

TODO: add a value constructor to signature.Type?

Goal: simplify and clarify the component.
