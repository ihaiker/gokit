package handler

import (
	"github.com/ihaiker/gokit/remoting"
)

type reg struct {
	def    remoting.Handler
	events map[remoting.EventType]func(event *remoting.Event)
}

func Reg() *reg {
	return &reg{
		events: make(map[remoting.EventType]func(event *remoting.Event)),
	}
}

func (self *reg) With(def remoting.Handler) *reg {
	self.def = def
	return self
}

func (self *reg) On(event remoting.EventType, handler func(event *remoting.Event)) *reg {
	self.events[event] = handler
	return self
}

func (self *reg) Ons(events ...remoting.EventType) func(func(*remoting.Event)) *reg {
	return func(handler func(*remoting.Event)) *reg {
		for _, event := range events {
			self.events[event] = handler
		}
		return self
	}
}

func (self *reg) OnEvent(event *remoting.Event) {
	if handle, has := self.events[event.Type]; has {
		handle(event)
	} else if self.def != nil {
		self.def.OnEvent(event)
	}
}
