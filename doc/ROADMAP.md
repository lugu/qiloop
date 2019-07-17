for v0.5
--------

- doc: diagram which shows a bus
- complete implementation of log manager
- fix go get and GOPATH doc: decide which version to support

for later
---------

- example: NAO say, Pepper say
- benchmark bandwidth (128, 512, 1024, 4096 bytes per message)
- benchmark latency: (50, 75, 90, 99, 99.9, 99.99, 99.999 percentil)
- transport: fifo / pipe
- transport: socketpair
- transport: websocket
- gateway implementation
- implement raw data type (signature, Type and IDL)
- get ride of implicit tokens when creating a session (Option)
- session: re-use established connections
- repair stats and traces
- actor: use Terminate message instead of OnTerminate.
- QUIC implementation: multiplex client
