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
	coverage      bool
	gobin         string
	reportsPath   string
	shortArgument string
}

func (self *Shell) GoTest(directory, packageName string) (output string, err error) {
	output, err = self.compilePackageDependencies(directory)
	if err == nil {
		output, err = self.goTest(directory, packageName)
	}
	return
}

func (self *Shell) compilePackageDependencies(directory string) (output string, err error) {
	return execute(directory, self.gobin, "test", "-i")
}

func (self *Shell) goTest(directory, packageName string) (output string, err error) {
	if !self.coverage {
		return self.runWithoutCoverage(directory, packageName)
	}

	return self.tryRunWithCoverage(directory, packageName)
}

func (self *Shell) tryRunWithCoverage(directory, packageName string) (output string, err error) {
	profileName := self.composeProfileName(packageName)
	output, err = self.runWithCoverage(directory, packageName, profileName+".txt")

	if err != nil && self.coverage {
		output, err = self.runWithoutCoverage(directory, packageName)
	} else if self.coverage {
		self.generateCoverageReports(directory, profileName+".txt", profileName+".html")
	}
	return
}

func (self *Shell) composeProfileName(packageName string) string {
	reportFilename := strings.Replace(packageName, string(os.PathSeparator), "-", -1)
	reportPath := filepath.Join(self.reportsPath, reportFilename)
	return reportPath
}

func (self *Shell) runWithCoverage(directory, packageName, profile string) (string, error) {
	arguments := []string{"test", "-v", self.shortArgument, "-covermode=set", "-coverprofile=" + profile}
	arguments = append(arguments, self.jsonOrNot(directory, packageName)...)
	return execute(directory, self.gobin, arguments...)
}

func (self *Shell) runWithoutCoverage(directory, packageName string) (string, error) {
	arguments := []string{"test", "-v", self.shortArgument}
	arguments = append(arguments, self.jsonOrNot(directory, packageName)...)
	return execute(directory, self.gobin, arguments...)
}

func (self *Shell) jsonOrNot(directory, packageName string) []string {
	imports, err := execute(directory, self.gobin, "list", "-f", "'{{.TestImports}}'", packageName)
	if !strings.Contains(imports, goconveyDSLImport) && err == nil {
		return []string{}
	}
	return []string{"-json"}
}

func (self *Shell) generateCoverageReports(directory, profile, html string) {
	execute(directory, self.gobin, "tool", "cover", "-html="+profile, "-o", html)
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
	self.shortArgument = fmt.Sprintf("-short=%t", short)
	self.coverage = coverageEnabled(cover, reports)
	self.reportsPath = reports
	return self
}

func coverageEnabled(cover bool, reports string) bool {
	return (cover &&
		goVersion_1_2_orGreater() &&
		coverToolInstalled() &&
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

func coverToolInstalled() bool {
	working, err := os.Getwd()
	if err != nil {
		working = "."
	}
	output, _ := execute(working, "go", "tool", "cover")
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

func execute(directory, name string, args ...string) (output string, err error) {
	command := exec.Command(name, args...)
	command.Dir = directory
	rawOutput, err := command.CombinedOutput()
	output = string(rawOutput)
	return
}

const (
	goconveyDSLImport          = "github.com/smartystreets/goconvey/convey " // note the trailing space: we don't want to target packages nested in the /convey package.
	pleaseUpgradeGoVersion     = "Go version is less that 1.2 (%s), please upgrade to the latest stable version to enable coverage reporting.\n"
	coverToolMissing           = "Go cover tool is not installed or not accessible: `go get code.google.com/p/go.tools/cmd/cover`\n"
	reportDirectoryUnavailable = "Could not find or create the coverage report directory (at: '%s'). You probably won't see any coverage statistics...\n"
)
