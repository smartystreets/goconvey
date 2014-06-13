package system

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDepthLimitingFileSystem(t *testing.T) {
	Convey("Subject: Depth-limited file system", t, func() {
		inner := NewMockFileSystem()

		Convey("When the depth limit is set to -1", func() {
			files := NewDepthLimit(inner, -1)
			files.Walk(root, inner.step)

			Convey("No depth limiting should be enforced", func() {
				So(inner.walked, ShouldResemble, []string{
					folder,
					subfolder,
					subsubfolder,
				})
			})
		})

		Convey("When the depth limit is not negative", func() {
			files := NewDepthLimit(inner, 1)
			files.Walk(root, inner.step)

			Convey("Directories outside the depth should be skipped", func() {
				So(inner.walked, ShouldResemble, []string{
					folder,
					subfolder,
				})
			})
		})

		Convey("When requesting a listing", func() {
			files := NewDepthLimit(inner, -1)
			listing, indicator := files.Listing("hi")

			Convey("The request should be forwarded to the inner file system", func() {
				So(inner.listingCalled, ShouldBeTrue)
			})

			Convey("The inner result should be returned", func() {
				So(listing, ShouldResemble, innerListing)
				So(indicator, ShouldEqual, listingIndicator)
			})
		})

		Convey("When checking the existence of a directory", func() {
			files := NewDepthLimit(inner, -1)
			exists := files.Exists("hi")

			Convey("The request should be forwarded to the inner file system", func() {
				So(inner.existsCalled, ShouldBeTrue)
			})

			Convey("The inner result should be returned", func() {
				So(exists, ShouldBeTrue)
			})
		})
	})
}

//////////////////////////////

type MockFileSystem struct {
	paths         []*FakeFileInfo
	walked        []string
	listingCalled bool
	existsCalled  bool
}

func (self *MockFileSystem) Walk(root string, step filepath.WalkFunc) {
	for _, path := range self.paths {
		step(path.path, path, nil)
	}
}

func (self *MockFileSystem) step(path string, info os.FileInfo, err error) error {
	self.walked = append(self.walked, path)
	return err
}

func (self *MockFileSystem) Listing(directory string) ([]os.FileInfo, error) {
	self.listingCalled = true
	return innerListing, listingIndicator
}

func (self *MockFileSystem) Exists(directory string) bool {
	self.existsCalled = true
	return true
}

func NewMockFileSystem() *MockFileSystem {
	self := new(MockFileSystem)
	self.paths = []*FakeFileInfo{
		newFileInfo(folder, 42, time.Now()),
		newFileInfo(subfolder, 42, time.Now()),
		newFileInfo(subsubfolder, 42, time.Now()),
	}
	return self
}

const (
	root         = slash                   // - /
	folder       = root + "1"              // - /1
	subfolder    = folder + slash + "2"    // - /1/2
	subsubfolder = subfolder + slash + "3" // - /1/2/3
)

var listingIndicator = errors.New("Listing was called.")
var innerListing = []os.FileInfo{newFileInfo(folder, 42, time.Now())}
