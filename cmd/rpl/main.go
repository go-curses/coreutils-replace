package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/pkg/profile"

	"github.com/go-corelibs/notify"
	rpl "github.com/go-corelibs/replace"
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/env"
	replace "github.com/go-curses/coreutils-replace"
	"github.com/go-curses/coreutils-replace/ui"
)

var (
	AppName    = "rpl"
	AppVersion = "0.9.7"
	AppRelease = "trunk"
	AppTag     = "rpl"
	AppTitle   = "rpl"
	AppUsage   = "text search and replace utility"
	AppDesc    = `
rpl is a command line utility for searching and replacing content within plain
text files.

* rpl supports long-form and short-form command-line flags
* provides diff and notices on different outputs (os.Stdout and os.Stderr)
* has a Go-Curses user interface for interactively selecting changes to apply
  (in the spirit of "git add --patch")
* can preserve the case of the per-instance strings being replaced
* supports regular expressions applied on a per-line or multi-line basis
* target files can be provided via os.Stdin, files or command-line arguments
`
	AppManual = `

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
 # flags: --nop (-n), --recurse (-R), --show-diff (-d)

 rpl -nRd "search" "replaced" .

 # same as above but save the diff output to a file
 #
 # flags: --nop (-n), --recurse (-R), --show-diff (-d)

 rpl -nRd "search" "replaced" . > /tmp/search-replaced.patch

 # same as above but interactively, which outputs the user-interface to
 # STDOUT and so the diff is output to STDERR
 #
 # flags: --nop (-n), --interactive (-e), --recurse (-R), --show-diff (-d)

 rpl -neRd "search" "replaced" . 2> /tmp/search-replaced.patch


Regular Expression operations:

 # rpl supports search and replace operations using the Go language version
 # of regular expressions, see: https://pkg.go.dev/regexp/syntax for the
 # supported syntax and limitations; all search patterns are prefixed with
 # a global (?m) flag because the default mode would be to only search and
 # replace within the first line of a file's content
 #
 # flags: --regex (-r), --ignore-case (-i)

 rpl -ri '([a-z])([-a-z0-9]+?)' '${1}_${2}' *
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

* maximum file size: ` + replace.MaxFileSizeLabel + `
* maximum number of files: ` + humanize.Comma(int64(rpl.MaxFileCount)) + `
* more than 10k changes per file can consume gigabytes of memory
`
)

var (
	notifier = notify.New(notify.Info).Make()
)

func init() {
	// set a very aggressive gc threshold to curb large-files with many
	// matches from consuming all available memory
	debug.SetGCPercent(10)
}

func main() {
	if profileType := env.Get("GO_CDK_PROFILE", ""); profileType != "" {
		var p func(p *profile.Profile)
		switch strings.ToLower(profileType) {
		case "cpu":
			p = profile.CPUProfile
		case "mem":
			p = profile.MemProfile
		default:
			_, _ = fmt.Fprintf(os.Stderr, "GO_CDK_PROFILE must be one of: cpu or mem\n")
			os.Exit(1)
		}
		if profilePath := env.Get("GO_CDK_PROFILE_PATH", ""); profilePath != "" {
			defer profile.Start(p, profile.ProfilePath(profilePath)).Stop()
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "GO_CDK_PROFILE requires GO_CDK_PROFILE_PATH to be set\n")
			os.Exit(1)
		}
	}
	defer func() {
		if v := recover(); v != nil {
			_, _ = fmt.Fprintf(os.Stderr, "# rpl panic: %q\n", v)
			if slices.Contains(os.Args, "-v") || slices.Contains(os.Args, "--verbose") {
				_, _ = fmt.Fprintln(os.Stderr, "#")
				for _, line := range strings.Split(string(debug.Stack()), "\n") {
					_, _ = fmt.Fprintln(os.Stderr, "#\t"+line)
				}
			}
			os.Exit(1)
		}
	}()
	cdk.Init()
	if err := ui.NewUI(
		AppName,
		AppUsage,
		AppDesc,
		AppManual,
		AppVersion,
		AppRelease,
		AppTag,
		AppTitle,
		"/dev/tty",
		notifier,
	).Run(os.Args); err != nil {
		notifier.Error("# %v\n", err)
		os.Exit(1)
	}
}
