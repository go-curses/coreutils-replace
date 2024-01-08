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
	"path/filepath"
	"time"

	"github.com/go-corelibs/diff"

	"github.com/go-curses/cdk/lib/math"
)

const (
	gFileMatch = `🗸`
	gFileSkip  = `🗴`
)

func (u *CUI) updateStatusLine() {
	status := fmt.Sprintf("%d/%d: ", u.iter.Pos()+1, len(u.worker.Matched))
	w, _ := u.Display.Screen().Size()
	alloc := u.QuitButton.GetAllocation()
	padding := 10
	maxLen := math.FloorI(w-len(status)-alloc.W-padding, 10)
	var name string
	file, _ := filepath.Abs(u.iter.Name())
	if size := len(file); size > maxLen {
		name = "..." + file[size-maxLen:]
	} else {
		name = file
	}
	u.StatusLabel.SetLabel(status + " " + name)
	u.ActionArea.Resize()
	u.requestDrawAndShow()
}

func (u *CUI) updateFileWorkStatus() {
	if u.iter == nil {
		u.StatusLabel.SetLabel("")
		return
	}
	u.MainLabel.SetLabel(filepath.Base(u.iter.Name()))
	u.updateStatusLine()
}

func (u *CUI) updateEditWorkStatus() {
	if u.iter == nil || u.delta == nil {
		u.StatusLabel.SetLabel("")
		return
	}
	u.MainLabel.SetLabel(fmt.Sprintf(
		"Change %d of %d in: %s",
		u.group+1, u.delta.EditGroupsLen(),
		filepath.Base(u.iter.Name()),
	))
	u.updateStatusLine()
}

func (u *CUI) initWork() {
	w, _ := u.Display.Screen().Size()
	maxLen := math.FloorI((w/2)-2, 10)
	u.worker.Init(func(file string, matched bool) {
		var name string
		if size := len(file); size > maxLen {
			name = "..." + file[size-maxLen:]
		} else {
			name = file
		}
		if matched {
			u.StatusLabel.SetLabel(gFileMatch + " " + name)
		} else {
			u.StatusLabel.SetLabel(gFileSkip + " " + name)
		}
		u.requestDrawAndShow()
		time.Sleep(time.Millisecond * 10)
	})
	u.StatusLabel.SetLabel("")
	u.StateSpinner.StopSpinning()
	u.StateSpinner.Hide()
	u.startWork()
}

func (u *CUI) startWork() {
	u.iter = u.worker.Start()
	if len(u.worker.Matched) > 0 {
		u.processNextWork()
	} else if len(u.worker.Files) > 0 {
		u.MainLabel.SetText(fmt.Sprintf("no files have %q to replace!", u.worker.Search))
	} else {
		u.MainLabel.SetText("no files found")
	}
	return
}

func (u *CUI) applyAndProcessNextWork() {
	if u.iter != nil && u.delta != nil {
		if u.worker.DryRun {
			u.notifier.Info(u.delta.UnifiedEdits())
		} else {
			if _, unified, err := u.iter.ApplyChanges(u.delta); err != nil {
				u.notifier.Error("# error applying changes to %q: %v", u.iter.Name(), err)
			} else {
				u.notifier.Info(unified)
			}
		}
		u.group = -1
		u.processNextWork()
	}
}

func (u *CUI) processNextWork() {
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
	if _, _, u.delta, err = u.iter.Replace(); err != nil {
		u.notifier.Error(err.Error())
		u.processNextWork()
		return
	}

	u.delta.KeepAll()
	unified := u.delta.UnifiedEdits()
	markup := diff.TangoRender.RenderDiff(unified)
	if err = u.DiffLabel.SetMarkup(markup); err != nil {
		u.DiffLabel.LogErr(err)
	}
	u.DiffLabel.SetSizeRequest(u.DiffLabel.GetPlainTextInfo())

	u.displayFileView()
}

func (u *CUI) keepCurrentEdit() {
	if u.iter != nil && u.delta != nil {
		if u.group >= 0 && u.group < u.delta.EditGroupsLen() {
			u.delta.KeepGroup(u.group)
		}
	}
}

func (u *CUI) skipCurrentEdit() {
	if u.iter != nil && u.delta != nil {
		if u.group >= 0 && u.group < u.delta.EditGroupsLen() {
			u.delta.SkipGroup(u.group)
		}
	}
}

func (u *CUI) processNextEdit() {
	if u.iter != nil && u.delta != nil {
		u.group += 1
		if u.group < u.delta.EditGroupsLen() {
			unified := u.delta.EditGroup(u.group)
			u.updateDiffLabel(unified)
			u.displayEditView()
		} else {
			unified := u.delta.UnifiedEdits()
			u.updateDiffLabel(unified)
			u.displayFileView()
		}
	}
}

func (u *CUI) updateDiffLabel(unified string) {
	markup := diff.TangoRender.RenderDiff(unified)
	if err := u.DiffLabel.SetMarkup(markup); err != nil {
		u.DiffLabel.LogErr(err)
	}
	u.DiffLabel.SetSizeRequest(u.DiffLabel.GetPlainTextInfo())
}

func (u *CUI) skipCurrentWork() {
	if u.iter != nil && u.delta != nil {
		u.delta.SkipAll()
	}
}

func (u *CUI) displayFileView() {
	u.view = FileView
	u.SkipEditButton.Hide()
	u.KeepEditButton.Hide()

	numMatched := len(u.worker.Matched)

	if numMatched > 1 {
		u.SkipButton.Show()
	} else {
		u.SkipButton.Hide()
	}

	numEditGroups := u.delta.EditGroupsLen()
	if numEditGroups > 0 {
		if numEditGroups > 1 {
			u.EditButton.Show()
		} else {
			u.EditButton.Hide()
		}
		u.ApplyButton.Show()
	} else {
		u.EditButton.Hide()
		u.ApplyButton.Hide()
	}

	u.Window.Resize()
	u.updateFileWorkStatus()
}

func (u *CUI) displayEditView() {
	u.view = EditView
	u.SkipButton.Hide()
	u.EditButton.Hide()
	u.ApplyButton.Hide()
	u.SkipEditButton.Show()
	u.KeepEditButton.Show()
	u.Window.Resize()
	u.updateEditWorkStatus()
}