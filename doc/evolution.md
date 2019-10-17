Possible evolutions
===================

CapabilityMap
-------------

Problem: the authentication procedure requires the client to send a
capability map. This data structure contains values which can be
arbitrary complex. The server is forced to parse those values.

Solution: Make the CapabilityMap type a string to string map.

Peer to peer authentication
---------------------------

Problem: only the service directory connection is authenticated.
Peer to peer connection requires authentication.

Solution: Once authenticated with the service directory, use the
session ID of the service info data as an authenticating token when
contacting other servers.

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

Problem: A client of a given service can guess the object ID
designated to another client and hijack it.

Solution: Use true random number for service and object ID.

Object isolation
----------------

Problem: The following methods take an object ID parameter for no good
reason: MetaObject, Terminate, RegisterEvent, UnregisterEvent,
RegisterEventWithSignature. Implementing them correctly is complex in
the case of client objects.

Solution: Remove the object ID parameter from the listed methods.

Object ID collision
-------------------

Problem: clients of the same service can create objects with the same
object ID. This breaks the assumption that the destination of a
message is unique.

Solution: let the service decide of the object ID using a generic
method such as: leaseObject -> uint32.

Client object routing
---------------------

Problem: client objects impose some kind of routing. Since
signal/property registration is done using method calls, the service
needs to inspect the message payload in order know who subscribes to
what.

Solution: add some new types of messages: subscribe and unsubscribe.
Using those message type an implementation of a router would not need
to analyze the payload anymore (such lightweight implementation do not
requires supports for signatures and serialization).
