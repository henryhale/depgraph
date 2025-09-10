package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/henryhale/depgraph/cmd"
	"github.com/henryhale/depgraph/internal/graph"
	"github.com/henryhale/depgraph/internal/lang"
	"github.com/henryhale/depgraph/internal/output"
	"github.com/henryhale/depgraph/internal/util"
)

// command name
const Name string = "depgraph"

// version number
var version = "(untracked)"

func main() {
	// setup a logger
	log.SetPrefix(Name + ": ")
	log.SetFlags(0)

	// parse command line flags
	config := cmd.ParseConfig()

	// help
	if *config.ShowHelp {
		fmt.Print("Usage: ", Name, " [options]\n\n")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// version
	if *config.ShowVersion {
		fmt.Println(Name, "version", version)
		os.Exit(0)
	}

	// select language
	if len(*config.Lang) == 0 {
		log.Fatal("programming language not specified")
	}
	pl, found := lang.Get(*config.Lang)
	if !found {
		log.Fatal("'" + *config.Lang + "' language is not yet supported")
	}

	// check output format
	if !output.FormatSupported(config.OutputFormat) {
		log.Fatal("'" + *config.OutputFormat + "' output format is not supported")
	}

	// read target directory
	// - filter out ignored file paths
	// - verify file extensions
	files, err := util.TraverseDirectory(config.Dir, &pl.Extensions, &config.IgnoredPaths)
	if err != nil {
		log.Fatal(err)
	}

	// build deps map - analyze each file
	deps := make(graph.DependencyGraph)
	// keep track of external dependencies
	external := make(map[string][]string)

	extractorOptions := new(lang.ExtractorOptions)
	extractorOptions.Replacers = &config.ReplacePaths

	for _, filePath := range *files {
		result := lang.SourceFile{
			Imports: make(map[string][]string),
			Exports: []string{},
			Local:   true,
		}

		extractorOptions.Result = &result
		extractorOptions.File = &filePath

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		sourceCode := util.Preprocess(string(fileContent), pl.Comments)

		for _, rule := range pl.Rules {
			re := regexp.MustCompile(rule.RegExp)
			matches := re.FindAllStringSubmatch(sourceCode, -1)
			if matches == nil {
				continue
			}

			extractorOptions.Rule = &rule

			for _, match := range matches {
				extractorOptions.Match = &match

				pl.Extract(extractorOptions)
			}
		}

		deps[filePath] = result

		// ensure all imports exist or atleast external
		_, isExternal := external[filePath]
		if isExternal {
			delete(external, filePath)
		}
		for importpath, items := range result.Imports {
			_, exists := deps[importpath]
			if !exists {
				_, alreadyExternal := external[importpath]
				if alreadyExternal {
					external[importpath] = append(external[importpath], items...)
				} else {
					external[importpath] = items
				}
			}
		}
	}

	// add externals
	for path, exports := range external {
		deps[path] = lang.SourceFile{
			Imports: make(map[string][]string),
			Exports: exports,
			Local:   false,
		}
	}

	// produce formatted output
	output := output.Format(config.OutputFormat, &deps)

	// done!
	if *config.OutputFile == "stdout" {
		fmt.Println(output)
	} else {
		err := os.WriteFile(*config.OutputFile, []byte(output), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

}
