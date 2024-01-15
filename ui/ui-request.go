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
	"github.com/go-curses/cdk/log"
)

func (u *CUI) requestEdit(argv ...interface{}) (handled bool) {
	if u.SelectGroupsButton.IsVisible() {
		log.DebugF("edit-accel called")
		u.SelectGroupsButton.GrabFocus()
		u.delta.SkipAll()
		u.group = -1
		u.processNextEdit()
	}
	return
}

func (u *CUI) requestSkipEdit(argv ...interface{}) (handled bool) {
	if u.SkipGroupButton.IsVisible() {
		log.DebugF("skip-edit-accel called")
		u.SkipGroupButton.GrabFocus()
		u.skipCurrentEdit()
		u.processNextEdit()
	}
	return
}

func (u *CUI) requestKeepEdit(argv ...interface{}) (handled bool) {
	if u.KeepGroupButton.IsVisible() {
		log.DebugF("keep-edit-accel called")
		u.KeepGroupButton.GrabFocus()
		u.keepCurrentEdit()
		u.processNextEdit()
	}
	return
}

func (u *CUI) requestSkip(argv ...interface{}) (handled bool) {
	if u.SkipButton.IsVisible() {
		log.DebugF("skip-accel called")
		u.SkipButton.GrabFocus()
		u.skipCurrentWork()
		u.processNextWork()
	}
	return
}

func (u *CUI) requestApply(argv ...interface{}) (handled bool) {
	if u.ApplyButton.IsVisible() {
		log.DebugF("apply-accel called")
		u.ApplyButton.GrabFocus()
		//u.ApplyButton.SetPressed(true)
		u.applyAndProcessNextWork()
	}
	return
}

func (u *CUI) requestQuit() {
	log.DebugDF(1, "requesting quit")
	u.QuitButton.GrabFocus()
	u.QuitButton.SetPressed(true)
	u.Display.RequestQuit()
}

func (u *CUI) requestDrawAndShow() {
	u.Window.Resize()
	u.Display.RequestDraw()
	u.Display.RequestShow()
}
