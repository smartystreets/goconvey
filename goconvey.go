// This executable provides an HTTP server that watches for file system changes
// to .go files within the working directory (and all nested go packages).
// Navigating to the configured host and port will show a web UI showing the
// results of running `go test` in each go package.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/smartystreets/goconvey/web/server/api"
	"github.com/smartystreets/goconvey/web/server/contract"
	exec "github.com/smartystreets/goconvey/web/server/executor"
	parse "github.com/smartystreets/goconvey/web/server/parser"
	"github.com/smartystreets/goconvey/web/server/system"
	watch "github.com/smartystreets/goconvey/web/server/watcher"
)

func init() {
	quarterSecond, _ := time.ParseDuration("250ms")
	flag.IntVar(&port, "port", 8080, "The port at which to serve http.")
	flag.StringVar(&host, "host", "127.0.0.1", "The host at which to serve http.")
	flag.DurationVar(&nap, "poll", quarterSecond, "The interval to wait between polling the file system for changes (default: 250ms).")
	flag.IntVar(&packages, "packages", 10, "The number of packages to test in parallel. Higher == faster but more costly in terms of computing. (default: 10)")
	flag.StringVar(&gobin, "gobin", "go", "The path to the 'go' binary (default: search on the PATH)")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	flag.Parse()
	log.Printf("Initial configuration: [host: %s] [port: %d] [poll %v]\n", host, port, nap)

	monitor, server := wireup()

	go monitor.ScanForever()

	serveHTTP(server)
}

func serveHTTP(server contract.Server) {
	serveStaticResources()
	serveAjaxMethods(server)
	activateServer()
}

func serveStaticResources() {
	_, file, _, _ := runtime.Caller(0)
	here := filepath.Dir(file)
	static := filepath.Join(here, "/web/client")
	http.Handle("/", http.FileServer(http.Dir(static)))
}

func serveAjaxMethods(server contract.Server) {
	http.HandleFunc("/watch", server.Watch)
	http.HandleFunc("/ignore", server.Ignore)
	http.HandleFunc("/reinstate", server.Reinstate)
	http.HandleFunc("/latest", server.Results)
	http.HandleFunc("/execute", server.Execute)
	http.HandleFunc("/status", server.Status)
	http.HandleFunc("/status/poll", server.LongPollStatus)
}

func activateServer() {
	log.Printf("Serving HTTP at: http://%s:%d\n", host, port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func wireup() (*contract.Monitor, contract.Server) {
	log.Println("Constructing components...")
	working, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fs := system.NewFileSystem()
	shell := system.NewShell(gobin)

	watcher := watch.NewWatcher(fs, shell)
	watcher.Adjust(working)

	parser := parse.NewParser(parse.ParsePackageResults)
	tester := exec.NewConcurrentTester(shell)
	tester.SetBatchSize(packages)

	statusNotif := make(chan bool, 1)
	executor := exec.NewExecutor(tester, parser, statusNotif)
	server := api.NewHTTPServer(watcher, executor, statusNotif)
	scanner := watch.NewScanner(fs, watcher)
	monitor := contract.NewMonitor(scanner, watcher, executor, server, sleeper)

	return monitor, server
}

func sleeper() {
	time.Sleep(nap)
}

var (
	port     int
	host     string
	gobin    string
	nap      time.Duration
	packages int
)
