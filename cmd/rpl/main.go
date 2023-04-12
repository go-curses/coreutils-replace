package main

import (
	"os"

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
	app.CLI().UsageText = APP_NAME + " [options] <search> <replace> <path> [path...]"
	app.CLI().HideHelpCommand = true
	app.CLI().EnableBashCompletion = true
	app.CLI().UseShortOptionHandling = true
	cli.HelpFlag = &cli.BoolFlag{
		Category: "General",
		Name:     "help",
		Usage:    "display command-line usage information",
		Aliases:  []string{"h", "usage"},
	}
	cli.VersionFlag = &cli.BoolFlag{
		Category: "General",
		Name:     "version",
		Usage:    "display the version",
		Aliases:  []string{"V"},
	}

	app.CLI().Flags = append(
		app.CLI().Flags,

		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "recurse",
			Usage:    "recurse into sub-directories",
			Aliases:  []string{"R"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "dry-run",
			Usage:    "report what would have otherwise been done",
			Aliases:  []string{"n"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "all",
			Usage:    "include things that start with a \".\"",
			Aliases:  []string{"a"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "interactive",
			Usage:    "selectively apply replacement edits",
			Aliases:  []string{"I"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "backup",
			Usage:    "make backups before replacing content",
			Aliases:  []string{"b"},
		},
		&cli.StringFlag{
			Category: "Configuration",
			Name:     "backup-extension",
			Usage:    "specify the backup file suffix to use",
			Aliases:  []string{"B"},
			Value:    ".bak",
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "show-diff",
			Usage:    "include unified diffs of all changes in the output",
			Aliases:  []string{"D"},
		},

		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "regex",
			Usage:    "search and replace arguments are regular expressions",
			Aliases:  []string{"P"},
		},
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "multi-line",
			Usage:    "set the multi-line (?m) regexp flag (implies -P)",
			Aliases:  []string{"m"},
		},
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "dot-match-nl",
			Usage:    "set the dot-match-nl (?s) regexp flag (implies -P)",
			Aliases:  []string{"s"},
		},
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "multi-line-dot-match-nl",
			Usage:    "convenience flag to set -m and -s (implies -P)",
			Aliases:  []string{"ms", "p"},
		},
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "multi-line-dot-match-nl-insensitive",
			Usage:    "convenience flag to set -m, -s and -i (implies -P)",
			Aliases:  []string{"msi"},
		},
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "ignore-case",
			Usage:    "perform a case-insensitive search",
			Aliases:  []string{"i"},
		},

		&cli.BoolFlag{
			Category: "General",
			Name:     "quiet",
			Usage:    "run silently, ignored if --dry-run is also used",
			Aliases:  []string{"q"},
		},
		&cli.BoolFlag{
			Category: "General",
			Name:     "verbose",
			Usage:    "run loudly, ignored if --quiet is also used",
			Aliases:  []string{"v"},
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
