package remoting

import (
	"net"
	"time"
)

type Client interface {
	Start() Client
	Close()
	Send(msg interface{}) error
	GetChannel() Channel
	Wait()
}

type tcpClient struct {
	channel *tcpChannel
}

func (self *tcpClient) Start() Client {
	go self.channel.do(nil, nil)
	return self
}

func (self *tcpClient) Close() {
	self.channel.Close()
}

func (self *tcpClient) Send(msg interface{}) error {
	return self.channel.Write(msg)
}

func (self *tcpClient) SyncSend(msg interface{}, timeout time.Duration, cb SendMessageCallBack) {
	self.channel.AsyncWrite(msg, timeout, cb)
}

func (self *tcpClient) GetChannel() Channel {
	return self.channel
}

func (self *tcpClient) Wait() {
	self.channel.wait()
}

func NewClientWith(connect net.Conn, config *Config, handler Handler, coder Coder) Client {
	channel := newChannel(config, connect)
	channel.handler = handler
	channel.coder = coder
	return &tcpClient{channel: channel}
}

func NewClient(address string, config *Config, handler Handler, coder Coder) (Client, error) {
	if conn, err := Dial(address); err != nil {
		return nil, err
	} else {
		return NewClientWith(conn, config, handler, coder), nil
	}
}
