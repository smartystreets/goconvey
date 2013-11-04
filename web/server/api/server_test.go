package api

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/contract"
	"github.com/smartystreets/goconvey/web/server/parser"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const initialRoot = "/root/gopath/src/github.com/smartystreets/project"
const nonexistentRoot = "I don't exist"

func TestHTTPServer(t *testing.T) {
	var fixture *ServerFixture

	Convey("Subject: HttpServer responds to requests appropriately", t, func() {
		fixture = newServerFixture()

		Convey("Given an update is received", func() {
			fixture.ReceiveUpdate(&parser.CompleteOutput{Revision: "asdf"})

			Convey("When the update is requested", func() {
				update, status := fixture.RequestLatest()

				Convey("The server returns it", func() {
					So(update, ShouldResemble, &parser.CompleteOutput{Revision: "asdf"})
				})

				Convey("The server returns 200", func() {
					So(status, ShouldEqual, http.StatusOK)
				})
			})
		})

		Convey("When the root watch is queried", func() {
			root, status := fixture.QueryRootWatch()

			Convey("The server returns it", func() {
				So(root, ShouldEqual, "/root/gopath/src/github.com/smartystreets/project")
			})

			Convey("The server returns HTTP 200 - OK", func() {
				So(status, ShouldEqual, http.StatusOK)
			})
		})

		Convey("When the root watch is adjusted", func() {

			Convey("But the request is empty", func() {
				status, body := fixture.AdjustRootWatch("")

				Convey("The server returns HTTP 400 - Bad Input", func() {
					So(status, ShouldEqual, http.StatusBadRequest)
				})

				Convey("The server should provide a helpful error message", func() {
					So(body, ShouldEqual, "You must provide a non-blank path to watch.")
				})

				Convey("The server should not change the existing root", func() {
					root, _ := fixture.QueryRootWatch()
					So(root, ShouldEqual, initialRoot)
				})
			})

			Convey("But the request cannot be read", func() {
				status, body := fixture.AdjustRootWatch("!!error!!")

				Convey("The server returns HTTP 400 - Bad Input", func() {
					So(status, ShouldEqual, http.StatusBadRequest)
				})

				Convey("The server should provide a helpful error message that includes the underlying read error", func() {
					So(body, ShouldEqual, fmt.Sprintf("The request body could not be read: (error: '%s')", readError))
				})

				Convey("The server should NOT change the existing root", func() {
					root, _ := fixture.QueryRootWatch()
					So(root, ShouldEqual, initialRoot)
				})
			})

			Convey("And the new root exists", func() {
				status, body := fixture.AdjustRootWatch(initialRoot + "/package")

				Convey("The server returns HTTP 200 - OK", func() {
					So(status, ShouldEqual, http.StatusOK)
				})

				Convey("The body should NOT contain any error message or content", func() {
					So(body, ShouldEqual, "")
				})

				Convey("The server informs the watcher of the new root", func() {
					root, _ := fixture.QueryRootWatch()
					So(root, ShouldEqual, initialRoot+"/package")
				})
			})

			Convey("And the new root does NOT exist", func() {
				status, body := fixture.AdjustRootWatch(nonexistentRoot)

				Convey("The server returns HTTP 404 - Not Found", func() {
					So(status, ShouldEqual, http.StatusNotFound)
				})

				Convey("The body should contain a helpful error message", func() {
					So(body, ShouldEqual, fmt.Sprintf("Directory does not exist: '%s'", nonexistentRoot))
				})

				Convey("The server should not change the existing root", func() {
					root, _ := fixture.QueryRootWatch()
					So(root, ShouldEqual, initialRoot)
				})
			})
		})

		Convey("When a packge is ignored", func() {
			Convey("But the request is malformed", func() {
				Convey("The server returns HTTP 400 - Bad Input", nil)
			})

			Convey("And the request is well formed", func() {
				Convey("The server informs the watcher", nil)
				Convey("The server returns HTTP 200 - OK", nil)
			})
		})

		Convey("When a package is reinstated", func() {
			Convey("But the request is malformed", func() {
				Convey("The server returns 400", nil)
			})

			Convey("And the request is well formed", func() {
				Convey("The server informs the watcher", nil)
				Convey("The server returns HTTP 200 - OK", nil)
			})
		})

		Convey("When the status is requested", func() {
			Convey("The server returns it", nil)
			Convey("The server returns HTTP 200 - OK", nil)
		})

		Convey("When a manual execution of the test packages is requested", func() {
			Convey("The server invokes the executor using the watcher's listing", nil)
			Convey("The server returns HTTP 200 - OK", nil)
		})
	})
}

