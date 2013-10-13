// reference:
// http://stackoverflow.com/questions/10171941/need-help-understanding-why-select-isnt-blocking-forever

package main

import (
	"crypto/md5"
	"fmt"
	"github.com/smartystreets/goconvey/web/server/parser"
	"github.com/smartystreets/goconvey/web/server/results"
	"io"
	"os"
	"os/exec"
	"strings"
)

func runTests(path string) *TestPackage {
	if err := os.Chdir(path); err != nil {
		panic(fmt.Sprintf("Could not change dir to %s", path))
	}

	exec.Command("go", "test", "-i").Run()
	output, _ := exec.Command("go", "test", "-v", "-timeout=-42s").CombinedOutput()
	stringOutput := string(output)

	packageIndex := strings.Index(path, "/src/")
	packageName := path[packageIndex+len("/src/"):]
	result := parser.ParsePackageResults(packageName, stringOutput)
	return &TestPackage{
		Path:   path,
		Output: stringOutput,
		Parsed: result,
	}
}

func worker(in chan string, out chan *TestPackage) {
	for path := range in {
		out <- runTests(path)
	}
}

func main() {
	folders := []string{
		"/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/assertions",
		"/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/convey",
		"/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server",
		"/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/parser",
		"/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/reporting",
		"/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/printing",
		"/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/execution",
	}

	numWorkers := len(folders)

	// spawn workers
	in, out := make(chan string), make(chan *TestPackage)
	for i := 0; i < numWorkers; i++ {
		go worker(in, out)
	}

	// schedule tasks
	go func() {
		for _, f := range folders {
			in <- f
		}
	}()

	results := []*TestPackage{}
	for _ = range folders {
		results = append(results, <-out)
	}

	revision := md5.New()

	for _, output := range results {
		io.WriteString(revision, output.Path)
	}

	fmt.Println(string(revision.Sum(nil)))
}

type TestPackage struct {
	Path   string
	Output string
	Parsed *results.PackageResult
}
