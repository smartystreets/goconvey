package watcher

import (
	"fmt"
	"github.com/smartystreets/goconvey/web/server/contract"
	"path/filepath"
	"strings"
)

type goPath struct {
	files   contract.FileSystem
	shell   contract.Shell
	ambient []string
}

func (self *goPath) Set(gopath string) {
	for _, workspace := range strings.Split(gopath, delimiter) {
		if self.isAmbientWorkspace(workspace) {
			self.shell.Setenv("GOPATH", strings.Join(self.ambient, delimiter))
			return
		}
	}
	self.shell.Setenv("GOPATH", gopath)
}
func (self *goPath) isAmbientWorkspace(workspace string) bool {
	for _, x := range self.ambient {
		if workspace == x {
			return true
		}
	}
	return false
}

func (self *goPath) ResolvePackageName(folder string) string {
	for _, workspace := range strings.Split(self.current(), delimiter) {
		if strings.HasPrefix(folder, workspace) {
			prefix := filepath.Join(workspace, "src") + separator
			return folder[len(prefix):]
		}
	}

	panic(fmt.Sprintln(resolutionError, self.current()))
}

func (self *goPath) current() string {
	return self.shell.Getenv("GOPATH")
}

func newGoPath(files contract.FileSystem, shell contract.Shell) *goPath {
	self := &goPath{}
	self.files = files
	self.shell = shell
	self.ambient = strings.Split(self.current(), delimiter)
	return self
}

const delimiter = string(filepath.ListSeparator)
const separator = string(filepath.Separator)
const resolutionError = "Package cannot be resolved as it is outside of any workspaces listed in the current $GOPATH:"
