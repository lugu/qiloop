# About QiMessaging

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Presentation](#presentation)
  - [Introduction](#introduction)
  - [Overview](#overview)
  - [Comparisons](#comparisons)
  - [OSI Model](#osi-model)
- [Messages](#messages)
  - [Message Header](#message-header)
    - [Magic](#magic)
    - [Message ID](#message-id)
    - [Size](#size)
    - [Version](#version)
    - [Type](#type)
    - [Flags](#flags)
    - [Service ID](#service-id)
    - [Object ID](#object-id)
    - [Action ID](#action-id)
  - [Payload](#payload)
- [Signatures](#signatures)
  - [Types](#types)
    - [Basic types](#basic-types)
    - [Value](#value)
    - [Composite types](#composite-types)
    - [Object](#object)
    - [Other](#other)
  - [Examples](#examples)
    - [Real world examples](#real-world-examples)
  - [Signature Grammar](#signature-grammar)
- [Serialization](#serialization)
  - [Basic types](#basic-types-1)
  - [Values](#values)
  - [Composite types (map, vect and struct)](#composite-types-map-vect-and-struct)
  - [Object](#object-1)
- [Objects](#objects)
  - [Methods](#methods)
    - [Generic methods](#generic-methods)
    - [Specific methods](#specific-methods)
  - [Signals](#signals)
  - [Properties](#properties)
  - [MetaObject](#metaobject)
- [Services](#services)
  - [Service Server (ID 0)](#service-server-id-0)
  - [Service Directory (ID 1)](#service-directory-id-1)
    - [Methods](#methods-1)
    - [Signals](#signals-1)
  - [Example (LogManager)](#example-logmanager)
- [Networking](#networking)
  - [Endpoints](#endpoints)
  - [TCP](#tcp)
  - [SSL](#ssl)
- [Authentication](#authentication)
  - [CapabilityMap](#capabilitymap)
    - [Authentication state](#authentication-state)
    - [Authentication credentials](#authentication-credentials)
    - [Protocol feature negotiation](#protocol-feature-negotiation)
- [Routing](#routing)
  - [Message destination](#message-destination)
  - [Message origin](#message-origin)
  - [Message transferred](#message-transferred)
- [Misc](#misc)
  - [Object Statistics and tracing](#object-statistics-and-tracing)
  - [Interoperability](#interoperability)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Presentation

### Introduction

QiMessaging is a network protocol used to build rich distributed
applications. It was created by Aldebaran Robotics (currently SoftBank
Robotics) and is the foundation of the NAOqi SDK. SoftBank Robotics
uses it inside Romeo, NAO and Pepper robots. QiMessaging is also used
in Choregraphe, an integrated development environment which helps
programming NAO and Pepper robots.

An open source implementation of QiMessaging is actively developed as
part of the libqi framework: https://github.com/aldebaran/libqi

Choregraphe can be downloaded
[here](https://developer.softbankrobotics.com/us-en/downloads/pepper).

### Overview

QiMessaging exposes a software bus on which services can be used.
Services have methods (to be called) and signals (to be watched).
Signals, method parameters can have different kind of type:

- basic type (int, float, string, bool, ...)
- aggregates (vector, map, structures, tuple)
- untyped (also referred as `value`)
- object (services are objects)
- other (pointer, raw, unknown, ...)

Services are objects registered to a naming service called the
service directory.

Objects send messages to communicate. Messages are composed of a
header and a payload. Different type of message allow different
interactions (some message types are "Call", "Reply", "Cancel",
"Event"). Those messages are transmitted using a transport protocol
defined by an string (such as: "tcp://localhost:9559"). Currently
libqi supports two transport protocols (TCP and SSL).

### Comparisons

How does QiMessaging compares with:

* **D-Bus**: Both QiMessaging and D-Bus allow introspection (i.e. it
  is possible to list of the services and their methods). Both have a
  concept of asynchronous notification. D-Bus have a well defined
  permissions system. QiMessaging doesn't includes permissions but has
  authentication. QiMessaging allows different applications to
  communicate across the network.

* **Binder**: Both Binder and QiMessaging follow an object model. And
  both Android and libqi provides tools to process an intermediate
  description of the service (IDL) and generate proxy object. In
  QiMessaging objects are first class citizens (a method can return an
  object).

* **Thrift**: Both Thrift and QiMessaging enable procedure call over
  the network. Thrift have a serialization format optimized to reduce
  the size of the packets. QiMessaging uses a simpler serialization
  format for discoverability.

* **Component Object Model (COM)**: Both COM and QiMessaging offer an
  IPC mechanism independent of a particular programming language which
  can be easily binded to various languages. libqi supports C++ and
  Python. NAOqi has JavaScript binding among many other binding. Both
  COM and QiMessaging can use reference counting to manage the life
  time of an object.

* **ROS**: Both ROS and QiMessaging offer a publisher/subscriber model
  as well as an RPC mecanism. Conceptually, the two solutions have a
  lot in common. ROS 1 uses XMLRPC (i.e. XML over HTTP) while
  QiMessaging use a binary format over TCP connections. QiMessaging
  does not timestamp messages as opposed to ROS.

### OSI Model

QiMessaging corresponds to the following layers of the OSI model:

- **Layer 4 (transport)**: QiMessaging does not specify a transport
  protocol and supports different address schemes (ex: `tcp://`, `tcps://`)

- **Layer 5 (session)**: QiMessaging includes a mandatory
  authentication procedure.

- **Layer 6 (presentation)**: QiMessaging uses two abstractions for
  presentation: a serialization format and a signature format
  describing the serialized types. Using signatures and the
  serialization format applications can exchange rich types.

- **Layer 7 (application)**: API are exposed using services (which are
  objects). Objects expose methods and signals.

## Messages

Every data shared with QiMessaging is encapsulated in a message.
Messages start with a header followed with an optional payload.

Example of messages are:

- A message of type `call` is sent to a service to initiate a remote
  procedure call.

- A service responds to a `call` message with a `reply` message.

- A message of type `error` is sent when a service is unable to
  respond a `call` message.

- A message of type `event` is sent when a signal state has changed.


![Example of type of messages](/doc/examples-message-type.png)


Messages are composed of two parts:

- a header of fixed size
- an optional payload

### Message Header

A header is composed of 28 bytes structured in the following way:

```
      0       1       2       3       4
      +-------------------------------+    ^
   0  |             Magic             |    |
      +-------------------------------+    |
   4  |           Message ID          |    |
      +-------------------------------+    |
   8  |           Data size           |    |
      +-------------------------------+    |
   12 |    Version    |  Type | Flags |  Header
      +-------------------------------+    |
   16 |           Service ID          |    |
      +-------------------------------+    |
   20 |           Object ID           |    |
      +-------------------------------+    |
   24 |           Action ID           |    |
      +-------------------------------+    V
   28 |                               |
      |        Data (Optional)        |
      |                               |
      +-------------------------------+

```

In the C programming language, this becomes:

```
struct header_t {
    uint32_t magic;      // constant value
    uint32_t message_id; // identifier to associate call/reply messages
    uint32_t size;       // size of the payload
    uint32_t version;    // protocol version (0)
    uint8_t  type;       // type of message (call, reply, event, ...)
    uint8_t  flags;      // flags
    uint32_t service_id; // service id
    uint32_t object_id;  // object id
    uint32_t action_id;  // function or event id
};
```

The header is also documented
[here](http://doc.aldebaran.com/libqi/design/network-binary-protocol.html)
as part of the libqi.

#### Magic

Messages start with a constant value of 4 bytes used to identify
QiMessaging headers. This value is always `0x42dead42` and it is
encoded in big endian.

#### Message ID

The message ID identify the transaction the message belongs to. A
`call` message is replied with a `reply` (or an `error`) message
having the same message ID as the `call` message. Transactions
composed of only one message (such as `event` and `post` messages)
have their own message ID.

#### Size

The size of the payload. Can be zero.

#### Version

The version of the protocol used. This document describes version 0.

#### Type

Each message have a type. Possible types are:

- **Unknown** (0): not used
- **Call** (1): Initiate a remote procedure call. The payload of the
  message is the parameters of the method. The object and the method
  are identified by the fields Service, Object and Action described
  below.
- **Reply** (2): Response to remote procedure call message. The payload
  has the returned type of the called method.
- **Error** (3): Signal an error. Can be used in response to a call
  message. Payload is a value (values are described below).
- **Post** (4): Call a method but without expecting an answer. Payload
  contains the arguments of the method.
- **Event** (5): Inform of a new signal state. Events are sent
  following a call to the registerEvent method. unregisterEvent stops
  the stream of events. The payload is the new value of the signal.
- **Capability** (6):
- **Cancel** (7): Request the interruption of the remote procedure
  call.
- **Cancelled** (8):

#### Flags

#### Service ID

Each service registered to the service directory (described below) is
given a unique identifier. The service identifier act as a namespace
for the objects associated with the service.

#### Object ID

The identifier of the destination object. This value is to be resolved
withing the namespace of the service. This means two objects which
belongs to different services can have the same object ID.

#### Action ID

The identifier of the method or the signal associated with the
message.

### Payload

The payload is optional. Its content depends on the type of message.

- `call` message: the payload contains the parameters of the method called.

- `reply` message: the payload contains the returned value of the
  call.

- `error` message: the payload contains a value associated with the
  error (like a string of character).

- `event` message: the payload contains the new value of the signal.

- `post` message: the payload contains the parameters of the method
  called.

In order to send a `call` message, one needs to know the prototype of
the method to be called. Since every method have a different
prototype each payload must be crafted accordingly.

For example the method `authenticate` of the service 0 takes a map of
string keys and values of type `value` and returns also a map of
string and `value`. On the other hand, the method `services` of
the service 1 takes no arguments and returns a list of a structure
called `ServiceInfo`. One needs to know what is structure is in order
read it.

Fortunately, every object (and service) can be introspected by calling
its method `metaObject`. Using this property, one can learn about the
methods (and their prototype). The prototype of a method is composed
of two parts:

- the type of the arguments
- the type of the returned value

Those types are described using a strings characters called
*signatures*. Understanding the format of the signature is required to
know what to read and write in the payload of a message. The next
section explains the signature format.


## Signatures

A signature is a string which describes a type. Signatures are used in
different places of the protocol. A signature can represents:

- the type of a signal

- the type of the arguments of a method

- the type of the data returned by a method call.

- the concrete type of a `value` when it is serialized (`value`s are
  detailed explained below).

### Types

libqi [documentation](http://doc.aldebaran.com/2-5/dev/libqi/api/cpp/type/signature.html) on the various type.

#### Basic types

- 'i': signed integer: 32 bits signed.
- 'I': unsigned integer: 32 bits signed.
- 'f': float: 32 bits (IEEE 754)
- 'd': double: 64 bits (IEEE 754)
- 'l': long: 64 bits signed value
- 'L': unsigned long: 64 bits signed value
- 'b': bool: boolean value
- 's': string: string of character
- 'r': raw data: array of bytes

When describing the return type of a method:
- 'v': void: the method is not returning a result other than the
  information of its successful completion.

#### Value

- 'm': value: a container type. It can contains any basic type.

#### Composite types

- '[' ... ']': vect: a sequence of elements of the same type.
- '{' ... ... '}': map: associative array from a key type to an element type.
- '(' ... ')<' ... .... '>': structure: a limited set of name element of different types.

#### Object

- 'o': an object. The signature of an object does not describes it
  methods or signal for two reasons:

  1- because object are passed by reference and not actually
  serialized.

  2- the list of methods and signals associated with an object can be
  query by calling its method "metaObject". See MetaObject section.

#### Other

- 'X': an unknown type.

### Examples

- i: an integer (`int`)
- f: a float (`float`)
- d: a double (`double`)
- [s]: a vector of string (`std::vector<string>`)
- [m]: a vector of value
- {lb}: a map from long to boolean (`std::map<long int, bool>`)
- (s)<Structure,field>: a structure represented C like:

```
struct Structure {
    string field;
};
```
- o: an object.
- [o]: a vector of objects.
- {s(s)<Struture,field>}: represented in C++ with::

```
struct Structure {
    string field;
};
std::map<string, struct Structure>
```


#### Real world examples

- `[(sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>]`: a vector of structure `ServiceInfo`.

- `(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>`:
  in the Go programming language this becomes:

```
type MetaMethodParameter struct {
    Name        string
    Description string
}
type MetaMethod struct {
    Uid                 uint32
    ReturnSignature     string
    Name                string
    ParametersSignature string
    Description         string
    Parameters          []MetaMethodParameter
    ReturnDescription   string
}
```

### Signature Grammar

```
type_string = "s"
type_boolan = "b"
type_int = "i"
type_uint = "I"
type_value = "m"
type_long = "l"
type_ulong = "L"
type_float = "f"
type_double = "d"
type_raw_data = "r"
type_object = "o"

type_basic = type_int | type_uint | type_string | type_float | type_double | type_long | type_ulong | type_raw_data | type_boolean | type_value | type_object

type_map = "{" type_declaration type_declaration "}"

type_vect = "[" type_declaration "]"

type_tuple = "(" type_declaration ")"

type_definition = type_tuple "<" name "," list_of_names ">"

type_declaration = type_basic | type_map | type_vect | type_definition

list_of_declarations = type_declaration | type_declaration list_of_declarations

name = alphanumeric

list_of_names = name | name "," list_of_names
```

## Serialization

**Note**: The magic value (0x42dead42) is written in big endian while
all values are transmitted in little endian.

### Basic types

- **integer**: 32 bits little endian signed (int32).
- **unsigned integer**: 32 bits little endian (uint32).
- **long**: 64 bits little endian signed (int64).
- **unsigned long**: 64 bits little endian (uint64).
- **float**: 32 bites IEEE 754 little endian (float32).
- **double**: 64 bites IEEE 754 little endian (float64).
- **boolean**: 1 byte, zero for false (bool)
- **string**: an integer (as defined above) followed the bytes of the
  string. Not finishing with a zero (str).
- **raw data**: an array of byte of a variable size

### Values

A value is serialized with a string (defined above) representing the
signature of the concrete type followed with the serialized value of
the type.

### Composite types (map, vect and struct)

- **vector**: an integer representing the size of the sequence followed
  with the concatenation of the serialized elements.

- **map**: an integer representing the size of the map followed with the
  concatenation of the pair of serialized key and value.

- **structure**: the concatenation of the serialized fields.

- **tuple**: the concatenation of the serialized members.

### Object

An object is a reference to a remote entity, therefore it is not
really serialized. What is serialized is the description of this
object.

This description contains the following fields:
- **bool**: unknown usage
- **MetaObject**: description of the object (explained later)
- **integer**: unknown usage
- **integer**: service id
- **integer**: object id

## Objects

An object is composed of:
- a set of method to be called
- a set of signals to be watched
- a set of properties to be queried

Within the libqi framework objects are called
[AnyObject](http://doc.aldebaran.com/2-5/dev/libqi/api/cpp/type/anyobject.html).

### Methods

At the heart of QiMessaging is the feature of remote procedure calls:

- authentication is done by calling the method `authenticate` to the
  service 0.

- asynchronous communication of events requires to subscribe to a
  signal using the method `registerEvent`.

- services name resolution is done using methods from the service
  directory (service 1).

- method and signal name resolution of an object is done by calling
  the method `metaObject` of that object.

#### Generic methods

Here is a list of methods shared by almost every object:

- 0: `fn registerEvent(objectID: uint32, signalID: uint32, handler:
  uint64) uint64`: subscribes to a signal. The new values of the
  signal will be sent the client using messages of type `event`.

- 1: `fn unregisterEvent(objectID: uint32, signalID: uint32, handler: uint64) void`: unsubscribes
  from a signal.

- 2: `metaObject(objectID: uint32) MetaObject`: introspects an object. It
  returns structure called `MetaObject` which describe an object. This
  structure contains the list of methods, signal and properties as
  well as the signature of the associated types. When communicating
  with an object, the method `metaObject` is often the first method
  called because it allows the client to associate the name of the
  method with the action ID.

- 3: `terminate(objectID: uint32) void`: informs a object it is not used
  anymore. This allows the object to be destroyed. It is used in the
  context of objects returned by methods. In such situation life cycle
  of the object is controlled by the client.

- 5: `property(Value) Value`: returns the value associated with the
  property.

- 6: `setProperty(Value,Value) void `: sets the value of a property.

- 7: `properties() Vect<str>`

- 8: `registerEventWithSignature(objectID: uint32, signalID: uint32, handler: uint64, signature: String) uint64`

Notice: one exception is the the object 0 of service 0 which does not
supports those methods.

The `MetaObject` is useful to interact with an object. Especially it
gives access to the list of methods and signals as well as their
signature. The next section describes this `MetaObject`.

#### Specific methods

Allong with the previously describe methods, objects can have as many
methods as needed. Those specific methods have indexes which start
from 100.

### Signals

A signal represents a variable whose state can change in time. When
the state of the signal change, notifications are sent. One can
subscribe to a signal using `registerEvent` and unsubscribe with
`unregisterEvent`. Once registered, a client will receive new values
taken by a signal by messages of type `event` which contain the new
value.

### Properties
### MetaObject

`MetaObject` is a structure which describes an object. The description
includes the list of the methods along with their parameters and
return type.

Since every object has a method which returns a MetaObect, every
instance of a MetaObect contains a description of the `MetaObject` type.

Here is the signature of `MetaObject`:

```
({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,r
eturnSignature,name,parametersSignature,description,parameters,returnDes
cription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,ui
d,name,signature>}s)<MetaObject,methods,signals,properties,description>
```

This expression is complicated because signature of `MetaObject`
embedded the signature of `MetaMethod` (which embedded
`MetaMethodParameter`) as well as `MetaSignal` and `MetaProperty`.

Translated in a programming language this becomes:

```
type MetaObject struct {
	Methods     map[uint32]MetaMethod
	Signals     map[uint32]MetaSignal
	Properties  map[uint32]MetaProperty
	Description string
}

type MetaMethod struct {
	Uid                 uint32
	ReturnSignature     string
	Name                string
	ParametersSignature string
	Description         string
	Parameters          []MetaMethodParameter
	ReturnDescription   string
}

type MetaMethodParameter struct {
	Name        string
	Description string
}

type MetaSignal struct {
	Uid       uint32
	Name      string
	Signature string
}

type MetaProperty struct {
	Uid       uint32
	Name      string
	Signature string
}
```


## Services

The list of the services on the bus can be query using a service
(service directory). This makes this bus discoverable. Each service
exposes the list of its methods, signals and properties.

### Service Server (ID 0)

Used for authentication.

### Service Directory (ID 1)

Used to register new services to the bus and to list the services.

#### Methods

Here is a list of the specific method offerred by the service
directory:

- 100: `service(str) ServiceInfo`: associate the name of a service
  with its identifier and its network endpoint.

- 101: `services() Vect<ServiceInfo>`: list all registered services.

- 102: `registerService(ServiceInfo) uint32`: add a new service to the
  service directory.

- 103: `unregisterService(serviceID: uint32) void`: remove the service.

- 104: `serviceReady(serviceID: uint32) void`: informs the service directory when
  a given service is ready to receive requests.

- 105: `updateServiceInfo(ServiceInfo) void`: update the service
  information associated with a service.

- 108: `machineId() str`: returns the unique identifier of the
  machine.

- 109: `_socketOfService(serviceID: uint32) Object`:

#### Signals

In order to monitor the list of services, one can use the following
signals:

- 106: `serviceAdded(serviceID: uint32, name: str)`: informs when a new service has
  registered to the bus.

- 107: `serviceRemoved(serviceID: uint32, name: str)`: informs when a service has
  quitted the bus.


### Example (LogManager)

## Networking

### Endpoints

The `ServiceInfo` data structure which describes a service contains a
list of addresses whre to contact the service.

### TCP

The default port to contact the service directory is 9559.

### SSL

SSL is supported using the addressing scheme `tcps://` in the
`ServiceInfo` description of the service.

## Authentication

Authentication is required to communicate with QiMessaging. The
service 0 is responsible for the authentication procedure and have a
method called `authenticate`. When a client tries to skip the
authentication procedure, it is sent a `capability` message.

### CapabilityMap

The method `authenticate` takes a map of string and `value` as
argument. This data structure is called the `CapabilityMap`:

```
type CapabilityMap map[string]value
```

The capability map is used to exchange information during the
authentication procedure.

#### Authentication state

The status of the authentication procedure is stored in this
capability map under the key:

- `"__qi_auth_state"`: integer `value`

Possible values are:

- `1`: Error: an error occurs during the authentication. Possible
  cause of error is an invalid credential.

- `2`: Continue: the server request the client to provide further
  information. Possible reason is to request an acknowledgment
  required for the client.

- `3`: Done: the authentication procedure is completed. The client can
  access the bus.

#### Authentication credentials

The authentication to a service may be required. If so, the client is
expected to complete the capability map with the following keys:

- `"auth_user"`: string `value`: the user (ex: `nao`, `tablet`)

- `"auth_token"`: string `value`: a password

If the token is correct, the authentication state pass to `3` (done)
else it become `1` (error).

It is possible, if the server has no password defined, to generate a
random password and replies the client with the key:

- `"auth_newToken"`: string `value`: the new password to use

In such case, the authentication state pass to `2` (continue) and the
client must authenticate again using the new password.

![Example of token generation](/doc/examples-token-generation.png)

#### Protocol feature negotiation

The list of capabilities is not fixed in the protocol. This allow a
client to announce the non mandatory feature of the protocol which are
supported. Possible values are:

- `"ClientServerSocket"`: boolean `value`

- `"MessageFlags"`: boolean `value`

- `"MetaObjectCache"`: boolean `value`

- `"ObjectPtrUID"`: boolean `value`

- `"RemoteCancelableCalls"`: boolean `value`

## Routing
### Message destination
### Message origin
### Message transferred

## Misc
### Object Statistics and tracing

There is a set of methods used for tracing (index ranging from 80 to
85):

- 80 `isStatsEnabled() bool`: returns true if the statistics are
  enabled.

- 81: `enableStats(bool) void`: enables statistics.

- 82: `stats() Map<uint32,MethodStatistics> `: returns the current
  statistics

- 83: `clearStats() void `: reset the counters.

- 84: `isTraceEnabled() bool`: returns true if tracing is enable.

- 85: `enableTrace(bool) void`: enables / disables tracing.

Most object have this signal:

- 86: `traceObject(EventTrace)`: signal when a method of the object
  have been called.


### Interoperability

