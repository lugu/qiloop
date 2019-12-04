// Package main demonstrates how to move the robot through ALMotion
package main

import (
	"flag"
	"math/rand"

	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/examples/posture/proxy"
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
	alPosture, err := proxy.ALRobotPosture(session)
	if err != nil {
		panic(err)
	}

	postures, err := alPosture.GetPostureList()
	if err != nil {
		panic(err)
	}

	// pick a random posture
	if len(postures) == 0 {
		panic("no posture available")
	}

	p := postures[rand.Intn(len(postures))]
	println("Trying posture", p)

	// Go to posture
	ok, err := alPosture.GoToPosture(p, 1.0)
	if err != nil {
		panic(err)
	}

	println("Posture success: ", ok)
}
