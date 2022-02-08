package replace

import (
	"io/ioutil"
	"regexp"
	"strings"
)

func String(search, replace, path string) (content string, err error) {
	var input []byte
	if input, err = ioutil.ReadFile(path); err == nil {
		content = strings.ReplaceAll(string(input), search, replace)
	}
	return
}

func Regexp(search, replace, path string) (content string, err error) {
	var input []byte
	if input, err = ioutil.ReadFile(path); err == nil {
		var rx *regexp.Regexp
		if rx, err = regexp.Compile(search); err == nil {
			content = rx.ReplaceAllString(string(input), replace)
		}
	}
	return
}
