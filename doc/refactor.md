## random notes on refactoring

### Goals
    - uniform Type handling of proxy and stub
        - TODO: rename service constructor with explicit name like
          NewBombService() etc..
        - TODO: merge proxy and stub together
        - TODO: new stub object constructor:
                Create<Name>(Service) ObjectReference
          get both a ServerObject and a ProxyObject: registering a
          ServerObject shall return an ObjectReference.


    - clarify the various Proxy and Object definitions and concepts
    - flatten package hirarchy for a simpler public API
    - unique qiloop command (scan, proxy, stub, info)
    - clean-up (stage 1, multiple sessions/client/proxy)
    - standalone server (log and directory)
    - compatible with lugu/audit package
    - reconsilliation: can pass local objects as parameter

### Experience writing qitop

1. Difficult to get an object from a ServiceInfo
        => implement decorator

2. Too many packages: do not know where to search. TermUI is already a

3. Confusing concepts object.Object, object.ProxyObject, .... Need
   some documentation
        => session
        => proxy
        => ProxyObject
        => server
        => ServerObject

4. missing standalone server ready to use (like audit/honey)

5. better bridge between ServerObject <= ObjectReference <=> ProxyObject

### New hierarchy

    ./meta/type.go // type and typeset definition
    ./meta/signature.go // parser
    ./meta/basic.go // basic types and composites
    ./meta/interface.go // object type
    ./meta/ref.go // scope and references
    ./meta/proxy.go: create proxy
    ./meta/stub.go: create stub
