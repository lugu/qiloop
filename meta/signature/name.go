package signature

import (
	"fmt"
	"regexp"
	"strings"
)

var reservedMethods = []string{
	"Subscribe", "MetaObject", "Properties", "Property",
	"RegisterEvent", "RegisterEventWithSignature", "SetProperty",
	"Terminate", "UnregisterEvent", "Call", "CallID", "MethodID",
	"ObjectID", "OnDisconnect", "PropertyID", "ProxyService",
	"ServiceID", "SignalID", "Subscribe", "SubscribeID",
}
var keywords = []string{
	"break", "default", "func", "interface", "select",
	"case", "defer", "go", "map", "struct", "chan",
	"else", "goto", "package", "switch", "const",
	"fallthrough", "if", "range", "type", "continue",
	"for", "import", "return", "var", "error",
}

func ValidName(name string) string {
	exp := regexp.MustCompile(`[^_a-zA-Z0-9]+`)
	return exp.ReplaceAllString(name, "")
}

// CleanName remove funky character from a name
func CleanName(name string) string {
	return strings.Title(ValidName(name))
}

func CleanMethodName(name string) string {
	name = CleanName(name)
	for _, keyword := range reservedMethods {
		if name == keyword {
			return fmt.Sprintf("Do%s", name)
		}
	}
	return name
}

func CleanVarName(i int, name string) string {
	if name == "" {
		return fmt.Sprintf("P%d", i)
	}
	name = ValidName(name)
	for _, keyword := range keywords {
		if name == keyword {
			return fmt.Sprintf("%s_%d", name, i)
		}
	}
	return name
}
