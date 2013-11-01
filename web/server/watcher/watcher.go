package watcher

import (
	"errors"
	"fmt"
	"github.com/smartystreets/goconvey/web/server/contract"
	"os"
	"path/filepath"
	"strings"
)

type Watcher struct {
	fs      contract.FileSystem
	shell   contract.Shell
	watched map[string]*contract.Package
}

func (self *Watcher) Adjust(root string) error {
	if !self.fs.Exists(root) {
		return errors.New(fmt.Sprintf("Directory does not exist: '%s'", root))
	}

	self.watched = make(map[string]*contract.Package)
	self.fs.Walk(root, self.includeFolders)
	self.setGoPath(root)

	return nil
}
func (self *Watcher) includeFolders(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		self.watched[path] = &contract.Package{Active: true, Path: path, Name: info.Name()}
	}
	return nil
}
func (self *Watcher) setGoPath(root string) {
	if rootGoPathEnd := strings.Index(root, "/src"); rootGoPathEnd >= 0 {
		self.shell.Setenv("GOPATH", root[:rootGoPathEnd])
	} else {
		self.shell.Setenv("GOPATH", root)
	}
}

func (self *Watcher) Deletion(folder string) {
	delete(self.watched, folder)
}

func (self *Watcher) Creation(folder string) {
	self.watched[folder] = &contract.Package{Active: true, Path: folder, Name: filepath.Base(folder)}
}

func (self *Watcher) Ignore(folder string) error {
	// TODO: a package name is all we have from the server!!!!
	if value, exists := self.watched[folder]; exists {
		value.Active = false
	}
	return nil
}
func (self *Watcher) Reinstate(folder string) error {
	// TODO: a package name is all we have from the server!!!!
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

func NewWatcher(fs contract.FileSystem, shell contract.Shell) *Watcher {
	self := &Watcher{}
	self.fs = fs
	self.shell = shell
	self.watched = map[string]*contract.Package{}
	return self
}
