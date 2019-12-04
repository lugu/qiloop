![qiloop](https://github.com/lugu/qiloop/blob/master/doc/logo.jpg)

[![Build Status](https://travis-ci.org/lugu/qiloop.svg?branch=master)](https://travis-ci.org/lugu/qiloop)
[![CircleCI](https://circleci.com/gh/lugu/qiloop/tree/master.svg?style=shield)](https://circleci.com/gh/lugu/qiloop)
[![Go Report Card](https://goreportcard.com/badge/github.com/lugu/qiloop)](https://goreportcard.com/report/github.com/lugu/qiloop)
[![codecov](https://codecov.io/gh/lugu/qiloop/branch/master/graph/badge.svg)](https://codecov.io/gh/lugu/qiloop)
[![Test Coverage](https://api.codeclimate.com/v1/badges/b192466a26dbced44274/test_coverage)](https://codeclimate.com/github/lugu/qiloop/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/b192466a26dbced44274/maintainability)](https://codeclimate.com/github/lugu/qiloop/maintainability)
[![Documentation](https://godoc.org/github.com/lugu/qiloop?status.svg)](http://godoc.org/github.com/lugu/qiloop)
[![license](https://img.shields.io/github/license/lugu/qiloop.svg?maxAge=2592000)](https://github.com/lugu/qiloop/blob/master/LICENSE)
[![stability-unstable](https://img.shields.io/badge/stability-unstable-yellow.svg)](https://github.com/emersion/stability-badges#unstable)
[![Release](https://img.shields.io/github/tag/lugu/qiloop.svg)](https://github.com/lugu/qiloop/releases)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [About](#about)
- [Status](#status)
- [Usage](#usage)
- [Go API](#go-api)
- [Tutorials](#tutorials)
- [Examples](#examples)
- [Command line interface](#command-line-interface)
- [Authentication](#authentication)
- [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## About

QiMessaging is a network protocol used to build rich distributed applications.
It is created by [SoftBank Robotics](https://www.softbankrobotics.com/emea/en/index)
and is the foundation of the [NAOqi SDK](http://doc.aldebaran.com/2-8/) and
the [Pepper SDK](https://qisdk.softbankrobotics.com/).

libqi is the implementation of QiMessaging used in Pepper and NAO.
It is open-source and developped here:
[https://github.com/aldebaran/libqi](https://github.com/aldebaran/libqi).

QiLoop is another implementation of QiMessaging. It has two main goals:
  - being compatible with libqi
  - being a platform for experimenting with the protocol

Disclaimer: QiLoop is not affiliated with SoftBank Robotics.

## Frequently asked questions

Read the [FAQ](https://github.com/lugu/qiloop/blob/master/doc/FAQ.md).

## Status

Client and server sides are functional.

Service directory and log manager are implemented as part of the
standalone server (launched with `qiloop server`).

Features:

  - type supported: object, struct, values, map, list, tuple

  - actions: method, signals and properties are fully supported

  - cancel: client support only (see motion example)

  - transport: TCP, TLS, UNIX socket and QUIC (experimental)

  - authentication: read the credentials from `$HOME/.qiloop-auth.conf`

  - service introspection: generate IDL from a running instance (use `qiloop scan`)

  - IDL files: generate specialized proxy and service stub (use `qiloop stub`)

  - stats and trace support

## Usage

QiMessaging exposes a software bus to interract with services. Services have
methods (to be called), signals (to be watched) and properties (signals with
state). A naming service (the service directory) is used to discover and
register services. For a detailed description of the protocol, please visit
this [analysis of
QiMessaging](https://github.com/lugu/qiloop/blob/master/doc/about-qimessaging.md).

To connect to a service, a Session object is required: it represents the
connection to the service directory. Several transport protocols are supported
(TCP, TLS and UNIX socket).

With a session, one can request a proxy object representing a remote service.
The proxy object contains the helper methods needed to make the remote calls
and to handle the incomming signal notifications.

Services have methods, signals and properties which are described in an IDL
(Interface Description Language) format. This IDL file is process by the
`qiloop` command to generate the Go code which allow remote access to the
service (i.e. the proxy object).

## Go API

Installation:

    go get -u github.com/lugu/qiloop/...

Documentation: [http://godoc.org/github.com/lugu/qiloop](http://godoc.org/github.com/lugu/qiloop)

## Tutorials

  - How to create a proxy to an existing service: follow the [ALVideoDevice tutorial](https://github.com/lugu/qiloop/blob/master/doc/tutorial-videodevice.md).

  - How to create your own service: follow the [clock tutorial](https://github.com/lugu/qiloop/blob/master/doc/tutorial-clock.md).

## Examples

Basic examples:

  - [hello world](https://github.com/lugu/qiloop/blob/master/examples/say)
    illustrates how to call a method of a service: this example calls
    the method 'say' of a text to speech service.

  - [signal registration](https://github.com/lugu/qiloop/blob/master/examples/signal)
    illustrates how to subscribe to a signal: this example prints a
    log each time a service is added to the service directory.

Examples for NAO and Pepper:

  - [animated say](https://github.com/lugu/qiloop/blob/master/examples/animated-say)
    uses ALAnimatedSpeech to animate the robot.

  - [posture](https://github.com/lugu/qiloop/blob/master/examples/posture)
    puts the robot in a random position.

  - [motion](https://github.com/lugu/qiloop/blob/master/examples/motion)
    move the robot forward and demonstrate how to cancel a call.

  - [memory](https://github.com/lugu/qiloop/blob/master/examples/memory)
    uses ALMemory to react on a touch event.

Examples of service implementation:

  - [ping pong service](https://github.com/lugu/qiloop/blob/master/examples/pong)
    illustrates how to implement a service.

  - [space service](https://github.com/lugu/qiloop/blob/master/examples/space)
    illustrates the client side objects creation.

  - [clock service](https://github.com/lugu/qiloop/blob/master/examples/clock)
    completed version of the clock tutorial.

## Command line interface

```
    $ qiloop -h                                                                                                                                                            /home/ludo/qiloop
    qiloop - an utility to explore QiMessaging

	 ___T_
	| 6=6 |
	|__`__|
     .-._/___\_.-.
     ;   \___/   ;
	 ]| |[
	[_| |_]

      Usage:
	qiloop [info|log|scan|proxy|stub|server|trace]

      Subcommands:
	info - Connect a server and display services info
	log - Connect a server and prints logs
	scan - Connect a server and introspect a service to generate an IDL file
	proxy - Parse an IDL file and generate the specialized proxy code
	stub - Parse an IDL file and generate the specialized server code
	server - Starts a service directory and a log manager
	trace - Connect a server and traces services

      Flags:
	   --version  Displays the program version string.
	-h --help  Displays help with available flag, subcommand, and positional value parameters.
```

## Authentication

If you need to provide a login and a password to authenticate yourself
to a server, create a file `$HOME/.qiloop-auth.conf` with you login on the
first line and your password on the second.

## Contributing

 1. Fork me
 2. Create your feature branch
 3. Make changes (hopefully with tests and, why not, with documentation)
 4. Create new pull request
