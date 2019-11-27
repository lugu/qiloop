// Package main demonstrates how to move the robot through ALMotion
package main

import (
	"flag"
	"fmt"

	"bitbucket.org/swoldt/pkg/xerrors/iferr"
	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/examples/motion/proxy"
)

func main() {
	flag.Parse()
	// session represents a connection to the service directory.
	session, err := app.SessionFromFlag()
	iferr.Exit(err)
	defer session.Terminate()

	// Access the specialized proxy constructor.
	services := proxy.Services(session)

	// Obtain a proxy to the service
	motion, err := services.ALMotion(nil)
	iferr.Exit(err)

	// Remote procedure call: call the method "walk init" of the service.
	err = motion.WalkInit()
	iferr.Exit(err)

	motion.WaitUntilMoveIsFinished()
	fmt.Println("init move is finished")

	// Remote procedure call: call the method "walk to" of the service.
	err = motion.WalkTo(0.2, 0, 0)
	iferr.Exit(err)
}
