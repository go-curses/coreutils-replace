package main

import (
	_ "embed"
	"strings"
	"time"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/ctk"
	"github.com/go-curses/ctk/lib/enums"
)

/*
	- can quit at any time
	- for each path given:
		- recurse and gather dictionary of files
		- generating TextEdit lists associated with each file
	- for each file found:
		- display a diff of the whole file
			- prompt to accept all or walk through
		- if accept all, queue all edits, proceed next
		- if walk through:
			- display each edit, one at a time
			- prompt to keep or skip each edit
			- queue all edits kept, proceed next
	- apply all queued changes, recording all errors
	- report any errors during shutdown/exit STDERR
*/

var (
	gDisplay        cdk.Display
	gWindow         ctk.Window
	gMainLabel      ctk.Label
	gDiffView       ctk.ScrolledViewport
	gDiffLabel      ctk.Label
	gWorkAccel      ctk.AccelGroup
	gEditButton     ctk.Button
	gKeepEditButton ctk.Button
	gSkipEditButton ctk.Button
	gSkipButton     ctk.Button
	gKeepButton     ctk.Button
	gApplyButton    ctk.Button
	gQuitButton     ctk.Button
)

//go:embed rpl.accelmap
var rplAccelMap string

func startup(data []interface{}, argv ...interface{}) cenums.EventFlag {
	var ok bool
	if _, gDisplay, _, _, _, ok = ctk.ArgvApplicationSignalStartup(argv...); ok {

		ctk.GetAccelMap().LoadFromString(rplAccelMap)

		gWindow = ctk.NewWindowWithTitle(strings.Join(gArgv, " "))
		gWindow.SetName("rpl-window")
		// gWindow.Freeze()

		initAccelGroups()

		vbox := gWindow.GetVBox()
		vbox.SetSpacing(1)

		gMainLabel = ctk.NewLabel("Starting up...")
		gMainLabel.Show()
		gMainLabel.SetUseMarkup(true)
		gMainLabel.SetSingleLineMode(false)
		// gMainLabel.SetLineWrap(true)
		// gMainLabel.SetLineWrapMode(cenums.WRAP_WORD)
		gMainLabel.SetSizeRequest(-1, 2)
		gMainLabel.SetAlignment(0.5, 0.5)
		gMainLabel.SetJustify(cenums.JUSTIFY_CENTER)
		// _ = gMainLabel.SetBoolProperty(ctk.PropertyDebug, true)
		vbox.PackStart(gMainLabel, false, false, 0)

		workArea := ctk.NewHBox(false, 1)
		workArea.Show()
		workArea.SetSizeRequest(-1, 1)
		vbox.PackStart(workArea, false, false, 0)

		waLeftSep := ctk.NewSeparator()
		waLeftSep.Show()
		workArea.PackStart(waLeftSep, true, true, 0)

		gEditButton = ctk.NewButtonWithMnemonic("_Edit <F2>")
		gEditButton.Hide()
		gEditButton.SetHasTooltip(true)
		gEditButton.SetTooltipText("edit the selected changes")
		gEditButton.Connect(ctk.SignalActivate, "rpl-edit-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			gEditButton.LogDebug("clicked")
			gWorkAccel.Activate(cdk.KeyF2, 0)
			return cenums.EVENT_PASS
		})
		workArea.PackStart(gEditButton, false, false, 0)

		gSkipEditButton = ctk.NewButtonWithMnemonic("_Skip Edit <F3>")
		gSkipEditButton.Hide()
		gSkipEditButton.SetHasTooltip(true)
		gSkipEditButton.SetTooltipText("skip the changes presented below")
		gSkipEditButton.Connect(ctk.SignalActivate, "rpl-skip-edit-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			gSkipEditButton.LogDebug("clicked")
			gWorkAccel.Activate(cdk.KeyF3, 0)
			return cenums.EVENT_PASS
		})
		workArea.PackStart(gSkipEditButton, false, false, 0)

		gKeepEditButton = ctk.NewButtonWithMnemonic("_Keep Edit <F4>")
		gKeepEditButton.Hide()
		gKeepEditButton.SetHasTooltip(true)
		gKeepEditButton.SetTooltipText("keep the changes presented below")
		gKeepEditButton.Connect(ctk.SignalActivate, "rpl-keep-edit-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			gKeepEditButton.LogDebug("clicked")
			gWorkAccel.Activate(cdk.KeyF4, 0)
			return cenums.EVENT_PASS
		})
		workArea.PackStart(gKeepEditButton, false, false, 0)

		gSkipButton = ctk.NewButtonWithMnemonic("_Skip <F7>")
		gSkipButton.Hide()
		gSkipButton.SetHasTooltip(true)
		gSkipButton.SetTooltipText("skip the selected changes and proceed")
		gSkipButton.Connect(ctk.SignalActivate, "rpl-skip-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			gSkipButton.LogDebug("clicked")
			gWorkAccel.Activate(cdk.KeyF7, 0)
			return cenums.EVENT_PASS
		})
		workArea.PackStart(gSkipButton, false, false, 0)

		gKeepButton = ctk.NewButtonWithMnemonic("_Keep <F8>")
		gKeepButton.Hide()
		gKeepButton.SetHasTooltip(true)
		gKeepButton.SetTooltipText("accept the selected changes and proceed")
		gKeepButton.Connect(ctk.SignalActivate, "rpl-next-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			gKeepButton.LogDebug("clicked")
			gWorkAccel.Activate(cdk.KeyF8, 0)
			return cenums.EVENT_PASS
		})
		workArea.PackStart(gKeepButton, false, false, 0)

		gApplyButton = ctk.NewButtonWithMnemonic("_Apply <F9>")
		gApplyButton.Hide()
		gApplyButton.SetHasTooltip(true)
		gApplyButton.SetTooltipText("write the selected changes and exit")
		gApplyButton.Connect(ctk.SignalActivate, "rpl-apply-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			gApplyButton.LogDebug("clicked")
			gWorkAccel.Activate(cdk.KeyF9, 0)
			return cenums.EVENT_PASS
		})
		workArea.PackStart(gApplyButton, false, false, 0)

		waRightSep := ctk.NewSeparator()
		waRightSep.Show()
		workArea.PackStart(waRightSep, true, true, 0)

		gDiffView = ctk.NewScrolledViewport()
		gDiffView.Show()
		gDiffView.SetPolicy(enums.PolicyAutomatic, enums.PolicyAutomatic)
		vbox.PackStart(gDiffView, true, true, 0)

		gDiffLabel = ctk.NewLabel("")
		gDiffLabel.Show()
		gDiffLabel.SetUseMarkup(true)
		gDiffLabel.SetSingleLineMode(false)
		gDiffLabel.SetJustify(cenums.JUSTIFY_NONE)
		gDiffLabel.SetLineWrap(false)
		gDiffLabel.SetLineWrapMode(cenums.WRAP_NONE)
		diffTheme := gDiffLabel.GetTheme()
		diffTheme.Content.Normal = diffTheme.Content.Normal.Background(paint.ColorDarkBlue)
		diffTheme.Content.Prelight = diffTheme.Content.Prelight.Background(paint.ColorDarkBlue)
		diffTheme.Content.Active = diffTheme.Content.Active.Background(paint.ColorDarkBlue)
		diffTheme.Content.Selected = diffTheme.Content.Selected.Background(paint.ColorDarkBlue)
		diffTheme.Content.Insensitive = diffTheme.Content.Insensitive.Background(paint.ColorDarkBlue)
		gDiffLabel.SetTheme(diffTheme)
		gDiffView.SetTheme(diffTheme)
		// _ = gDiffLabel.SetBoolProperty(ctk.PropertyDebug, true)
		gDiffView.Add(gDiffLabel)

		actionArea := ctk.NewHButtonBox(false, 1)
		actionArea.Show()
		actionArea.SetSizeRequest(-1, 1)
		vbox.PackEnd(actionArea, false, false, 0)

		primaryActionSep := ctk.NewSeparator()
		primaryActionSep.Show()
		actionArea.PackStart(primaryActionSep, true, true, 0)

		secondaryActionSep := ctk.NewSeparator()
		secondaryActionSep.Show()
		actionArea.PackEnd(secondaryActionSep, true, true, 0)

		gQuitButton = ctk.NewButtonWithMnemonic("_Quit <F10>")
		gQuitButton.Show()
		gQuitButton.Connect(ctk.SignalActivate, "rpl-quit-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			requestQuit()
			return cenums.EVENT_PASS
		})
		actionArea.PackEnd(gQuitButton, false, false, 0)

		// gWindow.Thaw()
		gWindow.Show()
		gApp.NotifyStartupComplete()
		cdk.AddTimeout(time.Millisecond*10, startWorkProcess)
		return cenums.EVENT_PASS
	}
	return cenums.EVENT_STOP
}