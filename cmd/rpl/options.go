package main

import (
	"fmt"

	"github.com/go-curses/cdk/lib/sync"
)

type options struct {
	regex                          bool
	multiLine                      bool
	dotMatchNl                     bool
	multiLineDotMatchNl            bool
	multiLineDotMatchNlInsensitive bool
	recurse                        bool
	dryRun                         bool
	all                            bool
	ignoreCase                     bool
	backup                         bool
	backupExtension                string
	showDiff                       bool
	interactive                    bool
	quiet                          bool
	verbose                        bool

	sync.RWMutex
}

func (o *options) String() string {
	o.RLock()
	defer o.RUnlock()
	var s string
	s += "options{"
	s += fmt.Sprintf("regex=%v;", o.regex)
	s += fmt.Sprintf("multiLine=%v;", o.multiLine)
	s += fmt.Sprintf("dotMatchNl=%v;", o.dotMatchNl)
	s += fmt.Sprintf("multiLineDotMatchNl=%v;", o.multiLineDotMatchNl)
	s += fmt.Sprintf("multiLineDotMatchNlInsensitive=%v;", o.multiLineDotMatchNlInsensitive)
	s += fmt.Sprintf("recurse=%v;", o.recurse)
	s += fmt.Sprintf("dryRun=%v;", o.dryRun)
	s += fmt.Sprintf("all=%v;", o.all)
	s += fmt.Sprintf("ignoreCase=%v;", o.ignoreCase)
	s += fmt.Sprintf("backup=%v;", o.backup)
	s += fmt.Sprintf("backupExtension=%v;", o.backupExtension)
	s += fmt.Sprintf("showDiff=%v;", o.showDiff)
	s += fmt.Sprintf("interactive=%v;", o.interactive)
	s += fmt.Sprintf("quiet=%v;", o.quiet)
	s += fmt.Sprintf("verbose=%v;", o.verbose)
	s += "}"
	return s
}