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
	_ "embed"
	"fmt"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk"
	"github.com/go-curses/ctk/lib/enums"
)

//go:embed rpl.accelmap
var gAccelMap string

//go:embed rpl.styles
var gStyles string

func (u *CUI) makeWindowTitle() (title string) {
	title = fmt.Sprintf("%s (%s)", u.App.Name(), u.App.Version())
	add := func(v string) {
		if len(v) > 0 && v[0] == '-' {
			title += " " + v
		} else {
			title += fmt.Sprintf(" %q", v)
		}
	}
	var foundSearch bool
	for _, arg := range u.Args[1:] {
		if foundSearch {
			add(arg)
			break
		}
		foundSearch = arg == u.worker.Search
		add(arg)
	}
	return
}

func (u *CUI) startup(_ []interface{}, argv ...interface{}) cenums.EventFlag {
	var ok bool
	if _, u.Display, _, _, _, ok = ctk.ArgvApplicationSignalStartup(argv...); ok {

		ctk.GetAccelMap().LoadFromString(gAccelMap)

		u.Window = ctk.NewWindowWithTitle(u.makeWindowTitle())
		u.Window.SetName("rpl-window")
		if err := u.Window.ImportStylesFromString(gStyles); err != nil {
			u.LastError = fmt.Errorf("error importing window styles: %w", err)
			return cenums.EVENT_STOP
		}
		u.initAccelGroups()

		vbox := u.Window.GetVBox()
		vbox.SetSpacing(1)

		u.HeaderLabel = ctk.NewLabel("Starting up...")
		u.HeaderLabel.Show()
		u.HeaderLabel.SetUseMarkup(true)
		u.HeaderLabel.SetSizeRequest(-1, -1)
		u.HeaderLabel.SetAlignment(0.5, 0.5)
		u.HeaderLabel.SetJustify(cenums.JUSTIFY_CENTER)
		u.HeaderLabel.SetLineWrap(true)
		u.HeaderLabel.SetLineWrapMode(cenums.WRAP_CHAR)
		vbox.PackStart(u.HeaderLabel, false, false, 0)

		u.DiffView = ctk.NewScrolledViewport()
		u.DiffView.SetName("diff-view")
		u.DiffView.Show()
		u.DiffView.SetPolicy(enums.PolicyAutomatic, enums.PolicyAutomatic)
		vbox.PackStart(u.DiffView, true, true, 0)

		u.DiffLabel = ctk.NewLabel("")
		u.DiffLabel.SetName("diff-text")
		u.DiffLabel.Show()
		u.DiffLabel.SetUseMarkup(true)
		u.DiffLabel.SetSingleLineMode(false)
		u.DiffLabel.SetJustify(cenums.JUSTIFY_NONE)
		u.DiffLabel.SetLineWrap(false)
		u.DiffLabel.SetLineWrapMode(cenums.WRAP_NONE)
		u.DiffView.Add(u.DiffLabel)

		u.FooterLabel = ctk.NewLabel("")
		u.FooterLabel.Hide()
		u.FooterLabel.SetUseMarkup(true)
		u.FooterLabel.SetSizeRequest(-1, -1)
		u.FooterLabel.SetAlignment(0.5, 0.5)
		u.FooterLabel.SetJustify(cenums.JUSTIFY_CENTER)
		u.FooterLabel.SetLineWrap(true)
		u.FooterLabel.SetLineWrapMode(cenums.WRAP_CHAR)
		vbox.PackStart(u.FooterLabel, false, false, 0)

		workButtonsArea := ctk.NewHBox(false, 1)
		workButtonsArea.Show()
		workButtonsArea.SetSizeRequest(-1, 1)
		vbox.PackStart(workButtonsArea, false, false, 0)

		waLeftSep := ctk.NewSeparator()
		waLeftSep.Show()
		workButtonsArea.PackStart(waLeftSep, true, true, 0)

		mkButton := func(name, label, tooltip, handle string, fn func()) (b ctk.Button) {
			b = ctk.NewButtonWithMnemonic(label)
			b.SetName(name)
			b.Hide()
			if tooltip != "" {
				b.SetHasTooltip(true)
				b.SetTooltipText(tooltip)
			}
			b.Connect(ctk.SignalActivate, handle, func(data []interface{}, argv ...interface{}) cenums.EventFlag {
				b.LogDebug("clicked")
				fn()
				return cenums.EVENT_PASS
			})
			return
		}

		u.ContinueButton = mkButton("continue", ContinueAccelLabel, ContinueAccelTooltip, ContinueAccelHandle, u.startWork)
		u.ContinueButton.Hide()
		workButtonsArea.PackStart(u.ContinueButton, false, false, 0)

		u.SelectGroupsButton = mkButton("select-groups", SelectGroupsAccelLabel, SelectGroupsAccelTooltip, SelectGroupsAccelHandle, func() {
			u.WorkAccel.Activate(SelectGroupsAccelKey, 0)
		})
		u.SelectGroupsButton.Hide()
		workButtonsArea.PackStart(u.SelectGroupsButton, false, false, 0)

		u.SkipGroupButton = mkButton("skip-group", SkipGroupAccelLabel, SkipGroupAccelTooltip, SkipGroupAccelHandle, func() {
			u.WorkAccel.Activate(SkipGroupAccelKey, 0)
		})
		u.SkipGroupButton.Hide()
		workButtonsArea.PackStart(u.SkipGroupButton, false, false, 0)

		u.KeepGroupButton = mkButton("keep-group", KeepGroupAccelLabel, KeepGroupAccelTooltip, KeepGroupAccelHandle, func() {
			u.WorkAccel.Activate(KeepGroupAccelKey, 0)
		})
		u.KeepGroupButton.Hide()
		workButtonsArea.PackStart(u.KeepGroupButton, false, false, 0)

		u.SkipFileButton = mkButton("skip-file", SkipFileAccelLabel, SkipFileAccelTooltip, SkipFileAccelHandle, func() {
			u.WorkAccel.Activate(SkipFileAccelKey, 0)
		})
		u.SkipFileButton.Hide()
		workButtonsArea.PackStart(u.SkipFileButton, false, false, 0)

		u.SaveFileButton = mkButton("save-file", SaveFileAccelLabel, SaveFileAccelTooltip, SaveFileAccelHandle, func() {
			u.WorkAccel.Activate(SaveFileAccelKey, 0)
		})
		u.SaveFileButton.Hide()
		workButtonsArea.PackStart(u.SaveFileButton, false, false, 0)

		waRightSep := ctk.NewSeparator()
		waRightSep.Show()
		workButtonsArea.PackStart(waRightSep, true, true, 0)

		u.ActionArea = ctk.NewHButtonBox(false, 1)
		u.ActionArea.Show()
		u.ActionArea.SetSizeRequest(-1, 1)
		vbox.PackEnd(u.ActionArea, false, false, 0)

		u.StateSpinner = ctk.NewSpinner()
		u.StateSpinner.Show()
		u.StateSpinner.SetSizeRequest(1, 1)
		u.StateSpinner.StartSpinning()
		u.ActionArea.PackStart(u.StateSpinner, false, false, 0)

		u.StatusLabel = ctk.NewLabel("starting up...")
		u.StatusLabel.Show()
		u.StatusLabel.SetSizeRequest(-1, 1)
		u.StatusLabel.SetLineWrap(false)
		u.StatusLabel.SetLineWrapMode(cenums.WRAP_NONE)
		u.StatusLabel.SetSingleLineMode(true)
		u.ActionArea.PackStart(u.StatusLabel, true, true, 1)

		u.QuitButton = mkButton("quit", QuitAccelLabel, QuitAccelTooltip, QuitAccelHandle, u.requestQuit)
		u.QuitButton.Show()
		u.ActionArea.PackStart(u.QuitButton, false, false, 0)

		u.Window.Show()
		u.App.NotifyStartupComplete()
		u.Display.Connect(cdk.SignalEventResize, "ui-resize-handler", u.resize)
		cdk.Go(u.initWork)
		return cenums.EVENT_PASS
	}
	return cenums.EVENT_STOP
}
