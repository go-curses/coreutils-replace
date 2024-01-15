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
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/go-corelibs/globs"
)

func TestFinders(t *testing.T) {
	Convey("IsIncluded", t, func() {
		So(IsIncluded(nil, nil, "/test"), ShouldEqual, true)
		include, err := globs.Parse("*.txt")
		So(err, ShouldEqual, nil)
		So(IsIncluded(include, nil, "_testing/test.txt"), ShouldEqual, true)
	})

	Convey("FindAllIncluded", t, func() {
		So(len(FindAllIncluded([]string{"_testing/.hello-world.hidden", "_testing/test.txt"}, false, false, false, nil, nil)), ShouldEqual, 1)
	})
}
