package watcher

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"os"
	"path/filepath"
	"strings"
)

type walkStep struct {
	root   string
	path   string
	folder string
	info   os.FileInfo
}

func (self *walkStep) sum() int64 {
	return self.info.Size() + self.info.ModTime().Unix()
}

func (self *walkStep) isIgnored(watcher contract.Watcher) bool {
	return watcher.IsIgnored(self.folder)
}

func (self *walkStep) isNewWatchedFolder(watcher contract.Watcher) bool {
	return !watcher.IsWatched(self.folder) &&
		strings.HasPrefix(self.path, self.root) &&
		self.info.IsDir()
}

func (self *walkStep) isWatchedFile(watcher contract.Watcher) bool {
	return watcher.IsActive(self.folder) && filepath.Ext(self.path) == ".go"
}

func newWalkStep(root, path string, info os.FileInfo) *walkStep {
	self := &walkStep{}
	self.root = root
	self.path = path
	self.info = info
	self.folder = deriveFolderName(path, info)
	return self
}

func deriveFolderName(path string, info os.FileInfo) string {
	if info.IsDir() {
		return path
	} else {
		return filepath.Dir(path)
	}
}
