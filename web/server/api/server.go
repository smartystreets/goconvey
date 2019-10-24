package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/smartystreets/goconvey/web/server/contract"
	"github.com/smartystreets/goconvey/web/server/messaging"
)

type HTTPServer struct {
	watcher     chan messaging.WatcherCommand
	executor    contract.Executor
	latest      *contract.CompleteOutput
	currentRoot string
	longpoll    chan chan string
	paused      bool
}

func (h *HTTPServer) ReceiveUpdate(root string, update *contract.CompleteOutput) {
	h.currentRoot = root
	h.latest = update
}

func (h *HTTPServer) Watch(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		h.adjustRoot(response, request)
	} else if request.Method == "GET" {
		response.Write([]byte(h.currentRoot))
	}
}

func (h *HTTPServer) adjustRoot(response http.ResponseWriter, request *http.Request) {
	newRoot := h.parseQueryString("root", response, request)
	if newRoot == "" {
		return
	}
	info, err := os.Stat(newRoot) // TODO: how to unit test?
	if !info.IsDir() || err != nil {
		http.Error(response, err.Error(), http.StatusNotFound)
		return
	}

	h.watcher <- messaging.WatcherCommand{
		Instruction: messaging.WatcherAdjustRoot,
		Details:     newRoot,
	}
}

func (h *HTTPServer) Ignore(response http.ResponseWriter, request *http.Request) {
	paths := h.parseQueryString("paths", response, request)
	if paths != "" {
		h.watcher <- messaging.WatcherCommand{
			Instruction: messaging.WatcherIgnore,
			Details:     paths,
		}
	}
}

func (h *HTTPServer) Reinstate(response http.ResponseWriter, request *http.Request) {
	paths := h.parseQueryString("paths", response, request)
	if paths != "" {
		h.watcher <- messaging.WatcherCommand{
			Instruction: messaging.WatcherReinstate,
			Details:     paths,
		}
	}
}

func (h *HTTPServer) parseQueryString(key string, response http.ResponseWriter, request *http.Request) string {
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

func (h *HTTPServer) Status(response http.ResponseWriter, request *http.Request) {
	status := h.executor.Status()
	response.Write([]byte(status))
}

func (h *HTTPServer) LongPollStatus(response http.ResponseWriter, request *http.Request) {
	if h.executor.ClearStatusFlag() {
		response.Write([]byte(h.executor.Status()))
		return
	}

	timeout, err := strconv.Atoi(request.URL.Query().Get("timeout"))
	if err != nil || timeout > 180000 || timeout < 0 {
		timeout = 60000 // default timeout is 60 seconds
	}

	myReqChan := make(chan string)

	select {
	case h.longpoll <- myReqChan: // this case means the executor's status is changing
	case <-time.After(time.Duration(timeout) * time.Millisecond): // this case means the executor hasn't changed status
		return
	}

	out := <-myReqChan

	if out != "" { // TODO: Why is this check necessary? Sometimes it writes empty string...
		response.Write([]byte(out))
	}
}

func (h *HTTPServer) Results(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	response.Header().Set("Pragma", "no-cache")
	response.Header().Set("Expires", "0")
	if h.latest != nil {
		h.latest.Paused = h.paused
	}
	stuff, _ := json.Marshal(h.latest)
	response.Write(stuff)
}

func (h *HTTPServer) Execute(response http.ResponseWriter, request *http.Request) {
	go h.execute()
}

func (h *HTTPServer) execute() {
	h.watcher <- messaging.WatcherCommand{Instruction: messaging.WatcherExecute}
}

func (h *HTTPServer) TogglePause(response http.ResponseWriter, request *http.Request) {
	instruction := messaging.WatcherPause
	if h.paused {
		instruction = messaging.WatcherResume
	}

	h.watcher <- messaging.WatcherCommand{Instruction: instruction}
	h.paused = !h.paused

	fmt.Fprint(response, h.paused) // we could write out whatever helps keep the UI honest...
}

func NewHTTPServer(
	root string,
	watcher chan messaging.WatcherCommand,
	executor contract.Executor,
	status chan chan string) *HTTPServer {

	self := new(HTTPServer)
	self.currentRoot = root
	self.watcher = watcher
	self.executor = executor
	self.longpoll = status
	return self
}
