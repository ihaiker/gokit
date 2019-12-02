package rpc

import (
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/coder/tlv"
	"github.com/ihaiker/gokit/remoting/handler"
	"time"
)

type OnMessage func(channel remoting.Channel, request *Request) *Response
type OnResponse func(response *Response)

func OK(channel remoting.Channel, request *Request) *Response {
	resp := NewResponse(request.ID())
	resp.Body = []byte("OK")
	return resp
}

func Error(channel remoting.Channel, request *Request, err error) *Response {
	resp := NewResponse(request.ID())
	resp.Error = err
	return resp
}

func Check(check func(channel remoting.Channel, request *Request) error, onMessage OnMessage) OnMessage {
	return func(channel remoting.Channel, request *Request) *Response {
		if err := check(channel, request); err != nil {
			resp := NewResponse(request.ID())
			resp.Error = err
			return resp
		} else {
			return onMessage(channel, request)
		}
	}
}

func makeHandlerMaker(onMessage OnMessage, onResponse OnResponse) remoting.HandlerMaker {
	return func(channel remoting.Channel) remoting.Handler {
		return newHandler(onMessage, onResponse)
	}
}

func newHandler(onMessage OnMessage, onResponse OnResponse) remoting.Handler {
	ping := new(Ping)
	pong := new(Pong)
	reg := handler.Reg()
	reg.WithOnIdle(func(ch remoting.Channel) {
		logger.Debug("write ping to:", ch.GetRemoteAddress())
		_ = ch.Write(ping, time.Second)
	}).WithOnDecodeError(func(ch remoting.Channel, err error) {
		logger.Debug("decoder on:", ch.GetRemoteAddress(), ", error:", err)
		ch.Close()
	}).WithOnEncodeError(func(ch remoting.Channel, msg interface{}, err error) {
		logger.Debug("encode on:", ch.GetRemoteAddress(), ", error:", err)
		ch.Close()
	}).WithOnError(func(ch remoting.Channel, msg interface{}, err error) {
		logger.Debug("error on:", ch.GetRemoteAddress(), ", error:", err)
		ch.Close()
	}).WithOnMessage(func(ch remoting.Channel, msg interface{}) {
		pkg := msg.(tlv.Message)
		switch pkg.TypeID() {
		case PING:
			_ = ch.Write(pong, time.Second)
		case PONG:
			//do nothing
		case REQUEST:
			req := msg.(*Request)
			logger.Debug("request: ", req.URL, ", ch:", ch.GetRemoteAddress())
			commons.Try(func() {
				if resp := onMessage(ch, req); resp != nil {
					if err := ch.Write(resp, time.Second); err != nil {
						logger.Errorf("write response %s error: %s", req.URL, err)
					}
				}
			}, func(e error) {
				logger.Errorf("dealwith request(%s) error: %s", req.URL, e)
				resp := new(Response)
				resp.id = msg.(*Request).id
				resp.Error = ErrSystemError
				if err := ch.Write(resp, time.Second); err != nil {
					logger.Errorf("write response %s error: %s", req.URL, err)
				}
			})
		case RESPONSE:
			logger.Debug("response: ", msg.(*Response).id, ", ch:", ch.GetRemoteAddress())
			onResponse(msg.(*Response))
		}
	})
	return reg
}
