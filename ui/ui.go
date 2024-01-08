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
	"sort"
	"sync"

	"github.com/urfave/cli/v2"

	"github.com/go-corelibs/diff"
	"github.com/go-corelibs/notify"
	"github.com/go-curses/cdk"
	"github.com/go-curses/ctk"

	"github.com/go-curses/coreutils-replace"
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

	Display        cdk.Display
	Window         ctk.Window
	MainLabel      ctk.Label
	DiffView       ctk.ScrolledViewport
	DiffLabel      ctk.Label
	WorkAccel      ctk.AccelGroup
	EditButton     ctk.Button
	KeepEditButton ctk.Button
	SkipEditButton ctk.Button
	SkipButton     ctk.Button
	ApplyButton    ctk.Button
	QuitButton     ctk.Button

	ActionArea ctk.HButtonBox

	StateSpinner ctk.Spinner
	StatusLabel  ctk.Label

	LastError error

	notifier notify.Notifier
	worker   *replace.Worker
	iter     *replace.Iterator
	delta    *diff.Diff
	group    int

	view ViewType

	sync.RWMutex
}

func NewUI(name, usage, description, version, release, tag, title, ttyPath string, notifier notify.Notifier) (u *CUI) {

	u = &CUI{
		App:      ctk.NewApplication(name, usage, description, version, tag, title, ttyPath),
		notifier: notifier,
	}
	c := u.App.CLI()
	c.Version = version + " (" + release + ")"
	c.ArgsUsage = ""
	c.UsageText = name + " [options] <search> <replace> <path> [path...]"
	c.HideHelpCommand = true
	c.EnableBashCompletion = true
	c.UseShortOptionHandling = true
	c.CustomAppHelpTemplate = gAppHelpTemplate

	cli.HelpFlag = &cli.BoolFlag{
		Category: "General",
		Name:     "help",
		Usage:    "display command-line usage information",
		Aliases:  []string{"h"},
	}
	cli.VersionFlag = &cli.BoolFlag{
		Category: "General",
		Name:     "version",
		Usage:    "display the version",
		Aliases:  []string{"V"},
	}

	c.Flags = append(c.Flags, replace.CliFlags...)
	sort.Sort(cli.FlagsByName(c.Flags))

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