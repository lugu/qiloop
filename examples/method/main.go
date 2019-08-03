// Package main illustrates how to make a method call to a remote
// object. It uses the specialized proxy of the service directory.
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

	// proxies is an helper to access the specialized proxy.
	proxies := services.Services(session)

	// obtain a representation of the service directory
	directory, err := proxies.ServiceDirectory()
	if err != nil {
		panic(err)
	}

	// call the method "services" of the service directory.
	serviceList, err := directory.Services()
	if err != nil {
		panic(err)
	}

	// print the list of services.
	for _, info := range serviceList {
		println("service " + info.Name)
	}
}
