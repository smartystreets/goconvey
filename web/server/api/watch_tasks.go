package api

import (
	"fmt"
	"github.com/smartystreets/goconvey/web/server/contract"
	"io/ioutil"
	"log"
	"net/http"
)

type WatchRequestHandler struct {
	request  *http.Request
	response http.ResponseWriter
	watcher  contract.Watcher
}

func (self *WatchRequestHandler) ProvideCurrentRoot() {
	root := self.watcher.Root()
	self.response.Write([]byte(root))
}

func (self *WatchRequestHandler) AdjustRoot() {
	root := self.extractPayload()
	if root == "" {
		return
	}
	log.Println("Adjusting root:", root)

	err := self.watcher.Adjust(root)
	if err != nil {
		http.Error(self.response, err.Error(), http.StatusNotFound)
	}
}

func (self *WatchRequestHandler) IgnorePackage() {
	if folder := self.extractPayload(); folder != "" {
		log.Println("Ignoring:", folder)
		self.watcher.Ignore(folder)
	}
}

func (self *WatchRequestHandler) ReinstatePackage() {
	if folder := self.extractPayload(); folder != "" {
		log.Println("Reinstating:", folder)
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
	log.Printf("HTTP %d (request: %v): %s (HTTP %d)\n", http.StatusBadRequest, self.request, message)
	http.Error(self.response, message, http.StatusBadRequest)
}
func (self *WatchRequestHandler) returnBlankError() {
	message := "You must provide a non-blank path."
	log.Printf("HTTP %d (%v): %s\n", http.StatusBadRequest, self.request, message)
	http.Error(self.response, message, http.StatusBadRequest)
}

func newWatchRequestHandler(request *http.Request, response http.ResponseWriter, watcher contract.Watcher) *WatchRequestHandler {
	self := &WatchRequestHandler{}
	self.request = request
	self.response = response
	self.watcher = watcher
	return self
}
