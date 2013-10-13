package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func updateWatch(root string) {
	addNewWatches(root)
	removeExpiredWatches()
}
func addNewWatches(root string) {
	if rootWatch != root {
		// TODO: set gopath...
		adjustRoot(root)
	}

	watchNestedPaths(root)
}
func adjustRoot(root string) {
	clearAllWatches()
	addWatch(root)
	rootWatch = root
	fmt.Println("Watching new root:", root)
}
func clearAllWatches() {
	for path, _ := range watched {
		removeWatch(path)
	}
}
func watchNestedPaths(root string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if matches, _ := filepath.Glob(filepath.Join(path, "*test.go")); len(matches) > 0 {
			addWatch(path)
		}
		return nil
	})
}
func addWatch(path string) {
	if !watching(path) {
		watched[path] = true
		watcher.Watch(path)
		fmt.Println("Watching:", path)
	}
}

func watching(path string) bool {
	for w, _ := range watched {
		if w == path {
			return true
		}
	}
	return false
}

func removeExpiredWatches() {
	for path, _ := range watched {
		if !exists(path) {
			removeWatch(path)
		}
	}
}
func removeWatch(path string) {
	delete(watched, path)
	watcher.RemoveWatch(path)
	fmt.Println("No longer watching:", path)
}

func exists(directory string) bool {
	info, err := os.Stat(directory)
	return err == nil && info.IsDir()
}
