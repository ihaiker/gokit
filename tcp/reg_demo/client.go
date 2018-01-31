package main

import (
    "github.com/ihaiker/gokit/commons/logs"
    "github.com/ihaiker/gokit/tcp"
    "github.com/ihaiker/gokit/runtime/signal"
    "time"
    "github.com/ihaiker/gokit/tcp/reg_demo/msg"
)


func main() {
    pkg := &msg.Package2{}
    var config = &tcpKit.Config{
        PacketReceiveChanLimit: 10, PacketSendChanLimit: 10,
        AcceptTimeout:          100,
        IdleTime:               0,
    }
    protocol := tcpKit.NewTVProtocol()
    protocol.Reg(pkg)

    logs.SetAllLevel(logs.DEBUG)
    s := tcpKit.NewClient(config, &tcpKit.HandlerWrapper{}, protocol)
    go func() {
        if err := s.StartAt("127.0.0.1:6379"); err != nil {
            logs.Fatal("启动连接错误：", err)
        }
    }()
    time.Sleep(time.Second)
    s.Write(&msg.Package2{Msg: "IDLE"})
    defer s.Close()

    signalKit.Signal(func() {
        logs.Info("重新加载")
    })
}
