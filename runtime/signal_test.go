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

	err := lis.WaitTimeout(time.Second, func() {
		logs.Debug("关闭")
	})
	t.Log(err)
}
