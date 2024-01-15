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
	"os"
	"sync"

	"github.com/urfave/cli/v2"

	"github.com/go-corelibs/diff"
	"github.com/go-corelibs/notify"
	"github.com/go-corelibs/slices"
	"github.com/go-curses/cdk"
	"github.com/go-curses/ctk"

	replace "github.com/go-curses/coreutils-replace"
)

//go:embed rpl.help.tmpl
var gAppHelpTemplate string

type ViewType uint8

const (
	NopeView ViewType = iota
	FileView
	EditView
)

type CUI struct {
	App  ctk.Application
	Args []string

	Display            cdk.Display
	Window             ctk.Window
	HeaderLabel        ctk.Label
	FooterLabel        ctk.Label
	DiffView           ctk.ScrolledViewport
	DiffLabel          ctk.Label
	WorkAccel          ctk.AccelGroup
	ContinueButton     ctk.Button
	SelectGroupsButton ctk.Button
	KeepGroupButton    ctk.Button
	SkipGroupButton    ctk.Button
	SkipButton         ctk.Button
	ApplyButton        ctk.Button
	QuitButton         ctk.Button

	ActionArea ctk.HButtonBox

	StateSpinner ctk.Spinner
	StatusLabel  ctk.Label

	LastError error

	notifier notify.Notifier
	worker   *replace.Worker
	iter     *replace.Iterator
	delta    *diff.Diff
	count    int
	group    int

	pause bool

	results cFindResults

	view ViewType

	sync.RWMutex
}

func NewUI(name, usage, description, manual, version, release, tag, title, ttyPath string, notifier notify.Notifier) (u *CUI) {

	u = &CUI{
		App:      ctk.NewApplication(name, usage, description, version, tag, title, ttyPath),
		notifier: notifier,
	}
	c := u.App.CLI()
	c.Version = version + " (" + release + ")"
	c.ArgsUsage = ""
	c.UsageText = name + " [options] <search> <replace> [path...]"
	c.HideHelpCommand = true
	c.EnableBashCompletion = true
	c.UseShortOptionHandling = true

	if slices.Within("--"+HelpFlag.Name, os.Args[1:]) {
		c.Description += manual
	} else {
		c.Description = ""
		NoLimitsFlag.Category = ""
		NoLimitsFlag.Hidden = true
	}

	cli.HelpFlag = UsageFlag
	cli.VersionFlag = VersionFlag

	c.Flags = append(c.Flags,
		BackupFlag,
		BackupExtensionFlag,
		IgnoreCaseFlag,
		PreserveCaseFlag,
		NopFlag,
		NoLimitsFlag,

		ShowDiffFlag,
		InteractiveFlag,
		PauseFlag,

		RecurseFlag,
		AllFlag,
		NullFlag,
		FileFlag,
		ExcludeFlag,
		IncludeFlag,

		RegexFlag,
		MultiLineFlag,
		DotMatchNlFlag,

		HelpFlag,
		QuietFlag,
		VerboseFlag,
	)

	u.App.Connect(cdk.SignalPrepareStartup, "ui-prepare-startup-handler", u.prepareStartup)
	u.App.Connect(cdk.SignalPrepare, "ui-prepare-handler", u.prepare)
	u.App.Connect(cdk.SignalStartup, "ui-startup-handler", u.startup)
	u.App.Connect(cdk.SignalShutdown, "ui-shutdown-handler", u.shutdown)
	return
}

func (u *CUI) Run(argv []string) (err error) {
	err = u.App.Run(argv)
	return
}
