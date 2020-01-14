package remoting

import (
	"github.com/ihaiker/gokit/concurrent/executors"
	"net"
	"time"
)

type Client interface {
	Start() error
	Close() error
	Write(msg Message, timeout time.Duration) error
	AsyncWrite(msg Message, timeout time.Duration, callback SendMessageResult)
	GetChannel() Channel
	Wait()
}

type tcpClient struct {
	address string   //监听地址
	connect net.Conn //连接器
	coder   Coder    //消息编码
	handler Handler  //消息处理器
	options *Options //配置管理器
	*tcpChannel
	worker executors.ExecutorService
}

func (self *tcpClient) Start() (err error) {
	if self.connect, err = Dial(self.address); err != nil {
		return
	}

	if self.options.WorkerGroup != 0 {
		self.worker = executors.Fixed(self.options.WorkerGroup)
	}

	self.tcpChannel = newChannel(self.options, self.worker, self.connect)

	self.tcpChannel.handler = self.handler
	self.tcpChannel.coder = self.coder

	c := make(chan interface{})
	go self.tcpChannel.do(func(channel Channel) { c <- channel })
	<-c
	close(c)
	return nil
}

func (self *tcpClient) Close() error {
	err := self.tcpChannel.Close()
	if self.worker != nil {
		self.worker.Shutdown()
	}
	return err
}

func (self *tcpClient) GetChannel() Channel {
	return self.tcpChannel
}

func NewClient(address string, options *Options, handler Handler, coder Coder) Client {
	return &tcpClient{
		address: address, options: options,
		coder: coder, handler: handler,
	}
}
