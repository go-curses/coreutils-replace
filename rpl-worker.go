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
	"path/filepath"
	"regexp"

	"github.com/go-corelibs/filewriter"
	"github.com/go-corelibs/globs"
	"github.com/go-corelibs/notify"
	"github.com/go-corelibs/path"
	rpl "github.com/go-corelibs/replace"
	"github.com/go-corelibs/scanners"
	"github.com/go-corelibs/slices"
)

type Worker struct {
	Regex           bool
	MultiLine       bool
	DotMatchNl      bool
	Recurse         bool
	Nop             bool
	All             bool
	IgnoreCase      bool
	PreserveCase    bool
	BinAsText       bool
	RelativePath    string
	Backup          bool
	BackupExtension string
	NoLimits        bool
	ShowDiff        bool
	Interactive     bool
	Pause           bool
	Quiet           bool
	Verbose         bool

	Argv []string
	Argc int

	Search      string
	Pattern     *regexp.Regexp
	Replace     string
	Stdin       bool
	Null        bool
	AddFile     []string
	Include     globs.Globs
	IncludeArgs []string
	Exclude     globs.Globs
	ExcludeArgs []string

	Paths   []string
	Targets []string
	Files   []string
	Matched []string

	Notifier notify.Notifier

	fwo filewriter.FileWriter
	fwe filewriter.FileWriter

	initLookup map[string]struct{}
}

func (w *Worker) getBackupExtension() (extension string) {
	if extension = w.BackupExtension; extension != "" {
		return
	}
	extension = DefaultBackupExtension
	return
}

func (w *Worker) Init() (err error) {

	if w.Regex {
		if w.Pattern, err = rpl.MakeRegexp(w.Search, w.MultiLine, w.DotMatchNl, w.IgnoreCase); err != nil {
			err = fmt.Errorf("error compiling %q: %w", w.Search, err)
			return
		}
	}

	if w.Exclude, err = globs.Parse(w.ExcludeArgs...); err != nil {
		err = fmt.Errorf("--exclude %w", err)
		return
	} else if w.Include, err = globs.Parse(w.IncludeArgs...); err != nil {
		err = fmt.Errorf("--include %w", err)
		return
	}

	if !w.All {
		more, _ := globs.Parse("*" + w.getBackupExtension())
		w.Exclude = append(w.Exclude, more...)
	}

	if err = w.setupNotifier(); err == nil && w.NoLimits && !w.Quiet {
		w.Notifier.Error("%s\n", gNoLimitsWarning)
	}

	return
}

func (w *Worker) setupNotifier() (err error) {
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
		// stdout is in use by curses, so the idea of stderr becomes the
		// notifier stdout with the notifier stderr always a temp file; if the
		// actual os.Stderr is a named pipe, use that for the notifier stdout,
		// otherwise use a temp file. Using the named pipe mode allows the
		// diffs to be output while the UI is doing stuff, which is nice
		var stat os.FileInfo
		if stat, err = os.Stderr.Stat(); err == nil {
			if stat.Mode()&os.ModeNamedPipe != 0 {
				// user has piped the output of Stderr
				o = os.Stderr
			}
		}

		if o == nil {
			if w.fwo, err = filewriter.New().UseTemp(TempOutPattern).Make(); err != nil {
				err = fmt.Errorf("stdout filewriter error: %w", err)
				return
			}
			o = w.fwo
		}

		// stderr is always to a temp file when interactive
		if w.fwe, err = filewriter.New().UseTemp(TempErrPattern).Make(); err != nil {
			err = fmt.Errorf("stderr filewriter error: %w", err)
			return
		}
		e = w.fwe

	} else {
		o = os.Stdout
		e = os.Stderr
	}

	if w.Notifier == nil {
		w.Notifier = notify.New(level).
			SetOut(o).
			SetErr(e).
			Make()
	} else {
		w.Notifier.
			ModifyLevel(level).
			ModifyOut(o).
			ModifyErr(e)
	}

	return
}

