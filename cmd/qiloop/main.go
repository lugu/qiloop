package main

import "github.com/integrii/flaggy"

// Set version with:
// % go build -ldflags='-X main.version=1.0'
var version = "unknown"

var (
	infoCommand  *flaggy.Subcommand
	logCommand   *flaggy.Subcommand
	scanCommand  *flaggy.Subcommand
	proxyCommand *flaggy.Subcommand
	stubCommand  *flaggy.Subcommand

	serverURL   string = "tcp://localhost:9559"
	serviceName string = ""
	inputFile   string = ""
	outputFile  string = "-"
	packageName string = ""
)

func init() {
	flaggy.SetName("qiloop")
	flaggy.SetDescription("Utility to process QiMessaing IDL files")

	infoCommand = flaggy.NewSubcommand("info")
	infoCommand.Description = "Connect a server and display services info"
	infoCommand.String(&serverURL, "r", "qi-url",
		"server URL (default: tcp://localhost:9559)")
	infoCommand.String(&serviceName, "s", "service", "optional service name")

	logCommand = flaggy.NewSubcommand("log")
	logCommand.Description = "Connect a server and prints logs"
	logCommand.String(&serverURL, "r", "qi-url",
		"server URL (default: tcp://localhost:9559)")

	scanCommand = flaggy.NewSubcommand("scan")
	scanCommand.Description =
		"Connect a server and introspect a service to generate an IDL file"
	scanCommand.String(&serverURL, "r", "qi-url",
		"server URL (default: tcp://localhost:9559)")
	scanCommand.String(&serviceName, "s", "service", "optional service name")
	scanCommand.String(&outputFile, "i", "idl", "IDL file (output)")

	proxyCommand = flaggy.NewSubcommand("proxy")
	proxyCommand.Description =
		"Parse an IDL file and generate the specialized proxy code"
	proxyCommand.String(&inputFile, "i", "idl", "IDL file (input)")
	proxyCommand.String(&outputFile, "o", "output", "proxy file (output)")

	stubCommand = flaggy.NewSubcommand("stub")
	stubCommand.Description =
		"Parse an IDL file and generate the specialized server code"
	stubCommand.String(&inputFile, "i", "idl", "IDL file (input)")
	stubCommand.String(&outputFile, "o", "output", "server stub file (output)")
	stubCommand.String(&packageName, "p", "path", "optional package name")

	flaggy.AttachSubcommand(infoCommand, 1)
	flaggy.AttachSubcommand(logCommand, 1)
	flaggy.AttachSubcommand(scanCommand, 1)
	flaggy.AttachSubcommand(proxyCommand, 1)
	flaggy.AttachSubcommand(stubCommand, 1)

	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.SetVersion(version)
	flaggy.Parse()
}

func main() {
	if infoCommand.Used {
		info(serverURL, serviceName)
	} else if scanCommand.Used {
		scan(serverURL, serviceName, outputFile)
	} else if proxyCommand.Used {
		proxy(inputFile, outputFile)
	} else if stubCommand.Used {
		stub(inputFile, outputFile, packageName)
	} else if logCommand.Used {
		logger(serverURL)
	} else {
		flaggy.DefaultParser.ShowHelpAndExit("missing command")
	}
}
