package cmd

import (
	"flag"
	"strings"
)

type Options struct {
	Dir *string
	Lang *string
	OutputFormat *string
	IgnoredPaths []string
	ReplacePaths map[string]string
	ShowHelp *bool
	ShowVersion *bool
}

func InitArgs() Options {
	options := Options{}

	options.Dir = flag.String("d", "", "Specifies the target `directory` to analyze.\n")
	options.Lang = flag.String("l", "", "The programming `language` of the files: js, c, cpp, go, php\n")
	options.OutputFormat = flag.String("f", "", "The output `format` of the analysis: json, jsoncanvas, mermaid\n")
	options.ShowHelp = flag.Bool("h", false, "Show information about the command-line options and exit.\n")
	options.ShowVersion = flag.Bool("v", false, "Show the current version information and exit.\n")
	ignoredPaths := flag.String("i", "", "A comma-separated list of `directories` to ignore.\n")
	replacePaths := flag.String("r", "", "A key:value  comma-separated list of `paths` to replace.\n")

	flag.Parse()

	options.IgnoredPaths = strings.Split(*ignoredPaths, ",")

	options.ReplacePaths = make(map[string]string)
	segments := strings.Split(*replacePaths, ",")
	for _, segment := range segments {
		v := strings.Split(segment, ":")
		if len(v) == 2 {
			options.ReplacePaths[v[0]] = v[1]
		}
	}

	return options
}
