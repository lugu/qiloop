package qiloop_test

import (
	"github.com/lugu/qiloop"
	"github.com/lugu/qiloop/bus/services"
)

// This example shows how to create a session which connect a server
// and instanciate the proxy object of the service directory to list
// the services.
func Example_basic() {
	// Create a new session.
	session, err := qiloop.NewSession(
		"tcps://localhost:9443", // service directory URL
		"nao",                   // user
		"nao",                   // token
	)
	if err != nil {
		panic(err)
	}

	// Access the specialized proxy generated.
	proxies := services.Services(session)

	// Access a proxy object of the service directory.
	directory, err := proxies.ServiceDirectory()
	if err != nil {
		panic(err)
	}

	// Remote procedure call: call the method "services" of the
	// service directory.
	serviceList, err := directory.Services()
	if err != nil {
		panic(err)
	}

	// Iterate over the list of services.
	for _, info := range serviceList {
		println("service " + info.Name)
	}
}
