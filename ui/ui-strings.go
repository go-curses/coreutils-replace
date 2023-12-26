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
	"html"
	"strings"
)

func tangoDiff(unified string) (markup string) {
	for _, line := range strings.Split(unified, "\n") {
		lineLength := len(line)
		if lineLength > 0 {
			switch line[0] {
			case '+':
				markup += "<span foreground=\"#ffffff\" background=\"#007700\">"
				markup += html.EscapeString(line)
				markup += "</span>\n"
			case '-':
				markup += "<span foreground=\"#ffffff\" background=\"#770000\">"
				markup += html.EscapeString(line)
				markup += "</span>\n"
			case '@', ' ':
				fallthrough
			default:
				markup += "<span weight=\"dim\">"
				markup += html.EscapeString(line)
				markup += "</span>\n"
			}
		} else {
			markup += "\n"
		}
	}
	return
}