package watcher

type Scanner struct {
	// contains the FileSystem
	// contains the Watcher
	// NOTE: will have to signal creation and deletion of packages to the Watcher
}

func (self *Scanner) Scan(root string) (changed bool) { return false }
