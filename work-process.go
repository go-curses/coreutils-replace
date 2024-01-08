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
	"os"
	"regexp"

	"github.com/go-corelibs/diff"
)

func ProcessTargetString(search, replace, target string) (original, modified string, delta *diff.Diff, err error) {
	var data []byte
	if data, err = os.ReadFile(target); err == nil {
		original = string(data)
		modified = StringReplace(search, replace, original)
		delta = diff.New(target, original, modified)
	}
	return
}

func ProcessTargetStringInsensitive(search, replace, target string) (original, modified string, delta *diff.Diff, err error) {
	var data []byte
	if data, err = os.ReadFile(target); err == nil {
		original = string(data)
		modified = StringReplaceInsensitive(search, replace, original)
		delta = diff.New(target, original, modified)
	}
	return
}

func ProcessTargetRegex(search *regexp.Regexp, replace, target string) (original, modified string, delta *diff.Diff, err error) {
	var data []byte
	if data, err = os.ReadFile(target); err == nil {
		original = string(data)
		modified = RegexpReplace(search, replace, original)
		delta = diff.New(target, original, modified)
	}
	return
}