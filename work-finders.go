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
type FindAllMatchingFn func(file string, matched bool, err error)

func IsIncluded(include, exclude []*glob.Glob, input string) (included bool) {
	if included = len(exclude) == 0 && len(include) == 0; included {
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

	// was not explicitly excluded
	if included = len(include) == 0; included {
		// and no include globs to constrain
		return
	}

	// must be explicitly included

	var err error
	for _, i := range include {
		if included, err = i.Match(input); err == nil && included {
			return
		} else {
			log.ErrorF("error matching: glob=%q; input=%q; err=%q", i.Pattern(), input, err)
		}
	}

	return
}

func FindAllIncluded(targets []string, includeHidden bool, include, exclude []*glob.Glob) (found []string) {
	for _, target := range targets {

		if path.IsFile(target) {
			// process file path
			if !includeHidden && path.IsHidden(target) {
				continue
			} else if path.IsPlainText(target) && IsIncluded(include, exclude, target) {
				found = append(found, target)
			}
		} else if path.IsDir(target) {
			// process dir path
			files, _ := path.ListFiles(target, includeHidden)
			for _, file := range files {
				if !includeHidden && path.IsHidden(file) {
					continue
				} else if path.IsPlainText(file) && IsIncluded(include, exclude, file) {
					found = append(found, file)
				}
			}
			dirs, _ := path.ListDirs(target, includeHidden)
			more := FindAllIncluded(dirs, includeHidden, include, exclude)
			found = append(found, more...)
		}

	}
	return
}

func FindAllMatcher(targets []string, includeHidden bool, include, exclude []*glob.Glob, fn FindAllMatchingFn, matcher FindAllMatcherFn) (files, matches []string) {
	if fn == nil {
		fn = func(file string, matched bool, err error) {}
	}
	for _, file := range FindAllIncluded(targets, includeHidden, include, exclude) {
		files = append(files, file)
		var err error
		var data []byte
		var matched bool
		if data, err = os.ReadFile(file); err == nil {
			if matched = matcher(data); matched {
				matches = append(matches, file)
			}
		}
		fn(file, matched, err)
	}
	return
}

func FindAllMatchingString(search string, targets []string, includeHidden bool, include, exclude []*glob.Glob, fn FindAllMatchingFn) (files, matches []string) {
	files, matches = FindAllMatcher(targets, includeHidden, include, exclude, fn, func(data []byte) (matched bool) {
		matched = strings.Contains(string(data), search)
		return
	})
	return
}

func FindAllMatchingStringInsensitive(search string, targets []string, includeHidden bool, include, exclude []*glob.Glob, fn FindAllMatchingFn) (files, matches []string) {
	files, matches = FindAllMatcher(targets, includeHidden, include, exclude, fn, func(data []byte) (matched bool) {
		matched = strings.Contains(strings.ToLower(string(data)), strings.ToLower(search))
		return
	})
	return
}

func FindAllMatchingRegexp(search *regexp.Regexp, targets []string, includeHidden bool, include, exclude []*glob.Glob, fn FindAllMatchingFn) (files, matches []string) {
	files, matches = FindAllMatcher(targets, includeHidden, include, exclude, fn, func(data []byte) (matched bool) {
		matched = search.Match(data)
		return
	})
	return
}
