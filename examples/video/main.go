// Package main demonstrates show to connect ALVideoDevice to retreive
// an image.
//
// You can use ffmpeg to convert the image in PNG with:
//
// 	ffmpeg -s 320x240 -pix_fmt rgb24 -i image.rgb image.png
//
package main

import (
	"flag"
	"log"
	"os"

	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/type/value"
)

const (
	// cameras
	topCam    = 0
	bottomCam = 1
	depthCam  = 2
	stereoCam = 3

	// cameras
	qvga = 1
	vga  = 2
	vga4 = 3

	// colorspace
	yuv  = 10
	rgb  = 11
	hsv  = 12
	dist = 21

	fps = 10
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
	defer videoDevice.Unsubscribe(id)

	// Request an image
	img, err := videoDevice.GetImageRemote(id)
	if err != nil {
		log.Fatalf("failed to retrieve image: %s", err)
	}

	// GetImageRemote returns an value, let's cast it into a list
	// of values:
	values, ok := img.(value.ListValue)
	if !ok {
		log.Fatalf("invalid return type: %#v", img)
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
