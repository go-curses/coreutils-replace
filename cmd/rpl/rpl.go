package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/go-curses/coreutils/path"
	"github.com/go-curses/coreutils/replace"
)

func rplFileContents(c *cli.Context, s, r, p, content string) (err error) {
	tmpname := c.String("tmpdir") + string(os.PathSeparator) + filepath.Base(p) + ".tmp"
	if err := path.Overwrite(tmpname, content); err != nil {
		if path.Exists(tmpname) {
			_ = os.Remove(tmpname)
		}
		return newError(err.Error())
	}
	var unified string
	if unified, err = path.Diff(p, tmpname); err != nil {
		_ = os.Remove(tmpname)
		return
	}
	if unified != "" {
		if c.Bool("dry-run") {
			_ = os.Remove(tmpname)
			if c.Bool("show-diff") {
				log.Info("# dry-run changes not applied:\n%v\n", unified)
			} else {
				log.Info("# dry-run did not modify: %v\n", p)
			}
		} else {
			if c.Bool("backup") {
				if err := path.CopyFile(p, p+"."+c.String("backup-extension")); err != nil {
					_ = os.Remove(tmpname)
					return newError(err.Error())
				}
			}
			if err := path.MoveFile(tmpname, p); err != nil {
				_ = os.Remove(tmpname)
				return newError(err.Error())
			}
			_ = os.Remove(tmpname)
			if c.Bool("show-diff") {
				log.Info("# changes applied:\n%v\n", unified)
			} else {
				log.Info("# modified: %v\n", p)
			}
		}
	}
	return nil
}

func rplFileString(c *cli.Context, s, r, p string) (err error) {
	var content string
	if c.Bool("ignore-case") {
		pattern := fmt.Sprintf(`(?i)\Q%v\E`, s)
		return rplFileRegexp(c, pattern, r, p)
	}
	if content, err = replace.String(s, r, p); err != nil {
		return
	}
	return rplFileContents(c, s, r, p, content)
}

func rplFileRegexp(c *cli.Context, s, r, p string) (err error) {
	var content string
	rxFlags := []string{}
	if c.Bool("multi-line") ||
		c.Bool("multi-line-dot-match-nl") ||
		c.Bool("multi-line-dot-match-nl-insensitive") {
		rxFlags = append(rxFlags, "m")
	}
	if c.Bool("dot-match-nl") ||
		c.Bool("multi-line-dot-match-nl") ||
		c.Bool("multi-line-dot-match-nl-insensitive") {
		rxFlags = append(rxFlags, "s")
	}
	if c.Bool("ignore-case") ||
		c.Bool("multi-line-dot-match-nl-insensitive") {
		rxFlags = append(rxFlags, "i")
	}
	if len(rxFlags) > 0 {
		s = "(?" + strings.Join(rxFlags, "") + ")" + s
	}
	if content, err = replace.Regexp(s, r, p); err != nil {
		return
	}
	return rplFileContents(c, s, r, p, content)
}

func rplPathString(c *cli.Context, s, r, p string) (errs []error) {
	paths := path.Ls(p, c.Bool("all"), c.Bool("recurse"))
	for _, pp := range paths {
		if err := rplFileString(c, s, r, pp); err != nil {
			errs = append(errs, err)
		}
	}
	return
}

func rplPathRegexp(c *cli.Context, s, r, p string) (errs []error) {
	paths := path.Ls(p, c.Bool("all"), c.Bool("recurse"))
	for _, pp := range paths {
		if err := rplFileRegexp(c, s, r, pp); err != nil {
			errs = append(errs, err)
		}
	}
	return nil
}
