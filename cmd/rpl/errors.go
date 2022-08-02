package main

import (
	"fmt"

	"github.com/go-curses/coreutils/errors"
)

func newError(messages ...string) error {
	return errors.NewPrefixed(
		fmt.Sprintf("%v error", APP_NAME),
		messages...,
	)
}

func newErrorF(format string, argv ...interface{}) error {
	return errors.NewPrefixedF(
		fmt.Sprintf("%v error", APP_NAME),
		format,
		argv...,
	)
}

func newUsageError(messages ...string) error {
	usage := fmt.Sprintf(
		"usage: %v [options] <search> <replace> <path> [...paths]",
		APP_NAME,
	)
	messages = append(messages, usage)
	return errors.NewPrefixed(
		fmt.Sprintf("%v error", APP_NAME),
		messages...,
	)
}