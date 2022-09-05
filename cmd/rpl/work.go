package main

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"

	"github.com/go-curses/corelibs/diff"
	"github.com/go-curses/corelibs/path"

	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/log"

	"github.com/gabriel-vasile/mimetype"
)

const (
	MaxInteractiveFileSize = 1024 * 1024 * 1 // 1 Mb
)

var (
	gTotalWork   = -1
	gSourceFiles []string
	gTargetIndex = -1
	gGroupsIndex = -1
	gDelta       []*diff.Diff
	gWorkErrors  []error
)

func findAllFiles(argv ...string) (files []string) {
	gOptions.RLock()
	defer gOptions.RUnlock()
	for _, src := range argv {
		var fileInfo os.FileInfo
		var err error
		if fileInfo, err = os.Stat(src); err != nil {
			gWorkErrors = append(gWorkErrors, err)
			continue
		}
		if fileInfo.Mode().IsDir() {
			if gOptions.recurse {
				files = append(files, findAllFiles(path.Ls(src, gOptions.all, false)...)...)
			} else if gOptions.verbose {
				gWorkErrors = append(gWorkErrors, fmt.Errorf("%v: ignoring directory", src))
			}
		} else if fileInfo.Mode().IsRegular() {
			if kind, err := mimetype.DetectFile(src); err == nil {
				if mime := kind.String(); len(mime) >= 4 && mime[:4] == "text" {
					if fileInfo.Size() >= MaxInteractiveFileSize {
						if gOptions.verbose {
							gWorkErrors = append(gWorkErrors, fmt.Errorf("%v: size of %d is greater than %d bytes", src, fileInfo.Size(), MaxInteractiveFileSize))
						}
					} else {
						files = append(files, src)
					}
				} else if gOptions.verbose {
					gWorkErrors = append(gWorkErrors, fmt.Errorf("%v: not plain text", src))
				}
			} else {
				gWorkErrors = append(gWorkErrors, fmt.Errorf("%v: %v", src, err))
			}
		} else if gOptions.verbose {
			gWorkErrors = append(gWorkErrors, fmt.Errorf("%v: not a file or directory", src))
		}
	}
	return
}

func regexReplace(content string) (modified string, err error) {
	var rxFlags []string

	if gOptions.multiLine || gOptions.multiLineDotMatchNl || gOptions.multiLineDotMatchNlInsensitive {
		rxFlags = append(rxFlags, "m")
	}
	if gOptions.dotMatchNl || gOptions.multiLineDotMatchNl || gOptions.multiLineDotMatchNlInsensitive {
		rxFlags = append(rxFlags, "s")
	}
	if gOptions.ignoreCase || gOptions.multiLineDotMatchNlInsensitive {
		rxFlags = append(rxFlags, "i")
	}

	s := gSearch
	if len(rxFlags) > 0 {
		s = "(?" + strings.Join(rxFlags, "") + ")" + s
	}
	rx := regexp.MustCompile(s)
	modified = rx.ReplaceAllString(content, gReplace)
	return
}

func stringReplace(content string) (modified string, err error) {
	if gOptions.ignoreCase {
		rx := regexp.MustCompile("(?msi)\\Q" + gSearch + "\\E")
		log.DebugF("strings.Regexp: %v, %v, %v", rx, gReplace, gSourceFiles[gTargetIndex])
		modified = rx.ReplaceAllString(content, gReplace)
	} else {
		log.DebugF("strings.ReplaceAll: %v, %v, %v", gSourceFiles[gTargetIndex], gSearch, gReplace)
		modified = strings.ReplaceAll(content, gSearch, gReplace)
	}
	return
}

func tangoDiff(unified string) (markup string) {
	for _, line := range strings.Split(unified, "\n") {
		lineLength := len(line)
		if lineLength > 0 {
			switch line[0] {
			case '+':
				markup += "<span foreground=\"#ffffff\" background=\"#007700\">"
				markup += html.EscapeString(line)
				markup += "</span>\n"
			case '-':
				markup += "<span foreground=\"#ffffff\" background=\"#770000\">"
				markup += html.EscapeString(line)
				markup += "</span>\n"
			case '@', ' ':
				fallthrough
			default:
				markup += "<span weight=\"dim\">"
				markup += html.EscapeString(line)
				markup += "</span>\n"
			}
		} else {
			markup += "\n"
		}
	}
	return
}

func displayEditView() {
	gSkipButton.Hide()
	gEditButton.Hide()
	gKeepButton.Hide()
	gApplyButton.Hide()

	_ = gMainLabel.SetMarkup(fmt.Sprintf(
		"<b>%d</b> of <b>%d</b> files, editing <b>%d</b> of <b>%d</b> changes:\n<b>%v</b>",
		gTargetIndex+1, len(gSourceFiles),
		gGroupsIndex+1, gDelta[gTargetIndex].EditGroupsLen(),
		gSourceFiles[gTargetIndex],
	))

	gSkipEditButton.Show()
	gKeepEditButton.Show()

	gWindow.Resize()
	gWindow.RequestDrawAndShow()
}

