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
	return
}

type cFindResults []cFindResult

func (r cFindResults) Tango() (markup string) {
	for idx, result := range r {
		if idx > 0 {
			markup += "\n"
		}
		markup += result.Tango()
	}
	return
}
