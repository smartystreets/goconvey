package parser

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"log"
)

type Parser struct {
	parser func(*contract.PackageResult, string)
}

func (self *Parser) Parse(packages []*contract.Package) {
	for _, p := range packages {
		if p.Active {
			self.parser(p.Result, p.Output)
		} else {
			// TODO: should set Outcome to ignored
			log.Printf("Skipping parsing: %s\n", p.Name)
		}
		log.Printf("[%s]: %s\n", p.Result.Outcome, p.Name)
	}
}

func NewParser(helper func(*contract.PackageResult, string)) *Parser {
	self := &Parser{}
	self.parser = helper
	return self
}
