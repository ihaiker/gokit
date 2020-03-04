package remoting

import (
	"fmt"
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/concurrent/executors"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type asyncMessage struct {
	msg      Message
	timeout  time.Time
	callback SendMessageResult
}

//发送消息回调
type SendMessageResult func(msg Message, err error)

//连接保持器
type Channel interface {

	//同步发送消息
	Write(msg interface{}, timeout time.Duration) (err error)

	//异步发送消息
	AsyncWrite(msg interface{}, timeout time.Duration, result SendMessageResult)

	//获取远程连接地址
	GetRemoteAddress() string

	//当前状态
	GetStatus() Status

	//关闭
	Close() error

	commons.Attributes
}

type tcpChannel struct {
	options *Options
	connect net.Conn

	status    Status
	closeOne  *sync.Once
	closeChan chan struct{} //关闭信号

	group    *sync.WaitGroup
	sendChan chan *asyncMessage //发送队列

	coder   Coder
	handler Handler

	commons.Attributes

	heartbeatDetection *time.Timer         //心跳检测时间
	lastBeatTime       *atomic.AtomicInt64 //上次心跳时间

	worker executors.ExecutorService
}

func newChannel(options *Options, worker executors.ExecutorService, conn net.Conn) *tcpChannel {
	if c, match := conn.(*net.TCPConn); match {
		_ = c.SetKeepAlive(true)
		_ = c.SetNoDelay(true)
		_ = c.SetReadBuffer(options.SendBuf)
		_ = c.SetWriteBuffer(options.RecvBuf)
	}

	return &tcpChannel{
		options: options, connect: conn, worker: worker,

		closeOne: new(sync.Once), status: Ready,
		group: new(sync.WaitGroup),

		closeChan: make(chan struct{}),
		sendChan:  make(chan *asyncMessage, options.SendChanLimit),

		Attributes: commons.NewAttributes(),
	}
}

func (self *tcpChannel) syncDo(fn func()) {
	defer self.group.Done()
	fn()
}

func (self *tcpChannel) onEvent(event *Event) {
	defer func() { _ = recover() }()
	self.handler.OnEvent(event)
}

//connected 连接后回调
func (self *tcpChannel) do(connected func(channel Channel)) {
	defer func() {
		self.closeChannel()
		self.notSendCallback()
		self.onEvent(NewEvent(CloseEvent, self))
	}()
	self.group.Add(3)

	go self.syncDo(self.readLoop)
	go self.syncDo(self.writeLoop)
	go self.syncDo(self.heartbeatLoop)

	self.status = Running
	self.onEvent(NewEvent(ConnectEvent, self))
	if connected != nil {
		connected(self)
	}
	self.group.Wait()
}

func (self *tcpChannel) readLoop() {
	defer self.closeChannel()

	for {
		select {
		case <-self.closeChan:
			return
		default:
			if msg, err := self.coder.Decode(self, self.connect); commons.NotNil(err) {
				if isCloseTCPConnect(err) { //连接已经关闭
					return
				}
				if self.GetStatus().IsStart() {
					self.onEvent(NewEvent(DecodeErrEvent, self, err))
				}
			} else {
				self.resetIdle()
				handlerMessage := func() {
					defer func() {
						if err := recover(); err != nil {
							self.onEvent(NewEvent(ErrEvent, self, fmt.Errorf("%v", err)))
						}
					}()
					self.onEvent(NewEvent(MessageEvent, self, msg))
				}
				if self.worker != nil { //异步执行
					_ = self.worker.Submit(handlerMessage)
				} else {
					handlerMessage()
				}
			}
		}
	}
}

func (self *tcpChannel) writeLoop() {
	defer self.closeChannel()

	for {
		select {
		case <-self.closeChan:
			return
		case asyncMsg := <-self.sendChan:
			if asyncMsg == nil {
				return //sendChan is close, asyncMsg is nil
			}
			if !self.GetStatus().IsStart() {
				if asyncMsg.callback != nil {
					asyncMsg.callback(asyncMsg.msg, ErrConnectClosed)
				}
				return
			}
			if time.Now().Before(asyncMsg.timeout) {
				if bs, err := self.coder.Encode(self, asyncMsg.msg); err != nil {
					self.onEvent(NewEvent(EncodeErrEvent, self, asyncMsg.msg, err))
				} else {
					_, err := self.connect.Write(bs)
					if asyncMsg.callback != nil {
						asyncMsg.callback(asyncMsg.msg, err)
					}
					if err != nil {
						self.onEvent(NewEvent(ErrEvent, self, asyncMsg.msg, err))
					}
				}
			} else if asyncMsg.callback != nil {
				asyncMsg.callback(asyncMsg.msg, ErrWriteTimeout)
			}
		}
	}
}

func (self *tcpChannel) heartbeatLoop() {
	if self.options.IdleTimeSeconds == 0 {
		return
	}
	self.heartbeatDetection = time.NewTimer(time.Second * time.Duration(self.options.IdleTimeout))
	self.lastBeatTime = atomic.NewAtomicInt64(time.Now().Unix())

	timeout := int64(self.options.IdleTimeout * self.options.IdleTimeSeconds)

	defer func() {
		self.heartbeatDetection.Stop()
		self.closeChannel()
	}()

	for {
		select {
		case <-self.closeChan:
			return
		case <-self.heartbeatDetection.C:
			self.heartbeatDetection.Reset(time.Second * time.Duration(self.options.IdleTimeSeconds))
			p := time.Now().Unix() - self.lastBeatTime.Get()
			if p >= timeout { //timeout
				return
			} else if p < int64(self.options.IdleTimeSeconds) {
				//没到时间呢
			} else {
				self.onEvent(NewEvent(IdleEvent, self))
			}
		}
	}
}

func (self *tcpChannel) resetIdle() {
	if self.heartbeatDetection != nil {
		self.lastBeatTime.Set(time.Now().Unix())
	}
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

	self.AsyncWrite(msg, timeout, func(msg interface{}, writerErr error) {
		result <- writerErr
	})

	err = <-result
	return
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

func (self *tcpChannel) notSendCallback() {
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
		self.status = Stop
		close(self.closeChan)
		_ = self.connect.Close()
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

func (self *tcpChannel) Wait() {
	self.group.Wait()
}
