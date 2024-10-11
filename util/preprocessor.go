package util

import (
	"regexp"
)

// rules for matching comments
var Comments = []string{
	// single line comment
	`//.*[\n\r]`,
	// multiple line commment
	`/\*[\s\S]*\*/`,
}

func Preprocess(code string, rules *[]string) string {
	for _, rule := range *rules {
		re := regexp.MustCompile(rule)
		code = re.ReplaceAllString(code, "")
	}
	return code
}