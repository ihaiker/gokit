package commons

import (
	"errors"
	"time"
)

var ErrAsyncTimeout = errors.New("async timeout")

type AsyncFun func() interface{}

//保证发送的时候不会出现: send on closed channel
func sendChannel(obj interface{}, c chan interface{}) (send error) {
	defer Catch(func(err error) {
		send = err
	})
	c <- obj
	return nil
}

//保证不会出现: closed 问题，这个恶心，没办法判断还补课已关闭
func closeChannel(c chan interface{}) {
	defer func() { _ = recover() }()
	close(c)
}

func Async(f AsyncFun) chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer func() {
			if err := recover(); err != nil {
				_ = sendChannel(err, ch)
			}
		}()
		_ = sendChannel(f(), ch)
	}()
	return ch
}

func AsyncTimeout(timeout time.Duration, f AsyncFun) (result interface{}) {
	ch := Async(f)
	select {
	case result = <-ch:
		closeChannel(ch)
	case <-time.After(timeout):
		closeChannel(ch)
		result = ErrAsyncTimeout
	}
	return result
}
