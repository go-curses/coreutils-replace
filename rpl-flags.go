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

package replace

import (
	"github.com/urfave/cli/v2"
)

var (
	DefaultBackupExtension = "~"
	DefaultBackupSeparator = "~"
)

const (
	NopFlag             = "nop"
	InteractiveFlag     = "interactive"
	BackupFlag          = "backup"
	BackupExtensionFlag = "backupExtension"
	ShowDiffFlag        = "showDiff"
	IgnoreCaseFlag      = "ignoreCase"
	RecurseFlag         = "recurse"
	AllFlag             = "all"
	NullFlag            = "null"
	FileFlag            = "file"
	ExcludeFlag         = "exclude"
	IncludeFlag         = "include"
	RegexFlag           = "regex"
	DotMatchNlFlag      = "dotMatchNl"
	QuietFlag           = "quiet"
	VerboseFlag         = "verbose"
)

var (
	CliFlags = []cli.Flag{

		// configurations
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     NopFlag,
			Aliases:  []string{"n"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     InteractiveFlag,
			Aliases:  []string{"e"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     BackupFlag,
			Aliases:  []string{"b"},
		},
		&cli.StringFlag{
			Category: "Configuration",
			Name:     BackupExtensionFlag,
			Aliases:  []string{"B"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     ShowDiffFlag,
			Aliases:  []string{"d"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     IgnoreCaseFlag,
			Aliases:  []string{"i"},
		},

		// file selection things
		&cli.BoolFlag{
			Category: "File Selection",
			Name:     RecurseFlag,
			Aliases:  []string{"r"},
		},
		&cli.BoolFlag{
			Category: "File Selection",
			Name:     AllFlag,
			Aliases:  []string{"a"},
		},
		&cli.BoolFlag{
			Category: "File Selection",
			Name:     NullFlag,
			Aliases:  []string{"0"},
		},
		&cli.StringSliceFlag{
			Category: "File Selection",
			Name:     FileFlag,
			Aliases:  []string{"f"},
		},
		&cli.StringSliceFlag{
			Category: "File Selection",
			Name:     ExcludeFlag,
			Aliases:  []string{"X"},
		},
		&cli.StringSliceFlag{
			Category: "File Selection",
			Name:     IncludeFlag,
			Aliases:  []string{"I"},
		},

		// regular expressions
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     RegexFlag,
			Aliases:  []string{"p"},
		},
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     DotMatchNlFlag,
			Aliases:  []string{"s"},
		},

		// general flags
		&cli.BoolFlag{
			Category: "General",
			Name:     QuietFlag,
			Aliases:  []string{"q"},
		},
		&cli.BoolFlag{
			Category: "General",
			Name:     VerboseFlag,
			Aliases:  []string{"v"},
		},
	}
)
