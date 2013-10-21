package watcher

import (
	"github.com/smartystreets/goconvey/web/server/system"
)

type Scanner struct {
	fs      system.FileSystem
	watcher Watcher
	// contains the FileSystem
	// contains the Watcher
	// NOTE: will have to signal creation and deletion of packages to the Watcher
}

func (self *Scanner) Scan(root string) (changed bool) { return false }

func NewScanner(fs system.FileSystem, watcher Watcher) *Scanner {
	self := &Scanner{}
	return self
}
