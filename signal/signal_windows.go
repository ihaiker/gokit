package signalkit

import (
	log "code.google.com/p/log4go"
	"os"
	"os/signal"
	"syscall"
)

// InitSignal register signals handler.
func InitSignal(reload func()) {

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("process get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
			//TODO reload
			if reload != nil {
				reload()
			}
		default:
			return
		}
	}
}
