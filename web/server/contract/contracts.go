package contract

import (
	"github.com/smartystreets/goconvey/web/server/parser"
	"net/http"
	"path/filepath"
)

type (
	Server interface {
		ReceiveUpdate(*parser.CompleteOutput)                      // internal
		Watch(writer http.ResponseWriter, request *http.Request)   // GET vs POST
		Status(writer http.ResponseWriter, request *http.Request)  // GET
		Results(writer http.ResponseWriter, request *http.Request) // GET
		Execute(writer http.ResponseWriter, request *http.Request) // POST
	}

	Executor interface {
		// contains the executor.Parser
		// contains the Server
		// contains the Shell

		ExecuteTests([]*Package) *parser.CompleteOutput
		IsRunning() bool
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
