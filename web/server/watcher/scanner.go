package watcher

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"os"
)

type Scanner struct {
	fs                 contract.FileSystem
	watcher            contract.Watcher
	currentChecksum    int64
	workingChecksum    int64
	latestFolders      map[string]bool
	preExistingFolders map[string]bool
}

func (self *Scanner) Scan(root string) (changed bool) {
	self.analyzeCurrentFileSystemState(root)
	self.accountForRecentlyDeletedFolders()
	return self.latestTestResultsAreStale()
}

func (self *Scanner) analyzeCurrentFileSystemState(root string) {
	self.fs.Walk(root, func(path string, info os.FileInfo, err error) error {
		step := newWalkStep(root, path, info, self.watcher)

		step.includeIn(self.latestFolders)

		if step.isIgnored() {
			return nil
		}

		if step.isNewWatchedFolder() {
			self.workingChecksum += step.sum()
			step.registerWatchedFolder()
		} else if step.isCurrentlyWatchedFolder() {
			self.workingChecksum += step.sum()
		} else if step.isWatchedFile() {
			self.workingChecksum += step.sum()
		}

		return nil
	})
}

func (self *Scanner) accountForRecentlyDeletedFolders() {
	for folder, _ := range self.preExistingFolders {
		if _, exists := self.latestFolders[folder]; !exists {
			self.watcher.Deletion(folder)
		}
	}
	self.preExistingFolders = self.latestFolders
	self.latestFolders = make(map[string]bool)
}

func (self *Scanner) latestTestResultsAreStale() bool {
	defer func() {
		self.currentChecksum = self.workingChecksum
		self.workingChecksum = 0
	}()

	return self.workingChecksum != self.currentChecksum
}

func NewScanner(fs contract.FileSystem, watcher contract.Watcher) *Scanner {
	self := &Scanner{}
	self.fs = fs
	self.watcher = watcher
	self.latestFolders = make(map[string]bool)
	self.preExistingFolders = make(map[string]bool)
	self.accountForWatchedFolders()

	return self
}
func (self *Scanner) accountForWatchedFolders() {
	for _, item := range self.watcher.WatchedFolders() {
		self.preExistingFolders[item.Path] = true
	}
}
