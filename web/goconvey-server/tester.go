package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/howeyc/fsnotify"
	"github.com/smartystreets/goconvey/web/goconvey-server/parser"
	"github.com/smartystreets/goconvey/web/goconvey-server/results"
)

func reactToChanges() {
	busy := true
	ready := make(chan bool)

	go runTests(done)

	for {
		select {
		case event := <-watcher.Event:
			if busy {
				continue
			}
			if isIgnored(event) {
				continue
			}
			if event.IsCreate() || watchRemoved(event) || goFileTouched(event) {
				busy = true
				go runTests(done)
			}

		case err := <-watcher.Error:
			panic(err)

		case <-done:
			time.AfterFunc(100*time.Millisecond, func() {
				ready <- true
			})

		case <-ready:
			busy = false
		}
	}
}

func isIgnored(event *fsnotify.FileEvent) bool {
	ignoredFilenames := []string{
		".DS_Store",
		"Thumbs.db",
		"__MAC_OSX",
	}

	for _, ignoredName := range ignoredFilenames {
		if strings.HasSuffix(event.Name, ignoredName) {
			return true
		}
	}
	return false
}

func goFileTouched(event *fsnotify.FileEvent) bool {
	return (event.IsModify() || event.IsRename()) && strings.HasSuffix(event.Name, ".go")
}

func watchRemoved(event *fsnotify.FileEvent) bool {
	return event.IsDelete() && watching(event.Name)
}

func runTests(done chan bool) {
	updateWatch(rootWatch)
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
		io.WriteString(revision, result.Output)
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
	packageName, _ := resolvePackageName(path)
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
		packageName, _ := resolvePackageName(path)
		exec.Command("go", "test", "-i", packageName).Run()
	}
}
func resolvePackageName(path string) (string, bool) {
	success := true
	const endGoPath = separator + "src" + separator
	index := strings.Index(path, endGoPath)
	if index < 0 {
		success = false
	}
	name := path[index+len(endGoPath):]
	return name, success
}
func testPackage(name string) string {
	fmt.Printf("Testing %s ...\n", name)
	output, err := exec.Command("go", "test", "-v", "-timeout=-42s", name).CombinedOutput()
	if len(output) == 0 && err != nil {
		panic(err)
	}
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
