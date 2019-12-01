// Package main illustrates how to make a method call to a remote
// object. It uses the specialized proxy of the animated speech service.
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

	// Access the specialized proxy constructor.
	proxies := services.Services(session)

	// Obtain a proxy to the service
	ans, err := proxies.ALAnimatedSpeech(nil)
	if err != nil {
		panic(err)
	}

	err = ans.SetBodyLanguageEnabled(true)
	if err != nil {
		panic(err)
	}

	err = ans.SetBodyTalkEnabled(true)
	if err != nil {
		panic(err)
	}

	// Remote procedure call: call the method "say" of the service.
	err = ans.Say("^startTag(explain) I have always been a very good explainer, did you know? ^waitTag(explain)")
	if err != nil {
		panic(err)
	}
}
