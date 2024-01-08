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
	"io"
	"regexp"

	"github.com/wneessen/go-fileperm"

	"github.com/go-corelibs/diff"
	"github.com/go-corelibs/path"
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
	if count := len(i.w.Matched); (i.pos + 1) < count {
		i.pos += 1
	} else if i.pos != count {
		i.pos = count
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

func (i *Iterator) Replace() (original, modified string, delta *diff.Diff, err error) {
	if !i.Valid() {
		err = io.EOF
		return
	}
	if i.w.Pattern == nil {
		if i.w.IgnoreCase {
			original, modified, delta, err = ProcessTargetStringInsensitive(i.w.Search, i.w.Replace, i.w.Matched[i.pos])
		} else {
			original, modified, delta, err = ProcessTargetString(i.w.Search, i.w.Replace, i.w.Matched[i.pos])
		}
	} else {
		var rx *regexp.Regexp
		if rx, err = MakeRegexp(i.w.Search, i.w); err != nil {
			return
		}
		original, modified, delta, err = ProcessTargetRegex(rx, i.w.Replace, i.w.Matched[i.pos])
	}
	return
}

func (i *Iterator) Apply() (count int, unified, backup string, err error) {
	if !i.Valid() {
		err = io.EOF
		return
	}
	var delta *diff.Diff
	if _, _, delta, err = i.Replace(); err != nil {
		return
	}
	delta.KeepAll()
	count, unified, backup, err = i.ApplyChanges(delta)
	return
}

func (i *Iterator) ApplyChanges(delta *diff.Diff) (count int, unified, backup string, err error) {
	if !i.Valid() {
		err = io.EOF
		return
	}

	if count = delta.KeepLen(); count == 0 {
		// nop
		return
	}

	var modified string
	if modified, err = delta.ModifiedEdits(); err != nil {
		return
	}
	unified = delta.UnifiedEdits()

	var backupExtension, backupSeparator string
	if i.w.Backup {
		if i.w.BackupExtension != "" {
			backupSeparator = "."
			backupExtension = i.w.BackupExtension
		} else {
			backupSeparator = DefaultBackupSeparator
			backupExtension = DefaultBackupExtension
		}
	}

	var fp fileperm.PermUser
	if fp, err = fileperm.New(i.w.Matched[i.pos]); err != nil {
		return
	} else if !fp.UserWritable() {
		err = errors.New("file write permission denied")
		return
	} else if i.w.Nop {
		if i.w.Backup { // simulate backup filename
			for backup = path.BackupName(i.w.Matched[i.pos], backupExtension, backupSeparator); path.Exists(backup); {
				backup = path.BackupName(backup, backupExtension, backupSeparator)
			}
		}
		return
	}

	if i.w.Backup {
		backup, err = path.BackupAndOverwrite(i.w.Matched[i.pos], modified, backupExtension, backupSeparator)
		return
	}

	err = path.Overwrite(i.w.Matched[i.pos], modified)
	return
}
