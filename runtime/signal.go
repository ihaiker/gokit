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
	onCloseFns []func()
}

func NewListener() *SignalListener {
	lis := &SignalListener{
		C:          make(chan os.Signal),
		onCloseFns: []func(){},
	}
	signal.Notify(lis.C, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	return lis
}

//正常退出
func (sl *SignalListener) Stop() {
	logger.Debug("send stop signal: ", syscall.SIGUSR1.String())
	sl.C <- syscall.SIGUSR1
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

func (sl *SignalListener) OnClose(fn func()) {
	sl.onCloseFns = append(sl.onCloseFns, fn)
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
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1:
			out := commons.AsyncTimeout(timeout, func() interface{} {
				logger.Debug("close self")
				for _, onCloseFn := range sl.onCloseFns {
					onCloseFn()
				}
				if closeFn != nil {
					closeFn()
				}
				return nil
			})
			if e := <-out; e == commons.ErrAsyncTimeout {
				logger.Info("close timeout: ", timeout.String())
				sl.Kill()
				return commons.ErrAsyncTimeout
			} else {
				return nil
			}
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
