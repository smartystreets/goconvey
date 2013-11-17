package api

import (
	"encoding/json"
	"fmt"
	"github.com/smartystreets/goconvey/web/server/contract"
	"net/http"
)

type HTTPServer struct {
	watcher  contract.Watcher
	executor contract.Executor
	latest   *contract.CompleteOutput
}

func (self *HTTPServer) ReceiveUpdate(update *contract.CompleteOutput) {
	self.latest = update
}

func (self *HTTPServer) Watch(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		self.adjustRoot(response, request)
	} else if request.Method == "GET" {
		response.Write([]byte(self.watcher.Root()))
	}
}

func (self *HTTPServer) adjustRoot(response http.ResponseWriter, request *http.Request) {
	newRoot := self.parseQueryString("root", response, request)
	if newRoot == "" {
		return
	}
	// TODO: make sure the newRoot is within the ambient $GOPATH. We probably need a shell now.
	err := self.watcher.Adjust(newRoot)
	if err != nil {
		http.Error(response, err.Error(), http.StatusNotFound)
	}
}

func (self *HTTPServer) Ignore(response http.ResponseWriter, request *http.Request) {
	path := self.parseQueryString("path", response, request)
	if path != "" {
		self.watcher.Ignore(path)
	}
}

func (self *HTTPServer) Reinstate(response http.ResponseWriter, request *http.Request) {
	path := self.parseQueryString("path", response, request)
	if path != "" {
		self.watcher.Reinstate(path)
	}
}

func (self *HTTPServer) parseQueryString(key string, response http.ResponseWriter, request *http.Request) string {
	value := request.URL.Query()[key]

	if len(value) == 0 {
		http.Error(response, fmt.Sprintf("No '%s' query string parameter included!", key), http.StatusBadRequest)
		return ""
	}

	path := value[0]
	if path == "" {
		http.Error(response, "You must provide a non-blank path.", http.StatusBadRequest)
	}
	return path
}

func (self *HTTPServer) Status(response http.ResponseWriter, request *http.Request) {
	status := self.executor.Status()
	response.Write([]byte(status))
}

func (self *HTTPServer) Results(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	response.Header().Set("Pragma", "no-cache")
	response.Header().Set("Expires", "0")
	stuff, _ := json.Marshal(self.latest)
	response.Write(stuff)
}

func (self *HTTPServer) Execute(response http.ResponseWriter, request *http.Request) {
	go self.execute()
}

func (self *HTTPServer) execute() {
	self.latest = self.executor.ExecuteTests(self.watcher.WatchedFolders())
}

func NewHTTPServer(watcher contract.Watcher, executor contract.Executor) *HTTPServer {
	self := &HTTPServer{}
	self.watcher = watcher
	self.executor = executor
	return self
}

var _ = fmt.Sprintf("Hi")
