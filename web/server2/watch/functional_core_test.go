package watch

import (
	"fmt"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCategorize(t *testing.T) {
	fileSystem := []*FileSystemItem{
		&FileSystemItem{
			Root:     "/",
			Path:     "/hello",
			Name:     "hello",
			Size:     12345566645,
			Modified: 712342134,
			IsFolder: true,
		},
		&FileSystemItem{
			Root:     "/",
			Path:     "/1/hello/world.txt",
			Name:     "world.txt",
			Size:     3,
			Modified: 5,
			IsFolder: false,
		},
		&FileSystemItem{
			Root:     "/",
			Path:     "/1/2/3/4/5/hello/world.go",
			Name:     "world.go",
			Size:     3,
			Modified: 5,
			IsFolder: false,
		},
		&FileSystemItem{
			Root:     "/",
			Path:     "/hello/.world.go",
			Name:     ".world.go",
			Size:     3,
			Modified: 5,
			IsFolder: false,
		},
		&FileSystemItem{
			Root:     "/",
			Path:     "/.hello",
			Name:     ".hello",
			Size:     3,
			Modified: 5,
			IsFolder: true,
		},
		&FileSystemItem{
			Root:     "/",
			Path:     "/.hello/hello",
			Name:     "hello",
			Size:     2,
			Modified: 3,
			IsFolder: true,
		},
		&FileSystemItem{
			Root:     "/",
			Path:     "/.hello/world.go",
			Name:     "world.go",
			Size:     3,
			Modified: 5,
			IsFolder: false,
		},
		&FileSystemItem{
			Root:     "/",
			Path:     "/hello/hi.goconvey",
			Name:     "hi.goconvey",
			Size:     2,
			Modified: 3,
			IsFolder: false,
		},
		&FileSystemItem{
			Root:     "/",
			Path:     "/hello2/.goconvey",
			Name:     ".goconvey",
			Size:     2,
			Modified: 3,
			IsFolder: false,
		},
	}

	Convey("A stream of file system items should be categorized correctly", t, func() {
		items := make(chan *FileSystemItem)

		go func() {
			for _, item := range fileSystem {
				items <- item
			}
			close(items)
		}()

		folders, profiles, goFiles := Categorize(items)

		So(folders, ShouldResemble, fileSystem[:1])
		So(profiles, ShouldResemble, fileSystem[7:8])
		So(goFiles, ShouldResemble, fileSystem[2:3])
	})
}

func TestParseProfile(t *testing.T) {
	var parseProfileTestCases = []struct {
		SKIP           bool
		description    string
		input          string
		resultIgnored  bool
		resultTestArgs []string
	}{
		{
			SKIP:           false,
			description:    "Blank profile",
			input:          "",
			resultIgnored:  false,
			resultTestArgs: []string{},
		},
		{
			SKIP:           false,
			description:    "All lines are blank or whitespace",
			input:          "\n \n \t\t\t  \n \n \n",
			resultIgnored:  false,
			resultTestArgs: []string{},
		},
		{
			SKIP:           false,
			description:    "Ignored package, no args included",
			input:          "IGNORE\n-timeout=4s",
			resultIgnored:  true,
			resultTestArgs: []string{},
		},
		{
			SKIP:          false,
			description:   "Ignore directive is commented, all args are included",
			input:         "#IGNORE\n-timeout=4s\n-parallel=5",
			resultIgnored: false,
			resultTestArgs: []string{
				"-timeout=4s",
				"-parallel=5",
			},
		},
		{
			SKIP:          false,
			description:   "No ignore directive, all args are included",
			input:         "-run=TestBlah\n-timeout=42s",
			resultIgnored: false,
			resultTestArgs: []string{
				"-run=TestBlah",
				"-timeout=42s",
			},
		},
		{
			SKIP:          false,
			description:   "Some args are commented, therefore ignored",
			input:         "-run=TestBlah\n #-timeout=42s",
			resultIgnored: false,
			resultTestArgs: []string{
				"-run=TestBlah",
			},
		},
		{
			SKIP:           false,
			description:    "All args are commented, therefore all are ignored",
			input:          "#-run=TestBlah\n//-timeout=42",
			resultIgnored:  false,
			resultTestArgs: []string{},
		},
		{
			SKIP:           false,
			description:    "We ignore certain flags like -v and -cover* because they are specified by the shell",
			input:          "-v\n-cover\n-coverprofile=blah.out",
			resultIgnored:  false,
			resultTestArgs: []string{},
		},
	}

	for i, test := range parseProfileTestCases {
		if test.SKIP {
			SkipConvey(fmt.Sprintf("Profile Parsing, Test Case #%d: %s (SKIPPED)", i, test.description), t, nil)
		} else {
			Convey(fmt.Sprintf("Profile Parsing, Test Case #%d: %s", i, test.description), t, func() {
				ignored, testArgs := ParseProfile(test.input)

				So(ignored, ShouldEqual, test.resultIgnored)
				So(testArgs, ShouldResemble, test.resultTestArgs)
			})
		}
	}
}
