package future

import (
    "github.com/ihaiker/gokit/concurrent/atomic"
)

type SetFuture struct {
	futureWapper
}

func (self *SetFuture) Set(result interface{}) {
    if !self.IsDone() {
        if self.status.CompareAndSet(_RUNNING, _OVER) {
            self.result = result
        }
        close(self.resultChan)
    }
}
func (self *SetFuture) Exception(err error) {
    if !self.IsDone() {
        if self.status.CompareAndSet(_RUNNING, _EXCEPTION) {
            self.err = err
        }
        close(self.resultChan)
    }
}

func Set() *SetFuture {
    f := &SetFuture{}
    f.resultChan = make(chan interface{}, 0)
    f.status = atomic.NewAtomicInt32(_RUNNING)
    return f
}
