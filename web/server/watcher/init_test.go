package watcher

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

const slash = string(os.PathSeparator)

// normalize replaces all slashes (back/forward) with os.PathSeparator.
func normalize(path string) string {
	return strings.Replace(strings.Replace(path, "/", slash, -1), "\\", slash, -1)
}
