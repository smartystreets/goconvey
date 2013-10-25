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
	watched map[string]*contract.Package
}

func (self *Watcher) Adjust(root string) error {
	if !self.fs.Exists(root) {
		return errors.New(fmt.Sprintf("Directory does not exist: '%s'", root))
	}

	self.watched = make(map[string]*contract.Package)
	self.fs.Walk(root, self.includeFolders)

	return nil
}
func (self *Watcher) includeFolders(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		self.watched[path] = &contract.Package{Active: true, Path: path, Name: info.Name()}
	}
	return nil
}

func (self *Watcher) Deletion(folder string) {
	delete(self.watched, folder)
}

func (self *Watcher) Creation(folder string) {
	self.watched[folder] = &contract.Package{Active: true, Path: folder, Name: filepath.Base(folder)}
}

func (self *Watcher) Ignore(folder string) error {
	if value, exists := self.watched[folder]; exists {
		value.Active = false
	}
	return nil
}
func (self *Watcher) Reinstate(folder string) error {
	if value, exists := self.watched[folder]; exists {
		value.Active = true
	}
	return nil
}
func (self *Watcher) WatchedFolders() []*contract.Package {
	i, watched := 0, make([]*contract.Package, len(self.watched))
	for _, item := range self.watched {
		watched[i] = item
		i++
	}
	return watched
}

func (self *Watcher) IsWatched(folder string) bool {
	if value, exists := self.watched[folder]; exists {
		return value.Active
	}
	return false
}

func (self *Watcher) IsIgnored(folder string) bool {
	if value, exists := self.watched[folder]; exists {
		return !value.Active
	}
	return false
}

func NewWatcher(fs contract.FileSystem) *Watcher {
	self := &Watcher{}
	self.fs = fs
	self.watched = map[string]*contract.Package{}
	return self
}
