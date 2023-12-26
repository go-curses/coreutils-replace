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

func MakeRegexp(search string, o *Worker) (rx *regexp.Regexp, err error) {

	var rxFlags []string

	if o.MultiLine /*|| o.MultiLineDotMatchNl || o.MultiLineDotMatchNlInsensitive*/ {
		rxFlags = append(rxFlags, "m")
	}
	if o.DotMatchNl /*|| o.MultiLineDotMatchNl || o.MultiLineDotMatchNlInsensitive*/ {
		rxFlags = append(rxFlags, "s")
	}
	if o.IgnoreCase /*|| o.MultiLineDotMatchNlInsensitive*/ {
		rxFlags = append(rxFlags, "i")
	}

	var pattern string
	if len(rxFlags) > 0 {
		pattern = "(?" + strings.Join(rxFlags, "") + ")" + search
	} else {
		pattern = search
	}

	rx, err = regexp.Compile(pattern)
	return
}