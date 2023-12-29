# replace

Command-line utility for searching and replacing text.

## INSTALLATION

### Debian APT installation method

```shell
> wget -c https://apt.go-curses.org/apt-go-curses-org_latest.deb
> sudo dpkg -i ./apt-go-curses-org_latest.deb
> sudo apt update
> sudo apt install replace
```

### Brew installation method

```shell
# this is a "head-only" formula, to install:
> brew install --HEAD go-curses/tap/replace
# and to upgrade:
> brew upgrade --fetch-HEAD replace
```

### Golang installation method

```shell
> go install github.com/go-curses/coreutils-replace/cmd/rpl@latest
```

## DOCUMENTATION

``` shell
> rpl --help
NAME:
   rpl - text search and replace utility

USAGE:
   rpl [options] <search> <replace> <path> [path...]

VERSION:
   v0.5.2 (trunk)

DESCRIPTION:

   rpl is a command line utility for searching and replacing text in one or more
   text files.

   When invoked with the --show-diff flag, a unified diff is generated and
   printed to STDOUT with any other info or error notices printed to STDERR.

   rpl has an interactive mode with a curses user interface for walking through
   all files matching the search parameter. When in interactive mode, and the
   unified diff output is requested, the diff is printed to STDERR upon rpl
   exiting the user interface.

   Examples:

    # recursively change "search" to "replaced" in all matching files:
    rpl -rd "search" "replaced" . > /tmp/search-replaced.patch

    # same as above but interactively:
    rpl -rdI "search" "replaced" . 2> /tmp/search-replaced.patch


OPTIONS:
   Configuration

   --backup, -b           make backups before replacing content
   --bak value, -B value  specify the backup file suffix to use (implies -b)
   --ignore-case, -i      perform a case-insensitive search (literal or regex)
   --interactive, -I      selectively apply changes per file
   --nop, -n              report what would have otherwise been done
   --show-diff, -d        include unified diffs of all changes in the output

   Expressions

   --dot-match-nl, -s     set the dot-match-nl (?s) regex flag (implies -p)
   --multi-line, -m       set the multi-line (?m) regex flag (implies -p)
   --regex, -p            search and replace arguments are regular expressions

   File Selection

   --all, -a              include files that start with a "."
   --recurse, -r          recurse into subdirectories
   -x value [ -x value ]  exclude files matching glob patterns

   General

   --help, -h             display command-line usage information
   --quiet, -q            run silently, ignored if --nop is also used
   --verbose, -v          run loudly, ignored if --quiet is also used
   --version, -V          display the version
```


## LICENSE

```
Copyright 2023  The Go-Curses Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use file except in compliance with the License.
You may obtain a copy of the license at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```