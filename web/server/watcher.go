package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/smartystreets/goconvey/web/server/parser"
	"github.com/smartystreets/goconvey/web/server/results"
	"io"
	"os"
	"os/exec"
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
	done, ready := make(chan bool), make(chan bool)
	busy := true

	go runTests(done)

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
	packageResults := []*results.PackageResult{}
	revision := md5.New()

	fmt.Println("")
	for path, _ := range watched {
		fmt.Printf("Running tests for: %s ...", path)
		if err := os.Chdir(path); err != nil {
			fmt.Println("Could not chdir to:", path)
			continue
		}

		exec.Command("go", "test", "-i").Run()
		output, _ := exec.Command("go", "test", "-v", "-timeout=-42s").CombinedOutput()
		stringOutput := string(output)
		io.WriteString(revision, stringOutput)
		packageIndex := strings.Index(path, "/src/")
		packageName := path[packageIndex+len("/src/"):]
		result := parser.ParsePackageResults(packageName, stringOutput)
		fmt.Printf("[%s]\n", result.Outcome)
		packageResults = append(packageResults, result)
	}

	output := results.CompleteOutput{
		Packages: packageResults,
		Revision: hex.EncodeToString(revision.Sum(nil)),
	}
	serialized, err := json.Marshal(output)
	if err != nil {
		fmt.Println("Problem serializing json test results!", err) // panic?
	} else {
		latestOutput = string(serialized)
	}
	done <- true
}
