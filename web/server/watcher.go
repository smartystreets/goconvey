package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/smartystreets/goconvey/web/server/parser"
	"github.com/smartystreets/goconvey/web/server/results"
	"hash"
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
			panic(err)

		case <-done:
			time.AfterFunc(500*time.Millisecond, func() {
				ready <- true
			})

		case <-ready:
			busy = false
		}
	}
}

func runTests(done chan bool) {
	revision := md5.New()
	packageResults := aggregateResults(revision)

	output := results.CompleteOutput{
		Packages: packageResults,
		Revision: hex.EncodeToString(revision.Sum(nil)),
	}

	remember(output)
	done <- true
}

func aggregateResults(revision hash.Hash) []*results.PackageResult {
	packageResults := []*results.PackageResult{}

	fmt.Println("")
	for path, _ := range watched {
		stringOutput := executeTests(path)
		io.WriteString(revision, stringOutput)
		result := parseTestOutput(path, stringOutput)
		packageResults = append(packageResults, result)
	}
	return packageResults
}

func executeTests(path string) string {
	fmt.Printf("Running tests for: %s ...", path)
	if err := os.Chdir(path); err != nil {
		panic(fmt.Sprintf("Could not chdir to: %s", path))
	}

	exec.Command("go", "test", "-i").Run()
	output, _ := exec.Command("go", "test", "-v", "-timeout=-42s").CombinedOutput()
	return string(output)
}

func parseTestOutput(path, stringOutput string) *results.PackageResult {
	packageIndex := strings.Index(path, "/src/")
	packageName := path[packageIndex+len("/src/"):]
	result := parser.ParsePackageResults(packageName, stringOutput)
	fmt.Printf("[%s]\n", result.Outcome)
	return result
}

func remember(output results.CompleteOutput) {
	serialized, err := json.Marshal(output)
	if err != nil {
		panic(err)
	} else {
		latestOutput = string(serialized)
	}
}
