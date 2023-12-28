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
	"regexp"
	"strings"
)

func RegexpReplace(search *regexp.Regexp, replace, contents string) (modified string) {
	modified = search.ReplaceAllString(contents, replace)
	return
}

func StringReplace(search, replace, contents string) (modified string) {
	modified = strings.ReplaceAll(contents[:], search, replace)
	return
}

func StringReplaceInsensitive(search, replace, contents string) string {
	// See: https://stackoverflow.com/questions/31348919/case-insensitive-string-replace-in-go/76512289

	if search == replace || search == "" {
		return contents
	}

	lowerContents := strings.ToLower(contents)
	lowerSearch := strings.ToLower(search)

	var count int
	if count = strings.Count(lowerContents, lowerSearch); count == 0 {
		return contents
	}

	var buffer strings.Builder
	buffer.Grow(len(contents) + (count * (len(replace) - len(search))))

	var start int
	for i := 0; i < count; i++ {
		j := start
		j += strings.Index(lowerContents[start:], lowerSearch)
		buffer.WriteString(contents[start:j])
		buffer.WriteString(replace)
		start = j + len(search)
	}
	buffer.WriteString(contents[start:])
	return buffer.String()
}