package app

import (
	"flag"
	"fmt"
	"log"
	"syscall"

	"github.com/lugu/qiloop"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/util"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	sessionURL string
	listenURL  string
	userName   string
	userToken  string
)

func init() {
	flag.StringVar(&sessionURL, "qi-url", "tcp://localhost:9559",
		"Service directory URL")
	flag.StringVar(&listenURL, "listen-url", "", "URL to listen to")
	flag.StringVar(&userName, "user", "", "user name")
	flag.StringVar(&userToken, "token", "", "user token")
}

// SessionFromFlag is an helper method to write client applications.
// It uses the flag package to get the session URL, the user and token
// strings to connect a session.
func SessionFromFlag() (bus.Session, error) {

	if userName != "" && userToken == "" {
		fmt.Println("Your token: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("Token input error: %s", err)
		}
		userToken = string(bytePassword)
	}

	return qiloop.NewSession(sessionURL, userName, userToken)
}

// ServerFromFlag is an helper method to write service applications.
// It uses the flag package to get the session URL, the URL to listen
// to, the user and token strings. It creates a new session, connect
// to the service directory and register the given service and listen
// for incoming connections.
func ServerFromFlag(serviceName string, object bus.Actor) (bus.Server, error) {

	session, err := SessionFromFlag()
	if err != nil {
		log.Fatalf("Failed to connect %s: %s", sessionURL, err)
	}

	if listenURL == "" {
		listenURL = util.NewUnixAddr()
	}

	var auth bus.Authenticator = bus.Yes{}
	if userName != "" {
		passwords := map[string]string{
			userName: userToken,
		}
		auth = bus.Dictionary(passwords)
	}
	server, err := services.NewServer(session, listenURL, auth)
	if err != nil {
		log.Fatalf("Failed to start server at %s: %s", listenURL, err)
	}

	_, err = server.NewService(serviceName, object)
	if err != nil {
		return nil, fmt.Errorf("Failed to register %s: %s", serviceName, err)
	}

	return server, nil
}
