package system

import (
	"os"
	"path/filepath"
	"time"
)

type FakeFileSystem struct {
	steps []*FakeFileInfo
}

func (self *FakeFileSystem) Create(path string, size int64, modified time.Time) {
	self.steps = append(self.steps, newFileInfo(path, size, modified))
}
func (self *FakeFileSystem) Modify(path string) {
	for _, step := range self.steps {
		if step.path == path {
			step.size++
		}
	}
}
func (self *FakeFileSystem) Delete(path string) {
	for i := 0; i < len(self.steps); i++ {
		if self.steps[i].path == path {
			self.steps = append(self.steps[:i], self.steps[i+1:]...)
			break
		}
	}
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
