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
	"github.com/go-curses/ctk"
)

func (u *CUI) initAccelGroups() {
	u.WorkAccel = ctk.NewAccelGroup()

	u.WorkAccel.ConnectByPath(
		"<rpl-window>/File/Edit",
		"edit-accel",
		u.requestEdit,
	)
	u.WorkAccel.ConnectByPath(
		"<rpl-window>/File/SkipEdit",
		"skip-edit-accel",
		u.requestSkipEdit,
	)
	u.WorkAccel.ConnectByPath(
		"<rpl-window>/File/KeepEdit",
		"keep-edit-accel",
		u.requestKeepEdit,
	)
	u.WorkAccel.ConnectByPath(
		"<rpl-window>/File/Skip",
		"skip-accel",
		u.requestSkip,
	)
	u.WorkAccel.ConnectByPath(
		"<rpl-window>/File/Apply",
		"accept-accel",
		u.requestApply,
	)
	u.Window.AddAccelGroup(u.WorkAccel)

	ag := ctk.NewAccelGroup()
	ag.ConnectByPath(
		"<rpl-window>/File/Quit",
		"quit-accel",
		u.quitAccel,
	)
	ag.ConnectByPath(
		"<rpl-window>/File/Exit",
		"ctrl-c-accel",
		u.sigintAccel,
	)
	u.Window.AddAccelGroup(ag)
}

func (u *CUI) quitAccel(argv ...interface{}) (handled bool) {
	log.DebugF("quit-accel called")
	u.requestQuit()
	return
}

func (u *CUI) sigintAccel(argv ...interface{}) (handled bool) {
	log.DebugF("ctrl-c-accel called")
	u.requestQuit()
	return
}
