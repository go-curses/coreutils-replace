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
	"strings"

	glob "github.com/ganbarodigital/go_glob"

	"github.com/go-corelibs/path"
	"github.com/go-curses/cdk/log"
)

type FindAllMatcherFn func(data []byte) (matched bool)
type FindAllMatchingFn func(file string, matched bool)

func IsIncluded(exclude []*glob.Glob, input string) (included bool) {
	if included = len(exclude) == 0; included {
		return
	}
	for _, g := range exclude {
		if excluded, err := g.Match(input); err == nil {
			if included = !excluded; !included {
				return
			}
		} else {
			log.ErrorF("error matching: glob=%q; input=%q; err=%q", g.Pattern(), input, err)
		}
	}
	return
}

func FindAllIncluded(exclude []*glob.Glob, targets []string) (found []string) {
	for _, target := range targets {

		if path.IsFile(target) {
			if path.IsPlainText(target) && IsIncluded(exclude, target) {
				found = append(found, target)
			}
		} else if path.IsDir(target) {
			more, _ := path.ListAllFiles(target)
			for _, file := range more {
				if path.IsPlainText(file) && IsIncluded(exclude, file) {
					found = append(found, file)
				}
			}
		}

	}
	return
}

func FindAllMatcher(targets []string, exclude []*glob.Glob, fn FindAllMatchingFn, matcher FindAllMatcherFn) (files, matches []string) {
	if fn == nil {
		fn = func(file string, matched bool) {}
	}
	for _, file := range FindAllIncluded(exclude, targets) {
		var matched bool
		files = append(files, file)
		if data, err := os.ReadFile(file); err == nil {
			if matched = matcher(data); matched {
				matches = append(matches, file)
			}
		}
		fn(file, matched)
	}
	return
}

func FindAllMatchingString(search string, targets []string, exclude []*glob.Glob, fn FindAllMatchingFn) (files, matches []string) {
	files, matches = FindAllMatcher(targets, exclude, fn, func(data []byte) (matched bool) {
		matched = strings.Contains(string(data), search)
		return
	})
	return
}

func FindAllMatchingStringInsensitive(search string, targets []string, exclude []*glob.Glob, fn FindAllMatchingFn) (files, matches []string) {
	files, matches = FindAllMatcher(targets, exclude, fn, func(data []byte) (matched bool) {
		matched = strings.Contains(strings.ToLower(string(data)), strings.ToLower(search))
		return
	})
	return
}

func FindAllMatchingRegexp(search *regexp.Regexp, targets []string, exclude []*glob.Glob, fn FindAllMatchingFn) (files, matches []string) {
	files, matches = FindAllMatcher(targets, exclude, fn, func(data []byte) (matched bool) {
		matched = search.Match(data)
		return
	})
	return
}