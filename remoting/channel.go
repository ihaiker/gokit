package remoting

import (
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/concurrent/executors"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

//发送消息回调
type SendMessageResult func(msg interface{}, err error)

type asyncMessage struct {
	msg      interface{}
	timeout  time.Time
	callback SendMessageResult
}

//连接保持器
type Channel interface {

	//同步发送消息
	Write(msg interface{}, timeout time.Duration) (err error)

	//异步发送消息
	AsyncWrite(msg interface{}, timeout time.Duration, result SendMessageResult)

	//获取远程连接地址
	GetRemoteAddress() string

	//获取远程的IP地址
	GetRemoteIp() string

	GetStatus() Status

	//关闭
	Close() error

	Wait()

	commons.Attributes
}

type tcpChannel struct {
	config *Config

	connect net.Conn

	closeOne *sync.Once

	status Status

	group *sync.WaitGroup

	closeChan chan struct{}
	sendChan  chan *asyncMessage

	coder   Coder
	handler Handler

	commons.Attributes

	idleTimer   *time.Timer
	idleTimeout *atomic.AtomicInt32

	worker *executors.GrPool
}

func newChannel(config *Config, worker *executors.GrPool, connect net.Conn) *tcpChannel {

	if tcpCon, match := connect.(*net.TCPConn); match {
		_ = tcpCon.SetKeepAlive(true)
		_ = tcpCon.SetNoDelay(true)
		_ = tcpCon.SetReadBuffer(config.ReadBufferSize)
		_ = tcpCon.SetWriteBuffer(config.WriteBufferSize)
	}

	return &tcpChannel{
		config: config, connect: connect, worker: worker,

		closeOne: new(sync.Once), status: Ready,
		group: new(sync.WaitGroup),

		closeChan: make(chan struct{}),
		sendChan:  make(chan *asyncMessage, config.SendChanLimit),

		Attributes: commons.NewAttributes(),
	}
}

//安全执行外部方法
func (self *tcpChannel) safeNotify(fn func(channel Channel)) {
	if fn == nil {
		return
	}
	_ = commons.AsyncTimeout(time.Second, func() interface{} {
		fn(self)
		return nil
	})
}

//安全执行外部方法
func (self *tcpChannel) safeNotifyError(fn func(Channel, error), err error) {
	if fn == nil {
		return
	}
	_ = commons.AsyncTimeout(time.Second, func() interface{} {
		return commons.SafeExec(func() {
			fn(self, err)
		})
	})
}

func (self *tcpChannel) syncDo(fn func()) {
	defer self.group.Done()
	fn()
}

//connected 连接后回调
//closed 关闭后回调
func (self *tcpChannel) do(connected, closed func(channel Channel)) {
	defer func() {
		logger.Debug("channel over: ", self)
		self.closeChannel()
		self.safeNotify(self.handler.OnClose)
		self.safeNotify(closed)
	}()

	logger.Debug("channel start: ", self)

	self.group.Add(3)

	go self.syncDo(self.heartbeatLoop)
	go self.syncDo(self.readLoop)
	go self.syncDo(self.writeLoop)

	self.safeNotify(self.handler.OnConnect)
	self.safeNotify(connected)

	self.status = Running

	self.Wait()
}

func (self *tcpChannel) readLoop() {
	defer func() {
		self.closeChannel()
		logger.Debug("channel reader close: ", self)
	}()
	logger.Debug("channel reader start: ", self)

	for {
		select {
		case <-self.closeChan:
			return
		default:
			if msg, err := self.coder.Decode(self, self.connect); commons.NotNil(err) {
				if isCloseTCPConnect(err) { //连接已经关闭
					return
				}
				if !strings.Contains(err.Error(), "i/o timeout") {
					if self.GetStatus().IsStart() {
						logger.Error("channel decode error: ", err)
						self.safeNotifyError(self.handler.OnDecodeError, err)
					}
				}
			} else {
				self.resetIdle()
				handlerMessage := func() {
					commons.Try(func() {
						self.handler.OnMessage(self, msg)
					}, func(err error) {
						defer func() { _ = recover() }()
						self.handler.OnError(self, msg, err)
					})
				}
				if self.worker != nil { //异步执行
					self.worker.Add(handlerMessage)
				} else {
					handlerMessage()
				}
			}
		}
	}
}

func (self *tcpChannel) writeLoop() {
	defer func() {
		logger.Debug("channel write stop: ", self)
		self.closeChannel()
	}()
	logger.Debug("channel write start: ", self)

	for {
		select {
		case <-self.closeChan:
			return
		case asyncMsg := <-self.sendChan:
			if !self.GetStatus().IsStart() {
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
		return
	}
	logger.Debug("channel ttl start:", self)
	self.idleTimer = time.NewTimer(time.Second * time.Duration(self.config.IdleDuration))
	self.idleTimeout = atomic.NewAtomicInt32(0)

	defer func() {
		logger.Debug("channel ttl stop:", self)
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
				logger.Debug("channel ttl timeout：", self)
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

func (self *tcpChannel) Wait() {
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

func (self *tcpChannel) AsyncWrite(msg interface{}, timeout time.Duration, sendCallback SendMessageResult) {
	if !self.GetStatus().IsStart() {
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

func (self *tcpChannel) unsendNotify() {
	close(self.sendChan)
	for msg := range self.sendChan {
		msg.callback(msg.msg, ErrConnectClosed)
	}
}

func (self *tcpChannel) GetStatus() Status {
	return self.status
}

//关闭服务
func (self *tcpChannel) closeChannel() {
	self.closeOne.Do(func() {
		logger.Debug("close channel：", self)
		_ = self.connect.Close()
		self.status = Stoped
		close(self.closeChan)
		self.unsendNotify()
	})
}

//关闭服务并等待退出
func (self *tcpChannel) Close() error {
	self.closeChannel()
	return nil
}

func (self *tcpChannel) String() string {
	return self.GetRemoteAddress()
}
