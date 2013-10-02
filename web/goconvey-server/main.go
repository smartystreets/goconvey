package main

import (
	"encoding/json"
	"fmt"
	"github.com/howeyc/fsnotify"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	go reactToChanges(watcher)

	working, _ := os.Getwd()
	watcher.Watch(working)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/latest", reportHandler)
	http.ListenAndServe(":8080", nil)
}

func reactToChanges(watcher *fsnotify.Watcher) {
	busy := false
	done := make(chan bool)
	ready := make(chan bool)

	for {
		select {
		case ev := <-watcher.Event:
			if strings.HasSuffix(ev.Name, ".go") && !busy {
				busy = true
				go test(done)
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

func test(done chan bool) {
	fmt.Println("Running tests...")

	// TODO: recurse into subdirectories and run tests...
	// oh yeah, and always check for new packages sprouting up,
	// or existing ones being removed...
	output, _ := exec.Command("go", "test", "-json").Output()
	result := parsePackageResult(string(output))

	serialized, _ := json.Marshal(result)
	// var buffer bytes.Buffer
	// json.Indent(&buffer, serialized, "", "  ")

	latest = string(serialized)
	done <- true
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "<html>...</html>") // TODO
}

func reportHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, latest)
}

var latest string
