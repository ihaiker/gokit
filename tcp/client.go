package tcpKit

import (
    "github.com/ihaiker/gokit/commons/logs"
    "net"
)

type Client struct {
    Connect
}

func NewClient(config *Config, handler Handler, protocol Protocol) *Client {
    client := &Client{}
    client.config = config
    client.logger = logs.Logger("tcpKit")
    client.closeChan = make(chan struct{})
    client.sendChan = make(chan interface{}, config.PacketSendChanLimit)
    client.Handler = handler
    client.Protocol = protocol
    return client
}

func NewClientWith(config *Config, wrapper *simpleWrapper) *Client {
    client := &Client{}
    client.config = config
    client.logger = logs.Logger("tcpKit")
    client.closeChan = make(chan struct{})
    client.sendChan = make(chan interface{}, config.PacketSendChanLimit)
    client.Handler = wrapper.handler
    client.Protocol = wrapper.protocol
    return client
}

func (c *Client) Start(conn *net.TCPConn) {
    c.connect = conn
    started := make(chan interface{})
    go c.Do(func(connect *Connect) {

    }, func(connect *Connect) {
        started <- 1
    })
    <-started
}

func (s *Client) StartAt(addr string) error {
    if listener, err := net.Dial("tcp", addr); err != nil {
        return err
    } else {
        s.Start(listener.(*net.TCPConn))
        return nil
    }
}
