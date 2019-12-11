package remoting

import (
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

//发送消息回调
type SendMessageCallBack func(msg interface{}, err error)

type asyncMessage struct {
	msg      interface{}
	timeout  time.Time
	callback SendMessageCallBack
}

//连接保持器
type Channel interface {
	//同步发送消息
	Write(msg interface{}, timeout time.Duration) (err error)

	//异步发送消息
	AsyncWrite(msg interface{}, timeout time.Duration, sendCallback SendMessageCallBack)

	//返回连接器
	GetConnect() net.Conn

	GetRemoteAddress() string

	GetRemoteIp() string

	Close()

	//是否正在运行
	IsRunning() bool

	commons.Attributes
}

const (
	ready uint32 = iota
	starting
	stoping
	stoped
)

type tcpChannel struct {
	config  *Config
	connect net.Conn

	closeOne *sync.Once
	status   *atomic.AtomicUint32

	group *sync.WaitGroup

	closeChan chan struct{}
	sendChan  chan *asyncMessage

	coder   Coder
	handler Handler

	commons.Attributes

	idleTimer   *time.Timer
	idleTimeout *atomic.AtomicInt32
}

func newChannel(config *Config, connect net.Conn) *tcpChannel {
	if tcpCon, match := connect.(*net.TCPConn); match {
		_ = tcpCon.SetKeepAlive(true)
		_ = tcpCon.SetNoDelay(true)
		_ = tcpCon.SetWriteBuffer(config.WriteBufferSize)
	}

	return &tcpChannel{
		config: config, connect: connect,

		closeOne: new(sync.Once), status: atomic.NewAtomicUint32(ready),
		group: new(sync.WaitGroup),

		closeChan: make(chan struct{}),
		sendChan:  make(chan *asyncMessage, config.SendChanLimit),

		Attributes: commons.NewAttributes(),
	}
}

//connected 连接后回调
//closed 关闭后回调
func (self *tcpChannel) do(connected, closed func(channel Channel)) {
	defer func() {
		logger.Debug("channel运行结束")
		self.closeChannel()
		if closed != nil {
			closed(self)
		}
		self.handler.OnClose(self)
	}()
	self.status.Set(starting)
	self.group.Add(3)

	go self.syncDo(self.heartbeatLoop)
	go self.syncDo(self.readLoop)
	go self.syncDo(self.writeLoop)

	if connected != nil {
		connected(self)
	}
	self.handler.OnConnect(self)

	self.wait()
}

func (c *tcpChannel) syncDo(fn func()) {
	defer c.group.Done()
	fn()
}

func (self *tcpChannel) readLoop() {
	logger.Debug("reader channel start: ", self.GetRemoteAddress())
	defer func() {
		self.closeChannel()
		logger.Debug("reader channel close: ", self.GetRemoteAddress())
	}()

	for {
		select {
		case <-self.closeChan:
			return
		default:
			_ = self.connect.SetReadDeadline(time.Now().Add(time.Second))
			if msg, err := self.coder.Decode(self, self.connect); commons.NotNil(err) {
				if isCloseTCPConnect(err) {
					return
				}
				if !strings.Contains(err.Error(), "i/o timeout") {
					if self.IsRunning() {
						logger.Errorf("服务监听错误：%s", err)
						self.handler.OnDecodeError(self, err)
					}
				}
			} else {
				_ = self.connect.SetReadDeadline(time.Now().Add(time.Second))
				self.resetIdle()
				if self.config.AsynHandlerGroup > 0 {
					//fixme 管理携程
					go commons.Try(func() {
						self.handler.OnMessage(self, msg)
					}, func(err error) {
						self.handler.OnError(self, msg, err)
					})
				} else {
					commons.Try(func() {
						self.handler.OnMessage(self, msg)
					}, func(err error) {
						self.handler.OnError(self, msg, err)
					})
				}
			}
		}
	}
}

func (self *tcpChannel) writeLoop() {
	logger.Debug("启动 write 携程")
	defer func() {
		logger.Debugf("关闭 writer 携程: %s", self.GetRemoteAddress())
		self.closeChannel()
	}()
	for {
		select {
		case <-self.closeChan:
			return
		case asyncMsg := <-self.sendChan:
			if !self.IsRunning() {
				return
			}
			if time.Now().Before(asyncMsg.timeout) {
				if bs, err := self.coder.Encode(self, asyncMsg.msg); err != nil {
					self.handler.OnEncodeError(self, asyncMsg.msg, err)
				} else {
					_, err := self.connect.Write(bs)
					if asyncMsg.callback != nil {
						asyncMsg.callback(asyncMsg.msg, err)
					}
					if err != nil {
						self.handler.OnError(self, asyncMsg.msg, err)
					}
				}
			} else if asyncMsg.callback != nil {
				asyncMsg.callback(asyncMsg.msg, ErrWriteTimeout)
			}
		}
	}
}

func (self *tcpChannel) heartbeatLoop() {
	if self.config.IdleDuration == 0 {
		logger.Debug(self.GetRemoteAddress(), " 不进行心跳检测！")
		return
	}

	logger.Debug("启动心跳检测携程:", self.GetRemoteAddress())

	self.idleTimer = time.NewTimer(time.Second * time.Duration(self.config.IdleDuration))
	self.idleTimeout = atomic.NewAtomicInt32(0)

	defer func() {
		logger.Debugf("关闭心跳检测: %s", self.GetRemoteAddress())
		self.idleTimer.Stop()
		self.closeChannel()
	}()

	for {
		select {
		case <-self.closeChan:
			return
		case <-self.idleTimer.C:
			self.idleTimer.Reset(time.Second * time.Duration(self.config.IdleDuration))
			if self.idleTimeout.GetAndIncrement(1) >= int32(self.config.IdleTimeout) {
				logger.Debug("连接超时未检测到：", self.GetRemoteAddress())
				return
			}
			self.handler.OnIdle(self)
		}
	}
}

func (self *tcpChannel) resetIdle() {
	if self.idleTimer != nil {
		self.idleTimer.Reset(time.Second * time.Duration(self.config.IdleDuration))
		self.idleTimeout.Set(0)
	}
}

func (self *tcpChannel) wait() {
	self.group.Wait()
}

func isCloseTCPConnect(err error) bool {
	if strings.Contains(err.Error(), ErrConnectClosed.Error()) || strings.Contains(err.Error(), "connection reset by peer") {
		return true
	} else if err == io.EOF {
		return true
	}
	return false
}

func (self *tcpChannel) Write(msg interface{}, timeout time.Duration) (err error) {
	result := make(chan error, 1)
	defer close(result)

	self.AsyncWrite(msg, timeout, func(msg interface{}, err error) {
		result <- err
	})

	return <-result
}

func (self *tcpChannel) AsyncWrite(msg interface{}, timeout time.Duration, sendCallback SendMessageCallBack) {
	if !self.IsRunning() {
		sendCallback(msg, ErrConnectClosed)
		return
	}

	select {
	case self.sendChan <- &asyncMessage{msg: msg, timeout: time.Now().Add(timeout), callback: sendCallback}:
		//do nothing
	case <-time.After(timeout):
		sendCallback(msg, ErrWriteTimeout)
	}
}

func (self *tcpChannel) GetConnect() net.Conn {
	return self.connect
}

func (self *tcpChannel) GetRemoteAddress() string {
	return self.connect.RemoteAddr().String()
}

func (self *tcpChannel) GetRemoteIp() string {
	address := self.GetRemoteAddress()
	idx := strings.Index(address, ":")
	if idx != -1 {
		address = address[0:idx]
	}
	return address
}

func (self *tcpChannel) unsendCallbackNotify() {
	close(self.sendChan)
	for msg := range self.sendChan {
		msg.callback(msg.msg, ErrConnectClosed)
	}
}

func (self *tcpChannel) closeChannel() {
	self.closeOne.Do(func() {
		logger.Debug("关闭 channel：", self.GetRemoteAddress())
		self.status.CompareAndSet(starting, stoping)

		close(self.closeChan)

		if err := self.connect.Close(); err != nil {
			logger.Error("关闭连接异常：", self.GetRemoteAddress(), ", error:", err)
		}

		self.status.CompareAndSet(stoping, stoped)

		self.unsendCallbackNotify()
	})
}

func (self *tcpChannel) Close() {
	self.closeChannel()
	self.group.Wait()
}

func (self *tcpChannel) IsRunning() bool {
	return self.status.Get() == starting
}
