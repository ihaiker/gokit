package main

import (
    "github.com/ihaiker/gokit/commons/logs"
    "github.com/ihaiker/gokit/tcp"
    "github.com/ihaiker/gokit/runtime/signal"
)

type TestClientHandlerWrapper struct {
    tcpKit.HandlerWrapper
}

func (h *TestClientHandlerWrapper) OnConnect(c *tcpKit.Connect) {
    h.HandlerWrapper.OnConnect(c)
    logs.Info(c.Write("GGGGGG 测试结果是什么？"))
}

func main() {
    var config = &tcpKit.Config{
        PacketReceiveChanLimit: 10, PacketSendChanLimit: 10,
        AcceptTimeout:          100,
        IdleTime:               0,
    }
    logs.SetAllLevel(logs.DEBUG)
    s := tcpKit.NewClient(config, &TestClientHandlerWrapper{}, &tcpKit.LineProtocol{})

    go func() {
        if err := s.StartAt("127.0.0.1:6379"); err != nil {
            logs.Fatal("启动连接错误：", err)
        }
    }()
    defer s.Close()

    signalKit.Signal(func() {
        logs.Info("重新加载")
    })
}
