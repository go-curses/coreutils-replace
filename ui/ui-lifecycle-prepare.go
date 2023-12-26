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
	"github.com/urfave/cli/v2"

	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/coreutils-replace"
)

// prepareStartup happens immediately upon cli action func
func (u *CUI) prepareStartup(data []interface{}, argv ...interface{}) cenums.EventFlag {
	var ok bool
	if ok = len(argv) == 2; ok {
		if u.Args, ok = argv[1].([]string); ok {
			return cenums.EVENT_PASS
		}
	}
	return cenums.EVENT_STOP
}

// prepare happens immediately upon cli action func, after prepareStartup and before everything else
func (u *CUI) prepare(data []interface{}, argv ...interface{}) cenums.EventFlag {
	var ok bool
	var ctx *cli.Context
	if len(argv) < 1 {
		log.ErrorF("internal error - prepare arguments not found")
		return cenums.EVENT_STOP
	} else if ctx, ok = argv[1].(*cli.Context); !ok {
		log.ErrorF("internal error - cli.Context not found (%T)", argv[0])
		return cenums.EVENT_STOP
	}

	if worker, eventFlag, err := replace.MakeWorker(ctx, u.notifier); err != nil {
		u.LastError = err
		return eventFlag
	} else if eventFlag == cenums.EVENT_STOP {
		return eventFlag
	} else {
		u.worker = worker
	}
	log.DebugF("prepared worker=%v", u.worker)

	if u.worker.Interactive {
		log.DebugF("starting interactive rpl")
		return cenums.EVENT_PASS
	}

	log.DebugF("processing non-interactive rpl")

	return cenums.EVENT_STOP
}