func (w *Worker) addTargetFile(target string) (tooMany bool, err error) {
	var resolved string

	if resolved, err = filepath.Abs(target); err != nil {
		err = fmt.Errorf("%q - %w", target, err)
		return
	}

	if w.RelativePath != "" {
		if w.RelativePath == "." {
			if w.RelativePath, err = os.Getwd(); err != nil {
				err = fmt.Errorf("%q - %w", target, err)
				return
			}
		}
		if resolved, err = filepath.Rel(w.RelativePath, resolved); err != nil {
			err = fmt.Errorf("%q - %w", target, err)
			return
		}
	}

	if _, present := w.initLookup[resolved]; !present {
		if path.Exists(resolved) {
			w.initLookup[resolved] = struct{}{}
			w.Targets = append(w.Targets, resolved)
		} else {
			err = fmt.Errorf("%w: %q", ErrNotFound, resolved)
			return
		}
	}
	tooMany = !w.NoLimits && len(w.Targets) > rpl.MaxFileCount
	return
}

func (w *Worker) scanTargetFn(line string) (stop bool) {
	var eee error
	if stop, eee = w.addTargetFile(line); eee != nil {
		w.Notifier.Error("# error: %v\n", eee)
	}
	return
}

func (w *Worker) InitTargets() (err error) {
	w.initLookup = make(map[string]struct{})

	// if not recursive, and "." is present, use the CWD files instead of "."
	if !w.Recurse && slices.Within(".", w.Paths) {
		w.Paths = slices.Prune(w.Paths, ".")
		var files []string
		if files, err = path.ListFiles(".", w.All); err != nil {
			return
		}
		w.Paths = append(w.Paths, files...)
	}

	// add any path arguments given
	for _, target := range w.Paths {
		if tooMany, ee := w.addTargetFile(target); tooMany {
			return ErrTooManyFiles
		} else if ee != nil {
			if w.Verbose {
				w.Notifier.Error("# error: %v\n", ee)
			}
		}
	}

	// scan and add any "additional files" given
	for _, target := range w.AddFile {
		if stopped, ee := scanners.ScanFileLines(target, w.scanTargetFn); ee != nil {
			w.Notifier.Error("# error scanning --file %q: %v", target, ee)
		} else if stopped {
			return ErrTooManyFiles
		}
	}

	if w.Stdin {
		// scan and add from os.Stdin
		if w.Null {
			// using null-terminated paths
			if scanners.ScanNulls(os.Stdin, w.scanTargetFn) {
				return ErrTooManyFiles
			}
		} else {
			// using one path per line
			if scanners.ScanLines(os.Stdin, w.scanTargetFn) {
				return ErrTooManyFiles
			}
		}
	}

	// free memory
	w.initLookup = nil
	return
}

func (w *Worker) FindMatching(fn rpl.FindAllMatchingFn) (err error) {
	if w.Regex {
		if w.MultiLine {
			w.Files, w.Matched, err = rpl.FindAllMatchingRegexp(w.Pattern, w.Targets, w.All, w.NoLimits, w.BinAsText, w.Recurse, w.Include, w.Exclude, fn)
		} else {
			w.Files, w.Matched, err = rpl.FindAllMatchingRegexpLines(w.Pattern, w.Targets, w.All, w.NoLimits, w.BinAsText, w.Recurse, w.Include, w.Exclude, fn)
		}
	} else if w.PreserveCase || w.IgnoreCase {
		w.Files, w.Matched, err = rpl.FindAllMatchingStringInsensitive(w.Search, w.Targets, w.All, w.NoLimits, w.BinAsText, w.Recurse, w.Include, w.Exclude, fn)
	} else {
		w.Files, w.Matched, err = rpl.FindAllMatchingString(w.Search, w.Targets, w.All, w.NoLimits, w.BinAsText, w.Recurse, w.Include, w.Exclude, fn)
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
	s += fmt.Sprintf("multiLine=%v;", w.MultiLine)
	s += fmt.Sprintf("dotMatchNl=%v;", w.DotMatchNl)
	s += fmt.Sprintf("recurse=%v;", w.Recurse)
	s += fmt.Sprintf("nop=%v;", w.Nop)
	s += fmt.Sprintf("all=%v;", w.All)
	s += fmt.Sprintf("ignoreCase=%v;", w.IgnoreCase)
	s += fmt.Sprintf("preserveCase=%v;", w.PreserveCase)
	s += fmt.Sprintf("binAsText=%v;", w.BinAsText)
	s += fmt.Sprintf("noLimits=%v;", w.NoLimits)
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
