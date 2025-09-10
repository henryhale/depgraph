package util

import (
	"regexp"
	"strings"
)

// split a string into segments - comma as a delimiter
func Explode(value string) *[]string {
	value = strings.TrimSpace(value)
	badChars := []string{" ", "\t", "\n", "\r", "\""}
	for _, char := range badChars {
		value = strings.ReplaceAll(value, char, "")
	}
	tokens := strings.Split(value, ",")
	return &tokens
}

// find all occurences/usage of an imported module/package
func LocateImports(prefix *string, code *string) []string {
	re := regexp.MustCompile(*prefix + `.(\w*)`)
	matches := re.FindAllStringSubmatch(*code, -1)
	result := []string{}
	if matches == nil {
		return result
	}
	for _, match := range matches {
		if len(match) == 2 && len(match[1]) > 0 {
			result = append(result, match[1])
		}
	}
	return result
}
