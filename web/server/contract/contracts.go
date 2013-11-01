package contract

import (
	"github.com/smartystreets/goconvey/web/server/parser"
	"net/http"
	"path/filepath"
)

type (
	Monitor interface {
		// contains the Server
		// contains the Watcher
		// contains the Scanner
		// contains the Executor
		Tick()   // one round of scanning and test execution
		Engage() // infinite for loop, calls Tick() between time.Sleep() (when no tests were run)
	}

	Server interface {
		ReceiveUpdate(*parser.CompleteOutput)
		Watch(writer http.ResponseWriter, request *http.Request)   // GET (query root) vs PUT (adjust root) vs POST (reinstate) vs DELETE (ignore)
		Status(writer http.ResponseWriter, request *http.Request)  // GET
		Results(writer http.ResponseWriter, request *http.Request) // GET
		Execute(writer http.ResponseWriter, request *http.Request) // POST
	}

	Executor interface {
		ExecuteTests([]*Package) *parser.CompleteOutput
		Status() string
	}

	Scanner interface {
		Scan(root string) (changed bool)
	}

	Watcher interface {
		Adjust(root string) error

		Deletion(folder string)
		Creation(folder string)

		Ignore(folder string) error
		Reinstate(folder string) error

		WatchedFolders() []*Package
		IsWatched(folder string) bool
		IsIgnored(folder string) bool
	}

	FileSystem interface {
		Walk(root string, step filepath.WalkFunc)
		Exists(directory string) bool
	}

	Shell interface {
		Execute(name string, args ...string) (output string, err error)
		Getenv(key string) string
		Setenv(key, value string) error
	}
)

type Package struct {
	Active       bool
	Path         string
	Name         string
	RawOutput    string
	ParsedOutput *parser.PackageResult
}
