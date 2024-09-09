# depgraph

A developer's tool to understanding new codebases

## Overview

This is a cli tool that tends to visualize a codebase right from inter-file
~~to in-file~~ dependency so you can understand how it is developed. It is basically
based on the fact that every programming language has known standard approaches
to handling modules/packages; imports and exports. Unlike existing dependency
resolution/visualization tools that are language-specific or parser-centric, 
depgraph is language agnostic since it uses regular expressions to match import 
and export statements from which a dependency graph is constructed. While this is
experimental, it produces amazing results for small vanilla projects. The graph 
is formatted and output as json, jsoncanvas, or mermaid. 

## Why?

I landed onto a codebase and couldn't figure out what was going on fast;
have you ever experienced that before? It is overwhelmingly a rewarding process that 
takes me many hours trying to figure out lots of stuff in a new codebasee.
So, I wanted to speed up the process a little bit such that I can easily fix a bug/issue 
or contribute to any project seamlessly. Additionally, I needed to learn ideas and 
patterns from several amazing opensource projects visually. Big codebases can scare,
but not anymore, depgraph is here!

## Features

- [x] cli
- [ ] multi-language support (planning on c, cpp, js, ts, go, php and others)
- [ ] interactive web interface (maybe d3.js or cyptoscape.js)
- [ ] generating images from the graph (like png, svg)
- [ ] export dependency graphs (probably mermaid)

... and more to come

## Installation

>:warning: work in progress

## Usage

Once installed, use
- `depgraph -v` to show version information
- `depgraph -h` to display help message

**Required arguments**
- `-d <path>` specifies the path to the directory containing source files
- `-l <language>` sets the programming language: `ts`, `js`, `go`, `c`, `cpp`, `php` 
- `-f <format>` specifies the output format of the result: `json`, `jsoncanvas`, `mermaid`

**Optional arguments**
- `-i [path1,path2, ...]` defines a list of comma-separated paths to ignore; for example `-i "tests,dist,build,node_modules"`
- `-r [old:new, ...]` defines a list of comma-separated key:value paths to replace; for example `-r "@:src"`

## Examples

- vanilla js project with tests, need mermaid visual
	```sh
	depgraph -d /path/to/project -f mermaid -l js -i tests
	```
- ts/js project with npm packages, tests, root directory alias (src -> @) and json output
	```sh
	depgraph -d /path/to/project -f json -l ts -i "tests,node_modules,dist" -r "@:src"
	```
 
To save output to a file, say `stats.json`, you can use i/o redirection, that is;
```sh
depgraph -d /path/to/project -f json -l go > stats.json
```

>Note: In case the output format is `mermaid`, you can use the mermaid vscode 
>extension or [mermaid.live](https://mermaid.live) to view the output.
>In case of `jsoncanvas`, you can use Obsidian.

## How it works

This section describes how depgraph works when you run the command;

1. Parsing CLI arguments: determine the target directory, output format, language and more
2. Directory traversal: build a list of files, ignoring filtered paths
3. Resolving dependencies: use regular expression to match imports and exports
4. Building a dependency graph: transform files into nodes and edges from imports
5. Formating the output: produce a stringified representation of the graph - json, jsoncanvas, mermaid
6. Printing the results: output the result to standard output

## Contributing

Thank you for looking into this amazing project. Incase of any issues, bugs, or proposing a new feature: [open a new issue](https://github.com/henryhale/depgraph/issues/new).

## Building

To build this project locally, ensure that you have [Go](https://go.dev/doc/install) installed.

Clone this repository using: `git clone https://github.com/henryhale/depgraph.git`

At the root of the repository, there exists a shell script which when executed yields a binary executable; 
using your shell, run

```sh
bash build.sh
```

or

```sh
chmod +x build.sh
./build.sh
```

# License

&copy; 2024 [Henry Hale](https://github.com/henryhale).

Release under [MIT License](https://github.com/henryhale/depgraph/blob/master/LICENSE.txt)

