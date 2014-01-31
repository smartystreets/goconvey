package contract

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMonitor(t *testing.T) {
	var fixture *MonitorFixture

	Convey("Subject: Monitor", t, func() {
		fixture = newMonitorFixture()

		Convey("When the file system has changed", func() {
			fixture.scanner.Modify("/root")

			Convey("As a result of scanning", func() {
				fixture.Scan()

				Convey("The watched packages should be executed and the results should be passed to the server", func() {
					So(fixture.server.latest, ShouldResemble, &CompleteOutput{Packages: []*PackageResult{NewPackageResult("1"), NewPackageResult("2")}})
				})
			})
		})

		Convey("When the file system has remained stagnant", func() {
			fixture.scanner.Reset("/root")

			Convey("As a result of scanning", func() {
				fixture.Scan()

				Convey("The process should take a nap", func() {
					So(fixture.nap, ShouldBeTrue)
				})

				Convey("The server should not receive any update", func() {
					So(fixture.server.latest, ShouldBeNil)
				})
			})
		})
	})
}

/******** MonitorFixture ********/

type MonitorFixture struct {
	monitor  *Monitor
	server   *FakeServer
	watcher  *FakeWatcher
	scanner  *FakeScanner
	executor *FakeExecutor
	nap      bool
}

func (self *MonitorFixture) Scan() {
	self.monitor.Scan()
}

func (self *MonitorFixture) sleep() {
	self.nap = true
}

func newMonitorFixture() *MonitorFixture {
	self := &MonitorFixture{}
	self.server = newFakeServer()
	self.watcher = newFakeWatcher()
	self.scanner = newFakeScanner()
	self.executor = newFakeExecutor()
	self.monitor = NewMonitor(self.scanner, self.watcher, self.executor, self.server, self.sleep)
	return self
}

/******** FakeServer ********/

type FakeServer struct {
	latest *CompleteOutput
}

func (self *FakeServer) ReceiveUpdate(update *CompleteOutput) {
	self.latest = update
}
func (self *FakeServer) Watch(http.ResponseWriter, *http.Request)          { panic("NOT SUPPORTED") }
func (self *FakeServer) Ignore(http.ResponseWriter, *http.Request)         { panic("NOT SUPPORTED") }
func (self *FakeServer) Reinstate(http.ResponseWriter, *http.Request)      { panic("NOT SUPPORTED") }
func (self *FakeServer) Status(http.ResponseWriter, *http.Request)         { panic("NOT SUPPORTED") }
func (self *FakeServer) LongPollStatus(http.ResponseWriter, *http.Request) { panic("NOT SUPPORTED") }
func (self *FakeServer) Results(http.ResponseWriter, *http.Request)        { panic("NOT SUPPORTED") }
func (self *FakeServer) Execute(http.ResponseWriter, *http.Request)        { panic("NOT SUPPORTED") }

func newFakeServer() *FakeServer {
	self := &FakeServer{}
	return self
}

/******** FakeWatcher ********/

type FakeWatcher struct{}

func (self *FakeWatcher) WatchedFolders() []*Package {
	return []*Package{
		&Package{Path: "/root", Result: NewPackageResult("1")},
		&Package{Path: "/root/nested", Result: NewPackageResult("2")},
	}
}

func (self *FakeWatcher) Root() string {
	return self.WatchedFolders()[0].Path
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

func (self *FakeScanner) Modify(path string) {
	self.dirty = true
}

func (self *FakeScanner) Reset(path string) {
	self.dirty = false
}

func (self *FakeScanner) Scan() (changed bool) {
	return self.dirty
}

func newFakeScanner() *FakeScanner {
	self := &FakeScanner{}
	return self
}

/******** FakeExecutor ********/

type FakeExecutor struct{}

func (self *FakeExecutor) ExecuteTests(packages []*Package) *CompleteOutput {
	complete := &CompleteOutput{}
	complete.Packages = []*PackageResult{}
	for _, p := range packages {
		complete.Packages = append(complete.Packages, p.Result)
	}
	return complete
}
func (self *FakeExecutor) Status() string { panic("NOT SUPPORTED") }

func newFakeExecutor() *FakeExecutor {
	self := &FakeExecutor{}
	return self
}
