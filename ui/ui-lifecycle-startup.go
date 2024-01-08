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
			title += " ..."
			break
		}
		foundSearch = arg == u.worker.Search
		add(arg)
	}
	return
}

func (u *CUI) startup(data []interface{}, argv ...interface{}) cenums.EventFlag {
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

		u.MainLabel = ctk.NewLabel("Starting up...")
		u.MainLabel.Show()
		u.MainLabel.SetUseMarkup(true)
		u.MainLabel.SetSizeRequest(-1, -1)
		u.MainLabel.SetAlignment(0.5, 0.5)
		u.MainLabel.SetJustify(cenums.JUSTIFY_CENTER)
		u.MainLabel.SetLineWrap(true)
		u.MainLabel.SetLineWrapMode(cenums.WRAP_CHAR)
		vbox.PackStart(u.MainLabel, false, false, 0)

		workButtonsArea := ctk.NewHBox(false, 1)
		workButtonsArea.Show()
		workButtonsArea.SetSizeRequest(-1, 1)
		vbox.PackStart(workButtonsArea, false, false, 0)

		waLeftSep := ctk.NewSeparator()
		waLeftSep.Show()
		workButtonsArea.PackStart(waLeftSep, true, true, 0)

		u.EditButton = ctk.NewButtonWithMnemonic("_Edit <F2>")
		u.EditButton.Hide()
		u.EditButton.SetHasTooltip(true)
		u.EditButton.SetTooltipText("edit the selected changes")
		u.EditButton.Connect(ctk.SignalActivate, "rpl-edit-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			u.EditButton.LogDebug("clicked")
			u.WorkAccel.Activate(cdk.KeyF2, 0)
			return cenums.EVENT_PASS
		})
		workButtonsArea.PackStart(u.EditButton, false, false, 0)

		u.SkipEditButton = ctk.NewButtonWithMnemonic("_Skip Change <F3>")
		u.SkipEditButton.Hide()
		u.SkipEditButton.SetHasTooltip(true)
		u.SkipEditButton.SetTooltipText("skip this group of changes")
		u.SkipEditButton.Connect(ctk.SignalActivate, "rpl-skip-edit-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			u.SkipEditButton.LogDebug("clicked")
			u.WorkAccel.Activate(cdk.KeyF3, 0)
			return cenums.EVENT_PASS
		})
		workButtonsArea.PackStart(u.SkipEditButton, false, false, 0)

		u.KeepEditButton = ctk.NewButtonWithMnemonic("_Keep Change <F4>")
		u.KeepEditButton.Hide()
		u.KeepEditButton.SetHasTooltip(true)
		u.KeepEditButton.SetTooltipText("keep this group of changes")
		u.KeepEditButton.Connect(ctk.SignalActivate, "rpl-keep-edit-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			u.KeepEditButton.LogDebug("clicked")
			u.WorkAccel.Activate(cdk.KeyF4, 0)
			return cenums.EVENT_PASS
		})
		workButtonsArea.PackStart(u.KeepEditButton, false, false, 0)

		u.SkipButton = ctk.NewButtonWithMnemonic("_Skip <F8>")
		u.SkipButton.Hide()
		u.SkipButton.SetHasTooltip(true)
		u.SkipButton.SetTooltipText("skip the file changes and proceed")
		u.SkipButton.Connect(ctk.SignalActivate, "rpl-skip-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			u.SkipButton.LogDebug("clicked")
			u.WorkAccel.Activate(cdk.KeyF8, 0)
			return cenums.EVENT_PASS
		})
		workButtonsArea.PackStart(u.SkipButton, false, false, 0)

		u.ApplyButton = ctk.NewButtonWithMnemonic("_Apply <F9>")
		u.ApplyButton.Hide()
		u.ApplyButton.SetHasTooltip(true)
		u.ApplyButton.SetTooltipText("write the file changes and proceed")
		u.ApplyButton.Connect(ctk.SignalActivate, "rpl-apply-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			u.ApplyButton.LogDebug("clicked")
			u.WorkAccel.Activate(cdk.KeyF9, 0)
			return cenums.EVENT_PASS
		})
		workButtonsArea.PackStart(u.ApplyButton, false, false, 0)

		waRightSep := ctk.NewSeparator()
		waRightSep.Show()
		workButtonsArea.PackStart(waRightSep, true, true, 0)

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

		//secondaryActionSep := ctk.NewSeparator()
		//secondaryActionSep.Show()
		//u.ActionArea.PackEnd(secondaryActionSep, true, true, 0)

		u.QuitButton = ctk.NewButtonWithMnemonic("_Quit <F10>")
		u.QuitButton.Show()
		u.QuitButton.Connect(ctk.SignalActivate, "rpl-quit-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			u.requestQuit()
			return cenums.EVENT_PASS
		})
		u.ActionArea.PackStart(u.QuitButton, false, false, 0)

		u.Window.Show()
		u.App.NotifyStartupComplete()
		u.Display.Connect(cdk.SignalEventResize, "ui-resize-handler", u.resize)
		cdk.Go(u.initWork)
		return cenums.EVENT_PASS
	}
	return cenums.EVENT_STOP
}
