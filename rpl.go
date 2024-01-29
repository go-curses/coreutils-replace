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

// Package replace provides a text replacement Worker concept for iteratively
// and selectively working with searching and replacing content.
package replace

import (
	"errors"
	"fmt"
	"os"

	"github.com/dustin/go-humanize"

	rpl "github.com/go-corelibs/replace"
)

var (
	DefaultBackupExtension = "~"
	DefaultBackupSeparator = "~"
)

var (
	TempErrPattern = fmt.Sprintf("rpl-%d.*.err", os.Getpid())
	TempOutPattern = fmt.Sprintf("rpl-%d.*.out", os.Getpid())
)

var (
	ErrNotFound      = errors.New("not found")
	ErrTooManyFiles  = fmt.Errorf("%w; try batches of %d or less", rpl.ErrTooManyFiles, rpl.MaxFileCount)
	gNoLimitsWarning = fmt.Sprintf("# WARNING: files larger than %s can consume all available memory\n", MaxFileSizeLabel)
)

var (
	MaxFileSizeLabel = humanize.Bytes(uint64(rpl.MaxFileSize))
)
