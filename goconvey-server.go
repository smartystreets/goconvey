package main

import (
	"flag"
	"fmt"
	"github.com/smartystreets/goconvey/web/server/api"
	"github.com/smartystreets/goconvey/web/server/contract"
	exec "github.com/smartystreets/goconvey/web/server/executor"
	parse "github.com/smartystreets/goconvey/web/server/parser"
	"github.com/smartystreets/goconvey/web/server/system"
	watch "github.com/smartystreets/goconvey/web/server/watcher"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	// dependencies:
	_ "github.com/jacobsa/oglematchers"
	_ "github.com/jacobsa/ogletest"
)

func init() {
	oneSecond, _ := time.ParseDuration("1s")
	flag.IntVar(&port, "port", 8081, "The port at which to serve http.")
	flag.StringVar(&host, "host", "127.0.0.1", "The host at which to serve http.")
	flag.DurationVar(&nap, "poll", oneSecond, "The interval to wait between polling the file system for changes (default: 1s)")
}

func main() {
	flag.Parse()

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
	http.HandleFunc("/latest", server.Results)
	http.HandleFunc("/execute", server.Execute)
	http.HandleFunc("/status", server.Status)
}

func activateServer() {
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func wireup() (*contract.Monitor, contract.Server) {
	working, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fs := system.NewFileSystem()
	shell := system.NewShell()

	watcher := watch.NewWatcher(fs, shell)
	watcher.Adjust(working)

	parser := parse.NewParser(parse.ParsePackageResults)
	tester := exec.NewConcurrentTester(shell)
	executor := exec.NewExecutor(tester, parser)
	server := api.NewHTTPServer(watcher, executor)
	scanner := watch.NewScanner(fs, watcher)
	monitor := contract.NewMonitor(scanner, watcher, executor, server, func() { time.Sleep(nap) })

	return monitor, server
}

var (
	port int
	host string
	nap  time.Duration
)
