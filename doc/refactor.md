## random notes on refactoring

### Goals
    - uniform Type handling of proxy and stub
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
