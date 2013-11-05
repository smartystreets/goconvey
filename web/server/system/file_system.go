package system

import (
	"os"
	"path/filepath"
	"strings"
)

type FileSystem struct{}

func (self *FileSystem) Walk(root string, step filepath.WalkFunc) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if self.shouldIgnore(path) {
			return nil
		}

		return step(path, info, err)
	})

	if err != nil {
		panic(err) // TODO?
	}
}

func (self *FileSystem) shouldIgnore(path string) bool {
	return strings.Contains(path, string(filepath.Separator)+".git")
}

func (self *FileSystem) Exists(directory string) bool {
	info, err := os.Stat(directory)
	return err == nil && info.IsDir()
}

func NewFileSystem() *FileSystem {
	return &FileSystem{}
}
