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

type signalListener struct {
	C chan os.Signal
}

func NewListener() *signalListener {
	lis := &signalListener{C: make(chan os.Signal)}
	signal.Notify(lis.C, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	return lis
}

//正常退出
func (sl *signalListener) Stop() {
	logger.Debug("send stop signal: ", syscall.SIGUSR1.String())
	sl.C <- syscall.SIGUSR1
}

//强制退出
func (sl *signalListener) Kill() {
	logger.Info("kill self")
	os.Exit(0)
}

//关闭程序，首先使用Stop正常退出，然后使用Kill直接退出程序
func (sl *signalListener) Shutdown(timeout time.Duration) {
	sl.Stop()
	<-time.After(timeout)
	sl.Kill()
}

//等待程序退出,如果close函数阻塞也将无法退出
func (sl *signalListener) Wait(close func()) error {
	for s := range sl.C {
		logs.Info("signal: ", s.String())
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1:
			close()
			return nil
		}
	}
}

//无调用的方式也可以退出
func (sl *signalListener) WaitTimeout(timeout time.Duration, close func()) error {
	for s := range sl.C {
		logs.Info("signal: ", s.String())
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1:
			out := commons.AsyncTimeout(timeout, func() interface{} {
				logger.Debug("close self")
				close()
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
}
