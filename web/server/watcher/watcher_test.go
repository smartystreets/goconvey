package watcher

import (
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/contract"
	"github.com/smartystreets/goconvey/web/server/system"
)

func TestWatcher(t *testing.T) {
	var fixture *watcherFixture

	Convey("Subject: Watcher", t, func() {
		fixture = newWatcherFixture()

		Convey("When initialized there should be ZERO watched folders", func() {
			So(len(fixture.watched()), ShouldEqual, 0)
		})

		Convey("When pointing to a folder", func() {
			actualWatches, expectedWatches := fixture.pointToExistingRoot("/root")

			Convey("That folder should be included as the first watched folder", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("When pointing to a folder that does not exist", func() {
			actualError, expectedError := fixture.pointToImaginaryRoot("/not/there")

			Convey("An appropriate error should be returned", func() {
				So(actualError, ShouldResemble, expectedError)
			})
		})

		Convey("When pointing to a folder with nested folders", func() {
			actualWatches, expectedWatches := fixture.pointToExistingRootWithNestedFolders()

			Convey("All nested folders should be added recursively to the watched folders", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("When the watcher is notified of a newly created folder", func() {
			actualWatches, expectedWatches := fixture.receiveNotificationOfNewFolder()

			Convey("The folder should be included in the watched folders", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("When the watcher is notified of a recently deleted folder", func() {
			actualWatches, expectedWatches := fixture.receiveNotificationOfDeletedFolder()

			Convey("The folder should no longer be included in the watched folders", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("When a watched folder is ignored", func() {
			actualWatches, expectedWatches := fixture.ignoreWatchedFolder()

			Convey("The folder should not be included in the watched folders", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("When a folder that is not being watched is ignored", func() {
			Convey("The request should be ignored", nil)
		})

		Convey("When a folder that does not exist is ignored", func() {
			Convey("There should be no change to the watched folders", nil)
		})

		Convey("When an ignored folder is reinstated", func() {
			Convey("The folder should be included once more in the watched folders", nil)
		})

		Convey("When an ignored folder is deleted and then reinstated", func() {
			Convey("The reinstatement request should be ignored", nil)
		})

		Convey("When a folder that is not being watched is reinstated", func() {
			Convey("The request should be ignored", nil)
		})
	})
}

type watcherFixture struct {
	watcher *Watcher
	fs      *system.FakeFileSystem
}

func (self *watcherFixture) watched() []*contract.Package {
	return self.watcher.WatchedFolders()
}
func (self *watcherFixture) pointToExistingRoot(folder string) (actual, expected interface{}) {
	self.fs.Create(folder, 1, time.Now())

	self.watcher.Adjust(folder)

	actual = self.watched()
	expected = []*contract.Package{&contract.Package{Active: true, Path: "/root", Name: "root"}}
	return
}
func (self *watcherFixture) pointToImaginaryRoot(folder string) (actual, expected interface{}) {
	actual = self.watcher.Adjust(folder)
	expected = errors.New("Directory does not exist: '/not/there'")
	return
}
func (self *watcherFixture) pointToExistingRootWithNestedFolders() (actual, expected interface{}) {
	self.fs.Create("/root", 1, time.Now())
	self.fs.Create("/root/sub", 2, time.Now())
	self.fs.Create("/root/sub2", 3, time.Now())
	self.fs.Create("/root/sub/subsub", 4, time.Now())

	self.watcher.Adjust("/root")

	actual = self.watched()
	expected = []*contract.Package{
		&contract.Package{Active: true, Path: "/root", Name: "root"},
		&contract.Package{Active: true, Path: "/root/sub", Name: "sub"},
		&contract.Package{Active: true, Path: "/root/sub2", Name: "sub2"},
		&contract.Package{Active: true, Path: "/root/sub/subsub", Name: "subsub"},
	}
	return
}
func (self *watcherFixture) receiveNotificationOfNewFolder() (actual, expected interface{}) {
	self.watcher.Creation("/root/sub")

	actual = self.watched()
	expected = []*contract.Package{&contract.Package{Active: true, Path: "/root/sub", Name: "sub"}}
	return
}
func (self *watcherFixture) receiveNotificationOfDeletedFolder() (actual, expected interface{}) {
	self.watcher.Creation("/root/sub2")
	self.watcher.Creation("/root/sub")

	self.watcher.Deletion("/root/sub")

	actual = self.watched()
	expected = []*contract.Package{&contract.Package{Active: true, Path: "/root/sub2", Name: "sub2"}}
	return
}
func (self *watcherFixture) ignoreWatchedFolder() (actual, expected interface{}) {
	self.watcher.Creation("/root/sub2")

	self.watcher.Ignore("/root/sub2")

	actual = self.watched()
	expected = []*contract.Package{}
	return
}

func newWatcherFixture() *watcherFixture {
	self := &watcherFixture{}
	self.fs = system.NewFakeFileSystem()
	self.watcher = NewWatcher(self.fs)
	return self
}

func init() {
	fmt.Sprintf("Keeps fmt in the import list...")
}
