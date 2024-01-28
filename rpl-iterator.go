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
	"io"

	"github.com/go-corelibs/diff"
	"github.com/go-corelibs/path"
	"github.com/go-corelibs/replace"
)

type Iterator struct {
	w   *Worker
	pos int
}

func (i *Iterator) Pos() (pos int) {
	pos = i.pos
	return
}

func (i *Iterator) Next() {
	if i.Valid() {
		i.pos += 1
	}
	return
}

func (i *Iterator) Valid() (valid bool) {
	valid = i != nil && i.pos >= 0 && i.pos < len(i.w.Matched)
	return
}

func (i *Iterator) Name() (path string) {
	if i.Valid() {
		path = i.w.Matched[i.pos]
	}
	return
}

func (i *Iterator) Replace() (original, modified string, count int, delta *diff.Diff, err error) {
	if !i.Valid() {
		err = io.EOF
	} else if i.w.Pattern != nil {
		if i.w.PreserveCase {
			original, modified, count, delta, err = replace.RegexPreserveFile(i.w.Pattern, i.w.Replace, i.w.Matched[i.pos])
		} else if i.w.MultiLine {
			original, modified, count, delta, err = replace.RegexFile(i.w.Pattern, i.w.Replace, i.w.Matched[i.pos])
		} else {
			original, modified, count, delta, err = replace.RegexLinesFile(i.w.Pattern, i.w.Replace, i.w.Matched[i.pos])
		}
	} else if i.w.PreserveCase {
		original, modified, count, delta, err = replace.StringPreserveFile(i.w.Search, i.w.Replace, i.w.Matched[i.pos])
	} else if i.w.IgnoreCase {
		original, modified, count, delta, err = replace.StringInsensitiveFile(i.w.Search, i.w.Replace, i.w.Matched[i.pos])
	} else {
		original, modified, count, delta, err = replace.StringFile(i.w.Search, i.w.Replace, i.w.Matched[i.pos])
	}
	return
}

func (i *Iterator) ApplyAll() (count int, unified, backup string, err error) {
	if !i.Valid() {
		err = io.EOF
		return
	}
	var delta *diff.Diff
	if _, _, count, delta, err = i.Replace(); err == nil {
		delta.KeepAll()
		_, unified, backup, err = i.ApplySpecific(delta)
	}
	return
}

func (i *Iterator) ApplySpecific(delta *diff.Diff) (count int, unified, backup string, err error) {
	if !i.Valid() {
		err = io.EOF
		return
	}

	if count = delta.KeepLen(); count == 0 {
		// nop
		return
	}

	var modified string
	if modified, err = delta.ModifiedEdits(); err == nil {

		unified = delta.UnifiedEdits()

		var backupExtension, backupSeparator string
		if i.w.Backup {
			// TODO: figure out a template pattern for backup extension
			//       that isn't as cumbersome as text/template and also
			//       not as terse as fmt.Sprintf
			backupExtension = i.w.getBackupExtension()
			if i.w.BackupExtension != "" {
				backupSeparator = "."
			} else {
				backupSeparator = DefaultBackupSeparator
			}
		}

		if i.w.Nop {
			if i.w.Backup { // simulate backup filename
				for backup = path.BackupName(i.w.Matched[i.pos], backupExtension, backupSeparator); path.Exists(backup); {
					backup = path.BackupName(backup, backupExtension, backupSeparator)
				}
			}
		} else if i.w.Backup {
			backup, err = path.BackupAndOverwrite(i.w.Matched[i.pos], modified, backupExtension, backupSeparator)
		} else {
			err = path.Overwrite(i.w.Matched[i.pos], modified)
		}

	}

	return
}
