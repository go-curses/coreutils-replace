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
	"os"
	"regexp"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/go-corelibs/chdirs"
	stdio "github.com/go-corelibs/mock-stdio"
	"github.com/go-corelibs/notify"
	rpl "github.com/go-corelibs/replace"
)

var (
	OriginalStdout = os.Stdout
	OriginalStderr = os.Stderr
)

func makeWorker() (outio, errio stdio.Stdio, w *Worker) {
	outio, errio = stdio.NewStdout(), stdio.NewStderr()
	_ = outio.Capture()
	_ = errio.Capture()
	w = &Worker{}
	return
}

func TestWorker(t *testing.T) {
	// only one test at a time due to capturing stdio
	m := &sync.Mutex{}

	Convey("Init", t, func() {

		Convey("Empty Settings", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			err := w.Init()
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(w.Notifier, ShouldNotEqual, nil)
			So(w.Notifier.Level(), ShouldEqual, notify.Info)
			So(w.Notifier.Stdout(), ShouldEqual, outio.Writer())
			So(w.Notifier.Stderr(), ShouldEqual, errio.Writer())
			So(w.FileWriterOut(), ShouldEqual, nil)
			So(w.FileWriterErr(), ShouldEqual, nil)
		})

		Convey("Empty Regex Set", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			w.Regex = true
			defer outio.Restore()
			defer errio.Restore()
			err := w.Init()
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=true;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(w.Pattern, ShouldNotEqual, nil)
		})

		Convey("Bad Regex Set", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			w.Regex = true
			w.Search = "[nope"
			defer outio.Restore()
			defer errio.Restore()
			err := w.Init()
			So(err, ShouldNotEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=true;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=[];`+
				`include=[];`+
				`}`)
			So(w.Pattern, ShouldEqual, (*regexp.Regexp)(nil))
		})

		Convey("Only Notifier Set", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			notifier := notify.New(notify.Debug).Make()
			w.Notifier = notifier
			err := w.Init()
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(w.Notifier, ShouldNotEqual, nil)
			So(w.Notifier.Level(), ShouldEqual, notify.Info)
			So(w.Notifier.Stdout(), ShouldEqual, outio.Writer())
			So(w.Notifier.Stderr(), ShouldEqual, errio.Writer())
		})

		Convey("Only Interactive Set", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.Interactive = true
			err := w.Init()
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=true;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(w.Notifier, ShouldNotEqual, nil)
			So(w.Notifier.Level(), ShouldEqual, notify.Info)
			So(w.Notifier.Stdout(), ShouldNotEqual, outio.Writer())
			So(w.Notifier.Stderr(), ShouldNotEqual, errio.Writer())
		})

		Convey("Only BackupExtension Set", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.BackupExtension = ".bak"
			err := w.Init()
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=.bak;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*.bak"];`+
				`include=[];`+
				`}`)
			So(w.Notifier, ShouldNotEqual, nil)
			So(w.Notifier.Level(), ShouldEqual, notify.Info)
			So(w.Notifier.Stdout(), ShouldEqual, outio.Writer())
			So(w.Notifier.Stderr(), ShouldEqual, errio.Writer())
		})

		Convey("Only Verbose Set", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.Verbose = true
			err := w.Init()
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=true;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(w.Notifier, ShouldNotEqual, nil)
			So(w.Notifier.Level(), ShouldEqual, notify.Debug)
			So(w.Notifier.Stdout(), ShouldEqual, outio.Writer())
			So(w.Notifier.Stderr(), ShouldEqual, errio.Writer())
		})

		Convey("Verbose and Quiet Set", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.Quiet = true
			w.Verbose = true
			err := w.Init()
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=true;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(w.Notifier, ShouldNotEqual, nil)
			So(w.Notifier.Level(), ShouldEqual, notify.Quiet)
			So(w.Notifier.Stdout(), ShouldEqual, outio.Writer())
			So(w.Notifier.Stderr(), ShouldEqual, errio.Writer())
		})

		Convey("Only Exclude Set", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.ExcludeArgs = []string{"*.txt"}
			err := w.Init()
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*.txt","*~"];`+
				`include=[];`+
				`}`)
			So(w.Notifier, ShouldNotEqual, nil)
			So(w.Notifier.Level(), ShouldEqual, notify.Info)
			So(w.Notifier.Stdout(), ShouldEqual, outio.Writer())
			So(w.Notifier.Stderr(), ShouldEqual, errio.Writer())
		})

		Convey("Exclude error", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.ExcludeArgs = []string{"[nope"}
			err := w.Init()
			So(err, ShouldNotEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=[];`+
				`include=[];`+
				`}`)
			So(w.Notifier, ShouldEqual, nil)
		})

		Convey("Only Include Set", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.IncludeArgs = []string{"*.txt"}
			err := w.Init()
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=["*.txt"];`+
				`}`)
			So(w.Notifier, ShouldNotEqual, nil)
			So(w.Notifier.Level(), ShouldEqual, notify.Info)
			So(w.Notifier.Stdout(), ShouldEqual, outio.Writer())
			So(w.Notifier.Stderr(), ShouldEqual, errio.Writer())
		})

		Convey("Include error", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.IncludeArgs = []string{"[nope"}
			err := w.Init()
			So(err, ShouldNotEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=[];`+
				`include=[];`+
				`}`)
			So(w.Notifier, ShouldEqual, nil)
		})

		Convey("Bad File Setting", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.AddFile = []string{"nope"}
			err := w.Init()
			outData := string(outio.Data())
			errData := string(errio.Data())
			So(err, ShouldEqual, nil)
			So(outData, ShouldEqual, ``)
			So(errData, ShouldEqual, ``)
		})

		Convey("No Limits Warning", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.NoLimits = true
			err := w.Init()
			outData := string(outio.Data())
			errData := string(errio.Data())
			So(err, ShouldEqual, nil)
			So(outData, ShouldEqual, ``)
			So(errData, ShouldEqual, gNoLimitsWarning+"\n")
		})

	})

	Convey("InitTargets", t, func() {

		Convey("Empty Settings", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			err := w.Init()
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(w.Notifier, ShouldNotEqual, nil)
			So(w.Notifier.Level(), ShouldEqual, notify.Info)
			So(w.Notifier.Stdout(), ShouldEqual, outio.Writer())
			So(w.Notifier.Stderr(), ShouldEqual, errio.Writer())
			err1 := w.InitTargets(nil)
			So(err1, ShouldEqual, nil)
		})

		Convey("One Path Settings", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.Paths = []string{"."}
			err := w.Init()
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(w.Notifier, ShouldNotEqual, nil)
			So(w.Notifier.Level(), ShouldEqual, notify.Info)
			So(w.Notifier.Stdout(), ShouldEqual, outio.Writer())
			So(w.Notifier.Stderr(), ShouldEqual, errio.Writer())
			err1 := w.InitTargets(nil)
			So(err1, ShouldEqual, nil)
		})

		Convey("Bad Path Settings", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.Paths = []string{"../..."}
			err := w.Init()
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(w.Notifier, ShouldNotEqual, nil)
			So(w.Notifier.Level(), ShouldEqual, notify.Info)
			So(w.Notifier.Stdout(), ShouldEqual, outio.Writer())
			So(w.Notifier.Stderr(), ShouldEqual, errio.Writer())
			err1 := w.InitTargets(nil)
			So(err1, ShouldEqual, nil)
		})

		Convey("Target Does Not Exist", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.Verbose = true
			w.Paths = []string{"not-a-thing"}
			err0 := w.Init()
			wString0 := w.String()
			err1 := w.InitTargets(nil)
			errData := string(errio.Data())
			So(err0, ShouldEqual, nil)
			So(wString0, ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=true;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(err1, ShouldEqual, nil)
			So(errData, ShouldEqual, "# error: not found: \"/quo/src/github.com/go-curses/coreutils-replace/not-a-thing\"\n")
		})

		Convey("Target Invalid Path", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			w.Verbose = true
			w.Paths = []string{"../not-a-thing"}
			err0 := w.Init()
			wString0 := w.String()
			var err1 error
			var errData string
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				chdirs.MockBadWD()
				defer chdirs.UnMockBadWD()
				err1 = w.InitTargets(nil)
				errData = string(errio.Data())
				wg.Done()
				return
			}()
			wg.Wait()
			So(err0, ShouldEqual, nil)
			So(wString0, ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=true;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(err1, ShouldEqual, nil)
			So(errData, ShouldEqual, "# error: \"../not-a-thing\" - getwd: no such file or directory\n")
		})
	})

	Convey("FindMatching", t, func() {

		Convey("Too Many Files (Regex)", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			oldMaxFiles := rpl.MaxFileCount
			defer func() { rpl.MaxFileCount = oldMaxFiles }()
			rpl.MaxFileCount = 1
			w.Regex = true
			w.Recurse = true
			w.Paths = append(w.Paths, "_testing")
			err0 := w.Init()
			wString0 := w.String()
			err1 := w.InitTargets(nil)
			err2 := w.FindMatching(nil)
			outData := string(outio.Data())
			errData := string(errio.Data())
			So(err0, ShouldEqual, nil)
			So(wString0, ShouldEqual, `Worker{`+
				`regex=true;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=true;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(err1, ShouldEqual, nil)
			So(err2, ShouldEqual, rpl.ErrTooManyFiles)
			So(outData, ShouldEqual, ``)
			So(errData, ShouldEqual, ``)
		})

		Convey("Too Many Files (Path)", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			oldMaxFiles := rpl.MaxFileCount
			defer func() { rpl.MaxFileCount = oldMaxFiles }()
			rpl.MaxFileCount = 1
			w.Recurse = true
			w.Paths = append(w.Paths, "_testing/test.txt", "_testing/hello.html")
			err0 := w.Init()
			wString0 := w.String()
			err1 := w.InitTargets(nil)
			outData := string(outio.Data())
			errData := string(errio.Data())
			So(err0, ShouldEqual, nil)
			So(wString0, ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=true;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(err1, ShouldEqual, ErrTooManyFiles)
			So(outData, ShouldEqual, ``)
			So(errData, ShouldEqual, ``)
		})

		Convey("Too Many Files (--file)", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			oldMaxFiles := rpl.MaxFileCount
			defer func() { rpl.MaxFileCount = oldMaxFiles }()
			rpl.MaxFileCount = 1
			w.AddFile = []string{"_testing/files.list"}
			err0 := w.Init()
			wString0 := w.String()
			err1 := w.InitTargets(nil)
			outData := string(outio.Data())
			errData := string(errio.Data())
			So(err0, ShouldEqual, nil)
			So(wString0, ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[_testing/files.list];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(err1, ShouldEqual, ErrTooManyFiles)
			So(outData, ShouldEqual, ``)
			So(errData, ShouldEqual, ``)
		})

		Convey("Error Scanning --file", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			fh, err := os.CreateTemp("", "coreutils-replace.*.tmp")
			tmpName := fh.Name()
			_, _ = fh.WriteString("Hello World\n")
			_ = fh.Close()
			_ = os.Chmod(tmpName, 0222)
			defer os.Remove(tmpName)
			w.AddFile = []string{tmpName}
			err0 := w.Init()
			wString0 := w.String()
			err1 := w.InitTargets(nil)
			outData := string(outio.Data())
			errData := string(errio.Data())
			So(err, ShouldEqual, nil)
			So(err0, ShouldEqual, nil)
			So(wString0, ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[`+tmpName+`];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(err1, ShouldEqual, nil)
			So(outData, ShouldEqual, ``)
			So(errData, ShouldEqual, "# error scanning --file \""+tmpName+"\": open "+tmpName+": permission denied")
		})

		Convey("Too Many Files (String)", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			oldMaxFiles := rpl.MaxFileCount
			defer func() { rpl.MaxFileCount = oldMaxFiles }()
			rpl.MaxFileCount = 1
			w.Recurse = true
			w.Paths = append(w.Paths, "_testing")
			err := w.Init()
			wString0 := w.String()
			err1 := w.InitTargets(nil)
			err2 := w.FindMatching(nil)
			outData := string(outio.Data())
			errData := string(errio.Data())
			So(err, ShouldEqual, nil)
			So(wString0, ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=true;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(err1, ShouldEqual, nil)
			So(err2, ShouldEqual, rpl.ErrTooManyFiles)
			So(outData, ShouldEqual, ``)
			So(errData, ShouldEqual, ``)
		})

		Convey("Too Many Files (StringInsensitive)", func() {
			m.Lock()
			defer m.Unlock()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			oldMaxFiles := rpl.MaxFileCount
			defer func() { rpl.MaxFileCount = oldMaxFiles }()
			rpl.MaxFileCount = 1
			w.IgnoreCase = true
			w.Recurse = true
			w.Paths = append(w.Paths, "_testing")
			err := w.Init()
			err1 := w.InitTargets(nil)
			err2 := w.FindMatching(nil)
			outData := string(outio.Data())
			errData := string(errio.Data())
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=true;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=true;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(err1, ShouldEqual, nil)
			So(err2, ShouldEqual, rpl.ErrTooManyFiles)
			So(outData, ShouldEqual, ``)
			So(errData, ShouldEqual, ``)
		})

		Convey("Too Many Files (os.Stdin)", func() {
			m.Lock()
			defer m.Unlock()
			inio := stdio.NewStdin([]byte(`_testing/hello.html
_testing/test.txt`))
			inio.Capture()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			oldMaxFiles := rpl.MaxFileCount
			defer func() { rpl.MaxFileCount = oldMaxFiles }()
			rpl.MaxFileCount = 1
			w.Stdin = true
			err := w.Init()
			err1 := w.InitTargets(nil)
			outData := string(outio.Data())
			errData := string(errio.Data())
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(err1, ShouldEqual, ErrTooManyFiles)
			So(outData, ShouldEqual, ``)
			So(errData, ShouldEqual, ``)
		})

		Convey("Too Many Files (os.Stdin --null)", func() {
			m.Lock()
			defer m.Unlock()
			inio := stdio.NewStdin([]byte(`_testing/hello.html` + string(rune(0)) + `_testing/test.txt`))
			inio.Capture()
			outio, errio, w := makeWorker()
			defer outio.Restore()
			defer errio.Restore()
			oldMaxFiles := rpl.MaxFileCount
			defer func() { rpl.MaxFileCount = oldMaxFiles }()
			rpl.MaxFileCount = 1
			w.Stdin = true
			w.Null = true
			err := w.Init()
			err1 := w.InitTargets(nil)
			outData := string(outio.Data())
			errData := string(errio.Data())
			So(err, ShouldEqual, nil)
			So(w.String(), ShouldEqual, `Worker{`+
				`regex=false;`+
				`multiLine=false;`+
				`dotMatchNl=false;`+
				`recurse=false;`+
				`nop=false;`+
				`all=false;`+
				`ignoreCase=false;`+
				`preserveCase=false;`+
				`binAsText=false;`+
				`noLimits=false;`+
				`backup=false;`+
				`backupExtension=;`+
				`showDiff=false;`+
				`interactive=false;`+
				`quiet=false;`+
				`verbose=false;`+
				`files=[];`+
				`exclude=["*~"];`+
				`include=[];`+
				`}`)
			So(err1, ShouldEqual, ErrTooManyFiles)
			So(outData, ShouldEqual, ``)
			So(errData, ShouldEqual, ``)
		})

	})

	Convey("StartIterating", t, func() {
		m.Lock()
		defer m.Unlock()
		outio, errio, w := makeWorker()
		defer outio.Restore()
		defer errio.Restore()
		w.Recurse = true
		w.Paths = []string{"_testing"}
		w.Search = "hello"
		w.Replace = "olleh"
		w.IgnoreCase = true
		err := w.Init()
		err1 := w.InitTargets(nil)
		err2 := w.FindMatching(nil)
		outData := string(outio.Data())
		errData := string(errio.Data())
		So(err, ShouldEqual, nil)
		So(w.String(), ShouldEqual, `Worker{`+
			`regex=false;`+
			`multiLine=false;`+
			`dotMatchNl=false;`+
			`recurse=true;`+
			`nop=false;`+
			`all=false;`+
			`ignoreCase=true;`+
			`preserveCase=false;`+
			`binAsText=false;`+
			`noLimits=false;`+
			`backup=false;`+
			`backupExtension=;`+
			`showDiff=false;`+
			`interactive=false;`+
			`quiet=false;`+
			`verbose=false;`+
			`files=[];`+
			`exclude=["*~"];`+
			`include=[];`+
			`}`)
		So(err1, ShouldEqual, nil)
		So(err2, ShouldEqual, nil)
		So(outData, ShouldEqual, ``)
		So(errData, ShouldEqual, ``)
		So(w.StartIterating(), ShouldNotEqual, nil)
	})

}
