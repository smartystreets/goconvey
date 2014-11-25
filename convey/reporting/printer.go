package reporting

import (
	"fmt"
	"io"
	"strings"
)

type Printer struct {
	out          io.Writer
	indent       int
	needSpace    bool
	inExpression bool
	firstLine    bool
}

func (self *Printer) Suite(name string) {
	self.Statement(name)
	self.indent++
	self.inExpression = true
	self.needSpace = true
}

func (self *Printer) Exit() {
	self.indent--
	if self.indent < 0 {
		self.indent = 0
	}
	self.inExpression = false
}

func (self *Printer) Statement(values ...interface{}) {
	var formatted string
	if self.firstLine {
		formatted = self.format(values...)
		self.firstLine = false
	} else {
		formatted = newline + self.format(values...)
	}
	self.inExpression = false
	self.out.Write([]byte(formatted))
}

func (self *Printer) Expression(values ...interface{}) {
	formatted := fmt.Sprint(values...)
	if !self.inExpression {
		formatted = newline + self.format(formatted)
	} else if self.needSpace {
		self.needSpace = false
		formatted = " " + formatted
	}
	self.inExpression = true
	self.out.Write([]byte(formatted))
}

func (self *Printer) Insert(text string) {
	self.out.Write([]byte(text))
}

func (self *Printer) format(values ...interface{}) string {
	pfx := strings.Repeat(pad, self.indent)
	formatted := pfx + fmt.Sprint(values...)
	indented := strings.Replace(formatted, newline, newline+pfx+pad, -1)
	return strings.TrimRight(indented, space+newline)
}

func NewPrinter(out io.Writer) *Printer {
	return &Printer{out: out, firstLine: true}
}

const space = " "
const pad = space + space
