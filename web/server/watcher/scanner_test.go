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
				Convey("The scanner should report a change in state", nil)
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
				Convey("The scanner should NOT report a change in state", nil)
			})

			Convey("When a go file is modified outside any watched folders", func() {
				Convey("The scanner should NOT report a change in state", nil)
			})

			Convey("When a go file is renamed outside any watched folders", func() {
				Convey("The scanner should NOT report a change in state", nil)
			})

			Convey("When a go file is deleted outside any watched folders", func() {
				Convey("The scanner should NOT report a change in state", nil)
			})

			Convey("When a miscellaneous file is created", func() {
				Convey("The scanner should NOT report a change in state", nil)
			})

			Convey("When a miscellaneous file is modified", func() {
				Convey("The scanner should NOT report a change in state", nil)
			})

			Convey("When a miscellaneous file is renamed", func() {
				Convey("The scanner should NOT report a change in state", nil)
			})

			Convey("When a miscellaneous file is deleted", func() {
				Convey("The scanner should NOT report a change in state", nil)
			})

			Convey("When a new folder is created inside a watched folder", func() {
				Convey("The scanner should report the change", nil)
				Convey("The scanner should notify the watcher of the creation", nil)
			})

			Convey("When a watched folder is deleted", func() {
				Convey("The scanner should report the change", nil)
				Convey("The scanner should notify the watcher of the deletion", nil)
			})

			Convey("When a folder is created outside any watched folders", func() {
				Convey("The scanner should NOT report the change", nil)
				Convey("The scanner should NOT notify the watcher of the change", nil)
			})

			Convey("When a folder that is not being watched is deleted", func() {
				Convey("The scanner should NOT report the change", nil)
				Convey("The scanner should NOT notify the watcher of the change", nil)
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
	return self.scanner.Scan("/root")
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
	self.watcher = NewWatcher(self.fs)
	self.watcher.Adjust("/root")
	self.scanner = NewScanner(self.fs, self.watcher)
	return self
}
