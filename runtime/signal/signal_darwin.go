package signalKit

import (
	"github.com/ihaiker/gokit/logs"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// InitSignal register signals handler.
func Signal(fn ...func(...os.Signal)) {
	logs.Infof("%s pid:%d", filepath.Base(os.Args[0]), os.Getpid())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		logs.Info("获取信号", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			if fn != nil && len(fn) == 2 {
				fn[1](s)
			}
			return
		case syscall.SIGHUP:
			if fn != nil && len(fn) > 0 {
				fn[0](s)
			}
		default:
			return
		}
	}
}
