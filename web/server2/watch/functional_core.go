package watch

import (
	"os"
	"path/filepath"
	"strings"
)

type FileSystemItem struct {
	Root     string
	Path     string
	Name     string
	Size     int64
	Modified int64
	IsFolder bool
}

func LimitDepth(items []FileSystemItem, depth int) []FileSystemItem {
	if depth < 0 {
		return items
	}

	filtered := []FileSystemItem{}
	for _, item := range items {
		nested := item.Path[len(item.Root):]
		if strings.Count(nested, slash) <= depth {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func Checksum(items []FileSystemItem) int64 {
	var sum int64

	for _, item := range items {
		if item.IsFolder && strings.HasPrefix(item.Name, ".") {
			continue
		}

		if item.IsFolder {
			sum++
			continue
		}

		extension := filepath.Ext(item.Path)

		if extension != ".go" && extension != ".goconvey" {
			continue
		}

		if strings.HasPrefix(item.Name, ".") {
			continue
		}

		if strings.HasPrefix(filepath.Base(filepath.Dir(item.Path)), ".") {
			continue
		}

		sum += item.Size + item.Modified
	}

	return sum
}

const slash = string(os.PathSeparator)
