package reporting

import (
	"fmt"
	"io"
	"strings"
)

type Printer struct {
	out    io.Writer
	prefix string
}

func (p *Printer) Println(message string, values ...interface{}) {
	formatted := p.format(message, values...) + newline
	p.out.Write([]byte(formatted))
}

func (p *Printer) Print(message string, values ...interface{}) {
	formatted := p.format(message, values...)
	p.out.Write([]byte(formatted))
}

func (p *Printer) Insert(text string) {
	p.out.Write([]byte(text))
}

func (p *Printer) format(message string, values ...interface{}) string {
	var formatted string
	if len(values) == 0 {
		formatted = p.prefix + message
	} else {
		formatted = p.prefix + fmt_Sprintf(message, values...)
	}
	indented := strings.Replace(formatted, newline, newline+p.prefix, -1)
	return strings.TrimRight(indented, space)
}

// Extracting fmt.Sprintf to a separate variable circumvents go vet, which, as of go 1.10 is run with go test.
var fmt_Sprintf = fmt.Sprintf

func (p *Printer) Indent() {
	p.prefix += pad
}

func (p *Printer) Dedent() {
	if len(p.prefix) >= padLength {
		p.prefix = p.prefix[:len(p.prefix)-padLength]
	}
}

func NewPrinter(out io.Writer) *Printer {
	self := new(Printer)
	self.out = out
	return self
}

const space = " "
const pad = space + space
const padLength = len(pad)
