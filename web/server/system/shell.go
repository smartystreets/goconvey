package system

import (
	"fmt"
	"log"
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
	short    bool
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
	if !self.coverage {
		return self.runWithoutCoverage(directory)
	}

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
	arguments := []string{
		"test", "-v", fmt.Sprintf("-short=%t", self.short),
		"-timeout=-42s", "-covermode=set", "-coverprofile=" + profile,
	}
	return self.execute(directory, self.gobin, arguments...)
}
func (self *Shell) runWithoutCoverage(directory string) (string, error) {
	arguments := []string{"test", "-v", fmt.Sprintf("-short=%t", self.short), "-timeout=-42s"}
	return self.execute(directory, self.gobin, arguments...)
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

func NewShell(gobin string, short bool, cover bool, reports string) *Shell {
	self := new(Shell)
	self.gobin = gobin
	self.short = short
	self.coverage = self.coverageEnabled(cover, reports)
	self.reports = reports
	return self
}

func (self *Shell) coverageEnabled(cover bool, reports string) bool {
	return (cover &&
		goVersion_1_2_orGreater() &&
		self.coverToolInstalled() &&
		ensureReportDirectoryExists(reports))
}

func goVersion_1_2_orGreater() bool {
	version := runtime.Version() // 'go1.2....'
	major, minor := version[2], version[4]
	version_1_2 := major >= byte('1') && minor >= byte('2')
	if !version_1_2 {
		log.Printf(pleaseUpgradeGoVersion, version)
		return false
	}
	return true
}

func (self *Shell) coverToolInstalled() bool {
	working, err := os.Getwd()
	if err != nil {
		working = "."
	}
	output, _ := self.execute(working, "go", "tool", "cover")
	installed := strings.Contains(output, "Usage of 'go tool cover':")
	if !installed {
		log.Print(coverToolMissing)
		return false
	}
	return true
}

func ensureReportDirectoryExists(reports string) bool {
	if exists(reports) {
		return true
	}

	if err := os.Mkdir(reports, 0755); err == nil {
		return true
	}

	log.Printf(reportDirectoryUnavailable, reports)
	return false
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

const (
	pleaseUpgradeGoVersion     = "Go version is less that 1.2 (%s), please upgrade to the latest stable version to enable coverage reporting.\n"
	coverToolMissing           = "Go cover tool is not installed or not accessible: `go get code.google.com/p/go.tools/cmd/cover`\n"
	reportDirectoryUnavailable = "Could not find or create the coverage report directory (at: '%s'). You probably won't see any coverage statistics...\n"
)
