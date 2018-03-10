# About QiMessaging

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Presentation](#presentation)
  - [Introduction](#introduction)
  - [Overview](#overview)
  - [Comparisons](#comparisons)
  - [OSI Model](#osi-model)
- [Message](#message)
  - [Message Header](#message-header)
    - [Magic number](#magic-number)
    - [Size](#size)
    - [Header Version](#header-version)
    - [Header Type](#header-type)
    - [Flags](#flags)
    - [Service ID](#service-id)
    - [Object ID](#object-id)
    - [Action ID](#action-id)
  - [Payload](#payload)
- [Signatures](#signatures)
  - [Types](#types)
    - [Basic types](#basic-types)
    - [Composite types](#composite-types)
    - [Object](#object)
  - [Examples](#examples)
    - [Real world examples](#real-world-examples)
  - [Signature Grammar](#signature-grammar)
- [Serialization](#serialization)
  - [Basic types](#basic-types-1)
  - [Values](#values)
  - [Composite types (map, list and struct)](#composite-types-map-list-and-struct)
  - [Object](#object-1)
- [Objects](#objects)
  - [Methods](#methods)
  - [Signaux](#signaux)
  - [Properties](#properties)
  - [MetaObject](#metaobject)
- [Services](#services)
  - [Service Server (ID 0)](#service-server-id-0)
  - [Service Directory (ID 1)](#service-directory-id-1)
  - [Example (LogManager)](#example-logmanager)
- [Networking](#networking)
  - [Endpoints](#endpoints)
  - [TCP](#tcp)
  - [SSL](#ssl)
- [Authentication](#authentication)
  - [CapabilityMap](#capabilitymap)
  - [Token](#token)
- [Routing](#routing)
  - [Message destination](#message-destination)
  - [Message origin](#message-origin)
  - [Message transferred](#message-transferred)
- [Misc](#misc)
  - [Object Statistics](#object-statistics)
  - [Object tracing](#object-tracing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Presentation

### Introduction

QiMessaging is a network protocol used to create rich distributed
applications. It was created by Aldebaran Robotics (currently SoftBank
Robotics) and is the foundation of the NAOqi SDK. SoftBank Robotics
uses it for the Romeo, NAO and Pepper robots. QiMessaging is also used
in Choregraphe, an integrated development environment which simplifies
programming of the NAO and Pepper robots.

An open source implementation of QiMessaging is actively developed as
part of the libqi framework: https://github.com/aldebaran/libqi

Choregraphe can be downloaded
[here](https://developer.softbankrobotics.com/us-en/downloads/pepper).

### Overview

QiMessaging exposes a software bus on which services can be used.
Services have methods (to be called) and signals (to be watched).
Signals, method parameters can have different kind of type:

- basic type (int, long, float, string and bool)
- aggregates (vector, hash map, structures)
- untyped (also referred as values)
- and even objects!

Actually services are plain objects registerred to a naming service
(service directory).

Objects exchange messages to communicate. Messages are composed of a
header and a payload. Different type of message allow different
interactions (some message types are "Call", "Reply", "Cancel",
"Event"). Those messages are transmitted using a transport protocol
defined by an string (such as: "tcp://localhost:9559"). Currently
libqi supports two transport protocols (TCP and SSL). 

FIXME: is QiMessaging a distributed system ? is there anything in the
protocol which prevent several services directory to exists
concurrently on the same bus ?

### Comparisons

How does QiMessaging compares with:

* **D-Bus**: Both QiMessaging and D-Bus allow introspection (i.e. it
  is possible to list of the services and their methods). Both have a
  concept of asynchroneous notification. D-Bus have a well defined
  permissions system. QiMessaging doesn't includes permissions but has
  authentication. QiMessaging allows different applications to
  communicate accross the network.

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
  IPC mechanism independant of a particular programming language which
  can be easily binded into various languages. libqi supports C++ and
  Python. NAOqi has JavaScript binding among many other binding. Both
  COM and QiMessaging can use reference counting to manage the life
  time of an object.

// TODO: ROS comparison

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

## Message

Objects exchange messages. A remote procedure call is initiate by
sending a "Call" message to an object. This object responds either
with an "Error" or a "Reply" message.

To publish a signal, an object send a "Event" message to the list of
object which have subscribed to the signal.

Messages are composed of two parts:

- a fixed size header describing the message.
- an optional payload containing serialized data.

### Message Header

The header is documented
[here](http://doc.aldebaran.com/libqi/design/network-binary-protocol.html).
Its size is 28 bytes and it is followed with an optional payload.

A message can be represented in Go with:

```
type Header struct {
    Id      uint32 // an identifier to match call/reply messages
    Size    uint32 // size of the payload
    Version uint16 // protocol version (0)
    Type    uint8  // type of the message
    Flags   uint8  // flags
    Service uint32 // service id
    Object  uint32 // object id
    Action  uint32 // function or event id
}

type Message struct {
    Header Header
    Payload []byte
}
```


#### Magic number

Messages start with a constant value of 4 bytes used to identify
QiMessaging headers. This value is always `0x42dead42` and it is
encoded in big endian.

#### Size

The size of the payload. Can be zero.

#### Header Version

The version of the protocol used. This document describes version 0.

#### Header Type

Each message have a type. Possible types are:

- **Unknown** (0): not used
- **Call** (1): Initiate a remote procedure call. The payload of the
  message is the parameters of the method. The object and the method
  are identified by the fields Service, Object and Action described
  bellow.
- **Reply** (2): Response to remote procedure call message. The payload
  has the returned type of the called method.
- **Error** (3): Signal an error. Can be used in response to a call
  message. Payload is a string value (values are described bellow).
- **Post** (4): // TODO
- **Event** (5): // TODO
- **Capability** (6): // TODO
- **Cancel** (7): // TODO
- **Cancelled** (8): // TODO

#### Flags

// TODO

#### Service ID

Each object registered to the service directory is given a unique
identifier. This identifier is used to address the objects associated
with this service.

#### Object ID

The identifier of the destination object. This value is to be resolved
withing the namespace of the service. This means two objects which
belongs to different services can have the same object ID.

#### Action ID

The identifier of the method or the signal associated with the
message.

### Payload

The payload is optional. Its content depends on the type of message.

- For example, in the case of a Call message, the payload is the
  binary serialization of the argument of the method to be called.
- In the case of a reply message, the payload is the returned value.
- In the case of an error message, the payload is a value. This value
  can be string a describing the error but does not have to be.

## Signatures

A signature is a string representing a type. Signatures are used in
different places of the protocol. It can represents:

- the type of a signal
- the type of the arguments to a method as well as the return type of
  the method.
- the concrete type of a value when it is serialized.

### Types

libqi documentation: http://doc.aldebaran.com/2-5/dev/libqi/api/cpp/type/signature.html

#### Basic types

- 'i' or 'I': integer: 32 bits signed.
- 'f': float: 32 bits (IEEE 754)
- 'l': long: 64 bits signed value
- 'b': bool: boolean value
- 's': string: string of character

// FIXME: is the string encoding fixed?

When describing the return type of a method:
- 'v': void: the method is not returning a result other than the
  information of its successful completion.

#### Composite types

- '[' ... ']': list: a vector of elements of the same type.
- '{' ... ... '}': map: associative array from a key type to an element type.
- '(' ... ')<' ... .... '>': structure: a limited set of name element of different types.

#### Object

- 'o': an object. The signature of an object does not describes it
  methods or signal for two reasons:

  1- because object are passed by reference and not actually
  serialized.

  2- the list of methods and signals associated with an object can be
  query by calling its method "metaObject". See MetaObject section.

### Examples

- i: an integer (`int`)
- f: a float (`float`)
- [s]: a vector of string (`std::vector<string>`)
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

- `[(sIsI[s]s)<ServiceInfo,name,serviceId,machineId,processId,endpoints,sessionId>]`: a vector of structure service Info.

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
type_int = "i" | "I"
type_value = "m"
type_long = "L"
type_float = "f"
type_object = "o"

type_basic = type_int | type_string | type_float | type_long | type_boolean | type_value | type_object

type_map = "{" type_declaration type_declaration "}"

type_list = "[" type_declaration "]"

list_of_declarations = type_declaration | type_declaration list_of_declarations

list_of_types = "(" list_of_declarations ")"

type_definition = list_of_types "<" name "," list_of_names ">"

type_declaration = type_basic | type_map | type_list | type_definition

name = alphanumeric

list_of_names = name | name "," list_of_names
```

## Serialization

**Note**: The magic value (0x42dead42) is written in big endian while
all values are transmitted in little endian.

### Basic types

- **integer**: 32 bits little endian signed.
- **long**: 64 bits little endian signed.
- **float**: IEEE 754 little endian.
- **boolean**: 1 byte (zero for false)
- **string**: an integer (as defined above) followed the bytes of the
  string. Not finishing with a zero.

### Values

A value is serialized with a string (a defined bellow) representing
the signature of the concrete type followed with the serialized value
of the type.

### Composite types (map, list and struct)

- **list**: an integer representing the size of the list followed with the
  concatenation of the serialized elements.

- **map**: an integer representing the size of the map followed with the
  concatenation of the pair of serialized key and value.

- **structure**: the concatenation of the serialized fields.

### Object

An object is a reference to a remote entity, therefore it is not
really serialized. What is serialized is the description of this
object.

This description contains the following fields:
- **bool**: unknown usage // FIXME: often true. Linked with terminate ?
- **MetaObject**: description of the object
- **integer**: unknown usage // FIXME: often 1.
- **integer**: service id
- **integer**: object id

## Objects

A service is associated with every object as well as a unique identifier
with the namespace of the service.

An object is composed of:
- a list of method to be called
- a list of signals to be watched
- a list of properties to be queried

libqi documentation: http://doc.aldebaran.com/2-5/dev/libqi/api/cpp/type/anyobject.html

### Methods
### Signaux
### Properties
### MetaObject

A MetaObject is a structure which describes an object. The description
includes the list of the methods along with their parameters and
return type.

Since every object has a method which returns a MetaObect, every
instance of a MetaObect contains a description of the MetaObject type.

## Services

The list of the services on the bus can be query using a service
(service directory). This makes this bus discoverable. Each service
exposes the list of its methods, signals and properties.

### Service Server (ID 0)

Used for authentication.

### Service Directory (ID 1)

Used to register new services to the bus and to list the services.

Description a service directory:

```
$ qicli info --hidden 1
001 [ServiceDirectory]
  * Info:
   machine   6126ad3c-2f1f-4e25-8ec9-8bb20bfd195e
   process   3082
   endpoints tcp://127.0.0.1:9559
  * Methods:
   000 registerEvent              UInt64 (UInt32,UInt32,UInt64)
   001 unregisterEvent            Void (UInt32,UInt32,UInt64)
   002 metaObject                 MetaObject (UInt32)
   003 terminate                  Void (UInt32)
   005 property                   Value (Value)
   006 setProperty                Void (Value,Value)
   007 properties                 List<String> ()
   008 registerEventWithSignature UInt64 (UInt32,UInt32,UInt64,String)

   080 isStatsEnabled             Bool ()
   081 enableStats                Void (Bool)
   082 stats                      Map<UInt32,MethodStatistics> ()
   083 clearStats                 Void ()
   084 isTraceEnabled             Bool ()
   085 enableTrace                Void (Bool)
   100 service                    ServiceInfo (String)
   101 services                   List<ServiceInfo> ()
   102 registerService            UInt32 (ServiceInfo)
   103 unregisterService          Void (UInt32)
   104 serviceReady               Void (UInt32)
   105 updateServiceInfo          Void (ServiceInfo)
   108 machineId                  String ()
   109 _socketOfService           Object (UInt32)
  * Signals:
   086 traceObject    (EventTrace)
   106 serviceAdded   (UInt32,String)
   107 serviceRemoved (UInt32,String)
```

### Example (LogManager)

## Networking
### Endpoints
### TCP
### SSL

## Authentication
### CapabilityMap
### Token

## Routing
### Message destination
### Message origin
### Message transferred

## Misc
### Object Statistics
### Object tracing
### Interroperability

