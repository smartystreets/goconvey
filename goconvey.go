// This executable provides an HTTP server that watches for file system changes
// to .go files within the working directory (and all nested go packages).
// Navigating to the configured host and port in a web browser will display the
// latest results of running `go test` in each go package.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go/build"

	"github.com/smartystreets/goconvey/web/server/api"
	"github.com/smartystreets/goconvey/web/server/contract"
	"github.com/smartystreets/goconvey/web/server/executor"
	"github.com/smartystreets/goconvey/web/server/messaging"
	"github.com/smartystreets/goconvey/web/server/parser"
	"github.com/smartystreets/goconvey/web/server/system"
	"github.com/smartystreets/goconvey/web/server/watch"
)

func init() {
	flags()
	folders()
}
func flags() {
	flag.IntVar(&port, "port", 8080, "The port at which to serve http.")
	flag.StringVar(&host, "host", "127.0.0.1", "The host at which to serve http.")
	flag.DurationVar(&nap, "poll", quarterSecond, "The interval to wait between polling the file system for changes (default: 250ms).")
	flag.IntVar(&packages, "packages", 10, "The number of packages to test in parallel. Higher == faster but more costly in terms of computing. (default: 10)")
	flag.StringVar(&gobin, "gobin", "go", "The path to the 'go' binary (default: search on the PATH).")
	flag.BoolVar(&cover, "cover", true, "Enable package-level coverage statistics. Requires Go 1.2+ and the go cover tool. (default: true)")
	flag.IntVar(&depth, "depth", -1, "The directory scanning depth. If -1, scan infinitely deep directory structures. 0: scan working directory. 1+: Scan into nested directories, limited to value. (default: -1)")
	flag.StringVar(&timeout, "timeout", "5s", "The test execution timeout if none is specified in the *.goconvey file (default: 5s).")

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func folders() {
	_, file, _, _ := runtime.Caller(0)
	here := filepath.Dir(file)
	static = filepath.Join(here, "/web/client")
	reports = filepath.Join(static, "reports")
}

func main() {
	flag.Parse()
	log.Printf(initialConfiguration, host, port, nap, cover)

	working, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cover = coverageEnabled(cover, reports)
	shell := system.NewShell(gobin, reports, cover, timeout)

	watcherInput := make(chan messaging.WatcherCommand)
	watcherOutput := make(chan messaging.Folders)
	watcher := watch.NewWatcher(working, depth, nap, watcherInput, watcherOutput)

	parser := parser.NewParser(parser.ParsePackageResults)
	tester := executor.NewConcurrentTester(shell)
	tester.SetBatchSize(packages)

	longpollChan := make(chan chan string)
	executor := executor.NewExecutor(tester, parser, longpollChan)
	server := api.NewHTTPServer(working, watcherInput, executor, longpollChan)

	go runTestOnUpdates(watcherOutput, executor, server)
	go watcher.Listen()
	serveHTTP(server)
}

func runTestOnUpdates(queue chan messaging.Folders, executor contract.Executor, server contract.Server) {
	for update := range queue {
		log.Println("Received request from watcher to execute tests...")
		root := ""
		packages := []*contract.Package{}
		for _, folder := range update {
			root = folder.Root
			hasImportCycle := testFilesImportTheirOwnPackage(folder.Path)
			packages = append(packages, contract.NewPackage(folder, hasImportCycle))
		}
		output := executor.ExecuteTests(packages)
		server.ReceiveUpdate(root, output)
	}
}

// This method exists because of a bug in the go cover tool that
// causes an infinite loop when you try to run `go test -cover`
// on a package that has an import cycle defined in one of it's
// test files. Yuck.
func testFilesImportTheirOwnPackage(packagePath string) bool {
	meta, err := build.ImportDir(packagePath, build.AllowBinary)
	if err != nil {
		return false
	}

	for _, dependency := range meta.TestImports {
		if dependency == meta.ImportPath {
			return true
		}
	}
	return false
}

func serveHTTP(server contract.Server) {
	serveStaticResources()
	serveAjaxMethods(server)
	activateServer()
}

func serveStaticResources() {
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
	http.HandleFunc("/pause", server.TogglePause)
}

func activateServer() {
	log.Printf("Serving HTTP at: http://%s:%d\n", host, port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Println(err)
	}
}

func coverageEnabled(cover bool, reports string) bool {
	return (cover &&
		goVersion_1_2_orGreater() &&
		coverToolInstalled() &&
		ensureReportDirectoryExists(reports))
}
func goVersion_1_2_orGreater() bool {
	version := runtime.Version() // 'go1.2....'
	major, minor := version[2], version[4]
	version_1_2 := major >= byte('1') && minor >= byte('2')
	if !version_1_2 {
		log.Printf(pleaseUpgradeGoVersion, version)
		return false
	}
	return true
}
func coverToolInstalled() bool {
	working, err := os.Getwd()
	if err != nil {
		working = "."
	}
	command := system.NewCommand(working, "go", "tool", "cover").Execute()
	installed := strings.Contains(command.Output, "Usage of 'go tool cover':")
	if !installed {
		log.Print(coverToolMissing)
		return false
	}
	return true
}
func ensureReportDirectoryExists(reports string) bool {
	if exists(reports) {
		return true
	}

	if err := os.Mkdir(reports, 0755); err == nil {
		return true
	}

	log.Printf(reportDirectoryUnavailable, reports)
	return false
}
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

var (
	port     int
	host     string
	gobin    string
	nap      time.Duration
	packages int
	cover    bool
	depth    int
	timeout  string

	static  string
	reports string

	quarterSecond = time.Millisecond * 250
)

const (
	initialConfiguration       = "Initial configuration: [host: %s] [port: %d] [poll: %v] [cover: %v]\n"
	pleaseUpgradeGoVersion     = "Go version is less that 1.2 (%s), please upgrade to the latest stable version to enable coverage reporting.\n"
	coverToolMissing           = "Go cover tool is not installed or not accessible: `go get code.google.com/p/go.tools/cmd/cover`\n"
	reportDirectoryUnavailable = "Could not find or create the coverage report directory (at: '%s'). You probably won't see any coverage statistics...\n"
)
