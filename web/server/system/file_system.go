package system

import (
	"os"
	"path/filepath"
)

type FileSystem struct{}

func (self *FileSystem) Walk(root string, step filepath.WalkFunc) {
	err := filepath.Walk(root, step)

	if err != nil {
		panic(err) // TODO?
	}
}
func (self *FileSystem) Exists(directory string) bool {
	info, err := os.Stat(directory)
	return err == nil && info.IsDir()
}
