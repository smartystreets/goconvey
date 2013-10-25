package watcher

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"os"
	"path/filepath"
	"strings"
)

type walkStep struct {
	root    string
	path    string
	folder  string
	info    os.FileInfo
	watcher contract.Watcher
}

func (self *walkStep) includeIn(walked map[string]bool) {
	walked[self.folder] = true
}

func (self *walkStep) sum() int64 {
	return self.info.Size() + self.info.ModTime().Unix()
}

func (self *walkStep) isIgnored() bool {
	return self.watcher.IsIgnored(self.folder)
}

func (self *walkStep) isNewWatchedFolder() bool {
	return !self.watcher.IsWatched(self.folder) &&
		strings.HasPrefix(self.path, self.root) &&
		self.info.IsDir()
}

func (self *walkStep) isCurrentlyWatchedFolder() bool {
	return self.watcher.IsActive(self.folder) &&
		strings.HasPrefix(self.path, self.root) &&
		self.info.IsDir()
}

func (self *walkStep) isWatchedFile() bool {
	return self.watcher.IsActive(self.folder) && filepath.Ext(self.path) == ".go"
}

func (self *walkStep) registerWatchedFolder() {
	self.watcher.Creation(self.folder)
}

func newWalkStep(root, path string, info os.FileInfo, watcher contract.Watcher) *walkStep {
	self := &walkStep{}
	self.root = root
	self.path = path
	self.info = info
	self.folder = deriveFolderName(path, info)
	self.watcher = watcher
	return self
}

func deriveFolderName(path string, info os.FileInfo) string {
	if info.IsDir() {
		return path
	} else {
		return filepath.Dir(path)
	}
}
