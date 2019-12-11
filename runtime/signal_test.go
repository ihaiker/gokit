package runtimeKit

import (
	"github.com/ihaiker/gokit/logs"
	"testing"
	"time"
)

func TestSignal(t *testing.T) {
	logs.SetDebugMode(true)
	lis := NewListener()
	logs.Info("start...")

	go func() {
		time.Sleep(time.Second)
		lis.Stop()
	}()

	lis.OnClose(func() {
		logger.Debug("关闭操作")
	})
	lis.OnClose(func() {
		logger.Debug("该你操作！！！")
	})
	err := lis.WaitWithTimeout(time.Second, func() {
		logs.Debug("关闭")
	})
	t.Log(err)
}
