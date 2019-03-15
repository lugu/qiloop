# proxy generation tutorial

This guide will show you how to generate a proxy for ALVideoDevice.

## Requirements

A running instance of QiMessaging. For example a NAOqi running on your
desktop inside Choregraph or a Pepper or Nao robot.

In this tutorial we assume the address to be tcp://127.0.0.1:9559.
If the port 9559 is not open, try 9503 or 9443 using tcps.

## Installation

Refer to [the Go website](https://golang.org) for the installation of
Go.

Install qiloop with:

```
go get github.com/lugu/qiloop/...
```

## Create a package

```
mkdir -p $GOPATH/src/demo
cd $GOPATH/src/demo
```

## Generate the proxy

Two steps:
- contact the running instance of ALVideoDevice to generate an IDL
  file. This is done using the `qiloop scan` command.

- generate the specialized proxy of ALVideoDevice from the IDL file.
  This is done using the `qiloop proxy` command.

The following example generates a proxy for ALVideoDevice (`demo/video_proxy.go`):

```
$GOPATH/bin/qiloop scan --qi-url tcp://127.0.0.1:9559 --idl demo/video_device.idl --service ALVideoDevice
```

Add one line at the top of the IDL file video_device.idl to specify a
package name:
'package main`

Then generate the proxy with:
```
$GOPATH/bin/qiloop proxy --idl demo/video_device.idl --output video_proxy.go
```


## Proxy usage

Finally, the main program (`demo/main.go`) which uses the proxy of
ALVideoDevice to obtain a image from a camera:

```golang
package main

import (
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/type/value"
	"log"
	"os"
)

const (
	topCam = 0
	vga    = 1
	rgb    = 13
)

func main() {
	sess, err := session.NewSession("tcp://127.0.0.1:9559")
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	services := Services(sess)
	videoDevice, err := services.ALVideoDevice()
	if err != nil {
		log.Fatalf("failed to create video device: %s", err)
	}

	id, err := videoDevice.Subscribe("me", topCam, vga, rgb)
	if err != nil {
		log.Fatalf("failed to initialize camera: %s", err)
	}

	img, err := videoDevice.GetImageRemote(id)
	if err != nil {
		log.Fatalf("failed to retrieve image: %s", err)
	}
	values, ok := img.(value.ListValue)
	if !ok {
		log.Fatalf("invalid return type")
	}
	width := values[0].(value.IntValue).Value()
	heigh := values[1].(value.IntValue).Value()
	pixels := values[6].(value.RawValue).Value()

	log.Printf("resolution: %dx%d", width, heigh)
	file, err := os.Create("image.rgb")
	if err != nil {
		log.Fatalf("can not create file")
	}
	file.Write(pixels)
        file.Close()
}
```

Since a lot of code is generated, use `go doc` to browse the
documentation of ALVideoDevice:
```
go doc demo

```
