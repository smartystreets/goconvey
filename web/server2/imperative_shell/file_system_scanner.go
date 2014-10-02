package imperative_shell

import (
	"os"
	"path/filepath"

	"github.com/smartystreets/goconvey/web/server2/messaging"
)

func ScanFileSystem(root string) []messaging.FileSystemItemFoundEvent {
	items := make(chan messaging.FileSystemItemFoundEvent)

	go func() {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			items <- messaging.FileSystemItemFoundEvent{
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

	list := []messaging.FileSystemItemFoundEvent{}
	for item := range items {
		list = append(list, item)
	}

	return list
}
