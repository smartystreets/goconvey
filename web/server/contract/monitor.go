package contract

type Monitor struct {
	// Scan()   // one round of scanning and test execution
	// Engage() // infinite for loop, calls Tick() between time.Sleep() (when no tests were run)
}

func NewMonitor(scanner Scanner, watcher Watcher, executor Executor, server Server) *Monitor {
	return nil
}
