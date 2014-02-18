package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Shell struct {
	coverage bool
	gobin    string
	reports  string
}

func (self *Shell) GoTest(directory, packageName string) (output string, err error) {
	output, err = self.compileDependencies(directory)
	if err == nil {
		output, err = self.goTest(directory, packageName)
	}
	return
}

func (self *Shell) compileDependencies(directory string) (output string, err error) {
	return self.execute(directory, self.gobin, "test", "-i")
}

func (self *Shell) goTest(directory, packageName string) (output string, err error) {
	reportFilename := strings.Replace(packageName, string(os.PathSeparator), "-", -1)
	reportPath := filepath.Join(self.reports, reportFilename)
	profile := reportPath + ".txt"
	output, err = self.runWithCoverage(directory, profile)

	if err != nil && self.coverage {
		output, err = self.runWithoutCoverage(directory)
	} else if self.coverage {
		self.generateCoverageReports(directory, profile, reportPath+".html")
	}
	return
}

func (self *Shell) runWithCoverage(directory, profile string) (string, error) {
	return self.execute(directory, self.gobin, "test", "-v", "-timeout=-42s", "-covermode=set", "-coverprofile="+profile)
}
func (self *Shell) runWithoutCoverage(directory string) (string, error) {
	return self.execute(directory, self.gobin, "test", "-v", "-timeout=-42s")
}

func (self *Shell) generateCoverageReports(directory, profile, html string) {
	self.execute(directory, self.gobin, "tool", "cover", "-html="+profile, "-o", html)
}

func (self *Shell) execute(directory, name string, args ...string) (output string, err error) {
	command := exec.Command(name, args...)
	command.Dir = directory
	rawOutput, err := command.CombinedOutput()
	output = string(rawOutput)
	return
}

func (self *Shell) Getenv(key string) string {
	return os.Getenv(key)
}

func (self *Shell) Setenv(key, value string) error {
	if self.Getenv(key) != value {
		return os.Setenv(key, value)
	}
	return nil
}

func NewShell(gobin string, cover bool, reports string) *Shell {
	self := new(Shell)
	self.gobin = gobin
	self.reports = reports
	if cover && goVersion_1_2_orGreater() {
		self.coverage = true
	}
	return self
}

func goVersion_1_2_orGreater() bool {
	version := runtime.Version() // 'go1.2....'
	major, minor := version[2], version[4]
	return major >= byte('1') && minor >= byte('2')
}
