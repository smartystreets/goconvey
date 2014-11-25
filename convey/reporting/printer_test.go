package reporting

import "testing"

func TestExpression(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const expected = "\nHello, World!"

	printer.Expression(expected[1:])

	if file.buffer != expected {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestExpressionPreservesEncodedStrings(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const expected = "\n= -> %3D"
	printer.Expression(expected[1:])

	if file.buffer != expected {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestStatement(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const expected = "Hello, World!"

	printer.Statement(expected)

	if file.buffer != expected {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestStatementPreservesEncodedStrings(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const expected = "= -> %3D"
	printer.Statement(expected)

	if file.buffer != expected {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestSuiteExpression(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const expected = "Hello, World! +\n  Goodbye, World!"

	printer.Suite("Hello, World!")
	printer.Expression("+")
	printer.Statement("Goodbye, World!")

	if file.buffer != expected {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestExpressionDedented(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const expected = "Hello, World!\nGoodbye, World!"

	printer.Suite("Hello, World!")
	printer.Exit()
	printer.Expression("Goodbye, World!")

	if file.buffer != expected {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestStatementIndented(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const expected = "Hello, World!\n  line\n    Goodbye, World!"

	printer.Suite("Hello, World!")
	printer.Statement("line\nGoodbye, World!")

	if file.buffer != expected {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestStatementDedented(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const expected = "Hello, World!\nGoodbye, World!"

	printer.Suite("Hello, World!")
	printer.Exit()
	printer.Statement("Goodbye, World!")

	if file.buffer != expected {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestExitTooFarShouldNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Should not have panicked!")
		}
	}()
	file := newMemoryFile()
	printer := NewPrinter(file)

	printer.Exit()

	t.Log("Getting to this point without panicking means we passed.")
}

func TestInsert(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)

	printer.Suite("Sup")
	printer.Expression("Hi")
	printer.Insert(" there")
	printer.Exit()

	expected := "Sup Hi there"
	if file.buffer != expected {
		t.Errorf("Should have written '%s' but instead wrote '%s'.", expected, file.buffer)
	}
}

func TestFullSuite(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)

	printer.Suite("Sup")
	printer.Expression("+")
	printer.Expression("+")
	printer.Statement("some\ntext\nwith\nbreaks")
	printer.Expression("X")
	printer.Suite("Other")
	printer.Statement("other")
	printer.Expression("+")
	printer.Exit()
	printer.Expression("+")
	printer.Exit()

	expected := (`Sup ++
  some
    text
    with
    breaks
  X
  Other
    other
    +
  +`)
	if file.buffer != expected {
		t.Errorf("Should have written '%s' but instead wrote '%s'.", expected, file.buffer)
	}
}

////////////////// memoryFile ////////////////////

type memoryFile struct {
	buffer string
}

func (self *memoryFile) Write(p []byte) (n int, err error) {
	self.buffer += string(p)
	return len(p), nil
}

func (self *memoryFile) String() string {
	return self.buffer
}

func newMemoryFile() *memoryFile {
	return new(memoryFile)
}
