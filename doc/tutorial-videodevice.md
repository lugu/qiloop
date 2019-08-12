# How to write a client application

This guide will show you how to write a client application which
connects to a service.

To illustrate this, we will write an application which connects to
the service ALVideoDevice and request an image from a camera. It is based on
the [video example](https://github.com/lugu/qiloop/blob/master/examples/video).

## Requirements

We need a running instance of the NAOqi SDK. It can either be on a
real robot (Pepper or NAO) or a virtual robot running on your desktop
(ex: using Choregraphe).

## Installation

Install the `qiloop` command line interface with

```
$ go get -u github.com/lugu/qiloop/cmd/qiloop
$ qiloop --version
Version: 0.8
```

## Prerequisite

On NAO, the service directory URL is `tcp://<ROBOT_IP>:9559`
(replace `<ROBOT_IP>` with the robot IP address).

On Pepper, the connection to the service directory is encrypted and
authenticated. The service directory URL is `tcps://<ROBOT_IP>:9503`.
The user login is `nao` and the password is the password of the UNIX
account `nao` on the robot. Create a file called ~/.qiloop-auth.conf with
the two lines (replace `<YOUR_PASSWORD>` with your password):

```
nao
<YOUR_PASSWORD>
```

Test the connection to the service ALVideoDevice on the robot:

On NAO:

```
$ qiloop --qi-url tcp://<ROBOT_IP>:9559 info --service ALVideoDevice
```

On Pepper:

```
$ qiloop --qi-url tcps://<ROBOT_IP>:9503 info --service ALVideoDevice
```

## Generate the proxy

Let's write this application!

First we will **scan** the service ALVideoDevice to extract its
interface in an IDL file. This is done using the `qiloop scan`
command:

```
qiloop scan --qi-url tcp://127.0.0.1:9559 --idl video_device.qi.idl --service ALVideoDevice
```

This produced a file called `video_device.qi.idl` with the list of
methods, signals and properties of ALVideoDevice.

Add one line at the top of the IDL file video_device.idl to specify a
package name: 'package main`

Using this file, we will generate the Go code to access to the
service. The command `qiloop proxy` will read the IDL file and
generate a specialized proxy object for ALVideoDevice.

```
qiloop proxy --idl video_device.qi.idl --output video_proxy.go
```

## Proxy usage

Finally, the main program which uses the proxy of ALVideoDevice to
obtain a image from a camera:

```golang
package main

import (
	"flag"
	"log"
	"os"

	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/type/value"
)

const (
	topCam = 0
	vga    = 1
	rgb    = 13
	fps    = 30
)

func main() {
	flag.Parse()
	// A Session object is used to connect the service directory.
	sess, err := app.SessionFromFlag()
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	// Using this session, let's instanciate our service
	// constructor.
	services := Services(sess)

	// Using the constructor, we request a proxy to ALVideoDevice
	videoDevice, err := services.ALVideoDevice()
	if err != nil {
		log.Fatalf("failed to create video device: %s", err)
	}

	// Configure the camera
	id, err := videoDevice.SubscribeCamera("me", topCam, vga, rgb, fps)
	if err != nil {
		log.Fatalf("failed to initialize camera: %s", err)
	}

	// Request an image
	img, err := videoDevice.GetImageRemote(id)
	if err != nil {
		log.Fatalf("failed to retrieve image: %s", err)
	}

	// GetImageRemote returns an value, let's cast it into a list
	// of values:
	values, ok := img.(value.ListValue)
	if !ok {
		log.Fatalf("invalid return type")
	}
	// Let's extract the image data.
	width := values[0].(value.IntValue).Value()
	heigh := values[1].(value.IntValue).Value()
	pixels := values[6].(value.RawValue).Value()

	log.Printf("camera resolution: %dx%d\n", width, heigh)
	file, err := os.Create("image.rgb")
	if err != nil {
		log.Fatalf("cannot create image: %s", err)
	}
	defer file.Close()
	file.Write(pixels)
}
```

You can convert the generate image into PNG format with:

```
ffmpeg -s 320x240 -pix_fmt rgb24 -i image.rgb image.png
```

Since a lot of code is generated, use `go doc` to browse the
documentation of ALVideoDevice:
```
go doc .

```
