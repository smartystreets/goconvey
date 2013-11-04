package contract

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/parser"
	"net/http"
	"testing"
)

func TestMonitor(t *testing.T) {
	var fixture *MonitorFixture

	Convey("Subject: Monitor", t, func() {
		fixture = newMonitorFixture()

		Convey("When the file system has changed", func() {
			fixture.scanner.dirty = true

			Convey("As a result of scanning", func() {
				fixture.Scan()

				Convey("The actively watched tests should be executed and the results should be passed to the server", nil)
			})
		})

		Convey("When the file system has remained stagnant", func() {
			Convey("As a result of scanning", func() {
				Convey("The process should take a nap", nil)
			})
		})
	})
}

type MonitorFixture struct {
	monitor  *Monitor
	server   *FakeServer
	watcher  *FakeWatcher
	scanner  *FakeScanner
	executor *FakeExecutor
}

func (self *MonitorFixture) Scan() {

}

func newMonitorFixture() *MonitorFixture {
	self := &MonitorFixture{}
	self.server = newFakeServer()
	self.watcher = newFakeWatcher()
	self.scanner = newFakeScanner()
	self.executor = newFakeExecutor()
	self.monitor = NewMonitor(self.scanner, self.watcher, self.executor, self.server)
	return self
}

/******** FakeServer ********/

type FakeServer struct {
}

func (self *FakeServer) ReceiveUpdate(*parser.CompleteOutput) {
	panic("NOT SUPPORTED")
}
func (self *FakeServer) Watch(writer http.ResponseWriter, request *http.Request) {
	panic("NOT SUPPORTED")
}
func (self *FakeServer) Status(writer http.ResponseWriter, request *http.Request) {
	panic("NOT SUPPORTED")
}
func (self *FakeServer) Results(writer http.ResponseWriter, request *http.Request) {
	panic("NOT SUPPORTED")
}
func (self *FakeServer) Execute(writer http.ResponseWriter, request *http.Request) {
	panic("NOT SUPPORTED")
}

func newFakeServer() *FakeServer {
	self := &FakeServer{}
	return self
}

/******** FakeWatcher ********/

type FakeWatcher struct {
}

func (self *FakeWatcher) WatchedFolders() []*Package {
	return []*Package{
		&Package{true, "/path/1", "1"},
		&Package{false, "/path/2", "2"},
		&Package{true, "/path/3", "3"},
	}
}

func (self *FakeWatcher) Adjust(root string) error     { panic("NOT SUPPORTED") }
func (self *FakeWatcher) Deletion(folder string)       { panic("NOT SUPPORTED") }
func (self *FakeWatcher) Creation(folder string)       { panic("NOT SUPPORTED") }
func (self *FakeWatcher) Ignore(folder string)         { panic("NOT SUPPORTED") }
func (self *FakeWatcher) Reinstate(folder string)      { panic("NOT SUPPORTED") }
func (self *FakeWatcher) IsWatched(folder string) bool { panic("NOT SUPPORTED") }
func (self *FakeWatcher) IsIgnored(folder string) bool { panic("NOT SUPPORTED") }

func newFakeWatcher() *FakeWatcher {
	self := &FakeWatcher{}
	return self
}

/******** FakeScanner ********/

type FakeScanner struct {
	dirty bool
}

func (self *FakeScanner) Scan(root string) (changed bool) { panic("NOT SUPPORTED") }

func newFakeScanner() *FakeScanner {
	self := &FakeScanner{}
	return self
}

/******** FakeExecutor ********/

type FakeExecutor struct {
}

func (self *FakeExecutor) ExecuteTests([]*Package) *parser.CompleteOutput { panic("NOT SUPPORTED") }
func (self *FakeExecutor) Status() string                                 { panic("NOT SUPPORTED") }

func newFakeExecutor() *FakeExecutor {
	self := &FakeExecutor{}
	return self
}
