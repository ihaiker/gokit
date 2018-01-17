/*
    事件处理器
 */
package tcpKit

import (
    "net"
    "github.com/ihaiker/gokit/commons/logs"
)

type Handler interface {
    //链接事件，当客户端链接时调用
    OnConnect(c *Connect)
    //新消息事件，当客户端发来新的消息时调用
    OnMessage(c *Connect, msg interface{})
    //编码异常事件，当编码消息时
    OnEncodeError(c *Connect, msg interface{}, err error)
    //处理错误事件，当OnMessage抛出未能处理的错误
    OnError(c *Connect, err error, msg interface{})
    //解码异常事件，解码时错误
    OnDecodeError(c *Connect, err error)
    //发送心跳包
    OnIdle(c *Connect)
    //关闭事件，当当前客户端关闭连接
    OnClose(c *Connect)
}

type HandlerWrapper struct {
}

func (h *HandlerWrapper) OnConnect(c *Connect) {
    c.connect.SetKeepAlive(true)
    c.connect.SetNoDelay(true)
    logs.Debugf("Handler OnConnect %s", c.connect.RemoteAddr().String())
}
func (h *HandlerWrapper) OnMessage(c *Connect, msg interface{}) {
    logs.Debugf("Handler OnMessage %s : msg:%s", c.connect.RemoteAddr().String(), msg)
}
func (h *HandlerWrapper) OnClose(c *Connect) {
    logs.Debugf("Handler OnClose %s ", c.connect.RemoteAddr().String())
}
func (h *HandlerWrapper) OnError(c *Connect, err error, msg interface{}) {
    logs.Debugf("Handler OnError %s : %s ,%s", c.connect.RemoteAddr().String(), err, msg)
}
func (h *HandlerWrapper) OnEncodeError(c *Connect, msg interface{}, err error) {
    logs.Debugf("Handler OnError %s : %s", c.connect.RemoteAddr().String(), err)
}
func (h *HandlerWrapper) OnDecodeError(c *Connect, err error) {
    logs.Debugf("Handler OnError %s : %s", c.connect.RemoteAddr().String(), err)
}
func (h *HandlerWrapper) OnIdle(c *Connect) {
    logs.Debugf("Handler OnIdle : %s", c.connect.RemoteAddr().String())
}

type HandlerMaker func(c *net.TCPConn) Handler
