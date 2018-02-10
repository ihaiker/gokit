package main

import (
    "github.com/ihaiker/gokit/tcp"
    "github.com/ihaiker/gokit/commons/logs"
    "github.com/ihaiker/gokit/runtime/signal"
    "time"
    "github.com/ihaiker/gokit/tcp/reg_demo/msg"
    "fmt"
)

func main() {
    var config = &tcpKit.Config{
        PacketReceiveChanLimit: 10, PacketSendChanLimit: 10,
        AcceptTimeout:          time.Second,
        IdleTime:               0,
        IdleTimeout:            2,
    }

    logs.SetAllLevel(logs.DEBUG)
    reg := tcpKit.NewRegister()
    reg.Reg(&msg.Package2{})

    s := tcpKit.NewServerWith(config, reg)
    err := s.StartAt("127.0.0.1:6379")
    fmt.Println(err)

    defer s.Stop()

    signalKit.Signal(func() {
        logs.Info("重新加载")
    })
}
