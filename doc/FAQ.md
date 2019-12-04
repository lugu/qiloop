# Frequently Asked Questions

## What is QiLoop ?

QiLoop is a Go package that implements QiMessaging, the network
protocol used to interact the robots NAO and Pepper.

## What is QiMessaging ?

QiMessaging is a network protocol which exposes a bus of services.
Those services are described with an IDL file. Think of it as a
mixture between d-bus and gRPC.

Pepper and NAO APIs are accessed through QiMessaging using NAOqi SDK
and Pepper SDK. Visit the official documentation for a description of
those APIs.

Here is an [description of QiMessaging](https://github.com/lugu/qiloop/blob/master/doc/about-qimessaging.md).

For a reference implementation, visit:
[https://github.com/aldebaran/libqi](https://github.com/aldebaran/libqi).

QiLoop provides the way to access this software bus in order to create
applications for Pepper of NAO in Go.

## What about libqi-go ?

libqi-go is Go's binding for libqi. It can be found here:
https://github.com/aldebaran/libqi-go

QiLoop is a pure Go re-implementation of QiMessaging and does not use
libqi-go.

## Why a re-implementation of QiMessaging is Go ?

The original motivation was to learn about QiMessaging and explore
what is possible with the protocol.

Some benefits of using QiLoop:

  - easy to understand and to debug
  - easy to test in isolation
  - easy to extend (contributions are welcome!)

Not to forget Go's inherent advantages:

  - portable (Linux, MacOS, Windows, Android and iOS)
  - easy to cross-compile (external cross toolchain)
  - easy to deploy (static binaries)

## Why should I use it ?

Probably not: you should not use QiLoop for anything critical.
SoftBank provides complete tools and IDE to program their robots.

This said, if you are a Go enthusiast and you are not afraid to
contribute: give QiLoop a try!

## Which robot is supported: Pepper, NAOv5, NAOv6 ?

The aim is to support Pepper, NAOv5 and NAOv6 if the protocol permit
it. So far, it seems possible. In terms of API, each robot
requires its own IDL file. So far those IDL files are not shipped as
part of QiLoop.

## Which programming languages are supported?

Only Go. Notice libqi have bindings in various languages.

## How do I start ?

You need a robot: either a physical robot (Pepper or NAO) or a virtual
robot (launched by Choregraphe).

Download qiloop:

        $ go get github.com/lugu/qiloop/cmd/qiloop

If you have a robot, try the `qiloop` command line tool to list and
scan the services.

        $ qiloop info --qi-url tcp://<ROBOT_IP>:9559

Or

        $ qiloop info --qi-url tcps://<ROBOT_IP>:9503

From there, take a look at the [examples
directory](https://github.com/lugu/qiloop/tree/master/examples).

## How to use QiLoop ?

The basic workflow is often the same:

 1. Use the program `qiloop scan` to generate an IDL file
    corresponding to the service your application targets.

 2. Use `qiloop proxy` to generate the Go code which gives access to
    this service.

 3. Code, test, debug and enjoy!

## Where is the documentation ?

You can find it at GoDoc:
[http://godoc.org/github.com/lugu/qiloop](http://godoc.org/github.com/lugu/qiloop)

## Can I use it for Android/iOS applications ?

At this stage, QiLoop is best suited for tooling or mobile
applications. See [QiTop](https://github.com/lugu/qitop) for an
example of application made with QiLoop.

Thanks to `gomobile` QiLoop is portable to Android (and maybe iOS?).
Take a look at [QiView](https://github.com/lugu/qiview) for a crude
mobile application.

## How can I submit a bug or a patch ?

Please open an issue or a pull request on Github.

## Why is it called QiLoop ?

The name refers to a meta description observed when analyzing the
binary protocol: Services have a method (called `metaObject`) which
returns the description (signature) of the available methods,
including the method used to describe the methods.

One can reverse engineer the protocol with zero knowledge!

## Where to discuss it ?

Please open an issue with your question on Github.

## My question is not answered here, what should I do ?

Please open an issue with your question on Github.
