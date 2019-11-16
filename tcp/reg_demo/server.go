package main

import (
	"fmt"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/runtime/signal"
	"github.com/ihaiker/gokit/tcp"
	"time"
)

func main() {
	var config = &tcpKit.Config{
		PacketReceiveChanLimit: 10, PacketSendChanLimit: 10,
		AcceptTimeout: time.Second,
		IdleTime:      0,
		IdleTimeout:   2,
	}

	logs.SetDebugMode(true)

	reg := tcpKit.NewRegister()
	//reg.Reg(&msg.Package2{})

	s := tcpKit.NewServerWith(config, reg)
	err := s.StartAt("127.0.0.1:6379")
	fmt.Println(err)

	defer s.Stop()

	signalKit.Signal(func() {
		logs.Info("重新加载")
	})
}
