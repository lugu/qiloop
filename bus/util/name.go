package util

import (
	"regexp"
	"strings"
)

// CleanName remove funky character from a name
func CleanName(name string) string {
	exp := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return strings.Title(exp.ReplaceAllString(name, ""))
}
