package system

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileSystem struct{}

func (self *FileSystem) Walk(root string, step filepath.WalkFunc) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if self.isMetaDirectory(info) {
			return filepath.SkipDir
		}

		return step(path, info, err)
	})

	if err != nil {
		log.Println("Error while walking file system:", err)
		panic(err)
	}
}

func (self *FileSystem) isMetaDirectory(info os.FileInfo) bool {
	return info.IsDir() && strings.HasPrefix(info.Name(), ".")
}

func (self *FileSystem) Exists(directory string) bool {
	info, err := os.Stat(directory)
	return err == nil && info.IsDir()
}

func NewFileSystem() *FileSystem {
	return &FileSystem{}
}
