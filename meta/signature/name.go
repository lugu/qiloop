package signature

import (
	"regexp"
	"strings"
)

// CleanName remove funky character from a name
func CleanName(name string) string {
	exp := regexp.MustCompile(`[^_a-zA-Z0-9]+`)
	return strings.Title(exp.ReplaceAllString(name, ""))
}
