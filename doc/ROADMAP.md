for v0.5
--------

- doc: diagram which shows a bus
- doc: clarify if implementor are actors and document it in the tutorial
  and the generated documentation of the interface.
- document internal pipeline
- fix go get and GOPATH doc: decide which version to support
- service tutorial
- complete implemetation of log manager

for later
---------

- example: NAO say, Pepper say
- benchmark bandwidth
- simplify service creation (see clock example)
- gateway implementation
- websocket protocol
- implement raw data type (signature, Type and IDL)
- get ride of implicit tokens when creating a session (Option)
- session: re-use established connections
- repair stats and traces
- actor: use Terminate message instead of OnTerminate.
- QUIC implementation: multiplex client
