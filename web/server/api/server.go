package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/smartystreets/goconvey/web/server/contract"
)

type HTTPServer struct {
	watcher      contract.Watcher
	executor     contract.Executor
	latest       *contract.CompleteOutput
	statusUpdate chan bool
}

func (self *HTTPServer) ReceiveUpdate(update *contract.CompleteOutput) {
	self.latest = update
}

func (self *HTTPServer) Watch(response http.ResponseWriter, request *http.Request) {

	// In case a web UI client disconnected (closed the tab),
	// the web UI will request, when it initially loads the page
	// and gets the watched directory, that the status channel
	// buffer be filled so that it can get the latest status updates
	// without missing a single beat.
	if request.URL.Query().Get("newclient") != "" {
		select {
		case self.statusUpdate <- true:
		default:
		}
	}

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
	err := self.watcher.Adjust(newRoot)
	if err != nil {
		http.Error(response, err.Error(), http.StatusNotFound)
	}
}

func (self *HTTPServer) Ignore(response http.ResponseWriter, request *http.Request) {
	paths := self.parseQueryString("paths", response, request)
	if paths != "" {
		self.watcher.Ignore(paths)
	}
}

func (self *HTTPServer) Reinstate(response http.ResponseWriter, request *http.Request) {
	paths := self.parseQueryString("paths", response, request)
	if paths != "" {
		self.watcher.Reinstate(paths)
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

func (self *HTTPServer) LongPollStatus(response http.ResponseWriter, request *http.Request) {
	select {
	case <-self.statusUpdate:
		self.Status(response, request)
	case <-time.After(1 * time.Minute): // MAJOR 'GOTCHA': This should be SHORTER than the client's timeout!
	}
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

func NewHTTPServer(watcher contract.Watcher, executor contract.Executor, status chan bool) *HTTPServer {
	self := new(HTTPServer)
	self.watcher = watcher
	self.executor = executor
	self.statusUpdate = status
	return self
}

var _ = fmt.Sprintf("Hi")
