package commons

import (
	"errors"
	"time"
)

var ErrAsyncTimeout = errors.New("async timeout")

type AsyncFun func() interface{}

func Async(f AsyncFun) chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer func() { _ = recover() }()
		ch <- f()
	}()
	return ch
}

func AsyncTimeout(timeout time.Duration, f AsyncFun) (result interface{}) {
	ch := Async(f)
	select {
	case result = <-ch:
		close(ch)
	case <-time.After(timeout):
		close(ch)
		result = ErrAsyncTimeout
	}
	return result
}
