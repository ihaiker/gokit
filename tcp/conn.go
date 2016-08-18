package gotcp

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Conn exposes a set of callbacks for the various events that occur on a connection
type Conn struct {
	srv               *Server
	conn              *net.TCPConn     // the raw connection
	extraData         interface{}      // to save extra data
	closeOnce         sync.Once        // close the conn, once, per instance
	closeFlag         int32            // close flag
	closeChan         chan struct{}    // close chanel
	packetSendChan    chan interface{} // packet send chanel
	packetReceiveChan chan interface{} // packet receive chanel

	heartbeat         *Heartbeat
	idleTimer         []*time.Timer
}

// newConn returns a wrapper of raw conn
func newConn(conn *net.TCPConn, srv *Server, heartbeat *Heartbeat) *Conn {
	return &Conn{
		srv:               srv,
		conn:              conn,
		closeChan:         make(chan struct{}),
		packetSendChan:    make(chan interface{}, srv.config.PacketSendChanLimit),
		packetReceiveChan: make(chan interface{}, srv.config.PacketReceiveChanLimit),
		heartbeat:  heartbeat,
	}
}

// GetExtraData gets the extra data from the Conn
func (c *Conn) GetExtraData() interface{} {
	return c.extraData
}

// PutExtraData puts the extra data with the Conn
func (c *Conn) PutExtraData(data interface{}) {
	c.extraData = data
}

// GetRawConn returns the raw net.TCPConn from the Conn
func (c *Conn) GetRawConn() *net.TCPConn {
	return c.conn
}

// Close closes the connection
func (c *Conn) Close() {
	c.closeOnce.Do(func() {
		atomic.StoreInt32(&c.closeFlag, 1)
		close(c.closeChan)
		close(c.packetSendChan)
		close(c.packetReceiveChan)
		c.conn.Close()
		c.srv.callback.OnClose(c)
	})
}

// IsClosed indicates whether or not the connection is closed
func (c *Conn) IsClosed() bool {
	return atomic.LoadInt32(&c.closeFlag) == 1
}

// AsyncWritePacket async writes a packet, this method will never block
func (c *Conn) Write(p interface{}, timeout time.Duration) (err error) {
	if c.IsClosed() {
		return ERR_CONN_CLOSING
	}

	defer func() {
		if e := recover(); e != nil {
			err = ERR_CONN_CLOSING
		}
	}()

	if timeout == 0 {
		select {
		case c.packetSendChan <- p:
			return nil

		//my add
		case <-c.closeChan:
			return ERR_CONN_CLOSING

		default:
			return ERR_CONN_CLOSING
		}

	} else {
		select {
		case c.packetSendChan <- p:
			return nil

		case <-c.closeChan:
			return ERR_CONN_CLOSING

		case <-time.After(timeout):
			return ERR_WRITE_BLOCKING
		}
	}
}

func asyncDo(fn func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()
}

// Do it
func (c *Conn) Do() {
	c.srv.callback.OnConnect(c)
	if c.IsClosed() {
		return
	}
	asyncDo(c.handleLoop, c.srv.waitGroup)
	asyncDo(c.readLoop, c.srv.waitGroup)
	asyncDo(c.writeLoop, c.srv.waitGroup)
	if c.heartbeat != nil {
		asyncDo(c.heartbeatLoop, c.srv.waitGroup)
	}
}

func (c *Conn) handleLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return

		case p := <-c.packetReceiveChan:
			if c.IsClosed() {
				return
			}
			fun := func() {
				defer func() {
					if err := recover(); err != nil {
						c.srv.callback.OnError(c, err)
					}
				}()
				c.srv.callback.OnMessage(c, p)
			}
		//async handler message
			if c.srv.config.AsyncMessageHand {
				go fun()
			} else {
				fun()
			}
		}
	}
}

func (c *Conn) readLoop() {

	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return

		default:
			readFun := func() {
				defer func() {
					if err := recover(); err != nil {
						c.srv.callback.OnError(c, DecodePackageError{Msg:err})
					}
				}()
				p, err := c.srv.protocol.Decode(c.conn)
				if err != nil {
					c.Close()
					return
				}
				c.readPackage()
				c.packetReceiveChan <- p
			}
			readFun()
		}
	}
}

func (c *Conn) writeLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return

		case p := <-c.packetSendChan:
			if c.IsClosed() {
				return
			}
			writeFun := func() {
				defer func() {
					if err := recover(); err != nil {
						c.srv.callback.OnError(c, EncodePackageError{Msg:err})
					}
				}()
				bytes, err := c.srv.protocol.Encode(p)
				if err != nil {
					c.srv.callback.OnError(c, EncodePackageError{Msg:err})
					return
				}
				if isIdle(p) == false {
					c.writePackage()
				}
				if _, err := c.conn.Write(bytes); err != nil {
					c.Close()
				}
			}
			writeFun()
		}
	}
}

func (self *Conn) heartbeatLoop() {
	self.idleTimer = make([]*time.Timer, 3)
	self.idleTimer[0] = time.NewTimer(self.heartbeat.Idle.ReadIdle)
	self.idleTimer[1] = time.NewTimer(self.heartbeat.Idle.WriteIdle)
	self.idleTimer[2] = time.NewTimer(self.heartbeat.Idle.AllIdle)

	defer func() {
		recover()
		for _, it := range self.idleTimer {
			it.Stop()
		}
		self.Close()
	}()

	for {
		select {
		case <-self.srv.exitChan: //服务器关闭
			return
		case <-self.closeChan:   //客户端被关闭
			return

		case <-self.idleTimer[0].C:
			self.idleTimer[0].Reset(self.heartbeat.Idle.ReadIdle)
			self.heartbeat.Handler.OnIdle(READ_IDLE_STATE, self)
		case <-self.idleTimer[1].C:
			self.idleTimer[1].Reset(self.heartbeat.Idle.WriteIdle)
			self.heartbeat.Handler.OnIdle(WRITER_IDLE_STATE, self)
		case <-self.idleTimer[2].C:
			self.idleTimer[2].Reset(self.heartbeat.Idle.AllIdle)
			self.heartbeat.Handler.OnIdle(ALL_IDLE_STATE, self)
		}
	}
}

func (self *Conn) readPackage()  {
	if self.heartbeat != nil {
		self.idleTimer[0].Reset(self.heartbeat.Idle.ReadIdle)
		self.idleTimer[2].Reset(self.heartbeat.Idle.AllIdle)
	}
}
func (self *Conn) writePackage()  {
	if self.heartbeat != nil {
		self.idleTimer[1].Reset(self.heartbeat.Idle.WriteIdle)
		self.idleTimer[2].Reset(self.heartbeat.Idle.AllIdle)
	}
}

func isIdle(i interface{}) bool {
	switch i.(type) {
	case IdlePackage:
		return true
	}
	return false
}