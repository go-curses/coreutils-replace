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
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/urfave/cli/v2"

	"github.com/go-corelibs/filewriter"
	"github.com/go-corelibs/notify"
	"github.com/go-corelibs/path"
	"github.com/go-corelibs/scanners"
	cenums "github.com/go-curses/cdk/lib/enums"
)

var (
	MaxFileCount   = 10000
	TempErrPattern = fmt.Sprintf("rpl-%d.*.err", os.Getpid())
	TempOutPattern = fmt.Sprintf("rpl-%d.*.out", os.Getpid())
)

var (
	ErrNotFound = errors.New("not found")
)

type Worker struct {
	Regex           bool
	DotMatchNl      bool
	Recurse         bool
	Nop             bool
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
	Stdin   bool
	Null    bool
	AddFile []string
	Include Globs
	Exclude Globs

	Paths   []string
	Targets []string
	Files   []string
	Matched []string

	Notifier notify.Notifier

	fwo filewriter.FileWriter
	fwe filewriter.FileWriter
}

func MakeWorker(ctx *cli.Context, notifier notify.Notifier) (w *Worker, eventFlag cenums.EventFlag, err error) {
	w = &Worker{
		Regex:           ctx.Bool(RegexFlag) || ctx.Bool(DotMatchNlFlag),
		DotMatchNl:      ctx.Bool(DotMatchNlFlag),
		Recurse:         ctx.Bool(RecurseFlag),
		Nop:             ctx.Bool(NopFlag),
		All:             ctx.Bool(AllFlag),
		IgnoreCase:      ctx.Bool(IgnoreCaseFlag),
		Backup:          ctx.Bool(BackupFlag) || ctx.String(BackupExtensionFlag) != "",
		BackupExtension: ctx.String(BackupExtensionFlag),
		ShowDiff:        ctx.Bool(ShowDiffFlag),
		Interactive:     ctx.Bool(InteractiveFlag),
		Quiet:           ctx.Bool(QuietFlag),
		Verbose:         ctx.Bool(VerboseFlag),
		Null:            ctx.Bool(NullFlag),
		AddFile:         ctx.StringSlice(FileFlag),
		Context:         ctx,
		Argv:            ctx.Args().Slice(),
		Argc:            ctx.NArg(),
		Notifier:        notifier,
	}

	if err = w.init(ctx); err != nil {
		return
	}

	if w.Argc < 2 {
		cli.ShowAppHelpAndExit(ctx, 1)
		eventFlag = cenums.EVENT_STOP
		return
	}

	w.Search, w.Replace = w.Argv[0], w.Argv[1]
	if w.Argc > 2 {
		if w.Stdin = w.Argv[2] == "-"; w.Stdin {
			w.Paths = w.Argv[3:]
		} else {
			w.Paths = w.Argv[2:]
		}
	}
	if len(w.Paths) == 0 && !ctx.IsSet(FileFlag) {
		w.Paths = []string{"."}
	}

	if w.Regex {
		if w.Pattern, err = MakeRegexp(w.Search, w); err != nil {
			err = fmt.Errorf("error compiling %q: %w", w.Search, err)
			eventFlag = cenums.EVENT_STOP
			return
		}
	}

	return
}

func (w *Worker) init(ctx *cli.Context) (err error) {

	if w.Exclude, err = ParseGlobs(ctx.StringSlice(ExcludeFlag)); err != nil {
		return
	} else if w.Include, err = ParseGlobs(ctx.StringSlice(IncludeFlag)); err != nil {
		return
	}

	if !w.All {
		more, _ := ParseGlobs([]string{"*~"})
		w.Exclude = append(w.Exclude, more...)
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
		// stdout is in use by curses, so if os.Stderr is piped, use that, otherwise write to a temp file
		var stat os.FileInfo
		if stat, err = os.Stderr.Stat(); err != nil {
			err = fmt.Errorf("error calling os.Stderr.Stat: %w", err)
			return
		} else if stat.Mode()&os.ModeNamedPipe != 0 {
			// user has piped the output of Stderr
			o = os.Stderr
		} else {
			// can't actually write to stderr because it's likely to the same terminal as curses
			if w.fwo, err = filewriter.New().UseTemp(TempOutPattern).Make(); err != nil {
				err = fmt.Errorf("error making new stdout temp file writer: %w", err)
				return
			}
			o = w.fwo
		}
		// stderr is always to a temp file when interactive
		if w.fwe, err = filewriter.New().UseTemp(TempErrPattern).Make(); err != nil {
			err = fmt.Errorf("error making new stderr temp file writer: %w", err)
			return
		}
		e = w.fwe
	} else {
		o = os.Stdout
		e = os.Stderr
	}

	w.Notifier.ModifyLevel(level).
		ModifyOut(o).
		ModifyErr(e)

	return
}

