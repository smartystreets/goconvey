package system

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type FakeFileSystem struct {
	steps []*FakeFileInfo
}

func (self *FakeFileSystem) Create(path string, size int64, modified time.Time) {
	self.steps = append(self.steps, newFileInfo(path, size, modified))
}

func (self *FakeFileSystem) Walk(root string, step filepath.WalkFunc) {
	for _, info := range self.steps {
		step(info.path, info, nil)
	}
}
func (self *FakeFileSystem) Exists(directory string) bool {
	for _, info := range self.steps {
		if info.IsDir() && info.path == directory {
			return true
		}
	}
	return false
}

func NewFakeFileSystem() *FakeFileSystem {
	self := &FakeFileSystem{}
	self.steps = []*FakeFileInfo{}
	return self
}

type FakeFileInfo struct {
	path     string
	size     int64
	modified time.Time
}

func (self *FakeFileInfo) Name() string       { return filepath.Base(self.path) }
func (self *FakeFileInfo) Size() int64        { return self.size }
func (self *FakeFileInfo) Mode() os.FileMode  { return 0 }
func (self *FakeFileInfo) ModTime() time.Time { return self.modified }
func (self *FakeFileInfo) IsDir() bool        { return filepath.Ext(self.path) == "" }
func (self *FakeFileInfo) Sys() interface{}   { return nil }

func newFileInfo(path string, size int64, modified time.Time) *FakeFileInfo {
	self := &FakeFileInfo{}
	self.path = path
	self.size = size
	self.modified = modified
	return self
}

func TestFakeFileSystem(t *testing.T) {
	var fs *FakeFileSystem

	Convey("Subject: FakeFileSystem", t, func() {
		fs = NewFakeFileSystem()

		Convey("When walking a barren file system", func() {
			step := func(path string, info os.FileInfo, err error) error { panic("Should NOT happen!") }

			Convey("The step function should never be called", func() {
				So(func() { fs.Walk("/", step) }, ShouldNotPanic)
			})
		})

		Convey("When a file system is populated...", func() {
			first, second, third := time.Now(), time.Now(), time.Now()
			fs.Create("/a", 1, first)
			fs.Create("/b", 2, second)
			fs.Create("/c", 3, third)

			Convey("...and then walked", func() {
				paths, names, sizes, times, errors := []string{}, []string{}, []int64{}, []time.Time{}, []error{}
				fs.Walk("/", func(path string, info os.FileInfo, err error) error {
					paths = append(paths, path)
					names = append(names, info.Name())
					sizes = append(sizes, info.Size())
					times = append(times, info.ModTime())
					errors = append(errors, err)
					return nil
				})

				Convey("Each path should be visited once", func() {
					So(paths, ShouldResemble, []string{"/a", "/b", "/c"})
					So(names, ShouldResemble, []string{"a", "b", "c"})
					So(sizes, ShouldResemble, []int64{1, 2, 3})
					So(times, ShouldResemble, []time.Time{first, second, third})
					So(errors, ShouldResemble, []error{nil, nil, nil})
				})
			})

		})

		Convey("When a directory does NOT exist it should NOT be found", func() {
			So(fs.Exists("/not/there"), ShouldBeFalse)
		})

		Convey("When a folder is created", func() {
			modified := time.Now()
			fs.Create("/path/to/folder", 3, modified)

			Convey("It should be visible as a folder", func() {
				So(fs.Exists("/path/to/folder"), ShouldBeTrue)
			})
		})

		Convey("When a file is created", func() {
			fs.Create("/path/to/file.txt", 3, time.Now())

			Convey("It should NOT be visible as a folder", func() {
				So(fs.Exists("/path/to/file.txt"), ShouldBeFalse)
			})
		})
	})
}
