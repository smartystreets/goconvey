package api

import (
	"github.com/smartystreets/goconvey/web/server/parser"
	"net/http"
)

type HTTPServer struct {
	// contains the Watcher
	// contains the Executor
	// contains the FileSystem
}

func (self *HTTPServer) ReceiveUpdate(*parser.CompleteOutput)                      {}
func (self *HTTPServer) Watch(writer http.ResponseWriter, request *http.Request)   {}
func (self *HTTPServer) Status(writer http.ResponseWriter, request *http.Request)  {}
func (self *HTTPServer) Results(writer http.ResponseWriter, request *http.Request) {}
func (self *HTTPServer) Execute(writer http.ResponseWriter, request *http.Request) {}

func NewHTTPServer() *HTTPServer {
	self := &HTTPServer{}
	return self
}
