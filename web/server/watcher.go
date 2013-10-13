package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/howeyc/fsnotify"
	"github.com/smartystreets/goconvey/web/server/parser"
	"github.com/smartystreets/goconvey/web/server/results"
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

func reactToChanges() {
	// TODO: encapsulate in a struct to reduce parameter passing (and facilitate testing?)

	done := make(chan bool)
	go runTests(done)
	engageWithFileSystem(done)
}
func engageWithFileSystem(done chan bool) {
	busy := true
	ready := make(chan bool)

	for {
		select {
		case event := <-watcher.Event:
			busy = runFullTestSuite(event, busy, done)

		case err := <-watcher.Error:
			panic(err)

		case <-done:
			prepareForNextTestRun(ready)

		case <-ready:
			busy = false
		}
	}
}
func runFullTestSuite(event *fsnotify.FileEvent, busy bool, done chan bool) bool {
	updateWatch(rootWatch)
	if strings.HasSuffix(event.Name, ".go") && !busy {
		go runTests(done)
		return true
	}
	return false
}
func prepareForNextTestRun(ready chan bool) {
	time.AfterFunc(500*time.Millisecond, func() {
		ready <- true
	})
}

func runTests(done chan bool) {
	// TODO: encapsulate in a struct to avoid parameter passing (and facilitate testing?)

	input, output := make(chan string), make(chan *TestPackage)
	spawnTestExecutors(input, output)
	go scheduleTestExecution(input)
	result := aggregateResults(output)
	remember(result)
	done <- true
}
func spawnTestExecutors(input chan string, output chan *TestPackage) {
	for i := 0; i < len(watched); i++ {
		go worker(input, output)
	}
}
func worker(in chan string, out chan *TestPackage) {
	for path := range in {
		out <- executeTests(path)
	}
}
func scheduleTestExecution(input chan string) {
	for folder, _ := range watched {
		input <- folder
	}
}
func aggregateResults(output chan *TestPackage) results.CompleteOutput {
	revision := md5.New()
	var packageResults []*results.PackageResult

	for _ = range watched {
		result := <-output
		io.WriteString(revision, result.Path)
		packageResults = append(packageResults, result.Parsed)
		fmt.Printf("Result for %s: [%s]\n", result.Parsed.PackageName, result.Parsed.Outcome)
	}

	return results.CompleteOutput{
		Packages: packageResults,
		Revision: hex.EncodeToString(revision.Sum(nil)),
	}
}

func executeTests(path string) *TestPackage {
	buildDependencies()
	packageName := resolvePackageName(path)
	stringOutput := testPackage(packageName)
	result := parser.ParsePackageResults(packageName, stringOutput)

	return &TestPackage{
		Path:   path,
		Output: stringOutput,
		Parsed: result,
	}
}
func buildDependencies() {
	for path, _ := range watched {
		packageName := resolvePackageName(path)
		exec.Command("go", "test", "-i", packageName).Run()
	}
}
func resolvePackageName(path string) string {
	index := strings.Index(path, "/src/")
	return path[index+len("/src/"):]
}
func testPackage(name string) string {
	fmt.Printf("Testing %s ...\n", name)
	output, _ := exec.Command("go", "test", "-v", "-timeout=-42s", name).CombinedOutput()
	return string(output)
}

func remember(output results.CompleteOutput) {
	serialized, err := json.Marshal(output)
	if err != nil {
		panic(err)
	} else {
		latestOutput = string(serialized)
	}
}

type TestPackage struct {
	Path   string
	Output string
	Parsed *results.PackageResult
}
