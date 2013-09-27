package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"os/exec"
	"strings"
	"time"
)

type runner struct {
	busy    bool
	ready   chan bool
	done    chan bool
	watcher *fsnotify.Watcher
	latest  string
}

func newRunner(watcher *fsnotify.Watcher) *runner {
	return &runner{false, make(chan bool), make(chan bool), watcher, ""}
}

func (self *runner) idle() {
	for {
		select {
		case event := <-self.watcher.Event:
			self.onEvent(event)
		case err := <-self.watcher.Error:
			self.onError(err)
		case <-self.done:
			self.onDone()
		case <-self.ready:
			self.busy = false
		}
	}
}

func (self *runner) onEvent(event *fsnotify.FileEvent) {
	if event == nil {
		return
	}
	fmt.Println(event)
	if strings.HasSuffix(event.Name, ".go") && !self.busy {
		self.busy = true
		go self.test(self.done)
	}
}

func (self *runner) onError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (self *runner) onDone() {
	time.AfterFunc(DELAY, func() {
		self.ready <- true
	})
}

func (self *runner) test(done chan bool) {
	fmt.Println("Running tests...")

	output, err := exec.Command("go", "test").Output() // "-json" arg
	if err != nil {
		fmt.Println(err)
	}
	self.latest = string(output)
	fmt.Println(self.latest)

	self.done <- true
}

var DELAY = 1500 * time.Millisecond
