package future

import (
    "github.com/ihaiker/gokit/commons"
    "github.com/ihaiker/gokit/commons/atomic"
)

type AsyncRunFuture struct {
    futureWapper
}

func (self *AsyncRunFuture) Run(fn func(Future) (interface{}, error)) {
    go func() {
        defer close(self.resultChan)
        defer func() {
            if e := recover(); e != nil {
                if self.status.CompareAndSet(_RUNNING, _EXCEPTION) {
                    self.err = commonKit.Catch(e)
                }
            }
        }()
        if r, e := fn(self); e != nil {
            if !self.IsCancelled() {
                self.status.CompareAndSet(_RUNNING, _EXCEPTION)
                self.err = e
            }
        } else {
            if !self.IsCancelled() {
                self.status.CompareAndSet(_RUNNING, _OVER)
                self.result = r
            }
        }
    }()
    self.status.CompareAndSet(_INIT, _RUNNING)
}

func Run(fn func(Future) (interface{}, error)) Future {
    f := &AsyncRunFuture{}
    f.resultChan = make(chan interface{}, 0)
    f.status = atomic.NewInt32V(_INIT)
    f.Run(fn)
    return f
}
