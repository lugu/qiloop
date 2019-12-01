BUG
---

- SD: unregister services on disconnection

FIXME
-----

- get ride of implicit tokens when creating a session (Option)
- client: service directory disconnect close all connections
- implement raw data type (signature, Type and IDL)
- remouve double message buffering (see doc/internal-design.md)
- remove concurrent actor access (see doc/internal-design.md)
- actor: use Terminate message instead of OnTerminate.
- QUIC implementation: multiplex client

TODO
----

- struct serialization: default, optional and out of order
- visual doc: diagram which shows a bus
- benchmark bandwidth (128, 512, 1024, 4096 bytes per message)
- benchmark latency: (50, 75, 90, 99, 99.9, 99.99, 99.999 percentil)
- transport: socketpair
- transport: websocket
- gateway implementation
- cancel: include context.Context in context
- refactor encoding: get it from the context
- dbus converter (namspace + encoding + transport)
