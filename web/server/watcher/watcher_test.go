package watcher

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/contract"
	"github.com/smartystreets/goconvey/web/server/system"
)

func TestWatcher(t *testing.T) {

	Convey("Subject: Watcher", t, func() {
		fixture := newWatcherFixture()

		Convey("When initialized there should be ZERO watched folders", func() {
			So(len(fixture.watched()), ShouldEqual, 0)
			So(fixture.watcher.Root(), ShouldBeBlank)
		})

		Convey("When pointing to a root folder", func() {
			actualWatches, expectedWatches := fixture.pointToExistingRoot(goProject)

			Convey("That folder should be included as the first watched folder", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})

			Convey("That folder should be the new root", func() {
				So(fixture.watcher.Root(), ShouldEqual, goProject)
			})
		})

		Convey("When pointing to a root folder that does not exist", func() {
			actualError, expectedError := fixture.pointToImaginaryRoot(normalize("/not/there"))

			Convey("An appropriate error should be returned", func() {
				So(actualError, ShouldResemble, expectedError)
			})

			Convey("The root should not be updated", func() {
				So(fixture.watcher.Root(), ShouldBeBlank)
			})
		})

		Convey("When pointing to a root folder with nested folders", func() {
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

			Convey("The folder should be marked as inactive in the watched folders listing", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("When multiple watched folders are ignored", func() {
			actualWatches, expectedWatches := fixture.ignoreWatchedFolders()
			Convey("The folders should be marked as inactive in the watched folders listing", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("When a folder that is not being watched is ignored", func() {
			actualWatches, expectedWatches := fixture.ignoreIrrelevantFolder()

			Convey("The request should be ignored", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("When a folder that does not exist is ignored", func() {
			actualWatches, expectedWatches := fixture.ignoreImaginaryFolder()

			Convey("There should be no change to the watched folders", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("When an ignored folder is reinstated", func() {
			actualWatches, expectedWatches := fixture.reinstateIgnoredFolder()

			Convey("The folder should be included once more in the watched folders", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("When multiple ignored folders are reinstated", func() {
			actualWatches, expectedWatches := fixture.reinstateIgnoredFolders()

			Convey("The folders should be included once more in the watched folders", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("When an ignored folder is deleted and then reinstated", func() {
			actualWatches, expectedWatches := fixture.reinstateDeletedFolder()

			Convey("The reinstatement request should be ignored", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("When a folder that is not being watched is reinstated", func() {
			actualWatches, expectedWatches := fixture.reinstateIrrelevantFolder()

			Convey("The request should be ignored", func() {
				So(actualWatches, ShouldResemble, expectedWatches)
			})
		})

		Convey("Regardless of the status of the watched folders", func() {
			folders := fixture.setupSeveralFoldersWithWatcher()

			Convey("The IsWatched query method should be correct", func() {
				So(fixture.watcher.IsWatched(folders["active"]), ShouldBeTrue)
				So(fixture.watcher.IsWatched(folders["reinstated"]), ShouldBeTrue)

				So(fixture.watcher.IsWatched(folders["ignored"]), ShouldBeFalse)
				So(fixture.watcher.IsWatched(folders["deleted"]), ShouldBeFalse)
				So(fixture.watcher.IsWatched(folders["irrelevant"]), ShouldBeFalse)
			})

			Convey("The IsIgnored query method should be correct", func() {
				So(fixture.watcher.IsIgnored(folders["ignored"]), ShouldBeTrue)

				So(fixture.watcher.IsIgnored(folders["active"]), ShouldBeFalse)
				So(fixture.watcher.IsIgnored(folders["reinstated"]), ShouldBeFalse)
				So(fixture.watcher.IsIgnored(folders["deleted"]), ShouldBeFalse)
				So(fixture.watcher.IsIgnored(folders["irrelevant"]), ShouldBeFalse)
			})
		})
	})
}

type watcherFixture struct {
	watcher *Watcher
	files   *system.FakeFileSystem
	shell   *system.FakeShell
}

func (self *watcherFixture) watched() []*contract.Package {
	return self.watcher.WatchedFolders()
}

func (self *watcherFixture) pointToExistingRoot(folder string) (actual, expected interface{}) {
	self.files.Create(folder, 1, time.Now())

	self.watcher.Adjust(folder)

	actual = self.watched()
	expected = []*contract.Package{
		&contract.Package{
			Active: true,
			Path:   goProject,
			Name:   goPackagePrefix,
			Result: contract.NewPackageResult(goPackagePrefix),
		},
	}
	return
}

func (self *watcherFixture) pointToImaginaryRoot(folder string) (actual, expected interface{}) {
	actual = self.watcher.Adjust(folder)
	expected = errors.New("Directory does not exist: '" + normalize("/not/there") + "'")
	return
}

func (self *watcherFixture) pointToExistingRootWithNestedFolders() (actual, expected interface{}) {
	self.files.Create(goProject, 1, time.Now())
	self.files.Create(goProject+normalize("/sub"), 2, time.Now())
	self.files.Create(goProject+normalize("/sub2"), 3, time.Now())
	self.files.Create(goProject+normalize("/sub/subsub"), 4, time.Now())

	self.watcher.Adjust(goProject)

	actual = self.watched()
	expected = []*contract.Package{
		&contract.Package{
			Active: true,
			Path:   goProject,
			Name:   goPackagePrefix,
			Result: contract.NewPackageResult(goPackagePrefix),
		},
		&contract.Package{
			Active: true,
			Path:   goProject + normalize("/sub"),
			Name:   goPackagePrefix + normalize("/sub"),
			Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub")),
		},
		&contract.Package{
			Active: true,
			Path:   goProject + normalize("/sub/subsub"),
			Name:   goPackagePrefix + normalize("/sub/subsub"),
			Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub/subsub")),
		},
		&contract.Package{
			Active: true,
			Path:   goProject + normalize("/sub2"),
			Name:   goPackagePrefix + normalize("/sub2"),
			Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub2")),
		},
	}
	return
}

func (self *watcherFixture) pointToRootOfGoPath() {
	self.files.Create(normalize("/root/gopath"), 5, time.Now())

	self.watcher.Adjust(normalize("/root/gopath"))
}

func (self *watcherFixture) pointToNestedPartOfGoPath() {
	self.files.Create(normalize("/root/gopath"), 5, time.Now())
	self.files.Create(normalize("/root/gopath/src/github.com/smartystreets/project"), 6, time.Now())

	self.watcher.Adjust(normalize("/root/gopath/src/github.com/smartystreets/project"))
}

func (self *watcherFixture) pointTo(path string) {
	self.files.Create(path, 5, time.Now())
	self.watcher.Adjust(path)
}

func (self *watcherFixture) setAmbientGoPath(path string) {
	self.shell.Setenv("GOPATH", path)
	self.files.Create(path, int64(42+len(path)), time.Now())
	self.watcher = NewWatcher(self.files, self.shell)
}

func (self *watcherFixture) receiveNotificationOfNewFolder() (actual, expected interface{}) {
	self.watcher.Creation(goProject + normalize("/sub"))

	actual = self.watched()
	expected = []*contract.Package{&contract.Package{Active: true, Path: goProject + normalize("/sub"), Name: goPackagePrefix + normalize("/sub"), Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub"))}}
	return
}

func (self *watcherFixture) receiveNotificationOfDeletedFolder() (actual, expected interface{}) {
	self.watcher.Creation(goProject + normalize("/sub2"))
	self.watcher.Creation(goProject + normalize("/sub"))

	self.watcher.Deletion(goProject + normalize("/sub"))

	actual = self.watched()
	expected = []*contract.Package{&contract.Package{Active: true, Path: goProject + normalize("/sub2"), Name: goPackagePrefix + normalize("/sub2"), Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub2"))}}
	return
}

func (self *watcherFixture) ignoreWatchedFolder() (actual, expected interface{}) {
	self.watcher.Creation(goProject + normalize("/sub2"))

	self.watcher.Ignore(goPackagePrefix + normalize("/sub2"))

	actual = self.watched()
	expected = []*contract.Package{&contract.Package{Active: false, Path: goProject + normalize("/sub2"), Name: goPackagePrefix + normalize("/sub2"), Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub2"))}}
	return
}

func (self *watcherFixture) ignoreWatchedFolders() (actual, expected interface{}) {
	self.watcher.Creation(goProject + normalize("/sub2"))
	self.watcher.Creation(goProject + normalize("/sub3"))
	self.watcher.Creation(goProject + normalize("/sub4"))

	self.watcher.Ignore(goPackagePrefix + normalize("/sub2;"+goPackagePrefix+normalize("/sub4")))

	actual = self.watched()
	expected = []*contract.Package{
		&contract.Package{Active: false, Path: goProject + normalize("/sub2"), Name: goPackagePrefix + normalize("/sub2"), Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub2"))},
		&contract.Package{Active: true, Path: goProject + normalize("/sub3"), Name: goPackagePrefix + normalize("/sub3"), Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub3"))},
		&contract.Package{Active: false, Path: goProject + normalize("/sub4"), Name: goPackagePrefix + normalize("/sub4"), Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub4"))},
	}
	return
}

func (self *watcherFixture) ignoreIrrelevantFolder() (actual, expected interface{}) {
	self.files.Create(goProject, 1, time.Now())
	self.files.Create(normalize("/something"), 1, time.Now())
	self.watcher.Adjust(goProject)

	self.watcher.Ignore(normalize("/something"))

	actual = self.watched()
	expected = []*contract.Package{&contract.Package{Active: true, Path: goProject, Name: goPackagePrefix, Result: contract.NewPackageResult(goPackagePrefix)}}
	return
}

func (self *watcherFixture) ignoreImaginaryFolder() (actual, expected interface{}) {
	self.files.Create(goProject, 1, time.Now())
	self.watcher.Adjust(goProject)

	self.watcher.Ignore(normalize("/not/there"))

	actual = self.watched()
	expected = []*contract.Package{&contract.Package{Active: true, Path: goProject, Name: goPackagePrefix, Result: contract.NewPackageResult(goPackagePrefix)}}
	return
}

func (self *watcherFixture) reinstateIgnoredFolder() (actual, expected interface{}) {
	self.files.Create(goProject, 1, time.Now())
	self.files.Create(goProject+normalize("/sub"), 2, time.Now())
	self.watcher.Adjust(goProject)
	self.watcher.Ignore(goPackagePrefix + normalize("/sub"))

	self.watcher.Reinstate(goProject + normalize("/sub"))

	actual = self.watched()
	expected = []*contract.Package{
		&contract.Package{Active: true, Path: goProject, Name: goPackagePrefix, Result: contract.NewPackageResult(goPackagePrefix)},
		&contract.Package{Active: true, Path: goProject + normalize("/sub"), Name: goPackagePrefix + normalize("/sub"), Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub"))},
	}
	return
}

func (self *watcherFixture) reinstateIgnoredFolders() (actual, expected interface{}) {
	self.files.Create(goProject, 1, time.Now())
	self.files.Create(goProject+normalize("/sub"), 2, time.Now())
	self.files.Create(goProject+normalize("/sub2"), 3, time.Now())
	self.files.Create(goProject+normalize("/sub3"), 4, time.Now())
	self.watcher.Adjust(goProject)
	self.watcher.Ignore(goPackagePrefix + normalize("/sub;") + goPackagePrefix + normalize("/sub2;") + goPackagePrefix + normalize("/sub3"))

	self.watcher.Reinstate(goProject + normalize("/sub;") + goPackagePrefix + normalize("/sub3"))

	actual = self.watched()
	expected = []*contract.Package{
		&contract.Package{Active: true, Path: goProject, Name: goPackagePrefix, Result: contract.NewPackageResult(goPackagePrefix)},
		&contract.Package{Active: true, Path: goProject + normalize("/sub"), Name: goPackagePrefix + normalize("/sub"), Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub"))},
		&contract.Package{Active: false, Path: goProject + normalize("/sub2"), Name: goPackagePrefix + normalize("/sub2"), Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub2"))},
		&contract.Package{Active: true, Path: goProject + normalize("/sub3"), Name: goPackagePrefix + normalize("/sub3"), Result: contract.NewPackageResult(goPackagePrefix + normalize("/sub3"))},
	}
	return
}

func (self *watcherFixture) reinstateDeletedFolder() (actual, expected interface{}) {
	self.files.Create(goProject, 1, time.Now())
	self.files.Create(goProject+normalize("/sub"), 2, time.Now())
	self.watcher.Adjust(goProject)
	self.watcher.Ignore(goPackagePrefix + normalize("/sub"))
	self.watcher.Deletion(goProject + normalize("/sub"))

	self.watcher.Reinstate(goPackagePrefix + normalize("/sub"))

	actual = self.watched()
	expected = []*contract.Package{&contract.Package{Active: true, Path: goProject, Name: goPackagePrefix, Result: contract.NewPackageResult(goPackagePrefix)}}
	return
}

func (self *watcherFixture) reinstateIrrelevantFolder() (actual, expected interface{}) {
	self.files.Create(goProject, 1, time.Now())
	self.files.Create(normalize("/irrelevant"), 2, time.Now())
	self.watcher.Adjust(goProject)

	self.watcher.Reinstate(normalize("/irrelevant"))

	actual = self.watched()
	expected = []*contract.Package{
		&contract.Package{
			Active: true,
			Path:   goProject,
			Name:   goPackagePrefix,
			Result: contract.NewPackageResult(goPackagePrefix),
		},
	}
	return
}

func (self *watcherFixture) setupSeveralFoldersWithWatcher() map[string]string {
	self.files.Create(goProject, 0, time.Now())
	self.files.Create(goProject+normalize("/active"), 1, time.Now())
	self.files.Create(goProject+normalize("/reinstated"), 2, time.Now())
	self.files.Create(goProject+normalize("/ignored"), 3, time.Now())
	self.files.Create(goProject+normalize("/deleted"), 4, time.Now())
	self.files.Create(normalize("/irrelevant"), 5, time.Now())

	self.watcher.Adjust(goProject)
	self.watcher.Ignore(goPackagePrefix + normalize("/ignored"))
	self.watcher.Ignore(goPackagePrefix + normalize("/reinstated"))
	self.watcher.Reinstate(goPackagePrefix + normalize("/reinstated"))
	self.watcher.Deletion(goProject + normalize("/deleted"))
	self.files.Delete(goProject + normalize("/deleted"))

	return map[string]string{
		"active":     goProject + normalize("/active"),
		"reinstated": goProject + normalize("/reinstated"),
		"ignored":    goProject + normalize("/ignored"),
		"deleted":    goProject + normalize("/deleted"),
		"irrelevant": normalize("/irrelevant"),
	}
}

func newWatcherFixture() *watcherFixture {
	self := new(watcherFixture)
	self.files = system.NewFakeFileSystem()
	self.shell = system.NewFakeShell()
	self.shell.Setenv("GOPATH", gopath)
	self.watcher = NewWatcher(self.files, self.shell)
	return self
}

var gopath = normalize("/root/gopath")
var goPackagePrefix = normalize("github.com/smartystreets/project")
var goProject = filepath.Join(gopath, normalize("/src/"), goPackagePrefix)

func init() {
	SetDefaultFailureMode(FailureContinues)
}
