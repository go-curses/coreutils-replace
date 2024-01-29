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
	"fmt"
	"html"
)

const (
	gFileSkipRune       = "-"
	gFileMatchRune      = "+"
	gFileErrorRune      = "x"
	gFileSkipCharacter  = `🗴`
	gFileMatchCharacter = `🗸`
	gFileErrorCharacter = `ℯ`
)

type cFindResult struct {
	target  string
	matched bool
	err     error
	markup  string
}

func (r cFindResult) Status(maxLen int) (line string) {
	var name string
	if size := len(r.target); size > maxLen {
		name = "..." + r.target[size-maxLen:]
	} else {
		name = r.target
	}
	if r.err != nil {
		line = gFileErrorCharacter + " " + name
	} else if r.matched {
		line = gFileMatchCharacter + " " + name
	} else {
		line = gFileSkipCharacter + " " + name
	}
	return
}

func (r cFindResult) Tango() (markup string) {
	if r.markup != "" {
		markup = r.markup
		return
	}
	if r.err != nil {
		markup = fmt.Sprintf(
			`[ <span foreground="red">%v</span> ]  %v <span foreground="#ffffff" background="red">/* %v */</span>`,
			gFileErrorRune, r.target, html.EscapeString(r.err.Error()),
		)
	} else if r.matched {
		markup = fmt.Sprintf(`[ <span foreground="gold">%v</span> ]  %v`, gFileMatchRune, r.target)
	} else {
		markup = fmt.Sprintf(`[ <span foreground="navy">%v</span> ]  %v`, gFileSkipRune, r.target)
	}
	r.markup = markup
	return
}

type cFindResults []cFindResult

func (r cFindResults) Status(maxLen int) (line string) {
	var errCount, matchedCount int
	for _, result := range r {
		if result.err != nil {
			errCount += 1
		} else if result.matched {
			matchedCount += 1
		}
	}
	line = fmt.Sprintf("%d errors; %d matched; %d total", errCount, matchedCount, len(r))
	if size := len(line); size >= maxLen {
		line = fmt.Sprintf("%d "+gFileMatchCharacter+"; %d "+gFileErrorCharacter+"; %d total", errCount, matchedCount, len(r))
	}
	return
}

func (r cFindResults) Tango() (markup string) {
	markup = r.TangoSlice(-1, -1)
	return
}

func (r cFindResults) TangoSlice(start, end int) (markup string) {
	var slice cFindResults
	if end > 0 {
		slice = r[start:end]
	} else if start > 0 {
		slice = r[start:]
	} else {
		slice = r
	}
	for idx, result := range slice {
		if idx > 0 {
			markup += "\n"
		}
		markup += result.Tango()
	}
	return
}

func (r cFindResults) TangoTail(size int) (markup string) {
	count := len(r)
	if size <= 0 || size > count {
		markup = r.Tango()
		return
	}
	markup = r.TangoSlice(count-size, -1)
	return
}
