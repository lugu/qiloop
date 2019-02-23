# qiloop

[![Build Status](https://travis-ci.org/lugu/qiloop.svg?branch=master)](https://travis-ci.org/lugu/qiloop)
[![CircleCI](https://circleci.com/gh/lugu/qiloop/tree/master.svg?style=shield)](https://circleci.com/gh/lugu/qiloop/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lugu/qiloop)](https://goreportcard.com/report/github.com/lugu/qiloop)
[![codecov](https://codecov.io/gh/lugu/qiloop/branch/master/graph/badge.svg)](https://codecov.io/gh/lugu/qiloop)
[![Test Coverage](https://api.codeclimate.com/v1/badges/b192466a26dbced44274/test_coverage)](https://codeclimate.com/github/lugu/qiloop/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/b192466a26dbced44274/maintainability)](https://codeclimate.com/github/lugu/qiloop/maintainability)
[![stability-unstable](https://img.shields.io/badge/stability-unstable-yellow.svg)](https://github.com/emersion/stability-badges#unstable)
[![Waffle.io - Columns and their card count](https://badge.waffle.io/lugu/qiloop.svg?columns=To%20Do,In%20Progress)](https://waffle.io/lugu/qiloop)

**`qiloop`** is an implementation of QiMessaging written in [Go](https://golang.org).

QiMessaging is a network protocol used to build rich distributed
applications. It was created by Aldebaran Robotics (currently
[SoftBank Robotics](https://www.softbankrobotics.com/emea/en/index))
and is the foundation of the NAOqi SDK. For more details about
QiMessaging, visit this [analysis of the
protocol](https://github.com/lugu/qiloop/blob/master/doc/NOTES.md).

## Installation

    go get github.com/lugu/qiloop/...

## Demo

Here is how to connect to a server and list the running services:

```golang
package main

import (
	"github.com/lugu/qiloop/bus/client/services"
	"github.com/lugu/qiloop/bus/session"
)

func main() {
	sess, err := session.NewSession("tcp://localhost:9559")
	if err != nil {
		panic(err)
	}

	proxies := services.Services(sess)

	directory, err := proxies.ServiceDirectory()
	if err != nil {
		panic(err)
	}

	serviceList, err := directory.Services()
	if err != nil {
		panic(err)
	}

	for _, info := range serviceList {
		println("service " + info.Name)
	}
}
```

## Proxy generation tutorial

By default, `qiloop` comes with two proxies: ServiceDirectory and
LogManager.

Follow [this tutorial](https://github.com/lugu/qiloop/blob/master/doc/TUTORIAL.md) to generate more proxy.

## Examples

-   [method call](https://github.com/lugu/qiloop/blob/master/bus/client/services/demo/cmd/method/main.go)
    illustrates how to call a method of a service: this example lists
    the services registered to the service directory.


-   [signal registration](https://github.com/lugu/qiloop/blob/master/bus/client/services/demo/cmd/signal/main.go)
    illustrates how to subscribe to a signal: this example prints a
    log each time a service is added to the service directory.

## Authentication

If you need to provide a login and a password to authenticate yourself
to a server, create a file `$HOME/.qi-auth.conf` with you login on the
first line and your password on the second.

## Status

This is work in progress, you have been warned.

The client and the server side is working: one can implement a service
from an IDL and generate a specialized proxy for this service.
A service directory is implemented as part of the standalone server.

What is working:

-   TCP and TLS connections
-   client proxy generation
-   service stub generation
-   method, signals and properties
-   Authentication
-   IDL parsing and generation