func displayFileView() {
	gSkipEditButton.Hide()
	gKeepEditButton.Hide()

	_ = gMainLabel.SetMarkup(fmt.Sprintf("<b>%d</b> of <b>%d</b> files, working on:\n<b>%v</b>", gTargetIndex+1, len(gSourceFiles), gSourceFiles[gTargetIndex]))

	numSourceFiles := len(gSourceFiles)
	if numSourceFiles > 1 {
		gSkipButton.Show()
	} else {
		gSkipButton.Hide()
	}

	numEditGroups := gDelta[gTargetIndex].EditGroupsLen()
	if numEditGroups > 0 {
		if numEditGroups > 1 {
			gEditButton.Show()
		} else {
			gEditButton.Hide()
		}
		gKeepButton.Show()
	} else {
		gEditButton.Hide()
		gApplyButton.Hide()
		gKeepButton.Hide()
	}

	gWindow.Resize()
	gWindow.RequestDrawAndShow()
}

func reviewFinalWork() {
	gSkipEditButton.Hide()
	gKeepEditButton.Hide()
	gEditButton.Hide()
	gKeepButton.Hide()
	gSkipButton.Hide()

	var unified string
	totalNumEdits := 0
	totalNumFiles := 0
	first := true
	for _, delta := range gDelta {
		numEdits := delta.KeepLen()
		totalNumEdits += numEdits
		if numEdits > 0 {
			if !first {
				unified += "\n"
			} else {
				first = false
			}
			unified += delta.UnifiedEdits()
			totalNumFiles += 1
		}
	}

	if totalNumFiles == 0 {
		_ = gMainLabel.SetMarkup(fmt.Sprintf("No changes selected for any of the %d files examined.", len(gSourceFiles)))
		w, h := gDisplay.Screen().Size()
		gMainLabel.SetSizeRequest(w-4, h-8)
		gDiffView.Hide()
		gQuitButton.GrabFocus()
	} else {
		_ = gMainLabel.SetMarkup(fmt.Sprintf("Reviewing %d pending changes across %d files:", totalNumEdits, totalNumFiles))
		gApplyButton.Show()
		gApplyButton.GrabFocus()
		if err := gDiffLabel.SetMarkup(tangoDiff(unified)); err != nil {
			gDiffLabel.LogErr(err)
		}
		gDiffLabel.SetSizeRequest(gDiffLabel.GetPlainTextInfo())
	}

	gWindow.Resize()
	gWindow.RequestDrawAndShow()
}

func keepCurrentEdit() {
	if gTargetIndex > -1 {
		numGroups := gDelta[gTargetIndex].EditGroupsLen()
		if gGroupsIndex >= 0 && gGroupsIndex < numGroups {
			gDelta[gTargetIndex].KeepGroup(gGroupsIndex)
		}
	}
}

func skipCurrentEdit() {
	if gTargetIndex > -1 {
		numGroups := gDelta[gTargetIndex].EditGroupsLen()
		if gGroupsIndex >= 0 && gGroupsIndex < numGroups {
			gDelta[gTargetIndex].SkipGroup(gGroupsIndex)
		}
	}
}

func processNextEdit() {
	if gTargetIndex > -1 {
		gGroupsIndex += 1

		numGroups := gDelta[gTargetIndex].EditGroupsLen()
		if gGroupsIndex < numGroups {
			unified := gDelta[gTargetIndex].EditGroup(gGroupsIndex)
			if err := gDiffLabel.SetMarkup(tangoDiff(unified)); err != nil {
				gDiffLabel.LogErr(err)
			}
			gDiffLabel.SetSizeRequest(gDiffLabel.GetPlainTextInfo())
			displayEditView()
		} else {
			gGroupsIndex = -1
			gTargetIndex -= 1
			processNextWork()
		}
	}
}

func skipCurrentWork() {
	if gTargetIndex > -1 {
		gDelta[gTargetIndex].SkipAll()
	}
}

func prepareWorkAt(index int) {
	var err error
	var source, changed string
	if gDelta[index] == nil {
		if source, err = path.ReadFile(gSourceFiles[index]); err != nil {
			gWorkErrors = append(gWorkErrors, fmt.Errorf("%v - %v", gSourceFiles[index], err))
			gDelta[index] = diff.New(gSourceFiles[index], "", "")
			return
		}

		if gOptions.regex || gOptions.dotMatchNl || gOptions.multiLine || gOptions.multiLineDotMatchNl || gOptions.multiLineDotMatchNlInsensitive {
			// regexp
			if changed, err = regexReplace(source); err != nil {
				gWorkErrors = append(gWorkErrors, fmt.Errorf("regexp - %v", err))
				changed = source
			}
		} else {
			// strings.Replace
			if changed, err = stringReplace(source); err != nil {
				gWorkErrors = append(gWorkErrors, fmt.Errorf("string - %v", err))
				changed = source
			}
		}

		gDelta[index] = diff.New(gSourceFiles[index], source, changed)
		gDelta[index].KeepAll()
	}
}

