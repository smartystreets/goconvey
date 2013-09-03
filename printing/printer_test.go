package printing

import "testing"

func TestPrint(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const expected = "Hello, World!"

	printer.Print(expected)

	if file.buffer != expected {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestPrintln(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const expected = "Hello, World!"

	printer.Println(expected)

	if file.buffer != expected+"\n" {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestPrintIndented(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const message = "Hello, World!\nGoodbye, World!"
	const expected = "  Hello, World!\n  Goodbye, World!"

	printer.Indent()
	printer.Print(message)

	if file.buffer != expected {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestPrintDedented(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const expected = "Hello, World!\nGoodbye, World!"

	printer.Indent()
	printer.Dedent()
	printer.Print(expected)

	if file.buffer != expected {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestPrintlnIndented(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const message = "Hello, World!\nGoodbye, World!"
	const expected = "  Hello, World!\n  Goodbye, World!\n"

	printer.Indent()
	printer.Println(message)

	if file.buffer != expected {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestPrintlnDedented(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)
	const expected = "Hello, World!\nGoodbye, World!"

	printer.Indent()
	printer.Dedent()
	printer.Println(expected)

	if file.buffer != expected+"\n" {
		t.Errorf("Expected '%s' to equal '%s'.", expected, file.buffer)
	}
}

func TestDedentTooFarShouldNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Should not have panicked!")
		}
	}()
	file := newMemoryFile()
	printer := NewPrinter(file)

	printer.Dedent()

	t.Log("Getting to this point without panicking means we passed.")
}

func TestInsert(t *testing.T) {
	file := newMemoryFile()
	printer := NewPrinter(file)

	printer.Indent()
	printer.Print("Hi")
	printer.Insert(" there")
	printer.Dedent()

	expected := "  Hi there"
	if file.buffer != expected {
		t.Errorf("Should have written '%s' but instead wrote '%s'.", expected, file.buffer)
	}
}

type memoryFile struct {
	buffer string
}

func (self *memoryFile) Write(p []byte) (n int, err error) {
	self.buffer += string(p)
	return len(p), nil
}

func newMemoryFile() *memoryFile {
	self := memoryFile{}
	return &self
}
