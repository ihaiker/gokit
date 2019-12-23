package runtimeKit

import (
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/logs"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger = logs.GetLogger("signal")

type SignalListener struct {
	C          chan os.Signal
	OnCloseFns []func()
}

func NewListener() *SignalListener {
	lis := &SignalListener{
		C:          make(chan os.Signal),
		OnCloseFns: []func(){},
	}
	signal.Notify(lis.C, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, )
	return lis
}

//正常退出
func (sl *SignalListener) Stop() {
	logger.Debug("send stop signal: ", syscall.SIGTERM.String())
	sl.C <- syscall.SIGTERM
}

//强制退出
func (sl *SignalListener) Kill() {
	logger.Info("kill self")
	os.Exit(0)
}

//关闭程序，首先使用Stop正常退出，然后使用Kill直接退出程序
func (sl *SignalListener) Shutdown(timeout time.Duration) {
	sl.Stop()
	<-time.After(timeout)
	sl.Kill()
}

//从最后添加一个
func (sl *SignalListener) OnClose(fn func()) {
	sl.OnCloseFns = append(sl.OnCloseFns, fn)
}

func (sl *SignalListener) PrependOnClose(fn func()) {
	sl.OnCloseFns = append([]func(){fn}, sl.OnCloseFns...)
}

//等待程序退出,如果close函数阻塞也将无法退出
func (sl *SignalListener) WaitWith(close func()) error {
	return sl.WaitWithTimeout(time.Hour, close)
}

//无调用的方式也可以退出
func (sl *SignalListener) WaitWithTimeout(timeout time.Duration, closeFn func()) error {
	for s := range sl.C {
		logs.Info("signal: ", s.String())
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			err := commons.AsyncTimeout(timeout, func() interface{} {
				logger.Debug("close by signal")
				for _, onCloseFn := range sl.OnCloseFns {
					onCloseFn()
				}
				if closeFn != nil {
					closeFn()
				}
				return nil
			})
			if err == commons.ErrAsyncTimeout {
				logger.Info("close timeout: ", timeout.String())
				sl.Kill()
				return commons.ErrAsyncTimeout
			}
			return nil
		}
	}
	return nil
}

func (sl *SignalListener) Wait() error {
	return sl.WaitTimeout(time.Hour)
}

func (sl *SignalListener) WaitTimeout(timeout time.Duration) error {
	return sl.WaitWithTimeout(timeout, func() {})
}
