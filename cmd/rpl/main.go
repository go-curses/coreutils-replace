package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/go-curses/coreutils/notify"
	"github.com/go-curses/coreutils/path"
)

var appName string = "rpl"
var log *notify.Notifier

func init() {
	appName = filepath.Base(os.Args[0])
	log = notify.New(notify.Quiet)
}

func processPath(c *cli.Context, s, r, p string) (errs []error) {
	p = strings.TrimRight(p, "/")

	if !path.Exists(p) {
		errs = append(errs, newError(fmt.Sprintf(`"%v" not found`, p)))
		return
	}

	if path.IsDir(p) {
		if c.Bool("regex") ||
			c.Bool("multi-line") ||
			c.Bool("dot-match-nl") ||
			c.Bool("multi-line-dot-match-nl") ||
			c.Bool("multi-line-dot-match-nl-insensitive") {
			return rplPathRegexp(c, s, r, p)
		}
		return rplPathString(c, s, r, p)
	}

	if c.Bool("regex") ||
		c.Bool("multi-line") ||
		c.Bool("dot-match-nl") ||
		c.Bool("multi-line-dot-match-nl") ||
		c.Bool("multi-line-dot-match-nl-insensitive") {
		if err := rplFileRegexp(c, s, r, p); err != nil {
			errs = append(errs, err)
		}
		return
	}
	if err := rplFileString(c, s, r, p); err != nil {
		errs = append(errs, err)
	}
	return
}

func action(c *cli.Context) error {
	if !c.Bool("quiet") {
		log = notify.New(notify.Info)
	}

	tmpdir := c.String("tmpdir")
	if !path.IsDir(tmpdir) {
		return newError(fmt.Sprintf(
			`"%v" not found or not a directory`,
			tmpdir,
		))
	}

	argv := c.Args().Slice()
	argc := len(argv)
	if argc == 0 {
		cli.ShowAppHelpAndExit(c, 1)
		return nil
	} else if argc < 3 {
		return newUsageError()
	}
	s, r := argv[0], argv[1]
	paths := argv[2:]

	errs := []error{}
	for _, p := range paths {
		if es := processPath(c, s, r, p); len(es) > 0 {
			errs = append(errs, es...)
			for _, e := range es {
				log.Error("%s\n", e.Error())
			}
		}
	}

	if len(errs) > 0 {
		return newError(fmt.Sprintf(
			"%d errors occurred",
			len(errs),
		))
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:            appName,
		Usage:           "command line search and replace",
		Version:         "0.1.0",
		Action:          action,
		HideHelpCommand: true,
		ArgsUsage:       "<search> <replace> <path> [...paths]",
		Description:     "Search and replace content within files using strings or regular expressions.",
		Flags: []cli.Flag{
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
				Name:    "quiet",
				Usage:   "run silently, ignored if --dry-run is also used",
				Aliases: []string{"q"},
			},
			&cli.BoolFlag{
				Name:    "backup",
				Usage:   "make backups before replacing content",
				Aliases: []string{"b"},
			},
			&cli.StringFlag{
				Name:    "backup-extension",
				Usage:   "specify the backup file extension to use",
				Aliases: []string{"B"},
				Value:   "bak",
			},
			&cli.BoolFlag{
				Name:    "show-diff",
				Usage:   "include unified diffs of all changes in the output",
				Aliases: []string{"D"},
			},
			&cli.StringFlag{
				Name:    "tmpdir",
				Usage:   "specify the tmpdir to use",
				Aliases: []string{"T"},
				EnvVars: []string{"TMPDIR"},
				Value:   os.TempDir(),
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Error("%v\n", err)
		os.Exit(1)
	}
}
