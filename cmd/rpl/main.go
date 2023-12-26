package main

import (
	"os"

	"github.com/go-curses/cdk"
	cstrings "github.com/go-curses/cdk/lib/strings"

	"github.com/go-curses/corelibs/notify"

	"github.com/go-curses/coreutils-replace/ui"
)

var (
	APP_NAME  = "rpl"
	APP_USAGE = "text search and replace utility"
	APP_DESC  = `rpl is a command line utility for searching and replacing text in one or more
text files.

When invoked with the --show-diff flag, a unified diff is generated and
printed to STDOUT with any other info or error notices printed to STDERR.

rpl has an interactive mode with a curses user interface for walking through
all files matching the search parameter. When in interactive mode, and the
unified diff output is requested, the diff is printed to STDERR upon rpl
exiting the user interface.

Examples:

 # recursively change "search" to "replaced" in all matching files:
 rpl -R -D "search" "replaced" > /tmp/search-replaced.patch

 # same as above but interactively:
 rpl -I -R -D "search" "replaced" 2> /tmp/search-replaced.patch
`
	APP_VERSION = "0.5.0"
	APP_RELEASE = "trunk"
	APP_TAG     = "rpl"
	APP_TITLE   = "rpl"
)

// Build Configuration Flags
// setting these will enable command line flags and their corresponding features
// use `go build -v -ldflags="-X 'main.IncludeLogFullPaths=false'"`
var (
	IncludeProfiling          = "false"
	IncludeLogFile            = "false"
	IncludeLogFormat          = "false"
	IncludeLogFullPaths       = "false"
	IncludeLogLevel           = "false"
	IncludeLogLevels          = "false"
	IncludeLogTimestamps      = "false"
	IncludeLogTimestampFormat = "false"
	IncludeLogOutput          = "false"
	notifier                  *notify.Notifier
)

func init() {
	cdk.Build.Profiling = cstrings.IsTrue(IncludeProfiling)
	cdk.Build.LogFile = cstrings.IsTrue(IncludeLogFile)
	cdk.Build.LogFormat = cstrings.IsTrue(IncludeLogFormat)
	cdk.Build.LogFullPaths = cstrings.IsTrue(IncludeLogFullPaths)
	cdk.Build.LogLevel = cstrings.IsTrue(IncludeLogLevel)
	cdk.Build.LogLevels = cstrings.IsTrue(IncludeLogLevels)
	cdk.Build.LogTimestamps = cstrings.IsTrue(IncludeLogTimestamps)
	cdk.Build.LogTimestampFormat = cstrings.IsTrue(IncludeLogTimestampFormat)
	cdk.Build.LogOutput = cstrings.IsTrue(IncludeLogOutput)
	notifier = notify.New(notify.Info)
}

func main() {
	cdk.Init()
	if err := ui.NewUI(
		APP_NAME,
		APP_USAGE,
		APP_DESC,
		APP_VERSION,
		APP_RELEASE,
		APP_TAG,
		APP_TITLE,
		"/dev/tty",
		notifier,
	).Run(os.Args); err != nil {
		notifier.Error("%v\n", err)
		os.Exit(1)
	}
}