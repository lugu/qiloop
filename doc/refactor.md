## random notes on refactoring

### Goals
    - uniform Type handling of proxy and stub
          get both a ServerObject and a ProxyObject: registering a
          ServerObject shall return an ObjectReference.

    - unique qiloop command (scan, proxy, stub, info)
    - get ride of implicit tokens when creating a session
    - remove bus/client/object.MetaObject: duplicated!!

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
