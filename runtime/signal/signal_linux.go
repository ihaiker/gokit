package signalKit

import (
	"os"
	"os/signal"
	"syscall"
    "github.com/ihaiker/gokit/commons/logs"
)

// InitSignal register signals handler.
func Signal(reload func()) {
    logs.Infof("%s pid:%d", filepath.Base(os.Args[0]), os.Getpid())
	c := make(chan os.Signal, 1)
	//linux
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		logs.Info("获取信号 %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			return
		case syscall.SIGHUP:
			if reload != nil {
				reload()
			}
		default:
			return
		}
	}
}
