package signalKit

import (
	log "code.google.com/p/log4go"
	"os"
	"os/signal"
	"syscall"
)

// InitSignal register signals handler.
func InitSignal(reload func()) {

	c := make(chan os.Signal, 1)

	//linux
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		log.Info("process get a signal %s", s.String())
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
