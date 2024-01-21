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
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
)

// Accelerator Button Labels
const (
	ContinueAccelLabel     = "_Continue"
	SelectGroupsAccelLabel = "Select _Groups <F2>"
	SkipGroupAccelLabel    = "_Skip Group <F3>"
	KeepGroupAccelLabel    = "_Keep Group <F4>"
	SkipFileAccelLabel     = "_Skip File <F8>"
	SaveFileAccelLabel     = "Save _File <F9>"
	QuitAccelLabel         = "_Quit <F10>"
)

// Accelerator Keys
const (
	SelectGroupsAccelKey = cdk.KeyF2
	SkipGroupAccelKey    = cdk.KeyF3
	KeepGroupAccelKey    = cdk.KeyF4
	SkipFileAccelKey     = cdk.KeyF8
	SaveFileAccelKey     = cdk.KeyF9
)

// Accelerator Paths
const (
	SelectGroupsAccelPath = "<rpl-window>/File/SelectGroups"
	SkipGroupAccelPath    = "<rpl-window>/File/SkipGroup"
	KeepGroupAccelPath    = "<rpl-window>/File/KeepGroup"
	SkipFileAccelPath     = "<rpl-window>/File/SkipFile"
	SaveFileAccelPath     = "<rpl-window>/File/SaveFile"
	QuitAccelPath         = "<rpl-window>/File/Quit"
	ExitAccelPath         = "<rpl-window>/File/Exit"
)

// Accelerator Handles
const (
	SelectGroupsAccelHandle = "select-groups-accel"
	SkipGroupAccelHandle    = "skip-group-accel"
	KeepGroupAccelHandle    = "keep-group-accel"
	SkipFileAccelHandle     = "skip-file-accel"
	SaveFileAccelHandle     = "save-file-accel"
	QuitAccelHandle         = "quit-accel"
	ExitAccelHandle         = "ctrl-c-accel"
)

func (u *CUI) initAccelGroups() {
	u.WorkAccel = ctk.NewAccelGroup()
	u.WorkAccel.ConnectByPath(SelectGroupsAccelPath, SelectGroupsAccelHandle, u.accelSelectGroups)
	u.WorkAccel.ConnectByPath(SkipGroupAccelPath, SkipGroupAccelHandle, u.accelSkipGroup)
	u.WorkAccel.ConnectByPath(KeepGroupAccelPath, KeepGroupAccelHandle, u.accelKeepGroup)
	u.WorkAccel.ConnectByPath(SkipFileAccelPath, SkipFileAccelHandle, u.accelSkipFile)
	u.WorkAccel.ConnectByPath(SaveFileAccelPath, SaveFileAccelHandle, u.accelSaveFile)
	u.Window.AddAccelGroup(u.WorkAccel)

	ag := ctk.NewAccelGroup()
	ag.ConnectByPath(QuitAccelPath, QuitAccelHandle, u.quitAccel)
	ag.ConnectByPath(ExitAccelPath, ExitAccelHandle, u.exitAccel)
	u.Window.AddAccelGroup(ag)
}

func (u *CUI) reportAccel(handle string) {
	log.DebugDF(1, "accelerator invoked: "+handle)
}

func (u *CUI) quitAccel(_ ...interface{}) (handled bool) {
	u.reportAccel(QuitAccelHandle)
	u.requestQuit()
	return
}

func (u *CUI) exitAccel(_ ...interface{}) (handled bool) {
	u.reportAccel(ExitAccelHandle + " (ctrl+c)")
	u.requestQuit()
	return
}

func (u *CUI) accelSelectGroups(_ ...interface{}) (handled bool) {
	if u.SelectGroupsButton.IsVisible() {
		u.reportAccel(SelectGroupsAccelHandle)
		u.SelectGroupsButton.GrabFocus()
		u.startSelectingGroups()
	}
	return
}

func (u *CUI) accelSkipGroup(_ ...interface{}) (handled bool) {
	if u.SkipGroupButton.IsVisible() {
		u.reportAccel(SkipGroupAccelHandle)
		u.SkipGroupButton.GrabFocus()
		u.skipCurrentGroup()
		u.processNextGroup()
	}
	return
}

func (u *CUI) accelKeepGroup(_ ...interface{}) (handled bool) {
	if u.KeepGroupButton.IsVisible() {
		u.reportAccel(KeepGroupAccelHandle)
		u.KeepGroupButton.GrabFocus()
		u.keepCurrentGroup()
		u.processNextGroup()
	}
	return
}

func (u *CUI) accelSkipFile(_ ...interface{}) (handled bool) {
	if u.SkipFileButton.IsVisible() {
		u.reportAccel(SkipFileAccelHandle)
		u.SkipFileButton.GrabFocus()
		u.skipCurrentFile()
		u.processNextFile()
	}
	return
}

func (u *CUI) accelSaveFile(_ ...interface{}) (handled bool) {
	if u.SaveFileButton.IsVisible() {
		u.reportAccel(SaveFileAccelHandle)
		u.SaveFileButton.GrabFocus()
		u.saveFileAndProcessNextFile()
	}
	return
}
