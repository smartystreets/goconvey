package main

import (
	"flag"
	"fmt"
	"github.com/howeyc/fsnotify"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func init() {
	flag.IntVar(&port, "port", 8080, "The port at which to serve http.")
	flag.StringVar(&host, "host", "127.0.0.1", "The host at which to serve http.")
	watched = make(map[string]bool)
	done = make(chan bool)

	startWatcher()
}
func startWatcher() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
}
func ensureEnoughOpenFiles() {
	exec.Command("ulimit", "-n", "4096").Run()
}

func main() {
	flag.Parse()
	defer watcher.Close()

	working, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	updateWatch(working)

	go reactToChanges()

	_, file, _, _ := runtime.Caller(0)
	here := filepath.Dir(file)
	static := filepath.Join(here, "..", "client")
	http.Handle("/", http.FileServer(http.Dir(static)))
	http.HandleFunc("/watch", watchHandler)
	http.HandleFunc("/latest", reportHandler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func reportHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
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
	if !isDirectory(newRoot) {
		http.Error(writer, "The 'root' value provided is not an existing directory.", http.StatusNotFound)
	} else {
		updateWatch(newRoot)
		go runTests(done)
	}
}

var (
	port         int
	host         string
	latestOutput string
	rootWatch    string
	watched      map[string]bool
	watcher      *fsnotify.Watcher
	done         chan bool
)

const separator = string(filepath.Separator)
