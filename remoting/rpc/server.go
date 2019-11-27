package rpc

import (
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/remoting"
	"time"
)

type RpcServer interface {
	Start()

	Shutdown()

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
	handler   OnMessage
	respCache map[uint32]*responseCache
	id        *atomic.AtomicUint32
}

func NewServer(address string, onMessage OnMessage) (RpcServer, error) {
	rpcServer := new(rpcServer)
	rpcServer.id = atomic.NewAtomicUint32(0)
	config := remoting.DefaultTCPConfig()
	if server, err := remoting.NewServer(address, config, makeHandlerMaker(onMessage, rpcServer.dealResponse), coderMaker); err != nil {
		return nil, err
	} else {
		rpcServer.server = server
	}
	rpcServer.handler = onMessage
	rpcServer.respCache = make(map[uint32]*responseCache)
	return rpcServer, nil
}

func (s *rpcServer) Start() {
	s.server.Start()
}

func (s *rpcServer) Shutdown() {
	s.server.Stop().Wait()
}

func (s *rpcServer) dealResponse(resp *Response) {
	if cache, has := s.respCache[resp.id]; has {
		delete(s.respCache, resp.id)
		commons.Try(func() {
			if cache.C != nil {
				cache.C <- resp
			} else {
				cache.callback(resp)
			}
		}, func(e error) { //防止并发问题正好删除关闭触发这里
			logger.Warn("dealResponse error:", e)
		})
	} else {
		logger.Debug("ignore response: ", resp.id)
	}
}

func (s *rpcServer) GetChannel(channel string) (remoting.Channel, bool) {
	return s.server.GetClientManager().Get(channel)
}

func (s *rpcServer) Send(channel string, request *Request, timeout time.Duration) *Response {
	request.id = s.id.IncrementAndGet(1)

	response := new(Response)
	if ch, has := s.GetChannel(channel); has {
		response.Error = ErrNotFount
	} else if err := ch.Write(request); err != nil {
		response.Error = err
	}

	rc := &responseCache{C: make(chan *Response)}
	s.respCache[request.id] = rc

	select {
	case resp := <-rc.C:
		response = resp
	case <-time.After(timeout):
		response.Error = ErrRpcTimeout
	}

	delete(s.respCache, request.id)
	rc.Close()

	return response
}

func (s *rpcServer) Async(channel string, request *Request, timeout time.Duration, callback func(response *Response)) {
	request.id = s.id.IncrementAndGet(1)

	if ch, has := s.GetChannel(channel); !has {
		callback(&Response{Error: ErrNotFount, id: request.id})
	} else {
		rc := &responseCache{callback: callback, timeout: time.Now().Add(timeout)}
		if err := ch.Write(request); err != nil {
			callback(&Response{Error: err, id: request.id})
		} else {
			s.respCache[request.id] = rc
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
