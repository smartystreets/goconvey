package watch

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

///////////////////////////////////////////////////////////////////////////////

type FileSystemItem struct {
	Root     string
	Path     string
	Name     string
	Size     int64
	Modified int64
	IsFolder bool

	ProfileDisabled  bool
	ProfileArguments []string
}

///////////////////////////////////////////////////////////////////////////////

func YieldFileSystemItems(root string) chan *FileSystemItem {
	items := make(chan *FileSystemItem)

	go func() {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			items <- &FileSystemItem{
				Root:     root,
				Path:     path,
				Name:     info.Name(),
				Size:     info.Size(),
				Modified: info.ModTime().Unix(),
				IsFolder: info.IsDir(),
			}

			return nil
		})
		close(items)
	}()

	return items
}

///////////////////////////////////////////////////////////////////////////////

func ReadContents(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return ""
	}

	return string(content)
}

///////////////////////////////////////////////////////////////////////////////
