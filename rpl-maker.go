// Copyright (c) 2024  The Go-Curses Authors
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
	"github.com/urfave/cli/v2"

	clcli "github.com/go-corelibs/cli"
	"github.com/go-corelibs/notify"
	"github.com/go-corelibs/slices"
	"github.com/go-curses/cdk/lib/enums"
)

func MakeWorker(ctx *cli.Context, notifier notify.Notifier) (w *Worker, eventFlag enums.EventFlag, err error) {
	w = &Worker{
		Regex:           ctx.Bool(RegexFlag.Name) || ctx.Bool(DotMatchNlFlag.Name) || ctx.Bool(MultiLineFlag.Name),
		MultiLine:       ctx.Bool(MultiLineFlag.Name),
		DotMatchNl:      ctx.Bool(DotMatchNlFlag.Name),
		Recurse:         ctx.Bool(RecurseFlag.Name),
		Nop:             ctx.Bool(NopFlag.Name),
		All:             ctx.Bool(AllFlag.Name),
		IgnoreCase:      ctx.Bool(IgnoreCaseFlag.Name),
		PreserveCase:    ctx.Bool(PreserveCaseFlag.Name),
		NoLimits:        ctx.Bool(NoLimitsFlag.Name),
		Backup:          ctx.Bool(BackupFlag.Name) || ctx.String(BackupExtensionFlag.Name) != "",
		BackupExtension: ctx.String(BackupExtensionFlag.Name),
		ShowDiff:        ctx.Bool(ShowDiffFlag.Name),
		Interactive:     ctx.Bool(InteractiveFlag.Name) || ctx.Bool(PauseFlag.Name),
		Pause:           ctx.Bool(PauseFlag.Name),
		Quiet:           ctx.Bool(QuietFlag.Name),
		Verbose:         ctx.Bool(VerboseFlag.Name),
		Null:            ctx.Bool(NullFlag.Name),
		AddFile:         ctx.StringSlice(FileFlag.Name),
		ExcludeArgs:     ctx.StringSlice(ExcludeFlag.Name),
		IncludeArgs:     ctx.StringSlice(IncludeFlag.Name),
		RelativePath:    ".",
		Argv:            ctx.Args().Slice(),
		Argc:            ctx.NArg(),
		Notifier:        notifier,
	}

	if w.Argc >= 2 {
		w.Search, w.Replace = w.Argv[0], w.Argv[1]
		if w.Argc > 2 {
			w.Argv = w.Argv[2:]
			if w.Stdin = slices.Within("-", w.Argv); w.Stdin {
				w.Argv = slices.Prune(w.Argv, "-")
			}
			w.Argc = len(w.Argv)
		}
	}
	w.Stdin = w.Stdin || w.Null

	if len(w.Paths) == 0 && !ctx.IsSet(FileFlag.Name) && !w.Stdin {
		// add CWD if no arguments and --file not present
		w.Paths = []string{"."}
	}

	if err = w.Init(); err != nil {
		return
	}

	if ctx.NArg() < 2 {
		if w.Verbose {
			clcli.ShowUsageOptionsAndExit(ctx, 1)
			return
		}
		clcli.ShowUsageAndExit(ctx, 1)
		return
	}

	return
}
