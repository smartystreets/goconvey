package watch

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/messaging"
)

func TestCategorize(t *testing.T) {
	fileSystem := []*FileSystemItem{
		&FileSystemItem{
			Root:     "/.hello",
			Path:     "/.hello",
			Name:     "hello",
			IsFolder: true,
		},
		&FileSystemItem{
			Root:     "/.hello",
			Path:     "/.hello/1/hello/world.txt",
			Name:     "world.txt",
			IsFolder: false,
		},
		&FileSystemItem{
			Root:     "/.hello",
			Path:     "/.hello/1/2/3/4/5/hello/world.go",
			Name:     "world.go",
			IsFolder: false,
		},
		&FileSystemItem{
			Root:     "/.hello",
			Path:     "/.hello/world.go",
			Name:     "world.go",
			IsFolder: false,
		},
		&FileSystemItem{
			Root:     "/.hello",
			Path:     "/.hello/hello/.world.go",
			Name:     ".world.go",
			IsFolder: false,
		},
		&FileSystemItem{
			Root:     "/.hello",
			Path:     "/.hello/.hello",
			Name:     ".hello",
			IsFolder: true,
		},
		&FileSystemItem{
			Root:     "/.hello",
			Path:     "/.hello/.hello/hello",
			Name:     "hello",
			IsFolder: true,
		},
		&FileSystemItem{
			Root:     "/.hello",
			Path:     "/.hello/.hello/world.go",
			Name:     "world.go",
			IsFolder: false,
		},
		&FileSystemItem{
			Root:     "/.hello",
			Path:     "/.hello/hello/hi.goconvey",
			Name:     "hi.goconvey",
			IsFolder: false,
		},
		&FileSystemItem{
			Root:     "/.hello",
			Path:     "/.hello/hello2/.goconvey",
			Name:     ".goconvey",
			IsFolder: false,
		},
		&FileSystemItem{
			Root:     "/.hello",
			Path:     "/.hello/_hello",
			Name:     "_hello",
			IsFolder: true,
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

		folders, profiles, goFiles := Categorize(items, "/.hello")

		So(folders, ShouldResemble, fileSystem[:1])
		So(profiles, ShouldResemble, fileSystem[8:9])
		So(goFiles, ShouldResemble, fileSystem[2:4])
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
			description:    "We ignore certain flags like -v and -cover and -coverprofile because they are specified by the shell",
			input:          "-v\n-cover\n-coverprofile=blah.out",
			resultIgnored:  false,
			resultTestArgs: []string{},
		},
		{
			SKIP:           false,
			description:    "We allow certain coverage flags like -coverpkg and -covermode",
			input:          "-coverpkg=blah\n-covermode=atomic",
			resultIgnored:  false,
			resultTestArgs: []string{"-coverpkg=blah", "-covermode=atomic"},
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

func TestCreateFolders(t *testing.T) {
	Convey("File system items that represent folders should be converted to folder structs correctly", t, func() {
		expected := map[string]*messaging.Folder{
			"/root/1":     &messaging.Folder{Path: "/root/1", Root: "/root"},
			"/root/1/2":   &messaging.Folder{Path: "/root/1/2", Root: "/root"},
			"/root/1/2/3": &messaging.Folder{Path: "/root/1/2/3", Root: "/root"},
		}

		inputs := []*FileSystemItem{
			&FileSystemItem{Path: "/root/1", Root: "/root", IsFolder: true},
			&FileSystemItem{Path: "/root/1/2", Root: "/root", IsFolder: true},
			&FileSystemItem{Path: "/root/1/2/3", Root: "/root", IsFolder: true},
		}

		actual := CreateFolders(inputs)

		for key, actualValue := range actual {
			So(actualValue, ShouldResemble, expected[key])
		}
	})
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

func TestAttachProfiles(t *testing.T) {
	Convey("Subject: Attaching profile information to a folder", t, func() {
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

		profiles := []*FileSystemItem{
			&FileSystemItem{
				Path:             "/root/too-shallow.goconvey",
				ProfileDisabled:  true,
				ProfileArguments: []string{"1", "2"},
			},
			&FileSystemItem{
				Path:             "/root/1/2/hi.goconvey",
				ProfileDisabled:  true,
				ProfileArguments: []string{"1", "2"},
			},
			&FileSystemItem{
				Path:             "/root/1/2/3/4/does-not-exist",
				ProfileDisabled:  true,
				ProfileArguments: []string{"1", "2", "3", "4"},
			},
		}

		Convey("Profiles that match folders should be merged with those folders", func() {
			AttachProfiles(folders, profiles)

			Convey("No profiles matched the first folder, so no assignments should occur", func() {
				So(folders["/root/1"].Disabled, ShouldBeFalse)
				So(folders["/root/1"].TestArguments, ShouldBeEmpty)
			})

			Convey("The second folder should match the first profile", func() {
				So(folders["/root/1/2"].Disabled, ShouldBeTrue)
				So(folders["/root/1/2"].TestArguments, ShouldResemble, []string{"1", "2"})
			})

			Convey("No profiles match the third folder so no assignments should occur", func() {
				So(folders["/root/1/2/3"].Disabled, ShouldBeFalse)
				So(folders["/root/1/2/3"].TestArguments, ShouldBeEmpty)
			})
		})
	})
}

func TestMarkIgnored(t *testing.T) {
	Convey("Subject: folders that have been ignored should be marked as such", t, func() {
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

		Convey("When there are no ignored folders", func() {
			ignored := map[string]struct{}{}
			MarkIgnored(folders, ignored)

			Convey("No folders should be marked as ignored", func() {
				So(folders["/root/1"].Ignored, ShouldBeFalse)
				So(folders["/root/1/2"].Ignored, ShouldBeFalse)
				So(folders["/root/1/2/3"].Ignored, ShouldBeFalse)
			})
		})
		Convey("When there are ignored folders", func() {
			ignored := map[string]struct{}{"1/2": struct{}{}}
			MarkIgnored(folders, ignored)

			Convey("The ignored folders should be marked as ignored", func() {
				So(folders["/root/1"].Ignored, ShouldBeFalse)
				So(folders["/root/1/2"].Ignored, ShouldBeTrue)
				So(folders["/root/1/2/3"].Ignored, ShouldBeFalse)
			})
		})
	})
}

func TestActiveFolders(t *testing.T) {
	Convey("Subject: Folders that are not ignored or disabled are active", t, func() {
		folders := map[string]*messaging.Folder{
			"/root/1": &messaging.Folder{
				Path:    "/root/1",
				Root:    "/root",
				Ignored: true,
			},
			"/root/1/2": &messaging.Folder{
				Path: "/root/1/2",
				Root: "/root",
			},
			"/root/1/2/3": &messaging.Folder{
				Path:     "/root/1/2/3",
				Root:     "/root",
				Disabled: true,
			},
		}

		active := ActiveFolders(folders)

		So(len(active), ShouldEqual, 1)
		So(active["/root/1/2"], ShouldResemble, folders["/root/1/2"])
	})
}

func TestSum(t *testing.T) {
	Convey("Subject: file system items within specified directores should be counted and summed", t, func() {
		folders := map[string]*messaging.Folder{
			"/root/1": &messaging.Folder{Path: "/root/1", Root: "/root", Ignored: true},
		}
		items := []*FileSystemItem{
			&FileSystemItem{Size: 1, Modified: 3, Path: "/root/1/hi.go"},
			&FileSystemItem{Size: 7, Modified: 13, Path: "/root/1/bye.go"},
			&FileSystemItem{Size: 33, Modified: 45, Path: "/root/1/2/salutations.go"}, // not counted
		}

		So(Sum(folders, items), ShouldEqual, 1+3+7+13)
	})
}
