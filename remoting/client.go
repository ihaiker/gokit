package remoting

import (
	"github.com/ihaiker/gokit/concurrent/executors"
	"net"
	"time"
)

type Client interface {
	Start() error
	Close() error
	Write(msg interface{}, timeout time.Duration) error
	AsyncWrite(msg interface{}, timeout time.Duration, callback SendMessageResult)
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
	worker  *executors.GrPool
}

func (self *tcpClient) Start() (err error) {
	if self.connect, err = Dial(self.address); err != nil {
		return
	}

	if self.config.AsyncHandlerGroup != 0 {
		self.worker = executors.NewPoolDefault(self.config.AsyncHandlerGroup)
	}
	self.channel = newChannel(self.config, self.worker, self.connect)

	self.channel.handler = self.handler
	self.channel.coder = self.coder

	c := make(chan interface{})
	go self.channel.do(func(channel Channel) { c <- channel }, nil)
	<-c
	close(c)
	return nil
}

func (self *tcpClient) Close() error {
	self.worker.Shutdown()
	return self.channel.Close()
}

func (self *tcpClient) Write(msg interface{}, timeout time.Duration) error {
	return self.channel.Write(msg, timeout)
}

func (self *tcpClient) AsyncWrite(msg interface{}, timeout time.Duration, cb SendMessageResult) {
	self.channel.AsyncWrite(msg, timeout, cb)
}

func (self *tcpClient) GetChannel() Channel {
	return self.channel
}

func (self *tcpClient) Wait() {
	self.channel.Wait()
}

func NewClient(address string, config *Config, handler Handler, coder Coder) Client {
	return &tcpClient{
		address: address,
		coder:   coder,
		handler: handler,
		config:  config,
	}
}
