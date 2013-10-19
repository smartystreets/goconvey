package system

import (
	"path/filepath"
)

type FileSystem struct{}

func (self *FileSystem) Walk(root string, step filepath.WalkFunc) {}
func (self *FileSystem) Exists(directory string) bool             { return false }
