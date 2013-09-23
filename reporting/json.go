package reporting

type assertion struct {
	file    string
	line    int
	failure string
	error   error
	skipped bool
}

type scope struct {
	completeOutput string // only on the parent scope?
	title          string
	file           string
	line           int
	assertions     []assertion
	children       []scope
	skipped        bool
}
