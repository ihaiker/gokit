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
		ch <- f()
	}()
	return ch
}

func AsyncTimeout(timeout time.Duration, f AsyncFun) chan interface{} {
	ch := make(chan interface{})
	go func() {
		select {
		case err := <-Async(f):
			ch <- err
		case <-time.After(timeout):
			ch <- ErrAsyncTimeout
		}
	}()
	return ch
}
