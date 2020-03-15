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

type (
	Service interface {
		Start() error
		Stop() error
	}

	SignalListener struct {
		C        chan os.Signal
		services []Service
		idx      int
	}
)

func NewListener() *SignalListener {
	lis := &SignalListener{
		C:        make(chan os.Signal),
		services: make([]Service, 0), idx: -1,
	}
	signal.Notify(lis.C, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, )
	return lis
}

//正常退出
func (sl *SignalListener) Shutdown() {
	logger.Debug("send stop signal: ", syscall.SIGTERM.String())
	sl.C <- syscall.SIGTERM
}

//强制退出
func (sl *SignalListener) Kill() {
	logger.Info("kill self")
	os.Exit(0)
}

func (sl *SignalListener) Add(services ...Service) *SignalListener {
	sl.services = append(sl.services, services...)
	return sl
}

//从最后添加一个
func (sl *SignalListener) AddStart(fn func() error) *SignalListener {
	sl.services = append(sl.services, &funcService{StartFn: fn})
	return sl
}

//从最后添加一个
func (sl *SignalListener) AddStop(fn func() error) *SignalListener {
	sl.services = append(sl.services, &funcService{StopFn: fn})
	return sl
}

func (sl *SignalListener) stop(timeout time.Duration) error {
	err := commons.AsyncTimeout(timeout, func() interface{} {
		for i := sl.idx; i >= 0; i-- {
			if err := sl.services[i].Stop(); err != nil {
				logger.Warn("stop service error: ", err)
			}
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

func (sl *SignalListener) start() error {
	for i, service := range sl.services {
		if err := service.Start(); err != nil {
			return err
		}
		sl.idx = i
	}
	return nil
}

func (sl *SignalListener) await(timeout time.Duration) error {
	for s := range sl.C {
		logger.Info("signal: ", s.String())
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			return sl.stop(timeout)
		}
	}
	return nil
}

func (sl *SignalListener) Wait() error {
	return sl.WaitTimeout(time.Second * 7)
}

func (sl *SignalListener) WaitTimeout(timeout time.Duration) error {
	if err := sl.start(); err != nil {
		_ = sl.stop(timeout)
		return err
	}
	return sl.await(timeout)
}

func Wait() error {
	return NewListener().WaitTimeout(time.Second * 7)
}

func WaitTimeout(timeout time.Duration) error {
	return NewListener().WaitTimeout(timeout)
}

func WaitTC(timeout time.Duration, fn func() error) error {
	return NewListener().AddStop(fn).WaitTimeout(timeout)
}
