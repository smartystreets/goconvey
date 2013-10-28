package executor

import (
	"github.com/smartystreets/goconvey/web/server/parser"
)

type Parser func(packageName, output string) *parser.PackageResult
