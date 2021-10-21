// This executable provides an HTTP server that watches for file system changes
// to .go files within the working directory (and all nested go packages).
// Navigating to the configured host and port in a web browser will display the
// latest results of running `go test` in each go package.
package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"golang.org/x/tools/go/packages"

	"github.com/smartystreets/goconvey/web/server/api"
	"github.com/smartystreets/goconvey/web/server/contract"
	"github.com/smartystreets/goconvey/web/server/executor"
	"github.com/smartystreets/goconvey/web/server/messaging"
	"github.com/smartystreets/goconvey/web/server/parser"
	"github.com/smartystreets/goconvey/web/server/system"
	"github.com/smartystreets/goconvey/web/server/watch"
)

func init() {
	flag.IntVar(&port, "port", 8080, "The port at which to serve http.")
	flag.StringVar(&host, "host", "127.0.0.1", "The host at which to serve http.")
	flag.DurationVar(&nap, "poll", quarterSecond, "The interval to wait between polling the file system for changes.")
	flag.IntVar(&parallelPackages, "packages", 10, "The number of packages to test in parallel. Higher == faster but more costly in terms of computing.")
	flag.StringVar(&gobin, "gobin", "go", "The path to the 'go' binary (default: search on the PATH).")
	flag.BoolVar(&cover, "cover", true, "Enable package-level coverage statistics.")
	flag.IntVar(&depth, "depth", -1, "The directory scanning depth. If -1, scan infinitely deep directory structures. 0: scan working directory. 1+: Scan into nested directories, limited to value.")
	flag.StringVar(&timeout, "timeout", "0", "The test execution timeout if none is specified in the *.goconvey file (default is '0', which is the same as not providing this option).")
	flag.StringVar(&watchedSuffixes, "watchedSuffixes", ".go", "A comma separated list of file suffixes to watch for modifications.")
	flag.StringVar(&excludedDirs, "excludedDirs", "vendor,node_modules", "A comma separated list of directories that will be excluded from being watched.")
	flag.StringVar(&workDir, "workDir", "", "set goconvey working directory (default current directory).")
	flag.BoolVar(&autoLaunchBrowser, "launchBrowser", true, "toggle auto launching of browser.")
	flag.BoolVar(&leakTemp, "leakTemp", false, "leak temp dir with coverage reports.")

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	flag.Parse()

	printHeader()

	tmpDir, err := ioutil.TempDir("", "*.goconvey")
	if err != nil {
		log.Fatal(err)
	}
	reports := filepath.Join(tmpDir, "coverage_out")
	if err := os.Mkdir(reports, 0700); err != nil {
		log.Fatal(err)
	}
	if leakTemp {
		log.Printf("leaking temporary directory %q\n", tmpDir)
	} else {
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				log.Printf("failed to clean temporary directory %q: %s\n", tmpDir, err)
			}
		}()
	}

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	working := getWorkDir()
	shell := system.NewShell(gobin, reports, cover, timeout)

	watcherInput := make(chan messaging.WatcherCommand)
	watcherOutput := make(chan messaging.Folders)
	excludedDirItems := strings.Split(excludedDirs, `,`)
	watcher := watch.NewWatcher(working, depth, nap, watcherInput, watcherOutput, watchedSuffixes, excludedDirItems)

	parser := parser.NewParser(parser.ParsePackageResults)
	tester := executor.NewConcurrentTester(shell)
	tester.SetBatchSize(parallelPackages)

	longpollChan := make(chan chan string)
	executor := executor.NewExecutor(tester, parser, longpollChan)
	server := api.NewHTTPServer(working, watcherInput, executor, longpollChan)
	listener := createListener()
	go runTestOnUpdates(watcherOutput, executor, server)
	go watcher.Listen()
	if autoLaunchBrowser {
		go launchBrowser(listener.Addr().String())
	}
	srv := serveHTTP(reports, server, listener)

	<-done
	log.Println("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown: %s\n", err)
	}
}

func printHeader() {
	log.Println("GoConvey server: ")
	serverVersion := "<unknown>"
	if binfo, ok := debug.ReadBuildInfo(); ok {
		serverVersion = binfo.Main.Version
	}
	log.Println("  version:", serverVersion)
	log.Println("  host:", host)
	log.Println("  port:", port)
	log.Println("  poll:", nap)
	log.Println("  cover:", cover)
	log.Println()
}

func browserCmd() (string, bool) {
	browser := map[string]string{
		"darwin":  "open",
		"linux":   "xdg-open",
		"windows": "start",
	}
	cmd, ok := browser[runtime.GOOS]
	return cmd, ok
}

func launchBrowser(addr string) {
	browser, ok := browserCmd()
	if !ok {
		log.Printf("Skipped launching browser for this OS: %s", runtime.GOOS)
		return
	}

	log.Printf("Launching browser on %s", addr)
	url := fmt.Sprintf("http://%s", addr)
	cmd := exec.Command(browser, url)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(output))
}

