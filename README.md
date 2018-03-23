[![Build Status](https://travis-ci.org/lugu/qiloop.svg?branch=master)](https://travis-ci.org/lugu/qiloop) [![Go Report Card](https://goreportcard.com/badge/github.com/lugu/qiloop)](https://goreportcard.com/report/github.com/lugu/qiloop) [![codecov](https://codecov.io/gh/lugu/qiloop/branch/master/graph/badge.svg)](https://codecov.io/gh/lugu/qiloop) [![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg)](https://github.com/emersion/stability-badges#experimental)

# qiloop

This is work in progress, you have been warned.

`qiloop` is an implementation of the QiMessaging protocol written in the
Go programming language.

For more information about QiMessaging, here are some notes [about QiMessaging](https://github.com/lugu/qiloop/blob/master/doc/NOTES.md).

Installation
------------

Please install [Go 1.10](https://golang.org/dl/).

```
go get -d github.com/lugu/qiloop/...
go generate github.com/lugu/qiloop/type/object
go generate github.com/lugu/qiloop/meta/stage2
```

Proxy generation
----------------

By default, `qiloop` comes with only two generated services:
ServiceManager and LogManager (package
github.com/lugu/qiloop/bus/services).

If you want to generate more proxies, you will need access to a
running QiMessage server. Here is how to generate proxies given a
running server on port 9559:

```
go get github.com/lugu/qiloop/meta/stage3/cmd/stage3
$GOPATH/bin/stage3 interfaces.go services.go localhost:9559
```

This will create `interface.go` and `services.go` which can be used to
access the service on the bus.

Examples
--------

- [info
  demo](https://github.com/lugu/qiloop/blob/master/bus/cmd/info/main.go)
  illustrates how to call a service: it lists the services registered
  to the service directory.


- [signal
  demo](https://github.com/lugu/qiloop/blob/master/bus/services/demo/cmd/signal/main.go)
  illustrates how to subscribe to a signal: it prints a log each time
  a service is removed from the service directory.

Status
------

The project is under rapid development. Currently, the client part if
mostly working (except the properties which are not yet supported). So
one should be able to use qiloop to call a service and subscribe to a
signal. Don't expect more than this.

What is working:

- TCP connection
- Proxy generation
- Call and signals
- IDL generation

What is under development:

- IDL parsing

What is yet to be done:

- Service stub generation
- Service properties
- SSL transport
- authentication
