package watch

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDepthLimit(t *testing.T) {
	fileSystem := []FileSystemItem{
		FileSystemItem{
			Root: "/hello/world",
			Path: "/hello/world/1/2/3",
		},
		FileSystemItem{
			Root: "/hello/world",
			Path: "/hello/world/1/2/3/4",
		},
		FileSystemItem{
			Root: "/hello/world",
			Path: "/hello/world/1/2/3/4/5",
		},
	}

	Convey("When a negative depth limit is specified", t, func() {
		filtered := LimitDepth(fileSystem, -1)

		Convey("The original collection should returned without modification", func() {
			So(filtered, ShouldResemble, fileSystem)
		})
	})

	Convey("When a nonzero depth limit is specified", t, func() {
		filtered := LimitDepth(fileSystem, 4)

		Convey("All items that are deeper than the limit should be filtered out", func() {
			So(filtered, ShouldResemble, fileSystem[:2])
		})
	})
}

func TestChecksum(t *testing.T) {
	fileSystem := []FileSystemItem{

		// directory; only ever counts as 1 'point'
		FileSystemItem{
			Root:     "/",
			Path:     "/hello",
			Name:     "hello",
			Size:     12345566645,
			Modified: 712342134,
			IsFolder: true,
		},

		// not go file; always ignored
		FileSystemItem{
			Root:     "/",
			Path:     "/1/hello/world.txt",
			Name:     "world.txt",
			Size:     3,
			Modified: 5,
			IsFolder: false,
		},

		// go file; will count as Size + Modified 'points'
		FileSystemItem{
			Root:     "/",
			Path:     "/1/2/3/4/5/hello/world.go",
			Name:     "world.go",
			Size:     3,
			Modified: 5,
			IsFolder: false,
		},

		// .dot file; always ignored
		FileSystemItem{
			Root:     "/",
			Path:     "/hello/.world.go",
			Name:     ".world.go",
			Size:     3,
			Modified: 5,
			IsFolder: false,
		},

		// .dot directory; always ignored
		FileSystemItem{
			Root:     "/",
			Path:     "/.hello",
			Name:     ".hello",
			Size:     3,
			Modified: 5,
			IsFolder: true,
		},

		// .dot directory contents; always ignored
		FileSystemItem{
			Root:     "/",
			Path:     "/.hello/world.go",
			Name:     "world.go",
			Size:     3,
			Modified: 5,
			IsFolder: false,
		},

		// .dot directory contents; always ignored
		FileSystemItem{
			Root:     "/",
			Path:     "/hello/hi.goconvey",
			Name:     "hi.goconvey",
			Size:     2,
			Modified: 3,
			IsFolder: false,
		},
	}

	Convey("The file system should be checksummed correctly", t, func() {
		So(Checksum(fileSystem), ShouldEqual, 14)
	})
}
