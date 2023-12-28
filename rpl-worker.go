// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package replace

import (
	"fmt"
	"io"
	"os"
	"regexp"

	glob "github.com/ganbarodigital/go_glob"
	"github.com/urfave/cli/v2"

	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/corelibs/filewriter"
	"github.com/go-curses/corelibs/notify"
)

type Worker struct {
	Regex           bool
	MultiLine       bool
	DotMatchNl      bool
	Recurse         bool
	DryRun          bool
	All             bool
	IgnoreCase      bool
	Backup          bool
	BackupExtension string
	ShowDiff        bool
	Interactive     bool
	Quiet           bool
	Verbose         bool

	Context *cli.Context
	Argv    []string
	Argc    int

	Search  string
	Pattern *regexp.Regexp
	Replace string
	Exclude []*glob.Glob

	Targets []string
	Files   []string
	Matched []string

	Notifier *notify.Notifier

	fwo filewriter.Writer
	fwe filewriter.Writer
}

func MakeWorker(ctx *cli.Context, notifier *notify.Notifier) (w *Worker, eventFlag cenums.EventFlag, err error) {
	w = &Worker{
		Regex:           ctx.Bool("regex"),
		MultiLine:       ctx.Bool("multi-line"),
		DotMatchNl:      ctx.Bool("dot-match-nl"),
		Recurse:         ctx.Bool("recurse"),
		DryRun:          ctx.Bool("nop"),
		All:             ctx.Bool("all"),
		IgnoreCase:      ctx.Bool("ignore-case"),
		Backup:          ctx.Bool("backup"),
		BackupExtension: ctx.String("bak"),
		ShowDiff:        ctx.Bool("show-diff"),
		Interactive:     ctx.Bool("interactive"),
		Quiet:           ctx.Bool("quiet"),
		Verbose:         ctx.Bool("verbose"),
		Context:         ctx,
		Argv:            ctx.Args().Slice(),
		Argc:            ctx.NArg(),
		Notifier:        notifier,
	}

	if w.BackupExtension != "" {
		w.Backup = true
	}

	if !w.Regex {
		w.Regex = w.DotMatchNl || w.MultiLine
	}

	var o, e io.Writer
	var level notify.Level
	if w.Quiet {
		w.Verbose = false
		level = notify.Quiet
	} else if w.Verbose {
		level = notify.Debug
	} else {
		level = notify.Info
	}
	if w.Interactive {
		// stdout is in use by curses, so if os.Stderr is piped, use that, otherwise need to write to a temp file
		var stat os.FileInfo
		if stat, err = os.Stderr.Stat(); err != nil {
			err = fmt.Errorf("error calling os.Stderr.Stat: %w", err)
			return
		} else if stat.Mode()&os.ModeNamedPipe != 0 {
			// user has piped the output of Stderr
			o = os.Stderr
		} else {
			// can't actually write to stderr because it's likely to the same terminal as curses
			if w.fwo, err = filewriter.NewTempFileWriter(fmt.Sprintf("rpl-%d-*.out", os.Getpid())); err != nil {
				err = fmt.Errorf("error making new stdout temp file writer: %w", err)
				return
			}
			o = w.fwo
		}
		// stderr is always to a temp file when interactive
		if w.fwe, err = filewriter.NewTempFileWriter(fmt.Sprintf("rpl-%d-*.err", os.Getpid())); err != nil {
			err = fmt.Errorf("error making new stderr temp file writer: %w", err)
			return
		}
		e = w.fwe
	} else {
		o = os.Stdout
		e = os.Stderr
	}
	w.Notifier.Set(level, o, e)

	if patterns := ctx.StringSlice("x"); len(patterns) > 0 {
		for _, pattern := range patterns {
			g := glob.NewGlob(pattern)
			if _, err = g.Match("./test/file.txt"); err != nil {
				err = fmt.Errorf("-x %q error: %w", pattern, err)
				return
			}
			w.Exclude = append(w.Exclude, g)
		}
	}

	if w.Argc < 3 {
		cli.ShowAppHelpAndExit(ctx, 1)
		eventFlag = cenums.EVENT_STOP
		return
	}
	w.Search, w.Replace = w.Argv[0], w.Argv[1]
	w.Targets = w.Argv[2:]

	if w.Regex {
		if w.Pattern, err = MakeRegexp(w.Search, w); err != nil {
			err = fmt.Errorf("error compiling regular expression: %w", err)
		}
	}

	return
}

func (w *Worker) Init(fn FindAllMatchingFn) {
	if w.Regex {
		w.Files, w.Matched = FindAllMatchingRegexp(w.Pattern, w.Targets, w.Exclude, fn)
	} else if w.IgnoreCase {
		w.Files, w.Matched = FindAllMatchingStringInsensitive(w.Search, w.Targets, w.Exclude, fn)
	} else {
		w.Files, w.Matched = FindAllMatchingString(w.Search, w.Targets, w.Exclude, fn)
	}
	return
}

func (w *Worker) Start() (iter *Iterator) {
	if len(w.Matched) > 0 {
		iter = &Iterator{
			w:   w,
			pos: 0,
		}
	}
	return
}

func (w *Worker) FileWriterOut() (fwo filewriter.Writer) {
	fwo = w.fwo
	return
}

func (w *Worker) FileWriterErr() (fwe filewriter.Writer) {
	fwe = w.fwe
	return
}

func (w *Worker) String() (s string) {
	s += "Worker{"
	s += fmt.Sprintf("regex=%v;", w.Regex)
	s += fmt.Sprintf("multiLine=%v;", w.MultiLine)
	s += fmt.Sprintf("dotMatchNl=%v;", w.DotMatchNl)
	s += fmt.Sprintf("recurse=%v;", w.Recurse)
	s += fmt.Sprintf("dryRun=%v;", w.DryRun)
	s += fmt.Sprintf("all=%v;", w.All)
	s += fmt.Sprintf("ignoreCase=%v;", w.IgnoreCase)
	s += fmt.Sprintf("backup=%v;", w.Backup)
	s += fmt.Sprintf("backupExtension=%v;", w.BackupExtension)
	s += fmt.Sprintf("showDiff=%v;", w.ShowDiff)
	s += fmt.Sprintf("interactive=%v;", w.Interactive)
	s += fmt.Sprintf("quiet=%v;", w.Quiet)
	s += fmt.Sprintf("verbose=%v;", w.Verbose)
	s += "}"
	return
}