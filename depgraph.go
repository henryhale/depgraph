package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/henryhale/depgraph/cmd"
	"github.com/henryhale/depgraph/export"
	"github.com/henryhale/depgraph/lang"
	"github.com/henryhale/depgraph/util"
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
		fmt.Print("Usage:", Name,"[options]\n\n")
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
		log.Fatal("'" + *config.Lang + "' is not yet supported")
	}

	// read target directory
	// - filter out ignored file paths
	// - verify file extensions
	files, err := util.TraverseDirectory(config.Dir, &pl.Extensions, &config.IgnoredPaths)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*files)

	// build deps map - analyze each file
	deps := make(lang.DependencyGraph)
	// keep track of external dependencies
	external := make(map[string][]string)

	extractorOptions := new(lang.ExtractorOptions)
	extractorOptions.Replacers = &config.ReplacePaths

	for _, file := range *files {
		result := lang.SourceFile{
			Imports: make(map[string][]string),
			Exports: []string{},
			Local: true,
		}

		extractorOptions.Result = &result
		extractorOptions.File = &file.Path
		
		file.Code = util.Preprocess(file.Code, pl.Comments)
		
		for _, rule := range pl.Rules {
			re := regexp.MustCompile(rule.RegExp)
			matches := re.FindAllStringSubmatch(file.Code, -1)
			if matches == nil {
				continue
			}

			extractorOptions.Rule = &rule

			for _, match := range matches {
				extractorOptions.Match = &match

				pl.Extract(extractorOptions)
			}
		}

		deps[file.Path] = result

		// ensure all imports exist or atleast external
		_, isExternal := external[file.Path]
		if isExternal {
			delete(external, file.Path)
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
			Local: false,
		}
	}

	// produce formatted output
	output := export.Format(config.OutputFormat, &deps)

	// done!
	if *config.OutputFile == "stdout" {
		fmt.Println(output)
	} else {
		err := os.WriteFile(*config.OutputFile, []byte(output), 644)
		if err != nil {
			log.Fatal(err)
		}
	}

}
