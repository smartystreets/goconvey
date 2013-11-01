package api

import (
	"encoding/json"
	"github.com/smartystreets/goconvey/web/server/parser"
	"net/http"
)

type HTTPServer struct {
	latest *parser.CompleteOutput
}

func (self *HTTPServer) ReceiveUpdate(update *parser.CompleteOutput) {
	self.latest = update
}

// GET (query root) vs PUT (adjust root) vs POST (reinstate) vs DELETE (ignore)
func (self *HTTPServer) Watch(writer http.ResponseWriter, request *http.Request) {}

func (self *HTTPServer) Status(writer http.ResponseWriter, request *http.Request) {}

func (self *HTTPServer) Results(writer http.ResponseWriter, request *http.Request) {
	stuff, _ := json.Marshal(self.latest)
	writer.Write(stuff)
}

func (self *HTTPServer) Execute(writer http.ResponseWriter, request *http.Request) {}

func NewHTTPServer() *HTTPServer {
	self := &HTTPServer{}
	return self
}
