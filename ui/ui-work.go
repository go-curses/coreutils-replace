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

	"github.com/go-curses/cdk/lib/math"
)

func (u *CUI) initWork() {
	u.setStateSpinner(true)
	u.DiffLabel.Show()
	u.DiffView.Show()

	w, _ := u.Display.Screen().Size()
	maxLen := math.FloorI((w/2)-2, 10)
	if u.LastError = u.worker.InitTargets(nil); u.LastError != nil {
		u.requestQuit()
		return
	}
	if u.LastError = u.worker.FindMatching(func(file string, matched bool, err error) {
		u.updateInitWorkStatus(maxLen, cFindResult{
			target:  file,
			matched: matched,
			err:     err,
		})
	}); u.LastError != nil {
		u.requestQuit()
		return
	}

	var count int
	if count = len(u.worker.Matched); count == 0 {
		u.setHeaderLabel(u.getSearchText("no files contain"))
		u.setFocusLabels(true)
	} else {
		if count == 1 {
			u.setHeaderLabel(u.getSearchText("one file contains"))
		} else {
			u.setHeaderLabel(u.getSearchText(fmt.Sprintf("%d files contain", count)))
		}
		u.setFooterLabel("(-) no matches; (+) has matches; (x) errors")
	}

	u.setStatusLabel("")
	u.setStateSpinner(false)

	if u.worker.Pause {
		if count > 0 {
			u.ContinueButton.Show()
			u.ContinueButton.GrabFocus()
		}
		u.requestDrawAndShow()
	} else {
		u.startWork()
	}
}

func (u *CUI) startWork() {
	u.setFooterLabel("")
	u.ContinueButton.Hide()
	u.setDiffLabel("", false)
	u.DiffView.ScrollTop()

	u.iter = u.worker.StartIterating()
	if len(u.worker.Matched) > 0 {
		// work to do
		u.processNextFile()
		u.DiffView.GrabFocus()
		return
	}

	// no work to do
	if len(u.worker.Files) > 0 {
		u.setHeaderLabel(fmt.Sprintf("no files match search: %q", u.worker.Search))
	} else {
		u.setHeaderLabel("no files to search")
	}

	u.QuitButton.GrabFocus()
	u.requestDrawAndShow()
	return
}

func (u *CUI) saveFileAndProcessNextFile() {
	if u.iter != nil && u.delta != nil {
		if u.worker.Nop {
			u.notifier.Info(u.delta.UnifiedEdits())
		} else {
			if _, unified, backup, err := u.iter.ApplySpecific(u.delta); err != nil {
				u.notifier.Error("# error applying changes to %q: %v\n", u.iter.Name(), err)
			} else {
				if u.worker.Verbose && backup != "" {
					u.notifier.Error("# backed up %q to %q\n", u.iter.Name(), backup)
				}
				u.notifier.Info(unified)
			}
		}
		u.group = -1
		u.processNextFile()
	}
}

func (u *CUI) processNextFile() {
	if u.delta == nil {
		// just started working
	} else {
		// move to the next
		u.iter.Next()
	}

	if !u.iter.Valid() {
		// all done!
		u.requestQuit()
		return
	}

	var err error
	if _, _, u.count, u.delta, err = u.iter.Replace(); err != nil {
		u.notifier.Error(err.Error())
		u.processNextFile()
		return
	}

	u.delta.KeepAll()
	unified := u.delta.UnifiedEdits()
	u.setDiffPatch(unified)

	if count := u.delta.EditGroupsLen(); count > 1 {
		u.setFooterLabel("one group of changes")
	} else {
		u.setFooterLabel(fmt.Sprintf("%d groups of changes", count))
	}

	u.displayFileView()
}

func (u *CUI) keepCurrentGroup() {
	if u.iter != nil && u.delta != nil {
		if u.group >= 0 && u.group < u.delta.EditGroupsLen() {
			u.delta.KeepGroup(u.group)
		}
	}
}

func (u *CUI) skipCurrentGroup() {
	if u.iter != nil && u.delta != nil {
		if u.group >= 0 && u.group < u.delta.EditGroupsLen() {
			u.delta.SkipGroup(u.group)
		}
	}
}

func (u *CUI) startSelectingGroups() {
	u.delta.SkipAll()
	u.group = -1
	u.processNextGroup()
}

func (u *CUI) processNextGroup() {
	if u.iter != nil && u.delta != nil {
		u.group += 1
		if u.group < u.delta.EditGroupsLen() {
			unified := u.delta.EditGroup(u.group)
			u.setDiffPatch(unified)
			u.displayEditView()
		} else {
			unified := u.delta.UnifiedEdits()
			u.setDiffPatch(unified)
			u.displayFileView()
		}
	}
}

func (u *CUI) skipCurrentFile() {
	if u.iter != nil && u.delta != nil {
		u.delta.SkipAll()
	}
}

func (u *CUI) displayFileView() {
	u.view = FileView
	u.SkipGroupButton.Hide()
	u.KeepGroupButton.Hide()

	numMatched := len(u.worker.Matched)

	if numMatched > 1 {
		u.SkipFileButton.Show()
	} else {
		u.SkipFileButton.Hide()
	}

	numEditGroups := u.delta.EditGroupsLen()
	if numEditGroups > 0 {
		if numEditGroups > 1 {
			u.SelectGroupsButton.Show()
		} else {
			u.SelectGroupsButton.Hide()
		}
		u.SaveFileButton.Show()
	} else {
		u.SelectGroupsButton.Hide()
		u.SaveFileButton.Hide()
	}

	u.setFooterLabel(u.getChangesText())

	u.DiffView.ScrollTop()
	u.Window.Resize()
	u.updateFileWorkStatus()
}

func (u *CUI) displayEditView() {
	u.view = SelectGroupsView
	u.SkipFileButton.Hide()
	u.SelectGroupsButton.Hide()
	u.SaveFileButton.Hide()
	u.SkipGroupButton.Show()
	u.KeepGroupButton.Show()
	u.Window.Resize()
	u.updateEditWorkStatus()
}
