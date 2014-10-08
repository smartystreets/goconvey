package watch

import (
	"fmt"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server2/messaging"
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
			SKIP:           false,
			description:    "Ignore directive is commented, all args are included",
			input:          "#IGNORE\n-timeout=4s\n-parallel=5",
			resultIgnored:  false,
			resultTestArgs: []string{"-timeout=4s", "-parallel=5"},
		},
		{
			SKIP:           false,
			description:    "No ignore directive, all args are included",
			input:          "-run=TestBlah\n-timeout=42s",
			resultIgnored:  false,
			resultTestArgs: []string{"-run=TestBlah", "-timeout=42s"},
		},
		{
			SKIP:           false,
			description:    "Some args are commented, therefore ignored",
			input:          "-run=TestBlah\n #-timeout=42s",
			resultIgnored:  false,
			resultTestArgs: []string{"-run=TestBlah"},
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

func TestLimitDepth(t *testing.T) {
	Convey("Subject: Limiting folders based on relative depth from a common root", t, func() {

		folders := map[string]*messaging.Folder{
			"/root/1": &messaging.Folder{
				Path: "/root/1",
				Root: "/root",
			},
			"/root/1/2": &messaging.Folder{
				Path: "/root/1/2",
				Root: "/root",
			},
			"/root/1/2/3": &messaging.Folder{
				Path: "/root/1/2/3",
				Root: "/root",
			},
		}

		Convey("When there is no depth limit", func() {
			LimitDepth(folders, -1)

			Convey("No folders should be excluded", func() {
				So(len(folders), ShouldEqual, 3)
			})
		})

		Convey("When there is a limit", func() {
			LimitDepth(folders, 2)

			Convey("The deepest folder (in this case) should be excluded", func() {
				So(len(folders), ShouldEqual, 2)
				_, exists := folders["/root/1/2/3"]
				So(exists, ShouldBeFalse)
			})
		})
	})
}
