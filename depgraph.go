package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/henryhale/depgraph/cmd"
	"github.com/henryhale/depgraph/export"
	"github.com/henryhale/depgraph/lang"
	"github.com/henryhale/depgraph/utils"
)

// version of depgraph
var version = "(untracked)"

func main() {
	// parse cli arguments
	options := cmd.InitArgs()

	// help
	if *options.ShowHelp {
		fmt.Print("Usage: depgraph [options]\n\n")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// version
	if *options.ShowVersion {
		fmt.Println("depgraph version", version)
		os.Exit(0)
	}

	// select parser
	if len(*options.Lang) == 0 {
		fmt.Println("error: programming language not specified")
		os.Exit(1)
	}
	pl, found := lang.Get(*options.Lang)
	if !found {
		fmt.Println("error: '" + *options.Lang + "' is not yet supported")
		os.Exit(1)
	}

	// read target directory
	// - filter out ignored file paths
	files, err := utils.TraverseDirectory(options.Dir, &pl.Extensions, &options.IgnoredPaths)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// build deps map - analyze each file
	deps := make(export.AnalysisResultMap)
	// keep track of external dependencies
	external := make(map[string][]string)

	for _, path := range *files {
		file, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		code := string(file)

		result := lang.AnalysisResult{
			Imports: make(map[string][]string),
			Exports: []string{},
		}

		for _, rule := range pl.Rules {
			re := regexp.MustCompile(rule.RegExp)
			matches := re.FindAllStringSubmatch(code, -1)
			if matches == nil {
				continue
			}

			for _, match := range matches {
				// exports
				if rule.Export && rule.Items > 0 {
					result.AddExport(*utils.Explode(match[rule.Items])...)
					continue
				}
				// imports
				if !rule.Export && rule.File > 0 {
					if rule.Items > 0 {
						importpath := utils.FullPath(match[rule.File], path, &options.ReplacePaths)
						result.AddImport(importpath, *utils.Explode(match[rule.Items]))
					} else {
						// incase of multiple line imports: go
						// incase of non-specific imported items: go, python, dart
						paths := strings.ReplaceAll(match[rule.File], "\n", ",")
						segments := utils.Explode(paths)
						for _, p := range *segments {
							if len(p) == 0 {
								continue
							}
							p = utils.FullPath(p, path, &options.ReplacePaths)

							// use regexp to match "p.xxx" - usage of import
							if pl.LocateImports {
								usedImports := utils.LocateImports(&p, &code)
								result.AddImport(p, usedImports)
							}
						}
					}
				}
			}
		}

		deps[path] = result

		// ensure all imports exist or atleast external
		_, isExternal := external[path]
		if isExternal {
			delete(external, path)
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
		deps[path] = lang.AnalysisResult{
			Imports: make(map[string][]string),
			Exports: exports,
		}
	}

	// produce formatted output
	output := export.Format(options.OutputFormat, &deps)

	// done!
	fmt.Println(output)

}
