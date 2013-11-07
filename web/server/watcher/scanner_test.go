package watcher

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/system"
	"testing"
	"time"
)

func TestScanner(t *testing.T) {
	var fixture *scannerFixture
	var changed bool

	Convey("To begin with, the scanner is provided a contrived file system environment", t, func() {
		fixture = newScannerFixture()

		Convey("When we call Scan() for the first time", func() {
			changed = fixture.scan()

			Convey("The scanner should report a change in state", func() {
				So(changed, ShouldBeTrue)
			})
		})

		Convey("Then, on subsequent calls to Scan()", func() {
			changed = fixture.scan()

			Convey("When the file system has not changed in any way", func() {

				Convey("The scanner should NOT report any change in state", func() {
					So(fixture.scan(), ShouldBeFalse)
				})
			})

			Convey("When a new go file is created within a watched folder", func() {
				fixture.fs.Create("/root/new_stuff.go", 42, time.Now())

				Convey("The scanner should report a change in state", func() {
					So(fixture.scan(), ShouldBeTrue)
				})
			})

			Convey("When an existing go file within a watched folder has been modified", func() {
				fixture.fs.Modify("/root/sub/file.go")

				Convey("The scanner should report a change in state", func() {
					So(fixture.scan(), ShouldBeTrue)
				})
			})

			Convey("When an existing go file within a watched folder has been renamed", func() {
				fixture.fs.Rename("/root/sub/file.go", "/root/sub/asdf.go")

				Convey("The scanner should report a change in state", func() {
					So(fixture.scan(), ShouldBeTrue)
				})
			})

			Convey("When an existing go file within a watched folder has been deleted", func() {
				fixture.fs.Delete("/root/sub/file.go")

				Convey("The scanner should report a change in state", func() {
					So(fixture.scan(), ShouldBeTrue)
				})
			})

			Convey("When a go file is created outside any watched folders", func() {
				fixture.fs.Create("/outside/new_stuff.go", 42, time.Now())

				Convey("The scanner should NOT report a change in state", func() {
					So(fixture.scan(), ShouldBeFalse)
				})
			})

			Convey("When a go file is modified outside any watched folders", func() {
				fixture.fs.Create("/outside/new_stuff.go", 42, time.Now())
				fixture.scan() // reset

				Convey("The scanner should NOT report a change in state", func() {
					So(fixture.scan(), ShouldBeFalse)
				})
			})

			Convey("When a go file is renamed outside any watched folders", func() {
				fixture.fs.Create("/outside/new_stuff.go", 42, time.Now())
				fixture.scan() // reset
				fixture.fs.Rename("/outside/new_stuff.go", "/outside/newer_stoff.go")

				Convey("The scanner should NOT report a change in state", func() {
					So(fixture.scan(), ShouldBeFalse)
				})
			})

			Convey("When a go file is deleted outside any watched folders", func() {
				fixture.fs.Create("/outside/new_stuff.go", 42, time.Now())
				fixture.scan() // reset
				fixture.fs.Delete("/outside/new_stuff.go")

				Convey("The scanner should NOT report a change in state", func() {
					So(fixture.scan(), ShouldBeFalse)
				})
			})

			Convey("When a miscellaneous file is created", func() {
				fixture.fs.Create("/root/new_stuff.MISC", 42, time.Now())

				Convey("The scanner should NOT report a change in state", func() {
					So(fixture.scan(), ShouldBeFalse)
				})
			})

			Convey("When a miscellaneous file is modified", func() {
				fixture.fs.Create("/root/new_stuff.MISC", 42, time.Now())
				fixture.scan() // reset

				Convey("The scanner should NOT report a change in state", func() {
					So(fixture.scan(), ShouldBeFalse)
				})
			})

			Convey("When a miscellaneous file is renamed", func() {
				fixture.fs.Create("/root/new_stuff.MISC", 42, time.Now())
				fixture.scan() // reset
				fixture.fs.Rename("/root/new_stuff.MISC", "/root/newer_stoff.MISC")

				Convey("The scanner should NOT report a change in state", func() {
					So(fixture.scan(), ShouldBeFalse)
				})
			})

			Convey("When a miscellaneous file is deleted", func() {
				fixture.fs.Create("/root/new_stuff.MISC", 42, time.Now())
				fixture.scan() // reset
				fixture.fs.Delete("/root/new_stuff.MISC")

				Convey("The scanner should NOT report a change in state", func() {
					So(fixture.scan(), ShouldBeFalse)
				})
			})

			Convey("When a new folder is created inside a watched folder", func() {
				fixture.fs.Create("/root/new", 41, time.Now())
				changed := fixture.scan()

				Convey("The scanner should report the change", func() {
					So(changed, ShouldBeTrue)
				})

				Convey("The scanner should notify the watcher of the creation", func() {
					So(fixture.wasCreated("/root/new"), ShouldBeTrue)
				})
			})

			Convey("When an empty watched folder is deleted", func() {
				fixture.fs.Delete("/root/sub/empty")
				changed := fixture.scan()

				Convey("The scanner should report the change", func() {
					So(changed, ShouldBeTrue)
				})

				Convey("The scanner should notify the watcher of the deletion", func() {
					So(fixture.wasDeleted("/root/sub/empty"), ShouldBeTrue)
				})
			})

			Convey("When a folder is created outside any watched folders", func() {
				fixture.fs.Create("/outside/asdf", 41, time.Now())
				changed := fixture.scan()

				Convey("The scanner should NOT report the change", func() {
					So(changed, ShouldBeFalse)
				})

				Convey("The scanner should NOT notify the watcher of the change", func() {
					So(fixture.wasCreated("/outside/asdf"), ShouldBeFalse)
				})
			})

			Convey("When an ignored folder is deleted", func() {
				fixture.watcher.Ignore("/root/sub/empty")
				fixture.fs.Delete("/root/sub/empty")
				changed := fixture.scan()

				Convey("The scanner should report the change", func() {
					So(changed, ShouldBeTrue)
				})

				Convey("The scanner should notify the watcher of the change", func() {
					So(fixture.wasDeleted("/root/sub/empty"), ShouldBeTrue)
				})
			})

			// Once upon a time the scanner didn't keep track of the root internally, it was
			// given as a parameter to the Scan() method. This meant that when the scanner
			// was instructed to scan a new root location it appeared to the scanner that
			// many of the internally stored folders had been deleted becuase they were not
			// part of the new root directory structure and they were reported as deletions
			// to the watcher, which was incorrect behavior.
			// TODO: use a mock for the watcher
			SkipConvey("When the watcher has adjusted the root", func() {
				fixture.fs.Create("/somewhere", 3, time.Now())
				fixture.fs.Create("/somewhere/else", 3, time.Now())
				fixture.watcher.Adjust("/somewhere")

				// This puts a previously watched folder back in the watcher,
				// so we can see if a deletion is signaled inadvertantly as a result of Scan():
				fixture.watcher.Creation("/root/sub")

				Convey("And the scanner reports a change", func() {
					changed := fixture.scan()

					Convey("The scanner should report the change", func() {
						So(changed, ShouldBeTrue)
					})

					Convey("The scanner should NOT notify the watcher of incorrect folder deletions", func() {
						watched := fixture.watcher.WatchedFolders()
						last := watched[len(watched)-1]
						So(last.Path, ShouldEqual, "/root/sub")
					})

					Convey("The scanner should NOT notify the watcher of incorrect folder creations", func() {

					})
				})
			})

		})
	})
}

type scannerFixture struct {
	scanner *Scanner
	fs      *system.FakeFileSystem
	watcher *Watcher
}

func (self *scannerFixture) scan() bool {
	return self.scanner.Scan()
}
func (self *scannerFixture) wasDeleted(folder string) bool {
	return !self.wasCreated(folder)
}
func (self *scannerFixture) wasCreated(folder string) bool {
	for _, w := range self.watcher.WatchedFolders() {
		if w.Path == folder {
			return true
		}
	}
	return false
}

func newScannerFixture() *scannerFixture {
	self := &scannerFixture{}
	self.fs = system.NewFakeFileSystem()
	self.fs.Create("/root", 0, time.Now())
	self.fs.Create("/root/file.go", 1, time.Now())
	self.fs.Create("/root/sub", 0, time.Now())
	self.fs.Create("/root/sub/file.go", 2, time.Now())
	self.fs.Create("/root/sub/empty", 0, time.Now())
	self.watcher = NewWatcher(self.fs, system.NewFakeShell())
	self.watcher.Adjust("/root")
	self.scanner = NewScanner(self.fs, self.watcher)
	return self
}
