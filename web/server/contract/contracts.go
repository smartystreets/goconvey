package contract

import (
	"github.com/smartystreets/goconvey/web/server/parser"
	"net/http"
	"path/filepath"
)

type (
	Server interface {
		ReceiveUpdate(*parser.CompleteOutput)
		Watch(writer http.ResponseWriter, request *http.Request)
		Status(writer http.ResponseWriter, request *http.Request)
		Results(writer http.ResponseWriter, request *http.Request)
		Execute(writer http.ResponseWriter, request *http.Request)
	}

	Scanner interface {
		Scan(root string) (changed bool)
	}

	Watcher interface {
		// contains the FileSystem

		Adjust(root string) error

		Deletion(folder string)
		Creation(folder string)

		Ignore(folder string) error
		Reinstate(folder string) error

		WatchedFolders() []*Package
		IsActive(folder string) bool
		IsIgnored(folder string) bool
		IsWatched(folder string) bool
	}

	Executor interface {
		// contains the executor.Parser
		// contains the Server
		// contains the Shell

		ExecuteTests([]*Package) *parser.CompleteOutput
		IsRunning() bool
	}

	FileSystem interface {
		Walk(root string, step filepath.WalkFunc)
		Exists(directory string) bool
	}

	Shell interface {
		Execute(name string, args ...string) (output string, err error)
	}
)

type Package struct {
	Active       bool
	Path         string
	Name         string
	RawOutput    string
	ParsedOutput *parser.PackageResult
}
