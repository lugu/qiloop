Possible evolutions
===================

Authentication
--------------

Problem: only the service directory connection is authenticated.
Peer to peer connection requires authentication.

Solution: Once authenticated with the service directory, use the
session id of the service info data as a token when contacting other
services.

Local transport
---------------

Problem: TCP communication escape the operating system permissions
mecanism.

Solution: Using a UNIX socket would allow the operating system to
attach it permissions in order to enforce a security model.

Wireless transport
------------------

Problem: Two concurent calls are serialized into a single TCP
connection. Once a packet is loss, both calls are delayed.

Solution: QUIC seems to offer an attractive alternative to TCP for
wireless connections.

Browser support
---------------

Problem: Cannot simply create a web page which interract with the
bus. A GUI need an heavy client.

Solution: WebSocket would allow pure JS implementations.

Object privacy
--------------

Problem: A client of a given service can guess the object id
destinated to another client and hijack it.

Solution: Use true random number for service and object id.
