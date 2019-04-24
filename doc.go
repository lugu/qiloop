// Package qiloop is an implementation of the protocol QiMessaging
// used to interract with the NAO and Pepper robots.
//
// Introduction
//
// QiMessaging is a software bus which exposes services. Services have
// methods (to be called), signals (to be watched). A naming service (the
// service directory) is used to discover and register services. For
// a detailed description of the protocol, please visit
// https://github.com/lugu/qiloop/blob/master/doc/NOTES.md
//
// To connect to a service, a Session object is required: it represents
// the connection to the service directory. Several transport
// protocols are supported (currently TCP, TLS and UNIX socket).
//
// With a session, on can request a proxy object representing a remote
// service. The proxy object contains the helper methods needed to make
// the remote calls and to handle the incomming signal notifications.
//
// Services have methods, signals and properties which are described
// in an IDL (Interface Description Language) format. This IDL file is
// process by the `qiloop` command to generate the go code which allow
// remote access to the service (i.e. the proxy object).
//
// For example, here is the IDL file which describes a service which
// have two methods, one signal and one property:
//
// 	package demo
//
// 	interface RoboticService
// 		fn move(x: int, y: int)
// 		fn say(sentence: str)
// 		sig obstacle(x: int, y: int)
// 		prop battery(level: int)
// 	end
//
// Use 'qiloop proxy' commmand to generate the go code which gives
// access to the service:
//
// 	$ go get github.com/lugu/qiloop/cmd/qiloop
// 	$ qiloop proxy --idl some_service.idl.qi --output proxy_gen.go
//
// The file proxy_gen.go contains a method called Services which gives
// access to the MyRobot service. The example bellow illustrate this.
//
// In order to communicate with an existing service for which the IDL
// file is unknown, the command scan can be use to introspect a
// running instance of the service and generate its IDL description.
//
// 	$ qiloop scan --qi-url "tcp://localhost:9559" --service LogManager --idl log_manager.idl.qi
//
// In order to implement a new service, create an IDL file and run
// the stub command to generate the helper method to register the
// service:
//
// 	$ qiloop stub --idl my_service.idl.qi --output stub_gen.go
//
// The file stub_gen.go the interface to implement as well as the
// helper methods to register the services.
//
// When offering a service, a Server is be used to handle incomming
// connection and to dispatch the requests.
//
// The actual implementation of a service is provided by a
// ServerObject which responds to the call requests and emits the
// signals. ServerObject are attached to a Server via a Service
// interface.
//
package qiloop

import (
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/directory"
	"github.com/lugu/qiloop/bus/session"
)

// NewSession creates a new connection to the service directory
// located at address addr. Use non empty string if credentials are
// required, else provide empty strings. Example of address:
// "tcp://localhost:9559", "tcps://localhost:9443".
func NewSession(addr, user, token string) (bus.Session, error) {
	// TODO: get ride of bus/session/token
	return session.NewSession(addr)
}

// Authenticator decides if a user/token tuple is valid. It is used to
// decide if an incomming connections is authorized to join the
// services.
type Authenticator interface {
	Authenticate(user, token string) bool
}

// Server listens to an interface handles incomming connection. It
// dispatches the message to the services and objects.
type Server interface {
	// NewService register a new service to the service directory.
	NewService(name string, object bus.ServerObject) (bus.Service, error)
	// Session returns a local session object which can be used to
	// access the server without authentication.
	Session() bus.Session
	// Terminate stops the server.
	Terminate() error
	// Returns a channel to wait for the server terminaison.
	WaitTerminate() chan error
}

// NewServiceDirectory starts a service directory listening at address
// addr. Possible addresses are "tcp://localhost:9559",
// "unix:///tmp/sock", "tcps://localhost:9443". Refer to qiloop/net
// for more details.
func NewServiceDirectory(addr string, auth Authenticator) (Server, error) {
	return directory.NewServer(addr, auth)
}