// Package main demonstrates how to move the robot through ALMotion
package main

import (
	"flag"
	"fmt"

	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/examples/motion/proxy"
)

func main() {
	flag.Parse()
	// session represents a connection to the service directory.
	session, err := app.SessionFromFlag()
	if err != nil {
		panic(err)
	}
	defer session.Terminate()

	// Access the specialized proxy constructor.
	services := proxy.Services(session)

	// Obtain a proxy to the service
	motion, err := services.ALMotion()
	if err != nil {
		panic(err)
	}

	// Remote procedure call: call the method "walk init" of the service.
	err = motion.MoveInit()
	if err != nil {
		panic(err)
	}

	motion.WaitUntilMoveIsFinished()
	fmt.Println("init move is finished")

	// Remote procedure call: call the method "move to" of the service.
	err = motion.MoveTo(0.2, 0, 0)
	if err != nil {
		panic(err)
	}
}
