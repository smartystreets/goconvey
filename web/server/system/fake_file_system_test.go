package system

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"time"
)

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

		Convey("When an existing file system item is modified", func() {
			fs.Create("/a.txt", 1, time.Now())
			fs.Modify("/a.txt")
			var size int64

			Convey("And the file system is then walked", func() {
				fs.Walk("/", func(path string, info os.FileInfo, err error) error {
					size = info.Size()
					return nil
				})
				Convey("The modification should be persistent", func() {
					So(size, ShouldEqual, 2)
				})
			})
		})

		Convey("When an existing file system item is deleted", func() {
			fs.Create("/a.txt", 1, time.Now())
			fs.Delete("/a.txt")
			var found bool

			Convey("And the file system is then walked", func() {
				fs.Walk("/", func(path string, info os.FileInfo, err error) error {
					if info.Name() == "a.txt" {
						found = true
					}
					return nil
				})
				Convey("The deleted entry should NOT be visited", func() {
					So(found, ShouldBeFalse)
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
