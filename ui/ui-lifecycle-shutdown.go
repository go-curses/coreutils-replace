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

package ui

import (
	"fmt"
	"os"
	"strings"

	cenums "github.com/go-curses/cdk/lib/enums"
)

// shutdown happens after the curses display screen is closed and the display itself shutdown, it is safe to use stdout
// and stderr for normal printing
func (u *CUI) shutdown(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if u.LastError != nil {
		u.notifier.Error("# error: %v\n", strings.TrimSuffix(u.LastError.Error(), "\n"))
		return cenums.EVENT_PASS
	}

	if u.worker.Interactive {
		// all work completed already, just need to output the file writers
		if o := u.worker.FileWriterOut(); o != nil {
			o.WalkFile(func(line string) (stop bool) {
				_, _ = fmt.Fprintf(os.Stderr, line+"\n")
				return
			})
			_ = o.Remove()
		}
		if e := u.worker.FileWriterErr(); e != nil {
			e.WalkFile(func(line string) (stop bool) {
				_, _ = fmt.Fprintf(os.Stdout, line+"\n")
				return
			})
			_ = e.Remove()
		}
		return cenums.EVENT_PASS
	}

	if tooMany := u.worker.InitTargets(nil); tooMany {
		return cenums.EVENT_PASS
	}

	u.worker.FindMatching(nil)

	if u.worker.Verbose {
		var format string
		if u.worker.Nop {
			format = "# [nop] would replace"
		} else {
			format = "# replacing"
		}
		if u.worker.Regex {
			format += " /%s/"
		} else {
			format += " %q"
		}
		format += " with %q in %d of %d files\n"
		u.notifier.Error(
			format,
			u.worker.Search,
			u.worker.Replace,
			len(u.worker.Matched),
			len(u.worker.Files),
		)
	}

	for iter := u.worker.StartIterating(); iter.Valid(); iter.Next() {
		var count int
		var unified, backup string
		var err error
		if count, unified, backup, err = iter.Apply(); err != nil {
			u.notifier.Error("# %q error: %v\n", err)
			continue
		}

		if u.worker.Nop {
			if backup != "" {
				u.notifier.Error("# [nop] would have backed up %q to %q\n", iter.Name(), backup)
			}
			u.notifier.Error("# [nop] would have made %d edits to: %q\n", count, iter.Name())
		} else {
			if backup != "" {
				u.notifier.Error("# backed up %q to %q\n", iter.Name(), backup)
			}
			u.notifier.Error("# made %d edits to: %q\n", count, iter.Name())
		}

		if u.worker.ShowDiff {
			u.notifier.Info(unified)
		}
	}

	return cenums.EVENT_PASS
}
