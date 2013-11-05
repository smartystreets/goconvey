package parser

import (
	"github.com/smartystreets/goconvey/web/server/contract"
)

type Parser struct {
	parser func(*contract.PackageResult, string)
}

func (self *Parser) Parse(packages []*contract.Package) {
	for _, p := range packages {
		if p.Active {
			self.parser(p.Result, p.Output)
		}
	}
}

func NewParser(helper func(*contract.PackageResult, string)) *Parser {
	self := &Parser{}
	self.parser = helper
	return self
}
