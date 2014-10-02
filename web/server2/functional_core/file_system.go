package functional_core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/smartystreets/goconvey/web/server2/messaging"
)

func LimitDepth(items []messaging.FileSystemItemFoundEvent, depth int) []messaging.FileSystemItemFoundEvent {
	if depth < 0 {
		return items
	}

	filtered := []messaging.FileSystemItemFoundEvent{}
	for _, item := range items {
		nested := item.Path[len(item.Root):]
		if strings.Count(nested, slash) <= depth {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func Checksum(items []messaging.FileSystemItemFoundEvent) int64 {
	var sum int64

	for _, item := range items {
		if item.IsFolder && strings.HasPrefix(item.Name, ".") {
			continue
		}

		if item.IsFolder {
			sum++
			continue
		}

		if filepath.Ext(item.Path) != ".go" {
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
