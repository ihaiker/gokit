package future

import (
	"errors"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"time"
)

const (
    _INIT      = int32(0)
    _RUNNING   = iota
    _CANCEL
    _EXCEPTION
    _OVER
)

type Future interface {
    Cancel() bool

    IsCancelled() bool

    IsDone() bool

    Get() (interface{}, error)

    GetWithTimeout(timeout time.Duration) (interface{}, error)
}

type futureWapper struct {
    result     interface{}
    err        error
    resultChan chan interface{}
    status     *atomic.AtomicInt32
}

func (self *futureWapper) Cancel() bool {
    if self.status.CompareAndSet(_RUNNING, _CANCEL) {
        self.err = errors.New("cancel")
        return true
    }
    return false
}

func (self *futureWapper) IsCancelled() bool {
    return self.status.Get() == _CANCEL
}

func (self *futureWapper) IsDone() bool {
    return self.status.Get() == _EXCEPTION || self.status.Get() == _OVER || self.status.Get() == _CANCEL
}

func (self *futureWapper) Get() (interface{}, error) {
    if self.IsDone() {
        return self.result, self.err
    }
    <-self.resultChan
    return self.result, self.err
}

func (self *futureWapper) GetWithTimeout(timeout time.Duration) (interface{}, error) {
    if self.IsDone() {
        return self.result, self.err
    }
    select {
    case <-self.resultChan:
        return self.result, self.err
    case <-time.After(timeout):
        return nil, errors.New("timeout")
    }
}
