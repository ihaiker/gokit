package signalKit

import (
	"os"
	"os/signal"
	"path/filepath"
	"github.com/ihaiker/gokit/logs"
	"syscall"
)

// InitSignal register signals handler.
func Signal(reload func()) {
	logs.Infof("%s pid:%d", filepath.Base(os.Args[0]), os.Getpid())
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logs.Info("获取信号 %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
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
