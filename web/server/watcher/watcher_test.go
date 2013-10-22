package watcher

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/contract"
	"github.com/smartystreets/goconvey/web/server/system"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	var (
		fs      *system.FakeFileSystem
		watcher *Watcher
	)

	Convey("Subject: Watcher", t, func() {
		fs = system.NewFakeFileSystem()
		watcher = NewWatcher(fs)

		Convey("When initialized there should be ZERO watched folders", func() {
			So(len(watcher.WatchedFolders()), ShouldEqual, 0)
		})

		Convey("When pointing to a folder", func() {
			fs.Create("/root", 1, time.Now())
			watcher.Adjust("/root")

			Convey("That folder should be included as the first watched folder", func() {
				So(watcher.WatchedFolders(), ShouldResemble, []*contract.Package{
					&contract.Package{
						Active: true,
						Path:   "/root",
						Name:   "root",
					},
				})
			})
		})

		Convey("When pointing to a folder that does not exist", func() {
			err := watcher.Adjust("/not/there")

			Convey("An appropriate error should be returned", func() {
				So(err, ShouldResemble, errors.New("Directory does not exist: '/not/there'"))
			})
		})

		Convey("When pointing to a folder with nested folders", func() {
			fs.Create("/root", 1, time.Now())
			fs.Create("/root/sub", 2, time.Now())
			fs.Create("/root/sub2", 3, time.Now())
			fs.Create("/root/sub/subsub", 4, time.Now())
			watcher.Adjust("/root")

			Convey("All nested folders should be added recursively to the watched folders", func() {
				So(watcher.WatchedFolders(), ShouldResemble, []*contract.Package{
					&contract.Package{Active: true, Path: "/root", Name: "root"},
					&contract.Package{Active: true, Path: "/root/sub", Name: "sub"},
					&contract.Package{Active: true, Path: "/root/sub2", Name: "sub2"},
					&contract.Package{Active: true, Path: "/root/sub/subsub", Name: "subsub"},
				})
			})
		})

		Convey("When creating a new folder", func() {
			watcher.Creation("/root/sub")

			Convey("The folder should be included in the watched folders", func() {
				So(watcher.WatchedFolders(), ShouldResemble, []*contract.Package{
					&contract.Package{Active: true, Path: "/root/sub", Name: "sub"},
				})
			})
		})

		Convey("When deleting an existing folder", func() {
			watcher.Creation("/root/sub2")
			watcher.Creation("/root/sub")
			watcher.Deletion("/root/sub")

			Convey("The folder should no longer be included in the watched folders", func() {
				So(watcher.WatchedFolders(), ShouldResemble, []*contract.Package{
					&contract.Package{Active: true, Path: "/root/sub2", Name: "sub2"},
				})
			})
		})

		Convey("When a watched folder is ignored", func() {
			watcher.Creation("/root/sub2")
			watcher.Ignore("/root/sub2")

			Convey("The folder should not be included in the watched folders", func() {
				So(len(watcher.WatchedFolders()), ShouldEqual, 0)
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

func init() {
	fmt.Sprintf("Keeps fmt in the import list...")
}
