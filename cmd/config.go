package cmd

import (
	"flag"
	"strings"
)

type Config struct {
	Dir          *string
	Lang         *string
	OutputFile   *string
	OutputFormat *string
	IgnoredPaths []string
	ReplacePaths map[string]string
	ShowHelp     *bool
	ShowVersion  *bool
}

func ParseConfig() Config {
	config := Config{}

	config.Dir = flag.String("d", ".", "Specifies the target `directory` to analyze (default: current working directory).\n")
	config.Lang = flag.String("l", "", "The programming `language` of the files: js, ts, c, cpp\n")
	config.OutputFormat = flag.String("f", "mermaid", "The output `format` of the analysis: mermaid (default), jsoncanvas, json\n")
	config.OutputFile = flag.String("o", "stdout", "Write output to the selected `file` (default: stdout)\n")
	config.ShowHelp = flag.Bool("h", false, "Show information about the command-line options and exit.\n")
	config.ShowVersion = flag.Bool("v", false, "Show the current version information and exit.\n")
	ignoredPaths := flag.String("i", "", "A comma-separated list of `directories` to ignore for example: 'tests,dist,build'\n")
	replacePaths := flag.String("r", "", "A key:value  comma-separated list of `paths` to replace for example: '@:src,demo:src/demo'\n")

	flag.Parse()

	config.IgnoredPaths = strings.Split(*ignoredPaths, ",")

	config.ReplacePaths = make(map[string]string)
	segments := strings.Split(*replacePaths, ",")
	for _, segment := range segments {
		v := strings.Split(segment, ":")
		if len(v) == 2 {
			config.ReplacePaths[v[0]] = v[1]
		}
	}

	return config
}
