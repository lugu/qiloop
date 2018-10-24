package util

import (
	"regexp"
)

// CleanName remove funky character from a name
func CleanName(name string) string {
	exp := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	return exp.ReplaceAllString(name, "_")
}
