package api

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/parser"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
					So(status, ShouldEqual, 200)
				})
			})
		})

		Convey("When the root watch is queried", func() {
			Convey("But the request is malformed", func() {
				Convey("The server returns 400", nil)
			})

			Convey("And the request is well formed", func() {
				Convey("The server returns it", nil)
				Convey("The server returns 200", nil)
			})
		})

		Convey("When the root watch is adjusted", func() {
			Convey("But the request is malformed", func() {
				Convey("The server returns 400", nil)
			})

			Convey("And the new root exists", func() {
				Convey("The server informs the watcher", nil)
				Convey("The server returns 200", nil)
			})

			Convey("And the new root does NOT exist", func() {
				Convey("The server returns 404", nil)
			})
		})

		Convey("When a packge is ignored", func() {
			Convey("But the request is malformed", func() {
				Convey("The server returns 400", nil)
			})

			Convey("And the request is well formed", func() {
				Convey("The server informs the watcher", nil)
				Convey("The server returns 200", nil)
			})
		})

		Convey("When a package is reinstated", func() {
			Convey("But the request is malformed", func() {
				Convey("The server returns 400", nil)
			})

			Convey("And the request is well formed", func() {
				Convey("The server informs the watcher", nil)
				Convey("The server returns 200", nil)
			})
		})

		Convey("When the status is requested", func() {
			Convey("The server returns it", nil)
			Convey("The server returns 200", nil)
		})

		Convey("When a manual execution of the test packages is requested", func() {
			Convey("The server invokes the executor using the watcher's listing", nil)
			Convey("The server returns 200", nil)
		})
	})
}

type ServerFixture struct {
	server *HTTPServer
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

func newServerFixture() *ServerFixture {
	self := &ServerFixture{}
	self.server = NewHTTPServer()
	return self
}

/*


TestServer
	Subject: Server responds to requests appropriately

		Given an update is received
			When the update is requested
				The server returns it
				The server returns 200

		When the root watch is queried
			But the request is malformed
				The server returns 400
			And the request is well formed
				The server returns it
				The server returns 200

		When the root watch is adjusted
			But the request is malformed
				The server returns 400
			And the new root exists
				The server informs the watcher
				The server returns 200
			And the new root does NOT exist
				The server returns 404

		When a packge is ignored
			But the request is malformed
				The server returns 400
			And the request is well formed
				The server informs the watcher
				The server returns 200

		When a package is reinstated
			But the request is malformed
				The server returns 400
			And the request is well formed
				The server informs the watcher
				The server returns 200

		When the status is requested
			The server returns it
			The server returns 200

		When a manual execution of the test packages is requested
			The server invokes the executor using the watcher's listing
			The server returns 200
*/

var _ = fmt.Sprintf("hi")
