package watcher

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"os"
)

type Scanner struct {
	fs      contract.FileSystem
	watcher contract.Watcher
	current int64
	working int64
}

// TODO: we must find out if a new directory has been added or if one has been removed!

func (self *Scanner) Scan(root string) (changed bool) {
	self.calculateWorkingChecksum(root)
	return self.determineIfChanged()
}

func (self *Scanner) determineIfChanged() bool {
	defer func() {
		self.current = self.working
		self.working = 0
	}()

	return self.working != self.current
}

func (self *Scanner) calculateWorkingChecksum(root string) {
	self.fs.Walk(root, func(path string, info os.FileInfo, err error) error {
		step := newWalkStep(root, path, info)

		if step.isIgnored(self.watcher) {
			return nil
		}

		if step.isNewWatchedFolder(self.watcher) {
			self.working += step.sum()
		} else if step.isWatchedFile(self.watcher) {
			self.working += step.sum()
		}

		return nil
	})
}

func NewScanner(fs contract.FileSystem, watcher contract.Watcher) *Scanner {
	self := &Scanner{}
	self.fs = fs
	self.watcher = watcher
	return self
}
