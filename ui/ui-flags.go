// Copyright (c) 2024  The Go-Curses Authors
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
	"github.com/urfave/cli/v2"

	clcli "github.com/go-corelibs/cli"
	"github.com/go-curses/cdk"
)

const (
	CaseSensitivityCategory    = "1. Case Sensitivity"
	RegularExpressionsCategory = "2. Regular Expressions"
	UserInterfaceCategory      = "3. User Interface"
	BackupsCategory            = "4. Backups"
	TargetSelectionCategory    = "5. Target Selection"
	GeneralCategory            = "6. General"
	GoCursesCategory           = "7. Go-Curses"
)

var (
	BackupFlag = &cli.BoolFlag{Category: BackupsCategory,
		Name: "backup", Aliases: []string{"b"},
		Usage: "make backups before replacing content",
	}
	BackupExtensionFlag = &cli.StringFlag{Category: BackupsCategory,
		Name: "backup-extension", Aliases: []string{"B"},
		Usage: "specify the backup file suffix to use (implies -b)",
	}

	IgnoreCaseFlag = &cli.BoolFlag{Category: CaseSensitivityCategory,
		Name: "ignore-case", Aliases: []string{"i"},
		Usage: "perform a case-insensitive search (plain or regex)",
	}
	PreserveCaseFlag = &cli.BoolFlag{Category: CaseSensitivityCategory,
		Name: "preserve-case", Aliases: []string{"P"},
		Usage: "try to preserve replacement string cases",
	}

	NoLimitsFlag = &cli.BoolFlag{Category: GeneralCategory,
		Name: "no-limits", Aliases: []string{"U"},
		Usage: "ignore max file count and size limits",
	}
	NopFlag = &cli.BoolFlag{Category: GeneralCategory,
		Name: "nope", Aliases: []string{"nop", "n"},
		Usage: "report what would otherwise have been done",
	}

	ShowDiffFlag = &cli.BoolFlag{Category: UserInterfaceCategory,
		Name: "show-diff", Aliases: []string{"d"},
		Usage: "output unified diffs for all changes",
	}
	InteractiveFlag = &cli.BoolFlag{Category: UserInterfaceCategory,
		Name: "interactive", Aliases: []string{"e"},
		Usage: "selectively apply changes per-file",
	}
	PauseFlag = &cli.BoolFlag{Category: UserInterfaceCategory,
		Name: "pause", Aliases: []string{"E"},
		Usage: "pause on file search results screen (implies -e)",
	}

	RecurseFlag = &cli.BoolFlag{Category: TargetSelectionCategory,
		Name: "recurse", Aliases: []string{"R"},
		Usage: "travel directory paths",
	}
	AllFlag = &cli.BoolFlag{Category: TargetSelectionCategory,
		Name: "all", Aliases: []string{"a"},
		Usage: "include backups and files that start with a dot",
	}
	NullFlag = &cli.BoolFlag{Category: TargetSelectionCategory,
		Name: "null", Aliases: []string{"0"},
		Usage: "read null-terminated paths from os.Stdin",
	}
	FileFlag = &cli.StringSliceFlag{Category: TargetSelectionCategory,
		Name: "file", Aliases: []string{"f"},
		Usage: "read paths listed in files",
	}
	ExcludeFlag = &cli.StringSliceFlag{Category: TargetSelectionCategory,
		Name: "exclude", Aliases: []string{"X"},
		Usage: "exclude files matching glob pattern",
	}
	IncludeFlag = &cli.StringSliceFlag{Category: TargetSelectionCategory,
		Name: "include", Aliases: []string{"I"},
		Usage: "include on files matching glob pattern",
	}

	RegexFlag = &cli.BoolFlag{Category: RegularExpressionsCategory,
		Name: "regex", Aliases: []string{"r"},
		Usage: "search and replace arguments are regular expressions",
	}
	MultiLineFlag = &cli.BoolFlag{Category: RegularExpressionsCategory,
		Name: "multi-line", Aliases: []string{"m"},
		Usage: "set the multiline (?m) global flag (implies -r)",
	}
	DotMatchNlFlag = &cli.BoolFlag{Category: RegularExpressionsCategory,
		Name: "dot-match-nl", Aliases: []string{"s"},
		Usage: "set the dot-match-nl (?s) global flag (implies -r)",
	}

	QuietFlag = &cli.BoolFlag{Category: GeneralCategory,
		Name: "quiet", Aliases: []string{"q"},
		Usage: "silence notices",
	}
	VerboseFlag = &cli.BoolFlag{Category: GeneralCategory,
		Name: "verbose", Aliases: []string{"v"},
		Usage: "verbose notices",
	}

	UsageFlag = &cli.BoolFlag{Category: GeneralCategory,
		Name: "usage", Aliases: []string{"h"},
		Usage: "display command-line usage information",
	}
	HelpFlag = &cli.BoolFlag{Category: GeneralCategory,
		Name:  "help",
		Usage: "display complete command-line help text",
	}
	VersionFlag = &cli.BoolFlag{Category: GeneralCategory,
		Name: "version", Aliases: []string{"V"},
		Usage: "display the version",
	}
)

func init() {
	cdk.AppCliProfileFlag.Category = GoCursesCategory
	cdk.AppCliProfilePathFlag.Category = GoCursesCategory
	cdk.AppCliLogFileFlag.Category = GoCursesCategory
	cdk.AppCliLogLevelFlag.Category = GoCursesCategory
	cdk.AppCliLogLevelsFlag.Category = GoCursesCategory
	cdk.AppCliTtyFlag.Category = GoCursesCategory

	cli.FlagStringer = clcli.NewFlagStringer().
		PruneDefaultBools(true).
		DetailsOnNewLines(true).
		PruneRepeats(true).
		Make()
}
