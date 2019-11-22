package handler

import (
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
)

var logger = logs.GetLogger("tcp")

type HandlerAdapter struct{}

func (h *HandlerAdapter) OnConnect(session remoting.Channel) {
	logger.Debug("OnConnect:", session.GetRemoteAddress())
}

func (h *HandlerAdapter) OnMessage(session remoting.Channel, msg interface{}) {
	logger.Debugf("OnMessage %s : msg:%v", session.GetRemoteAddress(), msg)
}

func (h *HandlerAdapter) OnClose(session remoting.Channel) {
	logger.Debugf("OnClose %s ", session.GetRemoteAddress())
}

func (h *HandlerAdapter) OnError(session remoting.Channel, msg interface{}, err error) {
	logger.Debugf("OnError %s : %s ,%v", session.GetRemoteAddress(), err, msg)
}

func (h *HandlerAdapter) OnEncodeError(session remoting.Channel, msg interface{}, err error) {
	h.OnError(session, msg, err)
}

func (h *HandlerAdapter) OnDecodeError(session remoting.Channel, err error) {
	h.OnError(session, nil, err)
}

func (h *HandlerAdapter) OnIdle(session remoting.Channel) {
	logger.Debug("OnIdle : ", session.GetRemoteAddress())
}

func Adapter() *HandlerAdapter {
	return &HandlerAdapter{}
}
