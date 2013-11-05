package contract

type Monitor struct {
	scanner  Scanner
	watcher  Watcher
	executor Executor
	server   Server
	sleep    func()
	// ScanForever() // infinite for loop, calls Scan() between time.Sleep() (when no tests were run)
}

func (self *Monitor) Scan() {
	root := self.watcher.WatchedFolders()[0].Path

	if self.scanner.Scan(root) {
		watched := self.watcher.WatchedFolders()
		output := self.executor.ExecuteTests(watched)
		self.server.ReceiveUpdate(output)
	} else {
		self.sleep()
	}
}

func NewMonitor(scanner Scanner, watcher Watcher, executor Executor, server Server, sleep func()) *Monitor {
	self := &Monitor{}
	self.scanner = scanner
	self.watcher = watcher
	self.executor = executor
	self.server = server
	self.sleep = sleep
	return self
}
