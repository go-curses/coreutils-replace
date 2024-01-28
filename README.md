[![Made with Go](https://img.shields.io/badge/go-v1.21+-blue.svg)](https://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-curses/coreutils-replace)](https://goreportcard.com/report/github.com/go-curses/coreutils-replace)

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
   rpl [options] <search> <replace> [path...]

VERSION:
   v0.9.4 (trunk)

DESCRIPTION:

   rpl is a command line utility for searching and replacing content within plain
   text files.

   * rpl supports long-form and short-form command-line flags
   * provides diff and notices on different outputs (os.Stdout and os.Stderr)
   * has a Go-Curses user interface for interactively selecting changes to apply
     (in the spirit of "git add --patch")
   * can preserve the case of the per-instance strings being replaced
   * supports regular expressions applied on a per-line or multi-line basis
   * target files can be provided via os.Stdin, files or command-line arguments


   Case operations:

    # change all instances of "search" exactly with "replace"
    #
    # flags: none (case-sensitive)

    rpl "search" "replace" *

    # change all instances of "search" with "replace"
    #
    # flags: --ignore-case (-i)

    rpl -i "search" "replace" *

    # change all instances of "SearchQuery" with "ReplaceValue", which will
    # ignores case while finding matches, detects the individual replacement
    # cases and maintains that case with the replacement value
    #
    # flags: --preserve-case (-P)

    rpl -P "SearchQuery" "ReplaceValue" *
    #
    # changes "SEARCHQUERY" to "REPLACEVALUE"
    # changes "searchquery" to "replacevalue"
    # changes "SearchQuery" to "ReplaceValue"
    # changes "searchQuery" to "replaceValue"

    # don't actually change all instances of "search" with "replace"
    #
    # flags: --ignore-case (-i), --nop (-n)

    rpl -in "search" "replace" *


   Interactive operations:

    # rpl has a curses based user-interface for interactively applying changes
    # to the matching files and works with all other command-line flags
    #
    # flags: --ignore-case (-i), --interactive (-e)

    rpl -ie "search" "replace" *
    #
    # once all the files have been filtered through any --include or --exclude
    # options, the user-interface walks through each file, prompting the user
    # with a unified diff of the changes to either Skip or Apply to the given
    # file. If there are more than one group of edits within the unified diff,
    # an additional option, Select, is added which allows the user to walk
    # through all the edit groups to pick and choose, similarly to how git
    # works with the "git add --patch ..." operation


   Backup operations:

    # backup and change all instances of "search" with "replace"; backup
    # files are named a trailing "~"; if the filename already exists, backup
    # file names end with a "~", an incremented number, and another "~"
    #
    # flags: --ignore-case (-i), --backup (-b)

    rpl -ib "search" "replace" *
    #
    # first backup filename: example.txt~
    # second backup filename: example.txt~1~

    # backup and change all instances of "search" with "replace"; backup
    # files are named a ".bak" extension; if the filename already exists,
    # backup file names end with a ".", an incremented number, and have a
    # ".bak" extension
    #
    # flags: --ignore-case (-i), --backup-extension (-B)

    rpl -i -B .bak "search" "replace" *
    #
    # first backup filename: example.txt.bak
    # second backup filename: example.txt.1.bak


   Unified diff output:

    # don't actually recursively change "search" to "replaced" and print a
    # universal diff of the changes that would be made to STDOUT, any errors
    # or other notices are printed to STDERR
    #
    # flags: --nop (-n), --recurse (-r), --show-diff (-d)

    rpl -nrd "search" "replaced" .

    # same as above but save the diff output to a file
    #
    # flags: --nop (-n), --recurse (-r), --show-diff (-d)

    rpl -nrd "search" "replaced" . > /tmp/search-replaced.patch

    # same as above but interactively, which outputs the user-interface to
    # STDOUT and so the diff is output to STDERR
    #
    # flags: --nop (-n), --interactive (-e), --recurse (-r), --show-diff (-d)

    rpl -nerd "search" "replaced" . 2> /tmp/search-replaced.patch


   Regular Expression operations:

    # rpl supports search and replace operations using the Go language version
    # of regular expressions, see: https://pkg.go.dev/regexp/syntax for the
    # supported syntax and limitations; all search patterns are prefixed with
    # a global (?m) flag because the default mode would be to only search and
    # replace within the first line of a file's content
    #
    # flags: --regex (-p), --ignore-case (-i)

    rpl -pi '([a-z])([-a-z0-9]+?)' '${1}_${2}' *
    #
    # this pattern captures two groups of characters, the first is a single
    # lower-case letter and the second is one or more dashes, lower-case
    # letters or numbers; because the --ignore-case flag is also present,
    # the search pattern is prefixed with a global (?i) flag to match text
    # case- insensitively
    #
    # the replacement pattern simply separates the two groups with an
    # underscore, note that the Perl \1 syntax is not supported and that
    # single-quotes are used to ensure the shell does not interpret the
    # regex variables as shell variables

    # one of the great regex flags is the (?s) option which changes the
    # interpretation of the (any character) ".", caret "^" and other syntax
    # to include newlines. The user can just add a leading (?s) to the search
    # pattern but rpl includes a --dot-match-nl flag which does this and when
    # used, the normal --regex flag is not required
    #
    # flags: --multi-line (m) --dot-match-nl (-s), --ignore-case (-i)

    rpl -msi '^func thing\(\) {\s^(.+?)^}$' 'func renamed() {\n${1}\n}' *.go
    #
    # this pattern captures the multi-line contents of static Go functions
    # named "thing" and simply renames them to "renamed"


   File selection operations:

    # rpl accepts a variable list of path arguments which are individually
    # converted to their absolute path equivalents and tested to see if they
    # even exist at all. Sometimes this method of supplying file names is not
    # good enough, such as when there are more files in the list than the OS
    # allows in a single command. For these sorts of cases, there are a few
    # more ways to supply paths, with all of them working together to build up
    # a list of all the files to work with

    # command line arguments method:
    #
    # flags: (none)

    rpl "search" "replace" ...

    # using newline-separated paths listed within one or more text files:
    #
    # flags: --file (-f)

    rpl -f filenames.txt "search" "replace"

    # using newline-separated paths derived from standard input:
    #
    # flags: (none), with single dash path

    find * -type f | rpl "search" "replace" -

    # using null-separated paths derived from standard input:
    #
    # flags: --null (-0), with single dash path

    find * -type f -print0 | rpl -0 "search" "replace" -

    # excluding globs of files:
    #
    # flags: --exclude (-X)

    rpl -X "*.bak" "search" "replace" *
    #
    # any file with the extension .bak is completely ignored

    # including globs of files:
    #
    # flags: --include (-I)

    rpl -I "*.txt" "search" "replace" *
    #
    # any file without the extension .txt is completely ignored

    # combining --exclude with --include requires the files to satisfy both
    # conditions, they must not be excluded and also must be included too
    #
    # flags: --exclude (-X), --include (-I)

    rpl -X "example.*" -I "*.txt" -I "*.md" "search" "replace" *
    #
    # replaces "search" with "replace" in all files that do not start with the
    # word "example" and also end with .txt or .md extensions

   Limitations:

   * maximum file size: 5.2 MB
   * maximum number of files: 100,000
   * more than 10k changes per file can consume gigabytes of memory


GLOBAL OPTIONS:

   1. Case Sensitivity

   --ignore-case, -i    perform a case-insensitive search (plain or regex)
   --preserve-case, -P  try to preserve replacement string cases

   2. Regular Expressions

   --dot-match-nl, -s  set the dot-match-nl (?s) global flag (implies -r)
   --multi-line, -m    set the multiline (?m) global flag (implies -r)
   --regex, -r         search and replace arguments are regular expressions

   3. User Interface

   --interactive, -e  selectively apply changes per-file
   --pause, -E        pause on file search results screen (implies -e)
   --show-diff, -d    output unified diffs for all changes

   4. Backups

   --backup, -b                        make backups before replacing content
   --backup-extension value, -B value  specify the backup file suffix to use (implies -b)

   5. Target Selection

   --all, -a                  include backups and files that start with a dot
   --exclude value, -X value  exclude files matching glob pattern
   --file value, -f value     read paths listed in files
   --include value, -I value  include on files matching glob pattern
   --null, -0                 read null-terminated paths from os.Stdin
   --recurse, -R              travel directory paths

   6. General

   --help             display complete command-line help text
   --no-limits, -U    ignore max file count and size limits
   --nope, --nop, -n  report what would otherwise have been done
   --quiet, -q        silence notices
   --usage, -h        display command-line usage information
   --verbose, -v      verbose notices
   --version, -V      display the version

```


## LICENSE

```
Copyright 2024  The Go-Curses Authors

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
