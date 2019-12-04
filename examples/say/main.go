// Package main illustrates how to make a method call to a remote
// object. It uses the specialized proxy of the text to speech
// service.
package main

import (
	"flag"

	"github.com/lugu/qiloop/app"
	"github.com/lugu/qiloop/bus/services"
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
	textToSpeech, err := services.ALTextToSpeech(session)
	if err != nil {
		panic(err)
	}

	// Remote procedure call: call the method "say" of the service.
	err = textToSpeech.Say("Hi there !")
	if err != nil {
		panic(err)
	}
}
