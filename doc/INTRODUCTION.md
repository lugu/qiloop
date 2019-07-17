# Introduction to QiMessaging

## Scope

This documents describes various aspects of the protocol QiMessaging.
It is destined to software engineers and system designers who want to
learn about the inner working of the protocol.

## Disclaimer

QiMessaging is the word historically used to describe the network
protocol of `libqi`. It have never found a formal definition nor a
clear specification. This document reflects the view of its author and
has no authoritative value.

## Overview

QiMessaging is a network protocol used to build rich distributed
applications.

In this context, rich means it allows developers to express
abstractions using an object oriented representation.

Distributed means the various part of the application can be located
on different nodes of a network. A naming service is included to
locate the various parts of the application.

## Usage

While QiMessaging was created in the context of a robotic application,
the problem it solves is not specific to this field.

QiMessaging provides an RPC mechanism. Services (objects) expose
methods. Each **method** has a **signature**: a name and a description of
the types of its parameters and returned value.

QiMessaging follows the publish-subscribe pattern. Publishers don't
need to know the subscribers. Topics are organized by object to which
they belong. Stateful topics are called **properties** and stateless
events are called **signals** in QiMessaging.

QiMessaging borrows some ideas from the actor model. Mainly, that an
actor is the main unit of computation. From an RPC and pub-sub
perspective, objects (and services) are first class citizens: one can
receive and send reference to objects. Objects exchange **messages**
of various types (call, post, cancel, event ...).

The rest of this document explores the relation between QiMessaging
on those models.

## Origin

QiMessaging was created at SoftBank robotics (formerly Aldebaran) with
the following concepts in mind:

- *Distribution and network discovery*: a robot can be composed of
  various component connected via a network. The components need to
  interact to fulfill the robot mission.

For example, the Pepper robot contains two CPUs: one located in the
head which controls the body, and one in the tablet to offer a
graphical interface. The robot can also be remotely controlled from a
desktop using the Choregraph IDE from SoftBank.

- *Reactive programming and data stream*: the robot evolves in an ever
  changing environment to which it must respond.

For example, the NAO robot can be pushed and loose balance anytime: in
which case it must quickly react. It also can detect and track faces
using its cameras.

- *Typed value and dynamic versioning*: the values involved in the APIs
  are typed, but the supported APIs are queried at runtime.

For example, between two software releases, deprecated API can be
removed. Since the life-cycle of the various components is not
identical. Being able to query the available API helps both forward
and backward compatibility.

## Interface Description Language

QiMessaging services are described using IDL files. Here is an
example: the service directory IDL:

        package directory

        interface ServiceDirectory
                fn service(name: str) -> ServiceInfo
                fn services() -> Vec<ServiceInfo>
                fn registerService(info: ServiceInfo) -> uint32
                fn unregisterService(serviceID: uint32)
                fn serviceReady(serviceID: uint32)
                fn updateServiceInfo(info: ServiceInfo)
                fn machineId() -> str
                sig serviceAdded(serviceID: uint32, name: str)
                sig serviceRemoved(serviceID: uint32, name: str)
        end

        struct ServiceInfo
                name: str
                serviceId: uint32
                machineId: str
                processId: uint32
                endpoints: Vec<str>
                sessionId: str
        end

It describe an object type (ServiceDirectory) which has methods
(service, services, registerService, ...) and signals (serviceAdded,
serviceRemoved). It also describe the `ServiceInfo` type used to
represent the information of a service.

## Generic object

As we have seen with the `ServiceDirectory` IDL, objects can have
methods (`fn`) and signals (`sig`). A method represents a RPC. A
signal represents a stateless event. Objects can also have properties
(`prop`) to represent observable states.

Every object have a set of methods not represented in the
`ServiceDirectory` example. Here is the common methods to every
object:

        interface Object
                fn registerEvent(objectID: uint32, actionID: uint32, handler: uint64) -> uint64
                fn unregisterEvent(objectID: uint32, actionID: uint32, handler: uint64)
                fn metaObject(objectID: uint32) -> MetaObject
                fn terminate(objectID: uint32)
                fn property(name: any) -> any
                fn setProperty(name: any, value: any)
                fn properties() -> Vec<str>
        end

The `registerEvent` and `unregisterEvent` are used to subscribe to a
signal or a properties. Once called, the object will send an event
message every time the property or the signal is updated.

The `property`, `setProperty` and `properties` methods are used to
access the state of a property.

The `metaObject` method returns the description of all the methods,
signals and properties supported by the object. It is used to
introspect an object. Here is the description of the MetaObject type:

        struct MetaObject
                methods: Map<uint32,MetaMethod>
                signals: Map<uint32,MetaSignal>
                properties: Map<uint32,MetaProperty>
                description: str
        end
        struct MetaSignal
                uid: uint32
                name: str
                signature: str
        end
        struct MetaProperty
                uid: uint32
                name: str
                signature: str
        end
        struct MetaMethod
                uid: uint32
                returnSignature: str
                name: str
                parametersSignature: str
                description: str
                parameters: Vec<MetaMethodParameter>
                returnDescription: str
        end
        struct MetaMethodParameter
                name: str
                description: str
        end

## logs
## directory
## wireshark
## stat
## trace
## syntax idl
## authentication
## gateway
## transport / scheme
## signature
## codec
## message origin