func runTestOnUpdates(queue chan messaging.Folders, executor contract.Executor, server contract.Server) {
	for update := range queue {
		log.Println("Received request from watcher to execute tests...")
		packages := extractPackages(update)
		output := executor.ExecuteTests(packages)
		root := extractRoot(update, packages)
		server.ReceiveUpdate(root, output)
	}
}

func extractPackages(folderList messaging.Folders) []*contract.Package {
	packageList := []*contract.Package{}
	for _, folder := range folderList {
		if isInsideTestdata(folder) {
			continue
		}
		hasImportCycle := testFilesImportTheirOwnPackage(folder.Path)
		packageName := resolvePackageName(folder.Path)
		packageList = append(
			packageList,
			contract.NewPackage(folder, packageName, hasImportCycle),
		)
	}
	return packageList
}

// For packages that operate on Go source code files, such as Go tooling, it is
// important to have a location that will not be considered part of package
// source to store those files. The official Go tooling selected the testdata
// folder for this purpose, so we need to ignore folders inside testdata.
func isInsideTestdata(folder *messaging.Folder) bool {
	relativePath, err := filepath.Rel(folder.Root, folder.Path)
	if err != nil {
		// There should never be a folder that's not inside the root, but if
		// there is, we can presumably count it as outside a testdata folder as
		// well
		return false
	}

	for _, directory := range strings.Split(filepath.ToSlash(relativePath), "/") {
		if directory == "testdata" {
			return true
		}
	}

	return false
}

func extractRoot(folderList messaging.Folders, packageList []*contract.Package) string {
	path := packageList[0].Path
	folder := folderList[path]
	return folder.Root
}

func createListener() net.Listener {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Println(err)
	}
	if l == nil {
		os.Exit(1)
	}
	return l
}

//go:embed web/client
var static embed.FS

func serveHTTP(reports string, server contract.Server, listener net.Listener) *http.Server {
	webclient, err := fs.Sub(static, "web/client")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.FS(webclient)))

	http.HandleFunc("/watch", server.Watch)
	http.HandleFunc("/ignore", server.Ignore)
	http.HandleFunc("/reinstate", server.Reinstate)
	http.HandleFunc("/latest", server.Results)
	http.HandleFunc("/execute", server.Execute)
	http.HandleFunc("/status", server.Status)
	http.HandleFunc("/status/poll", server.LongPollStatus)
	http.HandleFunc("/pause", server.TogglePause)

	http.Handle("/reports/", http.StripPrefix("/reports/", http.FileServer(http.Dir(reports))))

	log.Printf("Serving HTTP at: http://%s\n", listener.Addr())
	ret := &http.Server{}
	go func() {
		err := ret.Serve(listener)
		if err != nil {
			log.Println(err)
		}
	}()
	return ret
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getWorkDir() string {
	working := ""
	var err error
	if workDir != "" {
		working = workDir
	} else {
		working, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}
	result, err := exists(working)
	if err != nil {
		log.Fatal(err)
	}
	if !result {
		log.Fatalf("Path:%s does not exists", working)
	}
	return working
}

var (
	port              int
	host              string
	gobin             string
	nap               time.Duration
	parallelPackages  int
	cover             bool
	depth             int
	timeout           string
	watchedSuffixes   string
	excludedDirs      string
	autoLaunchBrowser bool
	leakTemp          bool

	quarterSecond = time.Millisecond * 250
	workDir       string
)

const (
	separator = string(filepath.Separator)
	endGoPath = separator + "src" + separator
)

// This method exists because of a bug in the go cover tool that
// causes an infinite loop when you try to run `go test -cover`
// on a package that has an import cycle defined in one of it's
// test files. Yuck.
func testFilesImportTheirOwnPackage(packagePath string) bool {
	meta, err := packages.Load(
		&packages.Config{
			Mode:  packages.NeedName | packages.NeedImports,
			Tests: true,
		},
		packagePath,
	)
	if err != nil {
		return false
	}

	testPackageID := fmt.Sprintf("%s [%s.test]", meta[0], meta[0])

	for _, testPackage := range meta[1:] {
		if testPackage.ID != testPackageID {
			continue
		}

		for dependency := range testPackage.Imports {
			if dependency == meta[0].PkgPath {
				return true
			}
		}
		break
	}
	return false
}

func resolvePackageName(path string) string {
	pkg, err := packages.Load(
		&packages.Config{
			Mode: packages.NeedName,
		},
		path,
	)
	if err == nil {
		return pkg[0].PkgPath
	}

	nameArr := strings.Split(path, endGoPath)
	return nameArr[len(nameArr)-1]
}
