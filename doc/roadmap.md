for v0.8
--------

- visual doc: diagram which shows a bus
- benchmark bandwidth (128, 512, 1024, 4096 bytes per message)
- benchmark latency: (50, 75, 90, 99, 99.9, 99.99, 99.999 percentil)
- repair stats
- example: NAO say, Pepper say

for later
---------

- logger: update listerner interface
- trace: implements statistics
- transport: socketpair
- transport: websocket
- remove concurrent actor access (see doc/internal-design.md)
- gateway implementation
- implement raw data type (signature, Type and IDL)
- get ride of implicit tokens when creating a session (Option)
- session: re-use established connections
- actor: use Terminate message instead of OnTerminate.
- QUIC implementation: multiplex client
