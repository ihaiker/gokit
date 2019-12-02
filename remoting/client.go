package remoting

import (
	"net"
	"time"
)

type Client interface {
	Start() error
	Close()
	Send(msg interface{}, timeout time.Duration) error
	Async(msg interface{}, timeout time.Duration, callback SendMessageCallBack)
	GetChannel() Channel
	Wait()
}

type tcpClient struct {
	address string   //监听地址
	connect net.Conn //连接器
	coder   Coder    //消息编码
	handler Handler  //消息处理器
	config  *Config  //配置管理器
	channel *tcpChannel
}

func (self *tcpClient) Start() (err error) {
	if self.connect, err = Dial(self.address); err != nil {
		return
	}

	self.channel = newChannel(self.config, self.connect)
	self.channel.handler = self.handler
	self.channel.coder = self.coder

	c := make(chan interface{})
	go self.channel.do(func(channel Channel) { c <- channel }, nil)
	<-c
	close(c)
	return nil
}

func (self *tcpClient) Close() {
	self.channel.Close()
}

func (self *tcpClient) Send(msg interface{}, timeout time.Duration) error {
	return self.channel.Write(msg, timeout)
}

func (self *tcpClient) Async(msg interface{}, timeout time.Duration, cb SendMessageCallBack) {
	self.channel.AsyncWrite(msg, timeout, cb)
}

func (self *tcpClient) GetChannel() Channel {
	return self.channel
}

func (self *tcpClient) Wait() {
	self.channel.wait()
}

func NewClient(address string, config *Config, handler Handler, coder Coder) Client {
	return &tcpClient{
		address: address,
		coder:   coder,
		handler: handler,
		config:  config,
	}
}
