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
	"io"
	"regexp"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIterator(t *testing.T) {
	m := &sync.Mutex{}
	Convey("StringInsensitive", t, func() {
		m.Lock()
		defer m.Unlock()
		outio, errio, w := makeWorker()
		defer outio.Restore()
		defer errio.Restore()
		w.Paths = []string{"_testing"}
		w.Search = "hello"
		w.Replace = "olleh"
		w.Recurse = true
		w.IgnoreCase = true
		err0 := w.Init()
		err1 := w.InitTargets(nil)
		err2 := w.FindMatching(nil)
		outData := string(outio.Data())
		errData := string(errio.Data())
		So(err0, ShouldEqual, nil)
		So(err1, ShouldEqual, nil)
		So(err2, ShouldEqual, nil)
		So(outData, ShouldEqual, ``)
		So(errData, ShouldEqual, ``)

		iter := w.StartIterating()
		So(iter, ShouldNotEqual, nil)
		So(iter.Valid(), ShouldEqual, true)
		So(iter.Pos(), ShouldEqual, 0)
		So(iter.Name(), ShouldEqual, w.Matched[0])

		original, modified, count, delta, err := iter.Replace()
		So(err, ShouldEqual, nil)
		So(original, ShouldNotEqual, modified)
		So(count, ShouldEqual, 2)
		So(delta.Len(), ShouldEqual, 4)

		w.IgnoreCase = false
		original, modified, count, delta, err = iter.Replace()
		So(err, ShouldEqual, nil)
		So(original, ShouldNotEqual, modified)
		So(count, ShouldEqual, 2)
		So(delta.Len(), ShouldEqual, 4)

		w.PreserveCase = true
		original, modified, count, delta, err = iter.Replace()
		So(err, ShouldEqual, nil)
		So(original, ShouldNotEqual, modified)
		So(count, ShouldEqual, 2)
		So(delta.Len(), ShouldEqual, 4)

		w.Regex = true
		w.Pattern = regexp.MustCompile(`hello`)
		original, modified, count, delta, err = iter.Replace()
		So(err, ShouldEqual, nil)
		So(original, ShouldNotEqual, modified)
		So(count, ShouldEqual, 2)
		So(delta.Len(), ShouldEqual, 4)

		w.PreserveCase = false
		original, modified, count, delta, err = iter.Replace()
		So(err, ShouldEqual, nil)
		So(original, ShouldNotEqual, modified)
		So(count, ShouldEqual, 2)
		So(delta.Len(), ShouldEqual, 4)

		w.Regex = false
		w.Pattern = nil
		iter.Next()
		So(iter.Valid(), ShouldEqual, true)
		iter.pos = len(w.Matched)
		iter.Next()
		So(iter.Valid(), ShouldEqual, false)
		So(iter.Pos(), ShouldEqual, len(w.Matched))

		original, modified, count, delta, err = iter.Replace()
		So(err, ShouldEqual, io.EOF)
	})
}
