package api

import (
	"encoding/json"
	"fmt"
	"github.com/smartystreets/goconvey/web/server/contract"
	"log"
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
	watch := newWatchRequestHandler(request, response, self.watcher)

	switch request.Method {
	case "PUT":
		watch.AdjustRoot()
	case "DELETE":
		watch.IgnorePackage()
	case "POST":
		watch.ReinstatePackage()
	case "GET":
		watch.ProvideCurrentRoot()
	}
}

func (self *HTTPServer) Status(response http.ResponseWriter, request *http.Request) {
	status := self.executor.Status()
	log.Println("Status requested:", status)
	response.Write([]byte(status))
}

func (self *HTTPServer) Results(response http.ResponseWriter, request *http.Request) {
	// TODO: include these headers (under test)
	// response.Header().Set("Content-Type", "application/json")
	// response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	// response.Header().Set("Pragma", "no-cache")
	// response.Header().Set("Expires", "0")
	stuff, _ := json.Marshal(self.latest)
	response.Write(stuff)
}

func (self *HTTPServer) Execute(response http.ResponseWriter, request *http.Request) {
	go func() {
		self.latest = self.executor.ExecuteTests(self.watcher.WatchedFolders())
	}()
}

func NewHTTPServer(watcher contract.Watcher, executor contract.Executor) *HTTPServer {
	self := &HTTPServer{}
	self.watcher = watcher
	self.executor = executor
	return self
}

var _ = fmt.Sprintf("Hi")