func processNextWork() {
	gTargetIndex += 1
	if gTargetIndex >= len(gSourceFiles) {
		reviewFinalWork()
		return
	}

	prepareWorkAt(gTargetIndex)

	unified := gDelta[gTargetIndex].UnifiedEdits()
	if unified == "" {
		if gOptions.verbose {
			unified = fmt.Sprintf("(no changes necessary, \"%v\" not found)", gSearch)
		} else {
			processNextWork()
			return
		}
	}
	if err := gDiffLabel.SetMarkup(tangoDiff(unified)); err != nil {
		gDiffLabel.LogErr(err)
	}
	gDiffLabel.SetSizeRequest(gDiffLabel.GetPlainTextInfo())

	displayFileView()
}

func startWorkProcess() cenums.EventFlag {
	gSourceFiles = findAllFiles(gTargets...)
	gTotalWork = len(gSourceFiles)
	if gTotalWork > 0 {
		gDelta = make([]*diff.Diff, gTotalWork)
		processNextWork()
	} else {
		gMainLabel.SetText("no files found")
	}
	return cenums.EVENT_STOP
}

func processCliWork() (err error) {
	gSourceFiles = findAllFiles(gTargets...)
	gTotalWork = len(gSourceFiles)
	if gTotalWork == 0 {
		if gOptions.verbose {
			err = fmt.Errorf("no files found")
		}
		return
	}
	gDelta = make([]*diff.Diff, gTotalWork)
	for idx, _ := range gSourceFiles {
		gTargetIndex = idx
		prepareWorkAt(idx)
	}
	return
}

func performWork() (o, e string) {
	if gAbortWork {
		if gOptions.verbose {
			e = "# rpl exited without making any changes\n"
		}
		return
	}

	if gOptions.dryRun {

		if gOptions.showDiff {
			var unified string
			for _, delta := range gDelta {
				if delta != nil && delta.KeepLen() > 0 {
					unified += delta.UnifiedEdits()
				}
			}
			if unified != "" {
				o += unified
			} else {
				e += fmt.Sprintf("# no changes made to any of the %d files examined (dry-run)\n", len(gSourceFiles))
			}
		} else {
			for idx, delta := range gDelta {
				if delta != nil {
					if delta.KeepLen() > 0 {
						o += fmt.Sprintf("# changed (dry-run): %v\n", gSourceFiles[idx])
					} else {
						o += fmt.Sprintf("# skipped (dry-run): %v\n", gSourceFiles[idx])
					}
				}
			}
		}

		e += "# rpl exited without making any changes (dry-run)\n"
		return
	}

	found := false
	totalNumEdits := 0
	totalFilesEdited := 0
	diffOut := ""

	for idx, delta := range gDelta {
		if delta != nil {
			if modified, err := delta.ModifiedEdits(); err != nil {
				e += fmt.Sprintf("error applying edits %v: %v\n", gSourceFiles[idx], err)
			} else {
				if gOptions.backup {
					if gOptions.backupExtension == "" {
						gOptions.backupExtension = ".bak"
					}
					if err := path.BackupAndOverwrite(gSourceFiles[idx], gSourceFiles[idx]+gOptions.backupExtension, modified); err != nil {
						e += fmt.Sprintf("error writing %v: %v\n", gSourceFiles[idx], err)
						continue
					}
				} else if err := path.Overwrite(gSourceFiles[idx], modified); err != nil {
					e += fmt.Sprintf("error writing %v: %v\n", gSourceFiles[idx], err)
					continue
				}

				numKept := delta.KeepLen()
				if numKept > 0 {
					found = true
					totalNumEdits += numKept
					totalFilesEdited += 1
					if gOptions.showDiff {
						unified := delta.UnifiedEdits()
						if unified != "" {
							diffOut += fmt.Sprintf(unified)
							if gOptions.verbose {
								e += fmt.Sprintf("# %v: %d changes made\n", gSourceFiles[idx], numKept)
							}
						} else if gOptions.verbose {
							e += fmt.Sprintf("# %v: no changes made\n", gSourceFiles[idx])
						}
					} else if !gOptions.quiet {
						o += fmt.Sprintf("# %v: %d changes made\n", gSourceFiles[idx], numKept)
					}
				} else if gOptions.verbose {
					o += fmt.Sprintf("# %v: no changes made\n", gSourceFiles[idx])
				}
			}
		}
	}

	if gOptions.showDiff && diffOut != "" {
		if gOptions.verbose {
			o += "## start of diff output\n"
		}
		o += diffOut
		diffOut = ""
		if gOptions.verbose {
			o += "## end of diff output\n"
		}
	}

	if gOptions.verbose {
		if found {
			e += fmt.Sprintf("# rpl exited after making %d changes across %d files\n", totalNumEdits, totalFilesEdited)
		} else if totalFilesEdited > 0 {
			e += fmt.Sprintf("# no changes made to any of the %d files examined\n", totalFilesEdited)
		} else {
			e += fmt.Sprintf("# nothing to do\n")
		}
	}

	return
}