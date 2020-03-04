package future

import (
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/concurrent/atomic"
)

type AsyncRunFuture struct {
	futureWapper
}

func (self *AsyncRunFuture) Run(fn func(Future) (interface{}, error)) {
	self.status.CompareAndSet(_INIT, _RUNNING)
	go func() {
		defer close(self.resultChan)
		defer commons.Catch(func(err error) {
			if !self.IsCancelled() {
				if self.status.CompareAndSet(_RUNNING, _EXCEPTION) {
					self.err = err
				}
			}
		})
		if result, err := fn(self); err != nil {
			if !self.IsCancelled() {
				self.status.CompareAndSet(_RUNNING, _EXCEPTION)
				self.err = err
			}
		} else {
			if !self.IsCancelled() {
				self.status.CompareAndSet(_RUNNING, _OVER)
				self.result = result
			}
		}
	}()
}

func Run(fn func(Future) (interface{}, error)) Future {
	f := &AsyncRunFuture{}
	f.resultChan = make(chan interface{}, 0)
	f.status = atomic.NewAtomicInt32(_INIT)
	f.Run(fn)
	return f
}
