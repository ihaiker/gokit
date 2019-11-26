package signalKit

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"github.com/ihaiker/gokit/logs"
)

// InitSignal register signals handler.
//idx:0 reload func
//idx:1 closed func
func Signal(fn ...func(...os.Signal)) {
	logs.Infof("%s pid:%d", filepath.Base(os.Args[0]), os.Getpid())
	c := make(chan os.Signal, 1)
	//linux
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		logs.Info("获取信号 %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			if fn != nil && len(fn) == 2 {
				fn[1](s)
			}
			return
		case syscall.SIGHUP:
			if fn != nil && len(fn) > 0 {
				fn[0]()
			}
		default:
			return
		}
	}
}
