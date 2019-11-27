// Package main demonstrates how to move the robot through ALMotion
package main

import (
	"flag"
	"fmt"

	"bitbucket.org/swoldt/pkg/xerrors/iferr"
	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/examples/posture/proxy"
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
	motion, err := services.ALRobotPosture(nil)
	iferr.Exit(err)

	ok, err := motion.GoToPosture("StandZero", 2.0)
	iferr.Exit(err)

	fmt.Println(ok)

	ok, err = motion.GoToPosture("Stand", 2.0)
	iferr.Exit(err)

}
