package system

import (
	"errors"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

///////////////////////////////////////////////////////////////////////////////
// Integration: ///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type Shell struct {
	coverage       bool
	gobin          string
	reportsPath    string
	currentCommand *Command
	abort          bool
}

func NewShell(gobin, reportsPath string, coverage bool) *Shell {
	return &Shell{
		coverage:    coverage,
		gobin:       gobin,
		reportsPath: reportsPath,
	}
}

func (self *Shell) GoTest(directory, packageName string, arguments []string) (output string, err error) {
	if self.currentCommand != nil {
		//reset if need be
		self.currentCommand.Kill()
		self.abort = false
	}
	reportFilename := strings.Replace(packageName, "/", "-", -1)
	reportPath := filepath.Join(self.reportsPath, reportFilename)
	reportData := reportPath + ".txt"
	reportHTML := reportPath + ".html"

	goconvey := findGoConvey(directory, self.gobin, packageName).Execute()
	compilation := compile(directory, self.gobin).Execute()

	withCoverage := runWithCoverage(compilation, goconvey, self.coverage, reportData, directory, self.gobin, arguments)
	self.currentCommand = withCoverage
	withCoverage.Execute()
	self.currentCommand = nil

	//early exit if we have an abort signal
	if self.abort {
		log.Print("Abort...")
		return withCoverage.Output, withCoverage.Error
	}

	final := runWithoutCoverage(compilation, withCoverage, goconvey, directory, self.gobin, arguments)
	self.currentCommand = final
	final.Execute()
	self.currentCommand = nil

	go generateReports(final, self.coverage, directory, self.gobin, reportData, reportHTML).Execute()

	return final.Output, final.Error
}

func (self *Shell) AbortGoTest() error {
	if self.currentCommand == nil {
		return errors.New("Not currently running a test command.")
	}
	self.abort = true
	return self.currentCommand.Kill()
}

///////////////////////////////////////////////////////////////////////////////
// Functional Core:////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

func findGoConvey(directory, gobin, packageName string) *Command {
	return NewCommand(directory, gobin, "list", "-f", "'{{.TestImports}}'", packageName)
}

func compile(directory, gobin string) *Command {
	return NewCommand(directory, gobin, "test", "-i")
}

func runWithCoverage(compile, goconvey *Command, coverage bool, reportPath, directory, gobin string, customArguments []string) *Command {
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

func runWithoutCoverage(compile, withCoverage, goconvey *Command, directory, gobin string, customArguments []string) *Command {
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

func generateReports(previous *Command, coverage bool, directory, gobin, reportData, reportHTML string) *Command {
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
	cmd        *exec.Cmd
	Output     string
	Error      error
}

func NewCommand(directory, executable string, arguments ...string) *Command {
	return &Command{
		directory:  directory,
		executable: executable,
		arguments:  arguments,
	}
}

func (this *Command) Execute() *Command {
	if len(this.executable) == 0 {
		return this
	}

	if len(this.Output) > 0 || this.Error != nil {
		return this
	}

	cmd := exec.Command(this.executable, this.arguments...)
	this.cmd = cmd
	cmd.Dir = this.directory
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	rawOutput, err := this.cmd.CombinedOutput()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if 2 == status.ExitStatus() {
					// don't pass back a "killed" message as an error, just return
					err = nil
				}
			}
		}
	}
	this.Error = err
	this.Output = string(rawOutput)
	return this
}

func (this *Command) Kill() error {
	// the this Command is usually a "go" command that has a subprocess, so we have to
	// use a negative PID for the syscall so it kills the whole tree
	if this.cmd != nil && this.cmd.Process != nil {
		return syscall.Kill(-this.cmd.Process.Pid, 15)
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////

const goconveyDSLImport = "github.com/smartystreets/goconvey/convey " // note the trailing space: we don't want to target packages nested in the /convey package.
var coverageStatementRE = regexp.MustCompile(`coverage: \d+\.\d% of statements(.*)\n`)
