package main

import (
	cenums "github.com/go-curses/cdk/lib/enums"
)

var (
	gAbortWork = false
)

func shutdown(data []interface{}, argv ...interface{}) cenums.EventFlag {
	o, e := performWork()
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
	if gErr != nil {
		msg := gErr.Error()
		msgLen := len(msg)
		if msgLen > 0 && msg[msgLen-1] == '\n' {
			msg = msg[:msgLen-2]
		}
		notifier.Error("%v\n", msg)
	}
	if len(gWorkErrors) > 0 {
		for _, err := range gWorkErrors {
			notifier.Error("%v\n", err)
		}
	}
	return cenums.EVENT_PASS
}