/********* Server Fixture *********/

type ServerFixture struct {
	server         *HTTPServer
	watcher        *FakeWatcher
	queriedRoot    string
	lastHTTPStatus int
}

func (self *ServerFixture) ReceiveUpdate(update *parser.CompleteOutput) {
	self.server.ReceiveUpdate(update)
}

func (self *ServerFixture) RequestLatest() (*parser.CompleteOutput, int) {
	request, _ := http.NewRequest("GET", "http://localhost:8080/results", nil)
	response := httptest.NewRecorder()

	self.server.Results(response, request)

	decoder := json.NewDecoder(strings.NewReader(response.Body.String()))
	update := &parser.CompleteOutput{}
	decoder.Decode(update)
	return update, response.Code
}

func (self *ServerFixture) QueryRootWatch() (string, int) {
	request, _ := http.NewRequest("GET", "http://localhost:8080/watch", nil)
	response := httptest.NewRecorder()

	self.server.Watch(response, request)

	return strings.TrimSpace(response.Body.String()), response.Code
}

func (self *ServerFixture) AdjustRootWatch(newRoot string) (status int, body string) {
	var reader io.Reader = strings.NewReader(newRoot)
	if newRoot == "!!error!!" {
		reader = &ErrorReadCloser{}
	}
	request, _ := http.NewRequest("PUT", "http://localhost:8080/watch", reader)
	response := httptest.NewRecorder()

	self.server.Watch(response, request)

	status, body = response.Code, strings.TrimSpace(response.Body.String())
	return
}

func newServerFixture() *ServerFixture {
	self := &ServerFixture{}
	self.watcher = newFakeWatcher()
	self.watcher.SetRootWatch(initialRoot)
	self.server = NewHTTPServer(self.watcher)
	return self
}

/********* Fake Watcher *********/

type FakeWatcher struct {
	root string
}

func (self *FakeWatcher) SetRootWatch(root string) {
	self.root = root
}

func (self *FakeWatcher) WatchedFolders() []*contract.Package {
	return []*contract.Package{&contract.Package{Path: self.root}}
}

func (self *FakeWatcher) Adjust(root string) error {
	if root == nonexistentRoot {
		return errors.New(fmt.Sprintf("Directory does not exist: '%s'", root))
	}
	self.root = root
	return nil
}
func (self *FakeWatcher) Deletion(folder string)       { panic("NOT SUPPORTED") }
func (self *FakeWatcher) Creation(folder string)       { panic("NOT SUPPORTED") }
func (self *FakeWatcher) Ignore(folder string)         { panic("NOT SUPPORTED") }
func (self *FakeWatcher) Reinstate(folder string)      { panic("NOT SUPPORTED") }
func (self *FakeWatcher) IsWatched(folder string) bool { panic("NOT SUPPORTED") }
func (self *FakeWatcher) IsIgnored(folder string) bool { panic("NOT SUPPORTED") }

func newFakeWatcher() *FakeWatcher {
	return &FakeWatcher{}
}

/********* Error Read Closer *********/

const readError = "Not sure why the request body would ever throw an error when read, " +
	"but this ensures that we know how to handle it..."

type ErrorReadCloser struct{}

func (self *ErrorReadCloser) Read(buffer []byte) (int, error) {
	return 0, errors.New(readError)
}

func (self *ErrorReadCloser) Close() error {
	return nil
}

var _ = fmt.Sprintf("hi")
