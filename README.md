<div align=center>

# depgraph

A developer's tool to understanding new codebases

[![](https://mermaid.ink/img/pako:eNptkE0OgjAQha9CZi0egIUrl650SQ2ZthNBAiX9STDA3W0pEhRXff3m9c10BhBKEmTw0NiVyeXK2iQxjscrHp8mgGRWBfY59vcAqJVfRr4agyp4n_P_RrEYf3mpnKGYMMuCK5-hYsiHCaXqivJ4bOPDdGl6Gr3JjOv7UOHbyvKJwMWWLzPv-Lbrrhh3AwdoSDdYSb_DIZgY2JIaYpB5KVHXDFg7eR86q26vVkBmtaMDuE6ipXOFfgNNhNMbLWiKbg?type=png)](https://mermaid.live/edit#pako:eNptkE0OgjAQha9CZi0egIUrl650SQ2ZthNBAiX9STDA3W0pEhRXff3m9c10BhBKEmTw0NiVyeXK2iQxjscrHp8mgGRWBfY59vcAqJVfRr4agyp4n_P_RrEYf3mpnKGYMMuCK5-hYsiHCaXqivJ4bOPDdGl6Gr3JjOv7UOHbyvKJwMWWLzPv-Lbrrhh3AwdoSDdYSb_DIZgY2JIaYpB5KVHXDFg7eR86q26vVkBmtaMDuE6ipXOFfgNNhNMbLWiKbg)

</div>

## Overview

This is a cli tool that tends to visualize a codebase right from inter-file
~~to in-file~~ dependency so you can understand how it is developed. It is basically
based on the fact that every programming language has known standard approaches
to handling modules/packages; imports and exports. Unlike existing dependency
resolution/visualization tools that are language-specific or parser-centric,
depgraph is language agnostic since it uses regular expressions to match import
and export statements from which a dependency graph is constructed. While this is
experimental, it produces amazing results for small vanilla projects. The graph
is formatted and output as json, [jsoncanvas](jsoncanvas.org), or [mermaid](mermaid.js.org).

> [!WARNING]
> depgraph is experimental and incomplete. It is under active
> development.

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
- [x] export dependency graphs as formatted text
  - [x] json
  - [x] [jsoncanvas](https://jsoncanvas.org)
  - [x] [mermaid](https://mermaid.js.org)
  - [x] [dot](https://graphviz.org/doc/info/lang.html)
- multi-language support - _**work in progress**_
  - [x] js, ts - _almost done_
  - [x] c, cpp - _in progress_
  - [ ] go
  - [ ] python
  - [ ] php
  - [ ] ....
- [ ] interactive web interface (maybe d3.js or cyptoscape.js)
- [ ] generating images from the graph (like png, svg)

... and more to come

## Installation

- Using a shell script: Linux/Mac/Termux/WSL
	```sh
	curl -fsSL https://raw.githubusercontent.com/henryhale/depgraph/master/scripts/install.sh | bash
	```
- Prebuilt Executables:

	Go to the [Github releases page](https://github.com/henryhale/depgraph/releases/latest) and download a prebuilt executable for your platform/machine.

## Usage

Once installed, use
- `depgraph -h` to display help message
- `depgraph -v` to show version information

**Required arguments**
- `-l <language>` sets the programming language: `ts`, `js`, `c`, `cpp` <!-- `go`, `php` -->

**Optional arguments**
- `-d <path>` specifies the path to the directory containing source files (_default: current working directory_)
- `-f <format>` specifies the output format of the result: `mermaid` - _default_, `dot`, `jsoncanvas`, `json`
- `-o <path>` write output to the selected `file` (_default: `stdout`_)
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

## Output

By default, the output is written to `stdout`.
To save output to a file, you can use;

- `-o <path>` option
  ```sh
  depgraph -d /path/to/project -l js -f json -o graph.json  
  ```
- I/O redirection
  ```sh
  depgraph -d /path/to/project -l js -f json > graph.json  
  ```

### Output Formats

- **mermaid** - _default_: In case the output format is `mermaid`, you can instantly use [mermaid.live](https://mermaid.live) to view the output otherwise checkout this [complete list](https://mermaid.js.org/ecosystem/integrations-community.html).
- **dot**: For a quick in-browser visualization & image export, check out: [Edotor](https://edotor.net/), [GraphvizOnline](https://dreampuf.github.io/GraphvizOnline/), [Graphviz Visual Editor](https://magjac.com/graphviz-visual-editor/). Otherwise [Graphviz](https://graphviz.org/download/) has a [command line utility](https://graphviz.org/doc/info/command.html) to generate images from the `dot` output.
- **jsoncanvas**: To view the visual output from `jsoncanvas` output, use one of the apps on this [list](https://jsoncanvas.org/docs/apps/).
- **json**: This option is meant from storage and usage with json viewer tools. While the search for a compatible and appropriate visualization tool for json output continues, this option may be _removed_ in future releases.

## How it works

This section describes how depgraph works when you run the command;

1. Parsing CLI arguments: determine the target directory, output format, language and more
2. Directory traversal: build a list of files, ignoring filtered paths
3. Resolving dependencies: use regular expression to match imports and exports
4. Building a dependency graph: transform files into nodes and edges from imports
5. Formating the output: produce a stringified representation of the graph - json, jsoncanvas, mermaid
6. Printing the results: output the result to standard output or file

## Contributing

Thank you for looking into this amazing project. Incase of any issues, bugs, or proposing a new feature: [open a new issue](https://github.com/henryhale/depgraph/issues/new).

## Building

To build this project locally, ensure that you have [Go v1.23.2](https://go.dev/doc/install) installed.

Clone this repository using: `git clone https://github.com/henryhale/depgraph.git`

At the root of the repository, there exists a shell script which when executed yields a binary executable;
using your shell, run

```sh
bash scripts/build.sh
```

or

```sh
chmod +x scripts/build.sh
./scripts/build.sh
```

# License

&copy; 2024-present [Henry Hale](https://github.com/henryhale).

Release under [MIT License](https://github.com/henryhale/depgraph/blob/master/LICENSE.txt)
