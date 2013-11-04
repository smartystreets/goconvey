package contract

import (
	"github.com/smartystreets/goconvey/web/server/parser"
	"net/http"
	"path/filepath"
	"strings"
)

/*

TODO: Watcher needs to return brand-new references of *Package in the watch list.
TODO: Concrete Executor needs to accept []*Package, not []string.
TODO: Tester needs to accept []*Package, not []string. (and filter on .Active accorrdingly)
TODO: Coorinator needs to accpt []*Packge, not []string. (and filter on .Active accorrdingly)
TODO: Parser needs to accept []*Package, not (packageName, output string) (and filter on .Active accordingly)

*/

type (
	Server interface {
		ReceiveUpdate(*parser.CompleteOutput)
		Watch(writer http.ResponseWriter, request *http.Request)
		Status(writer http.ResponseWriter, request *http.Request)
		Results(writer http.ResponseWriter, request *http.Request)
		Execute(writer http.ResponseWriter, request *http.Request)
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

		Ignore(folder string)
		Reinstate(folder string)

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
	Active bool
	Path   string
	Name   string
}

func ResolvePackageName(path string) string {
	index := strings.Index(path, endGoPath)
	if index < 0 {
		return path
	}
	packageBeginning := index + len(endGoPath)
	name := path[packageBeginning:]
	return name
}

const (
	separator = string(filepath.Separator)
	endGoPath = separator + "src" + separator
)
