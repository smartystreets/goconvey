package main

import (
	"encoding/json"
	"fmt"
	"os"
	_ "os/exec"
	"path/filepath"
	"strings"
	"time"
)

func updateWatch(root string) {
	addWatches(root)
	removeWatches()
}
func addWatches(root string) {
	if rootWatch != root {
		// TODO: set gopath...
		adjustRoot(root)
	}

	watchNestedPaths(root)
}
func adjustRoot(root string) {
	fmt.Println("Watching new root:", root)
	for path, _ := range watched {
		removeWatch(path)
	}
	rootWatch = root
	watch(root)
}
func watchNestedPaths(root string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if matches, _ := filepath.Glob(filepath.Join(path, "*test.go")); len(matches) > 0 {
			watch(path)
		}
		return nil
	})
}
func watch(path string) {
	if !watching(path) {
		fmt.Println("Watching:", path)
		watched[path] = true
		watcher.Watch(path)
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

func removeWatches() {
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

func reactToChanges() {
	busy := false
	done := make(chan bool)
	ready := make(chan bool)

	for {
		select {
		case ev := <-watcher.Event:
			updateWatch(rootWatch)
			if strings.HasSuffix(ev.Name, ".go") && !busy {
				busy = true
				go runTests(done)
			}

		case err := <-watcher.Error:
			fmt.Println(err)

		case <-done:
			// TODO: rethink this delay?
			time.AfterFunc(1500*time.Millisecond, func() {
				ready <- true
			})

		case <-ready:
			busy = false
		}
	}
}

func runTests(done chan bool) {
	results := []PackageResult{}
	// for path, _ := range watched {
	// fmt.Println("Running tests at:", path)
	// if err := os.Chdir(path); err != nil {
	// fmt.Println("Could not chdir to:", path)
	// continue
	// }
	// output, err := exec.Command("go", "test", "-json").Output()
	// if err != nil {
	// fmt.Printf("Error from test execution at %s. Error: %v\n", path, err)
	// continue // TODO: is the error expected on failure?
	// }
	// result := parsePackageResult(string(output))
	// fmt.Println("Result: ", result.Passed)
	// results = append(results, result)
	// }
	serialized, err := json.Marshal(results)
	if err != nil {
		fmt.Println("Problem serializing json test results!", err)
	} else {
		latestOutput = string(serialized)
	}
	done <- true
}
