package watcher

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"os"
)

type Scanner struct {
	fs                 contract.FileSystem
	watcher            contract.Watcher
	previous           int64
	latestFolders      map[string]bool
	preExistingFolders map[string]bool
}

func (self *Scanner) Scan(root string) bool {
	checksum, folders := self.analyzeCurrentFileSystemState(root)
	self.notifyWatcherOfChangesInFolderStructure(folders)
	return self.latestTestResultsAreStale(checksum)
}

func (self *Scanner) analyzeCurrentFileSystemState(root string) (checksum int64, folders map[string]bool) {
	folders = make(map[string]bool)

	self.fs.Walk(root, func(path string, info os.FileInfo, err error) error {
		step := newWalkStep(root, path, info, self.watcher)
		step.IncludeIn(folders)
		checksum += step.Sum()
		return nil
	})
	return checksum, folders
}

func (self *Scanner) notifyWatcherOfChangesInFolderStructure(latest map[string]bool) {
	self.accountForDeletedFolders(latest)
	self.accountForNewFolders(latest)
	self.preExistingFolders = latest
}
func (self *Scanner) accountForDeletedFolders(latest map[string]bool) {
	for folder, _ := range self.preExistingFolders {
		if _, exists := latest[folder]; !exists {
			self.watcher.Deletion(folder)
		}
	}
}
func (self *Scanner) accountForNewFolders(latest map[string]bool) {
	for folder, _ := range latest {
		if _, exists := self.preExistingFolders[folder]; !exists {
			self.watcher.Creation(folder)
		}
	}
}

func (self *Scanner) latestTestResultsAreStale(checksum int64) bool {
	defer func() { self.previous = checksum }()
	return self.previous != checksum
}

func NewScanner(fs contract.FileSystem, watcher contract.Watcher) *Scanner {
	self := &Scanner{}
	self.fs = fs
	self.watcher = watcher
	self.latestFolders = make(map[string]bool)
	self.preExistingFolders = make(map[string]bool)
	self.rememberCurrentlyWatchedFolders()

	return self
}
func (self *Scanner) rememberCurrentlyWatchedFolders() {
	for _, item := range self.watcher.WatchedFolders() {
		self.preExistingFolders[item.Path] = true
	}
}
