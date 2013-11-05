package contract

type Monitor struct {
	scanner  Scanner
	watcher  Watcher
	executor Executor
	server   Server
	sleep    func()
}

func (self *Monitor) ScanForever() {
	for {
		self.Scan()
	}
}

func (self *Monitor) Scan() {
	root := self.watcher.Root()

	if self.scanner.Scan(root) {
		self.executeTests()
	} else {
		self.sleep()
	}
}

func (self *Monitor) executeTests() {
	watched := self.watcher.WatchedFolders()
	output := self.executor.ExecuteTests(watched)
	self.server.ReceiveUpdate(output)
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
