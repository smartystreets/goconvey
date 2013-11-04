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
	root := self.extractPayload()
	if root == "" {
		return
	}

	err := self.watcher.Adjust(root)
	if err != nil {
		http.Error(self.response, err.Error(), http.StatusNotFound)
	}
}

func (self *WatchRequestHandler) IgnorePackage() {
	if folder := self.extractPayload(); folder != "" {
		self.watcher.Ignore(folder)
	}
}

func (self *WatchRequestHandler) ReinstatePackage() {
	if folder := self.extractPayload(); folder != "" {
		self.watcher.Reinstate(folder)
	}
}

func (self *WatchRequestHandler) extractPayload() (payload string) {
	rawBody, err := ioutil.ReadAll(self.request.Body)

	if err != nil {
		self.returnReadError(err)
	} else if len(rawBody) == 0 {
		self.returnBlankError()
	} else {
		payload = string(rawBody)
	}
	return
}
func (self *WatchRequestHandler) returnReadError(err error) {
	message := fmt.Sprintf("The request body could not be read: (error: '%s')", err.Error())
	http.Error(self.response, message, http.StatusBadRequest)
}
func (self *WatchRequestHandler) returnBlankError() {
	http.Error(self.response, "You must provide a non-blank path.", http.StatusBadRequest)
}

func newWatchRequestHandler(request *http.Request, response http.ResponseWriter, watcher contract.Watcher) *WatchRequestHandler {
	self := &WatchRequestHandler{}
	self.request = request
	self.response = response
	self.watcher = watcher
	return self
}
