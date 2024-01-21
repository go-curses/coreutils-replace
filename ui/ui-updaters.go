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

package ui

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-corelibs/diff"
	"github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/ctk/lib/enums"
)

func (u *CUI) getChangesText() (text string) {

	if u.delta != nil {
		if count := u.delta.EditGroupsLen(); count > 1 {
			if u.view == SelectGroupsView {
				text = fmt.Sprintf("%d of %d groups of changes", u.group+1, count)
			} else {
				text = fmt.Sprintf("%d groups of changes", count)
			}
		} else if count == 1 {
			text = "one group of changes"
		}
	}

	return
}

func (u *CUI) getSearchText(prefix string) (text string) {
	if u.worker.Regex {
		text = fmt.Sprintf("%s the pattern: %q", prefix, u.worker.Search)
	} else if u.worker.IgnoreCase {
		text = fmt.Sprintf("%s the relative text: %q", prefix, u.worker.Search)
	} else {
		text = fmt.Sprintf("%s the extact text: %q", prefix, u.worker.Search)
	}
	return
}

func (u *CUI) setFocusLabels(focused bool) {
	vbox := u.Window.GetVBox()
	if focused {
		u.DiffView.Hide()
		vbox.SetChildPacking(u.HeaderLabel, true, true, 0, enums.PackStart)
	} else {
		u.DiffView.Show()
		vbox.SetChildPacking(u.HeaderLabel, false, false, 0, enums.PackStart)
	}
}

func (u *CUI) setHeaderLabel(text string) {
	if text == "" {
		u.HeaderLabel.SetLabel("")
		u.HeaderLabel.Hide()
	} else {
		u.HeaderLabel.SetLabel(text)
		u.HeaderLabel.Show()
	}
}

func (u *CUI) setDiffLabel(text string, markup bool) {
	if markup {
		u.DiffLabel.SetUseMarkup(true)
		if err := u.DiffLabel.SetMarkup(text); err != nil {
			u.DiffLabel.LogErr(err)
		}
	} else {
		u.DiffLabel.SetUseMarkup(false)
		u.DiffLabel.SetLabel(text)
	}
	u.DiffLabel.SetSizeRequest(u.DiffLabel.GetPlainTextInfo())
	u.DiffLabel.Resize()
	u.DiffView.Resize()
}

func (u *CUI) setDiffPatch(unified string) {
	markup := diff.TangoRender.RenderDiff(unified)
	u.setDiffLabel(markup, true)
	u.DiffView.ScrollTop()
}

func (u *CUI) setFooterLabel(text string) {
	if text == "" {
		u.FooterLabel.SetLabel("")
		u.FooterLabel.Hide()
	} else {
		u.FooterLabel.SetLabel(text)
		u.FooterLabel.Show()
	}
}

func (u *CUI) setStatusLabel(text string) {
	if text == "" {
		u.StatusLabel.SetLabel("")
		//u.StatusLabel.Hide()
	} else {
		u.StatusLabel.SetLabel(text)
		//u.StatusLabel.Show()
	}
	return
}

func (u *CUI) setStateSpinner(active bool) {
	if active {
		u.StateSpinner.Show()
		u.StateSpinner.StartSpinning()
	} else {
		u.StateSpinner.StartSpinning()
		u.StateSpinner.Hide()
	}
}

func (u *CUI) updateInitWorkStatus(maxLen int, r cFindResult) {
	u.results = append(u.results, r)
	u.setStatusLabel(r.Status(maxLen))
	u.setDiffLabel(u.results.Tango(), true)
	u.DiffView.ScrollBottom()
	u.requestDrawAndShow()
	time.Sleep(time.Millisecond * 10)
}

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
	u.setStatusLabel(status + " " + name)
	u.ActionArea.Resize()
	u.requestDrawAndShow()
}

func (u *CUI) updateFileWorkStatus() {
	if u.iter == nil {
		u.setStatusLabel("")
		return
	}
	//u.setHeaderLabel("Viewing changes in: " + filepath.Base(u.iter.Name()))
	u.setHeaderLabel("")
	u.updateStatusLine()
}

func (u *CUI) updateEditWorkStatus() {
	if u.iter == nil || u.delta == nil {
		u.setStatusLabel("")
		return
	}
	u.setFooterLabel(u.getChangesText())
	u.updateStatusLine()
}
