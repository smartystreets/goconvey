package execution

func (parent *scope) adopt(child *scope) {
	if parent.hasChild(child) {
		return
	}
	name := functionName(child.action)
	parent.birthOrder = append(parent.birthOrder, name)
	parent.children[name] = child
}
func (parent *scope) hasChild(child *scope) bool {
	for _, name := range parent.birthOrder {
		if name == functionName(child.action) {
			return true
		}
	}
	return false
}

func (self *scope) registerReset(action func()) {
	self.resets[functionId(action)] = action
}

func (self *scope) visited() bool {
	return self.panicked || self.child >= len(self.birthOrder)
}

func (parent *scope) visit() {
	defer parent.exit()
	parent.enter()
	parent.action()
	parent.visitChildren()
}
func (parent *scope) enter() {
	parent.reporter.Enter(parent.name)
}
func (parent *scope) visitChildren() {
	if len(parent.children) == 0 {
		parent.cleanup()
	} else {
		parent.visitChild()
	}
}
func (parent *scope) visitChild() {
	child := parent.children[parent.birthOrder[parent.child]]
	child.visit()
	if child.visited() {
		parent.cleanup()
		parent.child++
	}
}
func (parent *scope) cleanup() {
	for _, reset := range parent.resets {
		reset()
	}
}
func (parent *scope) exit() {
	if problem := recover(); problem != nil {
		parent.panicked = true
		parent.reporter.Report(NewErrorReport(problem))
	}
	parent.reporter.Exit()
}

func newScope(name string, action func(), reporter Reporter) *scope {
	self := scope{name: name, action: action}
	self.children = make(map[string]*scope)
	self.birthOrder = []string{}
	self.resets = make(map[uintptr]func())
	self.reporter = reporter
	return &self
}

type scope struct {
	name       string
	action     func()
	children   map[string]*scope
	birthOrder []string
	child      int
	resets     map[uintptr]func()
	panicked   bool
	reporter   Reporter
}
