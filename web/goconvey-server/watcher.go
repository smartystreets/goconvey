package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func updateWatch(root string) {
	addNewWatches(root)
	removeExpiredWatches()
}
func addNewWatches(root string) {
	if rootWatch != root {
		setGoPath(root)
		adjustRoot(root)
	}

	watchNestedPaths(root)
}
func setGoPath(root string) {
	goPath := root
	fmt.Println("GOPATH (before):", os.Getenv("GOPATH"))
	index := strings.Index(root, separator+"src")
	if index > -1 {
		goPath = root[:index]
	}
	err := os.Setenv("GOPATH", goPath)
	if err != nil {
		fmt.Println("Error setting GOPATH:", err)
	}

	fmt.Println("GOPATH (after):", os.Getenv("GOPATH"))
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
		if matches, _ := filepath.Glob(filepath.Join(path, "*.go")); len(matches) > 0 {
			addWatch(path)
		}
		return nil
	})
}
func addWatch(path string) {
	if watching(path) {
		return
	}

	if !looksLikeGoPackage(path) {
		return
	}

	watched[path] = true
	watcher.Watch(path)
	fmt.Println("Watching:", path)
}

func looksLikeGoPackage(path string) bool {
	_, resolved := resolvePackageName(path)
	return resolved
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
		if !isDirectory(path) {
			removeWatch(path)
		}
	}
}
func removeWatch(path string) {
	delete(watched, path)
	watcher.RemoveWatch(path)
	fmt.Println("No longer watching:", path)
}

func isDirectory(directory string) bool {
	info, err := os.Stat(directory)
	return err == nil && info.IsDir()
}
