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

package replace

import (
	"fmt"

	glob "github.com/ganbarodigital/go_glob"
)

type Globs []*glob.Glob

func (g Globs) String() (list string) {
	list += "["
	for idx, gg := range g {
		if idx > 0 {
			list += ","
		}
		list += fmt.Sprintf("%q", gg.Pattern())
	}
	list += "]"
	return
}

func ParseGlobs(patterns []string) (globs []*glob.Glob, err error) {
	for _, pattern := range patterns {
		g := glob.NewGlob(pattern)
		if _, err = g.Match("./test/file.txt"); err != nil {
			err = fmt.Errorf("--exclude %q error: %w", pattern, err)
			return
		}
		globs = append(globs, g)
	}
	return
}