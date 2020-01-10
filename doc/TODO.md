## BUG

-   SD: unregister services on disconnection

## FIXME

-   get ride of implicit tokens when creating a session (Option)
-   client: service directory disconnect close all connections
-   implement raw data type (signature, Type and IDL)
-   remouve double message buffering (see doc/internal-design.md)
-   remove concurrent actor access (see doc/internal-design.md)
-   actor: use Terminate message instead of OnTerminate.
-   QUIC implementation: multiplex client

## TODO

-   struct serialization: default, optional and out of order
-   visual doc: diagram which shows a bus
-   benchmark bandwidth (128, 512, 1024, 4096 bytes per message)
-   benchmark latency: (50, 75, 90, 99, 99.9, 99.99, 99.999 percentil)
-   transport: socketpair
-   transport: websocket
-   gateway implementation
-   refactor encoding: get it from the context
-   dbus converter (namspace + encoding + transport)

# Encoding introduction

1. Make sure encoding does the job
2. Implement Proxy.Call using forced encoding
3. Update proxy gen to use Proxy.Call
4. Channel can build decoder from msg
5. Update stub gen to use Decoder
6. Remove Marshall/Unmarshall

# proxy abstractions

    // New method:
    Proxy.Call(action, signature string, arg, ret interface{}) error
    // keep backward compatibility
    Proxy.CallID(action uint32, payload []byte) ([]byte, error)

1. find the correct action id based on the name and the signature of
   argument.
2. type conversion and encoding
