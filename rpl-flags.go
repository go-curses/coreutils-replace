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
	DefaultBackupExtension = ".bak"
)

var (
	CliFlags = []cli.Flag{

		// configurations
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "nop",
			Usage:    "report what would have otherwise been done",
			Aliases:  []string{"n"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "interactive",
			Usage:    "selectively apply changes per file",
			Aliases:  []string{"I"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "backup",
			Usage:    "make backups before replacing content",
			Aliases:  []string{"b"},
		},
		&cli.StringFlag{
			Category: "Configuration",
			Name:     "bak",
			Usage:    "specify the backup file suffix to use (implies -b)",
			Aliases:  []string{"B"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "show-diff",
			Usage:    "include unified diffs of all changes in the output",
			Aliases:  []string{"d"},
		},
		&cli.BoolFlag{
			Category: "Configuration",
			Name:     "ignore-case",
			Usage:    "perform a case-insensitive search (literal or regex)",
			Aliases:  []string{"i"},
		},

		// file selection things
		&cli.BoolFlag{
			Category: "File Selection",
			Name:     "recurse",
			Usage:    "recurse into sub-directories",
			Aliases:  []string{"r"},
		},
		&cli.BoolFlag{
			Category: "File Selection",
			Name:     "all",
			Usage:    `include files that start with a "."`,
			Aliases:  []string{"a"},
		},
		&cli.StringSliceFlag{
			Category: "File Selection",
			Name:     "x",
			Usage:    "exclude files matching glob patterns",
		},

		// regular expressions
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "regex",
			Usage:    "search and replace arguments are regular expressions",
			Aliases:  []string{"P"},
		},
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "multi-line",
			Usage:    "set the multi-line (?m) regexp flag (implies -P)",
			Aliases:  []string{"m"},
		},
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "dot-match-nl",
			Usage:    "set the dot-match-nl (?s) regexp flag (implies -P)",
			Aliases:  []string{"s"},
		},
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "ms",
			Usage:    "convenience flag to set -m and -s",
			Aliases:  []string{"p"},
		},
		&cli.BoolFlag{
			Category: "Expressions",
			Name:     "msi",
			Usage:    "convenience flag to set -m, -s and -i",
			Aliases:  []string{"pi"},
		},

		// general flags
		&cli.BoolFlag{
			Category: "General",
			Name:     "quiet",
			Usage:    "run silently, ignored if --nop is also used",
			Aliases:  []string{"q"},
		},
		&cli.BoolFlag{
			Category: "General",
			Name:     "verbose",
			Usage:    "run loudly, ignored if --quiet is also used",
			Aliases:  []string{"v"},
		},
	}
)