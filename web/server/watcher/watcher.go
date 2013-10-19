package watcher

import (
	"github.com/smartystreets/goconvey/web/server/contract"
)

type Watcher struct {
	// contains the FileSystem
}

func (self *Watcher) Adjust(root string) error            { return nil }
func (self *Watcher) Deletion(path string)                {}
func (self *Watcher) Creation(path string)                {}
func (self *Watcher) Ignore(path string) error            { return nil }
func (self *Watcher) Reinstate(path string) error         { return nil }
func (self *Watcher) ActivePackages() []*contract.Package { return nil }
