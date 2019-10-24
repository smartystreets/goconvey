package system

import (
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// Integration: ///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type Shell struct {
	coverage       bool
	gobin          string
	reportsPath    string
	defaultTimeout string
}

func NewShell(gobin, reportsPath string, coverage bool, defaultTimeout string) *Shell {
	return &Shell{
		coverage:       coverage,
		gobin:          gobin,
		reportsPath:    reportsPath,
		defaultTimeout: defaultTimeout,
	}
}

func (s *Shell) GoTest(directory, packageName string, tags, arguments []string) (output string, err error) {
	reportFilename := strings.Replace(packageName, "/", "-", -1)
	reportPath := filepath.Join(s.reportsPath, reportFilename)
	reportData := reportPath + ".txt"
	reportHTML := reportPath + ".html"
	tagsArg := "-tags=" + strings.Join(tags, ",")

	goconvey := findGoConvey(directory, s.gobin, packageName, tagsArg).Execute()
	compilation := compile(directory, s.gobin, tagsArg).Execute()
	withCoverage := runWithCoverage(compilation, goconvey, s.coverage, reportData, directory, s.gobin, s.defaultTimeout, tagsArg, arguments).Execute()
	final := runWithoutCoverage(compilation, withCoverage, goconvey, directory, s.gobin, s.defaultTimeout, tagsArg, arguments).Execute()
	go generateReports(final, s.coverage, directory, s.gobin, reportData, reportHTML).Execute()

	return final.Output, final.Error
}

///////////////////////////////////////////////////////////////////////////////
// Functional Core:////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

func findGoConvey(directory, gobin, packageName, tagsArg string) Command {
	return NewCommand(directory, gobin, "list", "-f", "'{{.TestImports}}{{.XTestImports}}'", tagsArg, packageName)
}

func compile(directory, gobin, tagsArg string) Command {
	return NewCommand(directory, gobin, "test", "-i", tagsArg)
}

func runWithCoverage(compile, goconvey Command, coverage bool, reportPath, directory, gobin, defaultTimeout, tagsArg string, customArguments []string) Command {
	if compile.Error != nil || goconvey.Error != nil {
		return compile
	}

	if !coverage {
		return compile
	}

	arguments := []string{"test", "-v", "-coverprofile=" + reportPath, tagsArg}

	customArgsText := strings.Join(customArguments, "\t")
	if !strings.Contains(customArgsText, "-covermode=") {
		arguments = append(arguments, "-covermode=set")
	}

	if !strings.Contains(customArgsText, "-timeout=") {
		arguments = append(arguments, "-timeout="+defaultTimeout)
	}

	if strings.Contains(goconvey.Output, goconveyDSLImport) {
		arguments = append(arguments, "-convey-json")
	}

	arguments = append(arguments, customArguments...)

	return NewCommand(directory, gobin, arguments...)
}

func runWithoutCoverage(compile, withCoverage, goconvey Command, directory, gobin, defaultTimeout, tagsArg string, customArguments []string) Command {
	if compile.Error != nil {
		return compile
	}

	if goconvey.Error != nil {
		log.Println(gopathProblem, goconvey.Output, goconvey.Error)
		return goconvey
	}

	if coverageStatementRE.MatchString(withCoverage.Output) {
		return withCoverage
	}

	log.Printf("Coverage output: %v", withCoverage.Output)

	log.Print("Run without coverage")

	arguments := []string{"test", "-v", tagsArg}
	customArgsText := strings.Join(customArguments, "\t")
	if !strings.Contains(customArgsText, "-timeout=") {
		arguments = append(arguments, "-timeout="+defaultTimeout)
	}

	if strings.Contains(goconvey.Output, goconveyDSLImport) {
		arguments = append(arguments, "-convey-json")
	}
	arguments = append(arguments, customArguments...)
	return NewCommand(directory, gobin, arguments...)
}

func generateReports(previous Command, coverage bool, directory, gobin, reportData, reportHTML string) Command {
	if previous.Error != nil {
		return previous
	}

	if !coverage {
		return previous
	}

	return NewCommand(directory, gobin, "tool", "cover", "-html="+reportData, "-o", reportHTML)
}

///////////////////////////////////////////////////////////////////////////////
// Imperative Shell: //////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type Command struct {
	directory  string
	executable string
	arguments  []string

	Output string
	Error  error
}

func NewCommand(directory, executable string, arguments ...string) Command {
	return Command{
		directory:  directory,
		executable: executable,
		arguments:  arguments,
	}
}

func (c Command) Execute() Command {
	if len(c.executable) == 0 {
		return c
	}

	if len(c.Output) > 0 || c.Error != nil {
		return c
	}

	command := exec.Command(c.executable, c.arguments...)
	command.Dir = c.directory
	var rawOutput []byte
	rawOutput, c.Error = command.CombinedOutput()
	c.Output = string(rawOutput)
	return c
}

///////////////////////////////////////////////////////////////////////////////

const goconveyDSLImport = "github.com/smartystreets/goconvey/convey " // note the trailing space: we don't want to target packages nested in the /convey package.
const gopathProblem = "Please run goconvey from within $GOPATH/src (also, symlinks might be problematic). Output and Error: "

var coverageStatementRE = regexp.MustCompile(`(?m)^coverage: \d+\.\d% of statements(.*)$|^panic: test timed out after `)
