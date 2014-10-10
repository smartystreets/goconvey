package system

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// Integration: ///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type Shell struct {
	coverage    bool
	gobin       string
	reportsPath string
}

func NewShell(gobin, reportsPath string, coverage bool) *Shell {
	return &Shell{
		coverage:    coverage,
		gobin:       gobin,
		reportsPath: reportsPath,
	}
}

func (self *Shell) GoTest(directory, packageName string, arguments []string) (output string, err error) {
	reportFilename := strings.Replace(packageName, "/", "-", -1)
	reportPath := filepath.Join(self.reportsPath, reportFilename)
	reportData := reportPath + ".txt"
	reportHTML := reportPath + ".html"

	goconvey := findGoConvey(directory, self.gobin, packageName).Execute()
	compilation := compile(directory, self.gobin).Execute()
	withCoverage := runWithCoverage(compilation, goconvey, self.coverage, reportData, directory, self.gobin, arguments).Execute()
	final := runWithoutCoverage(compilation, withCoverage, goconvey, directory, self.gobin, arguments).Execute()
	go generateReports(final, self.coverage, directory, self.gobin, reportData, reportHTML).Execute()

	return final.Output, final.Error
}

///////////////////////////////////////////////////////////////////////////////
// Functional Core:////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

func findGoConvey(directory, gobin, packageName string) Command {
	return NewCommand(directory, gobin, "list", "-f", "'{{.TestImports}}'", packageName)
}

func compile(directory, gobin string) Command {
	return NewCommand(directory, gobin, "test", "-i")
}

func runWithCoverage(compile, goconvey Command, coverage bool, reportPath, directory, gobin string, customArguments []string) Command {
	if compile.Error != nil {
		return compile
	}

	if !coverage {
		return compile
	}

	arguments := []string{"test", "-v", "-coverprofile=" + reportPath}

	if !strings.Contains(strings.Join(customArguments, "\t"), "-covermode=") {
		arguments = append(arguments, "-covermode=set")
	}

	if strings.Contains(goconvey.Output, goconveyDSLImport) {
		arguments = append(arguments, "-json")
	}

	arguments = append(arguments, customArguments...)

	return NewCommand(directory, gobin, arguments...)
}

func runWithoutCoverage(compile, withCoverage, goconvey Command, directory, gobin string, customArguments []string) Command {
	if compile.Error != nil {
		return compile
	}

	if coverageStatementRE.MatchString(withCoverage.Output) {
		return withCoverage
	}

	arguments := []string{"test", "-v"}
	if strings.Contains(goconvey.Output, goconveyDSLImport) {
		arguments = append(arguments, "-json")
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

func (this Command) Execute() Command {
	if len(this.executable) == 0 {
		return this
	}

	if len(this.Output) > 0 || this.Error != nil {
		return this
	}

	command := exec.Command(this.executable, this.arguments...)
	command.Dir = this.directory
	var rawOutput []byte
	rawOutput, this.Error = command.CombinedOutput()
	this.Output = string(rawOutput)
	return this
}

///////////////////////////////////////////////////////////////////////////////

const goconveyDSLImport = "github.com/smartystreets/goconvey/convey " // note the trailing space: we don't want to target packages nested in the /convey package.
var coverageStatementRE = regexp.MustCompile(`coverage: \d+\.\d% of statements(.*)\n`)
