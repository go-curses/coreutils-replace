package main

import (
	"time"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
)

func initAccelGroups() {
	gWorkAccel = ctk.NewAccelGroup()
	gWorkAccel.ConnectByPath(
		"<rpl-window>/File/Edit",
		"edit-accel",
		func(argv ...interface{}) (handled bool) {
			if gEditButton.IsVisible() {
				log.DebugF("edit-accel called")
				gEditButton.GrabFocus()
				gDelta[gTargetIndex].SkipAll()
				gGroupsIndex = -1
				processNextEdit()
			}
			return
		},
	)
	gWorkAccel.ConnectByPath(
		"<rpl-window>/File/SkipEdit",
		"skip-edit-accel",
		func(argv ...interface{}) (handled bool) {
			if gSkipEditButton.IsVisible() {
				log.DebugF("skip-edit-accel called")
				gSkipEditButton.GrabFocus()
				skipCurrentEdit()
				processNextEdit()
			}
			return
		},
	)
	gWorkAccel.ConnectByPath(
		"<rpl-window>/File/KeepEdit",
		"keep-edit-accel",
		func(argv ...interface{}) (handled bool) {
			if gKeepEditButton.IsVisible() {
				log.DebugF("keep-edit-accel called")
				gKeepEditButton.GrabFocus()
				keepCurrentEdit()
				processNextEdit()
			}
			return
		},
	)
	gWorkAccel.ConnectByPath(
		"<rpl-window>/File/Skip",
		"skip-accel",
		func(argv ...interface{}) (handled bool) {
			if gSkipButton.IsVisible() {
				log.DebugF("skip-accel called")
				gSkipButton.GrabFocus()
				skipCurrentWork()
				processNextWork()
			}
			return
		},
	)
	gWorkAccel.ConnectByPath(
		"<rpl-window>/File/Keep",
		"next-accel",
		func(argv ...interface{}) (handled bool) {
			if gKeepButton.IsVisible() {
				log.DebugF("next-accel called")
				gKeepButton.GrabFocus()
				processNextWork()
			}
			return
		},
	)
	gWorkAccel.ConnectByPath(
		"<rpl-window>/File/Apply",
		"accept-accel",
		func(argv ...interface{}) (handled bool) {
			if gApplyButton.IsVisible() {
				log.DebugF("accept-accel called")
				gApplyButton.GrabFocus()
				gApplyButton.SetPressed(true)
				cdk.AddTimeout(time.Millisecond*100, func() enums.EventFlag {
					gDisplay.RequestQuit()
					return enums.EVENT_STOP
				})
			}
			return
		},
	)
	gWindow.AddAccelGroup(gWorkAccel)

	ag := ctk.NewAccelGroup()
	ag.ConnectByPath(
		"<rpl-window>/File/Quit",
		"quit-accel",
		func(argv ...interface{}) (handled bool) {
			ag.LogDebug("quit-accel called")
			gAbortWork = true
			requestQuit()
			return
		},
	)
	ag.ConnectByPath(
		"<rpl-window>/File/Exit",
		"ctrl-c-accel",
		func(argv ...interface{}) (handled bool) {
			ag.LogDebug("ctrl-c-accel called")
			gAbortWork = true
			requestQuit()
			return
		},
	)
	gWindow.AddAccelGroup(ag)
}

func requestQuit() {
	log.DebugDF(1, "requesting quit")
	gQuitButton.GrabFocus()
	gQuitButton.SetPressed(true)
	cdk.AddTimeout(time.Millisecond*100, func() enums.EventFlag {
		gDisplay.RequestQuit()
		return enums.EVENT_STOP
	})
}