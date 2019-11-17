package main

import (
	"fmt"
	"log"

	"github.com/integrii/flaggy"
	"github.com/lugu/qiloop/bus/session/token"
	asciibot "github.com/mattes/go-asciibot"
)

// Set version with:
// % go build -ldflags='-X main.version=1.0'
var version = "0.7"

var (
	infoCommand   *flaggy.Subcommand
	logCommand    *flaggy.Subcommand
	scanCommand   *flaggy.Subcommand
	proxyCommand  *flaggy.Subcommand
	stubCommand   *flaggy.Subcommand
	serverCommand *flaggy.Subcommand
	traceCommand  *flaggy.Subcommand

	serverURL   = "tcp://localhost:9559"
	serviceName = ""
	objectID    = uint32(1)
	logLevel    = uint32(5) // LogLevelVerbose
	inputFile   = ""
	outputFile  = "-"
	packageName = ""
)

func init() {
	flaggy.SetName("qiloop")
	description := fmt.Sprintf("%s\n\n%s",
		"an utility to explore QiMessaging",
		asciibot.Random())
	flaggy.SetDescription(description)

	authDescription := fmt.Sprintf("credentials (default: %s)", token.AuthFile)

	infoCommand = flaggy.NewSubcommand("info")
	infoCommand.Description = "Connect a server and display services info"
	infoCommand.String(&serverURL, "r", "qi-url",
		"server URL (default: tcp://localhost:9559)")
	infoCommand.String(&serviceName, "s", "service", "optional service name")
	infoCommand.String(&token.AuthFile, "a", "auth-file", authDescription)

	logCommand = flaggy.NewSubcommand("log")
	logCommand.Description = "Connect a server and prints logs"
	logCommand.String(&serverURL, "r", "qi-url",
		"server URL (default: tcp://localhost:9559)")
	logCommand.String(&token.AuthFile, "a", "auth-file", authDescription)
	levelInfo := "log level, 1:fatal, 2:error, 3:warning, 4:info, 5:verbose, 6:debug"
	logCommand.UInt32(&logLevel, "l", "level", levelInfo)

	scanCommand = flaggy.NewSubcommand("scan")
	scanCommand.Description =
		"Connect a server and introspect a service to generate an IDL file"
	scanCommand.String(&serverURL, "r", "qi-url",
		"server URL (default: tcp://localhost:9559)")
	scanCommand.String(&serviceName, "s", "service", "optional service name")
	scanCommand.String(&outputFile, "i", "idl", "IDL file (output)")
	scanCommand.String(&token.AuthFile, "a", "auth-file", authDescription)

	proxyCommand = flaggy.NewSubcommand("proxy")
	proxyCommand.Description =
		"Parse an IDL file and generate the specialized proxy code"
	proxyCommand.String(&inputFile, "i", "idl", "IDL file (input)")
	proxyCommand.String(&outputFile, "o", "output", "proxy file (output)")
	proxyCommand.String(&packageName, "p", "path", "optional package name")

	stubCommand = flaggy.NewSubcommand("stub")
	stubCommand.Description =
		"Parse an IDL file and generate the specialized server code"
	stubCommand.String(&inputFile, "i", "idl", "IDL file (input)")
	stubCommand.String(&outputFile, "o", "output", "server stub file (output)")
	stubCommand.String(&packageName, "p", "path", "optional package name")

	serverCommand = flaggy.NewSubcommand("server")
	serverCommand.Description =
		"Start a service directory and a log manager"
	serverCommand.String(&serverURL, "l", "qi-listen-url",
		"Listening URL (default: tcp://localhost:9559)")
	serverCommand.String(&token.AuthFile, "a", "auth-file", authDescription)

	traceCommand = flaggy.NewSubcommand("trace")
	traceCommand.Description = "Connect a server and traces services"
	traceCommand.String(&serverURL, "r", "qi-url",
		"server URL (default: tcp://localhost:9559)")
	traceCommand.String(&serviceName, "s", "service", "optional service name")
	traceCommand.UInt32(&objectID, "o", "object", "optional object id")
	traceCommand.String(&token.AuthFile, "a", "auth-file", authDescription)

	flaggy.AttachSubcommand(infoCommand, 1)
	flaggy.AttachSubcommand(logCommand, 1)
	flaggy.AttachSubcommand(scanCommand, 1)
	flaggy.AttachSubcommand(proxyCommand, 1)
	flaggy.AttachSubcommand(stubCommand, 1)
	flaggy.AttachSubcommand(serverCommand, 1)
	flaggy.AttachSubcommand(traceCommand, 1)

	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.SetVersion(version)
	flaggy.Parse()
}

func main() {
	log.SetFlags(0)

	if infoCommand.Used {
		info(serverURL, serviceName)
	} else if scanCommand.Used {
		scan(serverURL, serviceName, outputFile)
	} else if proxyCommand.Used {
		proxy(inputFile, outputFile, packageName)
	} else if stubCommand.Used {
		stub(inputFile, outputFile, packageName)
	} else if logCommand.Used {
		logger(serverURL, logLevel)
	} else if serverCommand.Used {
		server(serverURL)
	} else if traceCommand.Used {
		trace(serverURL, serviceName, objectID)
	} else {
		flaggy.DefaultParser.ShowHelpAndExit("missing command")
	}
}
