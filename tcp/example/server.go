package main

import (
    "net"
    "github.com/ihaiker/gokit/tcp"
    "github.com/ihaiker/gokit/commons/logs"
    "github.com/ihaiker/gokit/runtime/signal"
    "time"
)

var config = &tcpKit.Config{
    PacketReceiveChanLimit: 10, PacketSendChanLimit: 10,
    AcceptTimeout:          3,
    IdleTime:               1000,
}

type TestHandlerWrapper struct {
    tcpKit.HandlerWrapper
}

func (h *TestHandlerWrapper) OnMessage(c *tcpKit.Connect, msg interface{}) {
    newMsg := time.Now().String()
    logs.Debugf("新消息：%s，回复：%s", msg, newMsg)
    if err := c.AsyncWrite(newMsg, time.Second); err != nil {
        logs.Info("发送异常：",err)
    }
}

func (h *TestHandlerWrapper) OnConnect(c *tcpKit.Connect) {
    h.HandlerWrapper.OnConnect(c)
    logs.Info(c.Write("Server 测试结果是什么？"))
}

func (h *TestHandlerWrapper) OnIdle(c *tcpKit.Connect) {
    h.HandlerWrapper.OnIdle(c)
    c.Write("IDLE")
}

func (h *TestHandlerWrapper) OnClose(c *tcpKit.Connect) {
    for {
        if msg := c.PopUnSend(time.Millisecond * 20); msg != nil {
            logs.Debugf("未发送消息：%s", msg)
        } else {
            return
        }
    }
}

func makeHandler(c *net.TCPConn) tcpKit.Handler {
    return &TestHandlerWrapper{}
}

func makeProtocol(c *net.TCPConn) tcpKit.Protocol {
    return &tcpKit.LineProtocol{Delim: "\r\n"}
}

func main() {
    logs.SetAllLevel(logs.DEBUG)
    s := tcpKit.NewServer(config, makeHandler, makeProtocol)

    go s.StartAt("127.0.0.1:6379")
    defer s.Stop()

    signalKit.Signal(func() {
        logs.Info("重新加载")
    })
}
