package executor

import (
	"github.com/smartystreets/goconvey/web/server/parser"
)

// type Parser func(packageName, output string) *parser.PackageResult

type Parser interface {
	Parse(packageName, output string) *parser.PackageResult
}

type Tester interface {
	TestAll(folders []string) (output []string)
}
