package watch

import (
	"os"
	"path/filepath"
)

func ScanFileSystem(root string) []FileSystemItem {
	items := make(chan FileSystemItem)

	go func() {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			items <- FileSystemItem{
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

	list := []FileSystemItem{}
	for item := range items {
		list = append(list, item)
	}

	return list
}
