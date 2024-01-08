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

var (
	CliFlags = []cli.Flag{

		// configurations
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "nop",
			Aliases:  []string{"n"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "interactive",
			Aliases:  []string{"e"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "backup",
			Aliases:  []string{"b"},
		},
		&cli.StringFlag{
			Category: "Configuration",
			Name:     "backup-extension",
			Aliases:  []string{"B"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "show-diff",
			Aliases:  []string{"d"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "ignore-case",
			Aliases:  []string{"i"},
		},

		// file selection things
		&cli.BoolFlag{
			Category: "File Selection",
			Name:     "recurse",
			Aliases:  []string{"r"},
		},
		&cli.BoolFlag{
			Category: "File Selection",
			Name:     "all",
			Aliases:  []string{"a"},
		},
		&cli.BoolFlag{
			Category: "File Selection",
			Name:     "null",
			Aliases:  []string{"0"},
		},
		&cli.StringSliceFlag{
			Category: "File Selection",
			Name:     "file",
			Aliases:  []string{"f"},
		},
		&cli.StringSliceFlag{
			Category: "File Selection",
			Name:     "exclude",
			Aliases:  []string{"X"},
		},
		&cli.StringSliceFlag{
			Category: "File Selection",
			Name:     "include",
			Aliases:  []string{"I"},
		},

		// regular expressions
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "regex",
			Aliases:  []string{"p"},
		},
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "dot-match-nl",
			Aliases:  []string{"s"},
		},

		// general flags
		&cli.BoolFlag{
			Category: "General",
			Name:     "quiet",
			Aliases:  []string{"q"},
		},
		&cli.BoolFlag{
			Category: "General",
			Name:     "verbose",
			Aliases:  []string{"v"},
		},
	}
)
