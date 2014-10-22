package system

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestShellIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping potentially long-running integration test...")
		return
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	_, filename, _, _ := runtime.Caller(0)
	directory := filepath.Join(filepath.Dir(filename), "..", "watch", "integration_testing", "sub")
	packageName := "github.com/smartystreets/goconvey/web/server/integration_testing/sub"

	shell := NewShell("go", "", true)
	output, err := shell.GoTest(directory, packageName, []string{"-short"})

	if !strings.Contains(output, "PASS\n") || !strings.Contains(output, "ok") {
		t.Errorf("Expected output that resembed tests passing but got this instead: [%s]", output)
	}
	if err != nil {
		t.Error("Test run resulted in the following error:", err)
	}
}

func TestShellHangIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping potentially long-running integration test...")
		return
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	_, filename, _, _ := runtime.Caller(0)
	directory := filepath.Join(filepath.Dir(filename), "integration_testing", "hang")
	packageName := "github.com/smartystreets/goconvey/web/server/system/integration_testing/hang"

	shell := NewShell("go", "", true)

	const waitSeconds = 3
	startTime := time.Now()
	go func() {
		// wait for 5 seconds, after that kill the process
		timer := time.NewTimer(waitSeconds * time.Second)
		<-timer.C
		log.Print("abort go test!")
		err := shell.AbortGoTest()
		if err != nil {
			t.Error("Test run abort resulted in error:", err)
		}
	}()

	log.Print("start go test!")
	//this method will hang...
	output, err := shell.GoTest(directory, packageName, []string{})
	stopTime := time.Now()
	log.Print("go test returned")
	// confirm duration is between 5 and 6 seconds
	dur := stopTime.Sub(startTime)
	if dur < (waitSeconds*time.Second) || dur > ((waitSeconds+1)*time.Second) {
		t.Errorf("Expected %v second delay, was %s instead", waitSeconds, dur)
	}

	if !strings.Contains(output, "Hanging...") {
		t.Errorf("Unexpected output from test run: [%s]", output)
	}

	if err == nil || err.Error() != "exit status 2" {
		t.Errorf("Test run should have resulted in exit code error, was [%s] instead", err)
	}
}
