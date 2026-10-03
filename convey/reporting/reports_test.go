package reporting

import "testing"

func TestFailureReportParsesSerializedFailureViews(t *testing.T) {
	report := NewFailureReport(`{"Message":"Values differ.","Expected":"1","Actual":"2"}`, false)

	if report.Failure != "Values differ." {
		t.Errorf("Expected failure message 'Values differ.' but got '%s'.", report.Failure)
	}
	if report.Expected != "1" {
		t.Errorf("Expected expected value '1' but got '%s'.", report.Expected)
	}
	if report.Actual != "2" {
		t.Errorf("Expected actual value '2' but got '%s'.", report.Actual)
	}
}

// A custom assertion (convey.Assertion) is free to return any failure message
// it likes, including one that happens to be valid JSON. Only messages that
// actually carry a 'Message' field are FailureViews; anything else must be
// preserved verbatim so that it is not mistaken for a passing assertion.
func TestFailureReportRetainsMessagesThatAreNotFailureViews(t *testing.T) {
	messages := []string{
		`{"Expected":"even","Actual":"odd"}`,
		`{}`,
		`null`,
		`"not a failure view"`,
		"Expected: '1'\nActual:   '2'\n(Should be equal)",
	}

	for _, message := range messages {
		report := NewFailureReport(message, false)
		if report.Failure == "" {
			t.Errorf("Expected failure message '%s' to be retained, but it was empty.", message)
		}
	}
}