func (w *Worker) InitTargets(fn FindAllMatchingFn) (tooMany bool) {
	if fn == nil {
		fn = func(file string, matched bool, err error) {}
	}

	lookup := make(map[string]struct{})
	add := func(target string) (tooMany bool) {
		if abs, err := filepath.Abs(target); err == nil {
			if _, present := lookup[abs]; !present {
				if path.Exists(abs) {
					lookup[abs] = struct{}{}
					w.Targets = append(w.Targets, abs)
					//w.Notifier.Error("# does exist: %q\n", abs)
				} else {
					fn(abs, false, ErrNotFound)
					if w.Verbose {
						w.Notifier.Error("# does not exist: %q\n", abs)
					}
				}
			}
		} else {
			fn(target, false, err)
			if w.Verbose {
				w.Notifier.Error("# invalid path: %q: %s\n", abs, err)
			}
		}
		tooMany = len(w.Targets) > MaxFileCount
		return
	}

	// add any path arguments given
	for _, target := range w.Paths {
		add(target)
	}

	// scan and add any "additional files" given
	for _, target := range w.AddFile {
		var ee error
		if tooMany, ee = scanners.ScanFileLines(target, func(line string) (stop bool) {
			stop = add(line)
			return
		}); ee != nil {
			w.Notifier.Error("# error scanning --file %q: %v", target, ee)
		}
	}

	// scan and add from os.Stdin
	if w.Stdin {
		if w.Null {
			// using null-terminated paths
			tooMany = scanners.ScanNulls(os.Stdin, func(line string) (stop bool) {
				stop = add(line)
				return
			})
		} else {
			// using one path per line
			tooMany = scanners.ScanLines(os.Stdin, func(line string) (stop bool) {
				stop = add(line)
				return
			})
		}
	}

	if tooMany {
		w.Notifier.Error("# error: too many files, try batches of %d or less\n", MaxFileCount)
	}
	return
}

func (w *Worker) FindMatching(fn FindAllMatchingFn) {
	if w.Regex {
		w.Files, w.Matched = FindAllMatchingRegexp(w.Pattern, w.Targets, w.All, w.Include, w.Exclude, fn)
	} else if w.IgnoreCase {
		w.Files, w.Matched = FindAllMatchingStringInsensitive(w.Search, w.Targets, w.All, w.Include, w.Exclude, fn)
	} else {
		w.Files, w.Matched = FindAllMatchingString(w.Search, w.Targets, w.All, w.Include, w.Exclude, fn)
	}
	return
}

func (w *Worker) StartIterating() (iter *Iterator) {
	if len(w.Matched) > 0 {
		iter = &Iterator{
			w:   w,
			pos: 0,
		}
	}
	return
}

func (w *Worker) FileWriterOut() (fwo filewriter.FileWriter) {
	fwo = w.fwo
	return
}

func (w *Worker) FileWriterErr() (fwe filewriter.FileWriter) {
	fwe = w.fwe
	return
}

func (w *Worker) String() (s string) {
	s += "Worker{"
	s += fmt.Sprintf("regex=%v;", w.Regex)
	s += fmt.Sprintf("dotMatchNl=%v;", w.DotMatchNl)
	s += fmt.Sprintf("recurse=%v;", w.Recurse)
	s += fmt.Sprintf("nop=%v;", w.Nop)
	s += fmt.Sprintf("all=%v;", w.All)
	s += fmt.Sprintf("ignoreCase=%v;", w.IgnoreCase)
	s += fmt.Sprintf("backup=%v;", w.Backup)
	s += fmt.Sprintf("backupExtension=%v;", w.BackupExtension)
	s += fmt.Sprintf("showDiff=%v;", w.ShowDiff)
	s += fmt.Sprintf("interactive=%v;", w.Interactive)
	s += fmt.Sprintf("quiet=%v;", w.Quiet)
	s += fmt.Sprintf("verbose=%v;", w.Verbose)
	s += fmt.Sprintf("files=%+v;", w.AddFile)
	s += fmt.Sprintf("exclude=%v;", w.Exclude.String())
	s += fmt.Sprintf("include=%v;", w.Include.String())
	s += "}"
	return
}
