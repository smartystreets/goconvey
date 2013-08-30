package reporting

import (
	"time"
)

type Reporter struct {
	Successes int
	Failures  int
	Errors    int
	Duration  time.Duration
	Reports   map[string]*Report
}

type Report struct {
	Name       string
	SubReports []*Report
	Successes  int
	Failures   []error
	Errors     []error
	Duration   time.Duration
}
