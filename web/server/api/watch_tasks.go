package api

import (
	"fmt"
	"github.com/smartystreets/goconvey/web/server/contract"
	"io/ioutil"
	"net/http"
)

type WatchRequestHandler struct {
	request  *http.Request
	response http.ResponseWriter
	watcher  contract.Watcher
}

func (self *WatchRequestHandler) ProvideCurrentRoot() {
	watched := self.watcher.WatchedFolders() // TODO: what if len(watched) == 0? (can that even happen?)
	self.response.Write([]byte(watched[0].Path))
}

func (self *WatchRequestHandler) AdjustRoot() {
	rawBody, err := ioutil.ReadAll(self.request.Body)
	if err != nil {
		self.returnReadError(err)
	} else if len(rawBody) == 0 {
		self.returnBlankError()
	} else {
		self.adjust(rawBody)
	}
}
func (self *WatchRequestHandler) returnReadError(err error) {
	message := fmt.Sprintf("The request body could not be read: (error: '%s')", err.Error())
	http.Error(self.response, message, http.StatusBadRequest)
}
func (self *WatchRequestHandler) returnBlankError() {
	http.Error(self.response, "You must provide a non-blank path to watch.", http.StatusBadRequest)
}
func (self *WatchRequestHandler) adjust(rawBody []byte) {
	err := self.watcher.Adjust(string(rawBody))
	if err != nil {
		http.Error(self.response, err.Error(), http.StatusNotFound)
	}
}

func newWatchRequestHandler(request *http.Request, response http.ResponseWriter, watcher contract.Watcher) *WatchRequestHandler {
	self := &WatchRequestHandler{}
	self.request = request
	self.response = response
	self.watcher = watcher
	return self
}
