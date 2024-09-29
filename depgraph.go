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
	// parse command line flags
	config := cmd.ParseConfig()

	// help
	if *config.ShowHelp {
		fmt.Print("Usage: " + cmd.Name + " [options]\n\n")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// version
	if *config.ShowVersion {
		fmt.Println(cmd.Name + " version", version)
		os.Exit(0)
	}

	// select language
	if len(*config.Lang) == 0 {
		cmd.Fatal("programming language not specified")
	}
	pl, found := lang.Get(*config.Lang)
	if !found {
		cmd.Fatal("'" + *config.Lang + "' is not yet supported")
	}

	// read target directory
	// - filter out ignored file paths
	files, err := utils.TraverseDirectory(config.Dir, &pl.Extensions, &config.IgnoredPaths)
	if err != nil {
		cmd.Fatal(err)
	}

	// build deps map - analyze each file
	deps := make(lang.DependencyGraph)
	// keep track of external dependencies
	external := make(map[string][]string)

	for _, path := range *files {
		file, err := os.ReadFile(path)
		if err != nil {
			cmd.Fatal(err)
		}

		code := string(file)

		result := lang.SourceFile{
			Imports: make(map[string][]string),
			Exports: []string{},
			Local: true,
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
						importpath := utils.FullPath(match[rule.File], path, &config.ReplacePaths)
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
							p = utils.FullPath(p, path, &config.ReplacePaths)

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
			cmd.Fatal(err)
		}
	}

}
