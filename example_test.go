package qiloop_test

import (
	"github.com/lugu/qiloop"
	"github.com/lugu/qiloop/bus/services"
)

// This example llustrates how to make a method call to a remote
// object. It uses the specialized proxy of the text to speech
// service.
func Example_basic() {
	// imports the following packages:
	// 	"github.com/lugu/qiloop"
	// 	"github.com/lugu/qiloop/bus/services"

	// Create a new session.
	session, err := qiloop.NewSession(
		"tcp://localhost:9559", // Service directory URL
		"",                     // user
		"",                     // token
	)
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
