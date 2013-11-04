package api

import (
	"encoding/json"
	"fmt"
	"github.com/smartystreets/goconvey/web/server/contract"
	"github.com/smartystreets/goconvey/web/server/parser"
	"net/http"
)

type HTTPServer struct {
	watcher contract.Watcher
	latest  *parser.CompleteOutput
}

func (self *HTTPServer) ReceiveUpdate(update *parser.CompleteOutput) {
	self.latest = update
}

// GET (query root) vs PUT (adjust root) vs POST (reinstate) vs DELETE (ignore)
func (self *HTTPServer) Watch(response http.ResponseWriter, request *http.Request) {
	handler := newWatchRequestHandler(request, response, self.watcher)

	switch request.Method {
	case "PUT":
		handler.AdjustRoot()
	default: // GET
		handler.ProvideCurrentRoot()
	}
}

func (self *HTTPServer) Status(response http.ResponseWriter, request *http.Request) {}

func (self *HTTPServer) Results(response http.ResponseWriter, request *http.Request) {
	stuff, _ := json.Marshal(self.latest)
	response.Write(stuff)
}

func (self *HTTPServer) Execute(response http.ResponseWriter, request *http.Request) {}

func NewHTTPServer(watcher contract.Watcher) *HTTPServer {
	self := &HTTPServer{}
	self.watcher = watcher
	return self
}

var _ = fmt.Sprintf("Hi")
