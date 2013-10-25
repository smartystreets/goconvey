package watcher

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"os"
	"path/filepath"
	"strings"
)

type Scanner struct {
	fs      contract.FileSystem
	watcher contract.Watcher
	root    string
	watched map[string]bool
	ignored map[string]bool
	current int64
	working int64
}

// TODO: we must find out if a new directory has been added or if one has been removed!

func (self *Scanner) Scan(root string) (changed bool) {
	self.gatherWatchedFolders()
	self.fs.Walk(root, self.calculateWorkingChecksum)
	return self.determineIfChanged()
}

func (self *Scanner) determineIfChanged() bool {
	defer func() {
		self.current = self.working
		self.working = 0
	}()

	return self.working != self.current
}
func (self *Scanner) gatherWatchedFolders() {
	self.watched = map[string]bool{}
	self.ignored = map[string]bool{}

	for i, path := range self.watcher.WatchedFolders() {
		if i == 0 {
			self.root = path.Path
		}
		if path.Active {
			self.watched[path.Path] = true
		} else {
			self.ignored[path.Path] = true
		}
	}
}

func (self *Scanner) calculateWorkingChecksum(path string, info os.FileInfo, err error) error {
	folder := deriveFolderName(path, info)

	if self.ignored[folder] {
		return nil
	}

	if !self.watched[folder] && strings.HasPrefix(path, self.root) && info.IsDir() {
		self.working += info.Size() + info.ModTime().Unix()
	} else if self.watched[folder] && filepath.Ext(path) == ".go" {
		self.working += info.Size() + info.ModTime().Unix()
	}

	return nil
}
func deriveFolderName(path string, info os.FileInfo) string {
	if info.IsDir() {
		return path
	} else {
		return filepath.Dir(path)
	}
}

func NewScanner(fs contract.FileSystem, watcher contract.Watcher) *Scanner {
	self := &Scanner{}
	self.fs = fs
	self.watcher = watcher
	return self
}
