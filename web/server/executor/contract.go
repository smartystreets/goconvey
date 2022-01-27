package executor

import "github.com/luchuanbing123/goconvey/web/server/contract"

type Parser interface {
	Parse([]*contract.Package)
}

type Tester interface {
	SetBatchSize(batchSize int)
	TestAll(folders []*contract.Package)
}
