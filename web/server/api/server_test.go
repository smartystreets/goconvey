package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/contract"
)

const initialRoot = "/root/gopath/src/github.com/smartystreets/project"
const nonexistentRoot = "I don't exist"
const unreadableContent = "!!error!!"

func TestHTTPServer(t *testing.T) {
	var fixture *ServerFixture

	Convey("Subject: HttpServer responds to requests appropriately", t, func() {
		fixture = newServerFixture()

		Convey("Given an update is received", func() {
			fixture.ReceiveUpdate(&contract.CompleteOutput{Revision: "asdf"})

			Convey("When the update is requested", func() {
				update, response := fixture.RequestLatest()

				Convey("The server returns it", func() {
					So(update, ShouldResemble, &contract.CompleteOutput{Revision: "asdf"})
				})

				Convey("The server returns 200", func() {
					So(response.Code, ShouldEqual, http.StatusOK)
				})

				Convey("The server should include important cache-related headers", func() {
					So(len(response.HeaderMap), ShouldEqual, 4)
					So(response.HeaderMap["Content-Type"][0], ShouldEqual, "application/json")
					So(response.HeaderMap["Cache-Control"][0], ShouldEqual, "no-cache, no-store, must-revalidate")
					So(response.HeaderMap["Pragma"][0], ShouldEqual, "no-cache")
					So(response.HeaderMap["Expires"][0], ShouldEqual, "0")
				})
			})
		})

		lpRequests := 6 // Number of long-poll requests and status updates to try

		Convey("Given a long-polling request for a status update, when initially idle", func() {
			fixture.executor.status = "idle"
			lpDone := make(chan string)

			go func() {
				for i := 0; i < lpRequests; i++ {
					request, _ := http.NewRequest("GET", "http://localhost:8080/status/poll", nil)
					response := httptest.NewRecorder()
					fixture.server.LongPollStatus(response, request)
					_, newStatus := response.Code, strings.TrimSpace(response.Body.String())
					lpDone <- newStatus
				}
			}()

			Convey("When the status is changed by the executor, the response should immediately reflect that", func() {
				for i := 0; i < lpRequests; i++ {
					expectedStatus := statusRotation(i, lpRequests)
					fixture.SetExecutorStatus(expectedStatus)

					select {
					case actualStatus := <-lpDone:
						So(actualStatus, ShouldEqual, expectedStatus)
					case <-time.After(500 * time.Millisecond):
						So("TIMEOUT", ShouldEqual, expectedStatus)
					}

					/*Convey("The response should be sent immediately with the correct status", func() {
						// TODO: When issue #81 is fixed and Conveys can be nested
						// inside loops again, let's put the select {...} stuff
						// from the lines just above and put it inside its own convey
						// to actually make the assertions. Also see executor_test.go
						// for a similar problem.
					})*/
				}
			})
		})

		Convey("When the root watch is queried", func() {
			root, status := fixture.QueryRootWatch(false)

			Convey("The server returns it", func() {
				So(root, ShouldEqual, initialRoot)
			})

			Convey("The server returns HTTP 200 - OK", func() {
				So(status, ShouldEqual, http.StatusOK)
			})
		})

		Convey("When the root watch is queried as a new client", func() {
			fixture.QueryRootWatch(true)

			Convey("The status channel buffer should have a true value", func() {
				select {
				case val := <-fixture.server.statusNotif:
					So(val, ShouldBeTrue)
				default:
					So(false, ShouldBeTrue)
				}
			})
		})

		Convey("When the root watch is adjusted", func() {

			Convey("But the request has no root parameter", func() {
				status, body := fixture.AdjustRootWatchMalformed()

				Convey("The server returns HTTP 400 - Bad Input", func() {
					So(status, ShouldEqual, http.StatusBadRequest)
				})

				Convey("The body should contain a helpful error message", func() {
					So(body, ShouldEqual, "No 'root' query string parameter included!")
				})

				Convey("The server should not change the existing root", func() {
					root, _ := fixture.QueryRootWatch(false)
					So(root, ShouldEqual, initialRoot)
				})
			})

			Convey("But the root parameter is empty", func() {
				status, body := fixture.AdjustRootWatch("")

				Convey("The server returns HTTP 400 - Bad Input", func() {
					So(status, ShouldEqual, http.StatusBadRequest)
				})

				Convey("The server should provide a helpful error message", func() {
					So(body, ShouldEqual, "You must provide a non-blank path.")
				})

				Convey("The server should not change the existing root", func() {
					root, _ := fixture.QueryRootWatch(false)
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
					root, _ := fixture.QueryRootWatch(false)
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
					root, _ := fixture.QueryRootWatch(false)
					So(root, ShouldEqual, initialRoot)
				})
			})
		})

		Convey("When a packge is ignored", func() {

			Convey("But the request has no path parameter", func() {
				status, body := fixture.IgnoreMalformed()

				Convey("The server returns HTTP 400 - Bad Input", func() {
					So(status, ShouldEqual, http.StatusBadRequest)
				})

				Convey("The body should contain a helpful error message", func() {
					So(body, ShouldEqual, "No 'path' query string parameter included!")
				})

				Convey("The server should not ignore anything", func() {
					So(fixture.watcher.ignored, ShouldEqual, "")
				})
			})

			Convey("But the request is blank", func() {
				status, body := fixture.Ignore("")

				Convey("The server returns HTTP 400 - Bad Input", func() {
					So(status, ShouldEqual, http.StatusBadRequest)
				})

				Convey("The body should contain a helpful error message", func() {
					So(body, ShouldEqual, "You must provide a non-blank path.")
				})
			})

			Convey("And the request is well formed", func() {
				status, _ := fixture.Ignore(initialRoot)

				Convey("The server informs the watcher", func() {
					So(fixture.watcher.ignored, ShouldEqual, initialRoot)
				})
				Convey("The server returns HTTP 200 - OK", func() {
					So(status, ShouldEqual, http.StatusOK)
				})
			})
		})

		Convey("When a package is reinstated", func() {
			Convey("But the request has no path parameter", func() {
				status, body := fixture.ReinstateMalformed()

				Convey("The server returns HTTP 400 - Bad Input", func() {
					So(status, ShouldEqual, http.StatusBadRequest)
				})

				Convey("The body should contain a helpful error message", func() {
					So(body, ShouldEqual, "No 'path' query string parameter included!")
				})

				Convey("The server should not ignore anything", func() {
					So(fixture.watcher.reinstated, ShouldEqual, "")
				})
			})

			Convey("But the request is blank", func() {
				status, body := fixture.Reinstate("")

				Convey("The server returns HTTP 400 - Bad Input", func() {
					So(status, ShouldEqual, http.StatusBadRequest)
				})

				Convey("The body should contain a helpful error message", func() {
					So(body, ShouldEqual, "You must provide a non-blank path.")
				})
			})

			Convey("And the request is well formed", func() {
				status, _ := fixture.Reinstate(initialRoot)

				Convey("The server informs the watcher", func() {
					So(fixture.watcher.reinstated, ShouldEqual, initialRoot)
				})
				Convey("The server returns HTTP 200 - OK", func() {
					So(status, ShouldEqual, http.StatusOK)
				})
			})
		})

		Convey("When the status of the executor is requested", func() {
			fixture.executor.status = "blah blah blah"
			statusCode, statusBody := fixture.RequestExecutorStatus()

			Convey("The server asks the executor its status and returns it", func() {
				So(statusBody, ShouldEqual, "blah blah blah")
			})

			Convey("The server returns HTTP 200 - OK", func() {
				So(statusCode, ShouldEqual, http.StatusOK)
			})
		})

		Convey("When a manual execution of the test packages is requested", func() {
			status := fixture.ManualExecution()
			update, _ := fixture.RequestLatest()

			Convey("The server invokes the executor using the watcher's listing and save the result", func() {
				So(update, ShouldResemble, &contract.CompleteOutput{Revision: initialRoot})
			})

			Convey("The server returns HTTP 200 - OK", func() {
				So(status, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func statusRotation(i, total int) string {
	switch i % total {
	case 0:
		return "executing"
	case 1:
		return "parsing"
	default:
		return "idle"
	}
}

/********* Server Fixture *********/

type ServerFixture struct {
	server   *HTTPServer
	watcher  *FakeWatcher
	executor *FakeExecutor
}

func (self *ServerFixture) ReceiveUpdate(update *contract.CompleteOutput) {
	self.server.ReceiveUpdate(update)
}

func (self *ServerFixture) RequestLatest() (*contract.CompleteOutput, *httptest.ResponseRecorder) {
	request, _ := http.NewRequest("GET", "http://localhost:8080/results", nil)
	response := httptest.NewRecorder()

	self.server.Results(response, request)

	decoder := json.NewDecoder(strings.NewReader(response.Body.String()))
	update := &contract.CompleteOutput{}
	decoder.Decode(update)
	return update, response
}

func (self *ServerFixture) QueryRootWatch(newclient bool) (string, int) {
	url := "http://localhost:8080/watch"
	if newclient {
		url += "?newclient=1"
	}
	request, _ := http.NewRequest("GET", url, nil)
	response := httptest.NewRecorder()

	self.server.Watch(response, request)

	return strings.TrimSpace(response.Body.String()), response.Code
}

func (self *ServerFixture) AdjustRootWatchMalformed() (status int, body string) {
	request, _ := http.NewRequest("POST", "http://localhost:8080/watch", nil)
	response := httptest.NewRecorder()

	self.server.Watch(response, request)

	status, body = response.Code, strings.TrimSpace(response.Body.String())
	return
}

func (self *ServerFixture) AdjustRootWatch(newRoot string) (status int, body string) {
	escapedRoot := url.QueryEscape(newRoot)
	request, _ := http.NewRequest("POST", "http://localhost:8080/watch?root="+escapedRoot, nil)
	response := httptest.NewRecorder()

	self.server.Watch(response, request)

	status, body = response.Code, strings.TrimSpace(response.Body.String())
	return
}

func (self *ServerFixture) IgnoreMalformed() (status int, body string) {
	request, _ := http.NewRequest("POST", "http://localhost:8080/ignore", nil)
	response := httptest.NewRecorder()

	self.server.Ignore(response, request)

	status, body = response.Code, strings.TrimSpace(response.Body.String())
	return
}

func (self *ServerFixture) Ignore(folder string) (status int, body string) {
	escapedFolder := url.QueryEscape(folder)
	request, _ := http.NewRequest("POST", "http://localhost:8080/ignore?path="+escapedFolder, nil)
	response := httptest.NewRecorder()

	self.server.Ignore(response, request)

	status, body = response.Code, strings.TrimSpace(response.Body.String())
	return
}

func (self *ServerFixture) ReinstateMalformed() (status int, body string) {
	request, _ := http.NewRequest("POST", "http://localhost:8080/reinstate", nil)
	response := httptest.NewRecorder()

	self.server.Reinstate(response, request)

	status, body = response.Code, strings.TrimSpace(response.Body.String())
	return
}

func (self *ServerFixture) Reinstate(folder string) (status int, body string) {
	escapedFolder := url.QueryEscape(folder)
	request, _ := http.NewRequest("POST", "http://localhost:8080/reinstate?path="+escapedFolder, nil)
	response := httptest.NewRecorder()

	self.server.Reinstate(response, request)

	status, body = response.Code, strings.TrimSpace(response.Body.String())
	return
}

func (self *ServerFixture) SetExecutorStatus(status string) {
	self.executor.status = status
	select {
	case self.executor.statusNotif <- true:
	default:
	}
}

func (self *ServerFixture) RequestExecutorStatus() (code int, status string) {
	request, _ := http.NewRequest("GET", "http://localhost:8080/status", nil)
	response := httptest.NewRecorder()

	self.server.Status(response, request)

	code, status = response.Code, strings.TrimSpace(response.Body.String())
	return
}

func (self *ServerFixture) ManualExecution() int {
	request, _ := http.NewRequest("POST", "http://localhost:8080/execute", nil)
	response := httptest.NewRecorder()

	self.server.Execute(response, request)
	nap, _ := time.ParseDuration("100ms")
	time.Sleep(nap)
	return response.Code
}

func newServerFixture() *ServerFixture {
	self := &ServerFixture{}
	self.watcher = newFakeWatcher()
	self.watcher.SetRootWatch(initialRoot)
	statusNotif := make(chan bool, 1)
	self.executor = newFakeExecutor("", statusNotif)
	self.server = NewHTTPServer(self.watcher, self.executor, statusNotif)
	return self
}

/********* Fake Watcher *********/

type FakeWatcher struct {
	root       string
	ignored    string
	reinstated string
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

func (self *FakeWatcher) Root() string {
	return self.root
}

func (self *FakeWatcher) Ignore(folder string)    { self.ignored = folder }
func (self *FakeWatcher) Reinstate(folder string) { self.reinstated = folder }

func (self *FakeWatcher) Deletion(folder string)       { panic("NOT SUPPORTED") }
func (self *FakeWatcher) Creation(folder string)       { panic("NOT SUPPORTED") }
func (self *FakeWatcher) IsWatched(folder string) bool { panic("NOT SUPPORTED") }
func (self *FakeWatcher) IsIgnored(folder string) bool { panic("NOT SUPPORTED") }

func newFakeWatcher() *FakeWatcher {
	return &FakeWatcher{}
}

/********* Fake Executor *********/

type FakeExecutor struct {
	status      string
	executed    bool
	statusNotif chan bool
}

func (self *FakeExecutor) Status() string {
	return self.status
}

func (self *FakeExecutor) ExecuteTests(watched []*contract.Package) *contract.CompleteOutput {
	return &contract.CompleteOutput{Revision: watched[0].Path}
}

func newFakeExecutor(status string, ch chan bool) *FakeExecutor {
	return &FakeExecutor{
		status,
		false,
		ch,
	}
}

var _ = fmt.Sprintf("hi")
