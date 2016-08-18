package gotcp

import (
	"net"
	"sync"
	"sync/atomic"
)

type Client struct {
	callback          ClientCallback
	protocol          Protocol
	channel           *net.TCPConn

	config            *Config

	extraData         interface{}

	closeFlag         int32
	closeOnce         sync.Once       // close the conn, once, per instance
	closeChan         chan struct{}   // close chanel

	packetSendChan    chan interface{}
	packetReceiveChan chan interface{}
}

// GetExtraData gets the extra data from the Conn
func (self *Client) GetExtraData() interface{} {
	return self.extraData
}

// PutExtraData puts the extra data with the Conn
func (self *Client) PutExtraData(data interface{}) {
	self.extraData = data
}

// GetRawConn returns the raw net.TCPConn from the Conn
func (self *Client) GetRawConn() *net.TCPConn {
	return self.channel
}

func (self *Client) Start(conn *net.TCPConn) {
	self.channel = conn
	asyncClientDo(self.handlerLoop)
	asyncClientDo(self.writeLoop)
	asyncClientDo(self.readLoop)
	self.callback.OnConnect(self)
}

func (self *Client) Write(p interface{}) (err error) {
	if self.IsClosed() {
		return ERR_CONN_CLOSING
	}
	defer func() {
		if e := recover(); e != nil {
			err = ERR_CONN_CLOSING
		}
	}()

	select {
	case self.packetSendChan <- p:
		return nil

	case <-self.closeChan:
		return ERR_CONN_CLOSING
	}

	return nil
}

func (self *Client) IsClosed() bool {
	return atomic.LoadInt32(&self.closeFlag) == 1
}

func (self *Client) Close() {
	self.closeOnce.Do(func() {
		atomic.StoreInt32(&self.closeFlag, 1)
		self.channel.Close()
		close(self.closeChan)
		close(self.packetReceiveChan)
		close(self.packetSendChan)
		self.callback.OnClose(self)
	})
}

func (self *Client) readLoop() {
	defer func() {
		recover()
		self.Close()
	}()

	for {
		select {
		case <-self.closeChan:
			return
		}

		readFun := func() {
			defer func() {
				if err := recover(); err != nil {
					self.callback.OnError(self, DecodePackageError{Msg:err})
				}
			}()
			p, err := self.protocol.Decode(self.channel)
			if err != nil {
				self.Close()
				return
			}
			self.packetReceiveChan <- p
		}
		readFun()
	}
}

func (self *Client) writeLoop() {
	defer func() {
		recover()
		self.Close()
	}()

	for {
		select {
		case <-self.closeChan:
			return
		case p := <-self.packetSendChan:
			if self.IsClosed() {
				return
			}
			writeFun := func() {
				defer func() {
					if err := recover(); err != nil {
						self.callback.OnError(self, EncodePackageError{Msg:err})
					}
				}()
				bytes, err := self.protocol.Encode(p)
				if err != nil {
					self.callback.OnError(self, EncodePackageError{Msg:err})
					return
				}
				if _, err := self.channel.Write(bytes); err != nil {
					self.Close()
				}
			}
			writeFun()
		}
	}
}
func (self *Client) handlerLoop() {
	defer func() {
		recover()
		self.Close()
	}()

	for {
		select {
		case <-self.closeChan:
			return

		case p := <-self.packetReceiveChan:
			if self.IsClosed() {
				return
			}
			handlerFun := func() {
				defer func() {
					if err := recover(); err != nil {
						self.callback.OnError(self, err)
					}
				}()
				self.callback.OnMessage(self, p)
			}
			//async handler message
			if self.config.AsyncMessageHand {
				go handlerFun()
			} else {
				handlerFun()
			}
		}
	}
}

func asyncClientDo(fn func()) {
	go func() {
		fn()
	}()
}


//create new client
func NewClient(config *Config, callback ClientCallback, protocol Protocol) *Client {
	return &Client{
		callback: callback,
		protocol:protocol,
		config:config,
		closeChan:make(chan struct{}),
		packetSendChan:    make(chan interface{}, config.PacketSendChanLimit),
		packetReceiveChan: make(chan interface{}, config.PacketReceiveChanLimit),
	}
}