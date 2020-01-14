package handler

import (
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
)

type HandleAdapter struct{}

func (self *HandleAdapter) OnEvent(event *remoting.Event) {
	logs.GetLogger("remoting").Debug(event.Channel.GetRemoteAddress(), " event: ", event.Type, ", message: ", event.Values)
}

func Adapter() *HandleAdapter {
	return &HandleAdapter{}
}
