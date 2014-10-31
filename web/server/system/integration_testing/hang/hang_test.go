// This file's only purpose is to provide a realistic
// environment from which to run hanging tests
// against the system commands to kill.
package hang

import (
	"fmt"
	"testing"
)

func TestStuff(t *testing.T) {
	if testing.Short() {
		return
	}

	fmt.Println("Hanging...if this isn't killed shortly, there's a bug in GoConvey's command.kill")

	for {

	}

	blah := make(chan string)
	blah <- "hang!"

}
