package main

import (
    "io"
    "net"
    "github.com/ihaiker/gokit/tcp"
    "github.com/ihaiker/gokit/commons/logs"
    "github.com/ihaiker/gokit/runtime/signal"
    "time"
    "github.com/ihaiker/gokit/tcp/reg_demo/msg"
)


func main() {
    var config = &tcpKit.Config{
        PacketReceiveChanLimit: 10, PacketSendChanLimit: 10,
        AcceptTimeout:          time.Second,
        IdleTime:    0,
        IdleTimeout: 2,
    }

    logs.SetAllLevel(logs.DEBUG)
    protocol := tcpKit.NewSimpleProtocol()
    protocol.Reg(&msg.Package2{})

    s := tcpKit.NewServer(config, func(c *net.TCPConn) tcpKit.Handler {
        return &tcpKit.HandlerWrapper{}
    }, func(c io.Reader) tcpKit.Protocol {
        return protocol
    })

    go s.StartAt("127.0.0.1:6379")
    defer s.Stop()

    signalKit.Signal(func() {
        logs.Info("重新加载")
    })
}
