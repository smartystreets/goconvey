package main

import (
	"encoding/json"
	"fmt"
	"github.com/howeyc/fsnotify"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	watcher, _ = fsnotify.NewWatcher() // TODO: err
	defer watcher.Close()

	fmt.Println("Initialized watcher...")

	go reactToChanges()

	working, _ := os.Getwd() // TODO: err
	updateWatch(working)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/latest", reportHandler)
	http.HandleFunc("/watch", watchHandler)
	http.ListenAndServe(":8080", nil) // TODO: flag for port
}

func updateWatch(root string) {
	addWatches(root)
	removeWatches()
}
func addWatches(root string) {
	if rootWatch != root {
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

func exists(path string) bool {
	info, err := os.Stat(path)
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
	for path, _ := range watched {
		fmt.Println("Running tests at:", path)
		os.Chdir(path)                                            // TODO: err
		output, _ := exec.Command("go", "test", "-json").Output() // TODO: err
		result := parsePackageResult(string(output))
		fmt.Println("Result: ", result.Passed)
		results = append(results, result)
	}
	serialized, _ := json.Marshal(results) // TODO: err
	latestOutput = string(serialized)
	done <- true
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "<html>...</html>") // TODO: setup static handler for html and javascript?
}

func reportHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, latestOutput)
}

func watchHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		writer.Write([]byte(rootWatch))
		return
	}

	value := request.URL.Query()["root"]
	if len(value) == 0 {
		http.Error(writer, "No 'root' query string parameter included!", http.StatusBadRequest)
		return
	}
	newRoot := value[0]
	if !exists(newRoot) {
		http.Error(writer, "The 'root' value provided is not an existing directory.", http.StatusNotFound)
	} else {
		updateWatch(newRoot)
	}
}

var (
	latestOutput string
	rootWatch    string
	watched      map[string]bool
	watcher      *fsnotify.Watcher
)
