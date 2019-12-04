// Package main demonstrates how to move the robot through ALMotion
package main

import (
	"context"
	"flag"
	"time"

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

	// Obtain a proxy to the service
	motion, err := proxy.ALMotion(session)
	if err != nil {
		panic(err)
	}

	// Remote procedure call: call the method "walk init" of the service.
	err = motion.MoveInit()
	if err != nil {
		panic(err)
	}

	motion.WaitUntilMoveIsFinished()
	println("init move is finished")

	// Remote procedure call: call the method "move to" of the service.
	err = motion.MoveTo(0.2, 0, 0)
	if err != nil {
		panic(err)
	}

	// Demonstrate cancellation: create a context which expires in
	// two seconds.
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	// Move will be cancelled in two seconds.
	err = motion.WithContext(ctx).MoveTo(-0.4, 0, 0)
	if err != nil {
		panic(err)
	}
}
