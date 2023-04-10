package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/go-curses/cdk"
	cstrings "github.com/go-curses/cdk/lib/strings"

	"github.com/go-curses/corelibs/notify"

	"github.com/go-curses/ctk"
)

var (
	APP_NAME    = "rpl"
	APP_USAGE   = "search and replace"
	APP_DESC    = "command line search and replace"
	APP_VERSION = "0.2.3"
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
	APP_NAME = filepath.Base(os.Args[0])
	APP_TITLE = APP_NAME
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
	app := ctk.NewApplication(
		APP_NAME,
		APP_USAGE,
		APP_DESC,
		APP_VERSION+" ("+APP_RELEASE+")",
		APP_TAG,
		APP_TITLE,
		"/dev/tty",
	)

	app.CLI().ArgsUsage = ""
	app.CLI().HideHelpCommand = true
	app.CLI().UsageText = APP_NAME + " [options] <search> <replace> <path> [path...]"

	app.CLI().Flags = append(
		app.CLI().Flags,
		&cli.BoolFlag{
			Name:    "regex",
			Usage:   "search and replace arguments are regular expressions",
			Aliases: []string{"P"},
		},
		&cli.BoolFlag{
			Name:    "multi-line",
			Usage:   "set the multi-line (?m) regexp flag (implies -P)",
			Aliases: []string{"m"},
		},
		&cli.BoolFlag{
			Name:    "dot-match-nl",
			Usage:   "set the dot-match-nl (?s) regexp flag (implies -P)",
			Aliases: []string{"s"},
		},
		&cli.BoolFlag{
			Name:    "multi-line-dot-match-nl",
			Usage:   "convenience flag to set -m and -s (implies -P)",
			Aliases: []string{"ms", "p"},
		},
		&cli.BoolFlag{
			Name:    "multi-line-dot-match-nl-insensitive",
			Usage:   "convenience flag to set -m, -s and -i (implies -P)",
			Aliases: []string{"msi"},
		},
		&cli.BoolFlag{
			Name:    "recurse",
			Usage:   "recurse into sub-directories",
			Aliases: []string{"R"},
		},
		&cli.BoolFlag{
			Name:    "dry-run",
			Usage:   "report what would have otherwise been done",
			Aliases: []string{"n"},
		},
		&cli.BoolFlag{
			Name:    "all",
			Usage:   "include files and directories that start with a \".\"",
			Aliases: []string{"a"},
		},
		&cli.BoolFlag{
			Name:    "ignore-case",
			Usage:   "perform a case-insensitive search",
			Aliases: []string{"i"},
		},
		&cli.BoolFlag{
			Name:    "interactive",
			Usage:   "selectively apply replacement edits",
			Aliases: []string{"I"},
		},
		&cli.BoolFlag{
			Name:    "backup",
			Usage:   "make backups before replacing content",
			Aliases: []string{"b"},
		},
		&cli.StringFlag{
			Name:    "backup-extension",
			Usage:   "specify the backup file suffix to use",
			Aliases: []string{"B"},
			Value:   ".bak",
		},
		&cli.BoolFlag{
			Name:    "show-diff",
			Usage:   "include unified diffs of all changes in the output",
			Aliases: []string{"D"},
		},
		&cli.BoolFlag{
			Name:    "quiet",
			Usage:   "run silently, ignored if --dry-run is also used",
			Aliases: []string{"q"},
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Usage:   "run loudly, ignored if --quiet is also used",
			Aliases: []string{"v"},
		},
	)

	app.Connect(cdk.SignalPrepareStartup, "rpl-prepare-startup-handler", prepareStartup)
	app.Connect(cdk.SignalPrepare, "rpl-prepare-handler", prepare)
	app.Connect(cdk.SignalStartup, "rpl-startup-handler", startup)
	app.Connect(cdk.SignalShutdown, "rpl-shutdown-handler", shutdown)

	if err := app.Run(os.Args); err != nil {
		notifier.Error("%v\n", err)
		os.Exit(1)
	}
}
