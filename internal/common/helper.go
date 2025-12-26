package common

import (
	"regexp"
	"strings"
)

func CleanString(input string) string {
	result := regexp.MustCompile(`(\\r\\n|\\n|\\t|\s)+`).ReplaceAllString(input, " ")

	// Trim leading and trailing spaces
	return strings.TrimSpace(result)
}
