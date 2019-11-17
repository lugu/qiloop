
BUG
---

- session: re-use established connections
- clent: cache signal registration
- logger: update listerner interface
- client: service directory disconnect close all connections
- implement raw data type (signature, Type and IDL)

FIXME
-----

- remove concurrent actor access (see doc/internal-design.md)
- actor: use Terminate message instead of OnTerminate.
- QUIC implementation: multiplex client
- get ride of implicit tokens when creating a session (Option)

TODO
----

- document socket sementic
- visual doc: diagram which shows a bus
- benchmark bandwidth (128, 512, 1024, 4096 bytes per message)
- benchmark latency: (50, 75, 90, 99, 99.9, 99.99, 99.999 percentil)
- transport: socketpair
- transport: websocket
- gateway implementation
