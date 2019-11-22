package handler

import (
	"github.com/ihaiker/gokit/remoting"
)

type reg struct {
	def           remoting.Handler
	onConnect     func(session remoting.Channel)
	onMessage     func(session remoting.Channel, msg interface{})
	onEncodeError func(session remoting.Channel, msg interface{}, err error)
	onError       func(session remoting.Channel, msg interface{}, err error)
	onDecodeError func(session remoting.Channel, err error)
	onIdle        func(session remoting.Channel)
	onClose       func(session remoting.Channel)
}

func Reg() *reg {
	return &reg{}
}

func (self *reg) With(def remoting.Handler) *reg {
	self.def = def
	return self
}

func (self *reg) WithOnConnect(fn func(session remoting.Channel)) *reg {
	self.onConnect = fn
	return self
}

func (self *reg) WithOnMessage(fn func(session remoting.Channel, msg interface{})) *reg {
	self.onMessage = fn
	return self
}

func (self *reg) WithOnEncodeError(fn func(session remoting.Channel, msg interface{}, err error)) *reg {
	self.onEncodeError = fn
	return self
}

func (self *reg) WithOnError(fn func(session remoting.Channel, msg interface{}, err error)) *reg {
	self.onError = fn
	return self
}

func (self *reg) WithOnDecodeError(fn func(session remoting.Channel, err error)) *reg {
	self.onDecodeError = fn
	return self
}

func (self *reg) WithOnIdle(fn func(session remoting.Channel)) *reg {
	self.onIdle = fn
	return self
}

func (self *reg) WithOnClose(fn func(session remoting.Channel)) *reg {
	self.onClose = fn
	return self
}

func (self *reg) OnConnect(session remoting.Channel) {
	if self.onConnect == nil {
		if self.def != nil {
			self.def.OnConnect(session)
		}
	} else {
		self.onConnect(session)
	}
}
func (self *reg) OnMessage(session remoting.Channel, msg interface{}) {
	if self.onMessage == nil {
		if self.def != nil {
			self.def.OnMessage(session, msg)
		}
	} else {
		self.onMessage(session, msg)
	}
}
func (self *reg) OnClose(session remoting.Channel) {
	if self.onClose == nil {
		if self.def != nil {
			self.def.OnClose(session)
		}
	} else {
		self.onClose(session)
	}
}
func (self *reg) OnError(session remoting.Channel, msg interface{}, err error) {
	if self.onError == nil {
		if self.def != nil {
			self.def.OnError(session, msg, err)
		}
	} else {
		self.onError(session, msg, err)
	}
}
func (self *reg) OnEncodeError(session remoting.Channel, msg interface{}, err error) {
	if self.onEncodeError == nil {
		if self.def != nil {
			self.def.OnEncodeError(session, msg, err)
		}
	} else {
		self.onEncodeError(session, msg, err)
	}
}
func (self *reg) OnDecodeError(session remoting.Channel, err error) {
	if self.onDecodeError == nil {
		if self.def != nil {
			self.def.OnDecodeError(session, err)
		}
	} else {
		self.onDecodeError(session, err)
	}
}
func (self *reg) OnIdle(session remoting.Channel) {
	if self.onIdle == nil {
		if self.def != nil {
			self.def.OnIdle(session)
		}
	} else {
		self.onIdle(session)
	}
}
