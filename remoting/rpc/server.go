package rpc

import (
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/remoting"
	"time"
)

type RpcServer interface {
	Start() error

	Wait()

	Close() error

	GetChannelManager() remoting.ChannelManager
	SetChannelManager(remoting.ChannelManager)

	GetChannel(channel string) (ch remoting.Channel, has bool)

	Send(channel string, request *Request, timeout time.Duration) (response *Response)

	Async(channel string, request *Request, timeout time.Duration, callback func(response *Response))

	Oneway(address string, request *Request, timeout time.Duration)
}

type responseCache struct {
	C        chan *Response
	callback func(response *Response)
	timeout  time.Time
}

func (rc *responseCache) Close() {
	close(rc.C)
}

type rpcServer struct {
	server    remoting.Server
	respCache map[uint32]*responseCache
	id        *atomic.AtomicUint32
}

func NewServer(address string, onMessage OnMessage, onClose OnClose) RpcServer {
	return NewServerWithConfig(address, remoting.DefaultTCPConfig(), onMessage, onClose)
}

func NewServerWithConfig(address string, config *remoting.Config, onMessage OnMessage, onClose OnClose) RpcServer {
	rpcServer := new(rpcServer)
	rpcServer.id = atomic.NewAtomicUint32(0)
	rpcServer.server = remoting.NewServer(address, config, makeHandlerMaker(onMessage, rpcServer.onResponse, onClose), coderMaker)
	rpcServer.respCache = make(map[uint32]*responseCache)
	return rpcServer
}

func (s *rpcServer) Start() error {
	return s.server.Start()
}

func (s *rpcServer) Wait() {
	s.server.Wait()
}

func (s *rpcServer) Close() error {
	return s.server.Stop()
}

func (s *rpcServer) onResponse(resp *Response) {
	if cache, has := s.respCache[resp.id]; has {
		delete(s.respCache, resp.id)
		commons.Try(func() {
			if cache.C != nil {
				cache.C <- resp
			} else {
				cache.callback(resp)
			}
		}, func(e error) { //防止并发问题正好删除关闭触发这里
			logger.Warn("onResponse error:", e)
		})
	} else {
		logger.Debug("ignore response: ", resp.id)
	}
}

func (s *rpcServer) GetChannel(channel string) (remoting.Channel, bool) {
	return s.server.GetClientManager().Get(channel)
}

func (s *rpcServer) GetChannelManager() remoting.ChannelManager {
	return s.server.GetClientManager()
}

func (s *rpcServer) SetChannelManager(manager remoting.ChannelManager) {
	s.server.SetClientManager(manager)
}

func (s *rpcServer) Send(channel string, request *Request, timeout time.Duration) *Response {
	request.id = s.id.IncrementAndGet(1)

	response := new(Response)
	if ch, has := s.GetChannel(channel); !has {
		response.Error = ErrNotFount
	} else if err := ch.Write(request, timeout); err != nil {
		response.Error = err
	}
	if response.Error != nil {
		return response
	}

	rc := &responseCache{C: make(chan *Response)}
	defer rc.Close()

	s.respCache[request.id] = rc
	defer delete(s.respCache, request.id)

	select {
	case resp := <-rc.C:
		response = resp
	case <-time.After(timeout):
		response.Error = ErrRpcTimeout
	}

	return response
}

func (s *rpcServer) Async(channel string, request *Request, timeout time.Duration, callback func(response *Response)) {
	request.id = s.id.IncrementAndGet(1)

	if ch, has := s.GetChannel(channel); !has {
		callback(&Response{Error: ErrNotFount, id: request.id})
	} else {
		if err := ch.Write(request, timeout); err != nil {
			callback(&Response{Error: err, id: request.id})
		} else {
			s.respCache[request.id] = &responseCache{callback: callback, timeout: time.Now().Add(timeout)}
		}
	}
}

func (s *rpcServer) Oneway(address string, request *Request, timeout time.Duration) {
	s.Async(address, request, timeout, func(response *Response) {
		if response.Error != nil {
			logger.Error("send oneway error:", response.Error.Error())
		}
	})
}
