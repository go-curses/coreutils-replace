package main

import (
	"github.com/urfave/cli/v2"

	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/coreutils/notify"
	"github.com/go-curses/ctk"
)

var (
	gApp     ctk.Application
	gCtx     *cli.Context
	gErr     error
	gArgv    []string
	gSearch  string
	gReplace string
	gTargets []string
	gOptions = &options{}
)

func prepareStartup(data []interface{}, argv ...interface{}) cenums.EventFlag {
	var ok bool

	argc := len(argv)
	if argc > 0 {
		if gApp, ok = argv[0].(ctk.Application); !ok {
			log.ErrorF("internal error - ctk.Application not found (%T)", argv[0])
			return cenums.EVENT_STOP
		}
	} else {
		log.ErrorF("internal error - missing arguments")
		return cenums.EVENT_STOP
	}

	if argc > 1 {
		if gArgv, ok = argv[1].([]string); !ok {
			log.ErrorF("internal error - argv not []string (%T)", argv[1])
			return cenums.EVENT_STOP
		}
	} else {
		log.ErrorF("internal error - missing 2nd argument")
		return cenums.EVENT_STOP
	}

	return cenums.EVENT_PASS
}

func prepare(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if len(argv) < 1 {
		log.ErrorF("internal error - prepare arguments not found")
		return cenums.EVENT_STOP
	}

	var ok bool

	if gCtx, ok = argv[1].(*cli.Context); !ok {
		log.ErrorF("internal error - cli.Context not found (%T)", argv[0])
		return cenums.EVENT_STOP
	}

	gOptions.Lock()
	gOptions.regex = gCtx.Bool("regex")
	gOptions.multiLine = gCtx.Bool("multi-line")
	gOptions.dotMatchNl = gCtx.Bool("dot-match-nl")
	gOptions.multiLineDotMatchNl = gCtx.Bool("multi-line-dot-match-nl")
	gOptions.multiLineDotMatchNlInsensitive = gCtx.Bool("multi-line-dot-match-nl-insensitive")
	gOptions.recurse = gCtx.Bool("recurse")
	gOptions.dryRun = gCtx.Bool("dry-run")
	gOptions.all = gCtx.Bool("all")
	gOptions.ignoreCase = gCtx.Bool("ignore-case")
	gOptions.backup = gCtx.Bool("backup")
	gOptions.backupExtension = gCtx.String("backup-extension")
	gOptions.showDiff = gCtx.Bool("show-diff")
	gOptions.interactive = gCtx.Bool("interactive")
	gOptions.quiet = gCtx.Bool("quiet")
	gOptions.Unlock()
	log.DebugF("prepare options=%v", gOptions)

	if gOptions.quiet {
		notifier.Set(notify.Quiet)
	}

	cliArgv := gCtx.Args().Slice()
	cliArgc := len(cliArgv)
	if cliArgc == 0 {
		cli.ShowAppHelpAndExit(gCtx, 1)
		return cenums.EVENT_STOP
	} else if cliArgc < 3 {
		gErr = newUsageError()
		return cenums.EVENT_STOP
	}
	gSearch, gReplace = cliArgv[0], cliArgv[1]
	gTargets = cliArgv[2:]

	if gOptions.interactive {
		log.DebugF("starting interactive rpl")
		return cenums.EVENT_PASS
	}

	log.DebugF("processing non-interactive rpl")
	gErr = processCliWork()
	return cenums.EVENT_STOP
}