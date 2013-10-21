package watcher

import (
	"errors"
	"fmt"
	"github.com/smartystreets/goconvey/web/server/contract"
	"os"
	"path/filepath"
)

type Watcher struct {
	fs      contract.FileSystem
	watched []*contract.Package
}

func (self *Watcher) Adjust(root string) error {
	if !self.fs.Exists(root) {
		return errors.New(fmt.Sprintf("Directory does not exist: '%s'", root))
	}

	self.watched = self.watched[:0]
	self.fs.Walk(root, self.includeFolders)

	return nil
}
func (self *Watcher) includeFolders(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		self.watched = append(self.watched, &contract.Package{Active: true, Path: path, Name: info.Name()})
	}
	return nil
}

func (self *Watcher) Deletion(folder string) {
	for x := 0; x < len(self.watched); x++ {
		if self.watched[x].Path == folder {
			self.watched = append(self.watched[:x], self.watched[x+1:]...)
			break
		}
	}
}

func (self *Watcher) Creation(folder string) {
	self.watched = append(self.watched, &contract.Package{Active: true, Path: folder, Name: filepath.Base(folder)})
}

func (self *Watcher) Ignore(folder string) error    { return nil }
func (self *Watcher) Reinstate(folder string) error { return nil }
func (self *Watcher) WatchedFolders() []*contract.Package {
	return self.watched
}

func NewWatcher(fs contract.FileSystem) *Watcher {
	self := &Watcher{}
	self.fs = fs
	self.watched = []*contract.Package{}
	return self
}
