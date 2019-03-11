## random notes on refactoring

### Goals
    - uniform Type handling of proxy and stub
          get both a ServerObject and a ProxyObject: registering a
          ServerObject shall return an ObjectReference.

    - clarify the various Proxy and Object definitions and concepts
    - flatten package hirarchy for a simpler public API
    - unique qiloop command (scan, proxy, stub, info)
    - standalone server (log and directory)

### Experience writing qitop

2. Too many packages: do not know where to search. TermUI is already a

3. Confusing concepts object.Object, object.ProxyObject, .... Need
   some documentation
        => session
        => proxy
        => ProxyObject
        => server
        => ServerObject

5. better bridge between ServerObject <= ObjectReference <=> ProxyObject

### New hierarchy

    ./meta/type.go // type and typeset definition
    ./meta/signature.go // parser
    ./meta/basic.go // basic types and composites
    ./meta/interface.go // object type
    ./meta/ref.go // scope and references
    ./meta/proxy.go: create proxy
    ./meta/stub.go: create stub
