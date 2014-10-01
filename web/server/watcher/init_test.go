package watcher

import (
	"os"
	"strings"
)

const slash = string(os.PathSeparator)

// normalize replaces all slashes (back/forward) with os.PathSeparator.
func normalize(path string) string {
	return strings.Replace(strings.Replace(path, "/", slash, -1), "\\", slash, -1)
}
