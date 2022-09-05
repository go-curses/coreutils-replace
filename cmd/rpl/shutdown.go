package main

import (
	"strings"

	cenums "github.com/go-curses/cdk/lib/enums"
)

var (
	gAbortWork = false
)

func shutdown(data []interface{}, argv ...interface{}) cenums.EventFlag {
	o, e := performWork()
	if len(gWorkErrors) > 0 {
		for _, err := range gWorkErrors {
			notifier.Error("# %v\n", strings.TrimSuffix(err.Error(), "\n"))
		}
	}
	if gErr != nil {
		notifier.Error("# error: %v\n", strings.TrimSuffix(gErr.Error(), "\n"))
	}
	if gOptions.showDiff {
		if gOptions.interactive {
			// diff goes to stderr, everything else to stdout
			notifier.Error(o)
			notifier.Info(e)
		} else {
			notifier.Info(o)
			notifier.Error(e)
		}
	} else {
		notifier.Info(o)
		notifier.Info(e)
	}
	return cenums.EVENT_PASS
}